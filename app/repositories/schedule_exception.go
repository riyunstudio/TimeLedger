package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type ScheduleExceptionRepository struct {
	BaseRepository
	app   *app.App
	model *models.ScheduleException
}

func NewScheduleExceptionRepository(app *app.App) *ScheduleExceptionRepository {
	return &ScheduleExceptionRepository{
		app: app,
	}
}

func (rp *ScheduleExceptionRepository) GetByID(ctx context.Context, id uint) (models.ScheduleException, error) {
	var data models.ScheduleException
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *ScheduleExceptionRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.ScheduleException, error) {
	var data models.ScheduleException
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *ScheduleExceptionRepository) GetByRuleAndDate(ctx context.Context, ruleID uint, date time.Time) ([]models.ScheduleException, error) {
	var data []models.ScheduleException
	// 將 date 轉換為日期字串進行比較
	dateStr := date.Format("2006-01-02")
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("rule_id = ? AND DATE(original_date) = ?", ruleID, dateStr).Find(&data).Error
	return data, err
}

func (rp *ScheduleExceptionRepository) GetByRuleIDAndDateStr(ctx context.Context, ruleID uint, dateStr string) ([]models.ScheduleException, error) {
	var data []models.ScheduleException
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("rule_id = ? AND DATE(original_date) = ?", ruleID, dateStr).Find(&data).Error
	return data, err
}

func (rp *ScheduleExceptionRepository) GetByRuleAndDateAndCenterID(ctx context.Context, ruleID uint, centerID uint, date time.Time) ([]models.ScheduleException, error) {
	var data []models.ScheduleException
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("rule_id = ? AND center_id = ? AND original_date = ?", ruleID, centerID, date).Find(&data).Error
	return data, err
}

func (rp *ScheduleExceptionRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.ScheduleException, error) {
	var data []models.ScheduleException
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *ScheduleExceptionRepository) Create(ctx context.Context, data models.ScheduleException) (models.ScheduleException, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *ScheduleExceptionRepository) Update(ctx context.Context, data models.ScheduleException) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *ScheduleExceptionRepository) UpdateByIDAndCenterID(ctx context.Context, id uint, centerID uint, data models.ScheduleException) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).Save(&data).Error
}

func (rp *ScheduleExceptionRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Delete(&models.ScheduleException{}, id).Error
}

func (rp *ScheduleExceptionRepository) DeleteByIDAndCenterID(ctx context.Context, id uint, centerID uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).Delete(&models.ScheduleException{}).Error
}
