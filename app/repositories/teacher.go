package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/resources"

	"gorm.io/gorm"
)

type TeacherRepository struct {
	GenericRepository[models.Teacher]
	app *app.App
}

func NewTeacherRepository(app *app.App) *TeacherRepository {
	return &TeacherRepository{
		GenericRepository: NewGenericRepository[models.Teacher](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

// Transaction executes a function within a database transaction.
// This method creates a NEW TeacherRepository instance with transaction connections.
func (rp *TeacherRepository) Transaction(ctx context.Context, fn func(txRepo *TeacherRepository) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &TeacherRepository{
			GenericRepository: NewTransactionRepo[models.Teacher](ctx, tx, rp.table),
			app:               rp.app,
		}
		return fn(txRepo)
	})
}

// GetByLineUserID retrieves a teacher by their LINE user ID.
// This is the primary authentication method for teachers.
func (rp *TeacherRepository) GetByLineUserID(ctx context.Context, lineUserID string) (models.Teacher, error) {
	var data models.Teacher
	err := rp.dbRead.WithContext(ctx).
		Unscoped().
		Where("line_user_id = ?", lineUserID).
		First(&data).Error
	return data, err
}

// GetCenterID retrieves the center ID for a teacher with ACTIVE membership.
func (rp *TeacherRepository) GetCenterID(ctx context.Context, teacherID uint) (uint, error) {
	var membership models.CenterMembership
	err := rp.dbRead.WithContext(ctx).
		Table("center_memberships").
		Where("teacher_id = ? AND status = ?", teacherID, "ACTIVE").
		First(&membership).Error
	if err != nil {
		return 0, err
	}
	return membership.CenterID, nil
}

// BatchGetByIDs efficiently retrieves multiple teacher records by their IDs.
// This reduces N+1 query problems when loading multiple teachers.
func (rp *TeacherRepository) BatchGetByIDs(ctx context.Context, ids []uint) (map[uint]models.Teacher, error) {
	if len(ids) == 0 {
		return make(map[uint]models.Teacher), nil
	}

	var teachers []models.Teacher
	err := rp.dbRead.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&teachers).Error
	if err != nil {
		return nil, err
	}

	// Convert to map for O(1) lookup
	result := make(map[uint]models.Teacher, len(teachers))
	for _, teacher := range teachers {
		result[teacher.ID] = teacher
	}
	return result, nil
}

// ListPersonalHashtags retrieves personal hashtags for a teacher.
func (rp *TeacherRepository) ListPersonalHashtags(ctx context.Context, teacherID uint) ([]resources.PersonalHashtag, error) {
	var hashtags []resources.PersonalHashtag
	err := rp.dbRead.WithContext(ctx).
		Table("teacher_personal_hashtags").
		Select("teacher_personal_hashtags.hashtag_id, hashtags.name, teacher_personal_hashtags.sort_order").
		Joins("INNER JOIN hashtags ON teacher_personal_hashtags.hashtag_id = hashtags.id").
		Where("teacher_personal_hashtags.teacher_id = ?", teacherID).
		Order("teacher_personal_hashtags.sort_order ASC").
		Find(&hashtags).Error
	return hashtags, err
}

// DeleteAllPersonalHashtags removes all personal hashtag associations for a teacher.
func (rp *TeacherRepository) DeleteAllPersonalHashtags(ctx context.Context, teacherID uint) error {
	return rp.dbWrite.WithContext(ctx).
		Table("teacher_personal_hashtags").
		Where("teacher_id = ?", teacherID).
		Delete(&models.TeacherPersonalHashtag{}).Error
}

// CreatePersonalHashtag creates a personal hashtag association for a teacher.
func (rp *TeacherRepository) CreatePersonalHashtag(ctx context.Context, teacherID, hashtagID uint, sortOrder int) error {
	return rp.dbWrite.WithContext(ctx).
		Table("teacher_personal_hashtags").
		Create(&models.TeacherPersonalHashtag{
			TeacherID: teacherID,
			HashtagID: hashtagID,
			SortOrder: sortOrder,
		}).Error
}

// ListByCenter retrieves all teachers belonging to a specific center.
func (rp *TeacherRepository) ListByCenter(ctx context.Context, centerID uint) ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := rp.dbRead.WithContext(ctx).
		Table("teachers").
		Joins("INNER JOIN center_memberships ON center_memberships.teacher_id = teachers.id AND center_memberships.center_id = ? AND center_memberships.status = ?", centerID, "ACTIVE").
		Find(&teachers).Error
	return teachers, err
}

// SearchByName searches for teachers by name (partial match, case-insensitive).
func (rp *TeacherRepository) SearchByName(ctx context.Context, name string, limit int) ([]models.Teacher, error) {
	if limit <= 0 {
		limit = 20
	}
	var teachers []models.Teacher
	err := rp.dbRead.WithContext(ctx).
		Where("name LIKE ?", "%"+name+"%").
		Limit(limit).
		Find(&teachers).Error
	return teachers, err
}

// List retrieves all teachers (use with caution, prefer ListByCenter for specific use cases).
func (rp *TeacherRepository) List(ctx context.Context) ([]models.Teacher, error) {
	return rp.FindAll(ctx)
}
