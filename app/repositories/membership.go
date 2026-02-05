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
	// 使用 Pluck 直接提取 teacher_id，並過濾 active/invited 狀態與 soft delete
	// GORM 的 DeletedAt 會自動添加 deleted_at IS NULL 條件
	err := rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterMembership{}).
		Where("center_id = ? AND status IN ?", centerID, []string{"ACTIVE", "INVITED"}).
		Pluck("teacher_id", &membershipIDs).Error
	return membershipIDs, err
}

// ListTeachersByCenterPaginated 取得中心的老師列表（分頁）
func (rp *CenterMembershipRepository) ListTeachersByCenterPaginated(ctx context.Context, centerID uint, query string, page, limit int) ([]models.Teacher, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// 計算總數
	var total int64
	countQuery := rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterMembership{}).
		Joins("INNER JOIN teachers ON center_memberships.teacher_id = teachers.id").
		Where("center_memberships.center_id = ? AND center_memberships.status IN ?", centerID, []string{"ACTIVE", "INVITED"})
	if query != "" {
		countQuery = countQuery.Where("teachers.name LIKE ?", "%"+query+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查詢資料
	var memberships []models.CenterMembership
	dataQuery := rp.app.MySQL.RDB.WithContext(ctx).Model(&models.CenterMembership{}).
		Joins("INNER JOIN teachers ON center_memberships.teacher_id = teachers.id").
		Where("center_memberships.center_id = ? AND center_memberships.status IN ?", centerID, []string{"ACTIVE", "INVITED"})
	if query != "" {
		dataQuery = dataQuery.Where("teachers.name LIKE ?", "%"+query+"%")
	}
	if err := dataQuery.Order("teachers.name ASC").Offset(offset).Limit(limit).Find(&memberships).Error; err != nil {
		return nil, 0, err
	}

	// 提取老師 IDs 並取得老師資料
	teacherIDs := make([]uint, 0, len(memberships))
	for _, m := range memberships {
		teacherIDs = append(teacherIDs, m.TeacherID)
	}

	// 批次查詢老師資料
	var teachers []models.Teacher
	if len(teacherIDs) > 0 {
		err := rp.app.MySQL.RDB.WithContext(ctx).Where("id IN ?", teacherIDs).Find(&teachers).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return teachers, total, nil
}

type CenterMembershipRepositoryInterface interface {
	Create(ctx context.Context, data models.CenterMembership) (models.CenterMembership, error)
	Update(ctx context.Context, data models.CenterMembership) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (models.CenterMembership, error)
	GetByCenterAndTeacher(ctx context.Context, centerID, teacherID uint) (models.CenterMembership, error)
	ListByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error)
	ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.CenterMembership, error)
	GetActiveByTeacherID(ctx context.Context, teacherID uint) ([]models.CenterMembership, error)
	ListTeacherIDsByCenterID(ctx context.Context, centerID uint) ([]uint, error)
	ListTeachersByCenterPaginated(ctx context.Context, centerID uint, query string, page, limit int) ([]models.Teacher, int64, error)
}
