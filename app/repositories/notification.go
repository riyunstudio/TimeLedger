package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
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

	var query *gorm.DB

	if userType == "ADMIN" {
		// 管理員通知：user_id = 0，透過 center_id 過濾
		// 需要先取得 admin 的 center_id
		var centerID uint
		adminRepo := NewAdminUserRepository(r.app)
		admin, err := adminRepo.GetByID(ctx, userID)
		if err == nil {
			centerID = admin.CenterID
		}

		if centerID > 0 {
			query = r.app.MySQL.RDB.WithContext(ctx).
				Where("user_type = ? AND center_id = ?", "ADMIN", centerID).
				Order("created_at DESC")
		} else {
			query = r.app.MySQL.RDB.WithContext(ctx).
				Where("user_type = ? AND user_id = ?", "ADMIN", 0).
				Order("created_at DESC")
		}
	} else {
		// 老師通知：透過 user_id 和 user_type 過濾
		query = r.app.MySQL.RDB.WithContext(ctx).
			Where("user_id = ? AND user_type = ?", userID, userType).
			Order("created_at DESC")
	}

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

	if userType == "ADMIN" {
		// 管理員通知：user_id = 0，透過 center_id 過濾
		var centerID uint
		adminRepo := NewAdminUserRepository(r.app)
		admin, err := adminRepo.GetByID(ctx, userID)
		if err == nil {
			centerID = admin.CenterID
		}

		if centerID > 0 {
			err := r.app.MySQL.RDB.WithContext(ctx).
				Where("user_type = ? AND center_id = ? AND is_read = false", "ADMIN", centerID).
				Order("created_at DESC").
				Find(&notifications).Error
			return notifications, err
		}
	}

	// 老師通知或 fallback（找不到 admin 的 centerID）
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

	if userType == "ADMIN" {
		// 管理員通知：透過 center_id 過濾
		var centerID uint
		adminRepo := NewAdminUserRepository(r.app)
		admin, err := adminRepo.GetByID(ctx, userID)
		if err == nil {
			centerID = admin.CenterID
		}

		if centerID > 0 {
			return r.app.MySQL.WDB.WithContext(ctx).
				Model(&models.Notification{}).
				Where("user_type = ? AND center_id = ?", "ADMIN", centerID).
				Updates(map[string]interface{}{
					"is_read": true,
					"read_at": &now,
				}).Error
		}
	}

	// 老師通知或 fallback
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
