package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"type:bigint;not null;index" json:"user_id"`
	UserType  string         `gorm:"type:varchar(20);not null" json:"user_type"`
	CenterID  uint           `gorm:"type:bigint;index" json:"center_id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Type      string         `gorm:"type:varchar(20);not null" json:"type"`
	IsRead    bool           `gorm:"type:boolean;default:false;not null" json:"is_read"`
	ReadAt    *time.Time     `gorm:"type:datetime" json:"read_at"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;index" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}
