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

// BatchGetByCenterAndTeachers 批量取得中心對多個教師的備註（效能優化：減少 N+1 查詢）
func (r *CenterTeacherNoteRepository) BatchGetByCenterAndTeachers(ctx context.Context, centerID uint, teacherIDs []uint) (map[uint]models.CenterTeacherNote, error) {
	if len(teacherIDs) == 0 {
		return make(map[uint]models.CenterTeacherNote), nil
	}

	var notes []models.CenterTeacherNote
	err := r.app.MySQL.RDB.WithContext(ctx).Where("center_id = ? AND teacher_id IN ?", centerID, teacherIDs).Find(&notes).Error
	if err != nil {
		return nil, err
	}

	// 按教師 ID 分組
	result := make(map[uint]models.CenterTeacherNote, len(notes))
	for _, note := range notes {
		result[note.TeacherID] = note
	}
	return result, nil
}
