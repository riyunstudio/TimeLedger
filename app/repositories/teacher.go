package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TeacherRepository struct {
	BaseRepository
	app   *app.App
	model *models.Teacher
}

func NewTeacherRepository(app *app.App) *TeacherRepository {
	return &TeacherRepository{
		app: app,
	}
}

func (rp *TeacherRepository) GetByID(ctx context.Context, id uint) (models.Teacher, error) {
	var data models.Teacher
	err := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (rp *TeacherRepository) GetByLineUserID(ctx context.Context, lineUserID string) (models.Teacher, error) {
	var data models.Teacher

	println("[DEBUG] GetByLineUserID searching for:", lineUserID)
	println("[DEBUG] String length:", len(lineUserID))

	// 嘗試不同的查詢方式
	// 1. 先用 ID 測試（繞過 line_user_id）
	var byID models.Teacher
	idErr := rp.app.MySQL.RDB.WithContext(ctx).Where("id = ?", 21).First(&byID).Debug().Error
	println("[DEBUG] Query by ID 21, err:", idErr, "teacher name:", byID.Name)

	// 2. 用 line_user_id 查詢
	err := rp.app.MySQL.RDB.WithContext(ctx).Unscoped().Where("line_user_id = ?", lineUserID).First(&data).Debug().Error
	println("[DEBUG] Query by line_user_id, err:", err)

	// 3. 嘗試用原始 SQL
	var rawResult struct {
		ID         uint
		LineUserID string
		Name       string
	}
	rawErr := rp.app.MySQL.RDB.Raw("SELECT id, line_user_id, name FROM teachers WHERE line_user_id = ?", lineUserID).Scan(&rawResult).Debug().Error
	println("[DEBUG] Raw SQL query, err:", rawErr, "result:", rawResult.ID, rawResult.Name)

	return data, err
}

func (rp *TeacherRepository) List(ctx context.Context) ([]models.Teacher, error) {
	var data []models.Teacher
	err := rp.app.MySQL.RDB.WithContext(ctx).Find(&data).Error
	return data, err
}

func (rp *TeacherRepository) Create(ctx context.Context, data models.Teacher) (models.Teacher, error) {
	err := rp.app.MySQL.WDB.WithContext(ctx).Create(&data).Error
	return data, err
}

func (rp *TeacherRepository) Update(ctx context.Context, data models.Teacher) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Save(&data).Error
}

func (rp *TeacherRepository) DeleteByID(ctx context.Context, id uint) error {
	return rp.app.MySQL.WDB.WithContext(ctx).Where("id = ?", id).Delete(&models.Teacher{}).Error
}

func (rp *TeacherRepository) GetCenterID(ctx context.Context, teacherID uint) (uint, error) {
	var membership models.CenterMembership
	err := rp.app.MySQL.WDB.WithContext(ctx).
		Where("teacher_id = ? AND status = ?", teacherID, "ACTIVE").
		First(&membership).Error
	if err != nil {
		return 0, err
	}
	return membership.CenterID, nil
}

func (rp *TeacherRepository) ListPersonalHashtags(ctx context.Context, teacherID uint) ([]models.Hashtag, error) {
	var hashtags []models.Hashtag
	err := rp.app.MySQL.RDB.WithContext(ctx).
		Table("teacher_personal_hashtags").
		Select("hashtags.*").
		Joins("INNER JOIN hashtags ON teacher_personal_hashtags.hashtag_id = hashtags.id").
		Where("teacher_personal_hashtags.teacher_id = ?", teacherID).
		Order("teacher_personal_hashtags.sort_order ASC").
		Find(&hashtags).Error
	return hashtags, err
}
