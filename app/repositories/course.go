package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CourseRepository struct {
	BaseRepository
	app   *app.App
	model *models.Course
}

func NewCourseRepository(app *app.App) *CourseRepository {
	return &CourseRepository{
		app: app,
	}
}

func (rp *CourseRepository) GetByID(ctx context.Context, id uint) (models.Course, error) {
	var data models.Course
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *CourseRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.Course, error) {
	var data models.Course
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *CourseRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.Course, error) {
	var data []models.Course
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *CourseRepository) Create(ctx context.Context, data models.Course) (models.Course, error) {
	err := rp.app.Mysql.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *CourseRepository) Update(ctx context.Context, data models.Course) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *CourseRepository) DeleteByID(ctx context.Context, id uint, centerID uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).
		Where("id = ? AND center_id = ?", id, centerID).
		Delete(&models.Course{}).Error
}
