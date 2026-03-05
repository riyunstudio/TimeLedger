package mysql

import (
	"log"
)

// MigrateCourseCategories 新增課程類別表的遷移
func (db *DB) MigrateCourseCategories() {
	// 檢查表是否已存在
	var count int64
	if err := db.WDB.Raw(`
		SELECT COUNT(*) FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = 'course_categories'
	`).Scan(&count).Error; err != nil {
		log.Printf("Warning: Failed to check course_categories table: %v", err)
		return
	}

	if count > 0 {
		log.Println("course_categories table already exists, skipping")
		return
	}

	// 建立 course_categories 表
	if err := db.WDB.Exec(`
		CREATE TABLE course_categories (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			center_id BIGINT UNSIGNED NOT NULL,
			name VARCHAR(50) NOT NULL,
			INDEX idx_center_id (center_id),
			UNIQUE INDEX idx_center_name (center_id, name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`).Error; err != nil {
		log.Printf("Warning: Failed to create course_categories table: %v", err)
		return
	}

	log.Println("course_categories table created successfully")
}

// MigrateCourseCategoryToNewTable 將中心設定中的 course_categories 遷移到獨立表
// 這個遷移是一次性的，用於將現有資料從 CenterSettings JSON 遷移到獨立表
func (db *DB) MigrateCourseCategoryToNewTable() {
	// 檢查是否已經有遷移過的標記
	var count int64
	if err := db.WDB.Raw(`SELECT COUNT(*) FROM course_categories`).Scan(&count).Error; err != nil {
		// 表可能還沒建立，跳過
		log.Printf("course_categories table not found, skipping migration: %v", err)
		return
	}

	if count > 0 {
		// 已經有資料，跳過遷移
		log.Println("course_categories already migrated, skipping")
		return
	}

	// 從 centers 表讀取現有的 course_categories
	type CenterWithSettings struct {
		ID       uint
		Settings string `gorm:"column:settings"`
	}

	var centers []CenterWithSettings
	if err := db.WDB.Table("centers").Select("id", "settings").Find(&centers).Error; err != nil {
		log.Printf("Failed to read centers: %v", err)
		return
	}

	// 解析 settings 並遷移
	for _, center := range centers {
		if center.Settings == "" {
			continue
		}

		// 使用 GORM 的 JSON 解析
		type SettingsData struct {
			CourseCategories []string `json:"course_categories"`
		}

		var settings SettingsData
		if err := db.WDB.Raw(`SELECT ? AS settings`, center.Settings).Scan(&settings).Error; err == nil {
			// 手動解析
			var parsed map[string]interface{}
			if err := db.WDB.Raw(`SELECT settings FROM centers WHERE id = ?`, center.ID).Scan(&parsed).Error; err == nil {
				if cats, ok := parsed["course_categories"].([]interface{}); ok {
					for _, cat := range cats {
						if name, ok := cat.(string); ok && name != "" {
							db.WDB.Exec(`INSERT INTO course_categories (center_id, name) VALUES (?, ?)`, center.ID, name)
						}
					}
				}
			}
		}
	}

	log.Println("course_categories migration completed")
}

// RollbackCourseCategories 回滾遷移（刪除獨立表）
func (db *DB) RollbackCourseCategories() {
	if err := db.WDB.Exec(`DROP TABLE IF EXISTS course_categories`).Error; err != nil {
		log.Printf("Warning: Failed to drop course_categories table: %v", err)
		return
	}
	log.Println("course_categories table dropped")
}
