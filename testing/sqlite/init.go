package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Initialize() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// // 流程測試需要的Mock資料表
	// models := []any{
	// 	models.User{},
	// }

	// // 自動建立 schema
	// if len(models) > 0 {
	// 	if err := db.AutoMigrate(models...); err != nil {
	// 		return nil, err
	// 	}
	// }

	return db, nil
}
