package models

import (
	"time"

	"gorm.io/gorm"
)

type AdminUser struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CenterID          uint           `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	Email             string         `gorm:"type:varchar(255);not null" json:"email"`
	PasswordHash      string         `gorm:"type:varchar(255);not null" json:"-"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	Role              string         `gorm:"type:varchar(20);default:'STAFF';not null" json:"role"`
	Status            string         `gorm:"type:varchar(20);default:'ACTIVE';not null" json:"status"`
	// LINE 綁定相關欄位
	LineUserID        string         `gorm:"type:varchar(64);index" json:"-"`
	LineBindingCode   string         `gorm:"type:varchar(8)" json:"-"`
	LineBindingExpires *time.Time    `json:"-"`
	LineNotifyEnabled bool           `gorm:"default:true" json:"line_notify_enabled"`
	LineBoundAt       *time.Time     `json:"line_bound_at"`
	CreatedAt         time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}
