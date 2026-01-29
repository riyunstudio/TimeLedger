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

	// 恢復所有被軟刪除的排課規則
	result := db.Table("schedule_rules").
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		log.Fatalf("恢復排課規則失敗: %v", result.Error)
	}

	fmt.Printf("已恢復 %d 筆排課規則\n", result.RowsAffected)

	// 恢復假日
	result3 := db.Table("center_holidays").
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result3.Error != nil {
		log.Fatalf("恢復假日失敗: %v", result3.Error)
	}

	fmt.Printf("已恢復 %d 筆假日\n", result3.RowsAffected)

	// 驗證結果
	var activeRules int64
	db.Table("schedule_rules").Where("deleted_at IS NULL").Count(&activeRules)
	fmt.Printf("\n目前活躍的排課規則: %d 筆\n", activeRules)

	// 顯示恢復的規則
	var rules []map[string]interface{}
	db.Table("schedule_rules").Where("deleted_at IS NULL").Find(&rules)
	for _, r := range rules {
		fmt.Printf("  ID=%v, CenterID=%v, Weekday=%v, %v-%v\n", 
			r["id"], r["center_id"], r["weekday"], r["start_time"], r["end_time"])
	}
}
