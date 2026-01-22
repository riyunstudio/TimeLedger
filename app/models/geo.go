package models

type GeoCity struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Name      string        `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Districts []GeoDistrict `gorm:"foreignKey:CityID" json:"districts,omitempty"`
}

func (GeoCity) TableName() string {
	return "geo_cities"
}

type GeoDistrict struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	CityID uint   `gorm:"type:bigint;not null;index" json:"city_id"`
	Name   string `gorm:"type:varchar(100);not null" json:"name"`
}

func (GeoDistrict) TableName() string {
	return "geo_districts"
}
