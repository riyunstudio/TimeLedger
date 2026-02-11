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
