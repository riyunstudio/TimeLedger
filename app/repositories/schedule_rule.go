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
		Where("teacher_id = ? AND center_id = ?", teacherID, centerID).
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

	query := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("weekday = ?", weekday).
		Where("deleted_at IS NULL")

	// Check room overlap
	roomQuery := query.Where("room_id = ?", roomID)
	if excludeRuleID != nil {
		roomQuery = roomQuery.Where("id != ?", *excludeRuleID)
	}

	var roomRules []models.ScheduleRule
	if err := roomQuery.Find(&roomRules).Error; err != nil {
		return nil, nil, err
	}

	// Check time overlap for room rules
	for _, rule := range roomRules {
		if timesOverlap(rule.StartTime, rule.EndTime, startTime, endTime) {
			data = append(data, rule)
		}
	}

	// Check teacher overlap (if teacher_id is provided)
	if teacherID != nil && *teacherID != 0 {
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
			if timesOverlap(rule.StartTime, rule.EndTime, startTime, endTime) {
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

	return data, personalEventConflicts, nil
}

// timesOverlap checks if two time ranges overlap
func timesOverlap(start1, end1, start2, end2 string) bool {
	// Simple string comparison works for HH:MM format
	// Returns true if ranges overlap (not just adjacent)
	return start1 < end2 && end1 > start2
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

	// 取得事件的星期幾（1-7，週一到週日）
	eventWeekday := int(eventStartAt.Weekday())
	if eventWeekday == 0 {
		eventWeekday = 7 // 週日轉換為 7
	}

	// 將事件時間轉換為台北時區的 HH:MM 格式
	// 課程規則的時間是台北時區，需要統一比較
	loc, _ := time.LoadLocation("Asia/Taipei")
	eventStartAtTaipei := eventStartAt.In(loc)
	eventEndAtTaipei := eventEndAt.In(loc)
	eventStartTime := eventStartAtTaipei.Format("15:04")
	eventEndTime := eventEndAtTaipei.Format("15:04")

	for _, rule := range rules {
		// 只檢查同一天的規則
		if rule.Weekday != eventWeekday {
			continue
		}

		// 檢查事件日期是否在課程規則的有效範圍內
		eventDate := eventStartAtTaipei.Format("2006-01-02")
		ruleStartDate := rule.EffectiveRange.StartDate.Format("2006-01-02")
		ruleEndDate := rule.EffectiveRange.EndDate.Format("2006-01-02")

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
