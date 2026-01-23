package console

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
)

type ScheduleReminderJob struct {
	app                 *app.App
	notificationService services.NotificationService
	scheduleRuleRepo    *repositories.ScheduleRuleRepository
}

func NewScheduleReminderJob(app *app.App) *ScheduleReminderJob {
	return &ScheduleReminderJob{
		app:                 app,
		notificationService: services.NewNotificationService(app),
		scheduleRuleRepo:    repositories.NewScheduleRuleRepository(app),
	}
}

func (j *ScheduleReminderJob) Name() string {
	return "ScheduleReminderJob"
}

func (j *ScheduleReminderJob) Description() string {
	return "Send schedule reminders to teachers for tomorrow's classes"
}

func (j *ScheduleReminderJob) Repositories() {
	j.notificationService = services.NewNotificationService(j.app)
	j.scheduleRuleRepo = repositories.NewScheduleRuleRepository(j.app)
}

func (j *ScheduleReminderJob) Handle(cronExpr string) error {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)

	weekday := int(tomorrow.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	weekStartDate := tomorrow.Truncate(24 * time.Hour)
	weekEndDate := weekStartDate.Add(24 * time.Hour)

	var rules []models.ScheduleRule
	query := j.app.MySQL.RDB.
		Preload("Exceptions").
		Where("weekday = ? AND effective_range_start_date <= ? AND effective_range_end_date >= ?",
			weekday, weekStartDate, weekEndDate)

	if err := query.Find(&rules).Error; err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.TeacherID == nil {
			continue
		}

		hasCancellation := false
		for _, exc := range rule.Exceptions {
			if exc.Type == "CANCEL" && exc.Status == "APPROVED" {
				if exc.OriginalDate.Format("2006-01-02") == tomorrow.Format("2006-01-02") {
					hasCancellation = true
					break
				}
			}
		}

		if !hasCancellation {
			ctx := context.Background()
			if err := j.notificationService.SendScheduleReminder(ctx, rule.ID, tomorrow); err != nil {
				continue
			}
		}
	}

	return nil
}

type ExceptionReviewJob struct {
	app                   *app.App
	scheduleExceptionRepo *repositories.ScheduleExceptionRepository
}

func NewExceptionReviewJob(app *app.App) *ExceptionReviewJob {
	return &ExceptionReviewJob{
		app:                   app,
		scheduleExceptionRepo: repositories.NewScheduleExceptionRepository(app),
	}
}

func (j *ExceptionReviewJob) Name() string {
	return "ExceptionReviewJob"
}

func (j *ExceptionReviewJob) Description() string {
	return "Review pending exceptions that are older than 24 hours"
}

func (j *ExceptionReviewJob) Repositories() {
	j.scheduleExceptionRepo = repositories.NewScheduleExceptionRepository(j.app)
}

func (j *ExceptionReviewJob) Handle(cronExpr string) error {
	var exceptions []models.ScheduleException
	err := j.app.MySQL.RDB.
		Where("status = ? AND created_at < ?", "PENDING", time.Now().Add(-24*time.Hour)).
		Find(&exceptions).Error
	if err != nil {
		return err
	}

	for range exceptions {
	}

	return nil
}

type CleanupOldNotificationsJob struct {
	app              *app.App
	notificationRepo *repositories.NotificationRepository
}

func NewCleanupOldNotificationsJob(app *app.App) *CleanupOldNotificationsJob {
	return &CleanupOldNotificationsJob{
		app:              app,
		notificationRepo: repositories.NewNotificationRepository(app),
	}
}

func (j *CleanupOldNotificationsJob) Name() string {
	return "CleanupOldNotificationsJob"
}

func (j *CleanupOldNotificationsJob) Description() string {
	return "Delete read notifications older than 3 months"
}

func (j *CleanupOldNotificationsJob) Repositories() {
	j.notificationRepo = repositories.NewNotificationRepository(j.app)
}

func (j *CleanupOldNotificationsJob) Handle(cronExpr string) error {
	cutoffDate := time.Now().AddDate(0, -3, 0)

	err := j.app.MySQL.WDB.
		Where("created_at < ? AND is_read = ?", cutoffDate, true).
		Delete(&models.Notification{}).Error
	if err != nil {
		return err
	}

	return nil
}
