package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
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

func (rp *OfferingRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.Offering, error) {
	return rp.FindWithCenterScope(ctx, centerID, "is_active = ?", true)
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
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&newOffering).Error
	return newOffering, err
}

func (rp *OfferingRepository) ListByCenterIDPaginated(ctx context.Context, centerID uint, page, limit int) ([]models.Offering, int64, error) {
	return rp.FindPaged(ctx, page, limit, "created_at DESC", "center_id = ?", centerID)
}

func (rp *OfferingRepository) DeleteByID(ctx context.Context, id, centerID uint) error {
	return rp.DeleteByIDWithCenterScope(ctx, id, centerID)
}

func (rp *OfferingRepository) GetByIDAndCenterID(ctx context.Context, id, centerID uint) (models.Offering, error) {
	return rp.GetByIDWithCenterScope(ctx, id, centerID)
}
