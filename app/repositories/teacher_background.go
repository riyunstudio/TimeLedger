package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TeacherBackgroundRepository struct {
	GenericRepository[models.TeacherBackground]
	app *app.App
}

func NewTeacherBackgroundRepository(app *app.App) *TeacherBackgroundRepository {
	return &TeacherBackgroundRepository{
		GenericRepository: NewGenericRepository[models.TeacherBackground](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *TeacherBackgroundRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.TeacherBackground, error) {
	return r.Find(ctx, "teacher_id = ?", teacherID)
}

func (r *TeacherBackgroundRepository) GetByTeacherID(ctx context.Context, teacherID uint) ([]models.TeacherBackground, error) {
	var backgrounds []models.TeacherBackground
	err := r.app.MySQL.RDB.WithContext(ctx).
		Where("teacher_id = ?", teacherID).
		Order("created_at DESC").
		Find(&backgrounds).Error
	return backgrounds, err
}

func (r *TeacherBackgroundRepository) DeleteByTeacherID(ctx context.Context, teacherID uint) error {
	return r.app.MySQL.RDB.WithContext(ctx).
		Where("teacher_id = ?", teacherID).
		Delete(&models.TeacherBackground{}).Error
}

func (r *TeacherBackgroundRepository) Delete(ctx context.Context, id uint) error {
	return r.app.MySQL.RDB.WithContext(ctx).
		Delete(&models.TeacherBackground{}, id).Error
}
