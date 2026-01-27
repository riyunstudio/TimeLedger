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

// Create 创建邀请记录
func (rp *CenterInvitationRepository) Create(ctx context.Context, data models.CenterInvitation) (models.CenterInvitation, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

// GetByID 根据ID获取邀请
func (rp *CenterInvitationRepository) GetByID(ctx context.Context, id uint) (models.CenterInvitation, error) {
	var data models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

// GetByToken 根据Token获取邀请
func (rp *CenterInvitationRepository) GetByToken(ctx context.Context, token string) (models.CenterInvitation, error) {
	var data models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("token = ?", token).First(&data).Error
	return data, err
}

// GetByTeacherAndCenter 获取指定老师和中心的邀请记录
func (rp *CenterInvitationRepository) GetByTeacherAndCenter(ctx context.Context, teacherID, centerID uint) ([]models.CenterInvitation, error) {
	var data []models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("teacher_id = ? AND center_id = ?", teacherID, centerID).
		Order("created_at DESC").
		Find(&data).Error
	return data, err
}

// GetPendingByTeacher 获取老師所有待處理的邀請
func (rp *CenterInvitationRepository) GetPendingByTeacher(ctx context.Context, teacherID uint) ([]models.CenterInvitation, error) {
	var data []models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("teacher_id = ? AND status = ?", teacherID, models.InvitationStatusPending).
		Order("created_at DESC").
		Find(&data).Error
	return data, err
}

// HasPendingInvitation 检查是否有待处理的邀请
func (rp *CenterInvitationRepository) HasPendingInvitation(ctx context.Context, teacherID, centerID uint) (bool, error) {
	var count int64
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("teacher_id = ? AND center_id = ? AND status = ?", teacherID, centerID, models.InvitationStatusPending).
		Count(&count).Error
	return count > 0, err
}

// UpdateStatus 更新邀请状态
func (rp *CenterInvitationRepository) UpdateStatus(ctx context.Context, id uint, status models.InvitationStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == models.InvitationStatusAccepted || status == models.InvitationStatusDeclined {
		now := time.Now()
		updates["responded_at"] = &now
	}
	return rp.app.MySQL.WDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CountByCenter 统计中心的邀请数据
func (rp *CenterInvitationRepository) CountByCenter(ctx context.Context, centerID uint) (pending, accepted, declined int64, err error) {
	// 待處理
	if err = rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("center_id = ? AND status = ?", centerID, models.InvitationStatusPending).
		Count(&pending).Error; err != nil {
		return
	}
	// 已接受
	if err = rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("center_id = ? AND status = ?", centerID, models.InvitationStatusAccepted).
		Count(&accepted).Error; err != nil {
		return
	}
	// 已拒絕
	if err = rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("center_id = ? AND status = ?", centerID, models.InvitationStatusDeclined).
		Count(&declined).Error; err != nil {
		return
	}
	return
}

// GetByCenter 获取中心的所有邀请
func (rp *CenterInvitationRepository) GetByCenter(ctx context.Context, centerID uint, limit, offset int) ([]models.CenterInvitation, error) {
	var data []models.CenterInvitation
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	return data, err
}

// ExpireOldInvitations 使過期邀請失效
func (rp *CenterInvitationRepository) ExpireOldInvitations(ctx context.Context, beforeTime time.Time) (int64, error) {
	result := rp.app.MySQL.WDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("status = ? AND expires_at < ?", models.InvitationStatusPending, beforeTime).
		Update("status", models.InvitationStatusExpired)
	return result.RowsAffected, result.Error
}

// CountByCenterID 統計中心的所有邀請數量
func (rp *CenterInvitationRepository) CountByCenterID(ctx context.Context, centerID uint) (int64, error) {
	var count int64
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("center_id = ?", centerID).
		Count(&count).Error
	return count, err
}

// CountByStatus 統計特定狀態的邀請數量
func (rp *CenterInvitationRepository) CountByStatus(ctx context.Context, centerID uint, status models.InvitationStatus) (int64, error) {
	var count int64
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("center_id = ? AND status = ?", centerID, status).
		Count(&count).Error
	return count, err
}

// CountByDateRange 統計日期範圍內特定狀態的邀請數量
func (rp *CenterInvitationRepository) CountByDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time, status string) (int64, error) {
	var count int64
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("center_id = ? AND status = ? AND created_at BETWEEN ? AND ?", centerID, status, startDate, endDate).
		Count(&count).Error
	return count, err
}

// ListByCenterIDPaginated 分頁取得中心的邀請列表
func (rp *CenterInvitationRepository) ListByCenterIDPaginated(ctx context.Context, centerID uint, page, limit int64, status string) ([]models.CenterInvitation, int64, error) {
	var total int64
	var data []models.CenterInvitation

	// 計算總數
	query := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterInvitation{}).
		Where("center_id = ?", centerID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 計算偏移量
	offset := (int(page) - 1) * int(limit)

	// 取得資料
	findQuery := rp.app.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID)

	if status != "" {
		findQuery = findQuery.Where("status = ?", status)
	}

	err := findQuery.
		Order("created_at DESC").
		Limit(int(limit)).
		Offset(offset).
		Find(&data).Error

	return data, total, err
}
