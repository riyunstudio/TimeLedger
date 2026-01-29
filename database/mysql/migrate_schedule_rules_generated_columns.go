package mysql

import (
	"log"
)

// MigrateScheduleRulesGeneratedColumns 新增 schedule_rules 表的 Generated Columns
// 用於解決 effective_range JSON 欄位的全表掃描問題
//
// 問題說明：
// schedule_rules.effective_range 儲存為 JSON 格式（{"start_date": "...", "end_date": "..."}）
// 當查詢某日期範圍內的規則時，MySQL 無法使用索引，導致全表掃描
//
// 解決方案：
// 1. 先將 effective_range 中的日期正規化為 YYYY-MM-DD 格式
// 2. 新增 Generated Columns 從 JSON 中提取 start_date 和 end_date
// 3. 為這些欄位建立索引，大幅提升查詢效能
func (db *DB) MigrateScheduleRulesGeneratedColumns() {
	// 0. 先正規化 effective_range 中的日期格式
	// 將 ISO 格式（如 2026-01-01T00:00:00Z）轉換為 YYYY-MM-DD 格式
	log.Println("Normalizing effective_range date format...")
	if err := db.normalizeEffectiveRangeDateFormat(); err != nil {
		log.Printf("Warning: Failed to normalize effective_range dates: %v", err)
	} else {
		log.Println("Effective range date format normalized successfully")
	}

	// 1. 新增 start_date_generated 欄位（從 effective_range 提取）
	if err := db.addGeneratedColumnIfNotExists(
		"schedule_rules",
		"start_date_generated",
		"DATE GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date'))) STORED",
	); err != nil {
		log.Printf("Warning: Failed to add start_date_generated column: %v", err)
	}

	// 2. 新增 end_date_generated 欄位（從 effective_range 提取）
	if err := db.addGeneratedColumnIfNotExists(
		"schedule_rules",
		"end_date_generated",
		"DATE GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date'))) STORED",
	); err != nil {
		log.Printf("Warning: Failed to add end_date_generated column: %v", err)
	}

	// 3. 為 start_date_generated 建立索引
	if err := db.addIndexIfNotExists("schedule_rules", "idx_rule_start_date", "start_date_generated"); err != nil {
		log.Printf("Warning: Failed to add index on start_date_generated: %v", err)
	}

	// 4. 為 end_date_generated 建立索引
	if err := db.addIndexIfNotExists("schedule_rules", "idx_rule_end_date", "end_date_generated"); err != nil {
		log.Printf("Warning: Failed to add index on end_date_generated: %v", err)
	}

	// 5. 建立複合索引用於日期範圍查詢
	if err := db.addIndexIfNotExists("schedule_rules", "idx_rule_date_range", "start_date_generated, end_date_generated"); err != nil {
		log.Printf("Warning: Failed to add index on date range: %v", err)
	}

	log.Println("Schedule rules generated columns migration completed")
}

// normalizeEffectiveRangeDateFormat 正規化 effective_range 中的日期格式
// 將 ISO 8601 格式（如 2026-01-01T00:00:00Z）轉換為 MySQL DATE 格式（2026-01-01）
func (db *DB) normalizeEffectiveRangeDateFormat() error {
	// 使用 SUBSTRING 提取 YYYY-MM-DD 部分
	// 對於 ISO 格式如 "2026-01-01T00:00:00Z"，取前 10 個字元
	updateSQL := `
		UPDATE schedule_rules
		SET effective_range = JSON_SET(
			effective_range,
			'$.start_date',
			SUBSTRING(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')), 1, 10),
			'$.end_date',
			SUBSTRING(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')), 1, 10)
		)
		WHERE effective_range IS NOT NULL
		AND (JSON_EXTRACT(effective_range, '$.start_date') IS NOT NULL
			OR JSON_EXTRACT(effective_range, '$.end_date') IS NOT NULL)
	`
	return db.WDB.Exec(updateSQL).Error
}

// addGeneratedColumnIfNotExists 檢查並新增 Generated Column 欄位
func (db *DB) addGeneratedColumnIfNotExists(tableName, columnName, columnDef string) error {
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
		log.Printf("Generated column %s.%s already exists, skipping", tableName, columnName)
		return nil
	}

	// 新增 Generated Column 欄位
	addSQL := "ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + columnDef
	if err := db.WDB.Exec(addSQL).Error; err != nil {
		return err
	}

	log.Printf("Added generated column %s.%s", tableName, columnName)
	return nil
}

// addIndexIfNotExists 檢查並新增索引
func (db *DB) addIndexIfNotExists(tableName, indexName, columnName string) error {
	// 檢查索引是否存在
	var count int
	result := db.WDB.Raw(`
		SELECT COUNT(*) INTO @a
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
		AND INDEX_NAME = ?
	`, tableName, indexName).Scan(&count)

	if result.Error != nil {
		return result.Error
	}

	if count > 0 {
		log.Printf("Index %s on %s already exists, skipping", indexName, tableName)
		return nil
	}

	// 新增索引
	addSQL := "ALTER TABLE " + tableName + " ADD INDEX " + indexName + " (" + columnName + ")"
	if err := db.WDB.Exec(addSQL).Error; err != nil {
		return err
	}

	log.Printf("Added index %s on %s(%s)", indexName, tableName, columnName)
	return nil
}
