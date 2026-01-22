package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type PersonalEventRepository struct {
	BaseRepository
	app *app.App
}

func NewPersonalEventRepository(app *app.App) *PersonalEventRepository {
	return &PersonalEventRepository{app: app}
}

func (r *PersonalEventRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.PersonalEvent, error) {
	var events []models.PersonalEvent
	err := r.app.Mysql.RDB.WithContext(ctx).Where("teacher_id = ?", teacherID).Order("start_at ASC").Find(&events).Error
	return events, err
}

func (r *PersonalEventRepository) GetByID(ctx context.Context, id uint) (*models.PersonalEvent, error) {
	var event models.PersonalEvent
	err := r.app.Mysql.RDB.WithContext(ctx).First(&event, id).Error
	return &event, err
}

func (r *PersonalEventRepository) Create(ctx context.Context, event *models.PersonalEvent) error {
	return r.app.Mysql.WDB.WithContext(ctx).Create(event).Error
}

func (r *PersonalEventRepository) Update(ctx context.Context, event *models.PersonalEvent) error {
	return r.app.Mysql.WDB.WithContext(ctx).Save(event).Error
}

func (r *PersonalEventRepository) Delete(ctx context.Context, id uint) error {
	return r.app.Mysql.WDB.WithContext(ctx).Delete(&models.PersonalEvent{}, id).Error
}

func (r *PersonalEventRepository) GetByTeacherAndDateRange(ctx context.Context, teacherID uint, startDate, endDate time.Time) ([]models.PersonalEvent, error) {
	var events []models.PersonalEvent
	err := r.app.Mysql.RDB.WithContext(ctx).
		Where("teacher_id = ? AND start_at >= ? AND start_at < ?", teacherID, startDate, endDate).
		Order("start_at ASC").
		Find(&events).Error
	return events, err
}

type UpdateEventRequest struct {
	Title    *string
	StartAt  *time.Time
	EndAt    *time.Time
	IsAllDay *bool
	ColorHex *string
}

func (r *PersonalEventRepository) UpdateFutureOccurrences(ctx context.Context, eventID uint, teacherID uint, req UpdateEventRequest, updatedAt time.Time) (int64, error) {
	var event models.PersonalEvent
	if err := r.app.Mysql.RDB.WithContext(ctx).First(&event, eventID).Error; err != nil {
		return 0, err
	}

	updates := map[string]interface{}{
		"updated_at": updatedAt,
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.ColorHex != nil {
		updates["color_hex"] = *req.ColorHex
	}
	if req.IsAllDay != nil {
		updates["is_all_day"] = *req.IsAllDay
	}
	if req.StartAt != nil {
		updates["start_at"] = *req.StartAt
	}
	if req.EndAt != nil {
		updates["end_at"] = *req.EndAt
	}

	result := r.app.Mysql.WDB.WithContext(ctx).
		Model(&models.PersonalEvent{}).
		Where("id = ? AND teacher_id = ? AND start_at >= ?", eventID, teacherID, event.StartAt).
		Updates(updates)

	return result.RowsAffected, result.Error
}

func (r *PersonalEventRepository) UpdateAllOccurrences(ctx context.Context, eventID uint, teacherID uint, req UpdateEventRequest, updatedAt time.Time) (int64, error) {
	updates := map[string]interface{}{
		"updated_at": updatedAt,
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.ColorHex != nil {
		updates["color_hex"] = *req.ColorHex
	}
	if req.IsAllDay != nil {
		updates["is_all_day"] = *req.IsAllDay
	}
	if req.StartAt != nil {
		updates["start_at"] = *req.StartAt
	}
	if req.EndAt != nil {
		updates["end_at"] = *req.EndAt
	}

	result := r.app.Mysql.WDB.WithContext(ctx).
		Model(&models.PersonalEvent{}).
		Where("id = ? AND teacher_id = ?", eventID, teacherID).
		Updates(updates)

	return result.RowsAffected, result.Error
}
