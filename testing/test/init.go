package test

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitializeTestDB initializes a test database connection.
// It uses MySQL from .env configuration.
func InitializeTestDB() (*gorm.DB, error) {
	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 優化連接池配置
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// CloseDB closes the database connection safely.
func CloseDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// SetupTestDB creates the test database and tables if they don't exist.
func SetupTestDB() (*gorm.DB, error) {
	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// 創建資料庫（如果不存在）
	createDBSQL := "CREATE DATABASE IF NOT EXISTS timeledger_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	if err := db.Exec(createDBSQL).Error; err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// 斷開並重新連接到目標資料庫
	sqlDB, _ := db.DB()
	sqlDB.Close()

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reconnect to test database: %w", err)
	}

	return db, nil
}

// GetTestDB returns a database connection, creating it if necessary.
// Useful for long-running test processes.
var testDB *gorm.DB

func GetTestDB() (*gorm.DB, error) {
	if testDB == nil {
		var err error
		testDB, err = InitializeTestDB()
		if err != nil {
			return nil, err
		}
	}
	return testDB, nil
}

// CloseTestDB closes the cached test database connection.
func CloseTestDB() {
	if testDB != nil {
		sqlDB, _ := testDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		testDB = nil
	}
}

// CleanupTestData cleans up test data from the database.
// Call this after tests to remove test records.
func CleanupTestData(db *gorm.DB, tableNames []string) error {
	for _, table := range tableNames {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id > 0", table)).Error; err != nil {
			return fmt.Errorf("failed to cleanup table %s: %w", table, err)
		}
	}
	return nil
}

// PrintTestResult prints a formatted test result.
func PrintTestResult(name string, passed bool, duration time.Duration) {
	status := "✅ PASS"
	if !passed {
		status = "❌ FAIL"
	}
	fmt.Printf("  %s | %s | %s\n", status, name, duration.Round(time.Millisecond))
}

// EnsureEnv checks if required environment variables are set.
func EnsureEnv() error {
	required := []string{
		"MYSQL_MASTER_HOST",
		"MYSQL_MASTER_PORT",
		"MYSQL_MASTER_USER",
		"MYSQL_MASTER_PASS",
		"MYSQL_MASTER_NAME",
	}

	missing := []string{}
	for _, env := range required {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing environment variables: %v", missing)
	}

	return nil
}
