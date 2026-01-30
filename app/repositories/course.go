package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
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

// Transaction executes a function within a database transaction.
// This method creates a NEW CourseRepository instance with transaction connections
// to avoid race conditions in concurrent requests.
//
// Usage Example:
//
//	result, err := rp.Transaction(ctx, func(txRepo *CourseRepository) error {
//	    // All operations using txRepo will be within the same transaction
//	    // Custom methods like ListByCenterID are available
//	    if _, err := txRepo.Create(ctx, course1); err != nil {
//	        return err
//	    }
//	    if _, err := txRepo.Create(ctx, course2); err != nil {
//	        return err
//	    }
//	    return nil
//	})
func (rp *CourseRepository) Transaction(ctx context.Context, fn func(txRepo *CourseRepository) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new CourseRepository instance with transaction connections
		txRepo := &CourseRepository{
			GenericRepository: GenericRepository[models.Course]{
				dbRead:  tx.WithContext(ctx),
				dbWrite: tx.WithContext(ctx),
				table:   rp.table,
			},
			app: rp.app,
		}
		return fn(txRepo)
	})
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
