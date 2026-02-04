package models

import (
	"time"

	"gorm.io/gorm"
)

// AdminLoginHistory 管理員登入紀錄
type AdminLoginHistory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	AdminID   uint           `gorm:"type:bigint unsigned;not null;index" json:"admin_id"`
	Email     string         `gorm:"type:varchar(255);not null;index" json:"email"`
	Status    string         `gorm:"type:varchar(20);not null;index" json:"status"` // SUCCESS, FAILED
	IPAddress string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string         `gorm:"type:varchar(500)" json:"user_agent"`
	Reason    string         `gorm:"type:varchar(255)" json:"reason"` // 失敗原因
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	LoginStatusSuccess = "SUCCESS"
	LoginStatusFailed  = "FAILED"
)

func (AdminLoginHistory) TableName() string {
	return "admin_login_histories"
}
