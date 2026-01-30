package controllers

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	BaseController
	app                 *app.App
	notificationService services.NotificationService
}

type ListNotificationsRequest struct {
	Limit  int `form:"limit,default=20"`
	Offset int `form:"offset,default=0"`
}

type SendNotificationRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

type TeacherNotifyTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func NewNotificationController(app *app.App) *NotificationController {
	return &NotificationController{
		app:                 app,
		notificationService: services.NewNotificationService(app),
	}
}

// requireUserID 取得並驗證使用者 ID（通用模式）
func (ctl *NotificationController) requireUserID(helper *ContextHelper) uint {
	userID := helper.MustUserID()
	if userID == 0 {
		return 0
	}
	return userID
}

// ListNotifications 取得通知列表
// @Summary 取得通知列表
// @Description 取得目前使用者的通知列表
// @Tags Notifications
// @Accept json
// @Produce json
// @Param limit query int false "每頁筆數"
// @Param offset query int false "偏移量"
// @Success 200 {object} global.ApiResponse
// @Router /notifications [get]
func (ctl *NotificationController) ListNotifications(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	userID := ctl.requireUserID(helper)
	if userID == 0 {
		return
	}

	userType, _ := helper.UserType()

	var req ListNotificationsRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	notifications, err := ctl.notificationService.GetNotifications(ctx.Request.Context(), userID, userType, req.Limit, req.Offset)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	// 計算未讀通知數量
	unreadCount := 0
	for _, n := range notifications {
		if !n.IsRead {
			unreadCount++
		}
	}

	helper.Success(gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

// GetUnreadCount 取得未讀通知數量
// @Summary 取得未讀通知數量
// @Description 取得目前使用者的未讀通知數量
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200 {object} global.ApiResponse
// @Router /notifications/unread-count [get]
func (ctl *NotificationController) GetUnreadCount(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	userID := ctl.requireUserID(helper)
	if userID == 0 {
		return
	}

	userType, _ := helper.UserType()

	unreadCount, err := ctl.notificationService.GetUnreadCount(ctx.Request.Context(), userID, userType)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{
		"unread_count": unreadCount,
	})
}

// MarkAsRead 標記通知為已讀
// @Summary 標記通知為已讀
// @Description 將指定通知標記為已讀
// @Tags Notifications
// @Accept json
// @Produce json
// @Param id path uint true "通知ID"
// @Success 200 {object} global.ApiResponse
// @Router /notifications/{id}/read [post]
func (ctl *NotificationController) MarkAsRead(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	notificationID := helper.MustParamUint("id")
	if notificationID == 0 {
		return
	}

	err := ctl.notificationService.MarkAsRead(ctx.Request.Context(), notificationID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead 標記所有通知為已讀
// @Summary 標記所有通知為已讀
// @Description 將所有通知標記為已讀
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200 {object} global.ApiResponse
// @Router /notifications/read-all [post]
func (ctl *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	userID := ctl.requireUserID(helper)
	if userID == 0 {
		return
	}

	userType, _ := helper.UserType()

	err := ctl.notificationService.MarkAllAsRead(ctx.Request.Context(), userID, userType)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{"message": "All notifications marked as read"})
}

// SetNotifyToken 設定通知 Token
// @Summary 設定通知 Token
// @Description 設定老師的 LINE Notify Token
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body TeacherNotifyTokenRequest true "Token 資訊"
// @Success 200 {object} global.ApiResponse
// @Router /notifications/token [post]
func (ctl *NotificationController) SetNotifyToken(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	userID := ctl.requireUserID(helper)
	if userID == 0 {
		return
	}

	userType, _ := helper.UserType()
	if userType != "TEACHER" {
		helper.Forbidden("Only teachers can set notify token")
		return
	}

	var req TeacherNotifyTokenRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	err := ctl.notificationService.SetNotifyToken(ctx.Request.Context(), userID, req.Token)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{"message": "Notify token updated successfully"})
}

// SendTestNotification 發送測試通知
// @Summary 發送測試通知
// @Description 發送一條測試通知給目前使用者
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200 {object} global.ApiResponse
// @Router /notifications/test [post]
func (ctl *NotificationController) SendTestNotification(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	userID := ctl.requireUserID(helper)
	if userID == 0 {
		return
	}

	title := "測試通知"
	message := "這是一條測試通知\n時間: " + time.Now().Format("2006-01-02 15:04:05")

	err := ctl.notificationService.SendTeacherNotification(ctx.Request.Context(), userID, title, message)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{"message": "Test notification sent"})
}

// GetQueueStats 取得通知佇列統計（管理員專用）
// @Summary 取得通知佇列統計
// @Description 取得通知佇列的統計資訊
// @Tags Admin - Notifications
// @Accept json
// @Produce json
// @Success 200 {object} global.ApiResponse
// @Router /admin/notifications/queue-stats [get]
func (ctl *NotificationController) GetQueueStats(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	queueService := services.NewNotificationQueueService(ctl.app)
	stats := queueService.GetQueueStats(ctx.Request.Context())

	helper.Success(stats)
}
