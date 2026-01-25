package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CenterID   uint           `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	ActorType  string         `gorm:"type:varchar(20);not null" json:"actor_type"`
	ActorID    uint           `gorm:"type:bigint unsigned;not null" json:"actor_id"`
	Action     string         `gorm:"type:varchar(255);not null" json:"action"`
	TargetType string         `gorm:"type:varchar(255);not null" json:"target_type"`
	TargetID   uint           `gorm:"type:bigint unsigned;not null" json:"target_id"`
	Payload    AuditPayload   `gorm:"type:json" json:"payload"`
	Timestamp  time.Time      `gorm:"type:datetime;not null" json:"timestamp"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type AuditPayload struct {
	Before interface{} `json:"before,omitempty"`
	After  interface{} `json:"after,omitempty"`
}

func (ap *AuditPayload) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal AuditPayload value")
	}
	return json.Unmarshal(bytes, ap)
}

func (ap AuditPayload) Value() (driver.Value, error) {
	return json.Marshal(ap)
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
