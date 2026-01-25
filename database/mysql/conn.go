package mysql

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"timeLedger/configs"
	"timeLedger/global"

	"gitlab.en.mcbwvx.com/frame/zilean/logs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
			env.AppDebug,
		)

		// 從DB - 如果Slave設定與Master相同或為空，則使用同一個連線
		slaveHost := env.MysqlSlaveHost
		if slaveHost == "" || slaveHost == env.MysqlMasterHost {
			// 單資料庫模式：從DB與主DB使用同一個連線
			readDB = writeDB
			log.Println("MySQL: Single database mode (WDB = RDB)")
		} else {
			// 主從模式：使用不同的從DB連線
			readDB = connectDB(
				env.MysqlSlaveHost,
				env.MysqlSlavePort,
				env.MysqlSlaveUser,
				env.MysqlSlavePass,
				env.MysqlSlaveName,
				env.AppDebug,
			)
			log.Println("MySQL Master/Slave connected")
		}
	})

	return &DB{
		WDB: writeDB,
		RDB: readDB,
	}
}

func connectDB(host, port, user, pass, name string, debug bool) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &GormTraceLogger{debug: debug},
	})
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

// 實作 Gorm Log Interface
type GormTraceLogger struct {
	debug bool
}

// LogMode 實現 logger.Interface
func (t *GormTraceLogger) LogMode(level logger.LogLevel) logger.Interface {
	return t
}

func (t *GormTraceLogger) Info(ctx context.Context, msg string, data ...any) {
	// 只有 Debug模式 開啟才紀錄
	if !t.debug {
		return
	}

	// Tid
	traceID := getTraceID(ctx)
	logs.GormLogInit().
		SetEvent(logs.GORM_EVENT_DB_INFO).
		SetTraceID(traceID).
		SetExtraInfo(map[string]any{
			"message": msg,
			"data":    data,
		}).
		PrintInfo("GORM Info")
}

func (t *GormTraceLogger) Warn(ctx context.Context, msg string, data ...any) {
	// 只有 Debug模式 開啟才紀錄
	if !t.debug {
		return
	}

	// 取得 err
	errVal := extractErr(data...)
	if errVal == nil {
		return
	}

	// Tid
	traceID := getTraceID(ctx)
	logs.GormLogInit().
		SetEvent(logs.GORM_EVENT_DB_WARN).
		SetTraceID(traceID).
		SetError(errVal).
		SetExtraInfo(map[string]any{
			"message": msg,
			"data":    data,
		}).
		PrintWarn("GORM Warn")
}

func (t *GormTraceLogger) Error(ctx context.Context, msg string, data ...any) {
	// 只有 Debug模式 開啟才紀錄
	if !t.debug {
		return
	}

	// 取得 err
	errVal := extractErr(data...)
	if errVal == nil {
		return
	}

	// Tid
	traceID := getTraceID(ctx)

	logs.GormLogInit().
		SetEvent(logs.GORM_EVENT_DB_ERROR).
		SetTraceID(traceID).
		SetError(errVal).
		SetExtraInfo(map[string]any{
			"message": msg,
			"data":    data,
		}).
		PrintError("GORM Error")
}

func (t *GormTraceLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc() // 第二個回傳是筆數（因為只紀錄 err sql, 所以不需要存）

	// 只紀錄有錯的 SQL Error
	if err == nil {
		return
	}

	// 計算執行時間
	elapsed := time.Since(begin).Seconds()

	// Tid
	traceID := getTraceID(ctx)

	// 只記錄有錯誤的 SQL
	logs.GormLogInit().
		SetEvent(logs.GORM_EVENT_SQL_ERROR).
		SetSql(sql).
		SetRunTime(elapsed).
		SetTraceID(traceID).
		SetError(err).
		PrintError("SQL execution error")
}

func getTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(global.TraceIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func extractErr(data ...any) error {
	for _, d := range data {
		if err, ok := d.(error); ok {
			return err
		}
	}
	return nil
}
