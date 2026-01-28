package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	app    *app.App
	config RateLimitConfig
}

type RateLimitConfig struct {
	Enabled       bool          // 是否啟用速率限制
	RequestsPerIP int           // 每個 IP 的請求上限
	WindowSize    time.Duration // 時間窗口大小
	BlockDuration time.Duration // 封鎖時間（當超過限制時）
}

type RateLimiterService interface {
	CheckRateLimit(ctx context.Context, ip string) (allowed bool, remaining int, resetAt time.Time, err error)
	RecordRequest(ctx context.Context, ip string) error
	IsBlocked(ctx context.Context, ip string) (bool, time.Duration, error)
	ResetIP(ctx context.Context, ip string) error
}

type RedisRateLimiter struct {
	app    *app.App
	config RateLimitConfig
}

func NewRateLimiterService(app *app.App) RateLimiterService {
	// 預設配置
	config := RateLimitConfig{
		Enabled:       true,
		RequestsPerIP: 100,
		WindowSize:    time.Minute,
		BlockDuration: time.Minute * 5,
	}

	// 從環境變數讀取配置（如果可用）
	if app != nil && app.Env != nil {
		config.Enabled = app.Env.RateLimitEnabled
		config.RequestsPerIP = app.Env.RateLimitRequests

		// 從環境變數讀取配置（覆蓋預設值）
		if app.Env.RateLimitWindow != "" {
			if window, err := time.ParseDuration(app.Env.RateLimitWindow); err == nil {
				config.WindowSize = window
			}
		}
		if app.Env.RateLimitBlockDuration != "" {
			if block, err := time.ParseDuration(app.Env.RateLimitBlockDuration); err == nil {
				config.BlockDuration = block
			}
		}
	}

	return &RedisRateLimiter{
		app:    app,
		config: config,
	}
}

// CheckRateLimit 檢查是否允許請求
func (r *RedisRateLimiter) CheckRateLimit(ctx context.Context, ip string) (allowed bool, remaining int, resetAt time.Time, err error) {
	if !r.config.Enabled {
		return true, r.config.RequestsPerIP, time.Now().Add(r.config.WindowSize), nil
	}

	// 檢查是否被封鎖
	blocked, _, err := r.IsBlocked(ctx, ip)
	if err != nil {
		return false, 0, time.Time{}, err
	}
	if blocked {
		return false, 0, time.Now().Add(r.config.BlockDuration), nil
	}

	// 使用 Redis 滑動窗口計數器
	key := fmt.Sprintf("ratelimit:%s", ip)

	now := time.Now()
	windowStart := now.Add(-r.config.WindowSize)
	pipe := r.app.Redis.DB0.Pipeline()

	// 移除窗口外的舊請求
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixMilli(), 10))

	// 計算當前請求數
	countCmd := pipe.ZCard(ctx, key)

	// 執行pipeline
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, 0, time.Time{}, err
	}

	count := countCmd.Val()
	remaining = r.config.RequestsPerIP - int(count) - 1

	if count >= int64(r.config.RequestsPerIP) {
		return false, 0, now.Add(r.config.WindowSize), nil
	}

	return true, remaining, now.Add(r.config.WindowSize), nil
}

// RecordRequest 記錄請求
func (r *RedisRateLimiter) RecordRequest(ctx context.Context, ip string) error {
	if !r.config.Enabled {
		return nil
	}

	key := fmt.Sprintf("ratelimit:%s", ip)
	now := time.Now()
	member := redis.Z{
		Score:  float64(now.UnixMilli()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	}

	pipe := r.app.Redis.DB0.Pipeline()
	pipe.ZAdd(ctx, key, member)
	pipe.Expire(ctx, key, r.config.BlockDuration) // 設置過期時間為封鎖時間

	_, err := pipe.Exec(ctx)
	return err
}

// IsBlocked 檢查 IP 是否被封鎖
func (r *RedisRateLimiter) IsBlocked(ctx context.Context, ip string) (bool, time.Duration, error) {
	key := fmt.Sprintf("ratelimit:blocked:%s", ip)

	_, err := r.app.Redis.DB0.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	// 取得剩餘封鎖時間
	ttl, err := r.app.Redis.DB0.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	return true, ttl, nil
}

// ResetIP 重置 IP 的計數器
func (r *RedisRateLimiter) ResetIP(ctx context.Context, ip string) error {
	keys := []string{
		fmt.Sprintf("ratelimit:%s", ip),
		fmt.Sprintf("ratelimit:blocked:%s", ip),
	}
	return r.app.Redis.DB0.Del(ctx, keys...).Err()
}

// BlockIP 封鎖 IP
func (r *RedisRateLimiter) BlockIP(ctx context.Context, ip string, duration time.Duration) error {
	key := fmt.Sprintf("ratelimit:blocked:%s", ip)
	return r.app.Redis.DB0.Set(ctx, key, "1", duration).Err()
}

// RateLimitMiddleware 建立速率限制中介層
func RateLimitMiddleware(rateLimiter RateLimiterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 取得客戶端 IP
		ip := c.ClientIP()

		// 排除某些路徑（如健康檢查）
		path := c.Request.URL.Path
		if shouldSkipRateLimit(path) {
			c.Next()
			return
		}

		// 檢查速率限制
		allowed, remaining, resetAt, err := rateLimiter.CheckRateLimit(c.Request.Context(), ip)
		if err != nil {
			// 如果 Redis 錯誤，記錄日誌但允許請求（fail open）
			c.Next()
			return
		}

		// 設置 Rate Limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(remaining+1)) // 包含當前請求
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, remaining)))
		c.Header("X-RateLimit-Reset", resetAt.Format(time.RFC3339))

		if !allowed {
			// 記錄請求（用於計數）
			_ = rateLimiter.RecordRequest(c.Request.Context(), ip)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    errInfos.RATE_LIMIT_EXCEEDED,
				"message": "請求頻率過高，請稍後再試",
				"datas": gin.H{
					"retry_after": int(time.Until(resetAt).Seconds()),
				},
			})
			c.Abort()
			return
		}

		// 記錄請求
		if err := rateLimiter.RecordRequest(c.Request.Context(), ip); err != nil {
			// 記錄失敗不應該阻止請求
		}

		c.Next()
	}
}

// shouldSkipRateLimit 檢查是否應該跳過速率限制
func shouldSkipRateLimit(path string) bool {
	skipPaths := []string{
		"/health",
		"/ready",
		"/metrics",
		"/swagger",
		"/docs",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}
