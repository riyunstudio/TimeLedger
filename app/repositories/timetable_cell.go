package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type TimetableCellRepository struct {
	GenericRepository[models.TimetableCell]
	app *app.App
}

func NewTimetableCellRepository(app *app.App) *TimetableCellRepository {
	return &TimetableCellRepository{
		GenericRepository: NewGenericRepository[models.TimetableCell](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *TimetableCellRepository) ListByTemplateID(ctx context.Context, templateID uint) ([]models.TimetableCell, error) {
	return rp.Find(ctx, "template_id = ?", templateID)
}

func (rp *TimetableCellRepository) DeleteByTemplateID(ctx context.Context, templateID uint) error {
	_, err := rp.DeleteWhere(ctx, "template_id = ?", templateID)
	return err
}

func (rp *TimetableCellRepository) Delete(ctx context.Context, id uint) error {
	return rp.DeleteByID(ctx, id)
}
