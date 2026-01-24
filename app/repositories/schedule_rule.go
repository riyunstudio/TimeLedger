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
