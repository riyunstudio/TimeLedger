package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type AdminUserRepository struct {
	BaseRepository
	app   *app.App
	model *models.AdminUser
}

func NewAdminUserRepository(app *app.App) *AdminUserRepository {
	return &AdminUserRepository{
		app: app,
	}
}

func (rp *AdminUserRepository) GetByID(ctx context.Context, id uint) (models.AdminUser, error) {
	var data models.AdminUser
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *AdminUserRepository) GetByEmail(ctx context.Context, email string) (models.AdminUser, error) {
	var data models.AdminUser
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("email = ?", email).First(&data).Error
	return data, err
}

func (rp *AdminUserRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.AdminUser, error) {
	var data []models.AdminUser
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *AdminUserRepository) Create(ctx context.Context, data models.AdminUser) (models.AdminUser, error) {
	err := rp.app.Mysql.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *AdminUserRepository) Update(ctx context.Context, data models.AdminUser) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *AdminUserRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Delete(&models.AdminUser{}, id).Error
}
