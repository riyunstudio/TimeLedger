package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type HashtagRepository struct {
	BaseRepository
	app *app.App
}

func NewHashtagRepository(app *app.App) *HashtagRepository {
	return &HashtagRepository{app: app}
}

func (r *HashtagRepository) GetByID(ctx context.Context, id uint) (*models.Hashtag, error) {
	var hashtag models.Hashtag
	err := r.app.Mysql.RDB.WithContext(ctx).First(&hashtag, id).Error
	return &hashtag, err
}

func (r *HashtagRepository) GetByName(ctx context.Context, name string) (*models.Hashtag, error) {
	var hashtag models.Hashtag
	err := r.app.Mysql.RDB.WithContext(ctx).Where("name = ?", name).First(&hashtag).Error
	return &hashtag, err
}

func (r *HashtagRepository) Create(ctx context.Context, hashtag *models.Hashtag) error {
	return r.app.Mysql.WDB.WithContext(ctx).Create(hashtag).Error
}

func (r *HashtagRepository) List(ctx context.Context) ([]models.Hashtag, error) {
	var hashtags []models.Hashtag
	err := r.app.Mysql.RDB.WithContext(ctx).Find(&hashtags).Error
	return hashtags, err
}

func (r *HashtagRepository) IncrementUsage(ctx context.Context, name string) error {
	return r.app.Mysql.WDB.WithContext(ctx).Model(&models.Hashtag{}).Where("name = ?", name).UpdateColumn("usage_count", "usage_count + 1").Error
}
