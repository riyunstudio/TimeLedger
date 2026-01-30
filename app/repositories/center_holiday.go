package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

type CenterHolidayRepository struct {
	GenericRepository[models.CenterHoliday]
	app *app.App
}

func NewCenterHolidayRepository(app *app.App) *CenterHolidayRepository {
	return &CenterHolidayRepository{
		GenericRepository: NewGenericRepository[models.CenterHoliday](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// Transaction executes a function within a database transaction.
// This method creates a NEW CenterHolidayRepository instance with transaction connections.
func (r *CenterHolidayRepository) Transaction(ctx context.Context, fn func(txRepo *CenterHolidayRepository) error) error {
	return r.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &CenterHolidayRepository{
			GenericRepository: NewTransactionRepo[models.CenterHoliday](ctx, tx, r.table),
			app:               r.app,
		}
		return fn(txRepo)
	})
}

func (r *CenterHolidayRepository) GetByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (models.CenterHoliday, error) {
	return r.First(ctx, "center_id = ? AND date = ?", centerID, date.Format("2006-01-02"))
}

func (r *CenterHolidayRepository) Delete(ctx context.Context, id uint) error {
	return r.DeleteByID(ctx, id)
}

func (r *CenterHolidayRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterHoliday, error) {
	return r.FindWithCenterScope(ctx, centerID)
}

func (r *CenterHolidayRepository) ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.CenterHoliday, error) {
	return r.Find(ctx, "center_id = ? AND date >= ? AND date <= ?", centerID, start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (r *CenterHolidayRepository) ExistsByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (bool, error) {
	return r.Exists(ctx, "center_id = ? AND date = ?", centerID, date.Format("2006-01-02"))
}

func (r *CenterHolidayRepository) BulkCreate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, error) {
	if len(holidays) == 0 {
		return holidays, nil
	}
	return r.CreateBatch(ctx, holidays)
}

func (r *CenterHolidayRepository) BulkCreateWithSkipDuplicate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, int64, error) {
	if len(holidays) == 0 {
		return holidays, 0, nil
	}

	var created int64
	var createdHolidays []models.CenterHoliday

	err := r.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range holidays {
			var existing int64
			tx.Model(&models.CenterHoliday{}).
				Where("center_id = ? AND date = ?", holidays[i].CenterID, holidays[i].Date.Format("2006-01-02")).
				Count(&existing)

			if existing == 0 {
				if err := tx.Create(&holidays[i]).Error; err != nil {
					return err
				}
				createdHolidays = append(createdHolidays, holidays[i])
				created++
			}
		}
		return nil
	})

	return createdHolidays, created, err
}

type CenterHolidayRepositoryInterface interface {
	Create(ctx context.Context, holiday models.CenterHoliday) (models.CenterHoliday, error)
	GetByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (models.CenterHoliday, error)
	ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterHoliday, error)
	ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.CenterHoliday, error)
	ExistsByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (bool, error)
	BulkCreate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, error)
	BulkCreateWithSkipDuplicate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, int64, error)
}

var _ CenterHolidayRepositoryInterface = (*CenterHolidayRepository)(nil)
