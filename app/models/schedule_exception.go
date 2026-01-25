package models

import (
	"time"
)

type ScheduleException struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CenterID     uint           `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	RuleID       uint           `gorm:"type:bigint unsigned;not null;index:idx_rule_date" json:"rule_id"`
	OriginalDate time.Time      `gorm:"type:date;not null;index:idx_rule_date" json:"original_date"`
	Type         string         `gorm:"type:varchar(20);not null" json:"type"`
	Status       string         `gorm:"type:varchar(20);default:'PENDING';not null" json:"status"`
	NewStartAt   *time.Time     `gorm:"type:datetime" json:"new_start_at"`
	NewEndAt     *time.Time     `gorm:"type:datetime" json:"new_end_at"`
	NewTeacherID *uint          `gorm:"type:bigint unsigned" json:"new_teacher_id"`
	NewRoomID    *uint          `gorm:"type:bigint unsigned" json:"new_room_id"`
	Reason       string         `gorm:"type:text" json:"reason"`
	ReviewedBy   *uint          `gorm:"type:bigint unsigned" json:"reviewed_by"`
	ReviewedAt   *time.Time     `gorm:"type:datetime" json:"reviewed_at"`
	ReviewNote   string         `gorm:"type:text" json:"review_note"`
	CreatedAt    time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:datetime;not null" json:"updated_at"`

	// 關聯
	Rule  ScheduleRule `gorm:"foreignKey:RuleID" json:"rule,omitempty"`
	Teacher Teacher    `gorm:"foreignKey:TeacherID;references:ID;AssociationForeignKey:Rule.TeacherID" json:"teacher,omitempty"`
}

func (ScheduleException) TableName() string {
	return "schedule_exceptions"
}
