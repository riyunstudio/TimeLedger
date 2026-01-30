package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"

	"golang.org/x/crypto/bcrypt"
)

type AdminUserRepository struct {
	GenericRepository[models.AdminUser]
	app *app.App
}

func NewAdminUserRepository(app *app.App) *AdminUserRepository {
	return &AdminUserRepository{
		GenericRepository: NewGenericRepository[models.AdminUser](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *AdminUserRepository) GetByEmail(ctx context.Context, email string) (models.AdminUser, error) {
	return rp.First(ctx, "email = ?", email)
}

func (rp *AdminUserRepository) GetByLineUserID(ctx context.Context, lineUserID string) (*models.AdminUser, error) {
	data, err := rp.First(ctx, "line_user_id = ?", lineUserID)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (rp *AdminUserRepository) GetByCenterID(ctx context.Context, centerID uint) ([]models.AdminUser, error) {
	return rp.FindWithCenterScope(ctx, centerID)
}

func (rp *AdminUserRepository) GetByIDPtr(ctx context.Context, id uint) (*models.AdminUser, error) {
	data, err := rp.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (rp *AdminUserRepository) VerifyPassword(ctx context.Context, email string, password string) bool {
	data, err := rp.GetByEmail(ctx, email)
	if err != nil {
		return false
	}
	return rp.checkPassword(data.PasswordHash, password)
}

func (rp *AdminUserRepository) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (rp *AdminUserRepository) checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
