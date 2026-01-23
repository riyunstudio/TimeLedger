package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連接資料庫失敗: %v", err)
	}

	// 清理所有資料表
	tables := []string{
		"schedule_exceptions",
		"schedule_rules",
		"personal_events",
		"session_notes",
		"center_teacher_notes",
		"teacher_personal_hashtags",
		"teacher_skill_hashtags",
		"teacher_certificates",
		"teacher_skills",
		"hashtags",
		"center_memberships",
		"center_invitations",
		"center_holidays",
		"offerings",
		"courses",
		"rooms",
		"admin_users",
		"teachers",
		"centers",
	}

	fmt.Println("🧹 清理測試資料庫...")
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			fmt.Printf("  ⚠️  清理 %s 失敗: %v\n", table, err)
		} else {
			fmt.Printf("  ✅ 清理 %s 完成\n", table)
		}
	}

	// 重置自動遞增
	fmt.Println("\n🔄 重置自動遞增...")
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table)).Error; err != nil {
			fmt.Printf("  ⚠️  重置 %s 失敗: %v\n", table, err)
		}
	}

	fmt.Println("\n✅ 資料庫清理完成！")
}
