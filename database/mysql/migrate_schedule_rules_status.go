package mysql

import (
	"log"
	"time"

	"timeLedger/app/models"
)

// MigrateScheduleRulesStatus 遷移 schedule_rules 表的 status 欄位
// 確保所有現有的 schedule_rules 記錄都有正確的 status 值
func (db *DB) MigrateScheduleRulesStatus() {
	log.Println("Starting schedule_rules.status migration...")

	// 1. 檢查 status 欄位是否存在
	var columnExists bool
	checkSQL := `
		SELECT COUNT(*) INTO @a
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = 'schedule_rules'
		AND COLUMN_NAME = 'status'
	`
	if err := db.WDB.Raw(checkSQL).Scan(&columnExists).Error; err != nil {
		log.Printf("Warning: Failed to check status column existence: %v", err)
		return
	}

	if !columnExists {
		// 欄位不存在，先新增欄位
		addColumnSQL := `
			ALTER TABLE schedule_rules
			ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'CONFIRMED'
			AFTER suspended_dates
		`
		if err := db.WDB.Exec(addColumnSQL).Error; err != nil {
			log.Printf("Warning: Failed to add status column: %v", err)
			return
		}
		log.Println("Added status column to schedule_rules table")
	} else {
		log.Println("status column already exists in schedule_rules table")
	}

	// 2. 檢查是否有 NULL 或空的 status 值需要修復
	var nullOrEmptyCount int64
	nullCheckSQL := `
		SELECT COUNT(*) INTO @a
		FROM schedule_rules
		WHERE status IS NULL OR status = ''
	`
	if err := db.WDB.Raw(nullCheckSQL).Scan(&nullOrEmptyCount).Error; err != nil {
		log.Printf("Warning: Failed to check NULL/empty status values: %v", err)
	} else if nullOrEmptyCount > 0 {
		// 修復 NULL 或空的值為 CONFIRMED
		fixSQL := `
			UPDATE schedule_rules
			SET status = 'CONFIRMED'
			WHERE status IS NULL OR status = ''
		`
		if err := db.WDB.Exec(fixSQL).Error; err != nil {
			log.Printf("Warning: Failed to fix NULL/empty status values: %v", err)
		} else {
			log.Printf("Fixed %d records with NULL or empty status values", nullOrEmptyCount)
		}
	}

	// 3. 驗證所有現有的 status 值都是有效的
	var invalidCount int64
	validStatuses := []string{
		models.RuleStatusPlanned,
		models.RuleStatusConfirmed,
		models.RuleStatusSuspended,
		models.RuleStatusArchived,
	}

	// 建立 IN 子句
	inClause := "'" + validStatuses[0] + "'"
	for i := 1; i < len(validStatuses); i++ {
		inClause += ", '" + validStatuses[i] + "'"
	}

	invalidCheckSQL := `
		SELECT COUNT(*) INTO @a
		FROM schedule_rules
		WHERE status NOT IN (` + inClause + `)
	`
	if err := db.WDB.Raw(invalidCheckSQL).Scan(&invalidCount).Error; err != nil {
		log.Printf("Warning: Failed to check invalid status values: %v", err)
	} else if invalidCount > 0 {
		log.Printf("Found %d records with invalid status values, marking for review", invalidCount)
		// 記錄無效的記錄供管理員審查，但不自動修復
		// 因為這可能是業務邏輯問題
		var invalidRecords []struct {
			ID     uint
			Status  string
			CenterID uint
		}
		selectSQL := `
			SELECT id, status, center_id
			FROM schedule_rules
			WHERE status NOT IN (` + inClause + `)
			LIMIT 100
		`
		if err := db.WDB.Raw(selectSQL).Scan(&invalidRecords).Error; err == nil {
			for _, record := range invalidRecords {
				log.Printf(
					"Invalid status record: id=%d, center_id=%d, status='%s'",
					record.ID, record.CenterID, record.Status,
				)
			}
		}
	}

	// 4. 建立索引以優化狀態查詢（如果尚未存在）
	var indexExists bool
	indexCheckSQL := `
		SELECT COUNT(*) INTO @a
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = 'schedule_rules'
		AND INDEX_NAME = 'idx_schedule_rules_status'
	`
	if err := db.WDB.Raw(indexCheckSQL).Scan(&indexExists).Error; err != nil {
		log.Printf("Warning: Failed to check status index existence: %v", err)
	} else if !indexExists {
		addIndexSQL := `
			CREATE INDEX idx_schedule_rules_status
			ON schedule_rules (status)
		`
		if err := db.WDB.Exec(addIndexSQL).Error; err != nil {
			log.Printf("Warning: Failed to create status index: %v", err)
		} else {
			log.Println("Created index idx_schedule_rules_status on status column")
		}
	}

	log.Println("schedule_rules.status migration completed successfully")
}

// MigrateScheduleRulesStatusWithTimeout 带超时的迁移版本
// 用於在系統啟動時安全執行遷移
func (db *DB) MigrateScheduleRulesStatusWithTimeout(timeout time.Duration) {
	done := make(chan bool)

	go func() {
		db.MigrateScheduleRulesStatus()
		done <- true
	}()

	select {
	case <-done:
		log.Println("Status migration completed within timeout")
	case <-time.After(timeout):
		log.Println("Warning: Status migration timed out, continuing...")
	}
}
