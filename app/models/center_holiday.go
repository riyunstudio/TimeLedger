package models

import (
	"time"

	"gorm.io/gorm"
)

type CenterHoliday struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CenterID  uint           `gorm:"type:bigint unsigned;not null;index:idx_center_date" json:"center_id"`
	Date      time.Time      `gorm:"type:date;not null;index:idx_center_date" json:"date"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CenterHoliday) TableName() string {
	return "center_holidays"
}
