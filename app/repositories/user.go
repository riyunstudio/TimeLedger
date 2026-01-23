package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type UserRepository struct {
	BaseRepository
	app   *app.App
	model *models.User
}

func NewUserRepository(app *app.App) *UserRepository {
	return &UserRepository{
		app: app,
	}
}

func (rp *UserRepository) Get(ctx context.Context, cond models.User) (data models.User, err error) {
	query := rp.app.MySQL.RDB.WithContext(ctx).Model(&rp.model)
	if cond.ID != 0 {
		query.Where("id = ?", cond.ID)
	}
	if cond.Name != "" {
		query.Where("name = ?", cond.Name)
	}
	err = query.Find(&data).Error
	return
}

func (rp *UserRepository) Create(ctx context.Context, data models.User) (models.User, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Model(&rp.model).Create(&data).Error
	return data, err
}

func (rp *UserRepository) UpdateById(ctx context.Context, data models.User) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Model(&rp.model).Where("id", data.ID).Updates(&data).Error
}
