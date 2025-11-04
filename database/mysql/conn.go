package mysql

import (
	"akali/configs"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	WDB *gorm.DB // 寫入用
	RDB *gorm.DB // 讀取用
}

// 初始化主從資料庫
func Initialize(env *configs.Env) *DB {
	var (
		writeDB *gorm.DB
		readDB  *gorm.DB
		once    sync.Once
	)

	once.Do(func() {
		// 主DB
		writeDB = connectDB(
			env.MysqlMasterHost,
			env.MysqlMasterPort,
			env.MysqlMasterUser,
			env.MysqlMasterPass,
			env.MysqlMasterName,
		)

		// 從DB
		readDB = connectDB(
			env.MysqlSlaveHost,
			env.MysqlSlavePort,
			env.MysqlSlaveUser,
			env.MysqlSlavePass,
			env.MysqlSlaveName,
		)

		log.Println("MySQL Master/Slave connected")
	})

	return &DB{
		WDB: writeDB,
		RDB: readDB,
	}
}

func connectDB(host, port, user, pass, name string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("MySQL connection failed: %v", err))
	}

	// 設定寫入的連線池設定
	sqlDB, err := db.DB()

	if err != nil {
		panic(fmt.Errorf("MySQL set pool failed: %v", err))
	}
	// 設置閒置連線上限
	sqlDB.SetMaxIdleConns(50)
	// 設置開放連線上限
	sqlDB.SetMaxOpenConns(100)
	// 設置連線的最大閒置時間為 5 分鐘
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}
