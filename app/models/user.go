package models

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"type:varchar(30); default:''; not null; comment:'使用者名稱'" json:"name"`
	Ips        string `gorm:"type:json; not null; comment:'白名單IP'" json:"ips"`
	CreateTime int64  `gorm:"type:int(10); default:0; not null; comment:'建立時間'" json:"create_time"`
	UpdateTime int64  `gorm:"type:int(10); default:0; not null; comment:'更新時間'" json:"update_time"`
}
