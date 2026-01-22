package test

import (
	"fmt"
	"testing"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMySQLTestDB() (*gorm.DB, error) {
	dsn := "root:timeledger_root_2026@tcp(localhost:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("MySQL connection failed: %v", err)
	}
	return db, nil
}

func TestCenterRepository_CRUD_WithMySQL(t *testing.T) {
	t.Skip("Skipping - requires .env file and production DB")
}
