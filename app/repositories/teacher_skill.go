package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TeacherSkillRepository struct {
	GenericRepository[models.TeacherSkill]
	app *app.App
}

func NewTeacherSkillRepository(app *app.App) *TeacherSkillRepository {
	return &TeacherSkillRepository{
		GenericRepository: NewGenericRepository[models.TeacherSkill](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *TeacherSkillRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.TeacherSkill, error) {
	return r.Find(ctx, "teacher_id = ?", teacherID)
}

func (r *TeacherSkillRepository) BatchListByTeacherIDs(ctx context.Context, teacherIDs []uint) (map[uint][]models.TeacherSkill, error) {
	if len(teacherIDs) == 0 {
		return make(map[uint][]models.TeacherSkill), nil
	}

	var skills []models.TeacherSkill
	err := r.app.MySQL.RDB.WithContext(ctx).Preload("Hashtags.Hashtag").Where("teacher_id IN ?", teacherIDs).Find(&skills).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint][]models.TeacherSkill, len(teacherIDs))
	for _, skill := range skills {
		result[skill.TeacherID] = append(result[skill.TeacherID], skill)
	}
	return result, nil
}

func (r *TeacherSkillRepository) DeleteAllHashtags(ctx context.Context, skillID uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).
		Where("teacher_skill_id = ?", skillID).
		Delete(&models.TeacherSkillHashtag{}).Error
}

func (r *TeacherSkillRepository) CreateHashtag(ctx context.Context, skillID, hashtagID uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).
		Create(&models.TeacherSkillHashtag{
			TeacherSkillID: skillID,
			HashtagID:      hashtagID,
		}).Error
}
