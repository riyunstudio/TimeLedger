package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連線資料庫失敗: %v", err)
	}

	// 查詢陳大文的資料
	var teacher struct {
		ID             uint
		Name           string
		Email          string
		City           string
		District       string
		IsOpenToHiring bool
		IsActive       bool
	}

	db.Table("teachers").Where("id = 2").First(&teacher)

	log.Printf("陳大文資料:")
	log.Printf("  ID: %d", teacher.ID)
	log.Printf("  Name: %s", teacher.Name)
	log.Printf("  Email: %s", teacher.Email)
	log.Printf("  City: %s", teacher.City)
	log.Printf("  District: %s", teacher.District)
	log.Printf("  IsOpenToHiring: %v", teacher.IsOpenToHiring)
	log.Printf("  IsActive: %v", teacher.IsActive)

	// 查詢所有在新北市板橋區且 is_open_to_hiring=true 的老師
	var teachers []struct {
		ID       uint
		Name     string
		Email    string
		City     string
		District string
	}

	db.Table("teachers").
		Where("city = ? AND district = ? AND is_open_to_hiring = ?", "新北市", "板橋區", true).
		Find(&teachers)

	log.Printf("\n新北市板橋區 + is_open_to_hiring=true 的老師:")
	for _, t := range teachers {
		log.Printf("  ID=%d: %s (%s)", t.ID, t.Name, t.Email)
	}
}
