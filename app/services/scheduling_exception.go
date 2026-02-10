package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"

	"gorm.io/gorm"
)

type ScheduleExceptionServiceImpl struct {
	BaseService
	exceptionRepo     *repositories.ScheduleExceptionRepository
	ruleRepo          *repositories.ScheduleRuleRepository
	auditLogRepo      *repositories.AuditLogRepository
	centerRepo        *repositories.CenterRepository
	teacherRepo       *repositories.TeacherRepository
	validationService ScheduleValidationService
	notificationSvc   NotificationService
	notificationQueue NotificationQueueService
	cacheSvc          *CacheService
}

func NewScheduleExceptionService(app *app.App) ScheduleExceptionService {
	baseSvc := NewBaseService(app, "ScheduleExceptionService")
	svc := &ScheduleExceptionServiceImpl{
		BaseService: *baseSvc,
	}

	if app.MySQL != nil {
		svc.exceptionRepo = repositories.NewScheduleExceptionRepository(app)
		svc.ruleRepo = repositories.NewScheduleRuleRepository(app)
		svc.auditLogRepo = repositories.NewAuditLogRepository(app)
		svc.centerRepo = repositories.NewCenterRepository(app)
		svc.teacherRepo = repositories.NewTeacherRepository(app)
		svc.validationService = NewScheduleValidationService(app)
		svc.notificationSvc = NewNotificationService(app)
		svc.notificationQueue = NewNotificationQueueService(app)
		svc.cacheSvc = NewCacheService(app)
	}

	return svc
}

func (s *ScheduleExceptionServiceImpl) CheckExceptionDeadline(ctx context.Context, centerID uint, ruleID uint, exceptionDate time.Time) (bool, *errInfos.Res, error) {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		errInfo := s.App.Err.New(errInfos.SQL_ERROR)
		return false, errInfo, err
	}

	now := time.Now()

	if rule.LockAt != nil && now.After(*rule.LockAt) {
		errInfo := s.App.Err.New(errInfos.EXCEPTION_DEADLINE_EXCEEDED)
		return false, errInfo, nil
	}

	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err != nil {
		errInfo := s.App.Err.New(errInfos.SQL_ERROR)
		return false, errInfo, err
	}

	leadDays := center.Settings.ExceptionLeadDays

	// 0 = Unlimited（不限制截止日）
	if leadDays == 0 {
		return true, nil, nil
	}

	// 如果設定為負數，視為 14 天（向後相容）
	if leadDays < 0 {
		leadDays = 14
	}

	deadline := exceptionDate.AddDate(0, 0, -leadDays)
	if now.After(deadline) {
		errInfo := s.App.Err.New(errInfos.EXCEPTION_DEADLINE_EXCEEDED)
		return false, errInfo, nil
	}

	return true, nil, nil
}

func (s *ScheduleExceptionServiceImpl) CreateException(ctx context.Context, centerID uint, teacherID uint, ruleID uint, req *CreateExceptionRequest) (models.ScheduleException, *errInfos.Res, error) {
	allowed, errInfo, err := s.CheckExceptionDeadline(ctx, centerID, ruleID, req.OriginalDate)
	if err != nil {
		return models.ScheduleException{}, nil, err
	}
	if !allowed {
		return models.ScheduleException{}, errInfo, nil
	}

	exception := models.ScheduleException{
		CenterID:      centerID,
		RuleID:        ruleID,
		OriginalDate:  req.OriginalDate,
		ExceptionType: req.Type,
		Status:        "PENDING",
		NewStartAt:    req.NewStartAt,
		NewEndAt:      req.NewEndAt,
		NewTeacherID:  req.NewTeacherID,
		Reason:        req.Reason,
	}

	if req.NewTeacherName != "" {
		exception.Reason = fmt.Sprintf("[代課老師：%s] %s", req.NewTeacherName, req.Reason)
	}

	var createdException models.ScheduleException
	var teacherName, centerName string

	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var createErr error
		createdException, createErr = s.exceptionRepo.CreateWithDB(ctx, tx, exception)
		if createErr != nil {
			return fmt.Errorf("failed to create exception: %w", createErr)
		}

		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "CREATE_EXCEPTION",
			TargetType: "ScheduleException",
			TargetID:   createdException.ID,
			Payload:    models.AuditPayload{After: exception},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return models.ScheduleException{}, nil, txErr
	}

	teacher, _ := s.teacherRepo.GetByID(ctx, teacherID)
	teacherName = teacher.Name
	if teacherName == "" {
		teacherName = "老師"
	}

	center, _ := s.centerRepo.GetByID(ctx, centerID)
	centerName = center.Name

	createdException.ExceptionType = req.Type

	if s.notificationQueue != nil {
		_ = s.notificationQueue.NotifyExceptionSubmittedSync(ctx, &createdException, teacherName, centerName)
	}

	s.invalidateRelatedCaches(ctx, &createdException)

	return createdException, nil, nil
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

	s.invalidateRelatedCaches(ctx, &exception)

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

	status := action
	if action == "APPROVE" || action == "APPROVED" {
		status = "APPROVED"
	} else if action == "REJECT" || action == "REJECTED" {
		status = "REJECTED"
	}

	exception.Status = status
	exception.ReviewedBy = &adminID
	now := time.Now()
	exception.ReviewedAt = &now
	exception.ReviewNote = reason

	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if status == "APPROVED" {
			rule, err := s.ruleRepo.GetByID(ctx, exception.RuleID)
			if err != nil {
				return fmt.Errorf("failed to get rule: %w", err)
			}

			var startAt, endAt time.Time
			if exception.ExceptionType == "RESCHEDULE" && exception.NewStartAt != nil {
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
				nil,
				nil,
			)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
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

			if err := s.applyExceptionChangesWithTx(ctx, tx, &exception, &rule); err != nil {
				return fmt.Errorf("failed to apply exception changes: %w", err)
			}
		}

		if err := tx.Omit("created_at").Save(&exception).Error; err != nil {
			return fmt.Errorf("failed to update exception: %w", err)
		}

		auditLog := models.AuditLog{
			CenterID:   exception.CenterID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "REVIEW_EXCEPTION_" + action,
			TargetType: "ScheduleException",
			TargetID:   exceptionID,
			Payload:    models.AuditPayload{Before: oldStatus, After: action},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	if s.notificationQueue != nil {
		rule, _ := s.ruleRepo.GetByID(ctx, exception.RuleID)
		if rule.ID > 0 && rule.TeacherID != nil {
			teacher, _ := s.teacherRepo.GetByID(ctx, *rule.TeacherID)

			approved := status == "APPROVED"
			if teacher.ID > 0 {
				_ = s.notificationQueue.NotifyExceptionResultSync(ctx, &exception, &teacher, approved, reason)
			}
		}
	}

	s.invalidateRelatedCaches(ctx, &exception)

	return nil
}

func (s *ScheduleExceptionServiceImpl) applyExceptionChanges(ctx context.Context, exception *models.ScheduleException, rule *models.ScheduleRule) error {
	switch exception.ExceptionType {
	case "CANCEL":
	case "RESCHEDULE":
		cutoffDate := exception.OriginalDate.AddDate(0, 0, -1)
		rule.EffectiveRange.EndDate = cutoffDate
		if err := s.ruleRepo.Update(ctx, *rule); err != nil {
			return fmt.Errorf("failed to truncate original rule: %w", err)
		}

		newWeekday := int(exception.NewStartAt.Weekday())
		if newWeekday == 0 {
			newWeekday = 7
		}

		newRule := models.ScheduleRule{
			CenterID:   rule.CenterID,
			OfferingID: rule.OfferingID,
			TeacherID:  rule.TeacherID,
			RoomID:     rule.RoomID,
			Name:       rule.Name,
			Weekday:    newWeekday,
			StartTime:  exception.NewStartAt.Format("15:04"),
			EndTime:    exception.NewEndAt.Format("15:04"),
			Duration:   int(exception.NewEndAt.Sub(*exception.NewStartAt).Minutes()),
			EffectiveRange: models.DateRange{
				StartDate: *exception.NewStartAt,
				EndDate:   rule.EffectiveRange.EndDate,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if _, err := s.ruleRepo.Create(ctx, newRule); err != nil {
			return fmt.Errorf("failed to create reschedule rule: %w", err)
		}

	case "REPLACE_TEACHER":
		if exception.NewTeacherID != nil {
			rule.TeacherID = exception.NewTeacherID
			rule.UpdatedAt = time.Now()
			if err := s.ruleRepo.Update(ctx, *rule); err != nil {
				return fmt.Errorf("failed to update teacher for substitution: %w", err)
			}
		}
	}

	return nil
}

func (s *ScheduleExceptionServiceImpl) applyExceptionChangesWithTx(ctx context.Context, tx *gorm.DB, exception *models.ScheduleException, rule *models.ScheduleRule) error {
	switch exception.ExceptionType {
	case "CANCEL":
	case "RESCHEDULE":
		cutoffDate := exception.OriginalDate.AddDate(0, 0, -1)
		rule.EffectiveRange.EndDate = cutoffDate
		if err := tx.Omit("created_at").Save(rule).Error; err != nil {
			return fmt.Errorf("failed to truncate original rule: %w", err)
		}

		newWeekday := int(exception.NewStartAt.Weekday())
		if newWeekday == 0 {
			newWeekday = 7
		}

		newRule := models.ScheduleRule{
			CenterID:   rule.CenterID,
			OfferingID: rule.OfferingID,
			TeacherID:  rule.TeacherID,
			RoomID:     rule.RoomID,
			Name:       rule.Name,
			Weekday:    newWeekday,
			StartTime:  exception.NewStartAt.Format("15:04"),
			EndTime:    exception.NewEndAt.Format("15:04"),
			Duration:   int(exception.NewEndAt.Sub(*exception.NewStartAt).Minutes()),
			EffectiveRange: models.DateRange{
				StartDate: *exception.NewStartAt,
				EndDate:   rule.EffectiveRange.EndDate,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := tx.Create(&newRule).Error; err != nil {
			return fmt.Errorf("failed to create reschedule rule: %w", err)
		}

	case "REPLACE_TEACHER":
		if exception.NewTeacherID != nil {
			rule.TeacherID = exception.NewTeacherID
			rule.UpdatedAt = time.Now()
			if err := tx.Omit("created_at").Save(rule).Error; err != nil {
				return fmt.Errorf("failed to update teacher for substitution: %w", err)
			}
		}
	}

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
	err := s.App.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("original_date >= ?", startDate).
		Where("original_date <= ?", endDate).
		Where("status IN ('PENDING', 'APPROVED', 'REJECTED')").
		Find(&exceptions).Error

	return exceptions, err
}

func (s *ScheduleExceptionServiceImpl) GetPendingExceptions(ctx context.Context, centerID uint) ([]models.ScheduleException, error) {
	var exceptions []models.ScheduleException
	err := s.App.MySQL.RDB.WithContext(ctx).
		Preload("Rule").
		Preload("Rule.Teacher").
		Preload("Rule.Room").
		Where("center_id = ?", centerID).
		Where("status = ?", "PENDING").
		Order("created_at ASC").
		Find(&exceptions).Error

	return exceptions, err
}

func (s *ScheduleExceptionServiceImpl) GetAllExceptions(ctx context.Context, centerID uint, status string) ([]models.ScheduleException, error) {
	var exceptions []models.ScheduleException
	query := s.App.MySQL.RDB.WithContext(ctx).
		Preload("Rule").
		Preload("Rule.Teacher").
		Preload("Rule.Room").
		Where("center_id = ?", centerID)

	if status != "" {
		if status == "APPROVED" || status == "APPROVE" {
			query = query.Where("status IN (?, ?)", "APPROVED", "APPROVE")
		} else if status == "REJECTED" || status == "REJECT" {
			query = query.Where("status IN (?, ?)", "REJECTED", "REJECT")
		} else {
			query = query.Where("status = ?", status)
		}
	}

	err := query.Order("created_at DESC").Find(&exceptions).Error
	return exceptions, err
}

func (s *ScheduleExceptionServiceImpl) invalidateRelatedCaches(ctx context.Context, exception *models.ScheduleException) {
	rule, err := s.ruleRepo.GetByID(ctx, exception.RuleID)
	if err != nil {
		s.Logger.Warn("failed to get rule for cache invalidation", "error", err, "exception_id", exception.ID)
		return
	}

	pattern := fmt.Sprintf("schedule:expand:center:%d:*", exception.CenterID)
	_ = s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, pattern)

	if rule.ID != 0 && rule.TeacherID != nil {
		teacherPattern := fmt.Sprintf("schedule:expand:teacher:%d:center:%d:*", *rule.TeacherID, exception.CenterID)
		_ = s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, teacherPattern)
	}

	s.Logger.Info("exception-related cache invalidated",
		"center_id", exception.CenterID,
		"exception_id", exception.ID,
		"rule_id", exception.RuleID,
		"teacher_id", func() uint {
			if rule.ID != 0 && rule.TeacherID != nil {
				return *rule.TeacherID
			}
			return 0
		}())
}
