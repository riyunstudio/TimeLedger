package models

type Hashtag struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	UsageCount int    `gorm:"type:int;default:0;not null" json:"usage_count"`
}

func (Hashtag) TableName() string {
	return "hashtags"
}

type TeacherSkillHashtag struct {
	TeacherSkillID uint    `gorm:"type:bigint unsigned;not null;index" json:"teacher_skill_id"`
	HashtagID      uint    `gorm:"type:bigint unsigned;not null;index" json:"hashtag_id"`
	Hashtag        Hashtag `gorm:"foreignKey:HashtagID" json:"hashtag,omitempty"`
}

func (TeacherSkillHashtag) TableName() string {
	return "teacher_skill_hashtags"
}
