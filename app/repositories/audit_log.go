package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type AuditLogRepository struct {
	app *app.App
}

func NewAuditLogRepository(app *app.App) *AuditLogRepository {
	return &AuditLogRepository{
		app: app,
	}
}

func (r *AuditLogRepository) Create(ctx context.Context, log models.AuditLog) (models.AuditLog, error) {
	log.Timestamp = time.Now()
	err := r.app.Mysql.WDB.WithContext(ctx).Create(&log).Error
	return log, err
}

func (r *AuditLogRepository) ListByCenterID(ctx context.Context, centerID uint, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.app.Mysql.RDB.WithContext(ctx).Where("center_id = ?", centerID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) ListByActor(ctx context.Context, actorType string, actorID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.app.Mysql.RDB.WithContext(ctx).Where("actor_type = ? AND actor_id = ?", actorType, actorID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.app.Mysql.RDB.WithContext(ctx).
		Where("center_id = ? AND timestamp >= ? AND timestamp <= ?", centerID, start, end).
		Order("timestamp DESC").
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) GetByID(ctx context.Context, id uint) (models.AuditLog, error) {
	var log models.AuditLog
	err := r.app.Mysql.RDB.WithContext(ctx).First(&log, id).Error
	return log, err
}

func (r *AuditLogRepository) Delete(ctx context.Context, id uint) error {
	return r.app.Mysql.WDB.WithContext(ctx).Delete(&models.AuditLog{}, id).Error
}

func (r *AuditLogRepository) CountByCenter(ctx context.Context, centerID uint) (int64, error) {
	var count int64
	err := r.app.Mysql.RDB.WithContext(ctx).Model(&models.AuditLog{}).
		Where("center_id = ?", centerID).
		Count(&count).Error
	return count, err
}

type AuditLogRepositoryInterface interface {
	Create(ctx context.Context, log models.AuditLog) (models.AuditLog, error)
	ListByCenterID(ctx context.Context, centerID uint, limit, offset int) ([]models.AuditLog, error)
	ListByActor(ctx context.Context, actorType string, actorID uint, limit int) ([]models.AuditLog, error)
	ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.AuditLog, error)
	GetByID(ctx context.Context, id uint) (models.AuditLog, error)
	Delete(ctx context.Context, id uint) error
	CountByCenter(ctx context.Context, centerID uint) (int64, error)
}

var _ AuditLogRepositoryInterface = (*AuditLogRepository)(nil)
