package main

import (
	"fmt"
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

	// 查詢老師的 membership
	fmt.Println("=== center_memberships 表 ===")
	var memberships []map[string]interface{}
	db.Table("center_memberships").Find(&memberships)
	fmt.Printf("總筆數: %d\n", len(memberships))
	for _, m := range memberships {
		fmt.Printf("  TeacherID=%v, CenterID=%v, Status=%v\n",
			m["teacher_id"], m["center_id"], m["status"])
	}

	// 查詢 teachers 表
	fmt.Println("\n=== teachers 表 ===")
	var teachers []map[string]interface{}
	db.Table("teachers").Where("deleted_at IS NULL").Find(&teachers)
	fmt.Printf("總筆數: %d\n", len(teachers))
	for _, t := range teachers {
		fmt.Printf("  ID=%v, Name=%v\n", t["id"], t["name"])
	}

	// 查詢 centers 表
	fmt.Println("\n=== centers 表 ===")
	var centers []map[string]interface{}
	db.Table("centers").Find(&centers)
	fmt.Printf("總筆數: %d\n", len(centers))
	for _, c := range centers {
		fmt.Printf("  ID=%v, Name=%v\n", c["id"], c["name"])
	}
}
