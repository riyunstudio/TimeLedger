package repositories

import (
	"akali/app"
	"akali/app/models"
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

func (rp *UserRepository) Get(cond models.User) (data models.User, err error) {
	query := rp.app.Mysql.RDB.Model(&rp.model)
	if cond.ID != 0 {
		query.Where("id = ?", cond.ID)
	}
	err = query.Find(&data).Error
	return
}

func (rp *UserRepository) Create(data models.User) (models.User, error) {
	err := rp.app.Mysql.WDB.Model(&rp.model).Create(&data).Error
	return data, err
}

func (rp *UserRepository) UpdateById(data models.User) error {
	return rp.app.Mysql.WDB.Model(&rp.model).Where("id", data.ID).Updates(&data).Error
}
