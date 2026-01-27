package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"
	mockRedis "timeLedger/testing/redis"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupNotificationTestApp() (*app.App, *gorm.DB, func()) {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	// AutoMigrate required tables
	if err := mysqlDB.AutoMigrate(
		&models.Teacher{},
		&models.Notification{},
	); err != nil {
		panic(fmt.Sprintf("AutoMigrate error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)

	env := &configs.Env{
		JWTSecret:      "test-jwt-secret-key-for-testing-only",
		AppEnv:         "test",
		AppDebug:       true,
		AppTimezone:    "Asia/Taipei",
	}

	appInstance := &app.App{
		Env:   env,
		Err:   e,
		Tools: nil,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	cleanup := func() {
		mr.Close()
	}

	return appInstance, mysqlDB, cleanup
}

func TestNotificationService_SendTalentInvitationNotification(t *testing.T) {
	t.Run("SendTalentInvitationNotification_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupNotificationTestApp()
		defer cleanup()

		ctx := context.Background()

		// 1. 建立測試資料
		teacher := models.Teacher{
			Name:           "通知測試老師",
			Email:          "notifytest@test.com",
			LineUserID:     "line-user-for-notify-test",
			LineNotifyToken: "", // 測試時不發送 LINE Notify
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 清理函數
		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("user_id = ? AND user_type = ?", teacher.ID, "TEACHER").Delete(&models.Notification{})
		}()

		// 2. 測試發送通知
		svc := services.NewNotificationService(appInstance)

		centerName := "測試中心"
		inviteToken := "TEST-INVITE-123"

		err := svc.SendTalentInvitationNotification(ctx, teacher.ID, centerName, inviteToken)
		if err != nil {
			t.Fatalf("發送人才庫邀請通知失敗: %v", err)
		}

		// 3. 驗證通知記錄已建立
		notificationRepo := repositories.NewNotificationRepository(appInstance)
		notifications, err := notificationRepo.GetNotifications(ctx, teacher.ID, "TEACHER", 10, 0)
		if err != nil {
			t.Fatalf("取得通知列表失敗: %v", err)
		}

		// 找人才庫邀請通知
		foundInvitationNotification := false
		for _, n := range notifications {
			if n.Type == "TALENT_INVITATION" {
				foundInvitationNotification = true

				// 驗證通知內容
				if n.Title == "" {
					t.Error("通知標題不應該為空")
				}

				if n.Message == "" {
					t.Error("通知內容不應該為空")
				}

				// 驗證是否包含邀請連結
				if n.Message == "" || len(n.Message) < 50 {
					t.Error("通知內容應該包含邀請資訊")
				}

				break
			}
		}

		if !foundInvitationNotification {
			t.Error("應該找到人才庫邀請通知")
		}
	})

	t.Run("SendTalentInvitationNotification_WithLINEToken", func(t *testing.T) {
		appInstance, db, cleanup := setupNotificationTestApp()
		defer cleanup()

		ctx := context.Background()

		// 建立有 LINE Notify Token 的老師
		teacher := models.Teacher{
			Name:              "LINE 通知測試老師",
			Email:             "linenotify@test.com",
			LineUserID:        "line-user-with-token",
			LineNotifyToken:   "test-line-notify-token",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("user_id = ? AND user_type = ?", teacher.ID, "TEACHER").Delete(&models.Notification{})
		}()

		svc := services.NewNotificationService(appInstance)

		// 這裡測試時不會真的發送 LINE Notify（因為 token 是假的）
		// 但應該會建立通知記錄
		err := svc.SendTalentInvitationNotification(ctx, teacher.ID, "測試中心", "TEST-TOKEN-456")
		if err != nil {
			t.Fatalf("發送通知失敗: %v", err)
		}

		// 驗證通知已建立
		notificationRepo := repositories.NewNotificationRepository(appInstance)
		notifications, err := notificationRepo.GetNotifications(ctx, teacher.ID, "TEACHER", 10, 0)
		if err != nil {
			t.Fatalf("取得通知失敗: %v", err)
		}

		if len(notifications) == 0 {
			t.Error("應該有通知記錄")
		}

		// 驗證最新的通知是人才庫邀請
		if notifications[0].Type != "TALENT_INVITATION" {
			t.Errorf("通知類型應該是 TALENT_INVITATION，實際是 %s", notifications[0].Type)
		}
	})

	t.Run("SendTalentInvitationNotification_TeacherNotFound", func(t *testing.T) {
		appInstance, _, cleanup := setupNotificationTestApp()
		defer cleanup()

		ctx := context.Background()

		svc := services.NewNotificationService(appInstance)

		// 嘗試發送通知給不存在的老师
		nonExistentTeacherID := uint(999999)
		err := svc.SendTalentInvitationNotification(ctx, nonExistentTeacherID, "測試中心", "TEST-TOKEN")

		if err == nil {
			t.Error("預期錯誤，但沒有收到錯誤")
		}
	})
}

func TestNotificationRepository_CRUD(t *testing.T) {
	t.Run("CreateAndGetNotification", func(t *testing.T) {
		appInstance, db, cleanup := setupNotificationTestApp()
		defer cleanup()

		ctx := context.Background()

		teacher := models.Teacher{
			Name:      "通知倉庫測試老師",
			Email:     "repotest@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("user_id = ?", teacher.ID).Delete(&models.Notification{})
		}()

		// 建立通知
		notification := models.Notification{
			UserID:    teacher.ID,
			UserType:  "TEACHER",
			Title:     "測試通知",
			Message:   "這是測試通知內容",
			Type:      "TEST",
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		notificationRepo := repositories.NewNotificationRepository(appInstance)
		created, err := notificationRepo.Create(ctx, notification)
		if err != nil {
			t.Fatalf("建立通知失敗: %v", err)
		}

		if created.ID == 0 {
			t.Error("通知 ID 不應該為 0")
		}

		// 取得通知列表
		notifications, err := notificationRepo.GetNotifications(ctx, teacher.ID, "TEACHER", 10, 0)
		if err != nil {
			t.Fatalf("取得通知列表失敗: %v", err)
		}

		if len(notifications) < 1 {
			t.Error("應該至少有一筆通知")
		}

		// 標記為已讀
		err = notificationRepo.MarkAsRead(ctx, created.ID)
		if err != nil {
			t.Fatalf("標記已讀失敗: %v", err)
		}

		// 驗證已標記為已讀
		updated, err := notificationRepo.GetByID(ctx, created.ID)
		if err != nil {
			t.Fatalf("取得通知失敗: %v", err)
		}

		if !updated.IsRead {
			t.Error("通知應該已標記為已讀")
		}
	})
}
