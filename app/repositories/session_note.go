package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type SessionNoteRepository struct {
	BaseRepository
	app   *app.App
	model *models.SessionNote
}

func NewSessionNoteRepository(app *app.App) *SessionNoteRepository {
	return &SessionNoteRepository{
		app: app,
	}
}

func (rp *SessionNoteRepository) GetByRuleAndDate(ctx context.Context, ruleID uint, sessionDate time.Time) (models.SessionNote, error) {
	var data models.SessionNote
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("rule_id = ? AND session_date = ?", ruleID, sessionDate).First(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) GetByTeacherAndDate(ctx context.Context, teacherID uint, sessionDate time.Time) (models.SessionNote, error) {
	var data models.SessionNote
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("teacher_id = ? AND session_date = ?", teacherID, sessionDate).First(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) Create(ctx context.Context, data models.SessionNote) (models.SessionNote, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) Update(ctx context.Context, data models.SessionNote) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *SessionNoteRepository) GetOrCreate(ctx context.Context, teacherID uint, ruleID uint, sessionDate time.Time) (models.SessionNote, bool, error) {
	var data models.SessionNote
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("teacher_id = ? AND rule_id = ? AND session_date = ?", teacherID, ruleID, sessionDate).First(&data).Error

	if err == nil {
		return data, false, nil
	}

	if data.ID == 0 {
		newData := models.SessionNote{
			TeacherID:   teacherID,
			RuleID:      &ruleID,
			SessionDate: sessionDate,
			Content:     "",
			PrepNote:    "",
		}
		if createErr := rp.app.MySQL.WDB.WithContext(ctx).Create(&newData).Error; createErr != nil {
			return newData, true, createErr
		}
		return newData, true, nil
	}

	return data, false, err
}
