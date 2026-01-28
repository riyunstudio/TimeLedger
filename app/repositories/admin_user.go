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
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *AdminUserRepository) GetByEmail(ctx context.Context, email string) (models.AdminUser, error) {
	var data models.AdminUser
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("email = ?", email).First(&data).Error
	return data, err
}

func (rp *AdminUserRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.AdminUser, error) {
	var data []models.AdminUser
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

// GetByCenterID 依中心 ID 取得所有管理員（返回指標）
func (rp *AdminUserRepository) GetByCenterID(ctx context.Context, centerID uint) ([]*models.AdminUser, error) {
	var data []*models.AdminUser
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ? AND status = ?", centerID, "ACTIVE").
		Find(&data).Error
	return data, err
}

// GetByIDPtr 依 ID 取得管理員（返回指標）
func (rp *AdminUserRepository) GetByIDPtr(ctx context.Context, id uint) (*models.AdminUser, error) {
	var data models.AdminUser
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return &data, err
}

// GetByLineUserID 依 LINE User ID 取得管理員
func (rp *AdminUserRepository) GetByLineUserID(ctx context.Context, lineUserID string) (*models.AdminUser, error) {
	var data models.AdminUser
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("line_user_id = ?", lineUserID).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (rp *AdminUserRepository) Create(ctx context.Context, data models.AdminUser) (models.AdminUser, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *AdminUserRepository) Update(ctx context.Context, data models.AdminUser) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

// UpdateFields 更新特定欄位
func (rp *AdminUserRepository) UpdateFields(ctx context.Context, id uint, fields map[string]interface{}) error {
	return rp.app.MySQL.WDB.WithContext(ctx).
		Model(&models.AdminUser{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (rp *AdminUserRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Delete(&models.AdminUser{}, id).Error
}
