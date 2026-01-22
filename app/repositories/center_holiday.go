package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

type CenterHolidayRepository struct {
	app *app.App
}

func NewCenterHolidayRepository(app *app.App) *CenterHolidayRepository {
	return &CenterHolidayRepository{app: app}
}

func (r *CenterHolidayRepository) Create(ctx context.Context, holiday models.CenterHoliday) (models.CenterHoliday, error) {
	err := r.app.Mysql.WDB.WithContext(ctx).Create(&holiday).Error
	return holiday, err
}

func (r *CenterHolidayRepository) GetByID(ctx context.Context, id uint) (models.CenterHoliday, error) {
	var holiday models.CenterHoliday
	err := r.app.Mysql.RDB.WithContext(ctx).First(&holiday, id).Error
	return holiday, err
}

func (r *CenterHolidayRepository) GetByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (models.CenterHoliday, error) {
	var holiday models.CenterHoliday
	err := r.app.Mysql.RDB.WithContext(ctx).
		Where("center_id = ? AND date = ?", centerID, date.Format("2006-01-02")).
		First(&holiday).Error
	return holiday, err
}

func (r *CenterHolidayRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterHoliday, error) {
	var holidays []models.CenterHoliday
	err := r.app.Mysql.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Order("date ASC").
		Find(&holidays).Error
	return holidays, err
}

func (r *CenterHolidayRepository) ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.CenterHoliday, error) {
	var holidays []models.CenterHoliday
	err := r.app.Mysql.RDB.WithContext(ctx).
		Where("center_id = ? AND date >= ? AND date <= ?", centerID, start.Format("2006-01-02"), end.Format("2006-01-02")).
		Order("date ASC").
		Find(&holidays).Error
	return holidays, err
}

func (r *CenterHolidayRepository) Delete(ctx context.Context, id uint) error {
	return r.app.Mysql.WDB.WithContext(ctx).Delete(&models.CenterHoliday{}, id).Error
}

func (r *CenterHolidayRepository) ExistsByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (bool, error) {
	var count int64
	err := r.app.Mysql.RDB.WithContext(ctx).
		Model(&models.CenterHoliday{}).
		Where("center_id = ? AND date = ?", centerID, date.Format("2006-01-02")).
		Count(&count).Error
	return count > 0, err
}

func (r *CenterHolidayRepository) BulkCreate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, error) {
	if len(holidays) == 0 {
		return holidays, nil
	}
	err := r.app.Mysql.WDB.WithContext(ctx).Create(&holidays).Error
	return holidays, err
}

func (r *CenterHolidayRepository) BulkCreateWithSkipDuplicate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, int64, error) {
	if len(holidays) == 0 {
		return holidays, 0, nil
	}

	var created int64
	var createdHolidays []models.CenterHoliday

	err := r.app.Mysql.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
	GetByID(ctx context.Context, id uint) (models.CenterHoliday, error)
	GetByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (models.CenterHoliday, error)
	ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterHoliday, error)
	ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.CenterHoliday, error)
	Delete(ctx context.Context, id uint) error
	ExistsByCenterAndDate(ctx context.Context, centerID uint, date time.Time) (bool, error)
	BulkCreate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, error)
	BulkCreateWithSkipDuplicate(ctx context.Context, holidays []models.CenterHoliday) ([]models.CenterHoliday, int64, error)
}

var _ CenterHolidayRepositoryInterface = (*CenterHolidayRepository)(nil)
