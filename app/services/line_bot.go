package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"

	"gorm.io/gorm"
)

// LineBotService LINE Messaging API 服務
type LineBotService interface {
	// 基礎訊息發送
	PushMessage(ctx context.Context, userID string, message interface{}) error
	ReplyMessage(ctx context.Context, replyToken string, message interface{}) error
	Multicast(ctx context.Context, userIDs []string, message interface{}) error

	// Flex Message 發送（便捷方法）
	PushFlexMessage(ctx context.Context, userID string, altText string, flexContent interface{}) error
	ReplyFlexMessage(ctx context.Context, replyToken string, altText string, flexContent interface{}) error

	// 驗證
	VerifySignature(body []byte, signature string) bool

	// 身分識別
	GetCombinedIdentity(lineUserID string) (*CombinedIdentity, error)

	// 行程聚合
	GetAggregatedAgenda(lineUserID string, targetDate *time.Time) ([]AgendaItem, error)

	// 範本發送
	SendWelcomeTeacher(ctx context.Context, teacher *models.Teacher, centerName string) error
	SendWelcomeAdmin(ctx context.Context, admin *models.AdminUser, centerName string) error
	SendExceptionNotification(ctx context.Context, admin *models.AdminUser, exception *models.ScheduleException, teacherName string) error
	SendInvitationAcceptedNotification(ctx context.Context, admins []*models.AdminUser, teacher *models.Teacher, centerName string, role string) error
}

// LineBotServiceImpl LINE Messaging API 服務實現
type LineBotServiceImpl struct {
	app                  *app.App
	channelSecret        string
	channelToken         string
	apiURL               string
	profileURL           string
	replyURL             string
	multicastURL         string
	client               *http.Client
	templateService      LineBotTemplateService
	scheduleExpansionSvc ScheduleExpansionService
	personalEventSvc     *PersonalEventService
}

// NewLineBotService 建立 LINE Bot Service
func NewLineBotService(app *app.App) LineBotService {
	svc := &LineBotServiceImpl{
		app: app,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	if app.Env != nil {
		svc.channelSecret = app.Env.LineChannelSecret
		svc.channelToken = app.Env.LineChannelAccessToken
		svc.apiURL = "https://api.line.me/v2/bot/message/push"
		svc.profileURL = "https://api.line.me/v2/bot/profile"
		svc.replyURL = "https://api.line.me/v2/bot/message/reply"
		svc.multicastURL = "https://api.line.me/v2/bot/message/multicast"
		svc.templateService = NewLineBotTemplateService(app.Env.FrontendBaseURL)
	}

	if app.MySQL != nil {
		svc.scheduleExpansionSvc = NewScheduleExpansionService(app)
		svc.personalEventSvc = NewPersonalEventService(app)
	}

	return svc
}

// LineBotRequest LINE Bot API 請求結構
type LineBotRequest struct {
	Messages []interface{} `json:"messages"`
}

// LineBotReplyRequest LINE Bot Reply API 請求結構
type LineBotReplyRequest struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

// LineBotMulticastRequest LINE Bot Multicast API 請求結構
type LineBotMulticastRequest struct {
	To       []string      `json:"to"`
	Messages []interface{} `json:"messages"`
}

// LineBotErrorResponse LINE Bot API 錯誤回應
type LineBotErrorResponse struct {
	Message string `json:"message"`
	Details []struct {
		Property string `json:"property"`
		Message  string `json:"message"`
	} `json:"details"`
}

// AdminProfile 管理員簡化資料（用於顯示）
type AdminProfile struct {
	ID       uint   `json:"id"`
	CenterID uint   `json:"center_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// LineBotTeacherProfile 老師簡化資料（用於 LINE Bot 顯示）
// 注意：避免與 cache.go 中的 TeacherProfile 衝突
type LineBotTeacherProfile struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	City      string `json:"city"`
	District  string `json:"district"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// CombinedIdentity 整合身份資訊（用於 LINE Bot 回覆使用者身份相關問題）
// 組合管理員資料與老師資料，判斷主要角色
type CombinedIdentity struct {
	// 管理員資料列表（可能隸屬多個中心的管理員）
	AdminProfiles []AdminProfile `json:"admin_profiles"`
	// 老師個人檔案（若為老師角色）
	TeacherProfile *LineBotTeacherProfile `json:"teacher_profile,omitempty"`
	// 老師所屬中心會員關係
	Memberships []models.CenterMembership `json:"memberships"`
	// 主要角色：ADMIN 或 TEACHER（若有重疊，ADMIN 優先）
	PrimaryRole string `json:"primary_role"`
}

// AgendaSourceType 日程來源類型
type AgendaSourceType string

const (
	AgendaSourceTypeCenter   AgendaSourceType = "CENTER"   // 來自中心（排課）
	AgendaSourceTypePersonal AgendaSourceType = "PERSONAL" // 個人行程
)

// AgendaItem 日程項目（用於 LINE Bot 查詢今日/本週行程）
type AgendaItem struct {
	// 開始時間（格式：15:04）
	Time string `json:"time"`
	// 行程標題
	Title string `json:"title"`
	// 來源名稱（中心名稱 或 「個人」）
	SourceName string `json:"source_name"`
	// 來源類型
	SourceType AgendaSourceType `json:"source_type"`
}

// PushMessage 發送推播訊息給單一用戶
func (s *LineBotServiceImpl) PushMessage(ctx context.Context, userID string, message interface{}) error {
	messages := []interface{}{message}
	return s.pushToEndpoint(ctx, s.apiURL, map[string]interface{}{
		"to":       userID,
		"messages": messages,
	})
}

// ReplyMessage 回復訊息（用於 Webhook 回覆）
func (s *LineBotServiceImpl) ReplyMessage(ctx context.Context, replyToken string, message interface{}) error {
	messages := []interface{}{message}
	reqBody := LineBotReplyRequest{
		ReplyToken: replyToken,
		Messages:   messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	return s.sendRequest(ctx, s.replyURL, jsonBody)
}

// Multicast 群發訊息給多個用戶
func (s *LineBotServiceImpl) Multicast(ctx context.Context, userIDs []string, message interface{}) error {
	if len(userIDs) == 0 {
		return nil
	}

	messages := []interface{}{message}
	reqBody := LineBotMulticastRequest{
		To:       userIDs,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	return s.sendRequest(ctx, s.multicastURL, jsonBody)
}

// PushFlexMessage 發送 Flex Message 給單一用戶
func (s *LineBotServiceImpl) PushFlexMessage(ctx context.Context, userID string, altText string, flexContent interface{}) error {
	flexMessage := s.createFlexMessage(altText, flexContent)
	return s.PushMessage(ctx, userID, flexMessage)
}

// ReplyFlexMessage 回復 Flex Message
func (s *LineBotServiceImpl) ReplyFlexMessage(ctx context.Context, replyToken string, altText string, flexContent interface{}) error {
	flexMessage := s.createFlexMessage(altText, flexContent)
	return s.ReplyMessage(ctx, replyToken, flexMessage)
}

// createFlexMessage 建立 Flex Message 結構
func (s *LineBotServiceImpl) createFlexMessage(altText string, contents interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "flex",
		"altText":  altText,
		"contents": contents,
	}
}

// pushToEndpoint 發送到指定端點
func (s *LineBotServiceImpl) pushToEndpoint(ctx context.Context, endpoint string, body map[string]interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	return s.sendRequest(ctx, endpoint, jsonBody)
}

// sendRequest 發送 HTTP 請求到 LINE API
func (s *LineBotServiceImpl) sendRequest(ctx context.Context, endpoint string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.channelToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		var errResp LineBotErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return fmt.Errorf("LINE API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("LINE API error: %s", errResp.Message)
	}

	return nil
}

// VerifySignature 驗證 LINE Webhook 簽名
func (s *LineBotServiceImpl) VerifySignature(body []byte, signature string) bool {
	expectedSignature := s.generateSignature(body)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// generateSignature 生成簽名
func (s *LineBotServiceImpl) generateSignature(body []byte) string {
	hash := hmac.New(sha256.New, []byte(s.channelSecret))
	hash.Write(body)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// SendWelcomeTeacher 發送老師歡迎訊息
func (s *LineBotServiceImpl) SendWelcomeTeacher(ctx context.Context, teacher *models.Teacher, centerName string) error {
	if teacher.LineUserID == "" {
		return nil // 沒有 LINE ID，無法發送
	}

	flexMessage := s.templateService.GetWelcomeTeacherTemplate(teacher, centerName)
	return s.PushFlexMessage(ctx, teacher.LineUserID, "歡迎加入 TimeLedger！", flexMessage)
}

// SendWelcomeAdmin 發送管理員歡迎訊息
func (s *LineBotServiceImpl) SendWelcomeAdmin(ctx context.Context, admin *models.AdminUser, centerName string) error {
	if admin.LineUserID == "" {
		return nil // 沒有 LINE ID，無法發送
	}

	flexMessage := s.templateService.GetWelcomeAdminTemplate(admin, centerName)
	return s.PushFlexMessage(ctx, admin.LineUserID, "歡迎使用 TimeLedger！", flexMessage)
}

// SendExceptionNotification 發送例外申請通知給管理員
func (s *LineBotServiceImpl) SendExceptionNotification(ctx context.Context, admin *models.AdminUser, exception *models.ScheduleException, teacherName string) error {
	if admin.LineUserID == "" || !admin.LineNotifyEnabled {
		return nil // 沒有綁定或通知已關閉
	}

	flexMessage := s.templateService.GetExceptionSubmitTemplate(exception, teacherName, "")
	return s.PushFlexMessage(ctx, admin.LineUserID, "新的例外申請通知", flexMessage)
}

// SendInvitationAcceptedNotification 發送邀請接受通知給管理員
func (s *LineBotServiceImpl) SendInvitationAcceptedNotification(ctx context.Context, admins []*models.AdminUser, teacher *models.Teacher, centerName string, role string) error {
	if len(admins) == 0 {
		return nil // 沒有管理員可通知
	}

	flexMessage := s.templateService.GetInvitationAcceptedTemplate(teacher, centerName, role)

	// 收集所有已綁定 LINE 的管理員
	var lineUserIDs []string
	for _, admin := range admins {
		if admin.LineUserID != "" && admin.LineNotifyEnabled {
			lineUserIDs = append(lineUserIDs, admin.LineUserID)
		}
	}

	if len(lineUserIDs) == 0 {
		return nil // 沒有已綁定的管理員
	}

	// 群發通知給所有已綁定的管理員
	return s.Multicast(ctx, lineUserIDs, flexMessage)
}

// GetCombinedIdentity 查詢整合身份資訊
// 優先查詢管理員資料，若無則查詢老師資料
// 若兩者都找不到，回傳訪客身份
func (s *LineBotServiceImpl) GetCombinedIdentity(lineUserID string) (*CombinedIdentity, error) {
	identity := &CombinedIdentity{
		PrimaryRole: "GUEST",
	}

	// 優先查詢管理員（ADMIN 角色優先）
	adminRepo := repositories.NewAdminUserRepository(s.app)
	admin, err := adminRepo.GetByLineUserID(context.Background(), lineUserID)
	if err == nil && admin.ID != 0 {
		// 找到管理員身份
		identity.AdminProfiles = []AdminProfile{
			{
				ID:       admin.ID,
				CenterID: admin.CenterID,
				Name:     admin.Name,
				Email:    admin.Email,
				Role:     admin.Role,
			},
		}
		identity.PrimaryRole = "ADMIN"
		return identity, nil
	}

	// 若非資料庫錯誤（找不到資料視為正常流程）
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查詢管理員失敗: %w", err)
	}

	// 查詢老師資料
	teacherRepo := repositories.NewTeacherRepository(s.app)
	teacher, err := teacherRepo.GetByLineUserID(context.Background(), lineUserID)
	if err == nil && teacher.ID != 0 {
		// 找到老師身份
		teacherProfile := &LineBotTeacherProfile{
			ID:        teacher.ID,
			Name:      teacher.Name,
			Email:     teacher.Email,
			City:      teacher.City,
			District:  teacher.District,
			AvatarURL: teacher.AvatarURL,
		}

		// 預載入老師所屬中心會員關係（含中心名稱）
		var memberships []models.CenterMembership
		err = s.app.MySQL.RDB.WithContext(context.Background()).
			Table("center_memberships").
			Select("center_memberships.*, centers.name as center_name").
			Joins("LEFT JOIN centers ON centers.id = center_memberships.center_id").
			Where("center_memberships.teacher_id = ?", teacher.ID).
			Find(&memberships).Error
		if err != nil {
			return nil, fmt.Errorf("查詢老師會員關係失敗: %w", err)
		}

		identity.TeacherProfile = teacherProfile
		identity.Memberships = memberships
		identity.PrimaryRole = "TEACHER"
		return identity, nil
	}

	// 若非資料庫錯誤（找不到資料視為正常流程）
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查詢老師失敗: %w", err)
	}

	// 兩者都找不到，回傳訪客身份
	return identity, nil
}

// GetAggregatedAgenda 取得當日聚合行程
// 1. 先呼叫 GetCombinedIdentity 獲取身份與會員關係
// 2. 循環各中心調用排課擴展服務獲取當日課表
// 3. 調用 PersonalEventService.GetTodayOccurrences 獲取個人行程
// 4. 將兩者轉換為 AgendaItem 並按 StartTime 排序
func (s *LineBotServiceImpl) GetAggregatedAgenda(lineUserID string, targetDate *time.Time) ([]AgendaItem, error) {
	// 確定目標日期（預設為今天）
	date := time.Now()
	if targetDate != nil {
		date = *targetDate
	}

	// 取得身份資訊
	identity, err := s.GetCombinedIdentity(lineUserID)
	if err != nil {
		return nil, fmt.Errorf("取得身份資訊失敗: %w", err)
	}

	var allItems []AgendaItem

	// 如果是老師角色，取得老師 ID
	var teacherID uint
	if identity.TeacherProfile != nil {
		teacherID = identity.TeacherProfile.ID
	}

	// 遍歷各中心，獲取當日課表
	for _, membership := range identity.Memberships {
		centerID := membership.CenterID

		// 取得中心名稱
		centerName, err := s.getCenterName(context.Background(), centerID)
		if err != nil {
			centerName = "未知中心"
		}

		// 取得該中心當日的課表規則
		rules, err := s.getScheduleRulesForDate(context.Background(), centerID, teacherID, date)
		if err != nil {
			// 記錄錯誤但繼續處理其他中心
			continue
		}

		// 展開規則為當日課表
		var schedules []ExpandedSchedule
		if len(rules) > 0 {
			schedules = s.scheduleExpansionSvc.ExpandRules(context.Background(), rules, date, date, centerID)
		}

		// 將課表轉換為 AgendaItem
		for _, schedule := range schedules {
			item := AgendaItem{
				Time:       schedule.StartTime,
				Title:      schedule.OfferingName,
				SourceName: centerName,
				SourceType: AgendaSourceTypeCenter,
			}
			allItems = append(allItems, item)
		}
	}

	// 如果是老師角色，取得個人行程
	if teacherID > 0 {
		occurrences, _, err := s.personalEventSvc.GetTodayOccurrences(context.Background(), teacherID, date)
		if err != nil {
			// 記錄錯誤但繼續處理
			_ = err
		} else {
			// 將個人行程轉換為 AgendaItem
			for _, occ := range occurrences {
				item := AgendaItem{
					Time:       formatTimeForAgenda(occ.StartAt),
					Title:      occ.Title,
					SourceName: "個人",
					SourceType: AgendaSourceTypePersonal,
				}
				allItems = append(allItems, item)
			}
		}
	}

	// 按時間排序
	sortAgendaItemsByTime(allItems)

	return allItems, nil
}

// getCenterName 取得中心名稱
func (s *LineBotServiceImpl) getCenterName(ctx context.Context, centerID uint) (string, error) {
	centerRepo := repositories.NewCenterRepository(s.app)
	center, err := centerRepo.GetByID(ctx, centerID)
	if err != nil || center.ID == 0 {
		return "", fmt.Errorf("取得中心失敗: %w", err)
	}
	return center.Name, nil
}

// getScheduleRulesForDate 取得指定日期適用的排課規則
func (s *LineBotServiceImpl) getScheduleRulesForDate(ctx context.Context, centerID, teacherID uint, date time.Time) ([]models.ScheduleRule, error) {
	// 查詢該老師在該中心、指定日期的所有有效規則
	var rules []models.ScheduleRule
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // 週日對應 7
	}

	// 處理零值時間
	startDateStr := date.Format("2006-01-02")

	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("teacher_id = ?", teacherID).
		Where("weekday = ?", weekday).
		Where("COALESCE(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')), ''), 'null'), '0001-01-01') <= ?", startDateStr).
		Where("COALESCE(NULLIF(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')), ''), 'null'), '0001-01-01 00:00:00'), '9999-12-31') >= ?", startDateStr).
		Preload("Offering").
		Preload("Teacher").
		Preload("Room").
		Find(&rules).Error

	if err != nil {
		return nil, err
	}

	return rules, nil
}

// sortAgendaItemsByTime 按時間排序 AgendaItem
func sortAgendaItemsByTime(items []AgendaItem) {
	// 使用時間比較器排序
	sort.Slice(items, func(i, j int) bool {
		return compareTimeStrings(items[i].Time, items[j].Time) < 0
	})
}

// formatTimeForAgenda 將 time.Time 轉換為 HH:MM 格式
func formatTimeForAgenda(t time.Time) string {
	return t.Format("15:04")
}
