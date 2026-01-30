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
	CenterID       uint           `gorm:"type:bigint unsigned;not null;index:idx_center_weekday_time" json:"center_id"`
	OfferingID     uint           `gorm:"type:bigint unsigned;not null;index" json:"offering_id"`
	TeacherID      *uint          `gorm:"type:bigint unsigned;index:idx_teacher_time" json:"teacher_id"`
	RoomID         uint           `gorm:"type:bigint unsigned;not null;index:idx_room_time" json:"room_id"`
	Name           string         `gorm:"type:varchar(100)" json:"name"`
	Weekday        int            `gorm:"type:tinyint;not null;index:idx_center_weekday_time" json:"weekday"`
	StartTime      string         `gorm:"type:varchar(10);not null;index:idx_center_weekday_time" json:"start_time"`
	EndTime        string         `gorm:"type:varchar(10);not null" json:"end_time"`
	Duration       int            `gorm:"default:60" json:"duration"`
	IsCrossDay     bool           `gorm:"type:boolean;default:false;not null" json:"is_cross_day"` // 跨日課程標記（如 23:00-02:00）
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

// MarshalJSON 自訂 JSON 序列化，輸出 MySQL 相容格式
func (dr DateRange) MarshalJSON() ([]byte, error) {
	type Alias DateRange
	return json.Marshal(&struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{
		StartDate: dr.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:   dr.EndDate.Format("2006-01-02 15:04:05"),
	})
}

// UnmarshalJSON 自訂 JSON 反序列化，支援 ISO 8601 和 MySQL 格式
func (dr *DateRange) UnmarshalJSON(data []byte) error {
	type Alias DateRange
	aux := &struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// 嘗試解析 MySQL 格式
	loc := time.UTC
	startDate, err := time.ParseInLocation("2006-01-02 15:04:05", aux.StartDate, loc)
	if err != nil {
		// 嘗試解析 ISO 8601 格式
		startDate, err = time.ParseInLocation(time.RFC3339, aux.StartDate, loc)
		if err != nil {
			return errors.New("invalid start_date format")
		}
	}
	dr.StartDate = startDate

	if aux.EndDate != "" {
		endDate, err := time.ParseInLocation("2006-01-02 15:04:05", aux.EndDate, loc)
		if err != nil {
			endDate, err = time.ParseInLocation(time.RFC3339, aux.EndDate, loc)
			if err != nil {
				return errors.New("invalid end_date format")
			}
		}
		dr.EndDate = endDate
	}

	return nil
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
