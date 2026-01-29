package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterRepository struct {
	GenericRepository[models.Center]
	app *app.App
}

func NewCenterRepository(app *app.App) *CenterRepository {
	return &CenterRepository{
		GenericRepository: NewGenericRepository[models.Center](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// List retrieves all centers (deprecated, use Find instead).
// Kept for backwards compatibility with existing code.
func (rp *CenterRepository) List(ctx context.Context) ([]models.Center, error) {
	return rp.Find(ctx)
}

// GetByName retrieves a center by its name (exact match).
func (rp *CenterRepository) GetByName(ctx context.Context, name string) (models.Center, error) {
	var data models.Center
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("name = ?", name).
		First(&data).Error
	return data, err
}

// GetByPlanLevel retrieves all centers with a specific plan level.
func (rp *CenterRepository) GetByPlanLevel(ctx context.Context, planLevel string) ([]models.Center, error) {
	var data []models.Center
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("plan_level = ?", planLevel).
		Find(&data).Error
	return data, err
}

// ListActive retrieves all centers ordered by creation date.
func (rp *CenterRepository) ListActive(ctx context.Context) ([]models.Center, error) {
	var data []models.Center
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Order("created_at DESC").
		Find(&data).Error
	return data, err
}

// ExistsByName checks if a center with the given name already exists.
func (rp *CenterRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return rp.Exists(ctx, "name = ?", name)
}

// CountByPlanLevel counts centers by plan level.
func (rp *CenterRepository) CountByPlanLevel(ctx context.Context, planLevel string) (int64, error) {
	return rp.Count(ctx, "plan_level = ?", planLevel)
}

// UpdateSettings updates only the settings JSON field for a center.
func (rp *CenterRepository) UpdateSettings(ctx context.Context, id uint, settings models.CenterSettings) error {
	return rp.UpdateFields(ctx, id, map[string]interface{}{
		"settings": settings,
	})
}

// GetSettings retrieves the parsed settings for a center.
func (rp *CenterRepository) GetSettings(ctx context.Context, id uint) (models.CenterSettings, error) {
	center, err := rp.GetByID(ctx, id)
	if err != nil {
		return models.CenterSettings{}, err
	}
	return center.Settings, nil
}
