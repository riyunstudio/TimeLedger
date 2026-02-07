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
	RedisPassword string
	Concurrency   int
	MaxRetry      int
	RetryInterval time.Duration
	QueueName     string
}

// DefaultAsynqConfig 預設配置
func DefaultAsynqConfig() *AsynqConfig {
	return &AsynqConfig{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
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
	svc := &AsynqNotificationService{
		app: app,
		log: logger.GetLogger(),
	}

	if app.Env != nil {
		if cfg == nil {
			// 從環境變數讀取 Redis 配置
			cfg = &AsynqConfig{
				RedisAddr:     fmt.Sprintf("%s:%s", app.Env.RedisHost, app.Env.RedisPort),
				RedisPassword: app.Env.RedisPass,
				Concurrency:   10,
				MaxRetry:      3,
				RetryInterval: 5 * time.Second,
				QueueName:     "notifications",
			}
		}
		svc.config = cfg
		svc.client = asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddr, Password: cfg.RedisPassword})
		svc.templateService = NewLineBotTemplateService(app.Env.FrontendBaseURL)
	}

	if app.MySQL != nil {
		svc.adminRepo = repositories.NewAdminUserRepository(app)
		svc.teacherRepo = repositories.NewTeacherRepository(app)
		svc.lineBotService = NewLineBotService(app)
	}

	return svc
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

// AsynqProcessor Asynq 任務處理器（單例）
var asynqProcessor *AsynqTaskProcessor

// AsynqTaskProcessor Asynq 任務處理器
type AsynqTaskProcessor struct {
	app            *app.App
	adminRepo      *repositories.AdminUserRepository
	teacherRepo    *repositories.TeacherRepository
	lineBotService LineBotService
	log            *logger.Logger
}

// NewAsynqTaskProcessor 建立 Asynq 任務處理器
func NewAsynqTaskProcessor(appInstance *app.App) *AsynqTaskProcessor {
	return &AsynqTaskProcessor{
		app:            appInstance,
		adminRepo:      repositories.NewAdminUserRepository(appInstance),
		teacherRepo:    repositories.NewTeacherRepository(appInstance),
		lineBotService: NewLineBotService(appInstance),
		log:            logger.GetLogger(),
	}
}

// GetAsynqProcessor 取得或建立 Asynq 任務處理器（單例模式）
func GetAsynqProcessor(appInstance *app.App) *AsynqTaskProcessor {
	if asynqProcessor == nil {
		asynqProcessor = NewAsynqTaskProcessor(appInstance)
	}
	return asynqProcessor
}

// ProcessNotificationTask 處理通知任務（由 Worker 呼叫）
func (p *AsynqTaskProcessor) ProcessNotificationTask(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		p.log.Errorw("Failed to unmarshal task payload", "error", err)
		return err
	}

	p.log.Infow("Processing notification task",
		"type", payload.Type,
		"recipient_id", payload.RecipientID,
	)

	// 根據任務類型分發處理
	switch payload.Type {
	case TaskTypeExceptionSubmit:
		return p.processExceptionSubmit(ctx, &payload)
	case TaskTypeExceptionResult:
		return p.processExceptionResult(ctx, &payload)
	case TaskTypeWelcomeTeacher:
		return p.processWelcomeTeacher(ctx, &payload)
	case TaskTypeWelcomeAdmin:
		return p.processWelcomeAdmin(ctx, &payload)
	default:
		p.log.Warnw("Unknown task type", "type", payload.Type)
		return fmt.Errorf("unknown task type: %s", payload.Type)
	}
}

// processExceptionSubmit 處理例外申請通知
func (p *AsynqTaskProcessor) processExceptionSubmit(ctx context.Context, payload *TaskPayload) error {
	// 解析 payload
	var flexPayload map[string]interface{}
	if err := json.Unmarshal([]byte(payload.Payload), &flexPayload); err != nil {
		return fmt.Errorf("failed to parse flex payload: %w", err)
	}

	// 取得管理員的 LINE User ID
	if payload.RecipientType != "ADMIN" {
		return fmt.Errorf("invalid recipient type for exception submit: %s", payload.RecipientType)
	}

	admin, err := p.adminRepo.GetByIDPtr(ctx, payload.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	if admin.LineUserID == "" || !admin.LineNotifyEnabled {
		p.log.Infow("Admin not bound to LINE or notifications disabled",
			"admin_id", payload.RecipientID)
		return nil // 不視為錯誤，只是跳過
	}

	// 發送 LINE 訊息
	altText, _ := flexPayload["altText"].(string)
	contents, _ := flexPayload["contents"].(json.RawMessage)

	if err := p.lineBotService.PushFlexMessage(ctx, admin.LineUserID, altText, contents); err != nil {
		p.log.Errorw("Failed to send exception submit notification",
			"admin_id", payload.RecipientID,
			"error", err)
		return err
	}

	p.log.Infow("Exception submit notification sent successfully",
		"admin_id", payload.RecipientID)
	return nil
}

// processExceptionResult 處理例外審核結果通知
func (p *AsynqTaskProcessor) processExceptionResult(ctx context.Context, payload *TaskPayload) error {
	// 解析 payload
	var flexPayload map[string]interface{}
	if err := json.Unmarshal([]byte(payload.Payload), &flexPayload); err != nil {
		return fmt.Errorf("failed to parse flex payload: %w", err)
	}

	// 取得老師的 LINE User ID
	if payload.RecipientType != "TEACHER" {
		return fmt.Errorf("invalid recipient type for exception result: %s", payload.RecipientType)
	}

	teacher, err := p.teacherRepo.GetByID(ctx, payload.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to get teacher: %w", err)
	}

	if teacher.LineUserID == "" {
		p.log.Infow("Teacher not bound to LINE",
			"teacher_id", payload.RecipientID)
		return nil // 不視為錯誤，只是跳過
	}

	// 發送 LINE 訊息
	altText, _ := flexPayload["altText"].(string)
	contents, _ := flexPayload["contents"].(json.RawMessage)

	if err := p.lineBotService.PushFlexMessage(ctx, teacher.LineUserID, altText, contents); err != nil {
		p.log.Errorw("Failed to send exception result notification",
			"teacher_id", payload.RecipientID,
			"error", err)
		return err
	}

	p.log.Infow("Exception result notification sent successfully",
		"teacher_id", payload.RecipientID)
	return nil
}

// processWelcomeTeacher 處理老師歡迎通知
func (p *AsynqTaskProcessor) processWelcomeTeacher(ctx context.Context, payload *TaskPayload) error {
	// 解析 payload
	var flexPayload map[string]interface{}
	if err := json.Unmarshal([]byte(payload.Payload), &flexPayload); err != nil {
		return fmt.Errorf("failed to parse flex payload: %w", err)
	}

	// 取得老師的 LINE User ID
	if payload.RecipientType != "TEACHER" {
		return fmt.Errorf("invalid recipient type for welcome teacher: %s", payload.RecipientType)
	}

	teacher, err := p.teacherRepo.GetByID(ctx, payload.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to get teacher: %w", err)
	}

	if teacher.LineUserID == "" {
		p.log.Infow("Teacher not bound to LINE",
			"teacher_id", payload.RecipientID)
		return nil // 不視為錯誤，只是跳過
	}

	// 發送 LINE 訊息
	altText, _ := flexPayload["altText"].(string)
	contents, _ := flexPayload["contents"].(json.RawMessage)

	if err := p.lineBotService.PushFlexMessage(ctx, teacher.LineUserID, altText, contents); err != nil {
		p.log.Errorw("Failed to send welcome teacher notification",
			"teacher_id", payload.RecipientID,
			"error", err)
		return err
	}

	p.log.Infow("Welcome teacher notification sent successfully",
		"teacher_id", payload.RecipientID)
	return nil
}

// processWelcomeAdmin 處理管理員歡迎通知
func (p *AsynqTaskProcessor) processWelcomeAdmin(ctx context.Context, payload *TaskPayload) error {
	// 解析 payload
	var flexPayload map[string]interface{}
	if err := json.Unmarshal([]byte(payload.Payload), &flexPayload); err != nil {
		return fmt.Errorf("failed to parse flex payload: %w", err)
	}

	// 取得管理員的 LINE User ID
	if payload.RecipientType != "ADMIN" {
		return fmt.Errorf("invalid recipient type for welcome admin: %s", payload.RecipientType)
	}

	admin, err := p.adminRepo.GetByIDPtr(ctx, payload.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}

	if admin.LineUserID == "" {
		p.log.Infow("Admin not bound to LINE",
			"admin_id", payload.RecipientID)
		return nil // 不視為錯誤，只是跳過
	}

	// 發送 LINE 訊息
	altText, _ := flexPayload["altText"].(string)
	contents, _ := flexPayload["contents"].(json.RawMessage)

	if err := p.lineBotService.PushFlexMessage(ctx, admin.LineUserID, altText, contents); err != nil {
		p.log.Errorw("Failed to send welcome admin notification",
			"admin_id", payload.RecipientID,
			"error", err)
		return err
	}

	p.log.Infow("Welcome admin notification sent successfully",
		"admin_id", payload.RecipientID)
	return nil
}

// StartWorker 啟動 Asynq Worker
func (s *AsynqNotificationService) StartWorker(ctx context.Context) error {
	mux := asynq.NewServeMux()

	// 註冊任務處理器（使用共享的處理器實例）
	processor := GetAsynqProcessor(s.app)

	mux.HandleFunc(TaskTypeExceptionSubmit, processor.ProcessNotificationTask)
	mux.HandleFunc(TaskTypeExceptionResult, processor.ProcessNotificationTask)
	mux.HandleFunc(TaskTypeWelcomeTeacher, processor.ProcessNotificationTask)
	mux.HandleFunc(TaskTypeWelcomeAdmin, processor.ProcessNotificationTask)

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

// 保留舊的任務處理函數（向後兼容）
func processExceptionSubmit(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Deprecated] Processing exception submit notification: %s", string(t.Payload()))
	return nil
}

func processExceptionResult(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Deprecated] Processing exception result notification: %s", string(t.Payload()))
	return nil
}

func processWelcomeTeacher(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Deprecated] Processing welcome teacher notification: %s", string(t.Payload()))
	return nil
}

func processWelcomeAdmin(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Deprecated] Processing welcome admin notification: %s", string(t.Payload()))
	return nil
}
