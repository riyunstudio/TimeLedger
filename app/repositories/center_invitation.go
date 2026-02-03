package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterInvitationRepository struct {
	GenericRepository[models.CenterInvitation]
	app *app.App
}

func NewCenterInvitationRepository(app *app.App) *CenterInvitationRepository {
	return &CenterInvitationRepository{
		GenericRepository: NewGenericRepository[models.CenterInvitation](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *CenterInvitationRepository) GetByToken(ctx context.Context, token string) (models.CenterInvitation, error) {
	return rp.First(ctx, "token = ?", token)
}

func (rp *CenterInvitationRepository) GetByTeacherAndCenter(ctx context.Context, teacherID, centerID uint) ([]models.CenterInvitation, error) {
	return rp.Find(ctx, "teacher_id = ? AND center_id = ?", teacherID, centerID)
}

func (rp *CenterInvitationRepository) GetPendingByTeacher(ctx context.Context, teacherID uint) ([]models.CenterInvitation, error) {
	return rp.Find(ctx, "teacher_id = ? AND status = ?", teacherID, models.InvitationStatusPending)
}

func (rp *CenterInvitationRepository) HasPendingInvitation(ctx context.Context, teacherID, centerID uint) (bool, error) {
	return rp.Exists(ctx, "teacher_id = ? AND center_id = ? AND status = ?", teacherID, centerID, models.InvitationStatusPending)
}

func (rp *CenterInvitationRepository) GetPendingByCenter(ctx context.Context, centerID uint) ([]models.CenterInvitation, error) {
	return rp.Find(ctx, "center_id = ? AND status = ?", centerID, models.InvitationStatusPending)
}

func (rp *CenterInvitationRepository) ExpireOldInvitations(ctx context.Context, beforeTime time.Time) (int64, error) {
	result := rp.app.MySQL.WDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("status = ? AND expires_at < ?", models.InvitationStatusPending, beforeTime).
		Update("status", models.InvitationStatusExpired)
	return result.RowsAffected, result.Error
}

func (rp *CenterInvitationRepository) CountByCenter(ctx context.Context, centerID uint) (pending, accepted, declined int64, err error) {
	// Count pending
	pending, err = rp.Count(ctx, "center_id = ? AND status = ?", centerID, models.InvitationStatusPending)
	if err != nil {
		return 0, 0, 0, err
	}

	// Count accepted
	accepted, err = rp.Count(ctx, "center_id = ? AND status = ?", centerID, models.InvitationStatusAccepted)
	if err != nil {
		return 0, 0, 0, err
	}

	// Count declined
	declined, err = rp.Count(ctx, "center_id = ? AND status = ?", centerID, models.InvitationStatusDeclined)
	if err != nil {
		return 0, 0, 0, err
	}

	return pending, accepted, declined, nil
}

func (rp *CenterInvitationRepository) UpdateStatus(ctx context.Context, id uint, status models.InvitationStatus) error {
	now := time.Now()
	fields := map[string]interface{}{
		"status": status,
	}

	// 當狀態變更為 ACCEPTED 或 DECLINED 時，設定 RespondedAt
	if status == models.InvitationStatusAccepted || status == models.InvitationStatusDeclined {
		fields["responded_at"] = now
	}

	return rp.UpdateFields(ctx, id, fields)
}

func (rp *CenterInvitationRepository) CountByStatus(ctx context.Context, centerID uint, status models.InvitationStatus) (int64, error) {
	return rp.Count(ctx, "center_id = ? AND status = ?", centerID, status)
}

func (rp *CenterInvitationRepository) CountByCenterID(ctx context.Context, centerID uint) (int64, error) {
	return rp.Count(ctx, "center_id = ?", centerID)
}

func (rp *CenterInvitationRepository) CountByDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time) (int64, error) {
	return rp.Count(ctx, "center_id = ? AND created_at >= ? AND created_at <= ?", centerID, startDate, endDate)
}

func (rp *CenterInvitationRepository) ListByCenterIDPaginated(ctx context.Context, centerID uint, page, limit int, status string) ([]models.CenterInvitation, int64, error) {
	return rp.FindPaged(ctx, page, limit, "created_at DESC", "center_id = ?", centerID)
}

func (rp *CenterInvitationRepository) ListByTeacher(ctx context.Context, teacherID uint) ([]models.CenterInvitation, error) {
	return rp.Find(ctx, "teacher_id = ?", teacherID)
}

func (rp *CenterInvitationRepository) ListByTeacherWithStatus(ctx context.Context, teacherID uint, status string) ([]models.CenterInvitation, error) {
	if status != "" {
		return rp.Find(ctx, "teacher_id = ? AND status = ?", teacherID, status)
	}
	return rp.Find(ctx, "teacher_id = ?", teacherID)
}

// GetGeneralByCenterID 取得中心的通用邀請
func (rp *CenterInvitationRepository) GetGeneralByCenterID(ctx context.Context, centerID uint) (models.CenterInvitation, error) {
	return rp.First(ctx, "center_id = ? AND invite_type = ?", centerID, models.InvitationTypeGeneral)
}
