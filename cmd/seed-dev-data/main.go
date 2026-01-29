package main

import (
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

	wdb := db.Session(&gorm.Session{Logger: db.Logger.LogMode(4)})

	// 清理所有資料表
	fmt.Println("正在清理所有資料表...")
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
		"center_holidays",
	}

	for _, table := range tables {
		wdb.Exec("DELETE FROM " + table)
	}
	fmt.Println("資料表已清理")

	// 建立中心
	fmt.Println("\n正在建立中心...")
	center := map[string]interface{}{
		"name":       "莫札特音樂教室",
		"plan_level": "STARTER",
		"settings":   `{"allow_public_register":true,"default_language":"zh-TW"}`,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	wdb.Table("centers").Create(&center)
	fmt.Printf("已建立中心: %s (ID: %d)\n", center["name"], center["id"])

	// 建立地理資料
	fmt.Println("\n正在建立地理資料...")
	cities := []map[string]interface{}{
		{"name": "臺北市", "created_at": time.Now()},
		{"name": "新北市", "created_at": time.Now()},
		{"name": "桃園市", "created_at": time.Now()},
	}
	for _, city := range cities {
		wdb.Table("geo_cities").Create(&city)
	}

	var taipeiCity, newTaipeiCity map[string]interface{}
	wdb.Table("geo_cities").Where("name = ?", "臺北市").First(&taipeiCity)
	wdb.Table("geo_cities").Where("name = ?", "新北市").First(&newTaipeiCity)

	taipeiDistricts := []string{"中正區", "大同區", "中山區", "松山區", "大安區", "萬華區", "信義區", "士林區", "北投區", "內湖區", "南港區", "文山區"}
	newTaipeiDistricts := []string{"板橋區", "中和區", "永和區", "新莊區", "新店區", "樹林區", "三重區", "淡水區", "汐止區", "土城區", "蘆洲區"}

	for _, d := range taipeiDistricts {
		wdb.Table("geo_districts").Create(&map[string]interface{}{
			"city_id":    taipeiCity["id"],
			"name":       d,
			"created_at": time.Now(),
		})
	}
	for _, d := range newTaipeiDistricts {
		wdb.Table("geo_districts").Create(&map[string]interface{}{
			"city_id":    newTaipeiCity["id"],
			"name":       d,
			"created_at": time.Now(),
		})
	}
	fmt.Println("已建立地理資料")

	// 建立老師
	fmt.Println("\n正在建立老師...")
	teacher := map[string]interface{}{
		"line_user_id":        "LINE_TEACHER_001",
		"name":                "王小明",
		"email":               "wangxiaoming@example.com",
		"bio":                 "專業鋼琴教師，十年教學經驗",
		"city":                "臺北市",
		"district":            "大安區",
		"public_contact_info": "LineID: wangxiaoming",
		"is_open_to_hiring":   true,
		"is_active":           true,
		"created_at":          time.Now(),
		"updated_at":          time.Now(),
	}
	wdb.Table("teachers").Create(&teacher)
	fmt.Printf("已建立老師: %s (ID: %d)\n", teacher["name"], teacher["id"])

	// 建立會籍
	membership := map[string]interface{}{
		"center_id":  center["id"],
		"teacher_id": teacher["id"],
		"status":     "ACTIVE",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	wdb.Table("center_memberships").Create(&membership)
	fmt.Println("已建立會籍")

	// 建立管理員
	fmt.Println("\n正在建立管理員...")
	admin := map[string]interface{}{
		"email":         "admin@timeledger.com",
		"password_hash": "$2a$10$nZsYJrENRJoW1yLxuZPu0.H4L533HjUMU26pr1LiM0/4VppE02BpC",
		"name":          "系統管理員",
		"center_id":     center["id"],
		"role":          "OWNER",
		"status":        "ACTIVE",
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}
	wdb.Table("admin_users").Create(&admin)
	fmt.Println("已建立管理員")

	// 建立教室
	fmt.Println("\n正在建立教室...")
	rooms := []map[string]interface{}{
		{"center_id": center["id"], "name": "鋼琴教室A", "capacity": 2, "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "name": "鋼琴教室B", "capacity": 4, "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "name": "團體教室", "capacity": 15, "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
	}
	for _, room := range rooms {
		wdb.Table("rooms").Create(&room)
	}
	fmt.Printf("已建立 %d 間教室\n", len(rooms))

	// 建立課程
	fmt.Println("\n正在建立課程...")
	courses := []map[string]interface{}{
		{"center_id": center["id"], "name": "鋼琴基礎班", "default_duration": 60, "color_hex": "#10B981", "teacher_buffer_min": 10, "room_buffer_min": 10, "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "name": "鋼琴進階班", "default_duration": 60, "color_hex": "#3B82F6", "teacher_buffer_min": 10, "room_buffer_min": 10, "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "name": "兒童音樂啟蒙", "default_duration": 45, "color_hex": "#F59E0B", "teacher_buffer_min": 5, "room_buffer_min": 5, "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
	}
	for _, course := range courses {
		wdb.Table("courses").Create(&course)
	}
	fmt.Printf("已建立 %d 門課程\n", len(courses))

	// 建立班別 (Offering)
	fmt.Println("\n正在建立班別...")
	var coursesResult []map[string]interface{}
	var roomsResult []map[string]interface{}
	wdb.Table("courses").Where("center_id = ?", center["id"]).Find(&coursesResult)
	wdb.Table("rooms").Where("center_id = ?", center["id"]).Find(&roomsResult)

	teacherID := teacher["id"]
	offerings := []map[string]interface{}{
		{"center_id": center["id"], "course_id": coursesResult[0]["id"], "name": "鋼琴基礎班 - 每週一", "default_room_id": roomsResult[0]["id"], "default_teacher_id": teacherID, "allow_buffer_override": false, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "course_id": coursesResult[0]["id"], "name": "鋼琴基礎班 - 每週三", "default_room_id": roomsResult[0]["id"], "default_teacher_id": teacherID, "allow_buffer_override": false, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "course_id": coursesResult[1]["id"], "name": "鋼琴進階班", "default_room_id": roomsResult[1]["id"], "default_teacher_id": teacherID, "allow_buffer_override": false, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "course_id": coursesResult[2]["id"], "name": "兒童音樂啟蒙班", "default_room_id": roomsResult[2]["id"], "default_teacher_id": teacherID, "allow_buffer_override": false, "created_at": time.Now(), "updated_at": time.Now()},
	}
	for _, offering := range offerings {
		wdb.Table("offerings").Create(&offering)
	}
	fmt.Printf("已建立 %d 個班別\n", len(offerings))

	// 建立排課規則
	fmt.Println("\n正在建立排課規則...")
	var offeringsResult []map[string]interface{}
	wdb.Table("offerings").Where("center_id = ?", center["id"]).Find(&offeringsResult)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	rules := []map[string]interface{}{
		{"center_id": center["id"], "offering_id": offeringsResult[0]["id"], "teacher_id": teacherID, "room_id": roomsResult[0]["id"], "weekday": 1, "start_time": "10:00", "end_time": "11:00", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "offering_id": offeringsResult[0]["id"], "teacher_id": teacherID, "room_id": roomsResult[0]["id"], "weekday": 1, "start_time": "14:00", "end_time": "15:00", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "offering_id": offeringsResult[1]["id"], "teacher_id": teacherID, "room_id": roomsResult[0]["id"], "weekday": 3, "start_time": "10:00", "end_time": "11:00", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "offering_id": offeringsResult[1]["id"], "teacher_id": teacherID, "room_id": roomsResult[0]["id"], "weekday": 3, "start_time": "14:00", "end_time": "15:00", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "offering_id": offeringsResult[2]["id"], "teacher_id": teacherID, "room_id": roomsResult[1]["id"], "weekday": 2, "start_time": "16:00", "end_time": "17:00", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "offering_id": offeringsResult[2]["id"], "teacher_id": teacherID, "room_id": roomsResult[1]["id"], "weekday": 4, "start_time": "16:00", "end_time": "17:00", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "offering_id": offeringsResult[3]["id"], "teacher_id": teacherID, "room_id": roomsResult[2]["id"], "weekday": 6, "start_time": "09:00", "end_time": "10:00", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		{"center_id": center["id"], "offering_id": offeringsResult[3]["id"], "teacher_id": teacherID, "room_id": roomsResult[2]["id"], "weekday": 6, "start_time": "10:30", "end_time": "11:30", "effective_range": fmt.Sprintf(`{"start_date":"%s","end_date":"%s"}`, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z")), "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
	}
	for _, rule := range rules {
		wdb.Table("schedule_rules").Create(&rule)
	}
	fmt.Printf("已建立 %d 個排課規則\n", len(rules))

	fmt.Println("\n========================================")
	fmt.Println("測試資料初始化完成！")
	fmt.Println("========================================")
	fmt.Println("\n登入資訊：")
	fmt.Println("  老師端：使用 LINE 登入（測試帳號）")
	fmt.Println("  管理員：admin@timeledger.com / password123")
}
