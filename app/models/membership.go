package models

import (
	"time"

	"gorm.io/gorm"
)

type CenterMembership struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CenterID  uint           `gorm:"type:bigint;not null;index" json:"center_id"`
	TeacherID uint           `gorm:"type:bigint;not null;index" json:"teacher_id"`
	Status    string         `gorm:"type:varchar(20);default:'INVITED';not null" json:"status"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CenterMembership) TableName() string {
	return "center_memberships"
}
