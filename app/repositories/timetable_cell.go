package repositories

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// ReorderCellRequest 重新排序格子請求
type ReorderCellRequest struct {
	ID        uint `json:"id" binding:"required"`
	SortOrder int  `json:"sort_order" binding:"required"`
}

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

// BatchUpdateSortOrder 批次更新格子的排序順序
func (rp *TimetableCellRepository) BatchUpdateSortOrder(ctx context.Context, cells []ReorderCellRequest) error {
	for _, cell := range cells {
		updates := map[string]interface{}{
			"sort_order": cell.SortOrder,
			"updated_at": time.Now(),
		}
		if err := rp.UpdateFields(ctx, cell.ID, updates); err != nil {
			return err
		}
	}
	return nil
}
