package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type RoomRepository struct {
	GenericRepository[models.Room]
	app *app.App
}

func NewRoomRepository(app *app.App) *RoomRepository {
	return &RoomRepository{
		GenericRepository: NewGenericRepository[models.Room](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (rp *RoomRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.Room, error) {
	return rp.FindWithCenterScope(ctx, centerID, "is_active = ?", true)
}

func (rp *RoomRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.Room, error) {
	return rp.FindWithCenterScope(ctx, centerID)
}

func (rp *RoomRepository) ToggleActive(ctx context.Context, id uint, centerID uint, isActive bool) error {
	return rp.UpdateFieldsWithCenterScope(ctx, id, centerID, map[string]interface{}{
		"is_active": isActive,
	})
}
