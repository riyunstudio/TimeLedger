package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

type NotificationQueueRepository struct {
	GenericRepository[models.NotificationQueue]
	app *app.App
}

func NewNotificationQueueRepository(app *app.App) *NotificationQueueRepository {
	return &NotificationQueueRepository{
		GenericRepository: NewGenericRepository[models.NotificationQueue](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *NotificationQueueRepository) GetPending(ctx context.Context, limit int) ([]models.NotificationQueue, error) {
	var items []models.NotificationQueue
	err := r.app.MySQL.RDB.WithContext(ctx).
		Where("status = ?", models.NotificationStatusPending).
		Where("scheduled_at <= ?", time.Now()).
		Order("created_at ASC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

func (r *NotificationQueueRepository) IncrementRetryCount(ctx context.Context, id uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).
		Model(&models.NotificationQueue{}).
		Where("id = ?", id).
		Update("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

func (r *NotificationQueueRepository) DeleteOldEntries(ctx context.Context) error {
	cutoff := time.Now().AddDate(0, 0, -7)
	_, err := r.DeleteWhere(ctx, "created_at < ?", cutoff, models.NotificationStatusSent, models.NotificationStatusFailed)
	return err
}

func (r *NotificationQueueRepository) GetByRecipient(ctx context.Context, recipientID uint, recipientType string, limit, offset int) ([]models.NotificationQueue, error) {
	var items []models.NotificationQueue
	query := r.app.MySQL.RDB.WithContext(ctx).
		Where("recipient_id = ? AND recipient_type = ?", recipientID, recipientType).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&items).Error
	return items, err
}
