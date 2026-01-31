package controllers

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

// AdminUserController 管理員 Controller
type AdminUserController struct {
	BaseController
	app          *app.App
	adminService *services.AdminUserService
}

// NewAdminUserController 建立 AdminUserController
func NewAdminUserController(app *app.App) *AdminUserController {
	return &AdminUserController{
		BaseController: *NewBaseController(app),
		app:            app,
		adminService:   services.NewAdminUserService(app),
	}
}

// LINEBindingResponse LINE 綁定回應
type LINEBindingResponse struct {
	IsBound       bool       `json:"is_bound"`
	LineUserID    string     `json:"line_user_id,omitempty"`
	BoundAt       *time.Time `json:"bound_at,omitempty"`
	NotifyEnabled bool       `json:"notify_enabled"`
	WelcomeSent   bool       `json:"welcome_sent"`
}

// UpdateNotifySettingsRequest 更新通知設定請求
type UpdateNotifySettingsRequest struct {
	Enabled bool `json:"enabled" binding:"required"`
}

// requireAdminID 取得並驗證管理員 ID（通用模式）
func (ctl *AdminUserController) requireAdminID(helper *ContextHelper) uint {
	adminID := helper.MustUserID()
	if adminID == 0 {
		return 0
	}
	return adminID
}

// GetLINEBindingStatus 取得 LINE 綁定狀態
// @Summary 取得 LINE 綁定狀態
// @Description 取得目前管理員的 LINE 綁定狀態
// @Tags Admin - LINE
// @Accept json
// @Produce json
// @Success 200 {object} LINEBindingResponse
// @Router /admin/me/line-binding [get]
func (c *AdminUserController) GetLINEBindingStatus(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	status, errInfo, err := c.adminService.GetLINEBindingStatus(ctx.Request.Context(), adminID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(status)
}

// InitLINEBinding 初始化 LINE 綁定
// @Summary 初始化 LINE 綁定
// @Description 產生綁定驗證碼，開始 LINE 綁定流程
// @Tags Admin - LINE
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "code, expires_at"
// @Router /admin/me/line/bind [post]
func (c *AdminUserController) InitLINEBinding(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	code, expiresAt, errInfo, err := c.adminService.InitLINEBinding(ctx.Request.Context(), adminID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{
		"code":       code,
		"expires_at": expiresAt,
	})
}

// UnbindLINE 解除 LINE 綁定
// @Summary 解除 LINE 綁定
// @Description 解除管理員的 LINE 綁定
// @Tags Admin - LINE
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /admin/me/line/unbind [delete]
func (c *AdminUserController) UnbindLINE(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	errInfo, _ := c.adminService.UnbindLINE(ctx.Request.Context(), adminID)
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{"message": "LINE 帳號已解除綁定"})
}

// UpdateLINENotifySettings 更新 LINE 通知設定
// @Summary 更新 LINE 通知設定
// @Description 更新管理員的 LINE 通知開關
// @Tags Admin - LINE
// @Accept json
// @Produce json
// @Param request body UpdateNotifySettingsRequest true "通知設定"
// @Success 200 {object} map[string]string
// @Router /admin/me/line/notify-settings [patch]
func (c *AdminUserController) UpdateLINENotifySettings(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	var req UpdateNotifySettingsRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	errInfo, _ := c.adminService.UpdateLINENotifySettings(ctx.Request.Context(), adminID, req.Enabled)
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{
		"message":        "通知設定已更新",
		"notify_enabled": req.Enabled,
	})
}

// GetLINENotifySettings 取得 LINE 通知設定
// @Summary 取得 LINE 通知設定
// @Description 取得管理員的 LINE 通知開關狀態
// @Tags Admin - LINE
// @Accept json
// @Produce json
// @Success 200 {object} global.ApiResponse
// @Router /admin/me/line/notify-settings [get]
func (c *AdminUserController) GetLINENotifySettings(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	status, errInfo, err := c.adminService.GetLINEBindingStatus(ctx.Request.Context(), adminID)
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{
		"is_bound":        status.IsBound,
		"notify_enabled":  status.NotifyEnabled,
		"bound_at":        status.BoundAt,
	})
}

// GetAdminProfile 取得管理員個人資料
// @Summary 取得管理員個人資料
// @Description 取得目前登入管理員的個人資料
// @Tags Admin - Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.AdminUser
// @Router /admin/me/profile [get]
func (c *AdminUserController) GetAdminProfile(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	admin, errInfo, err := c.adminService.GetAdminProfile(ctx.Request.Context(), adminID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(admin)
}

// ChangePassword 修改管理員密碼
// @Summary 修改管理員密碼
// @Description 修改目前登入管理員的密碼
// @Tags Admin - Profile
// @Accept json
// @Produce json
// @Param request body services.ChangePasswordRequest true "密碼修改請求"
// @Success 200 {object} map[string]string
// @Router /admin/me/change-password [post]
func (c *AdminUserController) ChangePassword(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	var req services.ChangePasswordRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	errInfo, err := c.adminService.ChangePassword(ctx.Request.Context(), adminID, &req)
	if err != nil {
		if errInfo != nil && errInfo.Code == errInfos.PASSWORD_NOT_MATCH {
			helper.BadRequest(errInfo.Msg)
		} else {
			helper.ErrorWithInfo(errInfo)
		}
		return
	}

	helper.Success(gin.H{"message": "密碼已成功修改"})
}

// ToggleAdminStatusRequest 切換管理員狀態請求
type ToggleAdminStatusRequest struct {
	TargetAdminID uint   `json:"target_admin_id" binding:"required"`
	NewStatus     string `json:"new_status" binding:"required,oneof=ACTIVE INACTIVE"`
}

// ListAdmins 取得管理員列表
// @Summary 取得管理員列表
// @Description 取得目前中心的所有管理員列表
// @Tags Admin - Management
// @Accept json
// @Produce json
// @Success 200 {object} []services.AdminListItem
// @Router /admin/admins [get]
func (c *AdminUserController) ListAdmins(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	// 取得操作者資料以獲取 center_id
	operator, errInfo, err := c.adminService.GetAdminProfile(ctx.Request.Context(), adminID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	admins, listErrInfo, listErr := c.adminService.ListAdmins(ctx.Request.Context(), operator.CenterID)
	if listErr != nil {
		helper.ErrorWithInfo(listErrInfo)
		return
	}

	helper.Success(admins)
}

// ToggleAdminStatus 停用/啟用管理員
// @Summary 停用/啟用管理員
// @Description 停用或啟用指定的管理員帳號（僅 OWNER 可執行）
// @Tags Admin - Management
// @Accept json
// @Produce json
// @Param request body ToggleAdminStatusRequest true "狀態變更請求"
// @Success 200 {object} map[string]string
// @Router /admin/admins/toggle-status [post]
func (c *AdminUserController) ToggleAdminStatus(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	var req ToggleAdminStatusRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	errInfo, err := c.adminService.ToggleAdminStatus(ctx.Request.Context(), adminID, req.TargetAdminID, req.NewStatus)
	if err != nil {
		if errInfo != nil {
			if errInfo.Code == errInfos.FORBIDDEN || errInfo.Code == errInfos.ADMIN_CANNOT_DISABLE_SELF {
				helper.Forbidden(errInfo.Msg)
			} else {
				helper.ErrorWithInfo(errInfo)
			}
		} else {
			helper.InternalError(err.Error())
		}
		return
	}

	statusText := "已停用"
	if req.NewStatus == "ACTIVE" {
		statusText = "已啟用"
	}

	helper.Success(gin.H{"message": statusText})
}

// ResetAdminPasswordRequest 重設管理員密碼請求
type ResetAdminPasswordRequest struct {
	TargetAdminID uint   `json:"target_admin_id" binding:"required"`
	NewPassword   string `json:"new_password" binding:"required,min=6"`
}

// ResetAdminPassword 重設管理員密碼
// @Summary 重設管理員密碼
// @Description 重設指定管理員的密碼（僅 OWNER 可執行）
// @Tags Admin - Management
// @Accept json
// @Produce json
// @Param request body ResetAdminPasswordRequest true "密碼重設請求"
// @Success 200 {object} map[string]string
// @Router /admin/admins/reset-password [post]
func (c *AdminUserController) ResetAdminPassword(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	var req ResetAdminPasswordRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	errInfo, err := c.adminService.ResetAdminPassword(ctx.Request.Context(), adminID, req.TargetAdminID, req.NewPassword)
	if err != nil {
		if errInfo != nil && errInfo.Code == errInfos.FORBIDDEN {
			helper.Forbidden(errInfo.Msg)
		} else {
			helper.ErrorWithInfo(errInfo)
		}
		return
	}

	helper.Success(gin.H{"message": "密碼已重設"})
}

// ChangeAdminRoleRequest 修改管理員角色請求
type ChangeAdminRoleRequest struct {
	TargetAdminID uint   `json:"target_admin_id" binding:"required"`
	NewRole       string `json:"new_role" binding:"required,oneof=ADMIN STAFF OWNER"`
}

// ChangeAdminRole 修改管理員角色
// @Summary 修改管理員角色
// @Description 修改指定管理員的角色（僅 OWNER 可執行，且不能修改自己）
// @Tags Admin - Management
// @Accept json
// @Produce json
// @Param request body ChangeAdminRoleRequest true "角色變更請求"
// @Success 200 {object} map[string]string
// @Router /admin/admins/change-role [post]
func (c *AdminUserController) ChangeAdminRole(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	var req ChangeAdminRoleRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	errInfo, err := c.adminService.ChangeAdminRole(ctx.Request.Context(), adminID, req.TargetAdminID, req.NewRole)
	if err != nil {
		if errInfo != nil {
			if errInfo.Code == errInfos.FORBIDDEN || errInfo.Code == errInfos.ADMIN_NOT_FOUND {
				helper.BadRequest(errInfo.Msg)
			} else {
				helper.ErrorWithInfo(errInfo)
			}
		} else {
			helper.InternalError(err.Error())
		}
		return
	}

	helper.Success(gin.H{"message": "角色已更新"})
}

// CreateAdminRequest 建立管理員請求
type CreateAdminRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=1,max=255"`
	Role     string `json:"role" binding:"required,oneof=ADMIN STAFF OWNER"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateAdmin 建立管理員
// @Summary 建立管理員
// @Description 直接建立管理員帳號（僅 OWNER 可執行）
// @Tags Admin - Management
// @Accept json
// @Produce json
// @Param request body CreateAdminRequest true "管理員資料"
// @Success 200 {object} map[string]interface{}
// @Router /admin/admins [post]
func (c *AdminUserController) CreateAdmin(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := c.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	// 檢查操作者權限（必須是 OWNER）
	operator, errInfo, err := c.adminService.GetAdminProfile(ctx.Request.Context(), adminID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	if operator.Role != "OWNER" {
		helper.Forbidden("Only OWNER can create admins")
		return
	}

	var req CreateAdminRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	admin, createErrInfo, createErr := c.adminService.CreateAdmin(ctx.Request.Context(), operator.CenterID, req.Email, req.Name, req.Role, req.Password)
	if createErr != nil {
		if createErrInfo != nil && createErrInfo.Code == errInfos.ADMIN_EMAIL_EXISTS {
			helper.BadRequest(createErrInfo.Msg)
		} else {
			helper.ErrorWithInfo(createErrInfo)
		}
		return
	}

	helper.Success(gin.H{
		"id":        admin.ID,
		"email":     admin.Email,
		"name":      admin.Name,
		"role":      admin.Role,
		"status":    admin.Status,
		"center_id": admin.CenterID,
	})
}
