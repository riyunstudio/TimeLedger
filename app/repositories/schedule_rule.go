package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

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

func (rp *ScheduleRuleRepository) Create(ctx context.Context, data models.ScheduleRule) (models.ScheduleRule, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *ScheduleRuleRepository) Update(ctx context.Context, data models.ScheduleRule) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *ScheduleRuleRepository) BulkCreate(ctx context.Context, data []models.ScheduleRule) ([]models.ScheduleRule, error) {
	if len(data) == 0 {
		return data, nil
	}
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
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
	return rp.FindWithCenterScope(ctx, centerID)
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
