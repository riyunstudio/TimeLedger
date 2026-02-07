package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

// CenterTermRepository 學期資料存取層
type CenterTermRepository struct {
	GenericRepository[models.CenterTerm]
	app *app.App
}

// NewCenterTermRepository 建立 CenterTermRepository 實例
func NewCenterTermRepository(app *app.App) *CenterTermRepository {
	return &CenterTermRepository{
		GenericRepository: NewGenericRepository[models.CenterTerm](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// Transaction 執行資料庫交易
func (r *CenterTermRepository) Transaction(ctx context.Context, fn func(txRepo *CenterTermRepository) error) error {
	return r.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &CenterTermRepository{
			GenericRepository: NewTransactionRepo[models.CenterTerm](ctx, tx, r.table),
			app:               r.app,
		}
		return fn(txRepo)
	})
}

// GetByIDWithCenterScope 根據 ID 和中心 ID 查詢（帶權限檢查）
func (r *CenterTermRepository) GetByIDWithCenterScope(ctx context.Context, id, centerID uint) (models.CenterTerm, error) {
	return r.First(ctx, "id = ? AND center_id = ?", id, centerID)
}

// ListByCenterID 根據中心 ID 取得所有學期
func (r *CenterTermRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterTerm, error) {
	return r.FindWithCenterScope(ctx, centerID)
}

// ListActiveByCenterID 根據中心 ID 取得所有進行中的學期
func (r *CenterTermRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.CenterTerm, error) {
	now := time.Now()
	return r.Find(ctx, "center_id = ? AND start_date <= ? AND end_date >= ?", centerID, now, now)
}

// ExistsByCenterAndName 檢查學期名稱是否已存在
func (r *CenterTermRepository) ExistsByCenterAndName(ctx context.Context, centerID uint, name string) (bool, error) {
	return r.Exists(ctx, "center_id = ? AND name = ?", centerID, name)
}

// ExistsByCenterAndNameExcludingID 檢查學期名稱是否已存在（排除指定 ID）
func (r *CenterTermRepository) ExistsByCenterAndNameExcludingID(ctx context.Context, centerID uint, name string, excludeID uint) (bool, error) {
	return r.Exists(ctx, "center_id = ? AND name = ? AND id != ?", centerID, name, excludeID)
}

// ExistsByCenterAndDateRange 檢查日期範圍是否與現有學期重疊
func (r *CenterTermRepository) ExistsByCenterAndDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time, excludeID uint) (bool, error) {
	query := r.dbWrite.WithContext(ctx).
		Model(&models.CenterTerm{}).
		Where("center_id = ?", centerID).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate)

	// 排除更新的學期
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Delete 軟刪除學期
func (r *CenterTermRepository) Delete(ctx context.Context, id uint) error {
	return r.DeleteByID(ctx, id)
}

// Update 更新學期
func (r *CenterTermRepository) Update(ctx context.Context, term *models.CenterTerm) error {
	return r.GenericRepository.Update(ctx, *term)
}

// CenterTermRepositoryInterface 學期儲存庫介面
type CenterTermRepositoryInterface interface {
	Create(ctx context.Context, term models.CenterTerm) (models.CenterTerm, error)
	GetByIDWithCenterScope(ctx context.Context, id, centerID uint) (models.CenterTerm, error)
	ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterTerm, error)
	ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.CenterTerm, error)
	ExistsByCenterAndName(ctx context.Context, centerID uint, name string) (bool, error)
	ExistsByCenterAndNameExcludingID(ctx context.Context, centerID uint, name string, excludeID uint) (bool, error)
	ExistsByCenterAndDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time, excludeID uint) (bool, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, term *models.CenterTerm) error
}

var _ CenterTermRepositoryInterface = (*CenterTermRepository)(nil)
