package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type SessionNoteRepository struct {
	GenericRepository[models.SessionNote]
	app *app.App
}

func NewSessionNoteRepository(app *app.App) *SessionNoteRepository {
	return &SessionNoteRepository{
		GenericRepository: NewGenericRepository[models.SessionNote](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *SessionNoteRepository) GetByRuleAndDate(ctx context.Context, centerID uint, ruleID uint, sessionDate time.Time) (models.SessionNote, error) {
	var data models.SessionNote
	err := rp.app.MySQL.WDB.WithContext(ctx).
		Where("rule_id = ? AND session_date = ?", ruleID, sessionDate).
		Where("rule_id IN (SELECT id FROM schedule_rules WHERE center_id = ?)", centerID).
		First(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) GetByTeacherAndDate(ctx context.Context, centerID uint, teacherID uint, sessionDate time.Time) (models.SessionNote, error) {
	var data models.SessionNote
	err := rp.app.MySQL.WDB.WithContext(ctx).
		Where("teacher_id = ? AND session_date = ?", teacherID, sessionDate).
		Where("teacher_id IN (SELECT teacher_id FROM center_memberships WHERE center_id = ? AND status = 'ACTIVE')", centerID).
		First(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) GetOrCreate(ctx context.Context, centerID uint, teacherID uint, ruleID uint, sessionDate time.Time) (models.SessionNote, bool, error) {
	var data models.SessionNote
	normalizedDate := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 0, 0, 0, 0, sessionDate.Location())

	query := rp.app.MySQL.WDB.WithContext(ctx).
		Where("teacher_id = ? AND rule_id = ? AND session_date = ?", teacherID, ruleID, normalizedDate.Format("2006-01-02")).
		Where("rule_id IN (SELECT id FROM schedule_rules WHERE center_id = ?)", centerID).
		Where("teacher_id IN (SELECT teacher_id FROM center_memberships WHERE center_id = ? AND status = 'ACTIVE')", centerID)

	err := query.First(&data).Error

	if err == nil {
		return data, false, nil
	}

	if data.ID == 0 {
		newData := models.SessionNote{
			TeacherID:   teacherID,
			RuleID:      &ruleID,
			SessionDate: normalizedDate,
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
