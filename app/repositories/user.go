package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type UserRepository struct {
	GenericRepository[models.User]
	app *app.App
}

func NewUserRepository(app *app.App) *UserRepository {
	return &UserRepository{
		GenericRepository: NewGenericRepository[models.User](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *UserRepository) Get(ctx context.Context, cond models.User) (data models.User, err error) {
	query := rp.app.MySQL.RDB.WithContext(ctx).Model(&data)
	if cond.ID != 0 {
		query.Where("id = ?", cond.ID)
	}
	if cond.Name != "" {
		query.Where("name = ?", cond.Name)
	}
	err = query.Find(&data).Error
	return
}
