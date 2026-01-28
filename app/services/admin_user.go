package services

import (
	"context"
	"math/rand"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
)

// AdminUserService 管理員服務
type AdminUserService struct {
	BaseService
	app            *app.App
	adminRepo      *repositories.AdminUserRepository
	centerRepo     *repositories.CenterRepository
	lineBotService LineBotService
}

// NewAdminUserService 建立管理員服務
func NewAdminUserService(app *app.App) *AdminUserService {
	return &AdminUserService{
		app:            app,
		adminRepo:      repositories.NewAdminUserRepository(app),
		centerRepo:     repositories.NewCenterRepository(app),
		lineBotService: NewLineBotService(app),
	}
}

// LINEBindingStatus 管理員 LINE 綁定狀態
type LINEBindingStatus struct {
	IsBound       bool       `json:"is_bound"`
	LineUserID    string     `json:"line_user_id,omitempty"`
	AdminID       uint       `json:"admin_id,omitempty"`
	CenterID      uint       `json:"center_id,omitempty"`
	Role          string     `json:"role,omitempty"`
	BoundAt       *time.Time `json:"bound_at,omitempty"`
	NotifyEnabled bool       `json:"notify_enabled"`
	WelcomeSent   bool       `json:"welcome_sent"`
}

// GenerateBindingCode 生成綁定驗證碼
func GenerateBindingCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	const length = 6
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GetLINEBindingStatus 取得管理員 LINE 綁定狀態（依管理員 ID）
func (s *AdminUserService) GetLINEBindingStatus(ctx context.Context, adminID uint) (*LINEBindingStatus, *errInfos.Res, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	status := &LINEBindingStatus{
		IsBound:       admin.LineUserID != "",
		NotifyEnabled: admin.LineNotifyEnabled,
		BoundAt:       admin.LineBoundAt,
	}

	if admin.LineUserID != "" {
		if len(admin.LineUserID) > 8 {
			status.LineUserID = "***" + admin.LineUserID[len(admin.LineUserID)-5:]
		} else {
			status.LineUserID = "***"
		}
	}

	return status, nil, nil
}

// GetLINEBindingStatusByLineUserID 依 LINE User ID 取得管理員綁定狀態
func (s *AdminUserService) GetLINEBindingStatusByLineUserID(ctx context.Context, lineUserID string) (*LINEBindingStatus, *errInfos.Res, error) {
	admin, err := s.adminRepo.GetByLineUserID(ctx, lineUserID)
	if err != nil {
		return nil, nil, nil // 未找到，不視為錯誤
	}

	status := &LINEBindingStatus{
		IsBound:       true,
		AdminID:       admin.ID,
		CenterID:      admin.CenterID,
		Role:          admin.Role,
		NotifyEnabled: admin.LineNotifyEnabled,
		BoundAt:       admin.LineBoundAt,
	}

	if admin.LineUserID != "" {
		if len(admin.LineUserID) > 8 {
			status.LineUserID = "***" + admin.LineUserID[len(admin.LineUserID)-5:]
		} else {
			status.LineUserID = "***"
		}
	}

	return status, nil, nil
}

// InitLINEBinding 初始化 LINE 綁定（產生驗證碼）
func (s *AdminUserService) InitLINEBinding(ctx context.Context, adminID uint) (string, time.Time, *errInfos.Res, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return "", time.Time{}, s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	if admin.LineUserID != "" {
		return "", time.Time{}, s.app.Err.New(errInfos.LINE_ALREADY_BOUND), nil
	}

	code := GenerateBindingCode()
	expiresAt := time.Now().Add(10 * time.Minute)

	updates := map[string]interface{}{
		"line_binding_code":    code,
		"line_binding_expires": expiresAt,
	}

	if err := s.adminRepo.UpdateFields(ctx, adminID, updates); err != nil {
		return "", time.Time{}, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return code, expiresAt, nil, nil
}

// VerifyLINEBinding 驗證 LINE 綁定（由 Webhook 呼叫）
func (s *AdminUserService) VerifyLINEBinding(ctx context.Context, code string, lineUserID string) (uint, *errInfos.Res, error) {
	var admin models.AdminUser
	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("line_binding_code = ?", strings.ToUpper(code)).
		Where("line_binding_expires > ?", time.Now()).
		First(&admin).Error

	if err != nil {
		return 0, s.app.Err.New(errInfos.LINE_BINDING_CODE_INVALID), nil
	}

	now := time.Now()
	updates := map[string]interface{}{
		"line_user_id":         lineUserID,
		"line_binding_code":    "",
		"line_binding_expires": nil,
		"line_bound_at":        now,
		"line_notify_enabled":  true,
	}

	if err := s.adminRepo.UpdateFields(ctx, admin.ID, updates); err != nil {
		return 0, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return admin.ID, nil, nil
}

// UnbindLINE 解除 LINE 綁定
func (s *AdminUserService) UnbindLINE(ctx context.Context, adminID uint) (*errInfos.Res, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	if admin.LineUserID == "" {
		return s.app.Err.New(errInfos.LINE_NOT_BOUND), nil
	}

	updates := map[string]interface{}{
		"line_user_id":         "",
		"line_binding_code":    "",
		"line_binding_expires": nil,
		"line_notify_enabled":  false,
		"line_bound_at":        nil,
	}

	if err := s.adminRepo.UpdateFields(ctx, adminID, updates); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return nil, nil
}

// UpdateLINENotifySettings 更新 LINE 通知設定
func (s *AdminUserService) UpdateLINENotifySettings(ctx context.Context, adminID uint, enabled bool) (*errInfos.Res, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	if admin.LineUserID == "" {
		return s.app.Err.New(errInfos.LINE_NOT_BOUND), nil
	}

	updates := map[string]interface{}{
		"line_notify_enabled": enabled,
	}

	if err := s.adminRepo.UpdateFields(ctx, adminID, updates); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return nil, nil
}

// SendWelcomeMessageIfNeeded 發送歡迎訊息
func (s *AdminUserService) SendWelcomeMessageIfNeeded(ctx context.Context, adminID uint) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}

	if admin.LineUserID == "" || admin.LineBoundAt == nil {
		return nil
	}

	if time.Since(*admin.LineBoundAt) < 1*time.Minute && admin.LineNotifyEnabled {
		center, err := s.centerRepo.GetByID(ctx, admin.CenterID)
		if err != nil {
			return err
		}

		queueService := NewNotificationQueueService(s.app)
		return queueService.NotifyWelcomeAdmin(ctx, &admin, center.Name)
	}

	return nil
}
