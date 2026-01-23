package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterTeacherNoteRepository struct {
	BaseRepository
	app *app.App
}

func NewCenterTeacherNoteRepository(app *app.App) *CenterTeacherNoteRepository {
	return &CenterTeacherNoteRepository{app: app}
}

func (r *CenterTeacherNoteRepository) GetByCenterAndTeacher(ctx context.Context, centerID uint, teacherID uint) (models.CenterTeacherNote, error) {
	var note models.CenterTeacherNote
	err := r.app.MySQL.RDB.WithContext(ctx).Where("center_id = ? AND teacher_id = ?", centerID, teacherID).First(&note).Error
	return note, err
}

func (r *CenterTeacherNoteRepository) Create(ctx context.Context, note *models.CenterTeacherNote) error {
	return r.app.MySQL.WDB.WithContext(ctx).Create(note).Error
}

func (r *CenterTeacherNoteRepository) Update(ctx context.Context, note *models.CenterTeacherNote) error {
	return r.app.MySQL.WDB.WithContext(ctx).Save(note).Error
}

func (r *CenterTeacherNoteRepository) Delete(ctx context.Context, id uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).Delete(&models.CenterTeacherNote{}, id).Error
}

func (r *CenterTeacherNoteRepository) GetByID(ctx context.Context, id uint) (*models.CenterTeacherNote, error) {
	var note models.CenterTeacherNote
	err := r.app.MySQL.RDB.WithContext(ctx).First(&note, id).Error
	return &note, err
}
