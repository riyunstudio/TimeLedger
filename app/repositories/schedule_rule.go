package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type ScheduleRuleRepository struct {
	BaseRepository
	app   *app.App
	model *models.ScheduleRule
}

func NewScheduleRuleRepository(app *app.App) *ScheduleRuleRepository {
	return &ScheduleRuleRepository{
		app: app,
	}
}

func (rp *ScheduleRuleRepository) GetByID(ctx context.Context, id uint) (models.ScheduleRule, error) {
	var data models.ScheduleRule
	err := rp.app.Mysql.RDB.WithContext(ctx).Preload("Exceptions").Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.ScheduleRule, error) {
	var data models.ScheduleRule
	err := rp.app.Mysql.RDB.WithContext(ctx).Preload("Exceptions").
		Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) ListByTeacherID(ctx context.Context, teacherID uint, centerID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.Mysql.RDB.WithContext(ctx).
		Where("teacher_id = ? AND center_id = ?", teacherID, centerID).
		Find(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) ListByRoomID(ctx context.Context, roomID uint, centerID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.Mysql.RDB.WithContext(ctx).
		Where("room_id = ? AND center_id = ?", roomID, centerID).
		Find(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) Create(ctx context.Context, data models.ScheduleRule) (models.ScheduleRule, error) {
	err := rp.app.Mysql.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) Update(ctx context.Context, data models.ScheduleRule) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *ScheduleRuleRepository) UpdateByIDAndCenterID(ctx context.Context, id uint, centerID uint, data models.ScheduleRule) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).Save(&data).Error
}

func (rp *ScheduleRuleRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Delete(&models.ScheduleRule{}, id).Error
}

func (rp *ScheduleRuleRepository) DeleteByIDAndCenterID(ctx context.Context, id uint, centerID uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).Delete(&models.ScheduleRule{}).Error
}
