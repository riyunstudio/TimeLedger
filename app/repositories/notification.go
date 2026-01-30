package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type NotificationRepository struct {
	GenericRepository[models.Notification]
	app *app.App
}

func NewNotificationRepository(app *app.App) *NotificationRepository {
	return &NotificationRepository{
		GenericRepository: NewGenericRepository[models.Notification](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *NotificationRepository) ListUnread(ctx context.Context, userID uint, userType string) ([]models.Notification, error) {
	return r.Find(ctx, "user_id = ? AND user_type = ? AND is_read = ?", userID, userType, false)
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID uint) error {
	now := time.Now()
	return r.UpdateFields(ctx, notificationID, map[string]interface{}{
		"is_read": true,
		"read_at": &now,
	})
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uint, userType string) error {
	now := time.Now()
	result := r.app.MySQL.WDB.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})
	return result.Error
}
