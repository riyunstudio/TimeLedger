package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ScheduleRule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CenterID       uint           `gorm:"not null;index:idx_center_weekday_time" json:"center_id"`
	OfferingID     uint           `gorm:"not null;index" json:"offering_id"`
	TeacherID      *uint          `gorm:"index:idx_teacher_time" json:"teacher_id"`
	RoomID         uint           `gorm:"not null;index:idx_room_time" json:"room_id"`
	Weekday        int            `gorm:"type:tinyint;not null;index:idx_center_weekday_time" json:"weekday"`
	StartTime      string         `gorm:"type:varchar(10);not null;index:idx_center_weekday_time" json:"start_time"`
	EndTime        string         `gorm:"type:varchar(10);not null" json:"end_time"`
	EffectiveRange DateRange      `gorm:"type:json;not null" json:"effective_range"`
	LockAt         *time.Time     `gorm:"type:datetime;index" json:"lock_at"`
	CreatedAt      time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 關聯
	Offering Offering `gorm:"foreignKey:OfferingID" json:"offering,omitempty"`
	Teacher  Teacher  `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Room     Room     `gorm:"foreignKey:RoomID" json:"room,omitempty"`

	Exceptions []ScheduleException `gorm:"foreignKey:RuleID" json:"exceptions,omitempty"`
}

type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (dr *DateRange) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal DateRange value")
	}
	return json.Unmarshal(str, dr)
}

func (dr DateRange) Value() (driver.Value, error) {
	return json.Marshal(dr)
}

func (ScheduleRule) TableName() string {
	return "schedule_rules"
}
