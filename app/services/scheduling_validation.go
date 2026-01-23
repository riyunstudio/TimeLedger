package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

type ScheduleValidationServiceImpl struct {
	BaseService
	app              *app.App
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	roomRepo         *repositories.RoomRepository
	courseRepo       *repositories.CourseRepository
}

func NewScheduleValidationService(app *app.App) ScheduleValidationService {
	return &ScheduleValidationServiceImpl{
		app:              app,
		scheduleRuleRepo: repositories.NewScheduleRuleRepository(app),
		roomRepo:         repositories.NewRoomRepository(app),
		courseRepo:       repositories.NewCourseRepository(app),
	}
}

func (s *ScheduleValidationServiceImpl) CheckOverlap(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, excludeRuleID *uint) (ValidationResult, error) {
	result := ValidationResult{Valid: true}

	weekday := startTime.Weekday()
	if weekday == 0 {
		weekday = 7
	}

	startTimeStr := startTime.Format("15:04:05")
	endTimeStr := endTime.Format("15:04:05")

	query := s.app.MySQL.RDB.WithContext(ctx).Model(&models.ScheduleRule{}).
		Where("center_id = ?", centerID).
		Where("weekday = ?", weekday).
		Where("start_time < ?", endTimeStr).
		Where("end_time > ?", startTimeStr).
		Where("JSON_EXTRACT(effective_range, '$.start_date') <= ?", startTime.Format("2006-01-02")).
		Where("JSON_EXTRACT(effective_range, '$.end_date') >= ?", startTime.Format("2006-01-02"))

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

func (s *ScheduleValidationServiceImpl) ValidateFull(ctx context.Context, centerID uint, teacherID *uint, roomID uint, courseID uint, startTime, endTime time.Time, excludeRuleID *uint, allowBufferOverride bool) (ValidationResult, error) {
	result := ValidationResult{Valid: true}

	overlapResult, err := s.CheckOverlap(ctx, centerID, teacherID, roomID, startTime, endTime, excludeRuleID)
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

	if teacherID != nil {
		prevEndTime, _ := s.getPreviousSessionEndTime(ctx, centerID, *teacherID, startTime)
		if !prevEndTime.IsZero() {
			teacherBufferResult, err := s.CheckTeacherBuffer(ctx, centerID, *teacherID, prevEndTime, startTime, courseID)
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

	prevRoomEndTime, _ := s.getPreviousSessionEndTimeByRoom(ctx, centerID, roomID, startTime)
	if !prevRoomEndTime.IsZero() {
		roomBufferResult, err := s.CheckRoomBuffer(ctx, centerID, roomID, prevRoomEndTime, startTime, courseID)
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

	var rule models.ScheduleRule
	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("teacher_id = ?", teacherID).
		Where("weekday = ?", weekday).
		Where("end_time <= ?", startTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		return time.Time{}, nil
	}

	endTime, _ := time.Parse("2006-01-02 15:04:05", beforeTime.Format("2006-01-02")+" "+rule.EndTime)
	return endTime, nil
}

func (s *ScheduleValidationServiceImpl) getPreviousSessionEndTimeByRoom(ctx context.Context, centerID uint, roomID uint, beforeTime time.Time) (time.Time, error) {
	weekday := beforeTime.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	startTimeStr := beforeTime.Format("15:04:05")

	var rule models.ScheduleRule
	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("room_id = ?", roomID).
		Where("weekday = ?", weekday).
		Where("end_time <= ?", startTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		return time.Time{}, nil
	}

	endTime, _ := time.Parse("2006-01-02 15:04:05", beforeTime.Format("2006-01-02")+" "+rule.EndTime)
	return endTime, nil
}
