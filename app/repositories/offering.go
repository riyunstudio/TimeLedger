package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

type OfferingRepository struct {
	GenericRepository[models.Offering]
	app *app.App
}

func NewOfferingRepository(app *app.App) *OfferingRepository {
	return &OfferingRepository{
		GenericRepository: NewGenericRepository[models.Offering](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// Transaction executes a function within a database transaction.
// This method creates a NEW OfferingRepository instance with transaction connections
// to avoid race conditions in concurrent requests.
//
// Usage Example:
//
//	result, err := rp.Transaction(ctx, func(txRepo *OfferingRepository) error {
//	    // All operations using txRepo will be within the same transaction
//	    // Custom methods like Copy and ListActiveByCenterID are available
//	    if _, err := txRepo.Create(ctx, offering1); err != nil {
//	        return err
//	    }
//	    if _, err := txRepo.Create(ctx, offering2); err != nil {
//	        return err
//	    }
//	    return nil
//	})
func (rp *OfferingRepository) Transaction(ctx context.Context, fn func(txRepo *OfferingRepository) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new OfferingRepository instance with transaction connections
		txRepo := &OfferingRepository{
			GenericRepository: GenericRepository[models.Offering]{
				dbRead:  tx.WithContext(ctx),
				dbWrite: tx.WithContext(ctx),
				table:   rp.table,
			},
			app: rp.app,
		}
		return fn(txRepo)
	})
}

func (rp *OfferingRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.Offering, error) {
	var offerings []models.Offering
	err := rp.dbRead.WithContext(ctx).
		Preload("Course").
		Where("center_id = ? AND is_active = ?", centerID, true).
		Order("created_at DESC").
		Find(&offerings).Error
	return offerings, err
}

func (rp *OfferingRepository) ListByCenterIDPaginated(ctx context.Context, centerID uint, page, limit int) ([]models.Offering, int64, error) {
	return rp.FindPaged(ctx, page, limit, "created_at DESC", "center_id = ?", centerID)
}

// SearchByNamePaginated 搜尋班別（分頁 + 關鍵字）
func (rp *OfferingRepository) SearchByNamePaginated(ctx context.Context, centerID uint, query string, page, limit int) ([]models.Offering, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}

	offset := (page - 1) * limit

	var offerings []models.Offering
	var total int64

	baseQuery := rp.dbRead.WithContext(ctx).Model(&models.Offering{}).Where("center_id = ?", centerID)

	// 如果有搜尋關鍵字，加入模糊比對
	if query != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+query+"%")
	}

	// 取得總數
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 取得分頁資料
	if err := baseQuery.
		Preload("Course").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&offerings).Error; err != nil {
		return nil, 0, err
	}

	return offerings, total, nil
}

func (rp *OfferingRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.Offering, error) {
	return rp.GetByIDWithCenterScope(ctx, id, centerID)
}

func (rp *OfferingRepository) ToggleActive(ctx context.Context, id uint, centerID uint, isActive bool) error {
	return rp.UpdateFieldsWithCenterScope(ctx, id, centerID, map[string]interface{}{
		"is_active": isActive,
	})
}

func (rp *OfferingRepository) Copy(ctx context.Context, original models.Offering, newName string) (models.Offering, error) {
	newOffering := models.Offering{
		CenterID:            original.CenterID,
		CourseID:            original.CourseID,
		Name:                newName,
		DefaultRoomID:       original.DefaultRoomID,
		DefaultTeacherID:    original.DefaultTeacherID,
		AllowBufferOverride: original.AllowBufferOverride,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	err := rp.dbWrite.WithContext(ctx).Create(&newOffering).Error
	return newOffering, err
}
