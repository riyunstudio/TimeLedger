package repositories

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
)

type GeoRepository struct {
	BaseRepository
	app *app.App
}

func NewGeoRepository(app *app.App) *GeoRepository {
	return &GeoRepository{app: app}
}

func (r *GeoRepository) ListCities(ctx context.Context) ([]models.GeoCity, error) {
	var cities []models.GeoCity
	err := r.app.MySQL.RDB.WithContext(ctx).Preload("Districts").Find(&cities).Error
	return cities, err
}

func (r *GeoRepository) GetCityByName(ctx context.Context, name string) (*models.GeoCity, error) {
	var city models.GeoCity
	err := r.app.MySQL.RDB.WithContext(ctx).Preload("Districts").Where("name = ?", name).First(&city).Error
	return &city, err
}
