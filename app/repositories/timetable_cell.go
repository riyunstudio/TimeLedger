package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TimetableCellRepository struct {
	BaseRepository
	app   *app.App
	model *models.TimetableCell
}

func NewTimetableCellRepository(app *app.App) *TimetableCellRepository {
	return &TimetableCellRepository{
		app: app,
	}
}

func (rp *TimetableCellRepository) GetByID(ctx context.Context, id uint) (models.TimetableCell, error) {
	var data models.TimetableCell
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *TimetableCellRepository) ListByTemplateID(ctx context.Context, templateID uint) ([]models.TimetableCell, error) {
	var data []models.TimetableCell
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("template_id = ?", templateID).Find(&data).Error
	return data, err
}

func (rp *TimetableCellRepository) Create(ctx context.Context, data models.TimetableCell) (models.TimetableCell, error) {
	err := rp.app.Mysql.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *TimetableCellRepository) Update(ctx context.Context, data models.TimetableCell) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *TimetableCellRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Delete(&models.TimetableCell{}, id).Error
}

func (rp *TimetableCellRepository) DeleteByTemplateID(ctx context.Context, templateID uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Where("template_id = ?", templateID).Delete(&models.TimetableCell{}).Error
}
