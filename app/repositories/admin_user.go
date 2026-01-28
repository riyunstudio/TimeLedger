package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"

	"golang.org/x/crypto/bcrypt"
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

// GetByCenterID 依中心 ID 取得所有管理員（返回指標，不過濾狀態）
func (rp *AdminUserRepository) GetByCenterID(ctx context.Context, centerID uint) ([]*models.AdminUser, error) {
	var data []*models.AdminUser
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
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

// VerifyPassword 驗證密碼
func (rp *AdminUserRepository) VerifyPassword(ctx context.Context, email string, password string) bool {
	var data models.AdminUser
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("email = ?", email).First(&data).Error
	if err != nil {
		return false
	}
	// 使用bcrypt驗證密碼
	return rp.checkPassword(data.PasswordHash, password)
}

// HashPassword 產生密碼哈希
func (rp *AdminUserRepository) HashPassword(password string) (string, error) {
	// 使用bcrypt加密密碼
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkPassword 檢查密碼是否匹配
func (rp *AdminUserRepository) checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
