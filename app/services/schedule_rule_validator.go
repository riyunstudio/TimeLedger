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

// ScheduleRuleValidator 統一的排課規則驗證服務
// 整合所有檢查邏輯（重疊、緩衝、個人行程）
type ScheduleRuleValidator struct {
	app              *app.App
	validationService ScheduleValidationService
	scheduleRuleRepo *models.ScheduleRule
}

// NewScheduleRuleValidator 建立統一的驗證服務
func NewScheduleRuleValidator(app *app.App) *ScheduleRuleValidator {
	return &ScheduleRuleValidator{
		app:              app,
		validationService: NewScheduleValidationService(app),
	}
}

// ValidationSummary 驗證結果摘要
type ValidationSummary struct {
	Valid             bool                  `json:"valid"`
	OverlapConflicts  []OverlapInfo         `json:"overlap_conflicts,omitempty"`
	BufferConflicts   []BufferInfo          `json:"buffer_conflicts,omitempty"`
	AllConflicts      []ConflictInfo        `json:"all_conflicts,omitempty"`
}

// OverlapInfo 重疊衝突資訊
type OverlapInfo struct {
	Weekday   int    `json:"weekday"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	RuleID    uint   `json:"rule_id,omitempty"`
	RuleWeekday int  `json:"rule_weekday,omitempty"`
	ConflictType string `json:"conflict_type"` // "ROOM_OVERLAP", "TEACHER_OVERLAP"
	Message   string `json:"message"`
}

// BufferInfo 緩衝衝突資訊
type BufferInfo struct {
	Weekday        int    `json:"weekday"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	PrevEndTime    string `json:"prev_end_time,omitempty"`
	RequiredMinutes int   `json:"required_minutes"`
	GapMinutes     int    `json:"gap_minutes"`
	ConflictType   string `json:"conflict_type"` // "TEACHER_BUFFER", "ROOM_BUFFER"
	Message        string `json:"message"`
	CanOverride    bool   `json:"can_override"`
}

// ConflictInfo 完整衝突資訊（用於 ApplyTemplate）
type ConflictInfo struct {
	Weekday      int    `json:"weekday"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	ConflictType string `json:"conflict_type"` // "ROOM_OVERLAP", "TEACHER_OVERLAP", "PERSONAL_EVENT", "TEACHER_BUFFER", "ROOM_BUFFER"
	Message      string `json:"message"`
	RuleID       uint   `json:"rule_id,omitempty"`
	CanOverride  bool   `json:"can_override,omitempty"`
}

// ValidateForApplyTemplate 驗證模板套用的衝突
// 用於 ApplyTemplate API，返回詳細的衝突資訊
// allowOverride: 是否允許管理員覆蓋 Buffer 衝突
func (v *ScheduleRuleValidator) ValidateForApplyTemplate(ctx context.Context, centerID uint, offeringID uint, weekdays []int, cells []models.TimetableCell, startDate, endDate string, allowOverride bool) (*ValidationSummary, error) {
	summary := &ValidationSummary{
		Valid: true,
	}

	parsedStartDate, _ := time.Parse("2006-01-02", startDate)
	
	dayNames := []string{"", "週一", "週二", "週三", "週四", "週五", "週六", "週日"}

	for _, weekday := range weekdays {
		// 找到該 weekday 的第一個日期
		current := parsedStartDate
		weekdayDiff := weekday - int(current.Weekday())
		if weekdayDiff <= 0 {
			weekdayDiff += 7
		}
		targetDate := current.AddDate(0, 0, weekdayDiff)

		for _, cell := range cells {
			roomID := uint(0)
			if cell.RoomID != nil {
				roomID = *cell.RoomID
			}

			// 檢查 Room 和 Teacher Overlap
			overlappingRules, personalEventConflicts, err := v.checkOverlap(ctx, centerID, roomID, cell.TeacherID, weekday, cell.StartTime, cell.EndTime, targetDate)
			if err != nil {
				return nil, fmt.Errorf("failed to check overlap: %w", err)
			}

			// 處理 ScheduleRule 衝突
			for _, rule := range overlappingRules {
				conflictType := "ROOM_OVERLAP"
				if rule.TeacherID != nil && cell.TeacherID != nil && *rule.TeacherID == *cell.TeacherID {
					conflictType = "TEACHER_OVERLAP"
				}

				msg := fmt.Sprintf("%s %s-%s 教室已被佔用", dayNames[weekday], cell.StartTime, cell.EndTime)
				if conflictType == "TEACHER_OVERLAP" {
					msg = fmt.Sprintf("%s %s-%s 老師已有排課", dayNames[weekday], cell.StartTime, cell.EndTime)
				}

				conflict := ConflictInfo{
					Weekday:      weekday,
					StartTime:    cell.StartTime,
					EndTime:      cell.EndTime,
					ConflictType: conflictType,
					Message:      msg,
					RuleID:       rule.ID,
					CanOverride:  false, // Overlap 不可覆蓋
				}
				
				summary.AllConflicts = append(summary.AllConflicts, conflict)
				summary.Valid = false
			}

			// 處理 Personal Event 衝突
			for _, event := range personalEventConflicts {
				msg := fmt.Sprintf("%s %s-%s 「%s」老師個人行程",
					dayNames[weekday],
					event.StartAt.Format("15:04"),
					event.EndAt.Format("15:04"),
					event.Title)

				conflict := ConflictInfo{
					Weekday:      weekday,
					StartTime:    cell.StartTime,
					EndTime:      cell.EndTime,
					ConflictType: "PERSONAL_EVENT",
					Message:      msg,
					CanOverride:  false, // Personal Event 不可覆蓋
				}
				
				summary.AllConflicts = append(summary.AllConflicts, conflict)
				summary.Valid = false
			}

			// 檢查 Buffer（如果課程有設定）
			if offeringID > 0 {
				// 產生新課程的時間
				newStartTime, _ := time.Parse("2006-01-02 15:04", targetDate.Format("2006-01-02")+" "+cell.StartTime)

				// 檢查 Teacher Buffer
				if cell.TeacherID != nil && *cell.TeacherID > 0 {
					prevEndTime, _ := v.getPreviousSessionEndTime(ctx, centerID, *cell.TeacherID, weekday, cell.StartTime)
					if !prevEndTime.IsZero() {
						bufferResult, err := v.validationService.CheckTeacherBuffer(ctx, centerID, *cell.TeacherID, prevEndTime, newStartTime, offeringID)
						if err != nil {
							return nil, fmt.Errorf("failed to check teacher buffer: %w", err)
						}

						if !bufferResult.Valid {
							for _, bc := range bufferResult.Conflicts {
								conflict := ConflictInfo{
									Weekday:      weekday,
									StartTime:    cell.StartTime,
									EndTime:      cell.EndTime,
									ConflictType: "TEACHER_BUFFER",
									Message:      bc.Message,
									CanOverride:  bc.CanOverride && allowOverride,
								}
								summary.AllConflicts = append(summary.AllConflicts, conflict)
								// 如果不允許覆蓋，則標記為無效
								if !allowOverride || !bc.CanOverride {
									summary.Valid = false
								}
							}
						}
					}
				}

				// 檢查 Room Buffer
				if roomID > 0 {
					prevRoomEndTime, _ := v.getPreviousRoomSessionEndTime(ctx, centerID, roomID, weekday, cell.StartTime)
					if !prevRoomEndTime.IsZero() {
						bufferResult, err := v.validationService.CheckRoomBuffer(ctx, centerID, roomID, prevRoomEndTime, newStartTime, offeringID)
						if err != nil {
							return nil, fmt.Errorf("failed to check room buffer: %w", err)
						}

						if !bufferResult.Valid {
							for _, bc := range bufferResult.Conflicts {
								conflict := ConflictInfo{
									Weekday:      weekday,
									StartTime:    cell.StartTime,
									EndTime:      cell.EndTime,
									ConflictType: "ROOM_BUFFER",
									Message:      bc.Message,
									CanOverride:  bc.CanOverride && allowOverride,
								}
								summary.AllConflicts = append(summary.AllConflicts, conflict)
								// 如果不允許覆蓋，則標記為無效
								if !allowOverride || !bc.CanOverride {
									summary.Valid = false
								}
							}
						}
					}
				}
			}
		}
	}

	return summary, nil
}

// ValidateForCreateRule 驗證新規則的衝突
// 用於 CreateRule API，返回結構化的衝突資訊
// allowOverride: 是否允許管理員覆蓋 Buffer 衝突
func (v *ScheduleRuleValidator) ValidateForCreateRule(ctx context.Context, centerID uint, teacherID *uint, roomID uint, offeringID uint, weekdays []int, startDate, endDate, startTime, endTime string, allowOverride bool) (*ValidationSummary, error) {
	summary := &ValidationSummary{
		Valid: true,
	}

	parsedStartDate, _ := time.Parse("2006-01-02", startDate)
	
	dayNames := []string{"", "週一", "週二", "週三", "週四", "週五", "週六", "週日"}

	for _, weekday := range weekdays {
		// 找到該 weekday 的第一個日期
		current := parsedStartDate
		weekdayDiff := weekday - int(current.Weekday())
		if weekdayDiff <= 0 {
			weekdayDiff += 7
		}
		targetDate := current.AddDate(0, 0, weekdayDiff)

		// 檢查 Room 和 Teacher Overlap
		overlappingRules, personalEventConflicts, err := v.checkOverlap(ctx, centerID, roomID, teacherID, weekday, startTime, endTime, targetDate)
		if err != nil {
			return nil, fmt.Errorf("failed to check overlap: %w", err)
		}

		// 處理 ScheduleRule 衝突
		for _, rule := range overlappingRules {
			conflictType := "ROOM_OVERLAP"
			if rule.TeacherID != nil && teacherID != nil && *rule.TeacherID == *teacherID {
				conflictType = "TEACHER_OVERLAP"
			}

			msg := fmt.Sprintf("%s %s-%s", dayNames[weekday], startTime, endTime)
			if conflictType == "TEACHER_OVERLAP" {
				msg += " 老師已有排課"
			} else {
				msg += " 教室已被佔用"
			}

			summary.OverlapConflicts = append(summary.OverlapConflicts, OverlapInfo{
				Weekday:     weekday,
				StartTime:   startTime,
				EndTime:     endTime,
				RuleID:      rule.ID,
				ConflictType: conflictType,
				Message:     msg,
			})
			summary.Valid = false
		}

		// 處理 Personal Event 衝突
		for _, event := range personalEventConflicts {
			msg := fmt.Sprintf("%s %s-%s 「%s」老師個人行程",
				dayNames[weekday],
				event.StartAt.Format("15:04"),
				event.EndAt.Format("15:04"),
				event.Title)

			summary.OverlapConflicts = append(summary.OverlapConflicts, OverlapInfo{
				Weekday:     weekday,
				StartTime:   startTime,
				EndTime:     endTime,
				ConflictType: "PERSONAL_EVENT",
				Message:     msg,
			})
			summary.Valid = false
		}

		// 檢查 Buffer
		if offeringID > 0 {
			newStartTime, _ := time.Parse("2006-01-02 15:04", targetDate.Format("2006-01-02")+" "+startTime)

			// 檢查 Teacher Buffer
			if teacherID != nil && *teacherID > 0 {
				prevEndTime, _ := v.getPreviousSessionEndTime(ctx, centerID, *teacherID, weekday, startTime)
				if !prevEndTime.IsZero() {
					bufferResult, err := v.validationService.CheckTeacherBuffer(ctx, centerID, *teacherID, prevEndTime, newStartTime, offeringID)
					if err != nil {
						return nil, fmt.Errorf("failed to check teacher buffer: %w", err)
					}

					if !bufferResult.Valid {
						for _, bc := range bufferResult.Conflicts {
							summary.BufferConflicts = append(summary.BufferConflicts, BufferInfo{
								Weekday:        weekday,
								StartTime:      startTime,
								EndTime:        endTime,
								PrevEndTime:    prevEndTime.Format("15:04"),
								RequiredMinutes: bc.RequiredMinutes,
								GapMinutes:     bc.DiffMinutes,
								ConflictType:   "TEACHER_BUFFER",
								Message:        bc.Message,
								CanOverride:    bc.CanOverride && allowOverride,
							})
							// 如果不允許覆蓋，則標記為無效
							if !allowOverride || !bc.CanOverride {
								summary.Valid = false
							}
						}
					}
				}
			}

			// 檢查 Room Buffer
			if roomID > 0 {
				prevRoomEndTime, _ := v.getPreviousRoomSessionEndTime(ctx, centerID, roomID, weekday, startTime)
				if !prevRoomEndTime.IsZero() {
					bufferResult, err := v.validationService.CheckRoomBuffer(ctx, centerID, roomID, prevRoomEndTime, newStartTime, offeringID)
					if err != nil {
						return nil, fmt.Errorf("failed to check room buffer: %w", err)
					}

					if !bufferResult.Valid {
						for _, bc := range bufferResult.Conflicts {
							summary.BufferConflicts = append(summary.BufferConflicts, BufferInfo{
								Weekday:        weekday,
								StartTime:      startTime,
								EndTime:        endTime,
								PrevEndTime:    prevRoomEndTime.Format("15:04"),
								RequiredMinutes: bc.RequiredMinutes,
								GapMinutes:     bc.DiffMinutes,
								ConflictType:   "ROOM_BUFFER",
								Message:        bc.Message,
								CanOverride:    bc.CanOverride && allowOverride,
							})
							// 如果不允許覆蓋，則標記為無效
							if !allowOverride || !bc.CanOverride {
								summary.Valid = false
							}
						}
					}
				}
			}
		}
	}

	return summary, nil
}

// checkOverlap 檢查時間重疊（封裝 Repository 的 CheckOverlap）
func (v *ScheduleRuleValidator) checkOverlap(ctx context.Context, centerID uint, roomID uint, teacherID *uint, weekday int, startTime, endTime string, checkDate time.Time) ([]models.ScheduleRule, []models.PersonalEvent, error) {
	// 使用 ScheduleRuleRepository 的 CheckOverlap 方法
	// 這裡需要從 app 取得 repository
	// 由於 validator 不直接持有 repository，我們使用 validationService 的方式來處理
	
	// 簡化實作：直接使用資料庫查詢
	var overlappingRules []models.ScheduleRule
	var personalEventConflicts []models.PersonalEvent

	// 查詢教室重疊
	roomQuery := v.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("room_id = ?", roomID).
		Where("weekday = ?", weekday).
		Where("deleted_at IS NULL").
		Where("start_time < ?", endTime).
		Where("end_time > ?", startTime)

	var roomRules []models.ScheduleRule
	if err := roomQuery.Find(&roomRules).Error; err != nil {
		return nil, nil, err
	}
	overlappingRules = append(overlappingRules, roomRules...)

	// 查詢老師重疊
	if teacherID != nil && *teacherID > 0 {
		teacherQuery := v.app.MySQL.RDB.WithContext(ctx).
			Where("center_id = ?", centerID).
			Where("teacher_id = ?", *teacherID).
			Where("weekday = ?", weekday).
			Where("deleted_at IS NULL").
			Where("start_time < ?", endTime).
			Where("end_time > ?", startTime)

		var teacherRules []models.ScheduleRule
		if err := teacherQuery.Find(&teacherRules).Error; err != nil {
			return nil, nil, err
		}
		
		// 避免重複（如果 room 和 teacher 相同）
		existingIDs := make(map[uint]bool)
		for _, r := range overlappingRules {
			existingIDs[r.ID] = true
		}
		for _, r := range teacherRules {
			if !existingIDs[r.ID] {
				overlappingRules = append(overlappingRules, r)
			}
		}

		// 查詢個人行程衝突
		personalEventRepo := repositories.NewPersonalEventRepository(v.app)
		events, err := personalEventRepo.CheckPersonalEventConflict(ctx, *teacherID, weekday, startTime, endTime, checkDate)
		if err != nil {
			return nil, nil, err
		}
		personalEventConflicts = events
	}

	return overlappingRules, personalEventConflicts, nil
}

// getPreviousSessionEndTime 取得老師在指定 weekday 之前的最後一筆課程結束時間
func (v *ScheduleRuleValidator) getPreviousSessionEndTime(ctx context.Context, centerID, teacherID uint, weekday int, beforeTimeStr string) (time.Time, error) {
	weekdayVal := weekday
	if weekdayVal == 0 {
		weekdayVal = 7
	}

	var rule models.ScheduleRule
	err := v.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("teacher_id = ?", teacherID).
		Where("weekday = ?", weekdayVal).
		Where("end_time <= ?", beforeTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		// 如果找不到紀錄（gorm.ErrRecordNotFound），這是預期情況，回傳零值時間
		// 不要將此視為錯誤，這只是表示該老師在該時段之前沒有排課
		if err == gorm.ErrRecordNotFound {
			return time.Time{}, nil
		}
		// 其他資料庫錯誤才需要記錄
		return time.Time{}, err
	}

	endTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01"+" "+rule.EndTime)
	return endTime, nil
}

// getPreviousRoomSessionEndTime 取得教室在指定 weekday 之前的最後一筆課程結束時間
func (v *ScheduleRuleValidator) getPreviousRoomSessionEndTime(ctx context.Context, centerID, roomID uint, weekday int, beforeTimeStr string) (time.Time, error) {
	weekdayVal := weekday
	if weekdayVal == 0 {
		weekdayVal = 7
	}

	var rule models.ScheduleRule
	err := v.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("room_id = ?", roomID).
		Where("weekday = ?", weekdayVal).
		Where("end_time <= ?", beforeTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		// 如果找不到紀錄（gorm.ErrRecordNotFound），這是預期情況，回傳零值時間
		if err == gorm.ErrRecordNotFound {
			return time.Time{}, nil
		}
		// 其他資料庫錯誤才需要記錄
		return time.Time{}, err
	}

	endTime, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01"+" "+rule.EndTime)
	return endTime, nil
}
