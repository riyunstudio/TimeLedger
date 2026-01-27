package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

// NotificationQueueService 通知佇列服務
type NotificationQueueService interface {
	// 佇列管理（Redis）
	PushNotification(ctx context.Context, item *models.NotificationQueue) error
	ProcessQueue(ctx context.Context) error
	GetQueueStats(ctx context.Context) map[string]string

	// 便捷方法 - 發送例外通知給所有管理員
	NotifyExceptionSubmitted(ctx context.Context, exception *models.ScheduleException, teacherName string, centerName string) error
	NotifyExceptionResult(ctx context.Context, exception *models.ScheduleException, teacher *models.Teacher, approved bool, reason string) error

	// 同步發送方法（直接發送，不經佇列）
	NotifyExceptionSubmittedSync(ctx context.Context, exception *models.ScheduleException, teacherName string, centerName string) error
	NotifyExceptionResultSync(ctx context.Context, exception *models.ScheduleException, teacher *models.Teacher, approved bool, reason string) error

	// 便捷方法 - 發送歡迎訊息
	NotifyWelcomeTeacher(ctx context.Context, teacher *models.Teacher, centerName string) error
	NotifyWelcomeAdmin(ctx context.Context, admin *models.AdminUser, centerName string) error
}

// NotificationQueueServiceImpl 通知佇列服務實現
type NotificationQueueServiceImpl struct {
	BaseService
	app             *app.App
	adminRepo       *repositories.AdminUserRepository
	teacherRepo     *repositories.TeacherRepository
	lineBotService  LineBotService
	templateService LineBotTemplateService
	redisQueue      *RedisQueueService
}

// NewNotificationQueueService 建立通知佇列服務
func NewNotificationQueueService(app *app.App) NotificationQueueService {
	return &NotificationQueueServiceImpl{
		app:             app,
		adminRepo:       repositories.NewAdminUserRepository(app),
		teacherRepo:     repositories.NewTeacherRepository(app),
		lineBotService:  NewLineBotService(app),
		templateService: NewLineBotTemplateService(app.Env.FrontendBaseURL),
		redisQueue:      NewRedisQueueService(app),
	}
}

// PushNotification 將通知加入 Redis 佇列
func (s *NotificationQueueServiceImpl) PushNotification(ctx context.Context, item *models.NotificationQueue) error {
	queueItem := &NotificationItem{
		ID:            item.ID,
		Type:          item.Type,
		RecipientID:   item.RecipientID,
		RecipientType: item.RecipientType,
		Payload:       item.Payload,
		CreatedAt:     item.ScheduledAt,
		RetryCount:    0,
	}

	return s.redisQueue.PushNotification(ctx, queueItem)
}

// ProcessQueue 處理佇列（從 Redis 取出並發送）
func (s *NotificationQueueServiceImpl) ProcessQueue(ctx context.Context) error {
	// 先處理延遲重試佇列
	s.redisQueue.ProcessRetryQueue(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			item, err := s.redisQueue.PopNotification(ctx)
			if err != nil {
				fmt.Printf("[ERROR] Failed to pop notification: %v\n", err)
				time.Sleep(1 * time.Second) // 避免 busy loop
				continue
			}
			
			if item == nil {
				// 佇列為空，結束這輪處理
				return nil
			}

			// 處理通知
			if err := s.processRedisNotification(ctx, item); err != nil {
				fmt.Printf("[ERROR] Failed to process notification %d: %v\n", item.ID, err)
				// 加入重試佇列
				s.redisQueue.PushToRetry(ctx, item)
			} else {
				fmt.Printf("[INFO] Notification sent successfully: type=%s, recipient=%d\n", 
					item.Type, item.RecipientID)
			}
		}
	}
}

// processRedisNotification 處理單個 Redis 佇列通知
func (s *NotificationQueueServiceImpl) processRedisNotification(ctx context.Context, item *NotificationItem) error {
	// 解析 payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(item.Payload), &payload); err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	// 取得用戶的 LINE User ID
	var lineUserID string
	
	if item.RecipientType == "ADMIN" {
		admin, err := s.adminRepo.GetByIDPtr(ctx, item.RecipientID)
		if err != nil {
			return fmt.Errorf("failed to get admin: %w", err)
		}
		lineUserID = admin.LineUserID
	} else if item.RecipientType == "TEACHER" {
		teacher, err := s.teacherRepo.GetByID(ctx, item.RecipientID)
		if err != nil {
			return fmt.Errorf("failed to get teacher: %w", err)
		}
		lineUserID = teacher.LineUserID
	} else {
		return fmt.Errorf("unknown recipient type: %s", item.RecipientType)
	}

	if lineUserID == "" {
		return fmt.Errorf("user not bound LINE")
	}

	// 發送 LINE 訊息
	return s.lineBotService.PushMessage(ctx, lineUserID, payload)
}

// GetQueueStats 取得佇列統計
func (s *NotificationQueueServiceImpl) GetQueueStats(ctx context.Context) map[string]string {
	return s.redisQueue.GetStats(ctx)
}

// NotifyExceptionSubmitted 通知管理員有新的例外申請（加入 Redis 佇列）
func (s *NotificationQueueServiceImpl) NotifyExceptionSubmitted(ctx context.Context, exception *models.ScheduleException, teacherName string, centerName string) error {
	// 取得中心的所有管理員
	admins, err := s.adminRepo.GetByCenterID(ctx, exception.CenterID)
	if err != nil {
		return fmt.Errorf("failed to get admins: %w", err)
	}

	// 建立 Flex Message 範本
	flexContent := s.templateService.GetExceptionSubmitTemplate(exception, teacherName, centerName)
	payload, _ := json.Marshal(map[string]interface{}{
		"type":     "flex",
		"altText":  fmt.Sprintf("新的例外申請 - %s 老師", teacherName),
		"contents": flexContent,
	})

	// 為每個已綁定的管理員建立佇列項目
	for _, admin := range admins {
		if !admin.LineNotifyEnabled || admin.LineUserID == "" {
			continue
		}

		queueItem := &models.NotificationQueue{
			Type:          models.NotificationTypeExceptionSubmit,
			RecipientID:   admin.ID,
			RecipientType: "ADMIN",
			Payload:       string(payload),
			Status:        models.NotificationStatusPending,
			ScheduledAt:   time.Now(),
		}

		if err := s.PushNotification(ctx, queueItem); err != nil {
			fmt.Printf("[ERROR] Failed to queue notification for admin %d: %v\n", admin.ID, err)
		}
	}

	return nil
}

// NotifyExceptionSubmittedSync 同步發送例外申請通知給所有管理員（直接發送，不經佇列）
func (s *NotificationQueueServiceImpl) NotifyExceptionSubmittedSync(ctx context.Context, exception *models.ScheduleException, teacherName string, centerName string) error {
	// 取得中心的所有管理員
	admins, err := s.adminRepo.GetByCenterID(ctx, exception.CenterID)
	if err != nil {
		return fmt.Errorf("failed to get admins: %w", err)
	}

	// 建立 Flex Message 範本
	flexContent := s.templateService.GetExceptionSubmitTemplate(exception, teacherName, centerName)
	altText := fmt.Sprintf("新的例外申請 - %s 老師", teacherName)

	// 直接發送給每個已綁定的管理員
	for _, admin := range admins {
		if !admin.LineNotifyEnabled || admin.LineUserID == "" {
			continue
		}

		if err := s.lineBotService.PushFlexMessage(ctx, admin.LineUserID, altText, flexContent); err != nil {
			return fmt.Errorf("failed to send to admin %d: %w", admin.ID, err)
		}
	}

	return nil
}

// NotifyExceptionResult 通知老師例外審核結果（加入 Redis 佇列）
func (s *NotificationQueueServiceImpl) NotifyExceptionResult(ctx context.Context, exception *models.ScheduleException, teacher *models.Teacher, approved bool, reason string) error {
	if teacher.LineUserID == "" {
		return nil
	}

	var flexContent interface{}
	var altText string

	if approved {
		flexContent = s.templateService.GetExceptionApproveTemplate(exception, teacher.Name)
		altText = fmt.Sprintf("✅ 您的例外申請已核准 - %s", exception.GetDate().Format("2006/01/02"))
	} else {
		flexContent = s.templateService.GetExceptionRejectTemplate(exception, teacher.Name, reason)
		altText = fmt.Sprintf("❌ 您的例外申請已拒絕 - %s", exception.GetDate().Format("2006/01/02"))
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"type":     "flex",
		"altText":  altText,
		"contents": flexContent,
	})

	queueItem := &models.NotificationQueue{
		Type:          models.NotificationTypeExceptionResult,
		RecipientID:   teacher.ID,
		RecipientType: "TEACHER",
		Payload:       string(payload),
		Status:        models.NotificationStatusPending,
		ScheduledAt:   time.Now(),
	}

	return s.PushNotification(ctx, queueItem)
}

// NotifyExceptionResultSync 同步發送例外審核結果給老師（直接發送，不經佇列）
func (s *NotificationQueueServiceImpl) NotifyExceptionResultSync(ctx context.Context, exception *models.ScheduleException, teacher *models.Teacher, approved bool, reason string) error {
	if teacher.LineUserID == "" {
		return nil
	}

	var flexContent interface{}
	var altText string

	if approved {
		flexContent = s.templateService.GetExceptionApproveTemplate(exception, teacher.Name)
		altText = fmt.Sprintf("✅ 您的例外申請已核准 - %s", exception.GetDate().Format("2006/01/02"))
	} else {
		flexContent = s.templateService.GetExceptionRejectTemplate(exception, teacher.Name, reason)
		altText = fmt.Sprintf("❌ 您的例外申請已拒絕 - %s", exception.GetDate().Format("2006/01/02"))
	}

	// 直接發送給老師
	return s.lineBotService.PushFlexMessage(ctx, teacher.LineUserID, altText, flexContent)
}

// NotifyWelcomeTeacher 發送老師歡迎訊息（加入 Redis 佇列）
func (s *NotificationQueueServiceImpl) NotifyWelcomeTeacher(ctx context.Context, teacher *models.Teacher, centerName string) error {
	if teacher.LineUserID == "" {
		return nil
	}

	flexContent := s.templateService.GetWelcomeTeacherTemplate(teacher, centerName)
	payload, _ := json.Marshal(map[string]interface{}{
		"type":     "flex",
		"altText":  "👋 歡迎加入 TimeLedger！",
		"contents": flexContent,
	})

	queueItem := &models.NotificationQueue{
		Type:          models.NotificationTypeWelcomeTeacher,
		RecipientID:   teacher.ID,
		RecipientType: "TEACHER",
		Payload:       string(payload),
		Status:        models.NotificationStatusPending,
		ScheduledAt:   time.Now(),
	}

	return s.PushNotification(ctx, queueItem)
}

// NotifyWelcomeAdmin 發送管理員歡迎訊息（加入 Redis 佇列）
func (s *NotificationQueueServiceImpl) NotifyWelcomeAdmin(ctx context.Context, admin *models.AdminUser, centerName string) error {
	if admin.LineUserID == "" {
		return nil
	}

	flexContent := s.templateService.GetWelcomeAdminTemplate(admin, centerName)
	payload, _ := json.Marshal(map[string]interface{}{
		"type":     "flex",
		"altText":  "🎉 歡迎使用 TimeLedger！",
		"contents": flexContent,
	})

	queueItem := &models.NotificationQueue{
		Type:          models.NotificationTypeWelcomeAdmin,
		RecipientID:   admin.ID,
		RecipientType: "ADMIN",
		Payload:       string(payload),
		Status:        models.NotificationStatusPending,
		ScheduledAt:   time.Now(),
	}

	return s.PushNotification(ctx, queueItem)
}

// ProcessQueueHandler 處理佇列的定時任務（可由 cron 或 worker 呼叫）
func (s *NotificationQueueServiceImpl) ProcessQueueHandler() {
	ctx := context.Background()
	if err := s.ProcessQueue(ctx); err != nil {
		fmt.Printf("[ERROR] Failed to process notification queue: %v\n", err)
	}
}
