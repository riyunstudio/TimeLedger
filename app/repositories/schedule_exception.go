package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

type ScheduleExceptionRepository struct {
	GenericRepository[models.ScheduleException]
	app *app.App
}

func NewScheduleExceptionRepository(app *app.App) *ScheduleExceptionRepository {
	return &ScheduleExceptionRepository{
		GenericRepository: NewGenericRepository[models.ScheduleException](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// Transaction executes a function within a database transaction.
// This method creates a NEW ScheduleExceptionRepository instance with transaction connections.
func (rp *ScheduleExceptionRepository) Transaction(ctx context.Context, fn func(txRepo *ScheduleExceptionRepository) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &ScheduleExceptionRepository{
			GenericRepository: NewTransactionRepo[models.ScheduleException](ctx, tx, rp.table),
			app:               rp.app,
		}
		return fn(txRepo)
	})
}

// CreateWithDB creates a new record using the provided database connection (for transaction support).
func (rp *ScheduleExceptionRepository) CreateWithDB(ctx context.Context, db *gorm.DB, exception models.ScheduleException) (models.ScheduleException, error) {
	if err := db.WithContext(ctx).Table("schedule_exceptions").Create(&exception).Error; err != nil {
		return models.ScheduleException{}, err
	}
	return exception, nil
}

func (rp *ScheduleExceptionRepository) GetByRuleAndDate(ctx context.Context, ruleID uint, date time.Time) ([]models.ScheduleException, error) {
	return rp.Find(ctx, "rule_id = ? AND DATE(original_date) = ?", ruleID, date.Format("2006-01-02"))
}

func (rp *ScheduleExceptionRepository) GetByRuleAndDateAndCenterID(ctx context.Context, ruleID uint, centerID uint, date time.Time) ([]models.ScheduleException, error) {
	result, err := rp.FindWithCenterScope(ctx, centerID, "rule_id = ? AND original_date = ?", ruleID, date)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rp *ScheduleExceptionRepository) GetByTeacherID(ctx context.Context, teacherID uint, status string) ([]models.ScheduleException, error) {
	var exceptions []models.ScheduleException

	// 只取得該老師創建的例外申請（透過 schedule_rules.teacher_id 關聯）
	// 這樣確保老師只會看到自己申請的例外，不會看到同中心其他老師的申請
	query := rp.app.MySQL.RDB.WithContext(ctx).
		Table("schedule_exceptions").
		Select("schedule_exceptions.*").
		Joins("JOIN schedule_rules ON schedule_rules.id = schedule_exceptions.rule_id").
		Where("schedule_rules.teacher_id = ?", teacherID).
		Where("schedule_rules.deleted_at IS NULL")

	if status != "" {
		if status == "APPROVED" {
			query = query.Where("schedule_exceptions.status IN ('APPROVED', 'APPROVE')")
		} else if status == "REJECTED" {
			query = query.Where("schedule_exceptions.status IN ('REJECTED', 'REJECT')")
		} else {
			query = query.Where("schedule_exceptions.status = ?", status)
		}
	}

	err := query.Order("schedule_exceptions.created_at DESC").Find(&exceptions).Error
	return exceptions, err
}

func (rp *ScheduleExceptionRepository) BatchGetByRuleIDs(ctx context.Context, ruleIDs []uint, date time.Time) (map[uint][]models.ScheduleException, error) {
	if len(ruleIDs) == 0 {
		return make(map[uint][]models.ScheduleException), nil
	}

	exceptions, err := rp.Find(ctx, "rule_id IN ? AND DATE(original_date) = ?", ruleIDs, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	result := make(map[uint][]models.ScheduleException, len(ruleIDs))
	for _, exc := range exceptions {
		result[exc.RuleID] = append(result[exc.RuleID], exc)
	}
	return result, nil
}

func (rp *ScheduleExceptionRepository) GetByRuleIDsAndDateRange(ctx context.Context, ruleIDs []uint, startDate, endDate time.Time) (map[uint]map[string][]models.ScheduleException, error) {
	if len(ruleIDs) == 0 {
		return make(map[uint]map[string][]models.ScheduleException), nil
	}

	exceptions, err := rp.Find(ctx, "rule_id IN ? AND DATE(original_date) >= ? AND DATE(original_date) <= ?", ruleIDs, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	result := make(map[uint]map[string][]models.ScheduleException)
	for _, ruleID := range ruleIDs {
		result[ruleID] = make(map[string][]models.ScheduleException)
	}

	for _, exc := range exceptions {
		dateStr := exc.OriginalDate.Format("2006-01-02")
		if result[exc.RuleID] == nil {
			result[exc.RuleID] = make(map[string][]models.ScheduleException)
		}
		result[exc.RuleID][dateStr] = append(result[exc.RuleID][dateStr], exc)
	}

	return result, nil
}

func (rp *ScheduleExceptionRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.ScheduleException, error) {
	return rp.FindWithCenterScope(ctx, centerID)
}
