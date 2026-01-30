package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TimetableTemplateRepository struct {
	GenericRepository[models.TimetableTemplate]
	app *app.App
}

func NewTimetableTemplateRepository(app *app.App) *TimetableTemplateRepository {
	return &TimetableTemplateRepository{
		GenericRepository: NewGenericRepository[models.TimetableTemplate](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *TimetableTemplateRepository) GetByIDAndCenterID(ctx context.Context, id uint, centerID uint) (models.TimetableTemplate, error) {
	return rp.GetByIDWithCenterScope(ctx, id, centerID)
}

func (rp *TimetableTemplateRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.TimetableTemplate, error) {
	return rp.FindWithCenterScope(ctx, centerID)
}

func (rp *TimetableTemplateRepository) UpdateByIDAndCenterID(ctx context.Context, id uint, centerID uint, fields map[string]interface{}) error {
	return rp.UpdateFieldsWithCenterScope(ctx, id, centerID, fields)
}

func (rp *TimetableTemplateRepository) DeleteByIDAndCenterID(ctx context.Context, id uint, centerID uint) error {
	return rp.DeleteByIDWithCenterScope(ctx, id, centerID)
}
