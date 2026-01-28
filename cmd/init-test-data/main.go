package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 資料庫連線資訊
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?multiStatements=true&charset=utf8mb4&parseTime=True&loc=Local"

	fmt.Println("正在連線到 MySQL 資料庫...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("無法連線到資料庫: %v", err)
	}
	defer db.Close()

	// 測試連線
	if err := db.Ping(); err != nil {
		log.Fatalf("無法 Ping 資料庫: %v", err)
	}
	fmt.Println("✓ 成功連線到資料庫")

	// 讀取 SQL 檔案
	sqlFilePath := "database/mysql/init_test_data.sql"
	fmt.Printf("正在讀取 SQL 檔案: %s...\n", sqlFilePath)

	file, err := os.Open(sqlFilePath)
	if err != nil {
		log.Fatalf("無法讀取 SQL 檔案: %v", err)
	}
	defer file.Close()
	fmt.Println("✓ 成功讀取 SQL 檔案")

	// 解析 SQL 檔案，正確處理註釋
	fmt.Println("正在執行 SQL 指令...")
	
	scanner := bufio.NewScanner(file)
	var statements []string
	var currentStmt strings.Builder
	inBlockComment := false
	inMultiLineInsert := false

	for scanner.Scan() {
		line := scanner.Text()

		// 處理多行 INSERT 語句
		if inMultiLineInsert {
			// 檢查是否結束（沒有結尾括號）
			if !strings.HasSuffix(strings.TrimRight(line, " \t"), ";") {
				currentStmt.WriteString(line)
				currentStmt.WriteString("\n")
				continue
			}
		}

		// 處理行內註釋
		if idx := strings.Index(line, "--"); idx != -1 {
			// 如果在多行 INSERT 中，保留 -- 後面的內容（可能是 VALUES 子句的一部分）
			if !inMultiLineInsert {
				line = line[:idx]
			}
		}

		// 處理區塊註釋
		if idx := strings.Index(line, "/*"); idx != -1 {
			line = line[:idx]
			inBlockComment = true
		}
		if inBlockComment {
			if idx := strings.Index(line, "*/"); idx != -1 {
				line = line[idx+2:]
				inBlockComment = false
			} else {
				continue
			}
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString("\n")

		// 檢查是否為完整的 SQL 語句
		trimmedLine := strings.TrimSpace(line)
		if strings.HasSuffix(trimmedLine, ";") {
			stmt := currentStmt.String()
			stmt = strings.TrimSpace(stmt)
			if stmt != "" && !strings.HasPrefix(stmt, "--") {
				statements = append(statements, stmt)
			}
			currentStmt.Reset()
			inMultiLineInsert = false
		} else {
			// 檢查是否開始多行 INSERT
			if strings.HasPrefix(strings.ToUpper(trimmedLine), "INSERT") {
				inMultiLineInsert = true
			}
		}
	}

	// 執行所有 SQL 語句
	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, stmt := range statements {
		// 跳過純註釋
		trimmed := strings.TrimSpace(stmt)
		if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "/*") {
			skipCount++
			continue
		}

		// 跳過 SELECT 查詢
		if strings.HasPrefix(strings.ToUpper(trimmed), "SELECT ") {
			// 這是結尾的統計查詢，跳過
			skipCount++
			continue
		}

		// 跳過 USE 語句
		if strings.HasPrefix(strings.ToUpper(trimmed), "USE ") {
			skipCount++
			continue
		}

		_, err = db.Exec(stmt)
		if err != nil {
			// 忽略常見的無害錯誤
			errMsg := err.Error()
			ignoreErrors := []string{
				"Duplicate column",
				"already exists",
				"doesn't exist",
				"Unknown table",
				"Duplicate entry",
				"table '",
				"' doesn't exist",
				"Safe query",
			}
			
			shouldIgnore := false
			for _, pattern := range ignoreErrors {
				if strings.Contains(errMsg, pattern) {
					shouldIgnore = true
					break
				}
			}
			
			if shouldIgnore {
				skipCount++
				continue
			}
			
			// 顯示錯誤但不停止
			log.Printf("SQL %d Error: %v", i+1, err)
			errorCount++
		} else {
			successCount++
		}
	}

	fmt.Printf("✓ SQL 執行完成 (成功: %d, 跳過: %d, 錯誤: %d)\n", successCount, skipCount, errorCount)

	// 驗證資料
	fmt.Println("\n正在驗證資料...")
	
	tables := []struct {
		name  string
		query string
	}{
		{"centers", "SELECT COUNT(*) FROM centers"},
		{"admin_users", "SELECT COUNT(*) FROM admin_users"},
		{"teachers", "SELECT COUNT(*) FROM teachers"},
		{"rooms", "SELECT COUNT(*) FROM rooms"},
		{"courses", "SELECT COUNT(*) FROM courses"},
		{"offerings", "SELECT COUNT(*) FROM offerings"},
		{"schedule_rules", "SELECT COUNT(*) FROM schedule_rules"},
		{"personal_events", "SELECT COUNT(*) FROM personal_events"},
		{"schedule_exceptions", "SELECT COUNT(*) FROM schedule_exceptions"},
		{"center_holidays", "SELECT COUNT(*) FROM center_holidays"},
		{"teacher_skills", "SELECT COUNT(*) FROM teacher_skills"},
		{"teacher_certificates", "SELECT COUNT(*) FROM teacher_certificates"},
	}

	fmt.Println("\n========================================")
	fmt.Println("TimeLedger 測試資料初始化完成！")
	fmt.Println("========================================")
	fmt.Println("測試資料統計：")

	allZero := true
	for _, table := range tables {
		var count int
		err := db.QueryRow(table.query).Scan(&count)
		if err != nil {
			log.Printf("Warn: Cannot query %s: %v", table.name, err)
			continue
		}
		if count > 0 {
			allZero = false
		}
		fmt.Printf("- %s: %d筆\n", table.name, count)
	}

	fmt.Println("========================================")

	if allZero {
		fmt.Println("\n⚠️  資料可能未正確初始化")
		fmt.Println("執行指令: go test ./testing/test -v 來確認測試狀態")
	} else {
		fmt.Println("\n✓ 測試資料已準備就緒，可以執行測試了！")
		fmt.Println("執行指令: go test ./testing/test -v")
	}
}
