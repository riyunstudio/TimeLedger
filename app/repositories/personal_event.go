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
	err := r.app.MySQL.RDB.WithContext(ctx).Where("teacher_id = ?", teacherID).Order("start_at ASC").Find(&events).Error
	return events, err
}

func (r *PersonalEventRepository) GetByID(ctx context.Context, id uint) (*models.PersonalEvent, error) {
	var event models.PersonalEvent
	err := r.app.MySQL.RDB.WithContext(ctx).First(&event, id).Error
	return &event, err
}

func (r *PersonalEventRepository) Create(ctx context.Context, event *models.PersonalEvent) error {
	return r.app.MySQL.WDB.WithContext(ctx).Create(event).Error
}

func (r *PersonalEventRepository) Update(ctx context.Context, event *models.PersonalEvent) error {
	return r.app.MySQL.WDB.WithContext(ctx).Save(event).Error
}

func (r *PersonalEventRepository) Delete(ctx context.Context, id uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).Delete(&models.PersonalEvent{}, id).Error
}

func (r *PersonalEventRepository) GetByTeacherAndDateRange(ctx context.Context, teacherID uint, startDate, endDate time.Time) ([]models.PersonalEvent, error) {
	var events []models.PersonalEvent
	err := r.app.MySQL.RDB.WithContext(ctx).
		Where("teacher_id = ? AND start_at >= ? AND start_at < ?", teacherID, startDate, endDate).
		Order("start_at ASC").
		Find(&events).Error
	return events, err
}

// CheckPersonalEventConflict checks if there's a personal event conflict for a teacher at a specific weekday and time
func (r *PersonalEventRepository) CheckPersonalEventConflict(ctx context.Context, teacherID uint, weekday int, startTime string, endTime string, checkDate time.Time) ([]models.PersonalEvent, error) {
	// Get all personal events for this teacher
	events, err := r.ListByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	var conflicts []models.PersonalEvent
	
	for _, event := range events {
		// Skip non-recurring events (they should be checked against the specific date)
		if event.RecurrenceRule.Type == "" || event.RecurrenceRule.Type == "NONE" {
			// Check if the single event is on the same day and overlaps
			eventDate := event.StartAt
			if eventDate.Format("2006-01-02") == checkDate.Format("2006-01-02") {
				eventStartHour := event.StartAt.Hour()
				eventEndHour := event.EndAt.Hour()
				
				ruleStartHour := parseHour(startTime)
				ruleEndHour := parseHour(endTime)
				
				// Check time overlap
				if eventStartHour < ruleEndHour && eventEndHour > ruleStartHour {
					conflicts = append(conflicts, event)
				}
			}
			continue
		}

		// Handle recurring events
		if shouldIncludeRecurringEvent(event, weekday, checkDate) {
			eventStartHour := event.StartAt.Hour()
			eventEndHour := event.EndAt.Hour()
			
			ruleStartHour := parseHour(startTime)
			ruleEndHour := parseHour(endTime)
			
			// Check time overlap
			if eventStartHour < ruleEndHour && eventEndHour > ruleStartHour {
				conflicts = append(conflicts, event)
			}
		}
	}

	return conflicts, nil
}

// parseHour extracts hour from time string "HH:MM"
func parseHour(timeStr string) int {
	var hour int
	for _, c := range timeStr {
		if c >= '0' && c <= '9' {
			hour = hour*10 + int(c-'0')
			if hour >= 24 {
				break
			}
		}
	}
	return hour
}

// shouldIncludeRecurringEvent checks if a recurring event should be included for the given weekday and date
func shouldIncludeRecurringEvent(event models.PersonalEvent, weekday int, date time.Time) bool {
	rule := event.RecurrenceRule
	
	// Check if the date is within the valid range
	if rule.Until != nil {
		untilDate, err := time.Parse("2006-01-02", *rule.Until)
		if err == nil && date.After(untilDate) {
			return false
		}
	}

	if rule.Count != nil && *rule.Count > 0 {
		// This requires counting occurrences, which is complex
		// For simplicity, we'll check the weekday match
	}

	// Check weekday match
	eventWeekday := int(date.Weekday())
	if eventWeekday == 0 {
		eventWeekday = 7 // Convert Sunday from 0 to 7 for consistency
	}
	
	for _, w := range rule.Weekdays {
		if w == eventWeekday || (w == 0 && eventWeekday == 7) || (w == 7 && eventWeekday == 0) {
			return true
		}
	}
	
	// Also check the original event's weekday
	if eventWeekday == int(date.Weekday()) {
		return true
	}
	
	return false
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
	if err := r.app.MySQL.RDB.WithContext(ctx).First(&event, eventID).Error; err != nil {
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

	result := r.app.MySQL.WDB.WithContext(ctx).
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

	result := r.app.MySQL.WDB.WithContext(ctx).
		Model(&models.PersonalEvent{}).
		Where("id = ? AND teacher_id = ?", eventID, teacherID).
		Updates(updates)

	return result.RowsAffected, result.Error
}
