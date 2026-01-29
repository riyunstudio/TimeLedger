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

	// 查詢未刪除的排課規則
	var activeRules []map[string]interface{}
	result := db.Table("schedule_rules").
		Where("deleted_at IS NULL").
		Limit(10).
		Find(&activeRules)

	if result.Error != nil {
		log.Fatalf("查詢失敗: %v", result.Error)
	}

	fmt.Printf("活躍的排課規則: %d 筆\n", len(activeRules))
	for _, r := range activeRules {
		fmt.Printf("  ID=%v, CenterID=%v, Weekday=%v, %v-%v\n", 
			r["id"], r["center_id"], r["weekday"], r["start_time"], r["end_time"])
	}

	// 查詢已刪除的排課規則
	var deletedRules []map[string]interface{}
	db.Table("schedule_rules").
		Where("deleted_at IS NOT NULL").
		Limit(10).
		Find(&deletedRules)
	
	fmt.Printf("\n已刪除的排課規則: %d 筆\n", len(deletedRules))
	for _, r := range deletedRules {
		fmt.Printf("  ID=%v, CenterID=%v, Weekday=%v\n", 
			r["id"], r["center_id"], r["weekday"])
	}

	// 查詢假日
	var holidays []map[string]interface{}
	db.Table("center_holidays").Where("deleted_at IS NULL").Find(&holidays)
	fmt.Printf("\n活躍的假日: %d 筆\n", len(holidays))
	for _, h := range holidays {
		fmt.Printf("  CenterID=%v, Date=%v, Name=%v\n", 
			h["center_id"], h["date"], h["name"])
	}
}
