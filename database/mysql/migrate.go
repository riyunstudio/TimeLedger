package mysql

import (
	"fmt"
	"log"
	"timeLedger/app/models"
)

// AutoMigrate 自動建立資料表
func (db *DB) AutoMigrate() {
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
