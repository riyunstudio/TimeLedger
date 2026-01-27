package controllers

import (
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/global"

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

// GetLINEBindingStatus 取得 LINE 綁定狀態
// @Summary 取得 LINE 綁定狀態
// @Description 取得目前管理員的 LINE 綁定狀態
// @Tags Admin - LINE
// @Accept json
// @Produce json
// @Success 200 {object} LINEBindingResponse
// @Router /admin/me/line-binding [get]
func (c *AdminUserController) GetLINEBindingStatus(ctx *gin.Context) {
	// 從 gin context 取得 admin ID
	adminIDVal, exists := ctx.Get(string(global.UserIDKey))
	if !exists {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Unauthorized",
		})
		return
	}
	adminID := adminIDVal.(uint)

	status, eInfo, err := c.adminService.GetLINEBindingStatus(ctx.Request.Context(), adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    eInfo.Code,
			Message: eInfo.Msg,
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   status,
	})
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
	adminIDVal, exists := ctx.Get(string(global.UserIDKey))
	if !exists {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Unauthorized",
		})
		return
	}
	adminID := adminIDVal.(uint)

	code, expiresAt, eInfo, err := c.adminService.InitLINEBinding(ctx.Request.Context(), adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    eInfo.Code,
			Message: eInfo.Msg,
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas: gin.H{
			"code":       code,
			"expires_at": expiresAt,
		},
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
	adminIDVal, exists := ctx.Get(string(global.UserIDKey))
	if !exists {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Unauthorized",
		})
		return
	}
	adminID := adminIDVal.(uint)

	eInfo, err := c.adminService.UnbindLINE(ctx.Request.Context(), adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    eInfo.Code,
			Message: eInfo.Msg,
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas: gin.H{
			"message": "LINE 帳號已解除綁定",
		},
	})
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
	adminIDVal, exists := ctx.Get(string(global.UserIDKey))
	if !exists {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Unauthorized",
		})
		return
	}
	adminID := adminIDVal.(uint)

	var req UpdateNotifySettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	eInfo, err := c.adminService.UpdateLINENotifySettings(ctx.Request.Context(), adminID, req.Enabled)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    eInfo.Code,
			Message: eInfo.Msg,
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas: gin.H{
			"message":         "通知設定已更新",
			"notify_enabled": req.Enabled,
		},
	})
}
