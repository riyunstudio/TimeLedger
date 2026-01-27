package models

import (
	"time"
)

type NotificationQueue struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Type        string         `gorm:"type:varchar(32);not null;index" json:"type"`                     // exception_submit, exception_result, welcome, etc.
	RecipientID uint           `gorm:"not null;index" json:"recipient_id"`                              // admin_id 或 teacher_id
	RecipientType string       `gorm:"type:varchar(20);not null" json:"recipient_type"`                 // ADMIN, TEACHER
	Payload     string         `gorm:"type:text;not null" json:"-"`                                    // JSON 訊息內容
	Status      string         `gorm:"type:varchar(16);default:'pending';not null;index" json:"status"` // pending, sent, failed
	RetryCount  int            `gorm:"default:0" json:"retry_count"`
	ErrorMsg    string         `gorm:"type:text" json:"-"`
	ScheduledAt time.Time      `gorm:"type:datetime;not null" json:"scheduled_at"`                     // 預定發送時間
	SentAt      *time.Time     `json:"sent_at"`
	FailedAt    *time.Time     `json:"failed_at"`
	CreatedAt   time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
}

func (NotificationQueue) TableName() string {
	return "notification_queues"
}

// NotificationStatus constants
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

// NotificationType constants
const (
	NotificationTypeExceptionSubmit = "exception_submit" // 老師提交例外申請
	NotificationTypeExceptionResult = "exception_result" // 例外審核結果
	NotificationTypeWelcomeTeacher  = "welcome_teacher"  // 老師歡迎訊息
	NotificationTypeWelcomeAdmin    = "welcome_admin"    // 管理員歡迎訊息
)
