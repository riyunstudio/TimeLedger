package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/libs"

	"gorm.io/gorm"
)

type ScheduleRuleRepository struct {
	GenericRepository[models.ScheduleRule]
	app *app.App
}

func NewScheduleRuleRepository(app *app.App) *ScheduleRuleRepository {
	return &ScheduleRuleRepository{
		GenericRepository: NewGenericRepository[models.ScheduleRule](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// Transaction executes a function within a database transaction.
// This method creates a NEW ScheduleRuleRepository instance with transaction connections
// to avoid race conditions in concurrent requests.
//
// Usage Example:
//
//	result, err := rp.Transaction(ctx, func(txRepo *ScheduleRuleRepository) error {
//	    // All operations using txRepo will be within the same transaction
//	    // Custom methods like BulkCreate, ListByTeacherID are available
//	    if _, err := txRepo.BulkCreate(ctx, rules); err != nil {
//	        return err
//	    }
//	    if _, err := txRepo.DeleteByIDAndCenterID(ctx, oldRuleID, centerID); err != nil {
//	        return err
//	    }
//	    return nil
//	})
func (rp *ScheduleRuleRepository) Transaction(ctx context.Context, fn func(txRepo *ScheduleRuleRepository) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new ScheduleRuleRepository instance with transaction connections
		txRepo := &ScheduleRuleRepository{
			GenericRepository: GenericRepository[models.ScheduleRule]{
				dbRead:  tx.WithContext(ctx),
				dbWrite: tx.WithContext(ctx),
				table:   rp.table,
			},
			app: rp.app,
		}
		return fn(txRepo)
	})
}

func (rp *ScheduleRuleRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.ScheduleRule, error) {
	return rp.GetByIDWithCenterScope(ctx, id, centerID)
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
	return rp.FindWithCenterScope(ctx, centerID, "room_id = ?", roomID)
}

func (rp *ScheduleRuleRepository) ListByOfferingID(ctx context.Context, offeringID uint) ([]models.ScheduleRule, error) {
	return rp.Find(ctx, "offering_id = ?", offeringID)
}

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

func IsCrossDayTime(startTime, endTime string) bool {
	if startTime == endTime {
		return false
	}
	return endTime < startTime
}

func timesOverlap(start1, end1, start2, end2 string) bool {
	return start1 < end2 && end1 > start2
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

// ListByCenterIDPaginated 分頁取得排課規則（優化版）
// category 參數：用於過濾特定課程類別
// search 參數：用於搜尋課程名稱、老師名稱、教室名稱
// weekday 參數：用於過濾星期幾（1-7）
// status 參數：用於過濾狀態（upcoming/ongoing/ended）
func (rp *ScheduleRuleRepository) ListByCenterIDPaginated(ctx context.Context, centerID uint, page, limit int, category, search string, weekday int, status string) ([]models.ScheduleRule, int64, error) {
	var data []models.ScheduleRule
	var total int64

	loc := libs.GetTaiwanLocation()
	now := time.Now().In(loc)

	// ============================================
	// 第一步：取得符合條件的 ID 列表（優化 COUNT 查詢）
	// ============================================
	// 使用子查詢避免重複 JOIN，先取得符合條件的 rule IDs
	subQuery := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.ScheduleRule{}).
		Select("schedule_rules.id").
		Where("schedule_rules.center_id = ?", centerID)

	// 只有 category 需要 JOIN courses 表
	hasCategoryFilter := category != ""
	hasSearchFilter := search != ""
	hasWeekdayFilter := weekday > 0
	hasStatusFilter := status != ""

	if hasCategoryFilter {
		subQuery = subQuery.
			Joins("JOIN offerings ON schedule_rules.offering_id = offerings.id").
			Joins("JOIN courses ON offerings.course_id = courses.id").
			Where("courses.category = ?", category)
	}

	// 搜尋過濾：使用 EXISTS 子查詢或直接 JOIN
	if hasSearchFilter {
		searchPattern := "%" + search + "%"
		// 搜尋時需要 JOIN 相關表
		subQuery = subQuery.
			Joins("LEFT JOIN offerings ON schedule_rules.offering_id = offerings.id").
			Joins("LEFT JOIN rooms ON schedule_rules.room_id = rooms.id").
			Joins("LEFT JOIN teachers ON schedule_rules.teacher_id = teachers.id").
			Where(
				"offerings.name LIKE ? OR teachers.name LIKE ? OR rooms.name LIKE ?",
				searchPattern, searchPattern, searchPattern,
			)
	}

	// 星期過濾（不需要 JOIN）
	if hasWeekdayFilter {
		subQuery = subQuery.Where("schedule_rules.weekday = ?", weekday)
	}

	// 狀態過濾（使用 JSON_EXTRACT 查詢 effective_range）
	if hasStatusFilter {
		nowStr := now.Format("2006-01-02")
		switch status {
		case "upcoming":
			// 尚未開始：effective_range.start_date > now
			subQuery = subQuery.Where("JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.start_date')) > ?", nowStr)
		case "ongoing":
			// 進行中：effective_range.start_date <= now <= effective_range.end_date
			subQuery = subQuery.Where("JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.start_date')) <= ? AND (JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) IS NULL OR JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) >= ?)", nowStr, nowStr)
		case "ended":
			// 已結束：effective_range.end_date < now
			subQuery = subQuery.Where("JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) IS NOT NULL AND JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) < ?", nowStr)
		}
	}

	// 取得總數
	subQuery = subQuery.Order("schedule_rules.weekday ASC, schedule_rules.start_time ASC")
	if err := subQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 如果沒有符合條件的資料，直接返回
	if total == 0 {
		return data, total, nil
	}

	// ============================================
	// 第二步：取得分頁資料（使用 Preload 載入關聯）
	// ============================================
	offset := (page - 1) * limit

	// 取得符合條件的 ID 列表
	var ruleIDs []uint
	idQuery := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.ScheduleRule{}).
		Select("schedule_rules.id").
		Where("schedule_rules.center_id = ?", centerID)

	if hasCategoryFilter {
		idQuery = idQuery.
			Joins("JOIN offerings ON schedule_rules.offering_id = offerings.id").
			Joins("JOIN courses ON offerings.course_id = courses.id").
			Where("courses.category = ?", category)
	}

	if hasSearchFilter {
		searchPattern := "%" + search + "%"
		idQuery = idQuery.
			Joins("LEFT JOIN offerings ON schedule_rules.offering_id = offerings.id").
			Joins("LEFT JOIN rooms ON schedule_rules.room_id = rooms.id").
			Joins("LEFT JOIN teachers ON schedule_rules.teacher_id = teachers.id").
			Where(
				"offerings.name LIKE ? OR teachers.name LIKE ? OR rooms.name LIKE ?",
				searchPattern, searchPattern, searchPattern,
			)
	}

	if hasWeekdayFilter {
		idQuery = idQuery.Where("schedule_rules.weekday = ?", weekday)
	}

	// 狀態過濾（使用 JSON_EXTRACT 查詢 effective_range）
	if hasStatusFilter {
		nowStr := now.Format("2006-01-02")
		switch status {
		case "upcoming":
			idQuery = idQuery.Where("JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.start_date')) > ?", nowStr)
		case "ongoing":
			idQuery = idQuery.Where("JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.start_date')) <= ? AND (JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) IS NULL OR JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) >= ?)", nowStr, nowStr)
		case "ended":
			idQuery = idQuery.Where("JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) IS NOT NULL AND JSON_UNQUOTE(JSON_EXTRACT(schedule_rules.effective_range, '$.end_date')) < ?", nowStr)
		}
	}

	idQuery = idQuery.Order("schedule_rules.weekday ASC, schedule_rules.start_time ASC").
		Offset(offset).Limit(limit)

	if err := idQuery.Pluck("id", &ruleIDs).Error; err != nil {
		return nil, 0, err
	}

	// 如果沒有 ID，直接返回
	if len(ruleIDs) == 0 {
		return data, total, nil
	}

	// 使用 IN 查詢取得資料，並用 Preload 載入關聯（比 JOIN 更高效）
	query := rp.app.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Offering.Course").
		Preload("Room").
		Preload("Teacher").
		Where("schedule_rules.id IN ?", ruleIDs).
		Order("schedule_rules.weekday ASC, schedule_rules.start_time ASC")

	err := query.Find(&data).Error

	return data, total, err
}

// CheckPersonalEventConflict 檢查個人行程是否與排課規則衝突
func (rp *ScheduleRuleRepository) CheckPersonalEventConflict(ctx context.Context, teacherID, centerID uint, startAt, endAt time.Time) ([]models.ScheduleRule, error) {
	// 取得教師在該中心的所有規則
	rules, err := rp.ListByTeacherID(ctx, teacherID, centerID)
	if err != nil {
		return nil, err
	}

	var conflicts []models.ScheduleRule
	startTimeStr := startAt.Format("15:04")
	endTimeStr := endAt.Format("15:04")
	weekday := int(startAt.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	for _, rule := range rules {
		// 檢查星期是否匹配
		if rule.Weekday != weekday {
			continue
		}

		// 檢查時間是否重疊
		if rp.timesOverlap(rule.StartTime, rule.EndTime, startTimeStr, endTimeStr) {
			conflicts = append(conflicts, rule)
		}
	}

	return conflicts, nil
}

func (rp *ScheduleRuleRepository) timesOverlap(start1, end1, start2, end2 string) bool {
	return start1 < end2 && end1 > start2
}

// DeleteByIDAndCenterID 刪除規則（帶 center_id 檢查）
func (rp *ScheduleRuleRepository) DeleteByIDAndCenterID(ctx context.Context, id, centerID uint) error {
	return rp.DeleteByIDWithCenterScope(ctx, id, centerID)
}
