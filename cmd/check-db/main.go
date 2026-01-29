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

	// 查詢所有資料表的基本統計
	tables := []string{
		"centers",
		"teachers",
		"schedule_rules",
		"center_holidays",
		"personal_events",
	}

	for _, table := range tables {
		var count int64
		db.Table(table).Count(&count)
		fmt.Printf("%s: %d 筆資料\n", table, count)
	}

	// 查詢 schedule_rules 的欄位資訊
	var rules []map[string]interface{}
	db.Table("schedule_rules").Limit(5).Find(&rules)
	if len(rules) > 0 {
		fmt.Printf("\n排課規則範例:\n")
		for k, v := range rules[0] {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
}
