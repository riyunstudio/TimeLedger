package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

type NotificationServiceImpl struct {
	BaseService
	app                   *app.App
	notificationRepo      *repositories.NotificationRepository
	teacherRepo           *repositories.TeacherRepository
	scheduleRuleRepo      *repositories.ScheduleRuleRepository
	scheduleExceptionRepo *repositories.ScheduleExceptionRepository
	LINENotifyService     LINENotifyService
}

func NewNotificationService(app *app.App) NotificationService {
	return &NotificationServiceImpl{
		app:                   app,
		notificationRepo:      repositories.NewNotificationRepository(app),
		teacherRepo:           repositories.NewTeacherRepository(app),
		scheduleRuleRepo:      repositories.NewScheduleRuleRepository(app),
		scheduleExceptionRepo: repositories.NewScheduleExceptionRepository(app),
		LINENotifyService:     NewLINENotifyService(app),
	}
}

func (s *NotificationServiceImpl) SendTeacherNotification(ctx context.Context, teacherID uint, title, message string) error {
	return s.SendTeacherNotificationWithType(ctx, teacherID, title, message, "SYSTEM")
}

func (s *NotificationServiceImpl) SendTeacherNotificationWithType(ctx context.Context, teacherID uint, title, message string, notificationType string) error {
	notification := &models.Notification{
		UserID:    teacherID,
		UserType:  "TEACHER",
		Title:     title,
		Message:   message,
		Type:      notificationType,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	teacher, err := s.teacherRepo.GetByID(ctx, teacherID)
	if err != nil {
		return err
	}

	if teacher.LineNotifyToken != "" {
		go s.LINENotifyService.SendMessage(ctx, teacher.LineNotifyToken, title+"\n"+message)
	}

	return nil
}

func (s *NotificationServiceImpl) SendAdminNotification(ctx context.Context, centerID uint, title, message string, notificationType string) error {
	notification := &models.Notification{
		UserID:    0,
		UserType:  "ADMIN",
		CenterID:  centerID,
		Title:     title,
		Message:   message,
		Type:      notificationType,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *NotificationServiceImpl) SendScheduleReminder(ctx context.Context, ruleID uint, date time.Time) error {
	rule, err := s.scheduleRuleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return err
	}

	if rule.TeacherID == nil {
		return nil
	}

	title := "課程提醒"
	message := "您有課程即將開始\n\n時間: " + date.Format("2006-01-02 15:04")

	return s.SendTeacherNotification(ctx, *rule.TeacherID, title, message)
}

func (s *NotificationServiceImpl) SendExceptionNotification(ctx context.Context, exceptionID uint) error {
	exception, err := s.scheduleExceptionRepo.GetByID(ctx, exceptionID)
	if err != nil {
		return err
	}

	rule, err := s.scheduleRuleRepo.GetByID(ctx, exception.RuleID)
	if err != nil {
		return err
	}

	if rule.TeacherID == nil {
		return nil
	}

	title := "例外單通知"
	message := "有新的排課例外單\n類型: " + exception.ExceptionType + "\n日期: " + exception.OriginalDate.Format("2006-01-02")

	return s.SendTeacherNotification(ctx, *rule.TeacherID, title, message)
}

func (s *NotificationServiceImpl) SendReviewNotification(ctx context.Context, exceptionID uint, approved bool) error {
	exception, err := s.scheduleExceptionRepo.GetByID(ctx, exceptionID)
	if err != nil {
		return err
	}

	rule, err := s.scheduleRuleRepo.GetByID(ctx, exception.RuleID)
	if err != nil {
		return err
	}

	if rule.TeacherID == nil {
		return nil
	}

	title := "例外單審核結果"
	status := "已通過"
	if !approved {
		status = "已拒絕"
	}
	message := "您的例外單" + status + "\n日期: " + exception.OriginalDate.Format("2006-01-02")

	return s.SendTeacherNotificationWithType(ctx, *rule.TeacherID, title, message, "REVIEW_RESULT")
}

func (s *NotificationServiceImpl) CreateNotificationRecord(ctx context.Context, notification *models.Notification) error {
	return s.notificationRepo.Create(ctx, notification)
}

func (s *NotificationServiceImpl) GetNotifications(ctx context.Context, userID uint, userType string, limit int, offset int) ([]models.Notification, error) {
	return s.notificationRepo.List(ctx, userID, userType, limit, offset)
}

func (s *NotificationServiceImpl) MarkAsRead(ctx context.Context, notificationID uint) error {
	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

func (s *NotificationServiceImpl) MarkAllAsRead(ctx context.Context, userID uint, userType string) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID, userType)
}

// SendTalentInvitationNotification 發送人才庫邀請通知
func (s *NotificationServiceImpl) SendTalentInvitationNotification(ctx context.Context, teacherID uint, centerName string, inviteToken string) error {
	title := "🎉 人才庫邀請通知"
	message := fmt.Sprintf(`%s 邀請您加入人才庫！

點擊以下連結接受邀請：
%s

邀請碼：%s
（如非本人，請忽略此訊息）`, centerName, s.buildInvitationLink(inviteToken), inviteToken)

	// 建立通知記錄
	notification := &models.Notification{
		UserID:    teacherID,
		UserType:  "TEACHER",
		Title:     title,
		Message:   message,
		Type:      "TALENT_INVITATION",
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// 取得老師資料
	teacher, err := s.teacherRepo.GetByID(ctx, teacherID)
	if err != nil {
		return err
	}

	// 發送 LINE Notify（如果有的話）
	if teacher.LineNotifyToken != "" {
		go s.LINENotifyService.SendMessage(ctx, teacher.LineNotifyToken, title+"\n\n"+message)
	}

	return nil
}

// buildInvitationLink 建立邀請連結
func (s *NotificationServiceImpl) buildInvitationLink(token string) string {
	// 這裡應該從環境變數取得正確的前端 URL
	baseURL := "https://timeledger.app"
	return fmt.Sprintf("%s/teacher/invitation/accept?token=%s", baseURL, token)
}
