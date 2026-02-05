package models

import (
	"fmt"
	"time"
)

type ScheduleException struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CenterID      uint       `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	RuleID        uint       `gorm:"type:bigint unsigned;not null;index:idx_rule_date" json:"rule_id"`
	OriginalDate  time.Time  `gorm:"type:date;not null;index:idx_rule_date" json:"original_date"`
	ExceptionType string     `gorm:"column:exception_type;type:varchar(20);not null" json:"exception_type"` // LEAVE, RESCHEDULE, SWAP, CANCEL
	Status        string     `gorm:"type:varchar(20);default:'PENDING';not null" json:"status"`
	NewStartAt    *time.Time `gorm:"type:datetime" json:"new_start_at"`
	NewEndAt      *time.Time `gorm:"type:datetime" json:"new_end_at"`
	NewTeacherID  *uint      `gorm:"type:bigint unsigned" json:"new_teacher_id"`
	NewRoomID     *uint      `gorm:"type:bigint unsigned" json:"new_room_id"`
	Reason        string     `gorm:"type:text" json:"reason"`
	ReviewedBy    *uint      `gorm:"type:bigint unsigned" json:"reviewed_by"`
	ReviewedAt    *time.Time `gorm:"type:datetime" json:"reviewed_at"`
	ReviewNote    string     `gorm:"type:text" json:"review_note"`
	CreatedAt     time.Time  `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:datetime;not null" json:"updated_at"`

	// 關聯
	Rule ScheduleRule `gorm:"foreignKey:RuleID" json:"rule,omitempty"`
}

// GetDate 取得日期（用於顯示）
func (s *ScheduleException) GetDate() time.Time {
	return s.OriginalDate
}

// GetTimeRange 取得時間範圍字串
func (s *ScheduleException) GetTimeRange() string {
	// 使用 NewStartAt 和 NewEndAt，如果沒有則使用 Rule 的時間
	if s.NewStartAt != nil && s.NewEndAt != nil {
		return fmt.Sprintf("%s - %s",
			s.NewStartAt.Format("15:04"),
			s.NewEndAt.Format("15:04"),
		)
	}

	// fallback 到 Rule 的時間（使用 StartTime 和 EndTime 欄位）
	if s.Rule.StartTime != "" && s.Rule.EndTime != "" {
		return fmt.Sprintf("%s - %s", s.Rule.StartTime, s.Rule.EndTime)
	}

	return "時間未定"
}

func (ScheduleException) TableName() string {
	return "schedule_exceptions"
}
