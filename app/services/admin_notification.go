package services

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"

	"github.com/redis/go-redis/v9"
)

// BroadcastRateLimitConfig 廣播速率限制配置
type BroadcastRateLimitConfig struct {
	Enabled          bool          // 是否啟用
	MaxBroadcasts    int           // 最大廣播次數
	WindowSize       time.Duration // 時間窗口
}

// AdminNotificationService 管理員通知服務
type AdminNotificationService struct {
	BaseService
	app              *app.App
	teacherRepo      *repositories.TeacherRepository
	membershipRepo   *repositories.CenterMembershipRepository
	centerRepo       *repositories.CenterRepository
	lineBotService   LineBotService
	templateService   LineBotTemplateService
	rateLimiter      *BroadcastRateLimiter
}

// BroadcastRateLimiter 廣播專用速率限制器
type BroadcastRateLimiter struct {
	app    *app.App
	config BroadcastRateLimitConfig
}

// NewBroadcastRateLimiter 建立廣播速率限制器
func NewBroadcastRateLimiter(app *app.App) *BroadcastRateLimiter {
	config := BroadcastRateLimitConfig{
		Enabled:       true,
		MaxBroadcasts: 5, // 每分鐘最多 5 次廣播
		WindowSize:    time.Minute,
	}

	// 從環境變數讀取配置
	if app != nil && app.Env != nil {
		config.Enabled = app.Env.BroadcastRateLimitEnabled
		if app.Env.BroadcastMaxPerMinute > 0 {
			config.MaxBroadcasts = app.Env.BroadcastMaxPerMinute
		}
	}

	return &BroadcastRateLimiter{
		app:    app,
		config: config,
	}
}

// Check 檢查是否允許廣播
// 返回: allowed(是否允許), remaining(剩餘次數), resetAt(重置時間), error
func (r *BroadcastRateLimiter) Check(ctx context.Context, adminID uint, centerID uint) (bool, int, time.Time, error) {
	if !r.config.Enabled {
		return true, r.config.MaxBroadcasts, time.Now().Add(r.config.WindowSize), nil
	}

	// 使用 Redis 滑動窗口
	key := fmt.Sprintf("broadcast:ratelimit:%d:%d", centerID, adminID)
	now := time.Now()
	windowStart := now.Add(-r.config.WindowSize)

	// 移除窗口外的舊記錄
	r.app.Redis.DB0.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixMilli(), 10))

	// 計算當前請求數
	count, err := r.app.Redis.DB0.ZCard(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, 0, time.Time{}, fmt.Errorf("取得廣播計數失敗: %w", err)
	}

	remaining := r.config.MaxBroadcasts - int(count) - 1
	resetAt := now.Add(r.config.WindowSize)

	if count >= int64(r.config.MaxBroadcasts) {
		// 記錄警告日誌（使用 app.Logger 如果可用）
		fmt.Printf("[WARN] broadcast rate limit exceeded: admin_id=%d, center_id=%d, count=%d, max=%d\n",
			adminID, centerID, count, r.config.MaxBroadcasts)
		return false, remaining, resetAt, nil
	}

	return true, remaining, resetAt, nil
}

// Record 記錄廣播請求
func (r *BroadcastRateLimiter) Record(ctx context.Context, adminID uint, centerID uint) error {
	if !r.config.Enabled {
		return nil
	}

	key := fmt.Sprintf("broadcast:ratelimit:%d:%d", centerID, adminID)
	now := time.Now()

	member := redis.Z{
		Score:  float64(now.UnixNano()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	}

	pipe := r.app.Redis.DB0.Pipeline()
	pipe.ZAdd(ctx, key, member)
	pipe.Expire(ctx, key, r.config.WindowSize)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("記錄廣播失敗: %w", err)
	}

	return nil
}

// NewAdminNotificationService 建立管理員通知服務
func NewAdminNotificationService(app *app.App) *AdminNotificationService {
	templateService := NewLineBotTemplateService(app.Env.FrontendBaseURL)

	return &AdminNotificationService{
		BaseService:      *NewBaseService(app, "AdminNotificationService"),
		app:              app,
		teacherRepo:      repositories.NewTeacherRepository(app),
		membershipRepo:   repositories.NewCenterMembershipRepository(app),
		centerRepo:       repositories.NewCenterRepository(app),
		lineBotService:   NewLineBotService(app),
		templateService:  templateService,
		rateLimiter:      NewBroadcastRateLimiter(app),
	}
}

// BroadcastResult 廣播結果
type BroadcastResult struct {
	SuccessCount int              `json:"success_count"`
	FailedCount  int              `json:"failed_count"`
	TotalCount   int              `json:"total_count"`
	Message      string           `json:"message"`
	RateLimit    *RateLimitInfo   `json:"rate_limit,omitempty"`
}

// RateLimitInfo 速率限制資訊
type RateLimitInfo struct {
	Allowed     bool      `json:"allowed"`
	Remaining   int       `json:"remaining"`
	ResetAt     time.Time `json:"reset_at"`
	MaxRequests int       `json:"max_requests"`
	WindowSec   int       `json:"window_seconds"`
}

// BroadcastToTeachers 廣播訊息給中心老師
// centerID: 中心 ID（從 JWT 取得，確保資料隔離）
// adminID: 管理員 ID（用於記錄）
// messageType: 訊息類型（GENERAL 或 URGENT）
// title: 標題
// message: 訊息內容
// warning: 警告提示（可選）
// actionLabel: 按鈕文字（可選）
// actionURL: 按鈕連結（可選）
// teacherIDs: 指定老師 ID 清單（空白表示發送給所有老師）
func (s *AdminNotificationService) BroadcastToTeachers(
	ctx context.Context,
	centerID uint,
	adminID uint,
	messageType string,
	title string,
	message string,
	warning string,
	actionLabel string,
	actionURL string,
	teacherIDs []uint,
) (*BroadcastResult, *errInfos.Res, error) {
	// 【速率限制檢查】防止管理員連點
	allowed, remaining, resetAt, err := s.rateLimiter.Check(ctx, adminID, centerID)
	if err != nil {
		s.Logger.Error("rate limit check failed", "error", err, "admin_id", adminID, "center_id", centerID)
		// Redis 錯誤時記錄日誌但允許請求（fail open）
	}

	rateLimitInfo := &RateLimitInfo{
		Allowed:     allowed,
		Remaining:   remaining,
		ResetAt:     resetAt,
		MaxRequests: s.rateLimiter.config.MaxBroadcasts,
		WindowSec:   int(s.rateLimiter.config.WindowSize.Seconds()),
	}

	if !allowed {
		s.Logger.Warn("broadcast rate limited",
			"admin_id", adminID,
			"center_id", centerID,
			"remaining", remaining,
			"reset_at", resetAt,
		)
		return &BroadcastResult{
			SuccessCount: 0,
			FailedCount:  0,
			TotalCount:   0,
			Message:      fmt.Sprintf("廣播頻率過高，請在 %d 秒後再試", int(time.Until(resetAt).Seconds())),
			RateLimit:    rateLimitInfo,
		}, s.App.Err.New(errInfos.RATE_LIMIT_EXCEEDED), nil
	}

	s.Logger.Info("broadcasting message to teachers",
		"center_id", centerID,
		"admin_id", adminID,
		"specified_teachers", len(teacherIDs),
		"rate_limit_remaining", remaining,
	)

	// 取得目標老師清單
	teachers, err := s.getTargetTeachers(ctx, centerID, teacherIDs)
	if err != nil {
		s.Logger.Error("failed to get target teachers", "error", err)
		return nil, s.App.Err.New(errInfos.SQL_ERROR), err
	}

	if len(teachers) == 0 {
		return &BroadcastResult{
			SuccessCount: 0,
			FailedCount:  0,
			TotalCount:   0,
			Message:      "沒有符合條件的老師",
			RateLimit:    rateLimitInfo,
		}, nil, nil
	}

	// 過濾出有綁定 LINE 的老師
	var lineUserIDs []string
	for _, teacher := range teachers {
		if teacher.LineUserID != "" {
			lineUserIDs = append(lineUserIDs, teacher.LineUserID)
		}
	}

	if len(lineUserIDs) == 0 {
		return &BroadcastResult{
			SuccessCount: 0,
			FailedCount:  0,
			TotalCount:   len(teachers),
			Message:      "目標老師都尚未綁定 LINE",
			RateLimit:    rateLimitInfo,
		}, nil, nil
	}

	// 取得中心名稱
	centerName := ""
	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err == nil {
		centerName = center.Name
	}

	// 建立 Flex Message 結構
	flexMessage := s.templateService.GetBroadcastTemplate(
		centerName,
		title,
		message,
		warning,
		actionLabel,
		actionURL,
	)

	// 設定 altText
	altTextPrefix := "🔔 廣播通知"
	if messageType == "URGENT" {
		altTextPrefix = "🚨 緊急通知"
	}

	// 包裝為 Flex Message 格式
	lineMessage := map[string]interface{}{
		"type":     "flex",
		"altText":  fmt.Sprintf("%s - %s", altTextPrefix, title),
		"contents": flexMessage,
	}

	// 【記錄廣播請求】
	if err := s.rateLimiter.Record(ctx, adminID, centerID); err != nil {
		s.Logger.Error("failed to record broadcast", "error", err)
	}

	// 使用 Multicast 發送訊息
	err = s.lineBotService.Multicast(ctx, lineUserIDs, lineMessage)
	if err != nil {
		s.Logger.Error("multicast failed", "error", err, "user_count", len(lineUserIDs))
		return &BroadcastResult{
			SuccessCount: 0,
			FailedCount:  len(lineUserIDs),
			TotalCount:   len(teachers),
			Message:      fmt.Sprintf("發送失敗：%s", err.Error()),
			RateLimit:    rateLimitInfo,
		}, nil, err
	}

	s.Logger.Info("broadcast completed",
		"success_count", len(lineUserIDs),
		"total_teachers", len(teachers),
	)

	return &BroadcastResult{
		SuccessCount: len(lineUserIDs),
		FailedCount:  0,
		TotalCount:   len(teachers),
		Message:      fmt.Sprintf("成功發送給 %d 位老師", len(lineUserIDs)),
		RateLimit:    rateLimitInfo,
	}, nil, nil
}

// getTargetTeachers 取得目標老師清單
// 若 teacherIDs 為空，返回中心所有老師；否則返回指定的老師（需驗證屬於該中心）
func (s *AdminNotificationService) getTargetTeachers(
	ctx context.Context,
	centerID uint,
	teacherIDs []uint,
) ([]models.Teacher, error) {
	if len(teacherIDs) == 0 {
		// 取得中心所有老師
		return s.teacherRepo.ListByCenter(ctx, centerID)
	}

	// 取得指定的老師，並驗證屬於該中心
	var teachers []models.Teacher
	err := s.app.MySQL.RDB.WithContext(ctx).
		Table("teachers").
		Joins("INNER JOIN center_memberships ON center_memberships.teacher_id = teachers.id").
		Where("center_memberships.center_id = ?", centerID).
		Where("center_memberships.status = ?", "ACTIVE").
		Where("teachers.id IN ?", teacherIDs).
		Find(&teachers).Error

	if err != nil {
		return nil, err
	}

	return teachers, nil
}
