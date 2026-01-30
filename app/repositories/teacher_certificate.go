package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TeacherCertificateRepository struct {
	GenericRepository[models.TeacherCertificate]
	app *app.App
}

func NewTeacherCertificateRepository(app *app.App) *TeacherCertificateRepository {
	return &TeacherCertificateRepository{
		GenericRepository: NewGenericRepository[models.TeacherCertificate](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *TeacherCertificateRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.TeacherCertificate, error) {
	return r.Find(ctx, "teacher_id = ?", teacherID)
}

// BatchListByTeacherIDs 批次查詢多位老師的證照
func (r *TeacherCertificateRepository) BatchListByTeacherIDs(ctx context.Context, teacherIDs []uint) (map[uint][]models.TeacherCertificate, error) {
	if len(teacherIDs) == 0 {
		return make(map[uint][]models.TeacherCertificate), nil
	}

	var certificates []models.TeacherCertificate
	err := r.app.MySQL.RDB.WithContext(ctx).Where("teacher_id IN ?", teacherIDs).Find(&certificates).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint][]models.TeacherCertificate, len(teacherIDs))
	for _, cert := range certificates {
		result[cert.TeacherID] = append(result[cert.TeacherID], cert)
	}
	return result, nil
}
