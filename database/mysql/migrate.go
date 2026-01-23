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
		&models.SessionNote{},
		&models.AuditLog{},
		&models.Notification{},
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
