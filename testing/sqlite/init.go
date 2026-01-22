package sqlite

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Initialize() (*gorm.DB, error) {
	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("mysql init error: %w", err)
	}
	return db, nil
}
