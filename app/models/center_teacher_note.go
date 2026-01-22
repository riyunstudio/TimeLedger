package models

import (
	"time"
)

type CenterTeacherNote struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CenterID     uint      `gorm:"type:bigint;not null;index" json:"center_id"`
	TeacherID    uint      `gorm:"type:bigint;not null" json:"teacher_id"`
	InternalNote string    `gorm:"type:text" json:"internal_note"`
	Rating       int       `gorm:"type:tinyint" json:"rating"`
	CreatedAt    time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

func (CenterTeacherNote) TableName() string {
	return "center_teacher_notes"
}
