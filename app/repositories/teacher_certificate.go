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
