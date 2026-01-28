package test

import (
	"context"
	"fmt"
	"testing"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	mockRedis "timeLedger/testing/redis"

	"github.com/gin-gonic/gin"
)

func setupRateLimiterTestApp() (*app.App, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	testApp := &app.App{
		Env:   nil,
		Err:   e,
		Tools: tool,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	_ = mr
	return testApp, mysqlDB
}

// TestRateLimiter 測試速率限制功能
func TestRateLimiter(t *testing.T) {
	testApp, _ := setupRateLimiterTestApp()
	ctx := context.Background()

	rateLimiter := services.NewRateLimiterService(testApp)

	t.Run("CheckRateLimit_FirstRequest", func(t *testing.T) {
		ip := "192.168.1.100"

		allowed, remaining, _, err := rateLimiter.CheckRateLimit(ctx, ip)
		if err != nil {
			t.Fatalf("CheckRateLimit failed: %v", err)
		}

		if !allowed {
			t.Error("First request should be allowed")
		}

		if remaining <= 0 {
			t.Error("Remaining requests should be positive for first request")
		}
	})

	t.Run("CheckRateLimit_MultipleRequests", func(t *testing.T) {
		ip := "192.168.1.101"

		// 發送多個請求
		maxRequests := 5 // 使用較小的值進行測試
		for i := 0; i < maxRequests; i++ {
			err := rateLimiter.RecordRequest(ctx, ip)
			if err != nil {
				t.Fatalf("RecordRequest failed: %v", err)
			}
		}

		_, remaining, _, err := rateLimiter.CheckRateLimit(ctx, ip)
		if err != nil {
			t.Fatalf("CheckRateLimit failed: %v", err)
		}

		// 應該還允許一些請求
		if remaining < 0 {
			t.Error("Should have remaining requests")
		}
	})

	t.Run("ResetIP", func(t *testing.T) {
		ip := "192.168.1.102"

		// 記錄一些請求
		for i := 0; i < 3; i++ {
			_ = rateLimiter.RecordRequest(ctx, ip)
		}

		// 重置 IP
		err := rateLimiter.ResetIP(ctx, ip)
		if err != nil {
			t.Fatalf("ResetIP failed: %v", err)
		}

		// 重置後應該允許請求
		allowed, _, _, err := rateLimiter.CheckRateLimit(ctx, ip)
		if err != nil {
			t.Fatalf("CheckRateLimit after reset failed: %v", err)
		}

		if !allowed {
			t.Error("Request should be allowed after reset")
		}
	})
}
