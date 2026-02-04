package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

type AdminLoginHistoryRepository struct {
	GenericRepository[models.AdminLoginHistory]
	app *app.App
}

func NewAdminLoginHistoryRepository(app *app.App) *AdminLoginHistoryRepository {
	return &AdminLoginHistoryRepository{
		GenericRepository: NewGenericRepository[models.AdminLoginHistory](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// Transaction executes a function within a database transaction.
func (rp *AdminLoginHistoryRepository) Transaction(ctx context.Context, fn func(txRepo *AdminLoginHistoryRepository) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &AdminLoginHistoryRepository{
			GenericRepository: NewTransactionRepo[models.AdminLoginHistory](ctx, tx, rp.table),
			app:               rp.app,
		}
		return fn(txRepo)
	})
}

// GetByAdminID 取得指定管理員的登入紀錄
func (rp *AdminLoginHistoryRepository) GetByAdminID(ctx context.Context, adminID uint, limit int) ([]models.AdminLoginHistory, error) {
	if limit <= 0 {
		limit = 50
	}
	histories, err := rp.Find(ctx, "admin_id = ?", adminID)
	if err != nil {
		return nil, err
	}
	if len(histories) > limit {
		return histories[:limit], nil
	}
	return histories, nil
}

// GetRecentFailedLogins 取得最近的失敗登入紀錄
func (rp *AdminLoginHistoryRepository) GetRecentFailedLogins(ctx context.Context, adminID uint, since time.Time) ([]models.AdminLoginHistory, error) {
	histories, err := rp.Find(ctx, "admin_id = ? AND status = ? AND created_at >= ?",
		adminID, models.LoginStatusFailed, since)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// CountFailedLoginsSince 計算指定時間範圍內的失敗登入次數
func (rp *AdminLoginHistoryRepository) CountFailedLoginsSince(ctx context.Context, adminID uint, since time.Time) int64 {
	var count int64
	rp.dbRead.WithContext(ctx).Model(&models.AdminLoginHistory{}).
		Where("admin_id = ? AND status = ? AND created_at >= ?",
			adminID, models.LoginStatusFailed, since).
		Count(&count)
	return count
}
