package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterRepository struct {
	BaseRepository
	app   *app.App
	model *models.Center
}

func NewCenterRepository(app *app.App) *CenterRepository {
	return &CenterRepository{
		app: app,
	}
}

func (rp *CenterRepository) GetByID(ctx context.Context, id uint) (models.Center, error) {
	var data models.Center
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *CenterRepository) List(ctx context.Context) ([]models.Center, error) {
	var data []models.Center
	err := rp.app.MySQL.RDB.WithContext(ctx).Find(&data).Error
	return data, err
}

func (rp *CenterRepository) Create(ctx context.Context, data models.Center) (models.Center, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *CenterRepository) Update(ctx context.Context, data models.Center) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}
