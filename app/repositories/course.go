package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CourseRepository struct {
	GenericRepository[models.Course]
	app *app.App
}

func NewCourseRepository(app *app.App) *CourseRepository {
	return &CourseRepository{
		GenericRepository: NewGenericRepository[models.Course](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *CourseRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.Course, error) {
	return rp.FindWithCenterScope(ctx, centerID)
}

func (rp *CourseRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.Course, error) {
	return rp.FindWithCenterScope(ctx, centerID, "is_active = ?", true)
}

func (rp *CourseRepository) ToggleActive(ctx context.Context, id uint, centerID uint, isActive bool) error {
	return rp.UpdateFieldsWithCenterScope(ctx, id, centerID, map[string]interface{}{
		"is_active": isActive,
	})
}
