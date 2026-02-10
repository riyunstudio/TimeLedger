package models

import (
	"time"

	"gorm.io/gorm"
)

type Teacher struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	LineUserID        string         `gorm:"type:varchar(255);index" json:"line_user_id"`
	LineNotifyToken   string         `gorm:"type:varchar(255)" json:"line_notify_token"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	Email             string         `gorm:"type:varchar(255)" json:"email"`
	AvatarURL         string         `gorm:"type:varchar(512)" json:"avatar_url"`
	Bio               string         `gorm:"type:text" json:"bio"`
	IsOpenToHiring    bool           `gorm:"type:boolean;default:false;not null" json:"is_open_to_hiring"`
	IsPlaceholder     bool           `gorm:"type:boolean;default:false;not null" json:"is_placeholder"`
	City              string         `gorm:"type:varchar(100);index" json:"city"`
	District          string         `gorm:"type:varchar(100)" json:"district"`
	PublicContactInfo string         `gorm:"type:text" json:"public_contact_info"`
	CreatedAt         time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	Skills           []TeacherSkill           `gorm:"foreignKey:TeacherID" json:"skills,omitempty"`
	Certificates     []TeacherCertificate     `gorm:"foreignKey:TeacherID" json:"certificates,omitempty"`
	PersonalHashtags []TeacherPersonalHashtag `gorm:"foreignKey:TeacherID" json:"personal_hashtags,omitempty"`
}

func (Teacher) TableName() string {
	return "teachers"
}
