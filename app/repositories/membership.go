package repositories

import (
	"context"
	"fmt"
	"timeLedger/app"
	"timeLedger/app/models"
)

type CenterMembershipRepository struct {
	GenericRepository[models.CenterMembership]
	app *app.App
}

func NewCenterMembershipRepository(app *app.App) *CenterMembershipRepository {
	return &CenterMembershipRepository{
		GenericRepository: NewGenericRepository[models.CenterMembership](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *CenterMembershipRepository) GetByCenterAndTeacher(ctx context.Context, centerID, teacherID uint) (models.CenterMembership, error) {
	return rp.FirstWithCenterScope(ctx, centerID, "teacher_id = ?", teacherID)
}

func (rp *CenterMembershipRepository) GetActiveByTeacherAndCenter(ctx context.Context, teacherID uint, centerIDStr string) (*models.CenterMembership, error) {
	centerID := 0
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil || centerID == 0 {
		return nil, fmt.Errorf("invalid center_id")
	}
	data, err := rp.FirstWithCenterScope(ctx, uint(centerID), "teacher_id = ? AND status = ?", teacherID, "ACTIVE")
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (rp *CenterMembershipRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error) {
	return rp.Find(ctx, "teacher_id = ?", teacherID)
}

func (rp *CenterMembershipRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.CenterMembership, error) {
	return rp.FindWithCenterScope(ctx, centerID, "status = ?", "ACTIVE")
}

func (rp *CenterMembershipRepository) GetActiveByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error) {
	return rp.Find(ctx, "teacher_id = ? AND status = ?", teacherID, "ACTIVE")
}

func (rp *CenterMembershipRepository) ListTeacherIDsByCenterID(ctx context.Context, centerID uint) ([]uint, error) {
	var membershipIDs []uint
	_, err := rp.FindWithCenterScope(ctx, centerID, "status IN ?", []string{"ACTIVE", "INVITED"})
	if err != nil {
		return nil, err
	}
	err = rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterMembership{}).Where("center_id = ? AND status IN ?", centerID, []string{"ACTIVE", "INVITED"}).Pluck("teacher_id", &membershipIDs).Error
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
