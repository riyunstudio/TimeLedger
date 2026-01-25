package models

type TeacherPersonalHashtag struct {
	TeacherID uint    `gorm:"type:bigint unsigned;not null;index" json:"teacher_id"`
	HashtagID uint    `gorm:"type:bigint unsigned;not null;index" json:"hashtag_id"`
	SortOrder int     `gorm:"type:tinyint;not null" json:"sort_order"`
	Hashtag   Hashtag `gorm:"foreignKey:HashtagID" json:"hashtag,omitempty"`
}

func (TeacherPersonalHashtag) TableName() string {
	return "teacher_personal_hashtags"
}
