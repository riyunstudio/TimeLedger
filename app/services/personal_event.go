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
	// 檢查個人行程是否與中心課程衝突
	conflictMsg, err := s.checkPersonalEventConflicts(ctx, teacherID, req.StartAt, req.EndAt)
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
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
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

	// 檢查衝突（如果時間有變更）
	if req.StartAt != nil && req.EndAt != nil {
		conflictMsg, err := s.checkPersonalEventConflicts(ctx, teacherID, *req.StartAt, *req.EndAt)
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
			event.StartAt = *req.StartAt
		}
		if req.EndAt != nil {
			event.EndAt = *req.EndAt
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
			StartAt:  req.StartAt,
			EndAt:    req.EndAt,
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
			StartAt:  req.StartAt,
			EndAt:    req.EndAt,
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
			event.StartAt = *req.StartAt
		}
		if req.EndAt != nil {
			event.EndAt = *req.EndAt
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
