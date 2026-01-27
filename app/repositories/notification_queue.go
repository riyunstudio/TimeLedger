package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

type NotificationQueueRepository struct {
	BaseRepository
	app *app.App
}

func NewNotificationQueueRepository(app *app.App) *NotificationQueueRepository {
	return &NotificationQueueRepository{app: app}
}

// Create 建立通知佇列項目
func (r *NotificationQueueRepository) Create(ctx context.Context, item *models.NotificationQueue) error {
	return r.app.MySQL.WDB.WithContext(ctx).Create(item).Error
}

// GetByID 依 ID 取得通知佇列項目
func (r *NotificationQueueRepository) GetByID(ctx context.Context, id uint) (*models.NotificationQueue, error) {
	var item models.NotificationQueue
	err := r.app.MySQL.RDB.WithContext(ctx).First(&item, id).Error
	return &item, err
}

// GetPending 取得待處理的佇列項目
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

// UpdateStatus 更新佇列狀態
func (r *NotificationQueueRepository) UpdateStatus(ctx context.Context, id uint, updateData map[string]interface{}) error {
	return r.app.MySQL.WDB.WithContext(ctx).
		Model(&models.NotificationQueue{}).
		Where("id = ?", id).
		Updates(updateData).Error
}

// IncrementRetryCount 增加重試次數
func (r *NotificationQueueRepository) IncrementRetryCount(ctx context.Context, id uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).
		Model(&models.NotificationQueue{}).
		Where("id = ?", id).
		Update("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

// GetByRecipient 依收件者取得佇列項目
func (r *NotificationQueueRepository) GetByRecipient(ctx context.Context, recipientID uint, recipientType string, limit int, offset int) ([]models.NotificationQueue, error) {
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

// DeleteOldEntries 刪除過期的佇列項目（保留 7 天）
func (r *NotificationQueueRepository) DeleteOldEntries(ctx context.Context) error {
	cutoff := time.Now().AddDate(0, 0, -7)
	return r.app.MySQL.WDB.WithContext(ctx).
		Where("created_at < ?", cutoff).
		Where("status IN ?", []string{models.NotificationStatusSent, models.NotificationStatusFailed}).
		Delete(&models.NotificationQueue{}).Error
}

// GetStats 取得佇列統計
func (r *NotificationQueueRepository) GetStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	var pending int64
	if err := r.app.MySQL.RDB.WithContext(ctx).
		Model(&models.NotificationQueue{}).
		Where("status = ?", models.NotificationStatusPending).
		Count(&pending).Error; err != nil {
		return nil, err
	}
	stats["pending"] = pending

	var sent int64
	if err := r.app.MySQL.RDB.WithContext(ctx).
		Model(&models.NotificationQueue{}).
		Where("status = ?", models.NotificationStatusSent).
		Count(&sent).Error; err != nil {
		return nil, err
	}
	stats["sent"] = sent

	var failed int64
	if err := r.app.MySQL.RDB.WithContext(ctx).
		Model(&models.NotificationQueue{}).
		Where("status = ?", models.NotificationStatusFailed).
		Count(&failed).Error; err != nil {
		return nil, err
	}
	stats["failed"] = failed

	return stats, nil
}
