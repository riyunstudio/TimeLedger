package mysql

import (
	"log"
)

// MigrateAdditionalColumns 新增缺失欄位的遷移
// 用於補齊 production 環境缺少的欄位
func (db *DB) MigrateAdditionalColumns() {
	// 1. rooms.cleaning_time - 教室清潔緩衝時間
	if err := db.addColumnIfNotExists("rooms", "cleaning_time", "INT NOT NULL DEFAULT 0"); err != nil {
		log.Printf("Warning: Failed to add cleaning_time column: %v", err)
	}

	// 2. teachers.is_active - 教師是否已激活
	if err := db.addColumnIfNotExists("teachers", "is_active", "BOOLEAN NOT NULL DEFAULT FALSE"); err != nil {
		log.Printf("Warning: Failed to add is_active column: %v", err)
	}

	// 3. teachers.invited_at - 邀請時間
	if err := db.addColumnIfNotExists("teachers", "invited_at", "DATETIME NULL"); err != nil {
		log.Printf("Warning: Failed to add invited_at column: %v", err)
	}

	// 4. teachers.activated_at - 激活時間
	if err := db.addColumnIfNotExists("teachers", "activated_at", "DATETIME NULL"); err != nil {
		log.Printf("Warning: Failed to add activated_at column: %v", err)
	}

	// 5. schedule_rules.is_active - 規則是否有效
	if err := db.addColumnIfNotExists("schedule_rules", "is_active", "BOOLEAN NOT NULL DEFAULT TRUE"); err != nil {
		log.Printf("Warning: Failed to add is_active column to schedule_rules: %v", err)
	}

	// 6. schedule_rules.skip_holiday - 是否跳過假日（預設 true）
	if err := db.addColumnIfNotExists("schedule_rules", "skip_holiday", "BOOLEAN NOT NULL DEFAULT TRUE"); err != nil {
		log.Printf("Warning: Failed to add skip_holiday column to schedule_rules: %v", err)
	}

	// 7. center_holidays.force_cancel - 是否強制取消課堂（預設 false）
	if err := db.addColumnIfNotExists("center_holidays", "force_cancel", "BOOLEAN NOT NULL DEFAULT FALSE"); err != nil {
		log.Printf("Warning: Failed to add force_cancel column to center_holidays: %v", err)
	}

	// 8. schedule_rules.status - 課程狀態 (PLANNED/CONFIRMED)
	if err := db.addColumnIfNotExists("schedule_rules", "status", "VARCHAR(20) NOT NULL DEFAULT 'CONFIRMED'"); err != nil {
		log.Printf("Warning: Failed to add status column to schedule_rules: %v", err)
	}

	log.Println("Additional columns migration completed")
}

// addColumnIfNotExists 檢查並新增欄位
func (db *DB) addColumnIfNotExists(tableName, columnName, columnDef string) error {
	// 檢查欄位是否存在
	var count int
	result := db.WDB.Raw(`
		SELECT COUNT(*) INTO @a
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
		AND COLUMN_NAME = ?
	`, tableName, columnName).Scan(&count)

	if result.Error != nil {
		return result.Error
	}

	if count > 0 {
		log.Printf("Column %s.%s already exists, skipping", tableName, columnName)
		return nil
	}

	// 新增欄位
	addSQL := "ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + columnDef
	if err := db.WDB.Exec(addSQL).Error; err != nil {
		return err
	}

	log.Printf("Added column %s.%s", tableName, columnName)
	return nil
}
