package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
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

// Transaction executes a function within a database transaction.
// This method creates a NEW RoomRepository instance with transaction connections.
func (rp *RoomRepository) Transaction(ctx context.Context, fn func(txRepo *RoomRepository) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &RoomRepository{
			GenericRepository: NewTransactionRepo[models.Room](ctx, tx, rp.table),
			app:               rp.app,
		}
		return fn(txRepo)
	})
}

func (rp *RoomRepository) ListActiveByCenterID(ctx context.Context, centerID uint) ([]models.Room, error) {
	return rp.FindWithCenterScope(ctx, centerID, "is_active = ?", true)
}

func (rp *RoomRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.Room, error) {
	return rp.FindWithCenterScope(ctx, centerID)
}

// SearchByNamePaginated 搜尋教室（分頁）
func (rp *RoomRepository) SearchByNamePaginated(ctx context.Context, centerID uint, query string, page, limit int) ([]models.Room, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// 計算總數
	var total int64
	countQuery := rp.dbRead.Table(rp.table).
		Where("center_id = ?", centerID)
	if query != "" {
		countQuery = countQuery.Where("name LIKE ?", "%"+query+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查詢資料
	var rooms []models.Room
	dataQuery := rp.dbRead.Table(rp.table).
		Where("center_id = ?", centerID)
	if query != "" {
		dataQuery = dataQuery.Where("name LIKE ?", "%"+query+"%")
	}
	if err := dataQuery.Order("id DESC").Offset(offset).Limit(limit).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (rp *RoomRepository) ToggleActive(ctx context.Context, id uint, centerID uint, isActive bool) error {
	return rp.UpdateFieldsWithCenterScope(ctx, id, centerID, map[string]interface{}{
		"is_active": isActive,
	})
}
