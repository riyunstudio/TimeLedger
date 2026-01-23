package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

type ScheduleExpansionServiceImpl struct {
	BaseService
	app              *app.App
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	exceptionRepo    *repositories.ScheduleExceptionRepository
	auditLogRepo     *repositories.AuditLogRepository
	centerRepo       *repositories.CenterRepository
	holidayRepo      *repositories.CenterHolidayRepository
}

func NewScheduleExpansionService(app *app.App) ScheduleExpansionService {
	return &ScheduleExpansionServiceImpl{
		app:              app,
		scheduleRuleRepo: repositories.NewScheduleRuleRepository(app),
		exceptionRepo:    repositories.NewScheduleExceptionRepository(app),
		auditLogRepo:     repositories.NewAuditLogRepository(app),
		centerRepo:       repositories.NewCenterRepository(app),
		holidayRepo:      repositories.NewCenterHolidayRepository(app),
	}
}

func (s *ScheduleExpansionServiceImpl) ExpandRules(ctx context.Context, rules []models.ScheduleRule, startDate, endDate time.Time, centerID uint) []ExpandedSchedule {
	var schedules []ExpandedSchedule

	holidays, _ := s.holidayRepo.ListByDateRange(ctx, centerID, startDate, endDate)
	holidaySet := make(map[string]bool)
	for _, h := range holidays {
		holidaySet[h.Date.Format("2006-01-02")] = true
	}

	for _, rule := range rules {
		if rule.Weekday == 0 {
			continue
		}

		ruleStartDate := rule.EffectiveRange.StartDate
		ruleEndDate := rule.EffectiveRange.EndDate

		date := startDate
		for date.Before(endDate) || date.Equal(endDate) {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			if weekday == int(rule.Weekday) {
				isWithinEffectiveRange := true
				if !ruleStartDate.IsZero() && date.Before(ruleStartDate) {
					isWithinEffectiveRange = false
				}
				if !ruleEndDate.IsZero() && date.After(ruleEndDate) {
					isWithinEffectiveRange = false
				}

				if isWithinEffectiveRange {
					dateStr := date.Format("2006-01-02")
					isHoliday := holidaySet[dateStr]

					exceptions, _ := s.exceptionRepo.GetByRuleAndDate(ctx, rule.ID, date)
					hasException := len(exceptions) > 0

					schedule := ExpandedSchedule{
						RuleID:       rule.ID,
						Date:         date,
						StartTime:    rule.StartTime,
						EndTime:      rule.EndTime,
						RoomID:       rule.RoomID,
						TeacherID:    rule.TeacherID,
						IsHoliday:    isHoliday,
						HasException: hasException,
					}

					schedules = append(schedules, schedule)
				}
			}

			date = date.AddDate(0, 0, 1)
		}
	}

	return schedules
}

func (s *ScheduleExpansionServiceImpl) GetEffectiveRuleForDate(ctx context.Context, offeringID uint, date time.Time) (*models.ScheduleRule, error) {
	var rules []models.ScheduleRule
	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("offering_id = ?", offeringID).
		Where("weekday = ?", func() int {
			weekday := int(date.Weekday())
			if weekday == 0 {
				return 7
			}
			return weekday
		}()).
		Where("JSON_EXTRACT(effective_range, '$.start_date') <= ?", date.Format("2006-01-02")).
		Where(func() string {
			return "JSON_EXTRACT(effective_range, '$.end_date') >= ? OR JSON_EXTRACT(effective_range, '$.end_date') = '\"0001-01-01\"' OR JSON_EXTRACT(effective_range, '$.end_date') IS NULL"
		}()).
		Find(&rules).Error

	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return nil, nil
	}

	return &rules[0], nil
}

func (s *ScheduleExpansionServiceImpl) DetectPhaseTransitions(ctx context.Context, centerID uint, offeringID uint, startDate, endDate time.Time) ([]PhaseTransition, error) {
	rules, err := s.scheduleRuleRepo.ListByOfferingID(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	if len(rules) <= 1 {
		return []PhaseTransition{}, nil
	}

	var transitions []PhaseTransition
	date := startDate
	prevRule := (*models.ScheduleRule)(nil)

	for date.Before(endDate) || date.Equal(endDate) {
		currentRule, _ := s.GetEffectiveRuleForDate(ctx, offeringID, date)

		if prevRule != nil && currentRule != nil {
			prevRuleID := prevRule.ID
			currRuleID := currentRule.ID

			if prevRuleID != currRuleID ||
				(prevRule.RoomID != currentRule.RoomID) ||
				!ptrEqual(prevRule.TeacherID, currentRule.TeacherID) ||
				(prevRule.StartTime != currentRule.StartTime) ||
				(prevRule.EndTime != currentRule.EndTime) {

				transition := PhaseTransition{
					Date:          date,
					PrevRuleID:    &prevRuleID,
					PrevRoomID:    &prevRule.RoomID,
					PrevTeacherID: prevRule.TeacherID,
					PrevStartTime: prevRule.StartTime,
					PrevEndTime:   prevRule.EndTime,
					NextRuleID:    &currRuleID,
					NextRoomID:    &currentRule.RoomID,
					NextTeacherID: currentRule.TeacherID,
					NextStartTime: currentRule.StartTime,
					NextEndTime:   currentRule.EndTime,
					HasGap:        false,
				}
				transitions = append(transitions, transition)
			}
		} else if prevRule != nil && currentRule == nil {
			prevRuleID := prevRule.ID
			transition := PhaseTransition{
				Date:       date,
				PrevRuleID: &prevRuleID,
				PrevRoomID: &prevRule.RoomID,
				HasGap:     true,
			}
			transitions = append(transitions, transition)
		} else if prevRule == nil && currentRule != nil {
			currRuleID := currentRule.ID
			transition := PhaseTransition{
				Date:       date,
				NextRuleID: &currRuleID,
				NextRoomID: &currentRule.RoomID,
			}
			transitions = append(transitions, transition)
		}

		prevRule = currentRule
		date = date.AddDate(0, 0, 1)
	}

	return transitions, nil
}

func (s *ScheduleExpansionServiceImpl) GetRulesByEffectiveDateRange(ctx context.Context, centerID uint, offeringID uint, startDate, endDate time.Time) ([]models.ScheduleRule, error) {
	var rules []models.ScheduleRule
	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("offering_id = ?", offeringID).
		Where("JSON_EXTRACT(effective_range, '$.start_date') <= ?", endDate.Format("2006-01-02")).
		Where(func() string {
			return "(JSON_EXTRACT(effective_range, '$.end_date') >= ? OR JSON_EXTRACT(effective_range, '$.end_date') = '\"0001-01-01\"' OR JSON_EXTRACT(effective_range, '$.end_date') IS NULL)"
		}()).
		Order("JSON_EXTRACT(effective_range, '$.start_date') ASC").
		Find(&rules).Error

	return rules, err
}

func ptrEqual(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

type ScheduleExceptionServiceImpl struct {
	BaseService
	app               *app.App
	exceptionRepo     *repositories.ScheduleExceptionRepository
	ruleRepo          *repositories.ScheduleRuleRepository
	auditLogRepo      *repositories.AuditLogRepository
	centerRepo        *repositories.CenterRepository
	validationService ScheduleValidationService
}

func NewScheduleExceptionService(app *app.App) ScheduleExceptionService {
	svc := &ScheduleExceptionServiceImpl{
		app:               app,
		exceptionRepo:     repositories.NewScheduleExceptionRepository(app),
		ruleRepo:          repositories.NewScheduleRuleRepository(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
		centerRepo:        repositories.NewCenterRepository(app),
		validationService: NewScheduleValidationService(app),
	}
	return svc
}

func (s *ScheduleExceptionServiceImpl) CheckExceptionDeadline(ctx context.Context, centerID uint, ruleID uint, exceptionDate time.Time) (bool, string, error) {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return false, "", err
	}

	now := time.Now()

	if rule.LockAt != nil && now.After(*rule.LockAt) {
		return false, "已超過異動截止日", nil
	}

	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err != nil {
		return false, "", err
	}

	leadDays := center.Settings.ExceptionLeadDays
	if leadDays <= 0 {
		leadDays = 14
	}

	deadline := exceptionDate.AddDate(0, 0, -leadDays)
	if now.After(deadline) {
		return false, fmt.Sprintf("已超過異動截止日（需提前 %d 天申請）", leadDays), nil
	}

	return true, "", nil
}

func (s *ScheduleExceptionServiceImpl) CreateException(ctx context.Context, centerID uint, teacherID uint, ruleID uint, originalDate time.Time, exceptionType string, newStartAt, newEndAt *time.Time, newTeacherID *uint, reason string) (models.ScheduleException, error) {
	allowed, reasonStr, err := s.CheckExceptionDeadline(ctx, centerID, ruleID, originalDate)
	if err != nil {
		return models.ScheduleException{}, err
	}
	if !allowed {
		return models.ScheduleException{}, errors.New(reasonStr)
	}

	exception := models.ScheduleException{
		CenterID:     centerID,
		RuleID:       ruleID,
		OriginalDate: originalDate,
		Type:         exceptionType,
		Status:       "PENDING",
		NewStartAt:   newStartAt,
		NewEndAt:     newEndAt,
		NewTeacherID: newTeacherID,
		Reason:       reason,
	}

	createdException, err := s.exceptionRepo.Create(ctx, exception)
	if err != nil {
		return models.ScheduleException{}, err
	}

	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "CREATE_EXCEPTION",
		TargetType: "ScheduleException",
		TargetID:   createdException.ID,
		Payload:    models.AuditPayload{After: exception},
	})

	return createdException, nil
}

func (s *ScheduleExceptionServiceImpl) RevokeException(ctx context.Context, exceptionID uint, teacherID uint) error {
	exception, err := s.exceptionRepo.GetByID(ctx, exceptionID)
	if err != nil {
		return err
	}

	if exception.Status != "PENDING" {
		return errors.New("only pending exceptions can be revoked")
	}

	exception.Status = "REVOKED"
	if err := s.exceptionRepo.Update(ctx, exception); err != nil {
		return err
	}

	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   exception.CenterID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "REVOKE_EXCEPTION",
		TargetType: "ScheduleException",
		TargetID:   exceptionID,
		Payload:    models.AuditPayload{Before: "PENDING", After: "REVOKED"},
	})

	return nil
}

func (s *ScheduleExceptionServiceImpl) ReviewException(ctx context.Context, exceptionID uint, adminID uint, action string, overrideBuffer bool, reason string) error {
	exception, err := s.exceptionRepo.GetByID(ctx, exceptionID)
	if err != nil {
		return err
	}

	if exception.Status != "PENDING" {
		return errors.New("only pending exceptions can be reviewed")
	}

	oldStatus := exception.Status
	exception.Status = action
	exception.ReviewedBy = &adminID
	now := time.Now()
	exception.ReviewedAt = &now

	if action == "APPROVED" {
		rule, err := s.ruleRepo.GetByID(ctx, exception.RuleID)
		if err != nil {
			return err
		}

		var startAt, endAt time.Time
		if exception.Type == "RESCHEDULE" && exception.NewStartAt != nil {
			startAt = *exception.NewStartAt
			endAt = *exception.NewEndAt
		} else {
			startAt = exception.OriginalDate
			endAt = exception.OriginalDate
		}

		validateResult, err := s.validationService.ValidateFull(
			ctx,
			exception.CenterID,
			exception.NewTeacherID,
			rule.RoomID,
			rule.OfferingID,
			startAt,
			endAt,
			nil,
			overrideBuffer,
		)
		if err != nil {
			return err
		}

		if !validateResult.Valid {
			var hasHardOverlap bool
			var hasBufferConflict bool
			for _, c := range validateResult.Conflicts {
				if c.Type == "OVERLAP" || c.Type == "TEACHER_OVERLAP" || c.Type == "ROOM_OVERLAP" {
					hasHardOverlap = true
				}
				if c.Type == "TEACHER_BUFFER" || c.Type == "ROOM_BUFFER" {
					hasBufferConflict = true
				}
			}

			if hasHardOverlap {
				return errors.New("approval rejected: new time slot has hard overlap with existing schedule")
			}

			if hasBufferConflict && !overrideBuffer {
				return errors.New("approval rejected: new time slot has buffer conflict and override is not allowed")
			}
		}
	}

	if err := s.exceptionRepo.Update(ctx, exception); err != nil {
		return err
	}

	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   exception.CenterID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "REVIEW_EXCEPTION_" + action,
		TargetType: "ScheduleException",
		TargetID:   exceptionID,
		Payload:    models.AuditPayload{Before: oldStatus, After: action},
	})

	return nil
}

func (s *ScheduleExceptionServiceImpl) GetExceptionsByRule(ctx context.Context, ruleID uint) ([]models.ScheduleException, error) {
	exceptions, err := s.exceptionRepo.GetByRuleAndDate(ctx, ruleID, time.Time{})
	if err != nil {
		return nil, err
	}

	var filteredExceptions []models.ScheduleException
	for _, exc := range exceptions {
		filteredExceptions = append(filteredExceptions, exc)
	}

	return filteredExceptions, nil
}

func (s *ScheduleExceptionServiceImpl) GetExceptionsByDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]models.ScheduleException, error) {
	var exceptions []models.ScheduleException
	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("original_date >= ?", startDate).
		Where("original_date <= ?", endDate).
		Where("status IN ('PENDING', 'APPROVED', 'REJECTED')").
		Find(&exceptions).Error

	return exceptions, err
}

func (s *ScheduleExceptionServiceImpl) GetPendingExceptions(ctx context.Context, centerID uint) ([]models.ScheduleException, error) {
	var exceptions []models.ScheduleException
	err := s.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("status = ?", "PENDING").
		Order("created_at ASC").
		Find(&exceptions).Error

	return exceptions, err
}
