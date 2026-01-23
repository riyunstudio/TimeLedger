package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type OfferingRepository struct {
	BaseRepository
	app   *app.App
	model *models.Offering
}

func NewOfferingRepository(app *app.App) *OfferingRepository {
	return &OfferingRepository{
		app: app,
	}
}

func (rp *OfferingRepository) GetByID(ctx context.Context, id uint) (models.Offering, error) {
	var data models.Offering
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *OfferingRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.Offering, error) {
	var data models.Offering
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *OfferingRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.Offering, error) {
	var data []models.Offering
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *OfferingRepository) ListByCenterIDPaginated(ctx context.Context, centerID uint, page, limit int64) ([]models.Offering, int64, error) {
	var data []models.Offering
	var total int64

	rp.app.MySQL.RDB.WithContext(ctx).Model(&models.Offering{}).
		Where("center_id = ?", centerID).
		Count(&total)

	offset := (page - 1) * limit
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Order("id ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&data).Error

	return data, total, err
}

func (rp *OfferingRepository) Create(ctx context.Context, data models.Offering) (models.Offering, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *OfferingRepository) Update(ctx context.Context, data models.Offering) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *OfferingRepository) DeleteByID(ctx context.Context, id uint, centerID uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).
		Where("id = ? AND center_id = ?", id, centerID).
		Delete(&models.Offering{}).Error
}

func (rp *OfferingRepository) Copy(ctx context.Context, original models.Offering, newName string) (models.Offering, error) {
	newOffering := models.Offering{
		CenterID:            original.CenterID,
		CourseID:            original.CourseID,
		Name:                newName,
		DefaultRoomID:       original.DefaultRoomID,
		DefaultTeacherID:    original.DefaultTeacherID,
		AllowBufferOverride: original.AllowBufferOverride,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&newOffering).Error
	return newOffering, err
}
