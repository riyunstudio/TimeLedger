package models

import (
	"time"
)

type SessionNote struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TeacherID       uint      `gorm:"type:bigint unsigned;not null;index" json:"teacher_id"`
	RuleID          *uint     `gorm:"type:bigint unsigned" json:"rule_id"`
	PersonalEventID *uint     `gorm:"type:bigint unsigned" json:"personal_event_id"`
	SessionDate     time.Time `gorm:"type:date;not null" json:"session_date"`
	Content         string    `gorm:"type:text" json:"content"`
	PrepNote        string    `gorm:"type:text" json:"prep_note"`
	UpdatedAt       time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

func (SessionNote) TableName() string {
	return "session_notes"
}
