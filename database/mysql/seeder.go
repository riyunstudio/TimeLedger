package mysql

import (
	"log"
	"time"
	"timeLedger/app/models"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

// Seed 建立預設資料
func (db *DB) Seeds(tools *tools.Tools) {
	cleanAllTables(db)
	seedOneCenter(db)
	seedGeoData(db)
	seedOneTeacher(db)
	seedOneAdminUser(db)
	seedTestResources(db)
	log.Println("Database seed complete")
}

func cleanAllTables(db *DB) {
	tables := []string{
		"schedule_exceptions",
		"schedule_rules",
		"session_notes",
		"personal_events",
		"teacher_certificates",
		"teacher_skills",
		"center_memberships",
		"timetable_cells",
		"timetable_templates",
		"offerings",
		"audit_logs",
		"admin_users",
		"teachers",
		"courses",
		"rooms",
		"center_invitations",
		"notifications",
		"users",
		"geo_districts",
		"geo_cities",
		"centers",
	}

	for _, table := range tables {
		db.WDB.Exec("DELETE FROM " + table)
	}
	log.Println("All tables cleaned")
}

func seedOneCenter(db *DB) {
	center := models.Center{
		Name:      "莫札特音樂教室",
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			AllowPublicRegister: true,
			DefaultLanguage:     "zh-TW",
		},
		CreatedAt: time.Now(),
	}

	db.WDB.Create(&center)
	log.Printf("Created default center: %s (ID: %d)", center.Name, center.ID)
}

func seedGeoData(db *DB) {
	cities := []models.GeoCity{
		{
			Name: "台北市",
			Districts: []models.GeoDistrict{
				{Name: "大安區"},
				{Name: "信義區"},
				{Name: "中山區"},
			},
		},
		{
			Name: "新北市",
			Districts: []models.GeoDistrict{
				{Name: "板橋區"},
				{Name: "新莊區"},
			},
		},
	}

	for _, city := range cities {
		db.WDB.Create(&city)
	}
}

func seedOneTeacher(db *DB) {
	teacher := models.Teacher{
		LineUserID:        "LINE_TEACHER_001",
		Name:              "王小明",
		Email:             "wangxiaoming@example.com",
		Bio:               "專業鋼琴教師，十年教學經驗，國立台北藝術大學音樂系畢業",
		City:              "台北市",
		District:          "大安區",
		PublicContactInfo: "LineID: wangxiaoming",
		IsOpenToHiring:    true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	db.WDB.Create(&teacher)

	skill := models.TeacherSkill{
		TeacherID: teacher.ID,
		Category:  "音樂",
		SkillName: "鋼琴",
		Level:     "ADVANCED",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.WDB.Create(&skill)

	cert := models.TeacherCertificate{
		TeacherID: teacher.ID,
		Name:      "鋼琴演奏級認證",
		FileURL:   "https://example.com/certs/piano-license.pdf",
		IssuedAt:  time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.WDB.Create(&cert)

	log.Printf("Created teacher: %s (ID: %d)", teacher.Name, teacher.ID)
}

func seedOneAdminUser(db *DB) {
	admin := models.AdminUser{
		Email:        "admin@timeledger.com",
		PasswordHash: "$2a$10$wDC8I8iP0LJgkXoUEcxA0uy6S4O/KfDzExabt7YxpD6jtWMHzfyse",
		Name:         "系統管理員",
		CenterID:     1,
		Role:         "OWNER",
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	db.WDB.Create(&admin)
	log.Printf("Created admin: %s", admin.Email)
}

func seedTestResources(db *DB) {
	// 取得老師和中心的 ID
	var teacher models.Teacher
	db.WDB.Where("line_user_id = ?", "LINE_TEACHER_001").First(&teacher)

	var center models.Center
	db.WDB.Where("name = ?", "莫札特音樂教室").First(&center)

	// 建立教室
	rooms := []models.Room{
		{CenterID: center.ID, Name: "鋼琴教室A", Capacity: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, Name: "鋼琴教室B", Capacity: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, Name: "團體教室", Capacity: 15, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, room := range rooms {
		db.WDB.Create(&room)
	}
	log.Printf("Created %d rooms", len(rooms))

	// 建立課程
	courses := []models.Course{
		{CenterID: center.ID, Name: "鋼琴基礎班", DefaultDuration: 60, ColorHex: "#10B981", TeacherBufferMin: 10, RoomBufferMin: 10, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, Name: "鋼琴進階班", DefaultDuration: 60, ColorHex: "#3B82F6", TeacherBufferMin: 10, RoomBufferMin: 10, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, Name: "兒童音樂啟蒙", DefaultDuration: 45, ColorHex: "#F59E0B", TeacherBufferMin: 5, RoomBufferMin: 5, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, course := range courses {
		db.WDB.Create(&course)
	}
	log.Printf("Created %d courses", len(courses))

	// 取得課程 ID
	var coursesResult []models.Course
	db.WDB.Where("center_id = ?", center.ID).Find(&coursesResult)

	// 取得教室 ID
	var roomsResult []models.Room
	db.WDB.Where("center_id = ?", center.ID).Find(&roomsResult)

	// 建立班別 (Offering) - 關聯課程與老師
	teacherID := teacher.ID
	offerings := []models.Offering{
		{CenterID: center.ID, CourseID: coursesResult[0].ID, Name: "鋼琴基礎班 - 每週一", DefaultRoomID: &roomsResult[0].ID, DefaultTeacherID: &teacherID, AllowBufferOverride: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, CourseID: coursesResult[0].ID, Name: "鋼琴基礎班 - 每週三", DefaultRoomID: &roomsResult[0].ID, DefaultTeacherID: &teacherID, AllowBufferOverride: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, CourseID: coursesResult[1].ID, Name: "鋼琴進階班", DefaultRoomID: &roomsResult[1].ID, DefaultTeacherID: &teacherID, AllowBufferOverride: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, CourseID: coursesResult[2].ID, Name: "兒童音樂啟蒙班", DefaultRoomID: &roomsResult[2].ID, DefaultTeacherID: &teacherID, AllowBufferOverride: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, offering := range offerings {
		db.WDB.Create(&offering)
	}
	log.Printf("Created %d offerings", len(offerings))

	// 建立課程時段 (ScheduleRule)
	var offeringsResult []models.Offering
	db.WDB.Where("center_id = ?", center.ID).Find(&offeringsResult)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	rules := []models.ScheduleRule{
		{CenterID: center.ID, OfferingID: offeringsResult[0].ID, TeacherID: &teacherID, RoomID: roomsResult[0].ID, Weekday: 1, StartTime: "10:00", EndTime: "11:00", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, OfferingID: offeringsResult[0].ID, TeacherID: &teacherID, RoomID: roomsResult[0].ID, Weekday: 1, StartTime: "14:00", EndTime: "15:00", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, OfferingID: offeringsResult[1].ID, TeacherID: &teacherID, RoomID: roomsResult[0].ID, Weekday: 3, StartTime: "10:00", EndTime: "11:00", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, OfferingID: offeringsResult[1].ID, TeacherID: &teacherID, RoomID: roomsResult[0].ID, Weekday: 3, StartTime: "14:00", EndTime: "15:00", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, OfferingID: offeringsResult[2].ID, TeacherID: &teacherID, RoomID: roomsResult[1].ID, Weekday: 2, StartTime: "16:00", EndTime: "17:00", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, OfferingID: offeringsResult[2].ID, TeacherID: &teacherID, RoomID: roomsResult[1].ID, Weekday: 4, StartTime: "16:00", EndTime: "17:00", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, OfferingID: offeringsResult[3].ID, TeacherID: &teacherID, RoomID: roomsResult[2].ID, Weekday: 6, StartTime: "09:00", EndTime: "10:00", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{CenterID: center.ID, OfferingID: offeringsResult[3].ID, TeacherID: &teacherID, RoomID: roomsResult[2].ID, Weekday: 6, StartTime: "10:30", EndTime: "11:30", EffectiveRange: models.DateRange{StartDate: startDate, EndDate: endDate}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, rule := range rules {
		db.WDB.Create(&rule)
	}
	log.Printf("Created %d schedule rules", len(rules))
}
