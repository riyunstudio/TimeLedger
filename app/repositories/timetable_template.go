package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TimetableTemplateRepository struct {
	BaseRepository
	app   *app.App
	model *models.TimetableTemplate
}

func NewTimetableTemplateRepository(app *app.App) *TimetableTemplateRepository {
	return &TimetableTemplateRepository{
		app: app,
	}
}

func (rp *TimetableTemplateRepository) GetByID(ctx context.Context, id uint) (models.TimetableTemplate, error) {
	var data models.TimetableTemplate
	err := rp.app.Mysql.RDB.WithContext(ctx).Preload("Cells").Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *TimetableTemplateRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.TimetableTemplate, error) {
	var data models.TimetableTemplate
	err := rp.app.Mysql.RDB.WithContext(ctx).Preload("Cells").
		Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *TimetableTemplateRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.TimetableTemplate, error) {
	var data []models.TimetableTemplate
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *TimetableTemplateRepository) Create(ctx context.Context, data models.TimetableTemplate) (models.TimetableTemplate, error) {
	err := rp.app.Mysql.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *TimetableTemplateRepository) Update(ctx context.Context, data models.TimetableTemplate) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *TimetableTemplateRepository) UpdateByIDAndCenterID(ctx context.Context, id uint, centerID uint, data models.TimetableTemplate) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Model(&models.TimetableTemplate{}).Where("id = ? AND center_id = ?", id, centerID).Updates(map[string]interface{}{
		"name":       data.Name,
		"row_type":   data.RowType,
		"is_active":  data.IsActive,
		"updated_at": data.UpdatedAt,
	}).Error
}

func (rp *TimetableTemplateRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Delete(&models.TimetableTemplate{}, id).Error
}

func (rp *TimetableTemplateRepository) DeleteByIDAndCenterID(ctx context.Context, id uint, centerID uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).Delete(&models.TimetableTemplate{}).Error
}
