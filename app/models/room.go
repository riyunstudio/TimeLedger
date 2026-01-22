package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CenterID  uint           `gorm:"type:bigint;not null;index" json:"center_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Capacity  int            `gorm:"type:int;not null;default:1" json:"capacity"`
	IsActive  bool           `gorm:"type:boolean;default:true;not null" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Room) TableName() string {
	return "rooms"
}
