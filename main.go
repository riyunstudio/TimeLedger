package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"timeLedger/app"
	"timeLedger/app/console"
	"timeLedger/app/servers"
	"timeLedger/app/services"
	"timeLedger/global/logger"

	_ "timeLedger/docs"

	_ "github.com/joho/godotenv/autoload"
)

// @title timeLedger 阿卡莉 模板框架
// @version 1.0
// @description API 維護文件
// @schemes http
func main() {
	// recover 防止因服務 panic 直接關閉
	defer func() {
		if err := recover(); err != nil {
			zapLog := logger.GetLogger()
			zapLog.Errorw("Main panic recovered",
				"error", err,
				"stack", fmt.Sprintf("%v", err),
			)
		}
	}()

	// 初始化日誌器
	logConfig := logger.DefaultConfig()
	logConfig.Level = getLogLevel()
	logConfig.Format = "json"
	if _, err := logger.Initialize(logConfig); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	zapLog := logger.GetLogger()
	defer zapLog.Sync()

	zapLog.Info("Starting TimeLedger application...")

	ctx, cancel := context.WithCancel(context.Background())

	// 初始化 App
	appInstance := app.Initialize()

	// 排程（輕量級，佔用極少資源）
	scheduler := console.Initialize(appInstance)
	scheduler.Start()

	// 通知佇列 Worker（按環境變量開關，預設關閉以節省資源）
	if os.Getenv("NOTIFICATION_WORKER_ENABLED") == "true" {
		go startNotificationWorker(appInstance, ctx, zapLog)
	} else {
		zapLog.Info("Notification worker disabled (set NOTIFICATION_WORKER_ENABLED=true to enable)")
	}

	// 啟動 API server（主要服務）
	gin := servers.Initialize(appInstance)
	gin.Start()

	// 優雅退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go func() {
		<-signals
		cancel()
		scheduler.Stop()
		gin.Stop()
		zapLog.Info("Application shutting down gracefully")
		exit <- struct{}{}
	}()
	<-ctx.Done()
	<-exit
	zapLog.Info("Application exited")
}

// getLogLevel 從環境變數取得日誌級別
func getLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return "info"
	}
	return level
}

// startNotificationWorker 啟動通知佇列背景 worker
func startNotificationWorker(appInstance *app.App, ctx context.Context, zapLog *logger.Logger) {
	zapLog.Info("Starting notification queue worker...")

	queueService := services.NewNotificationQueueService(appInstance)

	// 檢查 Redis 連線
	redisQueue := services.NewRedisQueueService(appInstance)
	if !redisQueue.IsHealthy(context.Background()) {
		zapLog.Warn("Redis not connected, notification queue worker will not start")
		return
	}

	zapLog.Info("Notification queue worker started, listening on Redis queue...")

	// 定時打印統計
	go func() {
		ticker := time.NewTicker(60 * time.Second) // 每分鐘打印一次
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stats := queueService.GetQueueStats(context.Background())
				zapLog.Infow("Notification queue stats",
					"pending", stats["pending"],
					"retry", stats["retry"],
					"total", stats["total"],
					"retried", stats["retried"],
					"failed", stats["failed"],
				)
			}
		}
	}()

	// 持續處理佇列
	for {
		select {
		case <-ctx.Done():
			zapLog.Info("Notification worker stopped")
			return
		default:
			if err := queueService.ProcessQueue(ctx); err != nil {
				zapLog.Errorw("Queue processing error", "error", err)
				time.Sleep(1 * time.Second) // 避免 busy loop
			}
		}
	}
}
