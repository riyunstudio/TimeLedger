package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/global"

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

func (ctl *NotificationController) ListNotifications(ctx *gin.Context) {
	userID := ctx.GetUint(global.UserIDKey)
	userType := ctx.GetString(global.UserTypeKey)

	var req ListNotificationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid query parameters",
		})
		return
	}

	notifications, err := ctl.notificationService.GetNotifications(ctx, userID, userType, req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas: gin.H{
			"notifications": notifications,
			"unread_count":  len(notifications),
		},
	})
}

func (ctl *NotificationController) GetUnreadCount(ctx *gin.Context) {
	userID := ctx.GetUint(global.UserIDKey)
	userType := ctx.GetString(global.UserTypeKey)

	notificationRepo := repositories.NewNotificationRepository(ctl.app)
	notifications, err := notificationRepo.ListUnread(ctx, userID, userType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas: gin.H{
			"unread_count": len(notifications),
		},
	})
}

func (ctl *NotificationController) MarkAsRead(ctx *gin.Context) {
	id := ctx.Param("id")
	var idInt uint
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid notification ID",
		})
		return
	}

	if err := ctl.notificationService.MarkAsRead(ctx, idInt); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Notification marked as read",
	})
}

func (ctl *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID := ctx.GetUint(global.UserIDKey)
	userType := ctx.GetString(global.UserTypeKey)

	if err := ctl.notificationService.MarkAllAsRead(ctx, userID, userType); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "All notifications marked as read",
	})
}

func (ctl *NotificationController) SetNotifyToken(ctx *gin.Context) {
	var req TeacherNotifyTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	userID := ctx.GetUint(global.UserIDKey)
	userType := ctx.GetString(global.UserTypeKey)

	if userType != "TEACHER" {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    global.FORBIDDEN,
			Message: "Only teachers can set notify token",
		})
		return
	}

	teacherRepo := repositories.NewTeacherRepository(ctl.app)
	teacher, err := teacherRepo.GetByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	teacher.LineNotifyToken = req.Token
	if err := teacherRepo.Update(ctx, teacher); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Notify token updated successfully",
	})
}

func (ctl *NotificationController) SendTestNotification(ctx *gin.Context) {
	userID := ctx.GetUint(global.UserIDKey)

	title := "測試通知"
	message := "這是一條測試通知\n時間: " + time.Now().Format("2006-01-02 15:04:05")

	if err := ctl.notificationService.SendTeacherNotification(ctx, userID, title, message); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Test notification sent",
	})
}
