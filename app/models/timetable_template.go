package models

import (
	"time"

	"gorm.io/gorm"
)

type TimetableTemplate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CenterID  uint           `gorm:"type:bigint;not null;index" json:"center_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	RowType   string         `gorm:"type:varchar(20);not null" json:"row_type"`
	IsActive  bool           `gorm:"type:boolean;default:true;not null" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Cells []TimetableCell `gorm:"foreignKey:TemplateID" json:"cells,omitempty"`
}

func (TimetableTemplate) TableName() string {
	return "timetable_templates"
}
