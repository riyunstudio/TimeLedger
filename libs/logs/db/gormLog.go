package db

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type GormLog struct {
	service   string
	podName   string
	topic     string  // 主題或操作名稱
	event     string  // 事件名稱
	sql       string  // 語法
	runTime   float64 // 執行時間
	err       string  // 錯誤訊息
	extraInfo any     // 額外資訊
	traceID   string
}

func GormLogInit() *GormLog {
	return &GormLog{
		topic:   "Mysql",
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

const EVENT_DB_INFO = "db_info"
const EVENT_DB_WARN = "db_warn"
const EVENT_DB_ERROR = "db_error"
const EVENT_SQL_ERROR = "sql_error"
const EVENT_SQL_SLOW = "sql_slow"

func (l *GormLog) SetEvent(event string) *GormLog {
	l.event = event
	return l
}

func (l *GormLog) SetSql(sql string) *GormLog {
	l.sql = sql
	return l
}

func (l *GormLog) SetRunTime(runTime float64) *GormLog {
	l.runTime = runTime
	return l
}

func (l *GormLog) SetError(err error) *GormLog {
	if err != nil {
		l.err = err.Error()
	}
	return l
}

func (l *GormLog) SetExtraInfo(info any) *GormLog {
	l.extraInfo = info
	return l
}

func (l *GormLog) SetTraceID(traceID string) *GormLog {
	l.traceID = traceID
	return l
}

// 生成 logrus fields
func (l *GormLog) fields() logrus.Fields {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	f := logrus.Fields{}
	if l.service != "" {
		f["service"] = l.service
	}
	if l.podName != "" {
		f["podName"] = l.podName
	}
	if l.topic != "" {
		f["topic"] = l.topic
	}
	if l.event != "" {
		f["event"] = l.event
	}
	if l.sql != "" {
		f["sql"] = l.sql
	}
	if l.runTime != 0 {
		f["runTime"] = l.runTime
	}
	if l.err != "" {
		f["error"] = l.err
	}
	if l.extraInfo != nil {
		f["extraInfo"] = l.extraInfo
	}
	if l.traceID != "" {
		f["traceID"] = l.traceID
	}
	return f
}

// Info log
func (l *GormLog) PrintInfo(msg string) {
	logrus.WithFields(l.fields()).Info(msg)
}

// Warn log
func (l *GormLog) PrintWarn(msg string) {
	logrus.WithFields(l.fields()).Warn(msg)
}

// Error log
func (l *GormLog) PrintError(msg string) {
	logrus.WithFields(l.fields()).Error(msg)
}
