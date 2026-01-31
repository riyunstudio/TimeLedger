package models

import "time"

type TimetableCell struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TemplateID uint      `gorm:"type:bigint unsigned;not null;index" json:"template_id"`
	RowNo      int       `gorm:"type:int;not null" json:"row_no"`
	ColNo      int       `gorm:"type:int;not null" json:"col_no"`
	StartTime  string    `gorm:"type:varchar(10);not null" json:"start_time"`
	EndTime    string    `gorm:"type:varchar(10);not null" json:"end_time"`
	RoomID     *uint     `gorm:"type:bigint unsigned" json:"room_id"`
	TeacherID  *uint     `gorm:"type:bigint unsigned" json:"teacher_id"`
	IsActive   bool      `gorm:"type:boolean;default:true;not null" json:"is_active"`
	SortOrder  int       `gorm:"type:int;default:0;not null" json:"sort_order"`
	CreatedAt  time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

func (TimetableCell) TableName() string {
	return "timetable_cells"
}
