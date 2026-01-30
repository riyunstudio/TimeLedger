package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"

	"gorm.io/gorm"
)

type ScheduleExpansionServiceImpl struct {
	BaseService
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	exceptionRepo    *repositories.ScheduleExceptionRepository
	auditLogRepo     *repositories.AuditLogRepository
	centerRepo       *repositories.CenterRepository
	holidayRepo      *repositories.CenterHolidayRepository
}

func NewScheduleExpansionService(app *app.App) ScheduleExpansionService {
	baseSvc := NewBaseService(app, "ScheduleExpansionService")
	return &ScheduleExpansionServiceImpl{
		BaseService:      *baseSvc,
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

	// 批次取得所有規則在日期範圍內的例外資料（消除 N+1 查詢）
	ruleIDs := make([]uint, 0, len(rules))
	for _, rule := range rules {
		if rule.Weekday != 0 {
			ruleIDs = append(ruleIDs, rule.ID)
		}
	}

	// 建立例外資料 Map：RuleID -> DateString -> Exceptions
	exceptionsMap := make(map[uint]map[string][]models.ScheduleException)
	if len(ruleIDs) > 0 {
		exceptionsMap, _ = s.exceptionRepo.GetByRuleIDsAndDateRange(ctx, ruleIDs, startDate, endDate)
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

					// 如果是假日，跳過此 session（無感停課）
					if isHoliday {
						date = date.AddDate(0, 0, 1)
						continue
					}

					// 從預載入的 Map 中取得例外資料（批次查詢優化）
					ruleExceptions := exceptionsMap[rule.ID]
					exceptions := []models.ScheduleException{}
					if ruleExceptions != nil {
						exceptions = ruleExceptions[dateStr]
					}

					// 過濾並處理例外狀態
					skipSession := false
					var pendingException *models.ScheduleException
					var approvedException *models.ScheduleException

					for i := range exceptions {
						exc := &exceptions[i]
						// 跳過已取消的例外
						if exc.Status == "CANCELLED" {
							continue
						}

						// 記錄待處理的例外（用於顯示）
						if exc.Status == "PENDING" {
							if pendingException == nil {
								pendingException = exc
							}
						}

						// 處理已核准的例外
						if exc.Status == "APPROVED" {
							// 如果是停課例外，跳過此 sessions
							if exc.ExceptionType == "CANCEL" {
								skipSession = true
								approvedException = exc
								break
							}
							// 代課例外：更新老師
							if exc.ExceptionType == "REPLACE_TEACHER" && exc.NewTeacherID != nil {
								rule.TeacherID = exc.NewTeacherID
							}
							// 調課例外：由新規則處理（不會在這裡產生 sessions）
							if exc.ExceptionType == "RESCHEDULE" {
								// RESCHEDULE 由新規則處理，這裡標記但繼續產生原時段（用於顯示）
								if approvedException == nil {
									approvedException = exc
								}
							}
						}
					}

					if skipSession {
						date = date.AddDate(0, 0, 1)
						continue
					}

					// 解析開始和結束時間
					startParts := strings.Split(rule.StartTime, ":")
					endParts := strings.Split(rule.EndTime, ":")
					startHour, _ := strconv.Atoi(startParts[0])
					endHour, _ := strconv.Atoi(endParts[0])

					// 處理跨日課程（結束時間早於開始時間）
					isCrossDay := endHour < startHour

					// 創建課表項目
					createScheduleEntry := func(scheduleDate time.Time, sTime, eTime string) {
						schedule := ExpandedSchedule{
							RuleID:       rule.ID,
							Date:         scheduleDate,
							StartTime:    sTime,
							EndTime:      eTime,
							RoomID:       rule.RoomID,
							TeacherID:    rule.TeacherID,
							IsHoliday:    isHoliday,
							HasException: pendingException != nil || approvedException != nil,
							// 關聯資料
							OfferingName:   rule.Offering.Name,
							TeacherName:    rule.Teacher.Name,
							RoomName:       rule.Room.Name,
							OfferingID:     rule.OfferingID,
							EffectiveRange: &rule.EffectiveRange,
							IsCrossDayPart: isCrossDay, // 標記是否為跨日課程的一部分
						}

						// 添加例外資訊供前端顯示
						if pendingException != nil {
							schedule.ExceptionInfo = &ExpandedException{
								ID:           pendingException.ID,
								Type:         pendingException.ExceptionType,
								Status:       pendingException.Status,
								NewTeacherID: pendingException.NewTeacherID,
								NewStartAt:   pendingException.NewStartAt,
								NewEndAt:     pendingException.NewEndAt,
							}
						}

						schedules = append(schedules, schedule)
					}

					if isCrossDay {
						// 跨日課程：生成兩個條目
						// 條目1：開始日 22:00-24:00
						createScheduleEntry(date, rule.StartTime, "24:00")
						// 條目2：結束日（隔天）00:00-01:00
						nextDay := date.AddDate(0, 0, 1)
						createScheduleEntry(nextDay, "00:00", rule.EndTime)
					} else {
						// 普通課程：生成一個條目
						createScheduleEntry(date, rule.StartTime, rule.EndTime)
					}
				}
			}

			date = date.AddDate(0, 0, 1)
		}
	}

	return schedules
}

func (s *ScheduleExpansionServiceImpl) GetEffectiveRuleForDate(ctx context.Context, offeringID uint, date time.Time) (*models.ScheduleRule, error) {
	var rules []models.ScheduleRule
	err := s.App.MySQL.RDB.WithContext(ctx).
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
	// 批次查詢所有規則（消除 N+1：避免在迴圈中每次查詢資料庫）
	rules, err := s.scheduleRuleRepo.ListByOfferingID(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	if len(rules) <= 1 {
		return []PhaseTransition{}, nil
	}

	// 建立日期到規則的 Map：DateString -> *ScheduleRule
	// 這將 N 次查詢優化為 1 次查詢 + O(n) Map 建構
	ruleByDate := make(map[string]*models.ScheduleRule)
	for i := range rules {
		rule := &rules[i]
		ruleStart := rule.EffectiveRange.StartDate
		ruleEnd := rule.EffectiveRange.EndDate

		// 展開規則的有效日期範圍
		current := ruleStart
		for !current.After(ruleEnd) {
			// 檢查是否在查詢範圍內
			if !current.Before(startDate) && !current.After(endDate) {
				weekday := int(current.Weekday())
				if weekday == 0 {
					weekday = 7
				}

				// 檢查星期是否匹配
				if weekday == int(rule.Weekday) {
					dateStr := current.Format("2006-01-02")
					ruleByDate[dateStr] = rule
				}
			}
			current = current.AddDate(0, 0, 1)
		}
	}

	var transitions []PhaseTransition
	date := startDate
	prevRule := (*models.ScheduleRule)(nil)

	// 使用 Map 查詢替代資料庫查詢（O(1) vs O(n)）
	for date.Before(endDate) || date.Equal(endDate) {
		dateStr := date.Format("2006-01-02")
		currentRule := ruleByDate[dateStr]

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
	err := s.App.MySQL.RDB.WithContext(ctx).
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
		BaseService:       *baseSvc,
		exceptionRepo:     repositories.NewScheduleExceptionRepository(app),
		ruleRepo:          repositories.NewScheduleRuleRepository(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
		centerRepo:        repositories.NewCenterRepository(app),
		teacherRepo:       repositories.NewTeacherRepository(app),
		validationService: NewScheduleValidationService(app),
		notificationSvc:   NewNotificationService(app),
		notificationQueue: NewNotificationQueueService(app),
		cacheSvc:          NewCacheService(app),
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

func (s *ScheduleExceptionServiceImpl) CreateException(ctx context.Context, centerID uint, teacherID uint, ruleID uint, originalDate time.Time, exceptionType string, newStartAt, newEndAt *time.Time, newTeacherID *uint, newTeacherName string, reason string) (models.ScheduleException, error) {
	allowed, reasonStr, err := s.CheckExceptionDeadline(ctx, centerID, ruleID, originalDate)
	if err != nil {
		return models.ScheduleException{}, err
	}
	if !allowed {
		return models.ScheduleException{}, errors.New(reasonStr)
	}

	exception := models.ScheduleException{
		CenterID:      centerID,
		RuleID:        ruleID,
		OriginalDate:  originalDate,
		ExceptionType: exceptionType,
		Status:        "PENDING",
		NewStartAt:    newStartAt,
		NewEndAt:      newEndAt,
		NewTeacherID:  newTeacherID,
		Reason:        reason,
	}

	// 如果有代課老師名字，添加到 Reason 中（因為資料庫沒有這個欄位）
	if newTeacherName != "" {
		exception.Reason = fmt.Sprintf("[代課老師：%s] %s", newTeacherName, reason)
	}

	// 使用交易確保建立例外單和審核日誌的原子性
	var createdException models.ScheduleException
	var teacherName, centerName string

	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在交易中建立例外單
		var createErr error
		createdException, createErr = s.exceptionRepo.CreateWithDB(ctx, tx, exception)
		if createErr != nil {
			return fmt.Errorf("failed to create exception: %w", createErr)
		}

		// 在交易中記錄審核日誌
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
		return models.ScheduleException{}, txErr
	}

	// 取得老師和中心名稱用於 LINE 通知（交易外執行）
	teacher, _ := s.teacherRepo.GetByID(ctx, teacherID)
	teacherName = teacher.Name
	if teacherName == "" {
		teacherName = "老師"
	}

	center, _ := s.centerRepo.GetByID(ctx, centerID)
	centerName = center.Name

	// 更新 exception 的 ExceptionType 欄位（用於 LINE 通知）
	createdException.ExceptionType = exceptionType

	// 同步發送通知（直接發送，不經佇列）
	if s.notificationQueue != nil {
		_ = s.notificationQueue.NotifyExceptionSubmittedSync(ctx, &createdException, teacherName, centerName)
	}

	// 建立例外單後，使相關課表快取失效
	s.invalidateRelatedCaches(ctx, &createdException)

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

	// 撤回例外單後，使相關課表快取失效
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

	// 映射前端發送的 action 值到正確的狀態值
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

	// 使用交易確保審核操作的原子性
	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 如果核准，執行驗證和課程變更
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

			// 在交易中執行課程變更
			if err := s.applyExceptionChangesWithTx(ctx, tx, &exception, &rule); err != nil {
				return fmt.Errorf("failed to apply exception changes: %w", err)
			}
		}

		// 在交易中更新例外狀態
		if err := tx.Save(&exception).Error; err != nil {
			return fmt.Errorf("failed to update exception: %w", err)
		}

		// 在交易中記錄審核日誌
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

	// 發送 LINE 通知給老師（同步發送，交易外執行）
	// 取得老師資料
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

	// 審核例外單後，使相關課表快取失效
	s.invalidateRelatedCaches(ctx, &exception)

	return nil
}

// applyExceptionChanges 審核通過時同步執行課程變更
func (s *ScheduleExceptionServiceImpl) applyExceptionChanges(ctx context.Context, exception *models.ScheduleException, rule *models.ScheduleRule) error {
	switch exception.ExceptionType {
	case "CANCEL":
		// 停課：不修改規則本身，由 ExpandRules 根據已核准的例外狀態跳過該日期
		// 這樣只會影響特定的停課日期，不會影響未來的課程
		// 規則保持不變，ExpandRules 會在處理時檢查例外狀態

	case "RESCHEDULE":
		// 調課：分兩步驟
		// 1. 截斷原規則到前一天（確保原日期不再產生 sessions）
		// 2. 創建新規則段，從新日期開始
		cutoffDate := exception.OriginalDate.AddDate(0, 0, -1)
		rule.EffectiveRange.EndDate = cutoffDate
		if err := s.ruleRepo.Update(ctx, *rule); err != nil {
			return fmt.Errorf("failed to truncate original rule: %w", err)
		}

		// 轉換 weekday：Go 的 Weekday() 返回週日=0，但系統使用週日=7
		newWeekday := int(exception.NewStartAt.Weekday())
		if newWeekday == 0 {
			newWeekday = 7
		}

		// 創建新規則段
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
		// 代課：更新原規則的老師為代課老師
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

// applyExceptionChangesWithTx 審核通過時同步執行課程變更（交易版本）
func (s *ScheduleExceptionServiceImpl) applyExceptionChangesWithTx(ctx context.Context, tx *gorm.DB, exception *models.ScheduleException, rule *models.ScheduleRule) error {
	switch exception.ExceptionType {
	case "CANCEL":
		// 停課：不修改規則本身，由 ExpandRules 根據已核准的例外狀態跳過該日期
		// 這樣只會影響特定的停課日期，不會影響未來的課程
		// 規則保持不變，ExpandRules 會在處理時檢查例外狀態

	case "RESCHEDULE":
		// 調課：分兩步驟
		// 1. 截斷原規則到前一天（確保原日期不再產生 sessions）
		// 2. 創建新規則段，從新日期開始
		cutoffDate := exception.OriginalDate.AddDate(0, 0, -1)
		rule.EffectiveRange.EndDate = cutoffDate
		if err := tx.Save(rule).Error; err != nil {
			return fmt.Errorf("failed to truncate original rule: %w", err)
		}

		// 轉換 weekday：Go 的 Weekday() 返回週日=0，但系統使用週日=7
		newWeekday := int(exception.NewStartAt.Weekday())
		if newWeekday == 0 {
			newWeekday = 7
		}

		// 創建新規則段
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
		// 代課：更新原規則的老師為代課老師
		if exception.NewTeacherID != nil {
			rule.TeacherID = exception.NewTeacherID
			rule.UpdatedAt = time.Now()
			if err := tx.Save(rule).Error; err != nil {
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

// GetAllExceptions 取得所有例外單，可依狀態篩選
func (s *ScheduleExceptionServiceImpl) GetAllExceptions(ctx context.Context, centerID uint, status string) ([]models.ScheduleException, error) {
	var exceptions []models.ScheduleException
	query := s.App.MySQL.RDB.WithContext(ctx).
		Preload("Rule").
		Preload("Rule.Teacher").
		Preload("Rule.Room").
		Where("center_id = ?", centerID)

	// 如果有指定狀態，則過濾
	if status != "" {
		// 支援舊資料的 APPROVE/REJECT 狀態
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

// invalidateRelatedCaches 使例外相關的課表快取失效
func (s *ScheduleExceptionServiceImpl) invalidateRelatedCaches(ctx context.Context, exception *models.ScheduleException) {
	// 取得例外相關的規則
	rule, err := s.ruleRepo.GetByID(ctx, exception.RuleID)
	if err != nil {
		s.Logger.Warn("failed to get rule for cache invalidation", "error", err, "exception_id", exception.ID)
		return
	}

	// 清除中心課表快取
	pattern := fmt.Sprintf("schedule:expand:center:%d:*", exception.CenterID)
	_ = s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, pattern)

	// 清除老師課表快取
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
		app:           app,
		ruleRepo:      repositories.NewScheduleRuleRepository(app),
		exceptionRepo: repositories.NewScheduleExceptionRepository(app),
		expansionSvc:  NewScheduleExpansionService(app),
		auditLogRepo:  repositories.NewAuditLogRepository(app),
		offeringRepo:  repositories.NewOfferingRepository(app),
		cacheSvc:      NewCacheService(app),
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

	// 編輯循環課表後，使相關課表快取失效
	s.invalidateRelatedCaches(ctx, centerID, req.RuleID)

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

	// 使用交易確保建立例外單和審核日誌的原子性
	var createdCancel, createdAdd models.ScheduleException

	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在交易中建立取消例外單
		var createErr error
		createdCancel, createErr = s.exceptionRepo.CreateWithDB(ctx, tx, cancelExc)
		if createErr != nil {
			return fmt.Errorf("failed to create cancel exception: %w", createErr)
		}

		// 在交易中建立新增例外單
		createdAdd, createErr = s.exceptionRepo.CreateWithDB(ctx, tx, addExc)
		if createErr != nil {
			return fmt.Errorf("failed to create add exception: %w", createErr)
		}

		// 在交易中記錄審核日誌
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

	// 使用交易確保所有操作的原子性
	var createdRule *models.ScheduleRule

	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在交易中建立所有取消和新增例外單
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

			newTeacherID := req.NewTeacherID
			if newTeacherID == nil {
				newTeacherID = rule.TeacherID
			}
			newRoomID := req.NewRoomID
			if newRoomID == nil {
				newRoomIDVal := rule.RoomID
				newRoomID = &newRoomIDVal
			}

			addExc := models.ScheduleException{
				CenterID:      centerID,
				RuleID:        req.RuleID,
				OriginalDate:  date,
				ExceptionType: "ADD",
				Status:        "PENDING",
				NewStartAt:    &newStartAt,
				NewEndAt:      &newEndAt,
				NewTeacherID:  newTeacherID,
				NewRoomID:     newRoomID,
				Reason:        req.Reason,
			}
			createdAdd, createErr := s.exceptionRepo.CreateWithDB(ctx, tx, addExc)
			if createErr != nil {
				return fmt.Errorf("failed to create add exception: %w", createErr)
			}
			addExcs = append(addExcs, createdAdd)
		}

		// 在交易中建立新規則
		newRule := models.ScheduleRule{
			CenterID:       centerID,
			OfferingID:     rule.OfferingID,
			TeacherID:      req.NewTeacherID,
			RoomID:         *req.NewRoomID,
			Weekday:        rule.Weekday,
			StartTime:      req.NewStartTime,
			EndTime:        req.NewEndTime,
			EffectiveRange: models.DateRange{StartDate: req.EditDate, EndDate: rule.EffectiveRange.EndDate},
		}

		if req.NewTeacherID == nil {
			newRule.TeacherID = rule.TeacherID
		}
		if req.NewStartTime == "" {
			newRule.StartTime = rule.StartTime
		}
		if req.NewEndTime == "" {
			newRule.EndTime = rule.EndTime
		}
		if req.NewRoomID == nil {
			newRule.RoomID = rule.RoomID
		}

		if err := tx.Create(&newRule).Error; err != nil {
			return fmt.Errorf("failed to create new rule: %w", err)
		}
		createdRule = &newRule

		// 在交易中記錄審核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "EDIT_FUTURE_OCCURRENCES",
			TargetType: "ScheduleRule",
			TargetID:   newRule.ID,
			Payload:    models.AuditPayload{After: fmt.Sprintf("Create new rule %d for future occurrences from %s", newRule.ID, req.EditDate.Format("2006-01-02"))},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, cancelExcs, addExcs, txErr
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

	err := s.ruleRepo.Update(ctx, rule)
	if err != nil {
		return nil, err
	}

	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "EDIT_ALL_OCCURRENCES",
		TargetType: "ScheduleRule",
		TargetID:   rule.ID,
		Payload:    models.AuditPayload{After: fmt.Sprintf("Update all occurrences for rule %d", rule.ID)},
	})

	return &rule, nil
}

func (s *ScheduleRecurrenceServiceImpl) DeleteRecurringSchedule(ctx context.Context, centerID uint, teacherID uint, ruleID uint, editDate time.Time, mode RecurrenceEditMode, reason string) (RecurrenceEditResult, error) {
	result := RecurrenceEditResult{
		Mode:             mode,
		CancelExceptions: []models.ScheduleException{},
		AffectedCount:    0,
	}

	var affectedCount int
	var preview RecurrenceEditPreview
	var ruleDeletionID uint

	// 使用交易確保所有操作的原子性
	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
			created, createErr := s.exceptionRepo.CreateWithDB(ctx, tx, cancelExc)
			if createErr != nil {
				return fmt.Errorf("failed to create cancel exception: %w", createErr)
			}
			result.CancelExceptions = append(result.CancelExceptions, created)
			result.AffectedCount = 1
			affectedCount = 1

		case RecurrenceEditFuture:
			preview, _ = s.PreviewAffectedSessions(ctx, ruleID, editDate, RecurrenceEditFuture)
			for _, date := range preview.AffectedDates {
				cancelExc := models.ScheduleException{
					CenterID:      centerID,
					RuleID:        ruleID,
					OriginalDate:  date,
					ExceptionType: "CANCEL",
					Status:        "PENDING",
					Reason:        reason,
				}
				created, createErr := s.exceptionRepo.CreateWithDB(ctx, tx, cancelExc)
				if createErr != nil {
					return fmt.Errorf("failed to create cancel exception: %w", createErr)
				}
				result.CancelExceptions = append(result.CancelExceptions, created)
			}
			result.AffectedCount = preview.AffectedCount
			affectedCount = preview.AffectedCount

		case RecurrenceEditAll:
			preview, _ = s.PreviewAffectedSessions(ctx, ruleID, editDate, RecurrenceEditAll)
			// 獲取規則ID用於審核日誌
			rule, _ := s.ruleRepo.GetByID(ctx, ruleID)
			ruleDeletionID = rule.ID

			if err := tx.Unscoped().Where("id = ?", ruleID).Delete(&models.ScheduleRule{}).Error; err != nil {
				return fmt.Errorf("failed to delete rule: %w", err)
			}
			result.AffectedCount = preview.AffectedCount
			affectedCount = preview.AffectedCount
		}

		// 在交易中記錄審核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "DELETE_RECURRING_SCHEDULE",
			TargetType: "ScheduleRule",
			TargetID:   ruleDeletionID,
			Payload:    models.AuditPayload{After: fmt.Sprintf("Delete recurring schedule mode=%s, affected=%d", mode, affectedCount)},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return RecurrenceEditResult{}, txErr
	}

	// 刪除循環課表後，使相關課表快取失效
	s.invalidateRelatedCaches(ctx, centerID, ruleID)

	return result, nil
}

// invalidateRelatedCaches 使循環編輯相關的課表快取失效
func (s *ScheduleRecurrenceServiceImpl) invalidateRelatedCaches(ctx context.Context, centerID, ruleID uint) {
	// 取得規則以獲取老師ID
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		s.Logger.Warn("failed to get rule for cache invalidation", "error", err, "rule_id", ruleID)
		return
	}

	// 清除中心課表快取
	pattern := fmt.Sprintf("schedule:expand:center:%d:*", centerID)
	_ = s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, pattern)

	// 清除老師課表快取
	if rule.ID != 0 && rule.TeacherID != nil {
		teacherPattern := fmt.Sprintf("schedule:expand:teacher:%d:center:%d:*", *rule.TeacherID, centerID)
		_ = s.cacheSvc.DeleteByPattern(ctx, CacheCategorySchedule, teacherPattern)
	}

	s.Logger.Info("recurrence-edit cache invalidated",
		"center_id", centerID,
		"rule_id", ruleID,
		"teacher_id", func() uint {
			if rule.ID != 0 && rule.TeacherID != nil {
				return *rule.TeacherID
			}
			return 0
		}())
}
