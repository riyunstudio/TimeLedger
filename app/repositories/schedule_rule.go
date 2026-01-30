package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type ScheduleRuleRepository struct {
	BaseRepository
	app   *app.App
	model *models.ScheduleRule
}

func NewScheduleRuleRepository(app *app.App) *ScheduleRuleRepository {
	return &ScheduleRuleRepository{
		app: app,
	}
}

func (rp *ScheduleRuleRepository) GetByID(ctx context.Context, id uint) (models.ScheduleRule, error) {
	var data models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Room").
		Preload("Teacher").
		Preload("Exceptions").
		Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.ScheduleRule, error) {
	var data models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Room").
		Preload("Teacher").
		Preload("Exceptions").
		Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Room").
		Preload("Teacher").
		Where("center_id = ?", centerID).
		Order("weekday ASC, start_time ASC").
		Find(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) ListByTeacherID(ctx context.Context, teacherID uint, centerID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Room").
		Preload("Teacher").
		Where("teacher_id = ? AND center_id = ?", teacherID, centerID).
		Order("weekday ASC, start_time ASC").
		Find(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) ListByRoomID(ctx context.Context, roomID uint, centerID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("room_id = ? AND center_id = ?", roomID, centerID).
		Find(&data).Error
	return data, err
}

// CheckOverlap checks if there's an overlapping rule with the same room or teacher, or personal events
func (rp *ScheduleRuleRepository) CheckOverlap(ctx context.Context, centerID uint, roomID uint, teacherID *uint, weekday int, startTime string, endTime string, excludeRuleID *uint, checkDate time.Time) ([]models.ScheduleRule, []models.PersonalEvent, error) {
	var data []models.ScheduleRule
	var personalEventConflicts []models.PersonalEvent

	// 判斷是否為跨日課程
	isCrossDay := IsCrossDayTime(startTime, endTime)

	query := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("deleted_at IS NULL")

	// 如果是跨日課程，需要檢查：
	// 1. 當天（weekday）的課程中，開始時間在 20:00 以後的
	// 2. 下一天（weekday+1）的課程中，結束時間在 08:00 以前的
	if isCrossDay {
		// 檢查跨日衝突 - 當天晚間時段
		nextWeekday := weekday + 1
		if nextWeekday > 7 {
			nextWeekday = 1
		}

		// 查詢當天晚間課程（20:00 以後開始）
		eveningQuery := query.Where("weekday = ?", weekday)
		if excludeRuleID != nil {
			eveningQuery = eveningQuery.Where("id != ?", *excludeRuleID)
		}
		var eveningRules []models.ScheduleRule
		if err := eveningQuery.Find(&eveningRules).Error; err != nil {
			return nil, nil, err
		}

		for _, rule := range eveningRules {
			// 如果規則也是跨日課程，或開始時間在 20:00 以後
			if rule.IsCrossDay || ParseTimeToMinutes(rule.StartTime) >= 20*60 {
				if TimesOverlapCrossDay(startTime, endTime, true, rule.StartTime, rule.EndTime, rule.IsCrossDay) {
					data = append(data, rule)
				}
			}
		}

		// 查詢隔天凌晨課程（08:00 以前結束）
		morningQuery := query.Where("weekday = ?", nextWeekday)
		if excludeRuleID != nil {
			morningQuery = morningQuery.Where("id != ?", *excludeRuleID)
		}
		var morningRules []models.ScheduleRule
		if err := morningQuery.Find(&morningRules).Error; err != nil {
			return nil, nil, err
		}

		for _, rule := range morningRules {
			// 如果規則是跨日課程，或結束時間在 08:00 以前
			if rule.IsCrossDay || ParseTimeToMinutes(rule.EndTime) <= 8*60 {
				// 使用 TimesOverlapCrossDayWithNextDay 來比較隔天課程
				if TimesOverlapCrossDayWithNextDay(startTime, endTime, true, rule.StartTime, rule.EndTime) {
					// 避免重複
					exists := false
					for _, existing := range data {
						if existing.ID == rule.ID {
							exists = true
							break
						}
					}
					if !exists {
						data = append(data, rule)
					}
				}
			}
		}
	} else {
		// 普通課程的衝突檢測邏輯（保持原有行為）
		roomQuery := query.Where("weekday = ?", weekday)
		if excludeRuleID != nil {
			roomQuery = roomQuery.Where("id != ?", *excludeRuleID)
		}

		var roomRules []models.ScheduleRule
		if err := roomQuery.Find(&roomRules).Error; err != nil {
			return nil, nil, err
		}

		// Check time overlap for room rules
		for _, rule := range roomRules {
			if TimesOverlapCrossDay(startTime, endTime, false, rule.StartTime, rule.EndTime, rule.IsCrossDay) {
				data = append(data, rule)
			}
		}
	}

	// Check teacher overlap (if teacher_id is provided)
	if teacherID != nil && *teacherID != 0 {
		// 對於跨日課程，還需要檢查老師在隔天的課程
		if isCrossDay {
			nextWeekday := weekday + 1
			if nextWeekday > 7 {
				nextWeekday = 1
			}

			// 分開查詢當天和隔天的課程
			sameDayQuery := query.Where("teacher_id = ? AND weekday = ?", *teacherID, weekday)
			if excludeRuleID != nil {
				sameDayQuery = sameDayQuery.Where("id != ?", *excludeRuleID)
			}
			var sameDayRules []models.ScheduleRule
			if err := sameDayQuery.Find(&sameDayRules).Error; err != nil {
				return nil, nil, err
			}

			nextDayQuery := query.Where("teacher_id = ? AND weekday = ?", *teacherID, nextWeekday)
			if excludeRuleID != nil {
				nextDayQuery = nextDayQuery.Where("id != ?", *excludeRuleID)
			}
			var nextDayRules []models.ScheduleRule
			if err := nextDayQuery.Find(&nextDayRules).Error; err != nil {
				return nil, nil, err
			}

			// 檢查當天課程（使用 TimesOverlapCrossDay）
			for _, rule := range sameDayRules {
				if rule.IsCrossDay || ParseTimeToMinutes(rule.StartTime) >= 20*60 {
					if TimesOverlapCrossDay(startTime, endTime, true, rule.StartTime, rule.EndTime, rule.IsCrossDay) {
						// Avoid duplicate entries if room and teacher are the same
						exists := false
						for _, existing := range data {
							if existing.ID == rule.ID {
								exists = true
								break
							}
						}
						if !exists {
							data = append(data, rule)
						}
					}
				}
			}

			// 檢查隔天課程（使用 TimesOverlapCrossDayWithNextDay）
			for _, rule := range nextDayRules {
				if rule.IsCrossDay || ParseTimeToMinutes(rule.EndTime) <= 8*60 {
					if TimesOverlapCrossDayWithNextDay(startTime, endTime, true, rule.StartTime, rule.EndTime) {
						// Avoid duplicate entries if room and teacher are the same
						exists := false
						for _, existing := range data {
							if existing.ID == rule.ID {
								exists = true
								break
							}
						}
						if !exists {
							data = append(data, rule)
						}
					}
				}
			}

			// 檢查老師的個人行程（跨日）
			personalEventRepo := NewPersonalEventRepository(rp.app)
			events, err := personalEventRepo.CheckPersonalEventConflictCrossDay(ctx, *teacherID, weekday, startTime, endTime, checkDate)
			if err != nil {
				return nil, nil, err
			}
			personalEventConflicts = events
		} else {
			teacherQuery := query.Where("teacher_id = ?", *teacherID)
			if excludeRuleID != nil {
				teacherQuery = teacherQuery.Where("id != ?", *excludeRuleID)
			}

			var teacherRules []models.ScheduleRule
			if err := teacherQuery.Find(&teacherRules).Error; err != nil {
				return nil, nil, err
			}

			// Check time overlap for teacher rules
			for _, rule := range teacherRules {
				if TimesOverlapCrossDay(startTime, endTime, false, rule.StartTime, rule.EndTime, rule.IsCrossDay) {
					// Avoid duplicate entries if room and teacher are the same
					exists := false
					for _, existing := range data {
						if existing.ID == rule.ID {
							exists = true
							break
						}
					}
					if !exists {
						data = append(data, rule)
					}
				}
			}

			// Check personal events for this teacher
			personalEventRepo := NewPersonalEventRepository(rp.app)
			events, err := personalEventRepo.CheckPersonalEventConflict(ctx, *teacherID, weekday, startTime, endTime, checkDate)
			if err != nil {
				return nil, nil, err
			}
			personalEventConflicts = events
		}
	}

	return data, personalEventConflicts, nil
}

func (rp *ScheduleRuleRepository) Create(ctx context.Context, data models.ScheduleRule) (models.ScheduleRule, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) Update(ctx context.Context, data models.ScheduleRule) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *ScheduleRuleRepository) UpdateByIDAndCenterID(ctx context.Context, id uint, centerID uint, data models.ScheduleRule) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).Save(&data).Error
}

func (rp *ScheduleRuleRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Delete(&models.ScheduleRule{}, id).Error
}

func (rp *ScheduleRuleRepository) DeleteByIDAndCenterID(ctx context.Context, id uint, centerID uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).Delete(&models.ScheduleRule{}).Error
}

func (rp *ScheduleRuleRepository) ListByOfferingID(ctx context.Context, offeringID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("offering_id = ?", offeringID).
		Order("effective_range ASC").
		Find(&data).Error
	return data, err
}

// ListByOfferingIDWithPreload 批次查詢規則並預載入關聯資料（消除 N+1 查詢）
func (rp *ScheduleRuleRepository) ListByOfferingIDWithPreload(ctx context.Context, offeringID uint) ([]models.ScheduleRule, error) {
	var data []models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Room").
		Preload("Teacher").
		Where("offering_id = ?", offeringID).
		Order("effective_range ASC").
		Find(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) BulkCreate(ctx context.Context, data []models.ScheduleRule) ([]models.ScheduleRule, error) {
	if len(data) == 0 {
		return data, nil
	}
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

// CheckPersonalEventConflict 檢查個人行程是否與老師的課程衝突
// 用於在創建個人行程前檢查是否會與已排課程重疊
func (rp *ScheduleRuleRepository) CheckPersonalEventConflict(ctx context.Context, teacherID uint, centerID uint, eventStartAt time.Time, eventEndAt time.Time) ([]models.ScheduleRule, error) {
	var conflicts []models.ScheduleRule

	// 取得該老師在該中心的課程規則
	rules, err := rp.ListByTeacherID(ctx, teacherID, centerID)
	if err != nil {
		return nil, err
	}

	// 將事件時間轉換為台北時區，確保 weekday 計算正確
	loc := app.GetTaiwanLocation()
	eventStartAtTaipei := eventStartAt.In(loc)
	eventEndAtTaipei := eventEndAt.In(loc)

	// 取得事件的星期幾（1-7，週一到週日），使用台灣時區
	eventWeekday := int(eventStartAtTaipei.Weekday())
	if eventWeekday == 0 {
		eventWeekday = 7 // 週日轉換為 7
	}

	// 將事件時間轉換為台北時區的 HH:MM 格式
	eventStartTime := eventStartAtTaipei.Format("15:04")
	eventEndTime := eventEndAtTaipei.Format("15:04")

	// 事件日期（使用台灣時區）
	eventDate := eventStartAtTaipei.Format("2006-01-02")

	for _, rule := range rules {
		// 只檢查同一天的規則
		if rule.Weekday != eventWeekday {
			continue
		}

		// 檢查事件日期是否在課程規則的有效範圍內
		// 規則的有效範圍也應該轉換為台灣時區進行比較
		ruleStartDate := rule.EffectiveRange.StartDate.In(loc).Format("2006-01-02")
		ruleEndDate := rule.EffectiveRange.EndDate.In(loc).Format("2006-01-02")

		if eventDate < ruleStartDate || eventDate > ruleEndDate {
			// 事件日期不在課程規則的有效範圍內，跳過
			continue
		}

		// 檢查時間是否重疊
		overlaps := timesOverlap(rule.StartTime, rule.EndTime, eventStartTime, eventEndTime)
		if overlaps {
			conflicts = append(conflicts, rule)
		}
	}

	return conflicts, nil
}

// CheckPersonalEventConflictAllCenters 檢查個人行程是否與老師所有中心的課程衝突
func (rp *ScheduleRuleRepository) CheckPersonalEventConflictAllCenters(ctx context.Context, teacherID uint, centerIDs []uint, eventStartAt time.Time, eventEndAt time.Time) (map[uint][]models.ScheduleRule, error) {
	allConflicts := make(map[uint][]models.ScheduleRule)

	for _, centerID := range centerIDs {
		conflicts, err := rp.CheckPersonalEventConflict(ctx, teacherID, centerID, eventStartAt, eventEndAt)
		if err != nil {
			return nil, err
		}
		if len(conflicts) > 0 {
			allConflicts[centerID] = conflicts
		}
	}

	return allConflicts, nil
}

// IsCrossDayTime 檢查是否為跨日時間（結束時間早於開始時間）
func IsCrossDayTime(startTime, endTime string) bool {
	// 如果開始時間和結束時間相同，視為非跨日（全天課程的特例）
	if startTime == endTime {
		return false
	}
	return endTime < startTime
}

// timesOverlap checks if two time ranges overlap
func timesOverlap(start1, end1, start2, end2 string) bool {
	return start1 < end2 && end1 > start2
}

// GetLastSessionByTeacherAndWeekday 取得老師在指定星期幾的最後一堂課結束時間
func (rp *ScheduleRuleRepository) GetLastSessionByTeacherAndWeekday(ctx context.Context, centerID, teacherID uint, weekday int, beforeTimeStr string) (*models.ScheduleRule, error) {
	weekdayVal := weekday
	if weekdayVal == 0 {
		weekdayVal = 7
	}

	var rule models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("teacher_id = ?", teacherID).
		Where("weekday = ?", weekdayVal).
		Where("end_time <= ?", beforeTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetLastSessionByRoomAndWeekday 取得教室在指定星期幾的最後一堂課結束時間
func (rp *ScheduleRuleRepository) GetLastSessionByRoomAndWeekday(ctx context.Context, centerID, roomID uint, weekday int, beforeTimeStr string) (*models.ScheduleRule, error) {
	weekdayVal := weekday
	if weekdayVal == 0 {
		weekdayVal = 7
	}

	var rule models.ScheduleRule
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("room_id = ?", roomID).
		Where("weekday = ?", weekdayVal).
		Where("end_time <= ?", beforeTimeStr).
		Order("end_time DESC").
		First(&rule).Error

	if err != nil {
		return nil, err
	}
	return &rule, nil
}
