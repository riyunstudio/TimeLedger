package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"

	"golang.org/x/crypto/bcrypt"
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

// AdminProfileResponse 管理員個人資料回應
type AdminProfileResponse struct {
	ID                uint       `json:"id"`
	CenterID          uint       `json:"center_id"`
	CenterName        string     `json:"center_name"`
	Email             string     `json:"email"`
	Name              string     `json:"name"`
	Role              string     `json:"role"`
	Status            string     `json:"status"`
	LineNotifyEnabled bool       `json:"line_notify_enabled"`
	LineBoundAt       *time.Time `json:"line_bound_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
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

	if time.Since(*admin.LineBoundAt) < 1*time.Minute && admin.LineNotifyEnabled && admin.CenterID != 0 {
		center, err := s.centerRepo.GetByID(ctx, admin.CenterID)
		if err != nil {
			return err
		}

		queueService := NewNotificationQueueService(s.app)
		return queueService.NotifyWelcomeAdmin(ctx, &admin, center.Name)
	}

	return nil
}

// ChangePasswordRequest 修改密碼請求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改管理員密碼
func (s *AdminUserService) ChangePassword(ctx context.Context, adminID uint, req *ChangePasswordRequest) (*errInfos.Res, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 驗證舊密碼
	if !s.adminRepo.VerifyPassword(ctx, admin.Email, req.OldPassword) {
		return s.app.Err.New(errInfos.PASSWORD_NOT_MATCH), nil
	}

	// 加密新密碼
	hashedPassword, err := s.adminRepo.HashPassword(req.NewPassword)
	if err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 更新密碼
	updates := map[string]interface{}{
		"password": hashedPassword,
	}

	if err := s.adminRepo.UpdateFields(ctx, adminID, updates); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return nil, nil
}

// GetAdminProfile 取得管理員個人資料
func (s *AdminUserService) GetAdminProfile(ctx context.Context, adminID uint) (*AdminProfileResponse, *errInfos.Res, error) {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 查詢中心名稱
	centerName := ""
	center, err := s.centerRepo.GetByID(ctx, admin.CenterID)
	if err == nil {
		centerName = center.Name
	}

	profile := &AdminProfileResponse{
		ID:                admin.ID,
		CenterID:          admin.CenterID,
		CenterName:        centerName,
		Email:             admin.Email,
		Name:              admin.Name,
		Role:              admin.Role,
		Status:            admin.Status,
		LineNotifyEnabled: admin.LineNotifyEnabled,
		LineBoundAt:       admin.LineBoundAt,
		CreatedAt:         admin.CreatedAt,
	}

	return profile, nil, nil
}

// AdminListItem 管理員列表項目
type AdminListItem struct {
	ID                uint       `json:"id"`
	Email             string     `json:"email"`
	Name              string     `json:"name"`
	Role              string     `json:"role"`
	Status            string     `json:"status"`
	LineUserID        string     `json:"line_user_id,omitempty"`
	LineBoundAt       *time.Time `json:"line_bound_at,omitempty"`
	LineNotifyEnabled bool       `json:"line_notify_enabled"`
	CreatedAt         time.Time  `json:"created_at"`
}

// ListAdmins 取得管理員列表
func (s *AdminUserService) ListAdmins(ctx context.Context, centerID uint) ([]AdminListItem, *errInfos.Res, error) {
	admins, err := s.adminRepo.GetByCenterID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	result := make([]AdminListItem, 0, len(admins))
	for _, admin := range admins {
		item := AdminListItem{
			ID:                admin.ID,
			Email:             admin.Email,
			Name:              admin.Name,
			Role:              admin.Role,
			Status:            admin.Status,
			LineUserID:        admin.LineUserID,
			LineBoundAt:       admin.LineBoundAt,
			LineNotifyEnabled: admin.LineNotifyEnabled,
			CreatedAt:         admin.CreatedAt,
		}
		result = append(result, item)
	}

	return result, nil, nil
}

// ToggleAdminStatus 切換管理員狀態
func (s *AdminUserService) ToggleAdminStatus(ctx context.Context, adminID uint, targetAdminID uint, newStatus string) (*errInfos.Res, error) {
	// 檢查操作者權限
	operator, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 只有 OWNER 可以停用/啟用管理員
	if operator.Role != "OWNER" {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 不能停用自己
	if adminID == targetAdminID {
		return s.app.Err.New(errInfos.INVALID_STATUS), nil
	}

	// 檢查目標管理員
	target, err := s.adminRepo.GetByID(ctx, targetAdminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 不能停用 OWNER
	if target.Role == "OWNER" {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 檢查是否為同一中心
	if operator.CenterID != target.CenterID {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	updates := map[string]interface{}{
		"status": newStatus,
	}

	if err := s.adminRepo.UpdateFields(ctx, targetAdminID, updates); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return nil, nil
}

// ResetAdminPassword 重設管理員密碼（僅 OWNER 可執行）
func (s *AdminUserService) ResetAdminPassword(ctx context.Context, adminID uint, targetAdminID uint, newPassword string) (*errInfos.Res, error) {
	// 檢查操作者權限
	operator, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 只有 OWNER 可以重設密碼
	if operator.Role != "OWNER" {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 不能重設 OWNER 的密碼
	if targetAdminID == adminID {
		return s.app.Err.New(errInfos.INVALID_STATUS), nil
	}

	// 檢查目標管理員
	target, err := s.adminRepo.GetByID(ctx, targetAdminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 檢查是否為同一中心
	if operator.CenterID != target.CenterID {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 不能重設 OWNER 的密碼
	if target.Role == "OWNER" {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 加密新密碼
	hashedPassword, err := s.adminRepo.HashPassword(newPassword)
	if err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 更新密碼
	updates := map[string]interface{}{
		"password": hashedPassword,
	}

	if err := s.adminRepo.UpdateFields(ctx, targetAdminID, updates); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return nil, nil
}

// ChangeAdminRole 修改管理員角色（僅 OWNER 可執行）
func (s *AdminUserService) ChangeAdminRole(ctx context.Context, adminID uint, targetAdminID uint, newRole string) (*errInfos.Res, error) {
	// 檢查操作者權限
	operator, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 只有 OWNER 可以修改角色
	if operator.Role != "OWNER" {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 不能修改自己的角色
	if targetAdminID == adminID {
		return s.app.Err.New(errInfos.INVALID_STATUS), nil
	}

	// 檢查目標管理員
	target, err := s.adminRepo.GetByID(ctx, targetAdminID)
	if err != nil {
		return s.app.Err.New(errInfos.ADMIN_NOT_FOUND), err
	}

	// 檢查是否為同一中心
	if operator.CenterID != target.CenterID {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 不能修改 OWNER 的角色
	if target.Role == "OWNER" {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 更新角色
	updates := map[string]interface{}{
		"role": newRole,
	}

	if err := s.adminRepo.UpdateFields(ctx, targetAdminID, updates); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return nil, nil
}

// AdminInvitation 管理員邀請
type AdminInvitation struct {
	Code      string    `json:"code"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CenterID  uint      `json:"center_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateAdminInvitation 建立管理員邀請
func (s *AdminUserService) CreateAdminInvitation(ctx context.Context, centerID uint, email string, name string, role string) (*AdminInvitation, *errInfos.Res, error) {
	// 檢查 Email 是否已存在
	existingAdmin, err := s.adminRepo.GetByEmail(ctx, email)
	if err == nil && existingAdmin.ID > 0 {
		return nil, s.app.Err.New(errInfos.ADMIN_EMAIL_EXISTS), fmt.Errorf("email already exists")
	}

	// 生成邀請碼
	code := GenerateBindingCode()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天過期

	// 創建邀請記錄
	invitation := AdminInvitation{
		Code:      code,
		Email:     email,
		Name:      name,
		Role:      role,
		CenterID:  centerID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	// 這裡可以選擇發送 email 或 LINE 通知
	// 目前先返回邀請碼，由管理員手動發送

	return &invitation, nil, nil
}

// CreateAdmin 直接建立管理員
func (s *AdminUserService) CreateAdmin(ctx context.Context, centerID uint, email string, name string, role string, password string) (*models.AdminUser, *errInfos.Res, error) {
	// 檢查 Email 是否已存在
	existingAdmin, err := s.adminRepo.GetByEmail(ctx, email)
	if err == nil && existingAdmin.ID > 0 {
		return nil, s.app.Err.New(errInfos.ADMIN_EMAIL_EXISTS), fmt.Errorf("email already exists")
	}

	// 密碼加密
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 建立管理員
	admin := models.AdminUser{
		CenterID:     centerID,
		Email:        email,
		Name:         name,
		Role:         role,
		PasswordHash: string(passwordHash),
		Status:       "ACTIVE",
	}

	createdAdmin, err := s.adminRepo.Create(ctx, admin)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return &createdAdmin, nil, nil
}
