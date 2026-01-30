package repositories

import (
	"context"
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type PersonalEventRepository struct {
	GenericRepository[models.PersonalEvent]
	app *app.App
}

func NewPersonalEventRepository(app *app.App) *PersonalEventRepository {
	return &PersonalEventRepository{
		GenericRepository: NewGenericRepository[models.PersonalEvent](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *PersonalEventRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.PersonalEvent, error) {
	return r.Find(ctx, "teacher_id = ?", teacherID)
}

func (r *PersonalEventRepository) UpdateNote(ctx context.Context, eventID uint, content string) error {
	return r.UpdateFields(ctx, eventID, map[string]interface{}{
		"note": content,
	})
}

func (r *PersonalEventRepository) GetByTeacherAndDateRange(ctx context.Context, teacherID uint, startDate, endDate time.Time) ([]models.PersonalEvent, error) {
	return r.Find(ctx, "teacher_id = ? AND start_at >= ? AND start_at < ?", teacherID, startDate, endDate)
}

func (r *PersonalEventRepository) CheckPersonalEventConflict(ctx context.Context, teacherID uint, weekday int, startTime string, endTime string, checkDate time.Time) ([]models.PersonalEvent, error) {
	events, err := r.ListByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	var conflicts []models.PersonalEvent
	isCrossDay := endTime < startTime

	for _, event := range events {
		if event.RecurrenceRule.Type == "" || event.RecurrenceRule.Type == "NONE" {
			eventDate := event.StartAt
			if eventDate.Format("2006-01-02") == checkDate.Format("2006-01-02") {
				if TimesOverlapCrossDay(formatTimeForComparison(event.StartAt), formatTimeForComparison(event.EndAt), isCrossDay, startTime, endTime, isCrossDay) {
					conflicts = append(conflicts, event)
				}
			}
			continue
		}

		if shouldIncludeRecurringEvent(event, weekday, checkDate) {
			eventIsCrossDay := formatTimeForComparison(event.EndAt) < formatTimeForComparison(event.StartAt)
			if TimesOverlapCrossDay(formatTimeForComparison(event.StartAt), formatTimeForComparison(event.EndAt), eventIsCrossDay, startTime, endTime, isCrossDay) {
				conflicts = append(conflicts, event)
			}
		}
	}

	return conflicts, nil
}

func (r *PersonalEventRepository) CheckPersonalEventConflictCrossDay(ctx context.Context, teacherID uint, weekday int, startTime string, endTime string, checkDate time.Time) ([]models.PersonalEvent, error) {
	events, err := r.ListByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	var conflicts []models.PersonalEvent

	for _, event := range events {
		eventWeekday := int(event.StartAt.Weekday())
		if eventWeekday == 0 {
			eventWeekday = 7
		}

		eventDate := event.StartAt.Format("2006-01-02")
		checkDateStr := checkDate.Format("2006-01-02")
		nextDateStr := checkDate.AddDate(0, 0, 1).Format("2006-01-02")

		isSameDay := eventDate == checkDateStr
		isNextDay := eventDate == nextDateStr

		if !isSameDay && !isNextDay {
			continue
		}

		eventIsCrossDay := formatTimeForComparison(event.EndAt) < formatTimeForComparison(event.StartAt)

		if isSameDay {
			eventStartMin := ParseTimeToMinutes(formatTimeForComparison(event.StartAt))
			if eventIsCrossDay || eventStartMin >= 20*60 {
				if TimesOverlapCrossDay(startTime, endTime, true, formatTimeForComparison(event.StartAt), formatTimeForComparison(event.EndAt), eventIsCrossDay) {
					conflicts = append(conflicts, event)
				}
			}
		} else if isNextDay {
			eventEndMin := ParseTimeToMinutes(formatTimeForComparison(event.EndAt))
			if eventIsCrossDay || eventEndMin <= 8*60 {
				if TimesOverlapCrossDay(startTime, endTime, true, formatTimeForComparison(event.StartAt), formatTimeForComparison(event.EndAt), eventIsCrossDay) {
					conflicts = append(conflicts, event)
				}
			}
		}
	}

	return conflicts, nil
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

func formatTimeForComparison(t time.Time) string {
	return t.Format("15:04")
}

func ParseTimeToMinutes(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return 0
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	return hour*60 + minute
}

func TimesOverlapCrossDay(start1, end1 string, isCrossDay1 bool, start2, end2 string, isCrossDay2 bool) bool {
	start1Min := ParseTimeToMinutes(start1)
	end1Min := ParseTimeToMinutes(end1)
	start2Min := ParseTimeToMinutes(start2)
	end2Min := ParseTimeToMinutes(end2)

	if isCrossDay1 && end1Min < start1Min {
		end1Min += 24 * 60
	}
	if isCrossDay2 {
		if end2Min < start2Min {
			end2Min += 24 * 60
		} else {
			start2Min += 24 * 60
			end2Min += 24 * 60
		}
	}

	return start1Min < end2Min && end1Min > start2Min
}

func TimesOverlapCrossDayWithNextDay(start1, end1 string, isCrossDay1 bool, start2, end2 string) bool {
	start1Min := ParseTimeToMinutes(start1)
	end1Min := ParseTimeToMinutes(end1)
	start2Min := ParseTimeToMinutes(start2)
	end2Min := ParseTimeToMinutes(end2)

	if isCrossDay1 {
		end1Min += 24 * 60
	}

	start2Min += 24 * 60
	end2Min += 24 * 60

	return start1Min < end2Min && end1Min > start2Min
}

func shouldIncludeRecurringEvent(event models.PersonalEvent, weekday int, date time.Time) bool {
	rule := event.RecurrenceRule

	if rule.Until != nil {
		untilDate, err := time.Parse("2006-01-02", *rule.Until)
		if err == nil && date.After(untilDate) {
			return false
		}
	}

	if rule.Count != nil && *rule.Count > 0 {
		// Simplified weekday check for count-based rules
	}

	eventWeekday := int(date.Weekday())
	if eventWeekday == 0 {
		eventWeekday = 7
	}

	for _, w := range rule.Weekdays {
		if w == eventWeekday || (w == 0 && eventWeekday == 7) || (w == 7 && eventWeekday == 0) {
			return true
		}
	}

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
