package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

type ScheduleValidationServiceImpl struct {
	BaseService
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	roomRepo         *repositories.RoomRepository
	courseRepo       *repositories.CourseRepository
}

func NewScheduleValidationService(app *app.App) ScheduleValidationService {
	baseSvc := NewBaseService(app, "ScheduleValidationService")
	svc := &ScheduleValidationServiceImpl{
		BaseService: *baseSvc,
	}

	if app.MySQL != nil {
		svc.scheduleRuleRepo = repositories.NewScheduleRuleRepository(app)
		svc.roomRepo = repositories.NewRoomRepository(app)
		svc.courseRepo = repositories.NewCourseRepository(app)
	}

	return svc
}

func (s *ScheduleValidationServiceImpl) CheckOverlap(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, weekday int, excludeRuleID *uint) (ValidationResult, error) {
	result := ValidationResult{Valid: true}

	// 確保 weekday 正確轉換（Go 的 Weekday: Sunday=0, Monday=1, ..., Saturday=6）
	// 資料庫使用: Monday=1, Tuesday=2, ..., Sunday=7
	if weekday == 0 {
		weekday = 7
	}

	startTimeStr := startTime.Format("15:04:05")
	endTimeStr := endTime.Format("15:04:05")

	// 處理零值時間，避免 MySQL 8.0 報錯 (如 0000-01-01)
	startDateStr := "0001-01-01"
	if !startTime.IsZero() {
		startDateStr = startTime.Format("2006-01-02")
	}

	query := s.App.MySQL.RDB.WithContext(ctx).Model(&models.ScheduleRule{}).
		Where("center_id = ?", centerID).
		Where("weekday = ?", weekday).
		Where("start_time < ?", endTimeStr).
		Where("end_time > ?", startTimeStr).
		Where("COALESCE(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')), ''), 'null'), '0001-01-01') <= ?", startDateStr).
		Where("COALESCE(NULLIF(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')), ''), 'null'), '0001-01-01 00:00:00'), '9999-12-31') >= ?", startDateStr)

	if teacherID != nil {
		query = query.Where("teacher_id = ?", *teacherID)
	}

	if excludeRuleID != nil {
		query = query.Where("id != ?", *excludeRuleID)
	}

	var rules []models.ScheduleRule
	if err := query.Find(&rules).Error; err != nil {
		return ValidationResult{}, err
	}

	for _, rule := range rules {
		var conflictTypes []string
		var messages []string

		if teacherID != nil && rule.TeacherID != nil && *rule.TeacherID == *teacherID {
			conflictTypes = append(conflictTypes, "TEACHER_OVERLAP")
			messages = append(messages, "老師在該時段已有課程安排")
		}

		if rule.RoomID == roomID {
			conflictTypes = append(conflictTypes, "ROOM_OVERLAP")
			messages = append(messages, "教室在該時段已被佔用")
		}

		if len(conflictTypes) > 0 {
			result.Valid = false
			for i, ct := range conflictTypes {
				result.Conflicts = append(result.Conflicts, ValidationConflict{
					Type:    ct,
					Message: messages[i],
					Details: fmt.Sprintf("rule_id:%d, offering_id:%d", rule.ID, rule.OfferingID),
				})
			}
		}
	}

	return result, nil
}

func (s *ScheduleValidationServiceImpl) CheckTeacherBuffer(ctx context.Context, centerID uint, teacherID uint, prevEndTime, nextStartTime time.Time, courseID uint) (ValidationResult, error) {
	result := ValidationResult{Valid: true}

	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return ValidationResult{}, err
	}

	bufferMin := course.TeacherBufferMin
	if bufferMin == 0 {
		return result, nil
	}

	bufferDuration := time.Duration(bufferMin) * time.Minute
	minGap := prevEndTime.Add(bufferDuration)

	if nextStartTime.Before(minGap) {
		diffMinutes := int(nextStartTime.Sub(prevEndTime).Minutes())
		result.Valid = false
		result.Conflicts = append(result.Conflicts, ValidationConflict{
			Type:             "TEACHER_BUFFER",
			Message:          fmt.Sprintf("老師轉場時間不足，需間隔 %d 分鐘，實際間隔 %d 分鐘", bufferMin, diffMinutes),
			CanOverride:      true,
			RequiredMinutes:  bufferMin,
			DiffMinutes:      bufferMin - diffMinutes,
			ConflictSource:   "PREV_SESSION",
			ConflictSourceID: 0,
		})
	}

	return result, nil
}

func (s *ScheduleValidationServiceImpl) CheckRoomBuffer(ctx context.Context, centerID uint, roomID uint, prevEndTime, nextStartTime time.Time, courseID uint) (ValidationResult, error) {
	result := ValidationResult{Valid: true}

	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return ValidationResult{}, err
	}

	bufferMin := course.RoomBufferMin
	if bufferMin == 0 {
		return result, nil
	}

	bufferDuration := time.Duration(bufferMin) * time.Minute
	minGap := prevEndTime.Add(bufferDuration)

	if nextStartTime.Before(minGap) {
		diffMinutes := int(nextStartTime.Sub(prevEndTime).Minutes())
		result.Valid = false
		result.Conflicts = append(result.Conflicts, ValidationConflict{
			Type:             "ROOM_BUFFER",
			Message:          fmt.Sprintf("教室清潔時間不足，需間隔 %d 分鐘，實際間隔 %d 分鐘", bufferMin, diffMinutes),
			CanOverride:      true,
			RequiredMinutes:  bufferMin,
			DiffMinutes:      bufferMin - diffMinutes,
			ConflictSource:   "PREV_SESSION",
			ConflictSourceID: 0,
		})
	}

	return result, nil
}

func (s *ScheduleValidationServiceImpl) ValidateFull(ctx context.Context, centerID uint, teacherID *uint, roomID uint, courseID uint, startTime, endTime time.Time, excludeRuleID *uint, allowBufferOverride bool, prevEndTime, nextStartTime *time.Time) (ValidationResult, error) {
	result := ValidationResult{Valid: true}

	// 計算 weekday
	weekday := int(startTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	overlapResult, err := s.CheckOverlap(ctx, centerID, teacherID, roomID, startTime, endTime, weekday, excludeRuleID)
	if err != nil {
		return ValidationResult{}, err
	}

	if !overlapResult.Valid {
		result.Valid = false
		for _, c := range overlapResult.Conflicts {
			c.CanOverride = false
			result.Conflicts = append(result.Conflicts, c)
		}
	}

	// 處理教師緩衝時間
	// 如果提供了 prevEndTime 和 nextStartTime，直接使用
	// 否則自動計算上一堂課的結束時間
	if teacherID != nil {
		var teacherPrevEndTime time.Time
		if prevEndTime != nil && !prevEndTime.IsZero() {
			teacherPrevEndTime = *prevEndTime
		} else {
			// 自動計算：查詢該教師在同一天 startTime 之前的上一堂課
			calculatedPrevEndTime, _ := s.getPreviousSessionEndTime(ctx, centerID, *teacherID, startTime)
			teacherPrevEndTime = calculatedPrevEndTime
		}

		// 使用 nextStartTime 或 startTime 作為檢查點
		checkTime := startTime
		if nextStartTime != nil && !nextStartTime.IsZero() {
			checkTime = *nextStartTime
		}

		if !teacherPrevEndTime.IsZero() {
			teacherBufferResult, err := s.CheckTeacherBuffer(ctx, centerID, *teacherID, teacherPrevEndTime, checkTime, courseID)
			if err != nil {
				return ValidationResult{}, err
			}

			if !teacherBufferResult.Valid {
				if !allowBufferOverride {
					result.Valid = false
					result.Conflicts = append(result.Conflicts, teacherBufferResult.Conflicts...)
				} else {
					for i := range result.Conflicts {
						result.Conflicts[i].CanOverride = true
					}
				}
			}
		}
	}

	// 處理教室緩衝時間
	// 如果提供了 prevEndTime 和 nextStartTime，直接使用
	// 否則自動計算上一堂課的結束時間
	var roomPrevEndTime time.Time
	if prevEndTime != nil && !prevEndTime.IsZero() {
		roomPrevEndTime = *prevEndTime
	} else {
		// 自動計算：查詢該教室在同一天 startTime 之前的上一堂課
		calculatedPrevEndTime, _ := s.getPreviousSessionEndTimeByRoom(ctx, centerID, roomID, startTime)
		roomPrevEndTime = calculatedPrevEndTime
	}

	// 使用 nextStartTime 或 startTime 作為檢查點
	checkTime := startTime
	if nextStartTime != nil && !nextStartTime.IsZero() {
		checkTime = *nextStartTime
	}

	if !roomPrevEndTime.IsZero() {
		roomBufferResult, err := s.CheckRoomBuffer(ctx, centerID, roomID, roomPrevEndTime, checkTime, courseID)
		if err != nil {
			return ValidationResult{}, err
		}

		if !roomBufferResult.Valid {
			if !allowBufferOverride {
				result.Valid = false
				result.Conflicts = append(result.Conflicts, roomBufferResult.Conflicts...)
			} else {
				for i := range result.Conflicts {
					result.Conflicts[i].CanOverride = true
				}
			}
		}
	}

	return result, nil
}

func (s *ScheduleValidationServiceImpl) getPreviousSessionEndTime(ctx context.Context, centerID uint, teacherID uint, beforeTime time.Time) (time.Time, error) {
	weekday := beforeTime.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	startTimeStr := beforeTime.Format("15:04:05")

	// 處理零值時間
	startDateStr := "0001-01-01"
	if !beforeTime.IsZero() {
		startDateStr = beforeTime.Format("2006-01-02")
	}

	var rule models.ScheduleRule
	err := s.App.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("teacher_id = ?", teacherID).
		Where("weekday = ?", weekday).
		Where("COALESCE(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')), ''), 'null'), '0001-01-01') <= ?", startDateStr).
		Where("COALESCE(NULLIF(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')), ''), 'null'), '0001-01-01 00:00:00'), '9999-12-31') >= ?", startDateStr).
		Where("end_time <= ?", startTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		return time.Time{}, nil
	}

	// 使用中央時區解析時間
	// 注意：EndTime 格式是 "HH:mm"（沒有秒）
	loc := app.GetTaiwanLocation()

	// 解析 EndTime 並使用與 beforeTime 相同的日期
	// 這樣才能正確計算緩衝時間
	endTimeParts := strings.Split(rule.EndTime, ":")
	if len(endTimeParts) >= 2 {
		hour, _ := strconv.Atoi(endTimeParts[0])
		minute, _ := strconv.Atoi(endTimeParts[1])

		// 使用 beforeTime 的日期來構造 endTime，這樣才能正確比較
		endTime := time.Date(
			beforeTime.Year(), beforeTime.Month(), beforeTime.Day(),
			hour, minute, 0, 0, loc,
		)
		return endTime, nil
	}

	// Fallback: 使用舊的解析方式
	timeStr := "2000-01-01" + " " + rule.EndTime
	endTime, _ := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
	return endTime, nil
}

func (s *ScheduleValidationServiceImpl) getPreviousSessionEndTimeByRoom(ctx context.Context, centerID uint, roomID uint, beforeTime time.Time) (time.Time, error) {
	weekday := beforeTime.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	startTimeStr := beforeTime.Format("15:04:05")

	// 處理零值時間
	startDateStr := "0001-01-01"
	if !beforeTime.IsZero() {
		startDateStr = beforeTime.Format("2006-01-02")
	}

	var rule models.ScheduleRule
	err := s.App.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("room_id = ?", roomID).
		Where("weekday = ?", weekday).
		Where("COALESCE(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')), ''), 'null'), '0001-01-01') <= ?", startDateStr).
		Where("COALESCE(NULLIF(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')), ''), 'null'), '0001-01-01 00:00:00'), '9999-12-31') >= ?", startDateStr).
		Where("end_time <= ?", startTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		return time.Time{}, nil
	}

	// 使用中央時區解析時間
	// 注意：EndTime 格式是 "HH:mm"（沒有秒）
	loc := app.GetTaiwanLocation()

	// 解析 EndTime 並使用與 beforeTime 相同的日期
	// 這樣才能正確計算緩衝時間
	endTimeParts := strings.Split(rule.EndTime, ":")
	if len(endTimeParts) >= 2 {
		hour, _ := strconv.Atoi(endTimeParts[0])
		minute, _ := strconv.Atoi(endTimeParts[1])

		// 使用 beforeTime 的日期來構造 endTime，這樣才能正確比較
		endTime := time.Date(
			beforeTime.Year(), beforeTime.Month(), beforeTime.Day(),
			hour, minute, 0, 0, loc,
		)
		return endTime, nil
	}

	// Fallback: 使用舊的解析方式
	timeStr := "2000-01-01" + " " + rule.EndTime
	endTime, _ := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
	return endTime, nil
}
