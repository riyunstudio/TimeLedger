package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日誌器封裝
type Logger struct {
	sugar  *zap.SugaredLogger
	logger *zap.Logger
}

// 全局日誌器實例
var globalLogger *Logger
var once sync.Once

// Config 日誌配置
type Config struct {
	Level      string `json:"level"`      // debug, info, warn, error
	Format     string `json:"format"`     // json, console
	OutputPath string `json:"outputPath"` // 日誌輸出目錄
	Filename   string `json:"filename"`   // 日誌檔案名稱
	MaxSize    int    `json:"maxSize"`    // 日誌檔案最大大小 (MB)
	MaxAge     int    `json:"maxAge"`     // 日誌保留天數
	MaxBackups int    `json:"maxBackups"` // 最大備份數量
}

// DefaultConfig 預設配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Format:     "json",
		OutputPath: "logs",
		Filename:   "timeledger.log",
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 10,
	}
}

// Initialize 初始化全局日誌器
func Initialize(cfg *Config) (*Logger, error) {
	once.Do(func() {
		var err error
		globalLogger, err = NewLogger(cfg)
		if err != nil {
			panic(err)
		}
	})
	return globalLogger, nil
}

// GetLogger 取得全局日誌器
func GetLogger() *Logger {
	if globalLogger == nil {
		panic("logger not initialized, call Initialize() first")
	}
	return globalLogger
}

// NewLogger 建立新的日誌器實例
func NewLogger(cfg *Config) (*Logger, error) {
	// 解析日誌級別
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 設定日誌編碼器
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "console" {
		// 控制台格式（人類可讀）
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
	} else {
		// JSON 格式（ELK 分析用）
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
	}

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 設定寫入器
	var cores []zapcore.Core

	// 只有當 OutputPath 不為空時才建立檔案寫入器
	if cfg.OutputPath != "" {
		// 確保輸出目錄存在
		if err := os.MkdirAll(cfg.OutputPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// 開啟日誌檔案
		logFilePath := filepath.Join(cfg.OutputPath, cfg.Filename)
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer := zapcore.AddSync(logFile)

		// 檔案寫入器
		core := zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(writer),
			level,
		)
		cores = append(cores, core)
	}

	// 輸出到控制台（開發環境或無檔案輸出時）
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		level,
	)
	cores = append(cores, core)

	// 建立日誌器
	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{
		sugar:  logger.Sugar(),
		logger: logger,
	}, nil
}

// Sync 同步日誌（程式結束前呼叫）
func (l *Logger) Sync() {
	if l != nil {
		l.sugar.Sync()
		l.logger.Sync()
	}
}

// Debug  Debug 級別日誌
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debugf Debug 級別日誌（格式化）
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Debugw Debug 級別日誌（結構化）
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Info Info 級別日誌
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Infof Info 級別日誌（格式化）
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Infow Info 級別日誌（結構化）
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warn Warn 級別日誌
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Warnf Warn 級別日誌（格式化）
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Warnw Warn 級別日誌（結構化）
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Error Error 級別日誌
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Errorf Error 級別日誌（格式化）
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Errorw Error 級別日誌（結構化）
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatal Fatal 級別日誌
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Fatalf Fatal 級別日誌（格式化）
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

// Fatalw Fatal 級別日誌（結構化）
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// With 建立帶有上下文的新日誌器
func (l *Logger) With(keysAndValues ...interface{}) *zap.SugaredLogger {
	return l.sugar.With(keysAndValues...)
}

// ForComponent 為特定元件建立日誌器
func (l *Logger) ForComponent(component string) *zap.SugaredLogger {
	return l.sugar.With("component", component)
}
