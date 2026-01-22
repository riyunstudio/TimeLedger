package models

import (
	"time"

	"gorm.io/gorm"
)

type TeacherSkill struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeacherID uint           `gorm:"type:bigint;not null;index" json:"teacher_id"`
	Category  string         `gorm:"type:varchar(100);not null" json:"category"`
	SkillName string         `gorm:"type:varchar(255);not null" json:"skill_name"`
	Level     string         `gorm:"type:varchar(20);not null" json:"level"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Hashtags []TeacherSkillHashtag `gorm:"foreignKey:TeacherSkillID" json:"hashtags,omitempty"`
}

func (TeacherSkill) TableName() string {
	return "teacher_skills"
}
