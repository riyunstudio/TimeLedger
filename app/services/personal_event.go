package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global"
	"timeLedger/global/errInfos"
	"timeLedger/libs"
)

// PersonalEventService 教師個人行程相關業務邏輯
type PersonalEventService struct {
	BaseService
	app               *app.App
	personalEventRepo *repositories.PersonalEventRepository
	membershipRepo    *repositories.CenterMembershipRepository
	scheduleRuleRepo  *repositories.ScheduleRuleRepository
	auditLogRepo      *repositories.AuditLogRepository
}

// NewPersonalEventService 建立教師行程服務
func NewPersonalEventService(app *app.App) *PersonalEventService {
	return &PersonalEventService{
		app:               app,
		personalEventRepo: repositories.NewPersonalEventRepository(app),
		membershipRepo:    repositories.NewCenterMembershipRepository(app),
		scheduleRuleRepo:  repositories.NewScheduleRuleRepository(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
	}
}

// GetPersonalEvents 取得老師個人行程列表
func (s *PersonalEventService) GetPersonalEvents(ctx context.Context, teacherID uint, from, to time.Time) ([]models.PersonalEvent, *errInfos.Res, error) {
	var events []models.PersonalEvent
	var err error

	// 檢查 from 和 to 是否為零值
	fromIsZero := from.IsZero()
	toIsZero := to.IsZero()

	if !fromIsZero && !toIsZero {
		// Add one day to to date to make it inclusive
		toInclusive := to.AddDate(0, 0, 1)
		events, err = s.personalEventRepo.GetByTeacherAndDateRange(ctx, teacherID, from, toInclusive)
	} else {
		events, err = s.personalEventRepo.ListByTeacherID(ctx, teacherID)
	}

	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return events, nil, nil
}

// CreatePersonalEventRequest 新增個人行程請求
type CreatePersonalEventRequest struct {
	Title          string
	StartAt        time.Time
	EndAt          time.Time
	IsAllDay       bool
	ColorHex       string
	RecurrenceRule *models.RecurrenceRule
}

// CreatePersonalEventResult 新增個人行程結果
type CreatePersonalEventResult struct {
	Event       *models.PersonalEvent
	ConflictMsg string // 若有衝突，回傳衝突訊息
}

// CreatePersonalEvent 新增老師個人行程
func (s *PersonalEventService) CreatePersonalEvent(ctx context.Context, teacherID uint, req *CreatePersonalEventRequest) (*models.PersonalEvent, *errInfos.Res, error) {
	// 確保時間使用台灣時區
	startAt := libs.TimeToTaiwan(req.StartAt)
	endAt := libs.TimeToTaiwan(req.EndAt)

	// 檢查個人行程是否與中心課程衝突
	conflictMsg, err := s.checkPersonalEventConflicts(ctx, teacherID, startAt, endAt)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	if conflictMsg != "" {
		// 使用 INVALID_STATUS 錯誤碼，會映射為 409 Conflict
		return nil, &errInfos.Res{Code: errInfos.INVALID_STATUS, Msg: conflictMsg}, nil
	}

	event := &models.PersonalEvent{
		TeacherID: teacherID,
		Title:     req.Title,
		StartAt:   startAt,
		EndAt:     endAt,
		IsAllDay:  req.IsAllDay,
		ColorHex:  req.ColorHex,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.RecurrenceRule != nil {
		event.RecurrenceRule = *req.RecurrenceRule
	}

	createdEvent, err := s.personalEventRepo.Create(ctx, *event)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return &createdEvent, nil, nil
}

// UpdatePersonalEventRequest 更新個人行程請求
type UpdatePersonalEventRequest struct {
	Title          *string
	StartAt        *time.Time
	EndAt          *time.Time
	IsAllDay       *bool
	ColorHex       *string
	RecurrenceRule *models.RecurrenceRule
	UpdateMode     string // SINGLE, FUTURE, ALL
}

// UpdatePersonalEventResult 更新個人行程結果
type UpdatePersonalEventResult struct {
	UpdatedCount int64
	Message      string
}

// UpdatePersonalEvent 更新老師個人行程
func (s *PersonalEventService) UpdatePersonalEvent(ctx context.Context, eventID, teacherID uint, req *UpdatePersonalEventRequest) (*UpdatePersonalEventResult, *errInfos.Res, error) {
	// 取得行程並驗證權限
	event, err := s.personalEventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	if event.TeacherID != teacherID {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 確保時間使用台灣時區
	// 檢查衝突（如果時間有變更）
	if req.StartAt != nil && req.EndAt != nil {
		startAt := libs.TimeToTaiwan(*req.StartAt)
		endAt := libs.TimeToTaiwan(*req.EndAt)
		conflictMsg, err := s.checkPersonalEventConflicts(ctx, teacherID, startAt, endAt)
		if err != nil {
			return nil, s.app.Err.New(errInfos.SQL_ERROR), err
		}
		if conflictMsg != "" {
			// 使用 INVALID_STATUS 錯誤碼，會映射為 409 Conflict
			return nil, &errInfos.Res{Code: errInfos.INVALID_STATUS, Msg: conflictMsg}, nil
		}
	}

	now := time.Now()
	var updatedCount int64 = 1

	switch req.UpdateMode {
	case "SINGLE":
		if req.Title != nil {
			event.Title = *req.Title
		}
		if req.StartAt != nil {
			event.StartAt = libs.TimeToTaiwan(*req.StartAt)
		}
		if req.EndAt != nil {
			event.EndAt = libs.TimeToTaiwan(*req.EndAt)
		}
		if req.IsAllDay != nil {
			event.IsAllDay = *req.IsAllDay
		}
		if req.ColorHex != nil {
			event.ColorHex = *req.ColorHex
		}
		if req.RecurrenceRule != nil {
			event.RecurrenceRule = *req.RecurrenceRule
		}
		event.UpdatedAt = now
		if err := s.personalEventRepo.Update(ctx, event); err != nil {
			return nil, s.app.Err.New(errInfos.SQL_ERROR), err
		}

	case "FUTURE":
		if event.RecurrenceRule.Type == "" {
			return nil, &errInfos.Res{Code: global.BAD_REQUEST, Msg: "update_mode FUTURE requires a recurring event"}, nil
		}
		updatedCount, err = s.personalEventRepo.UpdateFutureOccurrences(ctx, eventID, teacherID, repositories.UpdateEventRequest{
			Title:    req.Title,
			StartAt:  libs.TimeToTaiwanOrNowIfNil(req.StartAt),
			EndAt:    libs.TimeToTaiwanOrNowIfNil(req.EndAt),
			IsAllDay: req.IsAllDay,
			ColorHex: req.ColorHex,
		}, now)
		if err != nil {
			return nil, s.app.Err.New(errInfos.SQL_ERROR), err
		}

	case "ALL":
		if event.RecurrenceRule.Type == "" {
			return nil, &errInfos.Res{Code: global.BAD_REQUEST, Msg: "update_mode ALL requires a recurring event"}, nil
		}
		updatedCount, err = s.personalEventRepo.UpdateAllOccurrences(ctx, eventID, teacherID, repositories.UpdateEventRequest{
			Title:    req.Title,
			StartAt:  libs.TimeToTaiwanOrNowIfNil(req.StartAt),
			EndAt:    libs.TimeToTaiwanOrNowIfNil(req.EndAt),
			IsAllDay: req.IsAllDay,
			ColorHex: req.ColorHex,
		}, now)
		if err != nil {
			return nil, s.app.Err.New(errInfos.SQL_ERROR), err
		}

	default:
		// 預設行為與 SINGLE 相同
		if req.Title != nil {
			event.Title = *req.Title
		}
		if req.StartAt != nil {
			event.StartAt = libs.TimeToTaiwan(*req.StartAt)
		}
		if req.EndAt != nil {
			event.EndAt = libs.TimeToTaiwan(*req.EndAt)
		}
		if req.IsAllDay != nil {
			event.IsAllDay = *req.IsAllDay
		}
		if req.ColorHex != nil {
			event.ColorHex = *req.ColorHex
		}
		event.UpdatedAt = now
		if err := s.personalEventRepo.Update(ctx, event); err != nil {
			return nil, s.app.Err.New(errInfos.SQL_ERROR), err
		}
	}

	return &UpdatePersonalEventResult{
		UpdatedCount: updatedCount,
		Message:      "Updated " + req.UpdateMode + " occurrence(s)",
	}, nil, nil
}

// DeletePersonalEvent 刪除老師個人行程
func (s *PersonalEventService) DeletePersonalEvent(ctx context.Context, eventID, teacherID uint) *errInfos.Res {
	event, err := s.personalEventRepo.GetByID(ctx, eventID)
	if err != nil {
		return s.app.Err.New(errInfos.NOT_FOUND)
	}

	if event.TeacherID != teacherID {
		return s.app.Err.New(errInfos.FORBIDDEN)
	}

	if err := s.personalEventRepo.DeleteByID(ctx, eventID); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR)
	}

	return nil
}

// GetPersonalEventNote 取得個人行程備註
func (s *PersonalEventService) GetPersonalEventNote(ctx context.Context, eventID, teacherID uint) (string, *errInfos.Res, error) {
	event, err := s.personalEventRepo.GetByID(ctx, eventID)
	if err != nil {
		return "", s.app.Err.New(errInfos.NOT_FOUND), err
	}

	if event.TeacherID != teacherID {
		return "", s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	return event.Note, nil, nil
}

// UpdatePersonalEventNoteRequest 更新備註請求
type UpdatePersonalEventNoteRequest struct {
	Content string
}

// UpdatePersonalEventNote 更新個人行程備註
func (s *PersonalEventService) UpdatePersonalEventNote(ctx context.Context, eventID, teacherID uint, req *UpdatePersonalEventNoteRequest) *errInfos.Res {
	event, err := s.personalEventRepo.GetByID(ctx, eventID)
	if err != nil {
		return s.app.Err.New(errInfos.NOT_FOUND)
	}

	if event.TeacherID != teacherID {
		return s.app.Err.New(errInfos.FORBIDDEN)
	}

	if err := s.personalEventRepo.UpdateNote(ctx, eventID, req.Content); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR)
	}

	return nil
}

// OccurrenceInstance 代表行程在特定日期的實例
type OccurrenceInstance struct {
    EventID       uint           `json:"event_id"`
    Title         string         `json:"title"`
    StartAt       time.Time      `json:"start_at"`
    EndAt         time.Time      `json:"end_at"`
    IsAllDay      bool           `json:"is_all_day"`
    ColorHex      string         `json:"color_hex"`
    Note          string         `json:"note"`
    Date          time.Time      `json:"date"`
    OriginalDate  time.Time      `json:"original_date"` // 原始行程開始日期
}

// GetTodayOccurrences 取得老師在目標日期的所有行程實例
func (s *PersonalEventService) GetTodayOccurrences(ctx context.Context, teacherID uint, targetDate time.Time) ([]OccurrenceInstance, *errInfos.Res, error) {
    // 取得該老師的所有行程
    events, err := s.personalEventRepo.ListByTeacherID(ctx, teacherID)
    if err != nil {
        return nil, s.app.Err.New(errInfos.SQL_ERROR), err
    }

    var occurrences []OccurrenceInstance

    // 取得目標日期的星期幾（1-7，週一到週日）
    targetWeekday := int(targetDate.Weekday())
    if targetWeekday == 0 {
        targetWeekday = 7
    }

    for _, event := range events {
        // 判斷行程是否在目標日期發生
        isOccurring, instance := s.checkEventOccurrenceOnDate(event, targetDate, targetWeekday)
        if isOccurring {
            occurrences = append(occurrences, instance)
        }
    }

    return occurrences, nil, nil
}

// checkEventOccurrenceOnDate 判斷行程是否在目標日期發生
func (s *PersonalEventService) checkEventOccurrenceOnDate(event models.PersonalEvent, targetDate time.Time, targetWeekday int) (bool, OccurrenceInstance) {
    // 處理單次行程（沒有循環規則）
    if event.RecurrenceRule.Type == "" {
        eventDate := time.Date(
            event.StartAt.Year(), event.StartAt.Month(), event.StartAt.Day(),
            0, 0, 0, 0, event.StartAt.Location(),
        )
        targetDateOnly := time.Date(
            targetDate.Year(), targetDate.Month(), targetDate.Day(),
            0, 0, 0, 0, targetDate.Location(),
        )

        if eventDate.Equal(targetDateOnly) {
            return true, OccurrenceInstance{
                EventID:      event.ID,
                Title:        event.Title,
                StartAt:      event.StartAt,
                EndAt:        event.EndAt,
                IsAllDay:     event.IsAllDay,
                ColorHex:     event.ColorHex,
                Note:         event.Note,
                Date:         targetDate,
                OriginalDate: event.StartAt,
            }
        }
        return false, OccurrenceInstance{}
    }

    // 處理有循環規則的行程
    rule := event.RecurrenceRule

    // 檢查 until 限制
    if rule.Until != nil {
        untilDate, parseErr := time.Parse("2006-01-02", *rule.Until)
        if parseErr == nil {
            targetDateOnly := time.Date(
                targetDate.Year(), targetDate.Month(), targetDate.Day(),
                0, 0, 0, 0, targetDate.Location(),
            )
            if targetDateOnly.After(untilDate) {
                return false, OccurrenceInstance{}
            }
        }
    }

    // 計算目標日期距離原始開始日期的天數
    originalDate := time.Date(
        event.StartAt.Year(), event.StartAt.Month(), event.StartAt.Day(),
        0, 0, 0, 0, event.StartAt.Location(),
    )
    targetDateOnly := time.Date(
        targetDate.Year(), targetDate.Month(), targetDate.Day(),
        0, 0, 0, 0, targetDate.Location(),
    )

    // 目標日期必須在原始日期當天或之後
    if targetDateOnly.Before(originalDate) {
        return false, OccurrenceInstance{}
    }

    daysDiff := int(targetDateOnly.Sub(originalDate).Hours() / 24)

    var isValidOccurrence bool

    switch rule.Type {
    case "DAILY":
        // 每日循環：檢查間隔
        isValidOccurrence = daysDiff%rule.Interval == 0

    case "WEEKLY":
        // 每週循環：檢查星期是否匹配
        isValidOccurrence = false
        for _, w := range rule.Weekdays {
            if w == targetWeekday {
                // 檢查間隔（每 N 週）
                if rule.Interval <= 1 {
                    isValidOccurrence = true
                    break
                }
                // 計算週數
                weeksDiff := daysDiff / 7
                if daysDiff%7 == 0 && weeksDiff%rule.Interval == 0 {
                    isValidOccurrence = true
                    break
                }
            }
        }

    case "CUSTOM":
        // 自訂循環：使用 interval 作為天數間隔
        isValidOccurrence = daysDiff%rule.Interval == 0

    default:
        // 未知的循環類型，不發生
        isValidOccurrence = false
    }

    // 檢查 count 限制
    if isValidOccurrence && rule.Count != nil && *rule.Count > 0 {
        occurrences := s.countOccurrences(originalDate, targetDateOnly, rule)
        if occurrences > *rule.Count {
            isValidOccurrence = false
        }
    }

    if !isValidOccurrence {
        return false, OccurrenceInstance{}
    }

    // 計算該實例的實際時間
    instanceStartAt := time.Date(
        targetDate.Year(), targetDate.Month(), targetDate.Day(),
        event.StartAt.Hour(), event.StartAt.Minute(), event.StartAt.Second(), event.StartAt.Nanosecond(),
        event.StartAt.Location(),
    )
    instanceEndAt := time.Date(
        targetDate.Year(), targetDate.Month(), targetDate.Day(),
        event.EndAt.Hour(), event.EndAt.Minute(), event.EndAt.Second(), event.EndAt.Nanosecond(),
        event.EndAt.Location(),
    )

    return true, OccurrenceInstance{
        EventID:      event.ID,
        Title:        event.Title,
        StartAt:      instanceStartAt,
        EndAt:        instanceEndAt,
        IsAllDay:     event.IsAllDay,
        ColorHex:     event.ColorHex,
        Note:         event.Note,
        Date:         targetDate,
        OriginalDate: event.StartAt,
    }
}

// countOccurrences 計算從原始日期到目標日期之間的發生次數
func (s *PersonalEventService) countOccurrences(originalDate, targetDate time.Time, rule models.RecurrenceRule) int {
    if targetDate.Before(originalDate) {
        return 0
    }

    daysDiff := int(targetDate.Sub(originalDate).Hours() / 24)

    switch rule.Type {
    case "DAILY", "CUSTOM":
        // 包括原始日期當天，所以 +1
        return daysDiff/rule.Interval + 1

    case "WEEKLY":
        // 計算完整的週數
        if daysDiff%7 != 0 {
            return 0
        }
        weeksDiff := daysDiff / 7
        if rule.Interval <= 1 {
            return weeksDiff + 1
        }
        // 計算符合間隔的週數
        return weeksDiff/rule.Interval + 1
    }

    return 1
}

// checkPersonalEventConflicts 檢查個人行程是否與中心課程衝突
// 回傳空字串表示無衝突，回傳非空字串表示有衝突（包含衝突訊息）
// 回傳 error 表示系統錯誤
func (s *PersonalEventService) checkPersonalEventConflicts(ctx context.Context, teacherID uint, startAt, endAt time.Time) (string, error) {
	memberships, err := s.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		return "", err
	}

	if len(memberships) == 0 {
		return "", nil // 沒有加入任何中心，直接通過
	}

	var centerIDs []uint
	for _, m := range memberships {
		centerIDs = append(centerIDs, m.CenterID)
	}

	// 檢查每個中心的課程衝突
	for _, centerID := range centerIDs {
		conflicts, err := s.scheduleRuleRepo.CheckPersonalEventConflict(ctx, teacherID, centerID, startAt, endAt)
		if err != nil {
			continue
		}

		if len(conflicts) > 0 {
			conflictMessages := make([]string, 0, len(conflicts))
			for _, rule := range conflicts {
				conflictMessages = append(conflictMessages, fmt.Sprintf(
					"%s %s-%s 在中心 %d 有課程「%s」的安排，時間衝突",
					startAt.Format("2006-01-02"),
					rule.StartTime,
					rule.EndTime,
					centerID,
					rule.Offering.Name,
				))
			}
			return conflictMessages[0], nil // 只顯示第一個衝突
		}
	}

	return "", nil
}
