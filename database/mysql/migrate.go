package mysql

import (
	"fmt"
	"log"
	"timeLedger/app/models"
)

// AutoMigrate 自動建立資料表
func (db *DB) AutoMigrate() {
	// 先刪除有問題的資料表，讓 GORM 重新建立（解決外鍵類型不相容問題）
	// dropTablesForReschema(db)

	if err := db.WDB.AutoMigrate(
		&models.User{},
		&models.Center{},
		&models.AdminUser{},
		&models.Teacher{},
		&models.CenterMembership{},
		&models.CenterInvitation{},
		&models.GeoCity{},
		&models.GeoDistrict{},
		&models.Course{},
		&models.Offering{},
		&models.Room{},
		&models.TimetableTemplate{},
		&models.TimetableCell{},
		&models.ScheduleRule{},
		&models.ScheduleException{},
		&models.PersonalEvent{},
		&models.TeacherSkill{},
		&models.Hashtag{},
		&models.TeacherSkillHashtag{},
		&models.TeacherPersonalHashtag{},
		&models.TeacherCertificate{},
		&models.CenterTeacherNote{},
		&models.CenterHoliday{},
		&models.SessionNote{},
		&models.AuditLog{},
		&models.Notification{},
		&models.NotificationQueue{},
	); err != nil {
		panic(fmt.Errorf("MySQL autoMigrate failed: %v", err))
	}
	log.Println("AutoMigrate done")
}

// dropTablesForReschema 刪除需要重建的資料表
func dropTablesForReschema(db *DB) {
	// 先刪除有外鍵依賴的表
	db.WDB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 刪除可能存在問題的表
	db.WDB.Exec("DROP TABLE IF EXISTS `schedule_rules`")
	db.WDB.Exec("DROP TABLE IF EXISTS `schedule_exceptions`")

	db.WDB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("Dropped schedule_rules and schedule_exceptions tables for reschema")
}

// MigrateScheduleExceptionsType 遷移 schedule_exceptions 表的 type 欄位
// 將舊的 type 欄位遷移到新的 exception_type 欄位
func (db *DB) MigrateScheduleExceptionsType() {
	// 檢查是否存在舊的 type 欄位
	var columnExists bool
	db.WDB.Raw("SELECT COUNT(*) INTO @a FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'schedule_exceptions' AND COLUMN_NAME = 'type'").Scan(&columnExists)

	// 檢查是否存在新的 exception_type 欄位
	var newColumnExists bool
	db.WDB.Raw("SELECT COUNT(*) INTO @a FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'schedule_exceptions' AND COLUMN_NAME = 'exception_type'").Scan(&newColumnExists)

	if columnExists && !newColumnExists {
		// 新增 exception_type 欄位
		if err := db.WDB.Exec("ALTER TABLE schedule_exceptions ADD COLUMN exception_type VARCHAR(20) NOT NULL DEFAULT 'CANCEL' AFTER original_date").Error; err != nil {
			log.Printf("Warning: Failed to add exception_type column: %v", err)
		}

		// 複製資料
		if err := db.WDB.Exec("UPDATE schedule_exceptions SET exception_type = type WHERE type IS NOT NULL AND type != ''").Error; err != nil {
			log.Printf("Warning: Failed to copy data from type to exception_type: %v", err)
		}

		log.Println("Migrated schedule_exceptions.type to exception_type")
	} else if newColumnExists {
		log.Println("schedule_exceptions.exception_type column already exists")
	} else {
		log.Println("No migration needed for schedule_exceptions type column")
	}
}
