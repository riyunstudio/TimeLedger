package models

import (
	"time"

	"gorm.io/gorm"
)

type CenterTerm struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CenterID  uint           `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	StartDate time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time      `gorm:"type:date;not null" json:"end_date"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CenterTerm) TableName() string {
	return "center_terms"
}
