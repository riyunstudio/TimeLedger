package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
)

func main() {
	a := app.Initialize()
	defer func() {
		sqlDB, _ := a.MySQL.WDB.DB()
		sqlDB.Close()
	}()

	// 檢查是否已存在
	var existing models.CenterMembership
	result := a.MySQL.WDB.Where("center_id = ? AND teacher_id = ?", 1, 40).First(&existing)
	if result.Error == nil {
		fmt.Println("會籍已存在，無需重複建立")
		os.Exit(0)
	}

	// 建立會籍
	membership := models.CenterMembership{
		CenterID:  1,
		TeacherID: 40,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.MySQL.WDB.Create(&membership).Error; err != nil {
		log.Fatal("建立會籍失敗:", err)
	}

	fmt.Println("老師中心會籍建立成功!")
	fmt.Printf("Teacher ID: %d, Center ID: %d\n", membership.TeacherID, membership.CenterID)
}
