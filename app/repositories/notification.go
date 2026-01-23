package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type NotificationRepository struct {
	BaseRepository
	app *app.App
}

func NewNotificationRepository(app *app.App) *NotificationRepository {
	return &NotificationRepository{app: app}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.app.MySQL.WDB.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.app.MySQL.RDB.WithContext(ctx).First(&notification, id).Error
	return &notification, err
}

func (r *NotificationRepository) List(ctx context.Context, userID uint, userType string, limit int, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	query := r.app.MySQL.RDB.WithContext(ctx).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) ListUnread(ctx context.Context, userID uint, userType string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.app.MySQL.RDB.WithContext(ctx).
		Where("user_id = ? AND user_type = ? AND is_read = false", userID, userType).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID uint) error {
	now := time.Now()
	return r.app.MySQL.WDB.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uint, userType string) error {
	now := time.Now()
	return r.app.MySQL.WDB.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *NotificationRepository) Delete(ctx context.Context, id uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).Delete(&models.Notification{}, id).Error
}
