package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterTeacherNoteRepository struct {
	GenericRepository[models.CenterTeacherNote]
	app *app.App
}

func NewCenterTeacherNoteRepository(app *app.App) *CenterTeacherNoteRepository {
	// 如果 app 為 nil，返回一個基礎結構（用於測試）
	if app == nil {
		return &CenterTeacherNoteRepository{}
	}
	return &CenterTeacherNoteRepository{
		GenericRepository: NewGenericRepository[models.CenterTeacherNote](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *CenterTeacherNoteRepository) GetByCenterAndTeacher(ctx context.Context, centerID uint, teacherID uint) (models.CenterTeacherNote, error) {
	return r.FirstWithCenterScope(ctx, centerID, "teacher_id = ?", teacherID)
}

func (r *CenterTeacherNoteRepository) BatchGetByCenterAndTeachers(ctx context.Context, centerID uint, teacherIDs []uint) (map[uint]models.CenterTeacherNote, error) {
	if len(teacherIDs) == 0 {
		return make(map[uint]models.CenterTeacherNote), nil
	}

	notes, err := r.FindWithCenterScope(ctx, centerID, "teacher_id IN ?", teacherIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]models.CenterTeacherNote, len(notes))
	for _, note := range notes {
		result[note.TeacherID] = note
	}
	return result, nil
}
