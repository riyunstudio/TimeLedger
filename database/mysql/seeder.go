package mysql

import (
	"fmt"
	"log"
	"time"
	"timeLedger/app/models"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

// Seed 建立測試資料
func (db *DB) Seeds(tools *tools.Tools) {
	cleanAllTables(db)
	seedCenters(db)
	seedGeoData(db)
	seedOneTeacher(db)
	seedOneAdminUser(db)
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

func seedCenters(db *DB) {
	centers := []models.Center{
		{
			Name:      "莫札特音樂教室",
			PlanLevel: "STARTER",
			Settings: models.CenterSettings{
				AllowPublicRegister: true,
				DefaultLanguage:     "zh-TW",
			},
			CreatedAt: time.Now(),
		},
		{
			Name:      "舞動舞蹈學院",
			PlanLevel: "PRO",
			Settings: models.CenterSettings{
				AllowPublicRegister: false,
				DefaultLanguage:     "zh-TW",
			},
			CreatedAt: time.Now(),
		},
		{
			Name:      "專業健身中心",
			PlanLevel: "TEAM",
			Settings: models.CenterSettings{
				AllowPublicRegister: true,
				DefaultLanguage:     "zh-TW",
			},
			CreatedAt: time.Now(),
		},
	}

	for _, center := range centers {
		var exists int64
		db.WDB.Model(&models.Center{}).Where("name = ?", center.Name).Count(&exists)
		if exists == 0 {
			db.WDB.Create(&center)
		}
	}
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
		var exists int64
		db.WDB.Model(&models.GeoCity{}).Where("name = ?", city.Name).Count(&exists)
		if exists == 0 {
			db.WDB.Create(&city)
		}
	}
}

func seedTeachers(db *DB, count int) {
	districts := []string{"大安區", "信義區", "中山區"}
	skills := []struct {
		category string
		name     string
		level    string
	}{
		{"音樂", "鋼琴", "ADVANCED"},
		{"音樂", "吉他", "INTERMEDIATE"},
		{"舞蹈", "街舞", "ADVANCED"},
		{"舞蹈", "芭蕾舞", "INTERMEDIATE"},
		{"瑜伽", "空中瑜伽", "ADVANCED"},
		{"瑜伽", "哈達瑜伽", "BASIC"},
	}

	for i := 1; i <= count; i++ {
		teacher := models.Teacher{
			LineUserID:        fmt.Sprintf("LINE_USER_%03d", i),
			Name:              fmt.Sprintf("老師%d", i),
			Email:             fmt.Sprintf("teacher%d@example.com", i),
			Bio:               fmt.Sprintf("專業教師%d，熱愛教學", i),
			City:              "台北市",
			District:          districts[i%len(districts)],
			PublicContactInfo: fmt.Sprintf("LineID: teacher%d", i),
			IsOpenToHiring:    i%3 == 0,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		var exists int64
		db.WDB.Model(&models.Teacher{}).Where("line_user_id = ?", teacher.LineUserID).Count(&exists)
		if exists == 0 {
			db.WDB.Create(&teacher)

			skill := skills[i%len(skills)]
			teacherSkill := models.TeacherSkill{
				TeacherID: teacher.ID,
				Category:  skill.category,
				SkillName: skill.name,
				Level:     skill.level,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			db.WDB.Create(&teacherSkill)

			if i%2 == 0 {
				cert := models.TeacherCertificate{
					TeacherID: teacher.ID,
					Name:      fmt.Sprintf("%s認證證書", skill.name),
					FileURL:   fmt.Sprintf("https://example.com/certs/%d.pdf", i),
					IssuedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				db.WDB.Create(&cert)
			}
		}
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

func seedUsers(db *DB) {
	users := []models.User{
		{
			Name: "阿卡莉",
			Ips:  `["192.168.1.10", "10.0.0.5"]`,
		},
	}

	for _, user := range users {
		var exists int64
		db.WDB.Model(&models.User{}).Where("name = ?", user.Name).Count(&exists)
		if exists == 0 {
			db.WDB.Create(&user)
		}
	}
}

func seedAdminUsers(db *DB) {
	adminUsers := []models.AdminUser{
		{
			Email:        "admin@timeledger.com",
			PasswordHash: "$2a$10$wDC8I8iP0LJgkXoUEcxA0uy6S4O/KfDzExabt7YxpD6jtWMHzfyse",
			Name:         "系統管理員",
			CenterID:     1,
			Role:         "OWNER",
			Status:       "ACTIVE",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, admin := range adminUsers {
		var exists int64
		db.WDB.Model(&models.AdminUser{}).Where("email = ?", admin.Email).Count(&exists)
		if exists == 0 {
			db.WDB.Create(&admin)
		}
	}
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
