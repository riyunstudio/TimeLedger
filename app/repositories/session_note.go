package repositories

import (
	"context"
	"fmt"
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

func (rp *SessionNoteRepository) GetByRuleAndDate(ctx context.Context, ruleID uint, sessionDate time.Time) (models.SessionNote, error) {
	var data models.SessionNote
	err := rp.app.MySQL.WDB.WithContext(ctx).Where("rule_id = ? AND session_date = ?", ruleID, sessionDate).First(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) GetByTeacherAndDate(ctx context.Context, teacherID uint, sessionDate time.Time) (models.SessionNote, error) {
	var data models.SessionNote
	err := rp.app.MySQL.WDB.WithContext(ctx).Where("teacher_id = ? AND session_date = ?", teacherID, sessionDate).First(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) GetOrCreate(ctx context.Context, teacherID uint, ruleID uint, sessionDate time.Time) (models.SessionNote, bool, error) {
	var data models.SessionNote
	normalizedDate := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 0, 0, 0, 0, sessionDate.Location())

	query := rp.app.MySQL.WDB.WithContext(ctx).
		Where("teacher_id = ? AND rule_id = ? AND session_date = ?", teacherID, ruleID, normalizedDate.Format("2006-01-02"))

	fmt.Printf("[SessionNote] GetOrCreate - teacherID: %d, ruleID: %d, sessionDate: %s, normalizedDate: %s\n",
		teacherID, ruleID, sessionDate.Format("2006-01-02"), normalizedDate.Format("2006-01-02"))

	err := query.First(&data).Error

	if err == nil {
		fmt.Printf("[SessionNote] Found existing note - id: %d, content: %s\n", data.ID, data.Content)
		return data, false, nil
	}

	fmt.Printf("[SessionNote] Not found, creating new note - error: %v\n", err)

	if data.ID == 0 {
		newData := models.SessionNote{
			TeacherID:   teacherID,
			RuleID:      &ruleID,
			SessionDate: normalizedDate,
			Content:     "",
			PrepNote:    "",
		}
		if createErr := rp.app.MySQL.WDB.WithContext(ctx).Create(&newData).Error; createErr != nil {
			fmt.Printf("[SessionNote] Failed to create note - error: %v\n", createErr)
			return newData, true, createErr
		}
		fmt.Printf("[SessionNote] Created new note - id: %d\n", newData.ID)
		return newData, true, nil
	}

	return data, false, err
}
