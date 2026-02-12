package models

import (
	"time"

	"gorm.io/gorm"
)

type TeacherBackground struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeacherID uint           `gorm:"type:bigint unsigned;not null;index:idx_teacher_background" json:"teacher_id"`
	FileURL   string         `gorm:"type:varchar(512);not null" json:"file_url"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TeacherBackground) TableName() string {
	return "teacher_backgrounds"
}
