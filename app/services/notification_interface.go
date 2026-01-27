package services

import (
	"context"
	"time"
	"timeLedger/app/models"
)

type NotificationService interface {
	// SendTeacherNotification 發送通知給老師
	SendTeacherNotification(ctx context.Context, teacherID uint, title, message string) error

	// SendTeacherNotificationWithType 發送通知給老師（帶類型）
	SendTeacherNotificationWithType(ctx context.Context, teacherID uint, title, message string, notificationType string) error

	// SendAdminNotification 發送通知給管理員
	SendAdminNotification(ctx context.Context, centerID uint, title, message string, notificationType string) error

	// SendScheduleReminder 發送排課提醒
	SendScheduleReminder(ctx context.Context, ruleID uint, date time.Time) error

	// SendExceptionNotification 發送例外單通知
	SendExceptionNotification(ctx context.Context, exceptionID uint) error

	// SendReviewNotification 發送審核通知
	SendReviewNotification(ctx context.Context, exceptionID uint, approved bool) error

	// SendTalentInvitationNotification 發送人才庫邀請通知
	SendTalentInvitationNotification(ctx context.Context, teacherID uint, centerName string, inviteToken string) error

	// CreateNotificationRecord 創建通知記錄
	CreateNotificationRecord(ctx context.Context, notification *models.Notification) error

	// GetNotifications 獲取用戶的通知列表
	GetNotifications(ctx context.Context, userID uint, userType string, limit int, offset int) ([]models.Notification, error)

	// MarkAsRead 標記通知為已讀
	MarkAsRead(ctx context.Context, notificationID uint) error

	// MarkAllAsRead 標記所有通知為已讀
	MarkAllAsRead(ctx context.Context, userID uint, userType string) error
}
