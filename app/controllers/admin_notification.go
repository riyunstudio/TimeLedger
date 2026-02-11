package controllers

import (
	"net/http"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

// AdminNotificationController 管理員通知控制器
type AdminNotificationController struct {
	app                 *app.App
	logger              *services.ServiceLogger
	lineBotService      services.LineBotService
	notificationService *services.AdminNotificationService
}

// NewAdminNotificationController 建立管理員通知控制器
func NewAdminNotificationController(app *app.App) *AdminNotificationController {
	return &AdminNotificationController{
		app:                 app,
		logger:              services.NewServiceLogger("AdminNotificationController"),
		lineBotService:      services.NewLineBotService(app),
		notificationService: services.NewAdminNotificationService(app),
	}
}

// BroadcastRequest 廣播訊息請求結構
type BroadcastRequest struct {
	// 訊息類型（必填）
	MessageType string `json:"message_type" binding:"required,oneof=GENERAL URGENT"`
	// 標題（必填，最多 50 字）
	Title string `json:"title" binding:"required,max=50"`
	// 訊息內容（必填，最多 2000 字）
	Message string `json:"message" binding:"required,max=2000"`
	// 警告提示（可選，最多 200 字）
	Warning string `json:"warning,omitempty" binding:"max=200"`
	// 按鈕文字（可選，最多 20 字）
	ActionLabel string `json:"action_label,omitempty" binding:"max=20"`
	// 按鈕連結（可選）
	ActionURL string `json:"action_url,omitempty"`
	// 接收對象過濾條件（可選）
	// 若為空則發送給中心所有已綁定 LINE 的老師
	TeacherIDs []uint `json:"teacher_ids,omitempty"`
}

// BroadcastResponse 廣播訊息回應結構
type BroadcastResponse struct {
	// 成功發送數量
	SuccessCount int `json:"success_count"`
	// 失敗發送數量
	FailedCount int `json:"failed_count"`
	// 總共發送數量
	TotalCount int `json:"total_count"`
	// 訊息
	Message string `json:"message"`
}

// Broadcast 廣播訊息給中心老師
// @Summary 廣播訊息給中心老師
// @Description 管理員發送 LINE 廣播訊息給所屬中心的老師
// @Tags Admin Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BroadcastRequest true "廣播訊息請求"
// @Success 200 {object} global.ApiResponse{data=BroadcastResponse}
// @Failure 400 {object} global.ApiResponse
// @Failure 401 {object} global.ApiResponse
// @Failure 403 {object} global.ApiResponse
// @Failure 500 {object} global.ApiResponse
// @Router /api/v1/admin/notifications/broadcast [POST]
func (c *AdminNotificationController) Broadcast(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	// 取得管理員 ID 和中心 ID
	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	// 解析請求
	var req BroadcastRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 驗證訊息內容
	if len(req.Message) == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "訊息內容不能為空",
		})
		return
	}

	// 呼叫 Service 層處理廣播
	result, eInfo, err := c.notificationService.BroadcastToTeachers(
		ctx.Request.Context(),
		centerID,
		adminID,
		req.MessageType,
		req.Title,
		req.Message,
		req.Warning,
		req.ActionLabel,
		req.ActionURL,
		req.TeacherIDs,
	)

	if err != nil {
		c.logger.Error("broadcast failed", "error", err, "center_id", centerID, "admin_id", adminID)
		if eInfo != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    eInfo.Code,
				Message: eInfo.Msg,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: "廣播訊息失敗",
		})
		return
	}

	// 回傳結果
	response := BroadcastResponse{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		TotalCount:   result.TotalCount,
		Message:      result.Message,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "廣播訊息發送完成",
		Datas:   response,
	})
}
