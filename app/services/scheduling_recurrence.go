package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"

	"gorm.io/gorm"
)

type ScheduleRecurrenceServiceImpl struct {
	BaseService
	app           *app.App
	ruleRepo      *repositories.ScheduleRuleRepository
	exceptionRepo *repositories.ScheduleExceptionRepository
	expansionSvc  ScheduleExpansionService
	auditLogRepo  *repositories.AuditLogRepository
	offeringRepo  *repositories.OfferingRepository
	cacheSvc      *CacheService
}

func NewScheduleRecurrenceService(app *app.App) ScheduleRecurrenceService {
	svc := &ScheduleRecurrenceServiceImpl{
		BaseService: *NewBaseService(app, "ScheduleRecurrenceService"),
		app:         app,
	}

	if app.MySQL != nil {
		svc.ruleRepo = repositories.NewScheduleRuleRepository(app)
		svc.exceptionRepo = repositories.NewScheduleExceptionRepository(app)
		svc.expansionSvc = NewScheduleExpansionService(app)
		svc.auditLogRepo = repositories.NewAuditLogRepository(app)
		svc.offeringRepo = repositories.NewOfferingRepository(app)
	}

	if app.Redis != nil {
		svc.cacheSvc = NewCacheService(app)
	}

	return svc
}

func (s *ScheduleRecurrenceServiceImpl) PreviewAffectedSessions(ctx context.Context, ruleID uint, editDate time.Time, mode RecurrenceEditMode) (RecurrenceEditPreview, error) {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return RecurrenceEditPreview{}, err
	}

	preview := RecurrenceEditPreview{
		Mode: mode,
	}

	switch mode {
	case RecurrenceEditSingle:
		preview.AffectedCount = 1
		preview.AffectedDates = []time.Time{editDate}
	case RecurrenceEditFuture:
		futureDates, err := s.getFutureOccurrences(rule, editDate)
		if err != nil {
			return RecurrenceEditPreview{}, err
		}
		preview.AffectedCount = len(futureDates)
		preview.AffectedDates = futureDates
		preview.WillCreateRule = true
	case RecurrenceEditAll:
		allDates, err := s.getAllOccurrences(rule)
		if err != nil {
			return RecurrenceEditPreview{}, err
		}
		preview.AffectedCount = len(allDates)
		preview.AffectedDates = allDates
	}

	return preview, nil
}

func (s *ScheduleRecurrenceServiceImpl) getFutureOccurrences(rule models.ScheduleRule, fromDate time.Time) ([]time.Time, error) {
	var dates []time.Time
	current := fromDate

	effectiveEnd := rule.EffectiveRange.EndDate
	if effectiveEnd.IsZero() {
		return dates, nil
	}

	for !current.After(effectiveEnd) {
		currentWeekday := int(current.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}
		if currentWeekday == rule.Weekday {
			dates = append(dates, current)
		}
		current = current.AddDate(0, 0, 1)
	}

	return dates, nil
}

func (s *ScheduleRecurrenceServiceImpl) getAllOccurrences(rule models.ScheduleRule) ([]time.Time, error) {
	var dates []time.Time
	current := rule.EffectiveRange.StartDate

	effectiveEnd := rule.EffectiveRange.EndDate
	if effectiveEnd.IsZero() {
		return dates, nil
	}

	for !current.After(effectiveEnd) {
		currentWeekday := int(current.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}
		if currentWeekday == rule.Weekday {
			dates = append(dates, current)
		}
		current = current.AddDate(0, 0, 1)
	}

	return dates, nil
}

func (s *ScheduleRecurrenceServiceImpl) EditRecurringSchedule(ctx context.Context, centerID uint, teacherID uint, req RecurrenceEditRequest) (RecurrenceEditResult, error) {
	rule, err := s.ruleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		return RecurrenceEditResult{}, err
	}

	result := RecurrenceEditResult{
		Mode:             req.Mode,
		CancelExceptions: []models.ScheduleException{},
		AddExceptions:    []models.ScheduleException{},
		AffectedCount:    0,
	}

	switch req.Mode {
	case RecurrenceEditSingle:
		cancelExc, addExc, err := s.editSingle(ctx, centerID, teacherID, rule, req)
		if err != nil {
			return RecurrenceEditResult{}, err
		}
		result.CancelExceptions = append(result.CancelExceptions, cancelExc)
		result.AddExceptions = append(result.AddExceptions, addExc)
		result.AffectedCount = 1

	case RecurrenceEditFuture:
		newRule, cancelExcs, addExcs, err := s.editFuture(ctx, centerID, teacherID, rule, req)
		if err != nil {
			return RecurrenceEditResult{}, err
		}
		result.NewRule = newRule
		result.CancelExceptions = cancelExcs
		result.AddExceptions = addExcs
		preview, _ := s.PreviewAffectedSessions(ctx, req.RuleID, req.EditDate, RecurrenceEditFuture)
		result.AffectedCount = preview.AffectedCount

	case RecurrenceEditAll:
		updatedRule, err := s.editAll(ctx, centerID, teacherID, rule, req)
		if err != nil {
			return RecurrenceEditResult{}, err
		}
		result.UpdatedRule = updatedRule
		preview, _ := s.PreviewAffectedSessions(ctx, req.RuleID, req.EditDate, RecurrenceEditAll)
		result.AffectedCount = preview.AffectedCount
	}

	s.invalidateRelatedCaches(ctx, centerID, rule)

	return result, nil
}

func (s *ScheduleRecurrenceServiceImpl) editSingle(ctx context.Context, centerID uint, teacherID uint, rule models.ScheduleRule, req RecurrenceEditRequest) (models.ScheduleException, models.ScheduleException, error) {
	cancelExc := models.ScheduleException{
		CenterID:      centerID,
		RuleID:        req.RuleID,
		OriginalDate:  req.EditDate,
		ExceptionType: "CANCEL",
		Status:        "PENDING",
		Reason:        req.Reason,
	}

	newStartAt := time.Date(req.EditDate.Year(), req.EditDate.Month(), req.EditDate.Day(), 0, 0, 0, 0, time.UTC)
	newEndAt := time.Date(req.EditDate.Year(), req.EditDate.Month(), req.EditDate.Day(), 0, 0, 0, 0, time.UTC)

	if req.NewStartTime != "" {
		parsedStart, _ := time.Parse("15:04:05", req.NewStartTime)
		newStartAt = time.Date(req.EditDate.Year(), req.EditDate.Month(), req.EditDate.Day(),
			parsedStart.Hour(), parsedStart.Minute(), 0, 0, time.UTC)
	}
	if req.NewEndTime != "" {
		parsedEnd, _ := time.Parse("15:04:05", req.NewEndTime)
		newEndAt = time.Date(req.EditDate.Year(), req.EditDate.Month(), req.EditDate.Day(),
			parsedEnd.Hour(), parsedEnd.Minute(), 0, 0, time.UTC)
	}

	addExc := models.ScheduleException{
		CenterID:      centerID,
		RuleID:        req.RuleID,
		OriginalDate:  req.EditDate,
		ExceptionType: "ADD",
		Status:        "PENDING",
		NewStartAt:    &newStartAt,
		NewEndAt:      &newEndAt,
		NewTeacherID:  req.NewTeacherID,
		NewRoomID:     req.NewRoomID,
		Reason:        req.Reason,
	}

	var createdCancel, createdAdd models.ScheduleException

	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var createErr error
		createdCancel, createErr = s.exceptionRepo.CreateWithDB(ctx, tx, cancelExc)
		if createErr != nil {
			return fmt.Errorf("failed to create cancel exception: %w", createErr)
		}

		createdAdd, createErr = s.exceptionRepo.CreateWithDB(ctx, tx, addExc)
		if createErr != nil {
			return fmt.Errorf("failed to create add exception: %w", createErr)
		}

		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "EDIT_SINGLE_OCCURRENCE",
			TargetType: "ScheduleException",
			TargetID:   createdCancel.ID,
			Payload:    models.AuditPayload{After: fmt.Sprintf("Edit single occurrence for rule %d on %s", req.RuleID, req.EditDate.Format("2006-01-02"))},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return models.ScheduleException{}, models.ScheduleException{}, txErr
	}

	return createdCancel, createdAdd, nil
}

func (s *ScheduleRecurrenceServiceImpl) editFuture(ctx context.Context, centerID uint, teacherID uint, rule models.ScheduleRule, req RecurrenceEditRequest) (*models.ScheduleRule, []models.ScheduleException, []models.ScheduleException, error) {
	preview, err := s.PreviewAffectedSessions(ctx, req.RuleID, req.EditDate, RecurrenceEditFuture)
	if err != nil {
		return nil, nil, nil, err
	}

	var cancelExcs []models.ScheduleException
	var addExcs []models.ScheduleException

	var createdRule *models.ScheduleRule

	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, date := range preview.AffectedDates {
			cancelExc := models.ScheduleException{
				CenterID:      centerID,
				RuleID:        req.RuleID,
				OriginalDate:  date,
				ExceptionType: "CANCEL",
				Status:        "PENDING",
				Reason:        req.Reason,
			}
			createdCancel, createErr := s.exceptionRepo.CreateWithDB(ctx, tx, cancelExc)
			if createErr != nil {
				return fmt.Errorf("failed to create cancel exception: %w", createErr)
			}
			cancelExcs = append(cancelExcs, createdCancel)

			newStartAt := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
			newEndAt := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

			startTime := req.NewStartTime
			if startTime == "" {
				startTime = rule.StartTime
			}
			endTime := req.NewEndTime
			if endTime == "" {
				endTime = rule.EndTime
			}

			parsedStart, _ := time.Parse("15:04:05", startTime)
			newStartAt = time.Date(date.Year(), date.Month(), date.Day(),
				parsedStart.Hour(), parsedStart.Minute(), 0, 0, time.UTC)

			parsedEnd, _ := time.Parse("15:04:05", endTime)
			newEndAt = time.Date(date.Year(), date.Month(), date.Day(),
				parsedEnd.Hour(), parsedEnd.Minute(), 0, 0, time.UTC)

			addExc := models.ScheduleException{
				CenterID:      centerID,
				RuleID:        req.RuleID,
				OriginalDate:  date,
				ExceptionType: "ADD",
				Status:        "PENDING",
				NewStartAt:    &newStartAt,
				NewEndAt:      &newEndAt,
				NewTeacherID:  req.NewTeacherID,
				NewRoomID:     req.NewRoomID,
				Reason:        req.Reason,
			}
			createdAdd, createErr := s.exceptionRepo.CreateWithDB(ctx, tx, addExc)
			if createErr != nil {
				return fmt.Errorf("failed to create add exception: %w", createErr)
			}
			addExcs = append(addExcs, createdAdd)
		}

		cutoffDate := req.EditDate.AddDate(0, 0, -1)
		rule.EffectiveRange.EndDate = cutoffDate
		if err := tx.Omit("created_at").Save(&rule).Error; err != nil {
			return fmt.Errorf("failed to update current rule: %w", err)
		}

		newRoomID := rule.RoomID
		if req.NewRoomID != nil {
			newRoomID = *req.NewRoomID
		}

		newRule := models.ScheduleRule{
			CenterID:   centerID,
			OfferingID: rule.OfferingID,
			TeacherID:  req.NewTeacherID,
			RoomID:     newRoomID,
			Name:       rule.Name,
			Weekday:    rule.Weekday,
			StartTime:  req.NewStartTime,
			EndTime:    req.NewEndTime,
			Duration:   rule.Duration,
			EffectiveRange: models.DateRange{
				StartDate: req.EditDate,
				EndDate:   time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		if err := tx.Create(&newRule).Error; err != nil {
			return fmt.Errorf("failed to create new rule: %w", err)
		}
		createdRule = &newRule

		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "EDIT_FUTURE_OCCURRENCES",
			TargetType: "ScheduleRule",
			TargetID:   req.RuleID,
			Payload:    models.AuditPayload{After: fmt.Sprintf("Edit future occurrences for rule %d starting from %s", req.RuleID, req.EditDate.Format("2006-01-02"))},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, nil, nil, txErr
	}

	return createdRule, cancelExcs, addExcs, nil
}

func (s *ScheduleRecurrenceServiceImpl) editAll(ctx context.Context, centerID uint, teacherID uint, rule models.ScheduleRule, req RecurrenceEditRequest) (*models.ScheduleRule, error) {
	if req.NewTeacherID != nil {
		rule.TeacherID = req.NewTeacherID
	}
	if req.NewRoomID != nil {
		rule.RoomID = *req.NewRoomID
	}
	if req.NewStartTime != "" {
		rule.StartTime = req.NewStartTime
	}
	if req.NewEndTime != "" {
		rule.EndTime = req.NewEndTime
	}

	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return nil, err
	}

	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "EDIT_ALL_OCCURRENCES",
		TargetType: "ScheduleRule",
		TargetID:   req.RuleID,
		Payload:    models.AuditPayload{After: fmt.Sprintf("Edit all occurrences for rule %d", req.RuleID)},
	})

	return &rule, nil
}

func (s *ScheduleRecurrenceServiceImpl) DeleteRecurringSchedule(ctx context.Context, centerID uint, teacherID uint, ruleID uint, editDate time.Time, mode RecurrenceEditMode, reason string) (RecurrenceEditResult, error) {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return RecurrenceEditResult{}, err
	}

	result := RecurrenceEditResult{
		Mode:             mode,
		CancelExceptions: []models.ScheduleException{},
		AffectedCount:    0,
	}

	switch mode {
	case RecurrenceEditSingle:
		cancelExc := models.ScheduleException{
			CenterID:      centerID,
			RuleID:        ruleID,
			OriginalDate:  editDate,
			ExceptionType: "CANCEL",
			Status:        "PENDING",
			Reason:        reason,
		}
		created, err := s.exceptionRepo.Create(ctx, cancelExc)
		if err != nil {
			return RecurrenceEditResult{}, err
		}
		result.CancelExceptions = append(result.CancelExceptions, created)
		result.AffectedCount = 1

	case RecurrenceEditFuture:
		preview, _ := s.PreviewAffectedSessions(ctx, ruleID, editDate, RecurrenceEditFuture)
		for _, date := range preview.AffectedDates {
			cancelExc := models.ScheduleException{
				CenterID:      centerID,
				RuleID:        ruleID,
				OriginalDate:  date,
				ExceptionType: "CANCEL",
				Status:        "PENDING",
				Reason:        reason,
			}
			created, _ := s.exceptionRepo.Create(ctx, cancelExc)
			result.CancelExceptions = append(result.CancelExceptions, created)
		}
		result.AffectedCount = preview.AffectedCount

		cutoffDate := editDate.AddDate(0, 0, -1)
		rule.EffectiveRange.EndDate = cutoffDate
		_ = s.ruleRepo.Update(ctx, rule)

	case RecurrenceEditAll:
		preview, _ := s.PreviewAffectedSessions(ctx, ruleID, rule.EffectiveRange.StartDate, RecurrenceEditAll)
		_ = s.ruleRepo.DeleteByID(ctx, ruleID)
		result.AffectedCount = preview.AffectedCount
	}

	s.invalidateRelatedCaches(ctx, centerID, rule)

	return result, nil
}

func (s *ScheduleRecurrenceServiceImpl) invalidateRelatedCaches(ctx context.Context, centerID uint, rule models.ScheduleRule) {
	pattern := fmt.Sprintf("schedule:expand:center:%d:*", centerID)
	_ = s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, pattern)

	if rule.ID != 0 && rule.TeacherID != nil {
		teacherPattern := fmt.Sprintf("schedule:expand:teacher:%d:center:%d:*", *rule.TeacherID, centerID)
		_ = s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, teacherPattern)
	}
}
