package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Center struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	PlanLevel string         `gorm:"type:varchar(20);default:'FREE';not null" json:"plan_level"`
	Settings  CenterSettings `gorm:"type:json;not null" json:"settings"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;autoCreateTime" json:"created_at"`
}

type CenterSettings struct {
	AllowPublicRegister      bool   `json:"allow_public_register"`
	DefaultLanguage         string `json:"default_language"`
	ExceptionLeadDays       int    `json:"exception_lead_days"`
	DefaultCourseDuration   int    `json:"default_course_duration"`
	OperatingStartTime      string `json:"operating_start_time"`
	OperatingEndTime        string `json:"operating_end_time"`
}

func (cs *CenterSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal CenterSettings value")
	}
	return json.Unmarshal(bytes, cs)
}

func (cs CenterSettings) Value() (driver.Value, error) {
	return json.Marshal(cs)
}

func (Center) TableName() string {
	return "centers"
}
