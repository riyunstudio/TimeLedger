package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"timeLedger/app"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis Key 前綴
	RedisKeyPending = "notification:pending"
	RedisKeyRetry   = "notification:retry"
	RedisKeyStats   = "notification:stats"

	// 最大重試次數
	MaxRetryCount = 3

	// 重試延遲（秒）
	RetryDelay = 5
)

// RedisQueueService Redis 佇列服務
type RedisQueueService struct {
	app   *app.App
	redis *redis.Client
}

// NewRedisQueueService 建立 Redis Queue Service
func NewRedisQueueService(app *app.App) *RedisQueueService {
	return &RedisQueueService{
		app:   app,
		redis: app.Redis.DB0,
	}
}

// NotificationItem 通知佇列項目
type NotificationItem struct {
	ID            uint      `json:"id"`
	Type          string    `json:"type"`                     // EXCEPTION_SUBMIT, EXCEPTION_RESULT, WELCOME
	RecipientID   uint      `json:"recipient_id"`
	RecipientType string    `json:"recipient_type"`           // ADMIN, TEACHER
	Payload       string    `json:"payload"`                  // JSON 字串
	RetryCount    int       `json:"retry_count"`
	CreatedAt     time.Time `json:"created_at"`
	LastTryAt     *time.Time `json:"last_try_at,omitempty"`
}

// PushNotification 將通知加入佇列
func (s *RedisQueueService) PushNotification(ctx context.Context, item *NotificationItem) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// 使用 LPUSH 加入佇列頭部（最新通知在前面）
	if err := s.redis.LPush(ctx, RedisKeyPending, data).Err(); err != nil {
		return fmt.Errorf("failed to push to redis: %w", err)
	}

	// 更新統計
	s.IncrementCounter("total")

	fmt.Printf("[INFO] Notification queued: type=%s, recipient=%d\n", 
		item.Type, item.RecipientID)

	return nil
}

// PopNotification 從佇列取出通知（Blocking）
func (s *RedisQueueService) PopNotification(ctx context.Context) (*NotificationItem, error) {
	// 使用 BRPOP 阻塞式取出，0 表示無超時
	result, err := s.redis.BRPop(ctx, 0, RedisKeyPending).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to pop from redis: %w", err)
	}

	if len(result) < 2 {
		return nil, nil
	}

	var item NotificationItem
	if err := json.Unmarshal([]byte(result[1]), &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification: %w", err)
	}

	// 標記處理時間
	now := time.Now()
	item.LastTryAt = &now

	return &item, nil
}

// PushToRetry 將通知加入延遲重試佇列
func (s *RedisQueueService) PushToRetry(ctx context.Context, item *NotificationItem) error {
	item.RetryCount++

	if item.RetryCount >= MaxRetryCount {
		// 超過最大重試次數，記錄失敗
		s.IncrementCounter("failed")
		fmt.Printf("[WARN] Notification %d exceeded max retries, discarded\n", item.ID)
		return nil
	}

	// 使用 sorted set 以時間戳排序，實現延遲重試
	data, _ := json.Marshal(item)
	score := float64(time.Now().Add(time.Second * time.Duration(RetryDelay*item.RetryCount)).Unix())
	
	if err := s.redis.ZAdd(ctx, RedisKeyRetry, redis.Z{
		Score:  score,
		Member: string(data),
	}).Err(); err != nil {
		return fmt.Errorf("failed to add to retry queue: %w", err)
	}

	s.IncrementCounter("retried")
	fmt.Printf("[INFO] Notification %d queued for retry (attempt %d)\n", 
		item.ID, item.RetryCount)

	return nil
}

// ProcessRetryQueue 處理延遲重試佇列
func (s *RedisQueueService) ProcessRetryQueue(ctx context.Context) error {
	// 取得已到期的重試項目
	now := float64(time.Now().Unix())
	
	result, err := s.redis.ZRangeByScore(ctx, RedisKeyRetry, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
		Count: 10,
	}).Result()
	
	if err != nil {
		return fmt.Errorf("failed to get retry queue: %w", err)
	}

	for _, data := range result {
		var item NotificationItem
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			continue
		}

		// 移出重試佇列
		s.redis.ZRem(ctx, RedisKeyRetry, data)

		// 重新加入待處理佇列
		item.RetryCount--
		s.PushNotification(ctx, &item)
	}

	return nil
}

// GetQueueLength 取得佇列長度
func (s *RedisQueueService) GetQueueLength(ctx context.Context) (int64, error) {
	return s.redis.LLen(ctx, RedisKeyPending).Result()
}

// GetStats 取得統計資訊
func (s *RedisQueueService) GetStats(ctx context.Context) map[string]string {
	stats := make(map[string]string)
	
	pending, _ := s.redis.LLen(ctx, RedisKeyPending).Result()
	retry, _ := s.redis.ZCard(ctx, RedisKeyRetry).Result()
	
	stats["pending"] = fmt.Sprintf("%d", pending)
	stats["retry"] = fmt.Sprintf("%d", retry)
	
	// 從 Hash 取得計數器
	hgetall, _ := s.redis.HGetAll(ctx, RedisKeyStats).Result()
	for k, v := range hgetall {
		stats[k] = v
	}

	return stats
}

// IncrementCounter 增加計數器
func (s *RedisQueueService) IncrementCounter(counter string) {
	s.redis.HIncrBy(context.Background(), RedisKeyStats, counter, 1)
}

// IsHealthy 檢查 Redis 連線是否正常
func (s *RedisQueueService) IsHealthy(ctx context.Context) bool {
	return s.redis.Ping(ctx).Err() == nil
}
