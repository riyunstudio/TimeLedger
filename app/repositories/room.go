package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type RoomRepository struct {
	BaseRepository
	app   *app.App
	model *models.Room
}

func NewRoomRepository(app *app.App) *RoomRepository {
	return &RoomRepository{
		app: app,
	}
}

func (rp *RoomRepository) GetByID(ctx context.Context, id uint) (models.Room, error) {
	var data models.Room
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *RoomRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.Room, error) {
	var data models.Room
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *RoomRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.Room, error) {
	var data []models.Room
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *RoomRepository) Create(ctx context.Context, data models.Room) (models.Room, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *RoomRepository) Update(ctx context.Context, data models.Room) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *RoomRepository) DeleteByID(ctx context.Context, id uint, centerID uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).
		Where("id = ? AND center_id = ?", id, centerID).
		Delete(&models.Room{}).Error
}
