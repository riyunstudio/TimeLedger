package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TeacherSkillRepository struct {
	BaseRepository
	app *app.App
}

func NewTeacherSkillRepository(app *app.App) *TeacherSkillRepository {
	return &TeacherSkillRepository{app: app}
}

func (r *TeacherSkillRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.TeacherSkill, error) {
	var skills []models.TeacherSkill
	err := r.app.MySQL.RDB.WithContext(ctx).Preload("Hashtags.Hashtag").Where("teacher_id = ?", teacherID).Find(&skills).Error
	return skills, err
}

func (r *TeacherSkillRepository) Create(ctx context.Context, skill *models.TeacherSkill) error {
	return r.app.MySQL.WDB.WithContext(ctx).Create(skill).Error
}

func (r *TeacherSkillRepository) Update(ctx context.Context, skill *models.TeacherSkill) error {
	return r.app.MySQL.WDB.WithContext(ctx).Save(skill).Error
}

func (r *TeacherSkillRepository) Delete(ctx context.Context, id uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).Delete(&models.TeacherSkill{}, id).Error
}

func (r *TeacherSkillRepository) GetByID(ctx context.Context, id uint) (*models.TeacherSkill, error) {
	var skill models.TeacherSkill
	err := r.app.MySQL.RDB.WithContext(ctx).Preload("Hashtags.Hashtag").First(&skill, id).Error
	return &skill, err
}

// BatchListByTeacherIDs 批量取得多個教師的技能資料（效能優化：減少 N+1 查詢）
func (r *TeacherSkillRepository) BatchListByTeacherIDs(ctx context.Context, teacherIDs []uint) (map[uint][]models.TeacherSkill, error) {
	if len(teacherIDs) == 0 {
		return make(map[uint][]models.TeacherSkill), nil
	}

	var skills []models.TeacherSkill
	err := r.app.MySQL.RDB.WithContext(ctx).Preload("Hashtags.Hashtag").Where("teacher_id IN ?", teacherIDs).Find(&skills).Error
	if err != nil {
		return nil, err
	}

	// 按教師 ID 分組
	result := make(map[uint][]models.TeacherSkill, len(teacherIDs))
	for _, skill := range skills {
		result[skill.TeacherID] = append(result[skill.TeacherID], skill)
	}
	return result, nil
}
