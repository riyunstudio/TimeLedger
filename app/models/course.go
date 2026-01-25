package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CenterID         uint           `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	DefaultDuration  int            `gorm:"type:int;not null;default:60" json:"default_duration"`
	ColorHex         string         `gorm:"type:varchar(7);not null" json:"color_hex"`
	RoomBufferMin    int            `gorm:"type:int;not null;default:0" json:"room_buffer_min"`
	TeacherBufferMin int            `gorm:"type:int;not null;default:0" json:"teacher_buffer_min"`
	IsActive         bool           `gorm:"type:boolean;default:true;not null" json:"is_active"`
	CreatedAt        time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Course) TableName() string {
	return "courses"
}
