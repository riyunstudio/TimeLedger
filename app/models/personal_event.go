package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PersonalEvent struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TeacherID      uint           `gorm:"type:bigint unsigned;not null;index:idx_teacher_time" json:"teacher_id"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	StartAt        time.Time      `gorm:"type:datetime;not null;index:idx_teacher_time" json:"start_at"`
	EndAt          time.Time      `gorm:"type:datetime;not null" json:"end_at"`
	RecurrenceRule RecurrenceRule `gorm:"type:json" json:"recurrence_rule"`
	IsAllDay       bool           `gorm:"type:boolean;default:false;not null" json:"is_all_day"`
	ColorHex       string         `gorm:"type:varchar(7)" json:"color_hex"`
	Note           string         `gorm:"type:text" json:"note"`
	CreatedAt      time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type RecurrenceRule struct {
	Type     string  `json:"type"`
	Interval int     `json:"interval"`
	Weekdays []int   `json:"weekdays,omitempty"`
	Until    *string `json:"until,omitempty"`
	Count    *int    `json:"count,omitempty"`
}

func (rr *RecurrenceRule) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal RecurrenceRule value")
	}
	return json.Unmarshal(bytes, rr)
}

func (rr RecurrenceRule) Value() (driver.Value, error) {
	return json.Marshal(rr)
}

func (PersonalEvent) TableName() string {
	return "personal_events"
}
