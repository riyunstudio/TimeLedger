package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type AuditLogRepository struct {
	GenericRepository[models.AuditLog]
	app *app.App
}

func NewAuditLogRepository(app *app.App) *AuditLogRepository {
	return &AuditLogRepository{
		GenericRepository: NewGenericRepository[models.AuditLog](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *AuditLogRepository) Create(ctx context.Context, log models.AuditLog) (models.AuditLog, error) {
	log.Timestamp = time.Now()
	return r.GenericRepository.Create(ctx, log)
}

func (r *AuditLogRepository) GetByID(ctx context.Context, id uint) (models.AuditLog, error) {
	return r.GenericRepository.GetByID(ctx, id)
}

func (r *AuditLogRepository) Delete(ctx context.Context, id uint) error {
	return r.GenericRepository.DeleteByID(ctx, id)
}

func (r *AuditLogRepository) ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.AuditLog, error) {
	return r.Find(ctx, "center_id = ? AND timestamp >= ? AND timestamp <= ?", centerID, start, end)
}

func (r *AuditLogRepository) CountByCenter(ctx context.Context, centerID uint) (int64, error) {
	return r.Count(ctx, "center_id = ?", centerID)
}

func (r *AuditLogRepository) ListByCenterID(ctx context.Context, centerID uint, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.app.MySQL.RDB.WithContext(ctx).Where("center_id = ?", centerID)
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
	query := r.app.MySQL.RDB.WithContext(ctx).Where("actor_type = ? AND actor_id = ?", actorType, actorID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

type AuditLogRepositoryInterface interface {
	Create(ctx context.Context, log models.AuditLog) (models.AuditLog, error)
	GetByID(ctx context.Context, id uint) (models.AuditLog, error)
	Delete(ctx context.Context, id uint) error
	ListByCenterID(ctx context.Context, centerID uint, limit, offset int) ([]models.AuditLog, error)
	ListByActor(ctx context.Context, actorType string, actorID uint, limit int) ([]models.AuditLog, error)
	ListByDateRange(ctx context.Context, centerID uint, start, end time.Time) ([]models.AuditLog, error)
	CountByCenter(ctx context.Context, centerID uint) (int64, error)
}

var _ AuditLogRepositoryInterface = (*AuditLogRepository)(nil)
