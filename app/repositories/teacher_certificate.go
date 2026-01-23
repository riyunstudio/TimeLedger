package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TeacherCertificateRepository struct {
	BaseRepository
	app *app.App
}

func NewTeacherCertificateRepository(app *app.App) *TeacherCertificateRepository {
	return &TeacherCertificateRepository{app: app}
}

func (r *TeacherCertificateRepository) ListByTeacherID(ctx context.Context, teacherID uint) ([]models.TeacherCertificate, error) {
	var certificates []models.TeacherCertificate
	err := r.app.MySQL.RDB.WithContext(ctx).Where("teacher_id = ?", teacherID).Find(&certificates).Error
	return certificates, err
}

func (r *TeacherCertificateRepository) Create(ctx context.Context, certificate *models.TeacherCertificate) error {
	return r.app.MySQL.WDB.WithContext(ctx).Create(certificate).Error
}

func (r *TeacherCertificateRepository) Update(ctx context.Context, certificate *models.TeacherCertificate) error {
	return r.app.MySQL.WDB.WithContext(ctx).Save(certificate).Error
}

func (r *TeacherCertificateRepository) Delete(ctx context.Context, id uint) error {
	return r.app.MySQL.WDB.WithContext(ctx).Delete(&models.TeacherCertificate{}, id).Error
}

func (r *TeacherCertificateRepository) GetByID(ctx context.Context, id uint) (*models.TeacherCertificate, error) {
	var certificate models.TeacherCertificate
	err := r.app.MySQL.RDB.WithContext(ctx).First(&certificate, id).Error
	return &certificate, err
}
