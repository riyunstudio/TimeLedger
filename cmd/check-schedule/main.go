package main

import (
	"context"
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

	ctx := context.Background()

	// 查詢所有有效的排課規則
	var rules []struct {
		ID        uint
		CenterID  uint
		Weekday   int
		StartTime string
		EndTime   string
	}

	// 簡化查詢
	result := db.WithContext(ctx).
		Table("schedule_rules").
		Where("deleted_at IS NULL").
		Find(&rules)

	if result.Error != nil {
		log.Fatalf("查詢失敗: %v", result.Error)
	}

	fmt.Printf("找到 %d 個排課規則:\n", len(rules))
	for _, r := range rules {
		fmt.Printf("  ID=%d, CenterID=%d, Weekday=%d, %s-%s\n", r.ID, r.CenterID, r.Weekday, r.StartTime, r.EndTime)
	}

	// 查詢假日
	var holidays []struct {
		CenterID uint
		Date     time.Time
		Name     string
	}

	holidayResult := db.WithContext(ctx).
		Table("center_holidays").
		Where("deleted_at IS NULL").
		Where("date >= ? AND date <= ?", "2026-01-01", "2026-01-31").
		Find(&holidays)

	if holidayResult.Error != nil {
		log.Fatalf("查詢假日失敗: %v", holidayResult.Error)
	}

	fmt.Printf("\n找到 %d 個假日:\n", len(holidays))
	for _, h := range holidays {
		fmt.Printf("  CenterID=%d, Date=%s, Name=%s\n", h.CenterID, h.Date.Format("2006-01-02"), h.Name)
	}
}
