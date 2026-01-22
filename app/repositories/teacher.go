package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TeacherRepository struct {
	BaseRepository
	app   *app.App
	model *models.Teacher
}

func NewTeacherRepository(app *app.App) *TeacherRepository {
	return &TeacherRepository{
		app: app,
	}
}

func (rp *TeacherRepository) GetByID(ctx context.Context, id uint) (models.Teacher, error) {
	var data models.Teacher
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *TeacherRepository) GetByLineUserID(ctx context.Context, lineUserID string) (models.Teacher, error) {
	var data models.Teacher
	err := rp.app.Mysql.RDB.WithContext(ctx).Where("line_user_id = ?", lineUserID).First(&data).Error
	return data, err
}

func (rp *TeacherRepository) List(ctx context.Context) ([]models.Teacher, error) {
	var data []models.Teacher
	err := rp.app.Mysql.RDB.WithContext(ctx).Find(&data).Error
	return data, err
}

func (rp *TeacherRepository) Create(ctx context.Context, data models.Teacher) (models.Teacher, error) {
	err := rp.app.Mysql.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *TeacherRepository) Update(ctx context.Context, data models.Teacher) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *TeacherRepository) DeleteByID(ctx context.Context, id uint) error {
	return rp.app.Mysql.WDB.WithContext(ctx).Where("id = ?", id).Delete(&models.Teacher{}).Error
}

func (rp *TeacherRepository) ListPersonalHashtags(ctx context.Context, teacherID uint) ([]models.Hashtag, error) {
	var hashtags []models.Hashtag
	err := rp.app.Mysql.RDB.WithContext(ctx).
		Table("teacher_personal_hashtags").
		Select("hashtags.*").
		Joins("INNER JOIN hashtags ON teacher_personal_hashtags.hashtag_id = hashtags.id").
		Where("teacher_personal_hashtags.teacher_id = ?", teacherID).
		Order("teacher_personal_hashtags.sort_order ASC").
		Find(&hashtags).Error
	return hashtags, err
}
