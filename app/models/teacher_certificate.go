package models

import (
	"time"

	"gorm.io/gorm"
)

type TeacherCertificate struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TeacherID  uint           `gorm:"type:bigint unsigned;not null;index" json:"teacher_id"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	FileURL    string         `gorm:"type:varchar(512);not null" json:"file_url"`
	IssuedAt   time.Time      `gorm:"type:date;not null" json:"issued_at"`
	IsVerified bool           `gorm:"type:boolean;default:false;not null" json:"is_verified"`
	CreatedAt  time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TeacherCertificate) TableName() string {
	return "teacher_certificates"
}
