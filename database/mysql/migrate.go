package mysql

import (
	"fmt"
	"log"
	"timeLedger/app/models"
)

// AutoMigrate 自動建立資料表
func (db *DB) AutoMigrate() {
	if err := db.WDB.AutoMigrate(
		&models.User{},
	); err != nil {
		panic(fmt.Errorf("MySQL autoMigrate failed: %v", err))
	}
	log.Println("AutoMigrate done")
}
