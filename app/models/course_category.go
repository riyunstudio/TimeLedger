package models

// CourseCategory 課程類別
type CourseCategory struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	CenterID uint   `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	Name     string `gorm:"type:varchar(50);not null" json:"name"`
}

func (CourseCategory) TableName() string {
	return "course_categories"
}

// CourseCategoryResponse 課程類別響應
type CourseCategoryResponse struct {
	ID       uint   `json:"id"`
	CenterID uint   `json:"center_id"`
	Name     string `json:"name"`
}
