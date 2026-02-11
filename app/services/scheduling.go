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
	"timeLedger/app/requests"
	"timeLedger/app/resources"
	"timeLedger/global/errInfos"
	"timeLedger/libs"
)

// ScheduleServiceInterface 排課服務接口
type ScheduleServiceInterface interface {
	// 衝突檢查
	CheckOverlap(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, weekday int, excludeRuleID *uint) (*OverlapCheckResult, error)
	CheckTeacherBuffer(ctx context.Context, centerID, teacherID uint, prevEndTime, nextStartTime time.Time, courseID uint) (*BufferCheckResult, error)
	CheckRoomBuffer(ctx context.Context, centerID, roomID uint, prevEndTime, nextStartTime time.Time, courseID uint) (*BufferCheckResult, error)
	// ValidateFull 完整驗證
	// 如果 prevEndTime 和 nextStartTime 為 nil，系統會自動計算上一堂課的結束時間
	ValidateFull(ctx context.Context, centerID uint, teacherID *uint, roomID, courseID uint, startTime, endTime time.Time, excludeRuleID *uint, allowBufferOverride bool, prevEndTime, nextStartTime *time.Time) (*FullValidationResult, error)

	// 規則管理
	GetRules(ctx context.Context, centerID uint) ([]models.ScheduleRule, error)
	CreateRule(ctx context.Context, centerID, adminID uint, req *CreateScheduleRuleRequest) ([]models.ScheduleRule, *errInfos.Res, error)
	UpdateRule(ctx context.Context, centerID, adminID, ruleID uint, req *UpdateScheduleRuleRequest) ([]models.ScheduleRule, *errInfos.Res, error)
	DeleteRule(ctx context.Context, centerID, adminID, ruleID uint) error

	// 例外管理
	CreateException(ctx context.Context, centerID, teacherID, ruleID uint, req *CreateExceptionRequest) (models.ScheduleException, *errInfos.Res, error)
	ReviewException(ctx context.Context, exceptionID, adminID uint, req *ReviewExceptionRequest) error
	GetExceptionsByRule(ctx context.Context, ruleID uint) ([]models.ScheduleException, error)
	GetExceptionsByDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]models.ScheduleException, error)
	GetPendingExceptions(ctx context.Context, centerID uint) ([]models.ScheduleException, error)
	GetAllExceptions(ctx context.Context, centerID uint, status string) ([]models.ScheduleException, error)

	// 展開與摘要
	ExpandRules(ctx context.Context, centerID uint, req *ExpandRulesRequest) ([]ExpandedSchedule, error)
	GetCachedExpandedSchedules(ctx context.Context, centerID uint, req *ExpandRulesRequest) ([]ExpandedSchedule, error)
	GetCachedTeacherSchedule(ctx context.Context, teacherID, centerID uint, startDate, endDate time.Time) ([]ExpandedSchedule, error)
	GetTodaySummary(ctx context.Context, centerID uint) (*TodaySummary, error)
	DetectPhaseTransitions(ctx context.Context, centerID, offeringID uint, startDate, endDate time.Time) ([]PhaseTransition, error)
	CheckRuleLockStatus(ctx context.Context, centerID, ruleID uint, exceptionDate time.Time) (*RuleLockStatus, error)

	// 快取管理
	InvalidateCenterScheduleCache(ctx context.Context, centerID uint) error
	InvalidateTeacherScheduleCache(ctx context.Context, teacherID, centerID uint) error
	InvalidateExceptionRelatedCache(ctx context.Context, centerID uint, exception *models.ScheduleException) error

	// 矩陣視圖
	GetMatrixViewData(ctx context.Context, centerID uint, req *requests.MatrixViewRequest) (*resources.MatrixViewResponse, error)
}

// CreateScheduleRuleRequest 建立排課規則請求
type CreateScheduleRuleRequest struct {
	Name           string  `json:"name" binding:"required"`
	OfferingID     uint    `json:"offering_id" binding:"required"`
	TeacherID      *uint   `json:"teacher_id"`
	RoomID         uint    `json:"room_id" binding:"required"`
	StartTime      string  `json:"start_time" binding:"required,time_format"`
	EndTime        string  `json:"end_time" binding:"required,time_format"`
	Duration       int     `json:"duration" binding:"required"`
	Weekdays       []int   `json:"weekdays" binding:"required,min=1"`
	StartDate      string  `json:"start_date" binding:"required,date_format"`
	EndDate        *string `json:"end_date"`
	OverrideBuffer bool    `json:"override_buffer"`
}

// UpdateScheduleRuleRequest 更新排課規則請求
type UpdateScheduleRuleRequest struct {
	Name            string     `json:"name"`
	OfferingID      uint       `json:"offering_id"`
	TeacherID       *uint      `json:"teacher_id"`
	RoomID          uint       `json:"room_id"`
	StartTime       string     `json:"start_time"`
	EndTime         string     `json:"end_time"`
	Duration        int        `json:"duration"`
	Weekdays        []int      `json:"weekdays"`
	StartDate       string     `json:"start_date"`
	EndDate         *string    `json:"end_date"`
	SuspendedDates  []string   `json:"suspended_dates"` // 停課日期列表
	UpdateMode      string     `json:"update_mode"`
	ExcludeRuleID   *uint      `json:"exclude_rule_id"` // 更新時排除自己，避免與自己衝突
}

// CreateExceptionRequest 建立例外請求
type CreateExceptionRequest struct {
	RuleID         uint       `json:"rule_id" binding:"required"`
	OriginalDate   time.Time  `json:"original_date" binding:"required"` // 格式: YYYY-MM-DD，服務層解析
	Type           string     `json:"type" binding:"required"`
	NewStartAt     *time.Time `json:"new_start_at"` // 服務層解析
	NewEndAt       *time.Time `json:"new_end_at"`   // 服務層解析
	NewTeacherID   *uint      `json:"new_teacher_id"`
	NewTeacherName string     `json:"new_teacher_name"`
	NewRoomID      *uint      `json:"new_room_id"`
	Reason         string     `json:"reason" binding:"required"`
}

// ReviewExceptionRequest 審核例外請求
type ReviewExceptionRequest struct {
	Action         string `json:"action" binding:"required"`
	OverrideBuffer bool   `json:"override_buffer"`
	Reason         string `json:"reason"`
}

// ExpandRulesRequest 展開規則請求
type ExpandRulesRequest struct {
	RuleIDs   []uint    `json:"rule_ids"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

// OverlapCheckResult 重疊檢查結果
type OverlapCheckResult = ValidationResult

// BufferCheckResult 緩衝時間檢查結果
type BufferCheckResult = ValidationResult

// FullValidationResult 完整驗證結果
type FullValidationResult = ValidationResult

// BufferConflictDetail 緩衝衝突詳細資訊
type BufferConflictDetail = ValidationConflict

// RuleLockStatus 規則鎖定狀態
type RuleLockStatus struct {
	IsLocked      bool       `json:"is_locked"`
	LockReason    string     `json:"lock_reason,omitempty"`
	LockAt        *time.Time `json:"lock_at,omitempty"`
	Deadline      time.Time  `json:"deadline"`
	DaysRemaining int        `json:"days_remaining"`
}

// TodaySummary 今日摘要
type TodaySummary struct {
	Sessions               []TodaySession `json:"sessions"`
	TotalSessions          int            `json:"total_sessions"`
	CompletedSessions      int            `json:"completed_sessions"`
	InProgressSessions     int            `json:"in_progress_sessions"`
	UpcomingSessions       int            `json:"upcoming_sessions"`
	InProgressTeacherNames []string       `json:"in_progress_teacher_names,omitempty"`
	PendingExceptions      int            `json:"pending_exceptions"`
	ChangesCount           int            `json:"changes_count"`
	HasScheduleChanges     bool           `json:"has_schedule_changes"`
}

// TodaySession 今日課程場次
type TodaySession struct {
	ID        uint          `json:"id"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Offering  TodayOffering `json:"offering"`
	Teacher   TodayTeacher  `json:"teacher"`
	Room      TodayRoom     `json:"room"`
	Status    string        `json:"status"`
}

// TodayOffering 今日課程
type TodayOffering struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TodayTeacher 今日老師
type TodayTeacher struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TodayRoom 今日教室
type TodayRoom struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// UpdateMode 更新模式常量
const (
	UpdateModeSingle = "SINGLE"
	UpdateModeFuture = "FUTURE"
	UpdateModeAll    = "ALL"
)

// ScheduleService 排課服務實現
type ScheduleService struct {
	BaseService
	ruleRepo          *repositories.ScheduleRuleRepository
	exceptionRepo     *repositories.ScheduleExceptionRepository
	offeringRepo      *repositories.OfferingRepository
	courseRepo        *repositories.CourseRepository
	centerRepo        *repositories.CenterRepository
	roomRepo          *repositories.RoomRepository
	auditLogRepo      *repositories.AuditLogRepository
	holidayRepo       *repositories.CenterHolidayRepository
	teacherRepo       *repositories.TeacherRepository
	validationSvc     ScheduleValidationService
	expansionSvc      ScheduleExpansionService
	exceptionSvc      ScheduleExceptionService
	notificationSvc   NotificationService
	notificationQueue NotificationQueueService
	cacheSvc          *CacheService
}

// NewScheduleService 建立排課服務
func NewScheduleService(app *app.App) ScheduleServiceInterface {
	baseSvc := NewBaseService(app, "ScheduleService")
	svc := &ScheduleService{
		BaseService:       *baseSvc,
		ruleRepo:          repositories.NewScheduleRuleRepository(app),
		exceptionRepo:     repositories.NewScheduleExceptionRepository(app),
		offeringRepo:      repositories.NewOfferingRepository(app),
		courseRepo:        repositories.NewCourseRepository(app),
		centerRepo:        repositories.NewCenterRepository(app),
		roomRepo:          repositories.NewRoomRepository(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
		holidayRepo:       repositories.NewCenterHolidayRepository(app),
		teacherRepo:       repositories.NewTeacherRepository(app),
		validationSvc:     NewScheduleValidationService(app),
		expansionSvc:      NewScheduleExpansionService(app),
		exceptionSvc:      NewScheduleExceptionService(app),
		notificationSvc:   NewNotificationService(app),
		notificationQueue: NewNotificationQueueService(app),
		cacheSvc:          NewCacheService(app),
	}
	return svc
}

// GetRules 取得排課規則列表
func (s *ScheduleService) GetRules(ctx context.Context, centerID uint) ([]models.ScheduleRule, error) {
	rules, err := s.ruleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// CreateRule 建立排課規則（使用交易確保原子性）
func (s *ScheduleService) CreateRule(ctx context.Context, centerID, adminID uint, req *CreateScheduleRuleRequest) ([]models.ScheduleRule, *errInfos.Res, error) {
	// 解析日期（使用台灣時區）
	loc := libs.GetTaiwanLocation()
	startDate, err := time.ParseInLocation("2006-01-02", req.StartDate, loc)
	if err != nil {
		return nil, s.App.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid start_date format: %w", err)
	}

	var endDate time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		endDate, err = time.ParseInLocation("2006-01-02", *req.EndDate, loc)
		if err != nil {
			return nil, s.App.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid end_date format: %w", err)
		}
	} else {
		endDate = time.Date(2099, 12, 31, 0, 0, 0, 0, loc)
	}

	// 檢查衝突
	startTimeParsed, _ := time.Parse("15:04", req.StartTime)
	endTimeParsed, _ := time.Parse("15:04", req.EndTime)

	// 使用請求中的第一個 weekday 進行重疊檢查
	var checkWeekday int
	if len(req.Weekdays) > 0 {
		checkWeekday = req.Weekdays[0]
	} else {
		checkWeekday = int(startTimeParsed.Weekday())
		if checkWeekday == 0 {
			checkWeekday = 7
		}
	}

	validationResult, err := s.validationSvc.CheckOverlap(ctx, centerID, req.TeacherID, req.RoomID, startTimeParsed, endTimeParsed, checkWeekday, nil)
	if err != nil {
		return nil, s.App.Err.New(errInfos.SQL_ERROR), fmt.Errorf("failed to check overlap: %w", err)
	}

	if !validationResult.Valid {
		return nil, s.App.Err.New(errInfos.SCHED_OVERLAP), fmt.Errorf("time slot conflict with existing rules or personal events")
	}

	// 取得課程設定
	offering, err := s.offeringRepo.GetByID(ctx, req.OfferingID)
	if err != nil {
		return nil, s.App.Err.New(errInfos.NOT_FOUND), fmt.Errorf("failed to get offering: %w", err)
	}

	// 檢查 Buffer
	bufferConflicts, err := s.checkBufferConflicts(ctx, centerID, req, &offering, startDate)
	if err != nil {
		return nil, s.App.Err.New(errInfos.SCHED_BUFFER), err
	}

	if len(bufferConflicts) > 0 {
		allOverridable := true
		for _, conflict := range bufferConflicts {
			if !conflict.CanOverride {
				allOverridable = false
				break
			}
		}

		if !req.OverrideBuffer || !allOverridable {
			return nil, s.App.Err.New(errInfos.SCHED_BUFFER), fmt.Errorf("insufficient buffer time")
		}
	}

	// 使用交易建立規則和審核日誌
	var createdRules []models.ScheduleRule

	// 使用 Repository 的 Transaction 方法，確保所有操作都在同一交易中
	txErr := s.ruleRepo.Transaction(ctx, func(txRepo *repositories.ScheduleRuleRepository) error {
		// 在交易中建立規則（使用 txRepo，它擁有交易連接）
		for _, weekday := range req.Weekdays {
			rule := models.ScheduleRule{
				CenterID:   centerID,
				OfferingID: req.OfferingID,
				TeacherID:  req.TeacherID,
				RoomID:     req.RoomID,
				Weekday:    weekday,
				StartTime:  req.StartTime,
				EndTime:    req.EndTime,
				Duration:   req.Duration,
				EffectiveRange: models.DateRange{
					StartDate: startDate,
					EndDate:   endDate,
				},
			}

			if _, err := txRepo.Create(ctx, rule); err != nil {
				return fmt.Errorf("failed to create schedule rule: %w", err)
			}
			createdRules = append(createdRules, rule)
		}

		// 在交易中記錄審核日誌
		// 從 txRepo 取得交易連接並傳遞給 auditLogRepo
		txDB := txRepo.GetDBWrite()
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "CREATE_SCHEDULE_RULE",
			TargetType: "ScheduleRule",
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"rules_created": len(createdRules),
				},
			},
		}
		if _, err := s.auditLogRepo.CreateWithTxDB(ctx, txDB, auditLog); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, s.App.Err.New(errInfos.ERR_TX_FAILED), txErr
	}

	// 建立規則後，使相關課表快取失效
	_ = s.InvalidateCenterScheduleCache(ctx, centerID)

	return createdRules, nil, nil
}

// checkBufferConflicts 檢查緩衝時間衝突
func (s *ScheduleService) checkBufferConflicts(ctx context.Context, centerID uint, req *CreateScheduleRuleRequest, offering *models.Offering, startDate time.Time) ([]BufferConflictDetail, error) {
	var conflicts []BufferConflictDetail

	// 計算檢查日期
	loc := app.GetTaiwanLocation()
	checkDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)

	for _, weekday := range req.Weekdays {
		current := checkDate
		weekdayDiff := weekday - int(current.Weekday())
		if weekdayDiff <= 0 {
			weekdayDiff += 7
		}
		targetDate := current.AddDate(0, 0, weekdayDiff)

		// 產生新课程的时间
		targetDateLocal := time.Date(
			targetDate.Year(), targetDate.Month(), targetDate.Day(),
			0, 0, 0, 0, loc,
		)
		newStartTime, _ := time.ParseInLocation("2006-01-02 15:04", targetDateLocal.Format("2006-01-02")+" "+req.StartTime, loc)

		// 檢查 Teacher Buffer
		if req.TeacherID != nil && *req.TeacherID > 0 {
			prevEndTime, _ := s.getTeacherPreviousSessionEndTime(ctx, centerID, *req.TeacherID, weekday, req.StartTime, newStartTime)
			if !prevEndTime.IsZero() {
				result, err := s.validationSvc.CheckTeacherBuffer(ctx, centerID, *req.TeacherID, prevEndTime, newStartTime, offering.CourseID)
				if err != nil {
					return nil, err
				}

				for _, conflict := range result.Conflicts {
					conflicts = append(conflicts, BufferConflictDetail{
						Type:            conflict.Type,
						Message:         conflict.Message,
						RequiredMinutes: conflict.RequiredMinutes,
						DiffMinutes:     conflict.DiffMinutes,
						CanOverride:     conflict.CanOverride,
					})
				}
			}
		}

		// 檢查 Room Buffer
		if req.RoomID > 0 {
			prevRoomEndTime, _ := s.getRoomPreviousSessionEndTime(ctx, centerID, req.RoomID, weekday, req.StartTime, newStartTime)
			if !prevRoomEndTime.IsZero() {
				result, err := s.validationSvc.CheckRoomBuffer(ctx, centerID, req.RoomID, prevRoomEndTime, newStartTime, offering.CourseID)
				if err != nil {
					return nil, err
				}

				for _, conflict := range result.Conflicts {
					conflicts = append(conflicts, BufferConflictDetail{
						Type:            conflict.Type,
						Message:         conflict.Message,
						RequiredMinutes: conflict.RequiredMinutes,
						DiffMinutes:     conflict.DiffMinutes,
						CanOverride:     conflict.CanOverride,
					})
				}
			}
		}
	}

	return conflicts, nil
}

// getTeacherPreviousSessionEndTime 取得老師在指定 weekday 之前的最後一筆課程結束時間
func (s *ScheduleService) getTeacherPreviousSessionEndTime(ctx context.Context, centerID, teacherID uint, weekday int, beforeTimeStr string, newStartTime time.Time) (time.Time, error) {
	rule, err := s.ruleRepo.GetLastSessionByTeacherAndWeekday(ctx, centerID, teacherID, weekday, beforeTimeStr)
	if err != nil {
		return time.Time{}, nil
	}

	if rule == nil {
		return time.Time{}, nil
	}

	loc := app.GetTaiwanLocation()
	endTimeParts := splitTime(rule.EndTime)
	if len(endTimeParts) >= 2 {
		hour := endTimeParts[0]
		minute := endTimeParts[1]

		endTime := time.Date(
			newStartTime.Year(), newStartTime.Month(), newStartTime.Day(),
			hour, minute, 0, 0, loc,
		)
		return endTime, nil
	}

	return time.Time{}, nil
}

// getRoomPreviousSessionEndTime 取得教室在指定 weekday 之前的最後一筆課程結束時間
func (s *ScheduleService) getRoomPreviousSessionEndTime(ctx context.Context, centerID, roomID uint, weekday int, beforeTimeStr string, newStartTime time.Time) (time.Time, error) {
	rule, err := s.ruleRepo.GetLastSessionByRoomAndWeekday(ctx, centerID, roomID, weekday, beforeTimeStr)
	if err != nil {
		return time.Time{}, nil
	}

	if rule == nil {
		return time.Time{}, nil
	}

	loc := app.GetTaiwanLocation()
	endTimeParts := splitTime(rule.EndTime)
	if len(endTimeParts) >= 2 {
		hour := endTimeParts[0]
		minute := endTimeParts[1]

		endTime := time.Date(
			newStartTime.Year(), newStartTime.Month(), newStartTime.Day(),
			hour, minute, 0, 0, loc,
		)
		return endTime, nil
	}

	return time.Time{}, nil
}

// UpdateRule 更新排課規則（使用交易確保原子性）
func (s *ScheduleService) UpdateRule(ctx context.Context, centerID, adminID, ruleID uint, req *UpdateScheduleRuleRequest) ([]models.ScheduleRule, *errInfos.Res, error) {
	// 檢查規則是否存在
	existingRule, err := s.ruleRepo.GetByIDAndCenterID(ctx, ruleID, centerID)
	if err != nil {
		return nil, s.App.Err.New(errInfos.NOT_FOUND), fmt.Errorf("rule not found")
	}

	// 解析日期
	var startDate, endDate time.Time
	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, s.App.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid start_date format")
		}
	}

	if req.EndDate != nil && *req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, s.App.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid end_date format")
		}
	}

	// 衝突檢查（使用 req.ExcludeRuleID 排除自己）
	if req.StartTime != "" && req.EndTime != "" && req.RoomID != 0 {
		startTimeParsed, _ := time.Parse("15:04", req.StartTime)
		endTimeParsed, _ := time.Parse("15:04", req.EndTime)

		// 決定要檢查的星期幾
		var checkWeekday int
		if len(req.Weekdays) > 0 {
			checkWeekday = req.Weekdays[0]
		} else if !startDate.IsZero() {
			checkWeekday = int(startDate.Weekday())
			if checkWeekday == 0 {
				checkWeekday = 7
			}
		} else {
			checkWeekday = existingRule.Weekday
		}

		// 使用 ExcludeRuleID 排除自己，避免與自己衝突
		validationResult, err := s.validationSvc.CheckOverlap(ctx, centerID, req.TeacherID, req.RoomID, startTimeParsed, endTimeParsed, checkWeekday, req.ExcludeRuleID)
		if err != nil {
			return nil, s.App.Err.New(errInfos.SQL_ERROR), fmt.Errorf("failed to check overlap: %w", err)
		}

		if !validationResult.Valid {
			return nil, s.App.Err.New(errInfos.SCHED_OVERLAP), fmt.Errorf("time slot conflict with existing rules")
		}
	}

	// 取得所有相關規則
	allRules, err := s.ruleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, s.App.Err.New(errInfos.SQL_ERROR), fmt.Errorf("failed to fetch rules: %w", err)
	}

	var relatedRules []models.ScheduleRule
	for _, rule := range allRules {
		if rule.OfferingID == existingRule.OfferingID &&
			rule.StartTime == existingRule.StartTime &&
			rule.EndTime == existingRule.EndTime &&
			rule.ID != existingRule.ID {
			relatedRules = append(relatedRules, rule)
		}
	}

	// 使用交易處理更新操作
	var resultRules []models.ScheduleRule

	// 使用 Repository 的 Transaction 方法，確保所有操作都在同一交易中
	txErr := s.ruleRepo.Transaction(ctx, func(txRepo *repositories.ScheduleRuleRepository) error {
		var handlerErr error
		var rules []models.ScheduleRule

		// 根據 update_mode 處理
		switch req.UpdateMode {
		case UpdateModeFuture:
			rules, handlerErr = s.handleFutureUpdateWithTx(txRepo, ctx, centerID, existingRule, relatedRules, req, startDate, endDate)
		case UpdateModeSingle:
			rules, handlerErr = s.handleSingleUpdateWithTx(txRepo, ctx, centerID, existingRule, req, startDate)
		default:
			rules, handlerErr = s.handleAllUpdateWithTx(txRepo, ctx, centerID, existingRule, relatedRules, req, startDate, endDate)
		}

		if handlerErr != nil {
			return handlerErr
		}
		resultRules = rules

		// 在交易中記錄審核日誌
		txDB := txRepo.GetDBWrite()
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "UPDATE_SCHEDULE_RULE",
			TargetType: "ScheduleRule",
			TargetID:   ruleID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"rules_updated": len(resultRules),
					"update_mode":   req.UpdateMode,
				},
			},
		}
		if _, err := s.auditLogRepo.CreateWithTxDB(ctx, txDB, auditLog); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, s.App.Err.New(errInfos.ERR_TX_FAILED), txErr
	}

	// 更新規則後，使相關課表快取失效
	_ = s.InvalidateCenterScheduleCache(ctx, centerID)

	return resultRules, nil, nil
}

// handleFutureUpdateWithTx 處理 FUTURE 模式（交易版本）
func (s *ScheduleService) handleFutureUpdateWithTx(txRepo *repositories.ScheduleRuleRepository, ctx context.Context, centerID uint, existingRule models.ScheduleRule, relatedRules []models.ScheduleRule, req *UpdateScheduleRuleRequest, startDate, endDate time.Time) ([]models.ScheduleRule, error) {
	if startDate.IsZero() {
		return nil, fmt.Errorf("start_date is required for FUTURE update mode")
	}

	cutoffDate := startDate.AddDate(0, 0, -1)
	var result []models.ScheduleRule

	// 更新現有規則
	for _, rule := range append(relatedRules, existingRule) {
		rule.EffectiveRange.EndDate = cutoffDate
		if err := txRepo.Update(ctx, rule); err != nil {
			return nil, err
		}
		result = append(result, rule)
	}

	// 建立新規則段
	newRules := s.createNewRuleSegment(centerID, existingRule, relatedRules, req, startDate, endDate)
	for _, newRule := range newRules {
		if _, err := txRepo.Create(ctx, newRule); err != nil {
			return nil, err
		}
		result = append(result, newRule)
	}

	return result, nil
}

// handleSingleUpdateWithTx 處理 SINGLE 模式（交易版本）
func (s *ScheduleService) handleSingleUpdateWithTx(txRepo *repositories.ScheduleRuleRepository, ctx context.Context, centerID uint, existingRule models.ScheduleRule, req *UpdateScheduleRuleRequest, targetDate time.Time) ([]models.ScheduleRule, error) {
	now := time.Now()

	// 取得交易連接以在交易中建立例外單
	txDB := txRepo.GetDBWrite()

	// 建立取消例外單（使用交易連接）
	cancelException := models.ScheduleException{
		CenterID:      centerID,
		RuleID:        existingRule.ID,
		OriginalDate:  targetDate,
		ExceptionType: "CANCEL",
		Status:        "APPROVED",
		Reason:        req.Name,
		ReviewedAt:    &now,
	}

	if _, err := s.exceptionRepo.CreateWithDB(ctx, txDB, cancelException); err != nil {
		return nil, fmt.Errorf("failed to create cancel exception: %w", err)
	}

	// 建立新時間的規則
	weekday := int(targetDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	newRule := models.ScheduleRule{
		CenterID:   centerID,
		OfferingID: req.OfferingID,
		TeacherID:  req.TeacherID,
		RoomID:     req.RoomID,
		Name:       req.Name,
		Weekday:    weekday,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Duration:   req.Duration,
		EffectiveRange: models.DateRange{
			StartDate: targetDate,
			EndDate:   targetDate,
		},
	}

	if _, err := txRepo.Create(ctx, newRule); err != nil {
		return nil, fmt.Errorf("failed to create new rule: %w", err)
	}

	return []models.ScheduleRule{existingRule, newRule}, nil
}

// handleAllUpdateWithTx 處理 ALL 模式（交易版本）
func (s *ScheduleService) handleAllUpdateWithTx(txRepo *repositories.ScheduleRuleRepository, ctx context.Context, centerID uint, existingRule models.ScheduleRule, relatedRules []models.ScheduleRule, req *UpdateScheduleRuleRequest, startDate, endDate time.Time) ([]models.ScheduleRule, error) {
	newWeekdayMap := make(map[int]bool)
	for _, w := range req.Weekdays {
		newWeekdayMap[w] = true
	}

	var updatedRules []models.ScheduleRule
	var deletedRuleIDs []uint

	// 處理現有規則
	for _, rule := range append(relatedRules, existingRule) {
		if newWeekdayMap[rule.Weekday] {
			updatedRule := s.applyUpdateToRule(rule, req, startDate, endDate)
			if err := txRepo.Update(ctx, updatedRule); err != nil {
				return nil, err
			}
			updatedRules = append(updatedRules, updatedRule)
			delete(newWeekdayMap, rule.Weekday)
		} else {
			deletedRuleIDs = append(deletedRuleIDs, rule.ID)
		}
	}

	// 建立新規則
	var createdRules []models.ScheduleRule
	for weekday := range newWeekdayMap {
		newRule := s.createSingleRule(centerID, existingRule, req, weekday, startDate, endDate)
		if _, err := txRepo.Create(ctx, *newRule); err != nil {
			return nil, err
		}
		createdRules = append(createdRules, *newRule)
	}

	// 執行硬刪除（使用交易連接）
	// 先刪除關聯的例外記錄，再刪除規則
	txDB := txRepo.GetDBWrite()
	for _, id := range deletedRuleIDs {
		// 先刪除關聯的例外記錄（避免外鍵約束衝突）
		if err := txDB.WithContext(ctx).Where("rule_id = ?", id).Delete(&models.ScheduleException{}).Error; err != nil {
			return nil, fmt.Errorf("failed to delete associated exceptions: %w", err)
		}
		// 再刪除規則
		if err := txDB.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.ScheduleRule{}).Error; err != nil {
			return nil, err
		}
	}

	return append(updatedRules, createdRules...), nil
}

// applyUpdateToRule 將更新应用到規則
func (s *ScheduleService) applyUpdateToRule(rule models.ScheduleRule, req *UpdateScheduleRuleRequest, startDate, endDate time.Time) models.ScheduleRule {
	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.OfferingID != 0 {
		rule.OfferingID = req.OfferingID
	}
	if req.TeacherID != nil {
		rule.TeacherID = req.TeacherID
	}
	if req.RoomID != 0 {
		rule.RoomID = req.RoomID
	}
	if req.StartTime != "" {
		rule.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		rule.EndTime = req.EndTime
	}
	if req.Duration != 0 {
		rule.Duration = req.Duration
	}
	if !startDate.IsZero() && startDate != rule.EffectiveRange.StartDate {
		rule.EffectiveRange.StartDate = startDate
	}
	if !endDate.IsZero() {
		rule.EffectiveRange.EndDate = endDate
	}
	// 處理 SuspendedDates（停課日期）
	if len(req.SuspendedDates) > 0 {
		loc := libs.GetTaiwanLocation()
		suspendedDates := make(models.SuspendedDates, 0, len(req.SuspendedDates))
		for _, dateStr := range req.SuspendedDates {
			date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
			if err != nil {
				s.Logger.Warn("invalid suspended_date format", "date", dateStr, "error", err)
				continue
			}
			suspendedDates = append(suspendedDates, date)
		}
		rule.SuspendedDates = suspendedDates
	}
	return rule
}

// createNewRuleSegment 為 FUTURE 模式創建新的規則段
func (s *ScheduleService) createNewRuleSegment(centerID uint, existingRule models.ScheduleRule, relatedRules []models.ScheduleRule, req *UpdateScheduleRuleRequest, startDate, endDate time.Time) []models.ScheduleRule {
	weekdaysToCreate := []int{existingRule.Weekday}
	for _, rule := range relatedRules {
		weekdaysToCreate = append(weekdaysToCreate, rule.Weekday)
	}

	var newRules []models.ScheduleRule
	for _, weekday := range weekdaysToCreate {
		newRule := s.createSingleRule(centerID, existingRule, req, weekday, startDate, endDate)
		newRules = append(newRules, *newRule)
	}

	return newRules
}

// createSingleRule 創建單個規則
func (s *ScheduleService) createSingleRule(centerID uint, existingRule models.ScheduleRule, req *UpdateScheduleRuleRequest, weekday int, startDate, endDate time.Time) *models.ScheduleRule {
	effectiveStartDate := startDate
	if startDate.IsZero() {
		effectiveStartDate = existingRule.EffectiveRange.StartDate
	}
	effectiveEndDate := endDate
	if endDate.IsZero() {
		effectiveEndDate = existingRule.EffectiveRange.EndDate
	}

	newRule := &models.ScheduleRule{
		CenterID:       centerID,
		OfferingID:     existingRule.OfferingID,
		TeacherID:      existingRule.TeacherID,
		RoomID:         existingRule.RoomID,
		Name:           existingRule.Name,
		Weekday:        weekday,
		StartTime:      existingRule.StartTime,
		EndTime:        existingRule.EndTime,
		Duration:       existingRule.Duration,
		SuspendedDates: existingRule.SuspendedDates, // 繼承現有的停課日期
		EffectiveRange: models.DateRange{
			StartDate: effectiveStartDate,
			EndDate:   effectiveEndDate,
		},
	}

	if req.Name != "" {
		newRule.Name = req.Name
	}
	if req.OfferingID != 0 {
		newRule.OfferingID = req.OfferingID
	}
	if req.TeacherID != nil {
		newRule.TeacherID = req.TeacherID
	}
	if req.RoomID != 0 {
		newRule.RoomID = req.RoomID
	}
	if req.StartTime != "" {
		newRule.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		newRule.EndTime = req.EndTime
	}
	if req.Duration != 0 {
		newRule.Duration = req.Duration
	}
	// 處理 SuspendedDates（停課日期）
	if len(req.SuspendedDates) > 0 {
		loc := libs.GetTaiwanLocation()
		suspendedDates := make(models.SuspendedDates, 0, len(req.SuspendedDates))
		for _, dateStr := range req.SuspendedDates {
			date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
			if err != nil {
				s.Logger.Warn("invalid suspended_date format", "date", dateStr, "error", err)
				continue
			}
			suspendedDates = append(suspendedDates, date)
		}
		newRule.SuspendedDates = suspendedDates
	}

	return newRule
}

// DeleteRule 刪除排課規則
func (s *ScheduleService) DeleteRule(ctx context.Context, centerID, adminID, ruleID uint) error {
	// 取得規則以獲取老師ID（用於清除老師課表快取）
	rule, _ := s.ruleRepo.GetByID(ctx, ruleID)

	if err := s.ruleRepo.DeleteByIDAndCenterID(ctx, ruleID, centerID); err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "DELETE_SCHEDULE_RULE",
		TargetType: "ScheduleRule",
		TargetID:   ruleID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
			},
		},
	})

	// 刪除規則後，使相關課表快取失效
	_ = s.InvalidateCenterScheduleCache(ctx, centerID)
	if rule.ID != 0 && rule.TeacherID != nil {
		_ = s.InvalidateTeacherScheduleCache(ctx, *rule.TeacherID, centerID)
	}

	return nil
}

// ExpandRules 展開排課規則
func (s *ScheduleService) ExpandRules(ctx context.Context, centerID uint, req *ExpandRulesRequest) ([]ExpandedSchedule, error) {
	rules, err := s.ruleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	var filteredRules []models.ScheduleRule
	if len(req.RuleIDs) > 0 {
		ruleIDSet := make(map[uint]bool)
		for _, id := range req.RuleIDs {
			ruleIDSet[id] = true
		}
		for _, rule := range rules {
			if ruleIDSet[rule.ID] {
				filteredRules = append(filteredRules, rule)
			}
		}
	} else {
		filteredRules = rules
	}

	expandedSchedules := s.expansionSvc.ExpandRules(ctx, filteredRules, req.StartDate, req.EndDate, centerID)
	return expandedSchedules, nil
}

// GetTodaySummary 取得今日摘要
func (s *ScheduleService) GetTodaySummary(ctx context.Context, centerID uint) (*TodaySummary, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, today.Location())

	// 展開今日的排課規則
	rules, err := s.ruleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	allSessions := s.expansionSvc.ExpandRules(ctx, rules, startOfDay, endOfDay, centerID)

	// 取得待審核例外申請
	pendingExceptions, _ := s.exceptionSvc.GetPendingExceptions(ctx, centerID)

	// 取得今日例外申請
	todayExceptions, _ := s.exceptionSvc.GetExceptionsByDateRange(ctx, centerID, startOfDay, endOfDay)

	response := &TodaySummary{
		PendingExceptions:  len(pendingExceptions),
		ChangesCount:       len(todayExceptions),
		HasScheduleChanges: len(todayExceptions) > 0,
		Sessions:           []TodaySession{},
	}

	now := today
	for _, session := range allSessions {
		teacherName := session.TeacherName
		if teacherName == "" {
			teacherName = "未安排老師"
		}

		roomName := session.RoomName
		if roomName == "" {
			roomName = "未安排教室"
		}

		offeringName := session.OfferingName
		if offeringName == "" {
			offeringName = "未知課程"
		}

		startDateTime, endDateTime := s.parseSessionTimes(session)

		// 判斷課程狀態
		status := s.determineSessionStatus(now, startDateTime, endDateTime)

		// 取得老師 ID（處理指針類型）
		var teacherID uint
		if session.TeacherID != nil {
			teacherID = *session.TeacherID
		}

		response.Sessions = append(response.Sessions, TodaySession{
			ID:        session.RuleID,
			StartTime: startDateTime,
			EndTime:   endDateTime,
			Offering:  TodayOffering{ID: session.OfferingID, Name: offeringName},
			Teacher:   TodayTeacher{ID: teacherID, Name: teacherName},
			Room:      TodayRoom{ID: session.RoomID, Name: roomName},
			Status:    status,
		})
	}

	// 計算統計數據
	response.TotalSessions = len(response.Sessions)
	for _, session := range response.Sessions {
		switch session.Status {
		case "completed":
			response.CompletedSessions++
		case "in_progress":
			response.InProgressSessions++
			if session.Teacher.Name != "" && session.Teacher.Name != "未安排老師" {
				found := false
				for _, name := range response.InProgressTeacherNames {
					if name == session.Teacher.Name {
						found = true
						break
					}
				}
				if !found {
					response.InProgressTeacherNames = append(response.InProgressTeacherNames, session.Teacher.Name)
				}
			}
		case "upcoming":
			response.UpcomingSessions++
		}
	}

	return response, nil
}

// parseSessionTimes 解析課程時間
func (s *ScheduleService) parseSessionTimes(session ExpandedSchedule) (time.Time, time.Time) {
	loc := app.GetTaiwanLocation()
	startDateTime := time.Date(
		session.Date.Year(), session.Date.Month(), session.Date.Day(),
		0, 0, 0, 0, loc,
	)
	endDateTime := time.Date(
		session.Date.Year(), session.Date.Month(), session.Date.Day(),
		0, 0, 0, 0, loc,
	)

	startParts := splitTime(session.StartTime)
	if len(startParts) >= 2 {
		startDateTime = startDateTime.Add(time.Duration(startParts[0])*time.Hour + time.Duration(startParts[1])*time.Minute)
	}

	endParts := splitTime(session.EndTime)
	if len(endParts) >= 2 {
		endDateTime = endDateTime.Add(time.Duration(endParts[0])*time.Hour + time.Duration(endParts[1])*time.Minute)
	}

	// 檢查是否為跨日課程
	isCrossDay := endParts[0] < startParts[0]
	if isCrossDay {
		endDateTime = endDateTime.Add(24 * time.Hour)
	}

	return startDateTime, endDateTime
}

// determineSessionStatus 判斷課程狀態
func (s *ScheduleService) determineSessionStatus(now, start, end time.Time) string {
	if now.After(end) {
		return "completed"
	} else if now.After(start) && now.Before(end) {
		return "in_progress"
	}
	return "upcoming"
}

// DetectPhaseTransitions 偵測階段轉換
func (s *ScheduleService) DetectPhaseTransitions(ctx context.Context, centerID, offeringID uint, startDate, endDate time.Time) ([]PhaseTransition, error) {
	return s.expansionSvc.DetectPhaseTransitions(ctx, centerID, offeringID, startDate, endDate)
}

// CheckRuleLockStatus 檢查規則鎖定狀態
func (s *ScheduleService) CheckRuleLockStatus(ctx context.Context, centerID, ruleID uint, exceptionDate time.Time) (*RuleLockStatus, error) {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("rule not found")
	}

	if rule.CenterID != centerID {
		return nil, fmt.Errorf("rule does not belong to this center")
	}

	now := time.Now()
	response := &RuleLockStatus{
		DaysRemaining: -1,
	}

	if rule.LockAt != nil && now.After(*rule.LockAt) {
		response.IsLocked = true
		response.LockReason = "已超過異動截止日"
		response.LockAt = rule.LockAt
	}

	center, _ := s.centerRepo.GetByID(ctx, centerID)
	leadDays := center.Settings.ExceptionLeadDays
	if leadDays <= 0 {
		leadDays = 14
	}

	deadline := exceptionDate.AddDate(0, 0, -leadDays)
	response.Deadline = deadline
	daysRemaining := int(deadline.Sub(now).Hours() / 24)
	response.DaysRemaining = daysRemaining

	if daysRemaining < 0 && !response.IsLocked {
		response.IsLocked = true
		response.LockReason = "已超過異動截止日（需提前 " + intToString(leadDays) + " 天申請）"
	}

	return response, nil
}

// intToString 整數轉字串
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

// parseInt 從字串中提取數字並轉換為整數
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var num int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		}
	}
	return num
}

// 衝突檢查方法 - 代理到 validationService
func (s *ScheduleService) CheckOverlap(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, weekday int, excludeRuleID *uint) (*OverlapCheckResult, error) {
	result, err := s.validationSvc.CheckOverlap(ctx, centerID, teacherID, roomID, startTime, endTime, weekday, excludeRuleID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ScheduleService) CheckTeacherBuffer(ctx context.Context, centerID, teacherID uint, prevEndTime, nextStartTime time.Time, courseID uint) (*BufferCheckResult, error) {
	result, err := s.validationSvc.CheckTeacherBuffer(ctx, centerID, teacherID, prevEndTime, nextStartTime, courseID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ScheduleService) CheckRoomBuffer(ctx context.Context, centerID, roomID uint, prevEndTime, nextStartTime time.Time, courseID uint) (*BufferCheckResult, error) {
	result, err := s.validationSvc.CheckRoomBuffer(ctx, centerID, roomID, prevEndTime, nextStartTime, courseID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ScheduleService) ValidateFull(ctx context.Context, centerID uint, teacherID *uint, roomID, courseID uint, startTime, endTime time.Time, excludeRuleID *uint, allowBufferOverride bool, prevEndTime, nextStartTime *time.Time) (*FullValidationResult, error) {
	result, err := s.validationSvc.ValidateFull(ctx, centerID, teacherID, roomID, courseID, startTime, endTime, excludeRuleID, allowBufferOverride, prevEndTime, nextStartTime)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 例外管理方法 - 代理到 exceptionService
func (s *ScheduleService) CreateException(ctx context.Context, centerID, teacherID, ruleID uint, req *CreateExceptionRequest) (models.ScheduleException, *errInfos.Res, error) {
	return s.exceptionSvc.CreateException(ctx, centerID, teacherID, ruleID, req)
}

func (s *ScheduleService) ReviewException(ctx context.Context, exceptionID, adminID uint, req *ReviewExceptionRequest) error {
	return s.exceptionSvc.ReviewException(ctx, exceptionID, adminID, req.Action, req.OverrideBuffer, req.Reason)
}

func (s *ScheduleService) GetExceptionsByRule(ctx context.Context, ruleID uint) ([]models.ScheduleException, error) {
	return s.exceptionSvc.GetExceptionsByRule(ctx, ruleID)
}

func (s *ScheduleService) GetExceptionsByDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]models.ScheduleException, error) {
	return s.exceptionSvc.GetExceptionsByDateRange(ctx, centerID, startDate, endDate)
}

func (s *ScheduleService) GetPendingExceptions(ctx context.Context, centerID uint) ([]models.ScheduleException, error) {
	return s.exceptionSvc.GetPendingExceptions(ctx, centerID)
}

func (s *ScheduleService) GetAllExceptions(ctx context.Context, centerID uint, status string) ([]models.ScheduleException, error) {
	return s.exceptionSvc.GetAllExceptions(ctx, centerID, status)
}

// =====================================================
// 課表快取方法 (Schedule Expansion Caching)
// =====================================================

// scheduleCacheKey 產生課表展開快取鍵
// 格式: schedule:expand:center:{center_id}:start:{start_date}:end:{end_date}
func (s *ScheduleService) scheduleCacheKey(centerID uint, startDate, endDate time.Time) string {
	return fmt.Sprintf("schedule:expand:center:%d:start:%s:end:%s",
		centerID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}

// teacherScheduleCacheKey 產生老師課表快取鍵
// 格式: schedule:expand:teacher:{teacher_id}:center:{center_id}:start:{start_date}:end:{end_date}
func (s *ScheduleService) teacherScheduleCacheKey(teacherID, centerID uint, startDate, endDate time.Time) string {
	return fmt.Sprintf("schedule:expand:teacher:%d:center:%d:start:%s:end:%s",
		teacherID, centerID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}

// GetCachedExpandedSchedules 取得帶快取的課表展開結果
// 先檢查快取，若未命中則展開後快取
func (s *ScheduleService) GetCachedExpandedSchedules(ctx context.Context, centerID uint, req *ExpandRulesRequest) ([]ExpandedSchedule, error) {
	cacheKey := s.scheduleCacheKey(centerID, req.StartDate, req.EndDate)

	// 嘗試從快取取得
	var cachedResult []ExpandedSchedule
	err := s.cacheSvc.GetJSON(ctx, &cachedResult, CacheCategorySchedule, cacheKey)
	if err == nil {
		s.Logger.Debug("cache hit for expanded schedules", "center_id", centerID, "start", req.StartDate.Format("2006-01-02"), "end", req.EndDate.Format("2006-01-02"))
		return cachedResult, nil
	}

	// 快取未命中，執行展開
	s.Logger.Debug("cache miss, expanding schedules", "center_id", centerID, "start", req.StartDate.Format("2006-01-02"), "end", req.EndDate.Format("2006-01-02"))

	result, err := s.ExpandRules(ctx, centerID, req)
	if err != nil {
		return nil, err
	}

	// 非同步寫入快取（不影響主要流程）
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 快取 10 分鐘
		_ = s.cacheSvc.SetWithTTL(cacheCtx, 10*time.Minute, CacheCategorySchedule, result, cacheKey)
	}()

	return result, nil
}

// GetCachedTeacherSchedule 取得帶快取的老師課表
func (s *ScheduleService) GetCachedTeacherSchedule(ctx context.Context, teacherID, centerID uint, startDate, endDate time.Time) ([]ExpandedSchedule, error) {
	cacheKey := s.teacherScheduleCacheKey(teacherID, centerID, startDate, endDate)

	// 嘗試從快取取得
	var cachedResult []ExpandedSchedule
	err := s.cacheSvc.GetJSON(ctx, &cachedResult, CacheCategorySchedule, cacheKey)
	if err == nil {
		s.Logger.Debug("cache hit for teacher schedule", "teacher_id", teacherID, "center_id", centerID)
		return cachedResult, nil
	}

	// 快取未命中，執行展開
	s.Logger.Debug("cache miss, expanding teacher schedule", "teacher_id", teacherID, "center_id", centerID)

	// 取得該老師的規則
	rules, err := s.ruleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	var teacherRules []models.ScheduleRule
	for _, rule := range rules {
		if rule.TeacherID != nil && *rule.TeacherID == teacherID {
			teacherRules = append(teacherRules, rule)
		}
	}

	result := s.expansionSvc.ExpandRules(ctx, teacherRules, startDate, endDate, centerID)

	// 非同步寫入快取
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.cacheSvc.SetWithTTL(cacheCtx, 10*time.Minute, CacheCategorySchedule, result, cacheKey)
	}()

	return result, nil
}

// InvalidateScheduleCache 使課表快取失效
// 當 ScheduleRule 異動時呼叫，精準清除受影響的快取
func (s *ScheduleService) InvalidateScheduleCache(ctx context.Context, centerID uint, ruleIDs []uint, startDate, endDate time.Time) error {
	// 清除指定中心的課表快取
	cacheKey := s.scheduleCacheKey(centerID, startDate, endDate)
	if err := s.cacheSvc.Delete(ctx, CacheCategorySchedule, cacheKey); err != nil {
		s.Logger.Warn("failed to delete schedule cache", "error", err)
	}

	s.Logger.Info("schedule cache invalidated", "center_id", centerID, "rule_ids", ruleIDs)
	return nil
}

// InvalidateCenterScheduleCache 使中心所有課表快取失效
// 用於大規模異動（如刪除課程、調整時間表）
func (s *ScheduleService) InvalidateCenterScheduleCache(ctx context.Context, centerID uint) error {
	pattern := fmt.Sprintf("schedule:expand:center:%d:*", centerID)
	if err := s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, pattern); err != nil {
		s.Logger.Warn("failed to delete center schedule cache pattern", "error", err)
	}

	s.Logger.Info("center schedule cache invalidated", "center_id", centerID)
	return nil
}

// InvalidateTeacherScheduleCache 使老師課表快取失效
// 當老師代課或異動時呼叫
func (s *ScheduleService) InvalidateTeacherScheduleCache(ctx context.Context, teacherID, centerID uint) error {
	pattern := fmt.Sprintf("schedule:expand:teacher:%d:center:%d:*", teacherID, centerID)
	if err := s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, pattern); err != nil {
		s.Logger.Warn("failed to delete teacher schedule cache pattern", "error", err)
	}

	s.Logger.Info("teacher schedule cache invalidated", "teacher_id", teacherID, "center_id", centerID)
	return nil
}

// InvalidateExceptionRelatedCache 使例外相關快取失效
// 當 ScheduleException 異動時呼叫
func (s *ScheduleService) InvalidateExceptionRelatedCache(ctx context.Context, centerID uint, exception *models.ScheduleException) error {
	// 取得例外相關的規則
	rule, err := s.ruleRepo.GetByID(ctx, exception.RuleID)
	if err != nil {
		s.Logger.Warn("failed to get rule for cache invalidation", "error", err)
		return nil
	}

	// 清除中心課表快取
	_ = s.InvalidateCenterScheduleCache(ctx, centerID)

	// 清除老師課表快取
	if rule.ID != 0 && rule.TeacherID != nil {
		_ = s.InvalidateTeacherScheduleCache(ctx, *rule.TeacherID, centerID)
	}

	s.Logger.Info("exception-related cache invalidated", "center_id", centerID, "exception_id", exception.ID, "rule_id", exception.RuleID)
	return nil
}

// ScheduleResource 排課資源轉換
type ScheduleResource struct {
	app *app.App
}

func NewScheduleResource(app *app.App) *ScheduleResource {
	return &ScheduleResource{app: app}
}

// ToRuleResponse 轉換為規則響應
func (r *ScheduleResource) ToRuleResponse(rule models.ScheduleRule) *resources.ScheduleRuleResponse {
	return &resources.ScheduleRuleResponse{
		ID:            rule.ID,
		CenterID:      rule.CenterID,
		OfferingID:    rule.OfferingID,
		TeacherID:     rule.TeacherID,
		RoomID:        rule.RoomID,
		Weekday:       rule.Weekday,
		StartTime:     rule.StartTime,
		EndTime:       rule.EndTime,
		Duration:      rule.Duration,
		EffectiveFrom: rule.EffectiveRange.StartDate.Format("2006-01-02"),
		EffectiveTo:   rule.EffectiveRange.EndDate.Format("2006-01-02"),
	}
}

// =========================================
// Matrix View Service Implementation
// =========================================

// GetMatrixViewData 取得矩陣視圖資料
// 實現 BFF 模式：後端直接回傳前端可直接渲染的矩陣結構
func (s *ScheduleService) GetMatrixViewData(ctx context.Context, centerID uint, req *requests.MatrixViewRequest) (*resources.MatrixViewResponse, error) {
	// 解析日期
	loc := libs.GetTaiwanLocation()
	startDate, err := time.ParseInLocation("2006-01-02", req.StartDate, loc)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}
	endDate, err := time.ParseInLocation("2006-01-02", req.EndDate, loc)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %w", err)
	}

	// 取得中心的營業時間設定
	var centerSettings *models.CenterSettings
	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err == nil {
		centerSettings = &center.Settings
	}

	// 取得展開後的課表
	expandReq := &ExpandRulesRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}
	expandedSchedules, err := s.ExpandRules(ctx, centerID, expandReq)
	if err != nil {
		return nil, fmt.Errorf("failed to expand rules: %w", err)
	}

	// 取得資源列表
	teacherMap := make(map[uint]*models.Teacher)
	roomMap := make(map[uint]*models.Room)

	// 收集所有資源 ID
	teacherIDs := make([]uint, 0)
	roomIDs := make([]uint, 0)
	for _, schedule := range expandedSchedules {
		if schedule.TeacherID != nil && *schedule.TeacherID > 0 {
			if _, exists := teacherMap[*schedule.TeacherID]; !exists {
				teacherIDs = append(teacherIDs, *schedule.TeacherID)
			}
		}
		if schedule.RoomID > 0 {
			if _, exists := roomMap[schedule.RoomID]; !exists {
				roomIDs = append(roomIDs, schedule.RoomID)
			}
		}
	}

	// 取得老師資料
	if len(teacherIDs) > 0 {
		teachers, _ := s.teacherRepo.BatchGetByIDs(ctx, teacherIDs)
		for _, t := range teachers {
			teacherMap[t.ID] = &t
		}
	}

	// 取得教室資料
	if len(roomIDs) > 0 {
		// 取得該中心的所有教室，然後過濾
		allRooms, _ := s.roomRepo.ListByCenterID(ctx, centerID)
		for _, r := range allRooms {
			// 只保留在 roomIDs 中的教室
			for _, id := range roomIDs {
				if r.ID == id {
					roomMap[r.ID] = &r
					break
				}
			}
		}
	}

	// 解析 include_suspended 參數
	includeSuspended := true
	if req.IncludeSuspended != nil {
		includeSuspended = *req.IncludeSuspended
	}

	// 解析 resource_ids 參數
	filterResourceIDs := make(map[uint]bool)
	if len(req.ResourceIDs) > 0 {
		for _, id := range req.ResourceIDs {
			filterResourceIDs[id] = true
		}
	}

	// 根據類型分組資源
	var matrixResources []resources.MatrixResource
	var usedResourceIDs map[uint]bool

	if req.Type == "teacher" || req.Type == "all" {
		matrixResources, usedResourceIDs = s.groupByTeacher(expandedSchedules, teacherMap, roomMap, filterResourceIDs, includeSuspended)
	}
	if req.Type == "room" || req.Type == "all" {
		roomResources, _ := s.groupByRoom(expandedSchedules, teacherMap, roomMap, filterResourceIDs, includeSuspended)
		// 如果是 all 類型，過濾掉已重複的資源
		if req.Type == "all" && usedResourceIDs != nil {
			var filteredRoomResources []resources.MatrixResource
			for _, r := range roomResources {
				if req.Type == "all" {
					// 檢查是否已在 teacher 中出現
					isDuplicate := false
					for _, mr := range matrixResources {
						if mr.ID == r.ID && mr.Type == r.Type {
							isDuplicate = true
							break
						}
					}
					if !isDuplicate {
						filteredRoomResources = append(filteredRoomResources, r)
					}
				}
			}
			matrixResources = append(matrixResources, filteredRoomResources...)
		} else {
			matrixResources = append(matrixResources, roomResources...)
		}
	}

	// 生成 time_slots（使用動態營業時間）
	timeSlots := s.generateTimeSlots(expandedSchedules, centerSettings)

	// 構建響應
	response := &resources.MatrixViewResponse{
		TimeSlots: timeSlots,
		Resources: matrixResources,
		DateRange: resources.MatrixDateRange{
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		},
	}

	return response, nil
}

// groupByTeacher 將課表按老師分組
func (s *ScheduleService) groupByTeacher(
	expandedSchedules []ExpandedSchedule,
	teacherMap map[uint]*models.Teacher,
	roomMap map[uint]*models.Room,
	filterResourceIDs map[uint]bool,
	includeSuspended bool,
) ([]resources.MatrixResource, map[uint]bool) {
	usedResourceIDs := make(map[uint]bool)
	teacherSchedules := make(map[uint][]ExpandedSchedule)

	for _, schedule := range expandedSchedules {
		if schedule.TeacherID == nil || *schedule.TeacherID == 0 {
			continue
		}

		// 過濾資源
		if len(filterResourceIDs) > 0 && !filterResourceIDs[*schedule.TeacherID] {
			continue
		}

		// 過濾停課
		if !includeSuspended && schedule.IsHoliday {
			continue
		}

		teacherSchedules[*schedule.TeacherID] = append(teacherSchedules[*schedule.TeacherID], schedule)
		usedResourceIDs[*schedule.TeacherID] = true
	}

	var matrixResources []resources.MatrixResource
	for teacherID, schedules := range teacherSchedules {
		teacher := teacherMap[teacherID]
		name := "未安排老師"
		if teacher != nil && teacher.Name != "" {
			name = teacher.Name
		}

		items := s.buildMatrixItems(schedules, teacherMap, roomMap)

		matrixResources = append(matrixResources, resources.MatrixResource{
			ID:    teacherID,
			Name:  name,
			Type:  "teacher",
			Items: items,
		})
	}

	return matrixResources, usedResourceIDs
}

// groupByRoom 將課表按教室分組
func (s *ScheduleService) groupByRoom(
	expandedSchedules []ExpandedSchedule,
	teacherMap map[uint]*models.Teacher,
	roomMap map[uint]*models.Room,
	filterResourceIDs map[uint]bool,
	includeSuspended bool,
) ([]resources.MatrixResource, map[uint]bool) {
	usedResourceIDs := make(map[uint]bool)
	roomSchedules := make(map[uint][]ExpandedSchedule)

	for _, schedule := range expandedSchedules {
		if schedule.RoomID == 0 {
			continue
		}

		// 過濾資源
		if len(filterResourceIDs) > 0 && !filterResourceIDs[schedule.RoomID] {
			continue
		}

		// 過濾停課
		if !includeSuspended && schedule.IsHoliday {
			continue
		}

		roomSchedules[schedule.RoomID] = append(roomSchedules[schedule.RoomID], schedule)
		usedResourceIDs[schedule.RoomID] = true
	}

	var matrixResources []resources.MatrixResource
	for roomID, schedules := range roomSchedules {
		room := roomMap[roomID]
		name := "未安排教室"
		if room != nil && room.Name != "" {
			name = room.Name
		}

		items := s.buildMatrixItems(schedules, teacherMap, roomMap)

		matrixResources = append(matrixResources, resources.MatrixResource{
			ID:    roomID,
			Name:  name,
			Type:  "room",
			Items: items,
		})
	}

	return matrixResources, usedResourceIDs
}

// buildMatrixItems 構建矩陣項目列表
func (s *ScheduleService) buildMatrixItems(
	schedules []ExpandedSchedule,
	teacherMap map[uint]*models.Teacher,
	roomMap map[uint]*models.Room,
) []resources.MatrixItem {
	items := make([]resources.MatrixItem, 0, len(schedules))

	for _, schedule := range schedules {
		// 解析時間
		startHour, startMinute := s.parseTimeToHourMinute(schedule.StartTime)
		_, _ = s.parseTimeToHourMinute(schedule.EndTime)

		// 計算持續分鐘數
		duration := s.calculateDuration(schedule.StartTime, schedule.EndTime)

		// 計算 CSS 位置和高度
		topOffset := float64(startHour*60+startMinute) / 1440 * 100 // 一天 1440 分鐘
		heightPercent := float64(duration) / 1440 * 100

		// 取得老師名稱
		teacherName := "未安排老師"
		if schedule.TeacherID != nil && *schedule.TeacherID > 0 {
			if teacher, exists := teacherMap[*schedule.TeacherID]; exists {
				teacherName = teacher.Name
			}
		}

		// 取得教室名稱
		roomName := "未安排教室"
		if schedule.RoomID > 0 {
			if room, exists := roomMap[schedule.RoomID]; exists {
				roomName = room.Name
			}
		}

		// 判斷是否為停課
		isSuspended := schedule.IsHoliday

		item := resources.MatrixItem{
			ID:            schedule.RuleID,
			RuleID:        schedule.RuleID,
			Title:         schedule.OfferingName,
			Date:          schedule.Date.Format("2006-01-02"),
			StartTime:     schedule.StartTime,
			EndTime:       schedule.EndTime,
			StartHour:     startHour,
			StartMinute:   startMinute,
			Duration:      duration,
			TopOffset:     topOffset,
			HeightPercent: heightPercent,
			OfferingID:    schedule.OfferingID,
			OfferingName:  schedule.OfferingName,
			TeacherID:     schedule.TeacherID,
			TeacherName:   teacherName,
			RoomID:        schedule.RoomID,
			RoomName:      roomName,
			IsHoliday:     schedule.IsHoliday,
			HasException:  schedule.HasException,
			IsSuspended:   isSuspended,
		}

		items = append(items, item)
	}

	return items
}

// parseTimeToHourMinute 解析時間字串為小時和分鐘
func (s *ScheduleService) parseTimeToHourMinute(timeStr string) (int, int) {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return 0, 0
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	return hour, minute
}

// calculateDuration 計算持續分鐘數
func (s *ScheduleService) calculateDuration(startTime, endTime string) int {
	startParts := strings.Split(startTime, ":")
	endParts := strings.Split(endTime, ":")

	if len(startParts) < 2 || len(endParts) < 2 {
		return 60 // 預設 60 分鐘
	}

	startHour, _ := strconv.Atoi(startParts[0])
	startMinute, _ := strconv.Atoi(startParts[1])
	endHour, _ := strconv.Atoi(endParts[0])
	endMinute, _ := strconv.Atoi(endParts[1])

	startTotal := startHour*60 + startMinute
	endTotal := endHour*60 + endMinute

	// 跨日處理
	if endTotal <= startTotal {
		endTotal += 24 * 60
	}

	return endTotal - startTotal
}

// generateTimeSlots 生成時段列表
// 根據中心的營業時間生成動態時段列表
// 如果 settings 為空或沒有設定，則使用預設值 (00:00 - 23:00)
func (s *ScheduleService) generateTimeSlots(schedules []ExpandedSchedule, settings *models.CenterSettings) []int {
	// 解析營業時間
	var startHour, endHour int

	if settings != nil && settings.OperatingStartTime != "" {
		startParts := strings.Split(settings.OperatingStartTime, ":")
		if len(startParts) >= 1 {
			fmt.Sscanf(startParts[0], "%d", &startHour)
		}
	} else {
		startHour = 0 // 預設開始時間
	}

	if settings != nil && settings.OperatingEndTime != "" {
		endParts := strings.Split(settings.OperatingEndTime, ":")
		if len(endParts) >= 1 {
			fmt.Sscanf(endParts[0], "%d", &endHour)
		}
	} else {
		endHour = 23 // 預設結束時間
	}

	// 確保結束時間大於開始時間
	if endHour <= startHour {
		endHour = 23
	}

	// 生成時段列表
	slots := make([]int, 0, endHour-startHour+1)
	for i := startHour; i <= endHour; i++ {
		slots = append(slots, i)
	}

	return slots
}
