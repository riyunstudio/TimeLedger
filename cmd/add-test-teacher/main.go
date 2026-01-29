package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連線資料庫失敗: %v", err)
	}

	// 確保資料表結構是最新的
	log.Println("正在同步資料表結構...")
	if err := db.AutoMigrate(&Teacher{}); err != nil {
		log.Printf("警告: AutoMigrate 失敗: %v", err)
	}

	// 新增測試老師（用於 SearchTalent_WithFilters 測試）
	teacher := Teacher{
		Name:           "測試老師 - SearchTalent - TestData",
		Email:          "test.searchtalent.20260129@test.com",
		Phone:          "0912345678",
		City:           "新北市",
		District:       "板橋區",
		Bio:            "這是為了測試搜尋功能新增的老師",
		IsActive:       true,
		IsOpenToHiring: true,
		LineUserID:     fmt.Sprintf("line_test_searchtalent_%d", time.Now().Unix()),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result := db.Table("teachers").Create(&teacher)
	if result.Error != nil {
		log.Fatalf("新增老師失敗: %v", result.Error)
	}

	log.Printf("✅ 已新增測試老師: %s (ID: %d)", teacher.Name, teacher.ID)
}

// Teacher 模型（符合開發資料庫結構）
type Teacher struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	LineUserID        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"line_user_id"`
	Name              string    `gorm:"type:varchar(255);not null" json:"name"`
	Email             string    `gorm:"type:varchar(255)" json:"email"`
	Phone             string    `gorm:"type:varchar(20)" json:"phone"`
	Bio               string    `gorm:"type:text" json:"bio"`
	IsOpenToHiring    bool      `gorm:"type:boolean;default:false;not null" json:"is_open_to_hiring"`
	City              string    `gorm:"type:varchar(100);index" json:"city"`
	District          string    `gorm:"type:varchar(100)" json:"district"`
	PublicContactInfo string    `gorm:"type:text" json:"public_contact_info"`
	IsActive          bool      `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt         time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

func (Teacher) TableName() string {
	return "teachers"
}
