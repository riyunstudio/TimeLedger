package models

import (
	"time"

	"gorm.io/gorm"
)

type Offering struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	CenterID            uint           `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	CourseID            uint           `gorm:"type:bigint unsigned;not null;index" json:"course_id"`
	Name                string         `gorm:"type:varchar(255);not null" json:"name"`
	DefaultRoomID       *uint          `gorm:"type:bigint unsigned;index" json:"default_room_id"`
	DefaultTeacherID    *uint          `gorm:"type:bigint unsigned;index" json:"default_teacher_id"`
	DefaultStartTime    string         `gorm:"type:varchar(5);not null;default:'09:00'" json:"default_start_time"` // 預設開始時間 (HH:MM)
	DefaultEndTime      string         `gorm:"type:varchar(5);not null;default:'10:00'" json:"default_end_time"`   // 預設結束時間 (HH:MM)
	AllowBufferOverride bool           `gorm:"type:boolean;default:false;not null" json:"allow_buffer_override"`
	IsActive            bool           `gorm:"type:boolean;default:true;not null" json:"is_active"`
	CreatedAt           time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯課程（含課程時長）
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (Offering) TableName() string {
	return "offerings"
}
