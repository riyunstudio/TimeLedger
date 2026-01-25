package models

import (
	"time"

	"gorm.io/gorm"
)

type CenterInvitation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CenterID  uint           `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	Email     string         `gorm:"type:varchar(255)" json:"email"`
	Token     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	Status    string         `gorm:"type:varchar(20);default:'PENDING';not null" json:"status"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	ExpiresAt time.Time      `gorm:"type:datetime;not null;index" json:"expires_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CenterInvitation) TableName() string {
	return "center_invitations"
}
