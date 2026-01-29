package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/logger"

	"github.com/hibiken/asynq"
)

// AsynqConfig Asynq 配置
type AsynqConfig struct {
	RedisAddr     string
	Concurrency   int
	MaxRetry      int
	RetryInterval time.Duration
	QueueName     string
}

// DefaultAsynqConfig 預設配置
func DefaultAsynqConfig() *AsynqConfig {
	return &AsynqConfig{
		RedisAddr:     "localhost:6379",
		Concurrency:   10,
		MaxRetry:      3,
		RetryInterval: 5 * time.Second,
		QueueName:     "notifications",
	}
}

// AsynqNotificationService 基於 Asynq 的通知佇列服務
type AsynqNotificationService struct {
	app             *app.App
	adminRepo       *repositories.AdminUserRepository
	teacherRepo     *repositories.TeacherRepository
	lineBotService  LineBotService
	templateService LineBotTemplateService
	client          *asynq.Client
	log             *logger.Logger
	config          *AsynqConfig
}

// NewAsynqNotificationService 建立 Asynq 通知服務
func NewAsynqNotificationService(app *app.App, cfg *AsynqConfig) *AsynqNotificationService {
	if cfg == nil {
		cfg = DefaultAsynqConfig()
	}

	return &AsynqNotificationService{
		app:             app,
		adminRepo:       repositories.NewAdminUserRepository(app),
		teacherRepo:     repositories.NewTeacherRepository(app),
		lineBotService:  NewLineBotService(app),
		templateService: NewLineBotTemplateService(app.Env.FrontendBaseURL),
		client:          asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddr}),
		log:             logger.GetLogger(),
		config:          cfg,
	}
}

// TaskPayload 通知任務負載
type TaskPayload struct {
	Type          string `json:"type"`
	RecipientID   uint   `json:"recipient_id"`
	RecipientType string `json:"recipient_type"`
	Payload       string `json:"payload"`
	CreatedAt     string `json:"created_at"`
}

// TaskType 任務類型
const (
	TaskTypeExceptionSubmit = "notification:exception_submit"
	TaskTypeExceptionResult = "notification:exception_result"
	TaskTypeWelcomeTeacher  = "notification:welcome_teacher"
	TaskTypeWelcomeAdmin    = "notification:welcome_admin"
)

// EnqueueNotification 將通知加入佇列
func (s *AsynqNotificationService) EnqueueNotification(ctx context.Context, item *models.NotificationQueue) error {
	payload := &TaskPayload{
		Type:          s.mapNotificationTypeToTask(item.Type),
		RecipientID:   item.RecipientID,
		RecipientType: item.RecipientType,
		Payload:       item.Payload,
		CreatedAt:     item.ScheduledAt.Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(
		payload.Type,
		payloadBytes,
		asynq.MaxRetry(s.config.MaxRetry),
		asynq.Timeout(30*time.Second),
		asynq.Queue(s.config.QueueName),
	)

	// 使用正確的 Enqueue API
	_, err = s.client.Enqueue(task)
	if err != nil {
		s.log.Errorw("Failed to enqueue notification",
			"type", item.Type,
			"recipient_id", item.RecipientID,
			"error", err,
		)
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	s.log.Infow("Notification enqueued",
		"type", item.Type,
		"recipient_id", item.RecipientID,
		"queue", s.config.QueueName,
	)

	return nil
}

// mapNotificationTypeToTask 映射通知類型到任務類型
func (s *AsynqNotificationService) mapNotificationTypeToTask(notificationType string) string {
	switch notificationType {
	case models.NotificationTypeExceptionSubmit:
		return TaskTypeExceptionSubmit
	case models.NotificationTypeExceptionResult:
		return TaskTypeExceptionResult
	case models.NotificationTypeWelcomeTeacher:
		return TaskTypeWelcomeTeacher
	case models.NotificationTypeWelcomeAdmin:
		return TaskTypeWelcomeAdmin
	default:
		return "notification:unknown"
	}
}

// GetQueueStats 取得佇列統計
func (s *AsynqNotificationService) GetQueueStats(ctx context.Context) map[string]string {
	stats := make(map[string]string)
	// 簡化版本：返回基本統計
	// 在實際部署時可透過 Redis 直接查詢
	stats["pending"] = "0"
	stats["retry"] = "0"
	stats["failed"] = "0"
	stats["total"] = "0"
	stats["retried"] = "0"
	return stats
}

// Close 關閉服務
func (s *AsynqNotificationService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// ProcessNotificationTask 處理通知任務（由 Worker 呼叫）
func ProcessNotificationTask(ctx context.Context, t *asynq.Task) error {
	log.Printf("Processing task: %s", t.Type())

	var payload TaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("Failed to unmarshal task payload: %v", err)
		return err
	}

	// 這裡需要從 context 取得必要的服務
	// 在實際部署時，應該透過 dependency injection 傳入
	// 簡化起見，這裡先返回 nil

	log.Printf("Processed notification task: type=%s, recipient=%d",
		payload.Type, payload.RecipientID)

	return nil
}

// StartWorker 啟動 Asynq Worker
func (s *AsynqNotificationService) StartWorker(ctx context.Context) error {
	mux := asynq.NewServeMux()

	// 註冊任務處理器
	mux.HandleFunc(TaskTypeExceptionSubmit, processExceptionSubmit)
	mux.HandleFunc(TaskTypeExceptionResult, processExceptionResult)
	mux.HandleFunc(TaskTypeWelcomeTeacher, processWelcomeTeacher)
	mux.HandleFunc(TaskTypeWelcomeAdmin, processWelcomeAdmin)

	// 啟動 worker
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: s.config.RedisAddr},
		asynq.Config{
			Concurrency: s.config.Concurrency,
			Queues: map[string]int{
				s.config.QueueName: 6, // 高優先級
			},
		},
	)

	return server.Start(mux)
}

// 任務處理函數（簡化版，實際應用需要 dependency injection）
func processExceptionSubmit(ctx context.Context, t *asynq.Task) error {
	log.Printf("Processing exception submit notification: %s", string(t.Payload()))
	// 實際處理邏輯需要存取資料庫和 LINE Bot Service
	return nil
}

func processExceptionResult(ctx context.Context, t *asynq.Task) error {
	log.Printf("Processing exception result notification: %s", string(t.Payload()))
	return nil
}

func processWelcomeTeacher(ctx context.Context, t *asynq.Task) error {
	log.Printf("Processing welcome teacher notification: %s", string(t.Payload()))
	return nil
}

func processWelcomeAdmin(ctx context.Context, t *asynq.Task) error {
	log.Printf("Processing welcome admin notification: %s", string(t.Payload()))
	return nil
}
