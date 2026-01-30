package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type HashtagRepository struct {
	GenericRepository[models.Hashtag]
	app *app.App
}

func NewHashtagRepository(app *app.App) *HashtagRepository {
	return &HashtagRepository{
		GenericRepository: NewGenericRepository[models.Hashtag](app.MySQL.RDB, app.MySQL.WDB),
		app:               app,
	}
}

func (r *HashtagRepository) GetByName(ctx context.Context, name string) (*models.Hashtag, error) {
	data, err := r.First(ctx, "name = ?", name)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *HashtagRepository) Search(ctx context.Context, query string) ([]models.Hashtag, error) {
	return r.Find(ctx, "name LIKE ?", "%"+query+"%")
}

func (r *HashtagRepository) IncrementUsage(ctx context.Context, name string) error {
	return r.app.MySQL.WDB.WithContext(ctx).Model(&models.Hashtag{}).Where("name = ?", name).UpdateColumn("usage_count", "usage_count + 1").Error
}
