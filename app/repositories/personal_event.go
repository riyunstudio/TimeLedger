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
	isCrossDay := endTime < startTime

	for _, event := range events {
		// Skip non-recurring events (they should be checked against the specific date)
		if event.RecurrenceRule.Type == "" || event.RecurrenceRule.Type == "NONE" {
			// Check if the single event is on the same day and overlaps
			eventDate := event.StartAt
			if eventDate.Format("2006-01-02") == checkDate.Format("2006-01-02") {
				// Use proper time comparison with minutes
				if TimesOverlapCrossDay(formatTimeForComparison(event.StartAt), formatTimeForComparison(event.EndAt), isCrossDay, startTime, endTime, isCrossDay) {
					conflicts = append(conflicts, event)
				}
			}
			continue
		}

		// Handle recurring events
		if shouldIncludeRecurringEvent(event, weekday, checkDate) {
			// Use proper time comparison with minutes
			eventIsCrossDay := formatTimeForComparison(event.EndAt) < formatTimeForComparison(event.StartAt)
			if TimesOverlapCrossDay(formatTimeForComparison(event.StartAt), formatTimeForComparison(event.EndAt), eventIsCrossDay, startTime, endTime, isCrossDay) {
				conflicts = append(conflicts, event)
			}
		}
	}

	return conflicts, nil
}

// CheckPersonalEventConflictCrossDay 檢查個人行程是否與跨日課程衝突
func (r *PersonalEventRepository) CheckPersonalEventConflictCrossDay(ctx context.Context, teacherID uint, weekday int, startTime string, endTime string, checkDate time.Time) ([]models.PersonalEvent, error) {
	// Get all personal events for this teacher
	events, err := r.ListByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	var conflicts []models.PersonalEvent

	for _, event := range events {
		// 計算事件的星期幾
		eventWeekday := int(event.StartAt.Weekday())
		if eventWeekday == 0 {
			eventWeekday = 7
		}

		// 檢查事件是否在當天或隔天
		eventDate := event.StartAt.Format("2006-01-02")
		checkDateStr := checkDate.Format("2006-01-02")
		nextDateStr := checkDate.AddDate(0, 0, 1).Format("2006-01-02")

		isSameDay := eventDate == checkDateStr
		isNextDay := eventDate == nextDateStr

		if !isSameDay && !isNextDay {
			continue
		}

		// 判斷事件是否為跨日
		eventIsCrossDay := formatTimeForComparison(event.EndAt) < formatTimeForComparison(event.StartAt)

		if isSameDay {
			// 當天：檢查晚上時段（20:00 以後）
			eventStartMin := ParseTimeToMinutes(formatTimeForComparison(event.StartAt))
			if eventIsCrossDay || eventStartMin >= 20*60 {
				if TimesOverlapCrossDay(startTime, endTime, true, formatTimeForComparison(event.StartAt), formatTimeForComparison(event.EndAt), eventIsCrossDay) {
					conflicts = append(conflicts, event)
				}
			}
		} else if isNextDay {
			// 隔天：檢查凌晨時段（08:00 以前）
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

// formatTimeForComparison formats a time.Time to "HH:MM" string for comparison
func formatTimeForComparison(t time.Time) string {
	return t.Format("15:04")
}

// ParseTimeToMinutes 將 HH:MM 格式轉換為當天分鐘數
func ParseTimeToMinutes(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return 0
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	return hour*60 + minute
}

// TimesOverlapCrossDay 處理跨日時間重疊檢測
// isCrossDay2 表示第二個課程的時間是否在隔天
// - 如果課程本身跨日（如 23:00-02:00），只有結束時間在隔天，加 end2Min
// - 如果課程在隔天（如 03:00-06:00），則開始和結束時間都在隔天，加 both
func TimesOverlapCrossDay(start1, end1 string, isCrossDay1 bool, start2, end2 string, isCrossDay2 bool) bool {
	// 將時間轉換為分鐘數
	start1Min := ParseTimeToMinutes(start1)
	end1Min := ParseTimeToMinutes(end1)
	start2Min := ParseTimeToMinutes(start2)
	end2Min := ParseTimeToMinutes(end2)

	// 對於跨日課程，將時間轉換到 48 小時制進行比較
	// 如果結束時間早於開始時間（表示跨日），則結束時間加 24*60
	if isCrossDay1 && end1Min < start1Min {
		end1Min += 24 * 60
	}
	// 如果第二個課程也是跨日課程
	if isCrossDay2 {
		// 如果結束時間早於開始時間（課程本身跨日，如 23:00-02:00），只加 end2Min
		if end2Min < start2Min {
			end2Min += 24 * 60
		} else {
			// 否則（課程在隔天，如 03:00-06:00），加 both
			start2Min += 24 * 60
			end2Min += 24 * 60
		}
	}

	// 檢查重疊
	return start1Min < end2Min && end1Min > start2Min
}

// TimesOverlapCrossDayWithNextDay 專門用於比較跨日課程與隔天課程的衝突檢測
// 參數：
// - start1, end1, isCrossDay1: 跨日課程的時間和標記
// - start2, end2: 隔天課程的時間（發生在跨日課程的隔天部分）
func TimesOverlapCrossDayWithNextDay(start1, end1 string, isCrossDay1 bool, start2, end2 string) bool {
	// 將時間轉換為分鐘數
	start1Min := ParseTimeToMinutes(start1)
	end1Min := ParseTimeToMinutes(end1)
	start2Min := ParseTimeToMinutes(start2)
	end2Min := ParseTimeToMinutes(end2)

	// 跨日課程的結束時間加 24*60 分鐘
	if isCrossDay1 {
		end1Min += 24 * 60
	}

	// 隔天課程的時間往回調整 24 小時（因為我們要與跨日課程的隔天部分比較）
	// 例如：隔天 01:00-03:00 調整為 25:00-27:00 (1+24, 3+24)
	start2Min += 24 * 60
	end2Min += 24 * 60

	// 檢查重疊
	return start1Min < end2Min && end1Min > start2Min
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
