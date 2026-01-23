package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterInvitationRepository struct {
	BaseRepository
	app *app.App
}

func NewCenterInvitationRepository(app *app.App) *CenterInvitationRepository {
	return &CenterInvitationRepository{
		app: app,
	}
}

func (rp *CenterInvitationRepository) GetByID(ctx context.Context, id uint) (models.CenterInvitation, error) {
	var data models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *CenterInvitationRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.CenterInvitation, error) {
	var data models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ? AND center_id = ?", id, centerID).First(&data).Error
	return data, err
}

func (rp *CenterInvitationRepository) GetByToken(ctx context.Context, token string) (models.CenterInvitation, error) {
	var data models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("token = ?", token).First(&data).Error
	return data, err
}

func (rp *CenterInvitationRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterInvitation, error) {
	var data []models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *CenterInvitationRepository) ListByCenterIDPaginated(ctx context.Context, centerID uint, page, limit int64, status string) ([]models.CenterInvitation, int64, error) {
	var data []models.CenterInvitation
	var total int64

	query := rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterInvitation{}).Where("center_id = ?", centerID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&data).Error

	return data, total, err
}

func (rp *CenterInvitationRepository) CountByCenterID(ctx context.Context, centerID uint) (int64, error) {
	var count int64
	err := rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterInvitation{}).Where("center_id = ?", centerID).Count(&count).Error
	return count, err
}

func (rp *CenterInvitationRepository) CountByStatus(ctx context.Context, centerID uint, status string) (int64, error) {
	var count int64
	err := rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterInvitation{}).
		Where("center_id = ? AND status = ?", centerID, status).
		Count(&count).Error
	return count, err
}

func (rp *CenterInvitationRepository) CountByDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time, status string) (int64, error) {
	var count int64
	query := rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterInvitation{}).
		Where("center_id = ?", centerID).
		Where("created_at >= ?", startDate).
		Where("created_at <= ?", endDate)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&count).Error
	return count, err
}

func (rp *CenterInvitationRepository) Create(ctx context.Context, data models.CenterInvitation) (models.CenterInvitation, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *CenterInvitationRepository) Update(ctx context.Context, data models.CenterInvitation) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *CenterInvitationRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return rp.app.MySQL.WDB.WithContext(ctx).
		Where("id = ?", id).
		Update("status", status).Error
}

func (rp *CenterInvitationRepository) DeleteByID(ctx context.Context, id uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Delete(&models.CenterInvitation{}, id).Error
}
