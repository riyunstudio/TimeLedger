package repositories

import (
	"context"
	"fmt"
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
	// 使用 WDB 查詢，確保能讀取到剛寫入的資料
	err := rp.app.MySQL.WDB.WithContext(ctx).Where("rule_id = ? AND session_date = ?", ruleID, sessionDate).First(&data).Error
	return data, err
}

func (rp *SessionNoteRepository) GetByTeacherAndDate(ctx context.Context, teacherID uint, sessionDate time.Time) (models.SessionNote, error) {
	var data models.SessionNote
	// 使用 WDB 查詢，確保能讀取到剛寫入的資料
	err := rp.app.MySQL.WDB.WithContext(ctx).Where("teacher_id = ? AND session_date = ?", teacherID, sessionDate).First(&data).Error
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

	// 標準化日期：確保時間部分為 00:00:00
	normalizedDate := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 0, 0, 0, 0, sessionDate.Location())

	// 使用 WDB 查詢，確保能讀取到剛寫入的資料
	query := rp.app.MySQL.WDB.WithContext(ctx).
		Where("teacher_id = ? AND rule_id = ? AND session_date = ?", teacherID, ruleID, normalizedDate.Format("2006-01-02"))

	// 打印調試日誌
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
