package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"timeLedger/libs"
)

type ScheduleRule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CenterID       uint           `gorm:"type:bigint unsigned;not null;index:idx_center_weekday_time" json:"center_id"`
	OfferingID     uint           `gorm:"type:bigint unsigned;not null;index" json:"offering_id"`
	TeacherID      *uint          `gorm:"type:bigint unsigned;index:idx_teacher_time" json:"teacher_id"`
	RoomID         uint           `gorm:"type:bigint unsigned;not null;index:idx_room_time" json:"room_id"`
	Name           string         `gorm:"type:varchar(100)" json:"name"`
	Code           string         `gorm:"type:varchar(50)" json:"code"`
	Weekday        int            `gorm:"type:tinyint;not null;index:idx_center_weekday_time" json:"weekday"`
	StartTime      string         `gorm:"type:varchar(10);not null;index:idx_center_weekday_time" json:"start_time"`
	EndTime        string         `gorm:"type:varchar(10);not null" json:"end_time"`
	Duration       int            `gorm:"default:60" json:"duration"`
	IsCrossDay     bool           `gorm:"type:boolean;default:false;not null" json:"is_cross_day"` // 跨日課程標記（如 23:00-02:00）
	SkipHoliday    bool           `gorm:"type:boolean;default:true;not null" json:"skip_holiday"`
	EffectiveRange DateRange      `gorm:"type:json;not null" json:"effective_range"`
	SuspendedDates SuspendedDates  `gorm:"type:json" json:"suspended_dates"`
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

// SuspendedDates 自訂類型，用於 suspension_dates 欄位的 JSON 序列化
type SuspendedDates []time.Time

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

// UnmarshalJSON 自訂 JSON 反序列化，支援多種日期格式
func (dr *DateRange) UnmarshalJSON(data []byte) error {
	type Alias DateRange
	aux := &struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// 使用台灣時區解析
	loc := libs.GetTaiwanLocation()

	// 格式1: YYYY-MM-DD HH:MM:SS+08:00 (帶時區的 ISO 8601)
	startDate, err := time.ParseInLocation("2006-01-02T15:04:05Z07:00", aux.StartDate, loc)
	if err == nil {
		dr.StartDate = startDate
	} else {
		// 格式2: YYYY-MM-DD HH:MM:SS (MySQL datetime)
		startDate, err = time.ParseInLocation("2006-01-02 15:04:05", aux.StartDate, loc)
		if err == nil {
			dr.StartDate = startDate
		} else {
			// 格式3: YYYY-MM-DD (日期 only)
			startDate, err = time.ParseInLocation("2006-01-02", aux.StartDate, loc)
			if err == nil {
				dr.StartDate = startDate
			} else {
				return errors.New("invalid start_date format: " + aux.StartDate)
			}
		}
	}

	if aux.EndDate != "" {
		endDate, err := time.ParseInLocation("2006-01-02T15:04:05Z07:00", aux.EndDate, loc)
		if err == nil {
			dr.EndDate = endDate
		} else {
			endDate, err = time.ParseInLocation("2006-01-02 15:04:05", aux.EndDate, loc)
			if err == nil {
				dr.EndDate = endDate
			} else {
				endDate, err = time.ParseInLocation("2006-01-02", aux.EndDate, loc)
				if err == nil {
					dr.EndDate = endDate
				} else {
					return errors.New("invalid end_date format: " + aux.EndDate)
				}
			}
		}
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

// SuspendedDates 自訂 JSON 序列化，將 time.Time 陣列序列化为 MySQL DATETIME 格式字串陣列
func (sd SuspendedDates) MarshalJSON() ([]byte, error) {
	dateStrings := make([]string, len(sd))
	for i, t := range sd {
		dateStrings[i] = t.Format("2006-01-02 15:04:05")
	}
	return json.Marshal(dateStrings)
}

// SuspendedDates 自訂 JSON 反序列化，支援多種日期格式
func (sd *SuspendedDates) UnmarshalJSON(data []byte) error {
	var dateStrings []string
	if err := json.Unmarshal(data, &dateStrings); err != nil {
		return err
	}

	if dateStrings == nil {
		*sd = nil
		return nil
	}

	// 使用台灣時區解析
	loc := libs.GetTaiwanLocation()
	dates := make([]time.Time, 0, len(dateStrings))

	for _, dateStr := range dateStrings {
		var parsedDate time.Time
		var err error

		// 格式1: YYYY-MM-DD HH:MM:SS+08:00 (帶時區的 ISO 8601)
		parsedDate, err = time.ParseInLocation("2006-01-02T15:04:05Z07:00", dateStr, loc)
		if err == nil {
			dates = append(dates, parsedDate)
			continue
		}

		// 格式2: YYYY-MM-DD HH:MM:SS (MySQL datetime)
		parsedDate, err = time.ParseInLocation("2006-01-02 15:04:05", dateStr, loc)
		if err == nil {
			dates = append(dates, parsedDate)
			continue
		}

		// 格式3: YYYY-MM-DD (日期 only)
		parsedDate, err = time.ParseInLocation("2006-01-02", dateStr, loc)
		if err == nil {
			dates = append(dates, parsedDate)
			continue
		}

		// 格式4: ISO 8601 標準格式
		parsedDate, err = time.ParseInLocation(time.RFC3339, dateStr, loc)
		if err == nil {
			dates = append(dates, parsedDate)
			continue
		}

		return errors.New("invalid date format in suspended_dates: " + dateStr)
	}

	*sd = dates
	return nil
}

// SuspendedDates Scan 方法，支援從資料庫讀取 JSON
func (sd *SuspendedDates) Scan(value interface{}) error {
	if value == nil {
		*sd = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal SuspendedDates value")
	}

	return json.Unmarshal(bytes, sd)
}

// SuspendedDates Value 方法，支援寫入資料庫
func (sd SuspendedDates) Value() (driver.Value, error) {
	if sd == nil {
		return nil, nil
	}
	return json.Marshal(sd)
}

func (ScheduleRule) TableName() string {
	return "schedule_rules"
}
