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
			Name: "臺北市",
			Districts: []models.GeoDistrict{
				{Name: "中正區"}, {Name: "大同區"}, {Name: "中山區"}, {Name: "松山區"},
				{Name: "大安區"}, {Name: "萬華區"}, {Name: "信義區"}, {Name: "士林區"},
				{Name: "北投區"}, {Name: "內湖區"}, {Name: "南港區"}, {Name: "文山區"},
			},
		},
		{
			Name: "新北市",
			Districts: []models.GeoDistrict{
				{Name: "板橋區"}, {Name: "中和區"}, {Name: "永和區"}, {Name: "新莊區"},
				{Name: "新店區"}, {Name: "樹林區"}, {Name: "三重區"}, {Name: "淡水區"},
				{Name: "汐止區"}, {Name: "土城區"}, {Name: "蘆洲區"}, {Name: "五股區"},
				{Name: "泰山區"}, {Name: "林口區"}, {Name: "八里區"}, {Name: "三峽區"},
				{Name: "鶯歌區"}, {Name: "瑞芳區"}, {Name: "貢寮區"}, {Name: "雙溪區"},
				{Name: "平溪區"}, {Name: "石門區"}, {Name: "金山區"}, {Name: "萬里區"},
				{Name: "坪林區"}, {Name: "烏來區"},
			},
		},
		{
			Name: "桃園市",
			Districts: []models.GeoDistrict{
				{Name: "桃園區"}, {Name: "中壢區"}, {Name: "平鎮區"}, {Name: "八德區"},
				{Name: "楊梅區"}, {Name: "蘆竹區"}, {Name: "龜山區"}, {Name: "龍潭區"},
				{Name: "大溪區"}, {Name: "大園區"}, {Name: "觀音區"}, {Name: "新屋區"},
				{Name: "復興區"}, {Name: "新豐鄉"}, {Name: "湖口鄉"}, {Name: "長安鄉"},
			},
		},
		{
			Name: "臺中市",
			Districts: []models.GeoDistrict{
				{Name: "中區"}, {Name: "東區"}, {Name: "南區"}, {Name: "西區"},
				{Name: "北區"}, {Name: "北屯區"}, {Name: "西屯區"}, {Name: "南屯區"},
				{Name: "太平區"}, {Name: "大里區"}, {Name: "霧峰區"}, {Name: "烏日區"},
				{Name: "豐原區"}, {Name: "后里區"}, {Name: "東勢區"}, {Name: "石岡區"},
				{Name: "新社區"}, {Name: "潭子區"}, {Name: "大雅區"}, {Name: "神岡區"},
				{Name: "大肚區"}, {Name: "沙鹿區"}, {Name: "龍井區"}, {Name: "梧棲區"},
				{Name: "清水區"}, {Name: "大甲區"}, {Name: "外埔區"}, {Name: "大安區"},
				{Name: "和平區"},
			},
		},
		{
			Name: "臺南市",
			Districts: []models.GeoDistrict{
				{Name: "中西區"}, {Name: "東區"}, {Name: "南區"}, {Name: "北區"},
				{Name: "安平區"}, {Name: "安南區"}, {Name: "永康區"}, {Name: "歸仁區"},
				{Name: "新化區"}, {Name: "左鎮區"}, {Name: "玉井區"}, {Name: "楠西區"},
				{Name: "南化區"}, {Name: "仁德區"}, {Name: "關廟區"}, {Name: "龍崎區"},
				{Name: "官田區"}, {Name: "麻豆區"}, {Name: "佳里區"}, {Name: "西港區"},
				{Name: "七股區"}, {Name: "將軍區"}, {Name: "學甲區"}, {Name: "北門區"},
				{Name: "新營區"}, {Name: "後壁區"}, {Name: "白河區"}, {Name: "東山區"},
				{Name: "六甲區"}, {Name: "下營區"}, {Name: "柳營區"}, {Name: "鹽水區"},
				{Name: "善化區"}, {Name: "大內區"}, {Name: "山上區"}, {Name: "新市區"},
				{Name: "安定區"},
			},
		},
		{
			Name: "高雄市",
			Districts: []models.GeoDistrict{
				{Name: "新興區"}, {Name: "前金區"}, {Name: "苓雅區"}, {Name: "鹽埕區"},
				{Name: "鼓山區"}, {Name: "旗津區"}, {Name: "前鎮區"}, {Name: "三民區"},
				{Name: "左營區"}, {Name: "楠梓區"}, {Name: "小港區"}, {Name: "鳳山區"},
				{Name: "林園區"}, {Name: "大寮區"}, {Name: "大樹區"}, {Name: "大社區"},
				{Name: "仁武區"}, {Name: "鳥松區"}, {Name: "岡山區"}, {Name: "橋頭區"},
				{Name: "燕巢區"}, {Name: "田寮區"}, {Name: "阿蓮區"}, {Name: "路竹區"},
				{Name: "湖內區"}, {Name: "茄萣區"}, {Name: "永安區"}, {Name: "彌陀區"},
				{Name: "梓官區"}, {Name: "旗山區"}, {Name: "美濃區"}, {Name: "六龜區"},
				{Name: "甲仙區"}, {Name: "杉林區"}, {Name: "內門區"}, {Name: "茂林區"},
				{Name: "桃源區"}, {Name: "那瑪夏區"},
			},
		},
		{
			Name: "基隆市",
			Districts: []models.GeoDistrict{
				{Name: "仁愛區"}, {Name: "信義區"}, {Name: "中正區"}, {Name: "中山區"},
				{Name: "安樂區"}, {Name: "暖暖區"}, {Name: "七堵區"},
			},
		},
		{
			Name: "新竹市",
			Districts: []models.GeoDistrict{
				{Name: "東區"}, {Name: "北區"}, {Name: "香山區"},
			},
		},
		{
			Name: "新竹縣",
			Districts: []models.GeoDistrict{
				{Name: "竹北市"}, {Name: "竹東鎮"}, {Name: "新埔鎮"}, {Name: "關西鎮"},
				{Name: "湖口鄉"}, {Name: "新豐鄉"}, {Name: "峨眉鄉"}, {Name: "五峰鄉"},
				{Name: "橫山鄉"}, {Name: "芎林鄉"}, {Name: "寶山鄉"}, {Name: "北埔鄉"},
				{Name: "尖石鄉"},
			},
		},
		{
			Name: "苗栗縣",
			Districts: []models.GeoDistrict{
				{Name: "苗栗市"}, {Name: "苑裡鎮"}, {Name: "通霄鎮"}, {Name: "竹南鎮"},
				{Name: "頭份市"}, {Name: "後龍鎮"}, {Name: "卓蘭鎮"}, {Name: "大湖鄉"},
				{Name: "公館鄉"}, {Name: "銅鑼鄉"}, {Name: "三義鄉"}, {Name: "西湖鄉"},
				{Name: "造橋鄉"}, {Name: "頭屋鄉"}, {Name: "三灣鄉"}, {Name: "南庄鄉"},
				{Name: "泰安鄉"},
			},
		},
		{
			Name: "彰化縣",
			Districts: []models.GeoDistrict{
				{Name: "彰化市"}, {Name: "鹿港鎮"}, {Name: "和美鎮"}, {Name: "線西鄉"},
				{Name: "伸港鄉"}, {Name: "福興鄉"}, {Name: "秀水鄉"}, {Name: "花壇鄉"},
				{Name: "芬園鄉"}, {Name: "大村鄉"}, {Name: "員林市"}, {Name: "永靖鄉"},
				{Name: "社頭鄉"}, {Name: "二水鄉"}, {Name: "北斗鎮"}, {Name: "田尾鄉"},
				{Name: "埤頭鄉"}, {Name: "溪州鄉"}, {Name: "竹塘鄉"}, {Name: "大城鄉"},
				{Name: "芳苑鄉"}, {Name: "二林鎮"},
			},
		},
		{
			Name: "南投縣",
			Districts: []models.GeoDistrict{
				{Name: "南投市"}, {Name: "埔里鎮"}, {Name: "草屯鎮"}, {Name: "竹山鎮"},
				{Name: "集集鎮"}, {Name: "名間鄉"}, {Name: "鹿谷鄉"}, {Name: "中寮鄉"},
				{Name: "魚池鄉"}, {Name: "國姓鄉"}, {Name: "水里鄉"}, {Name: "信義鄉"},
				{Name: "仁愛鄉"},
			},
		},
		{
			Name: "雲林縣",
			Districts: []models.GeoDistrict{
				{Name: "斗六市"}, {Name: "斗南鎮"}, {Name: "虎尾鎮"}, {Name: "西螺鎮"},
				{Name: "土庫鎮"}, {Name: "北港鎮"}, {Name: "林內鄉"}, {Name: "古坑鄉"},
				{Name: "大埤鄉"}, {Name: "莿桐鄉"}, {Name: "二崙鄉"}, {Name: "崙背鄉"},
				{Name: "麥寮鄉"}, {Name: "東勢鄉"}, {Name: "褒忠鄉"}, {Name: "臺西鄉"},
				{Name: "元長鄉"}, {Name: "四湖鄉"}, {Name: "口湖鄉"}, {Name: "水林鄉"},
			},
		},
		{
			Name: "嘉義市",
			Districts: []models.GeoDistrict{
				{Name: "東區"}, {Name: "西區"},
			},
		},
		{
			Name: "嘉義縣",
			Districts: []models.GeoDistrict{
				{Name: "太保市"}, {Name: "朴子市"}, {Name: "布袋鎮"}, {Name: "民雄鄉"},
				{Name: "新港鄉"}, {Name: "六腳鄉"}, {Name: "東石鄉"}, {Name: "義竹鄉"},
				{Name: "鹿草鄉"}, {Name: "水上鄉"}, {Name: "中埔鄉"}, {Name: "竹崎鄉"},
				{Name: "梅山鄉"}, {Name: "番路鄉"}, {Name: "大埔鄉"}, {Name: "阿里山鄉"},
			},
		},
		{
			Name: "屏東縣",
			Districts: []models.GeoDistrict{
				{Name: "屏東市"}, {Name: "潮州鎮"}, {Name: "東港鎮"}, {Name: "恆春鎮"},
				{Name: "萬丹鄉"}, {Name: "長治鄉"}, {Name: "麟洛鄉"}, {Name: "九如鄉"},
				{Name: "里港鄉"}, {Name: "鹽埔鄉"}, {Name: "高樹鄉"}, {Name: "萬巒鄉"},
				{Name: "內埔鄉"}, {Name: "竹田鄉"}, {Name: "新埤鄉"}, {Name: "枋寮鄉"},
				{Name: "新園鄉"}, {Name: "崁頂鄉"}, {Name: "林邊鄉"}, {Name: "南州鄉"},
				{Name: "佳冬鄉"}, {Name: "琉球鄉"}, {Name: "車城鄉"}, {Name: "滿州鄉"},
				{Name: "枋山鄉"}, {Name: "三地門鄉"}, {Name: "霧臺鄉"}, {Name: "瑪家鄉"},
				{Name: "泰武鄉"}, {Name: "來義鄉"}, {Name: "春日鄉"}, {Name: "獅子鄉"},
				{Name: "牡丹鄉"},
			},
		},
		{
			Name: "宜蘭縣",
			Districts: []models.GeoDistrict{
				{Name: "宜蘭市"}, {Name: "羅東鎮"}, {Name: "蘇澳鎮"}, {Name: "頭城鎮"},
				{Name: "礁溪鄉"}, {Name: "壯圍鄉"}, {Name: "員山鄉"}, {Name: "大同鄉"},
				{Name: "三星鄉"}, {Name: "五結鄉"}, {Name: "冬山鄉"}, {Name: "南澳鄉"},
				{Name: "三星鄉"},
			},
		},
		{
			Name: "花蓮縣",
			Districts: []models.GeoDistrict{
				{Name: "花蓮市"}, {Name: "鳳林鎮"}, {Name: "玉里鎮"}, {Name: "新城鄉"},
				{Name: "吉安鄉"}, {Name: "壽豐鄉"}, {Name: "秀林鄉"}, {Name: "新城鄉"},
				{Name: "光復鄉"}, {Name: "豐濱鄉"}, {Name: "瑞穗鄉"}, {Name: "富里鄉"},
				{Name: "卓溪鄉"},
			},
		},
		{
			Name: "臺東縣",
			Districts: []models.GeoDistrict{
				{Name: "臺東市"}, {Name: "成功鎮"}, {Name: "關山鎮"}, {Name: "卑南鄉"},
				{Name: "鹿野鄉"}, {Name: "池上鄉"}, {Name: "東河鄉"}, {Name: "長濱鄉"},
				{Name: "太麻里鄉"}, {Name: "金峰鄉"}, {Name: "大武鄉"}, {Name: "達仁鄉"},
				{Name: "海端鄉"}, {Name: "延平鄉"}, {Name: "綠島鄉"}, {Name: "蘭嶼鄉"},
			},
		},
		{
			Name: "澎湖縣",
			Districts: []models.GeoDistrict{
				{Name: "馬公市"}, {Name: "湖西鄉"}, {Name: "白沙鄉"}, {Name: "西嶼鄉"},
				{Name: "望安鄉"}, {Name: "七美鄉"},
			},
		},
		{
			Name: "金門縣",
			Districts: []models.GeoDistrict{
				{Name: "金城鎮"}, {Name: "金湖鎮"}, {Name: "金沙鎮"}, {Name: "金寧鄉"},
				{Name: "烈嶼鄉"}, {Name: "烏坵鄉"},
			},
		},
		{
			Name: "連江縣",
			Districts: []models.GeoDistrict{
				{Name: "南竿鄉"}, {Name: "北竿鄉"}, {Name: "莒光鄉"}, {Name: "東引鄉"},
			},
		},
	}

	for _, city := range cities {
		db.WDB.Create(&city)
	}
	log.Printf("Created %d cities with districts", len(cities))
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

	// 建立中心會籍（將老師與中心關聯）
	var center models.Center
	db.WDB.Where("name = ?", "莫札特音樂教室").First(&center)

	membership := models.CenterMembership{
		CenterID:  center.ID,
		TeacherID: teacher.ID,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.WDB.Create(&membership)

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
		PasswordHash: "$2a$10$nZsYJrENRJoW1yLxuZPu0.H4L533HjUMU26pr1LiM0/4VppE02BpC",
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
