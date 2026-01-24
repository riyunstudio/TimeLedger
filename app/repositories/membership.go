package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterMembershipRepository struct {
	BaseRepository
	app   *app.App
	model *models.CenterMembership
}

func NewCenterMembershipRepository(app *app.App) *CenterMembershipRepository {
	return &CenterMembershipRepository{
		app: app,
	}
}

func (rp *CenterMembershipRepository) GetByID(ctx context.Context, id uint) (models.CenterMembership, error) {
	var data models.CenterMembership
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *CenterMembershipRepository) GetByCenterAndTeacher(ctx context.Context, centerID, teacherID uint) (models.CenterMembership, error) {
	var data models.CenterMembership
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ? AND teacher_id = ?", centerID, teacherID).First(&data).Error
	return data, err
}

func (rp *CenterMembershipRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterMembership, error) {
	var data []models.CenterMembership
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ?", centerID).Find(&data).Error
	return data, err
}

func (rp *CenterMembershipRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error) {
	var data []models.CenterMembership
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("teacher_id = ?", teacherID).Find(&data).Error
	return data, err
}

func (rp *CenterMembershipRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.CenterMembership, error) {
	var data []models.CenterMembership
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("center_id = ? AND status = ?", centerID, "ACTIVE").Find(&data).Error
	return data, err
}

func (rp *CenterMembershipRepository) GetActiveByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error) {
	var data []models.CenterMembership
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("teacher_id = ? AND status = ?", teacherID, "ACTIVE").Find(&data).Error
	return data, err
}

func (rp *CenterMembershipRepository) ListTeacherIDsByCenterID(ctx context.Context, centerID uint) ([]uint, error) {
	var membershipIDs []uint
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Model(&models.CenterMembership{}).
		Where("center_id = ? AND status = ?", centerID, "ACTIVE").
		Pluck("teacher_id", &membershipIDs).Error
	return membershipIDs, err
}

type CenterMembershipRepositoryInterface interface {
	Create(ctx context.Context, data models.CenterMembership) (models.CenterMembership, error)
	Update(ctx context.Context, data models.CenterMembership) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (models.CenterMembership, error)
	GetByCenterAndTeacher(ctx context.Context, centerID, teacherID uint) (models.CenterMembership, error)
	ListByCenterID(ctx context.Context, centerID uint) ([]models.CenterMembership, error)
	ListByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error)
	ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.CenterMembership, error)
	GetActiveByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error)
	ListTeacherIDsByCenterID(ctx context.Context, centerID uint) ([]uint, error)
}

func (rp *CenterMembershipRepository) Create(ctx context.Context, data models.CenterMembership) (models.CenterMembership, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *CenterMembershipRepository) Update(ctx context.Context, data models.CenterMembership) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *CenterMembershipRepository) Delete(ctx context.Context, id uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Delete(&models.CenterMembership{}, id).Error
}
