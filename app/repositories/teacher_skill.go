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
	err := r.app.Mysql.RDB.WithContext(ctx).Preload("Hashtags.Hashtag").Where("teacher_id = ?", teacherID).Find(&skills).Error
	return skills, err
}

func (r *TeacherSkillRepository) Create(ctx context.Context, skill *models.TeacherSkill) error {
	return r.app.Mysql.WDB.WithContext(ctx).Create(skill).Error
}

func (r *TeacherSkillRepository) Update(ctx context.Context, skill *models.TeacherSkill) error {
	return r.app.Mysql.WDB.WithContext(ctx).Save(skill).Error
}

func (r *TeacherSkillRepository) Delete(ctx context.Context, id uint) error {
	return r.app.Mysql.WDB.WithContext(ctx).Delete(&models.TeacherSkill{}, id).Error
}

func (r *TeacherSkillRepository) GetByID(ctx context.Context, id uint) (*models.TeacherSkill, error) {
	var skill models.TeacherSkill
	err := r.app.Mysql.RDB.WithContext(ctx).Preload("Hashtags.Hashtag").First(&skill, id).Error
	return &skill, err
}
