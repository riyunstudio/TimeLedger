package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"
	mockRedis "timeLedger/testing/redis"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupCenterInvitationTestApp() (*app.App, *gorm.DB, func()) {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	// AutoMigrate required tables
	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.Teacher{},
		&models.CenterInvitation{},
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

func TestCenterInvitationRepository_CRUD(t *testing.T) {
	t.Run("CreateInvitation", func(t *testing.T) {
		appInstance, db, cleanup := setupCenterInvitationTestApp()
		defer cleanup()

		ctx := context.Background()

		// 建立測試資料
		center := models.Center{
			Name:      "Repo 測試中心",
			PlanLevel: "STARTER",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&center).Error; err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher := models.Teacher{
			Name:      "Repo 測試老師",
			Email:     "repoteacher@test.com",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 清理
		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("teacher_id = ?", teacher.ID).Delete(&models.CenterInvitation{})
		}()

		// 建立邀請
		invitation := models.CenterInvitation{
			CenterID:    center.ID,
			TeacherID:   teacher.ID,
			InvitedBy:   1,
			Token:       "TEST-TOKEN-001",
			Status:      models.InvitationStatusPending,
			InviteType:  models.InvitationTypeTalentPool,
			Message:     "歡迎加入！",
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:   time.Now(),
		}

		repo := repositories.NewCenterInvitationRepository(appInstance)
		created, err := repo.Create(ctx, invitation)
		if err != nil {
			t.Fatalf("建立邀請失敗: %v", err)
		}

		if created.ID == 0 {
			t.Error("邀請 ID 不應該為 0")
		}

		if created.Token != "TEST-TOKEN-001" {
			t.Errorf("預期 Token 為 TEST-TOKEN-001，實際為 %s", created.Token)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		appInstance, db, cleanup := setupCenterInvitationTestApp()
		defer cleanup()

		ctx := context.Background()

		center := models.Center{
			Name:      "GetByID 測試中心",
			PlanLevel: "STARTER",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&center).Error; err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher := models.Teacher{
			Name:      "GetByID 測試老師",
			Email:     "getbyid@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		invitation := models.CenterInvitation{
			CenterID:    center.ID,
			TeacherID:   teacher.ID,
			InvitedBy:   1,
			Token:       "GETBYID-TOKEN",
			Status:      models.InvitationStatusPending,
			InviteType:  models.InvitationTypeTalentPool,
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:   time.Now(),
		}
		if err := db.WithContext(ctx).Create(&invitation).Error; err != nil {
			t.Fatalf("建立邀請失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("teacher_id = ?", teacher.ID).Delete(&models.CenterInvitation{})
		}()

		repo := repositories.NewCenterInvitationRepository(appInstance)

		// 測試 GetByID
		fetched, err := repo.GetByID(ctx, invitation.ID)
		if err != nil {
			t.Fatalf("取得邀請失敗: %v", err)
		}

		if fetched.ID != invitation.ID {
			t.Errorf("預期 ID 為 %d，實際為 %d", invitation.ID, fetched.ID)
		}

		if fetched.Token != "GETBYID-TOKEN" {
			t.Errorf("預期 Token 為 GETBYID-TOKEN，實際為 %s", fetched.Token)
		}
	})

	t.Run("GetByTeacherAndCenter", func(t *testing.T) {
		appInstance, db, cleanup := setupCenterInvitationTestApp()
		defer cleanup()

		ctx := context.Background()

		center := models.Center{
			Name:      "查詢測試中心",
			PlanLevel: "STARTER",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&center).Error; err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher := models.Teacher{
			Name:      "查詢測試老師",
			Email:     "query@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立多筆邀請
		for i := 0; i < 3; i++ {
			invitation := models.CenterInvitation{
				CenterID:    center.ID,
				TeacherID:   teacher.ID,
				InvitedBy:   1,
				Token:       fmt.Sprintf("QUERY-TOKEN-%d", i),
				Status:      models.InvitationStatusPending,
				InviteType:  models.InvitationTypeTalentPool,
				ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
				CreatedAt:   time.Now(),
			}
			if err := db.WithContext(ctx).Create(&invitation).Error; err != nil {
				t.Fatalf("建立邀請 %d 失敗: %v", i, err)
			}
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("teacher_id = ?", teacher.ID).Delete(&models.CenterInvitation{})
		}()

		repo := repositories.NewCenterInvitationRepository(appInstance)

		invitations, err := repo.GetByTeacherAndCenter(ctx, teacher.ID, center.ID)
		if err != nil {
			t.Fatalf("取得邀請列表失敗: %v", err)
		}

		if len(invitations) != 3 {
			t.Errorf("預期找到 3 筆邀請，實際找到 %d 筆", len(invitations))
		}
	})

	t.Run("HasPendingInvitation", func(t *testing.T) {
		appInstance, db, cleanup := setupCenterInvitationTestApp()
		defer cleanup()

		ctx := context.Background()

		center := models.Center{
			Name:      "待處理測試中心",
			PlanLevel: "STARTER",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&center).Error; err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher := models.Teacher{
			Name:      "待處理測試老師",
			Email:     "pending@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立待處理邀請
		pendingInvitation := models.CenterInvitation{
			CenterID:    center.ID,
			TeacherID:   teacher.ID,
			InvitedBy:   1,
			Token:       "PENDING-TOKEN",
			Status:      models.InvitationStatusPending,
			InviteType:  models.InvitationTypeTalentPool,
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:   time.Now(),
		}
		if err := db.WithContext(ctx).Create(&pendingInvitation).Error; err != nil {
			t.Fatalf("建立邀請失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("teacher_id = ?", teacher.ID).Delete(&models.CenterInvitation{})
		}()

		repo := repositories.NewCenterInvitationRepository(appInstance)

		// 測試有待處理邀請
		hasPending, err := repo.HasPendingInvitation(ctx, teacher.ID, center.ID)
		if err != nil {
			t.Fatalf("檢查待處理邀請失敗: %v", err)
		}
		if !hasPending {
			t.Error("預期有待處理邀請")
		}

		// 測試沒有待處理邀請的組合
		hasPending2, err := repo.HasPendingInvitation(ctx, teacher.ID, 999)
		if err != nil {
			t.Fatalf("檢查待處理邀請失敗: %v", err)
		}
		if hasPending2 {
			t.Error("不應該有待處理邀請")
		}
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		appInstance, db, cleanup := setupCenterInvitationTestApp()
		defer cleanup()

		ctx := context.Background()

		center := models.Center{
			Name:      "狀態更新測試中心",
			PlanLevel: "STARTER",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&center).Error; err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher := models.Teacher{
			Name:      "狀態更新測試老師",
			Email:     "updatestatus@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		invitation := models.CenterInvitation{
			CenterID:    center.ID,
			TeacherID:   teacher.ID,
			InvitedBy:   1,
			Token:       "STATUS-TOKEN",
			Status:      models.InvitationStatusPending,
			InviteType:  models.InvitationTypeTalentPool,
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:   time.Now(),
		}
		if err := db.WithContext(ctx).Create(&invitation).Error; err != nil {
			t.Fatalf("建立邀請失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("teacher_id = ?", teacher.ID).Delete(&models.CenterInvitation{})
		}()

		repo := repositories.NewCenterInvitationRepository(appInstance)

		// 更新為已接受
		err := repo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusAccepted)
		if err != nil {
			t.Fatalf("更新狀態失敗: %v", err)
		}

		// 驗證更新結果
		updated, err := repo.GetByID(ctx, invitation.ID)
		if err != nil {
			t.Fatalf("取得邀請失敗: %v", err)
		}

		if updated.Status != models.InvitationStatusAccepted {
			t.Errorf("預期狀態為 ACCEPTED，實際為 %s", updated.Status)
		}

		if updated.RespondedAt == nil {
			t.Error("回覆時間不應該為空")
		}
	})

	t.Run("CountByCenter", func(t *testing.T) {
		appInstance, db, cleanup := setupCenterInvitationTestApp()
		defer cleanup()

		ctx := context.Background()

		center := models.Center{
			Name:      "統計測試中心",
			PlanLevel: "STARTER",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&center).Error; err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher := models.Teacher{
			Name:      "統計測試老師",
			Email:     "count@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立不同狀態的邀請
		invitations := []models.CenterInvitation{
			{Status: models.InvitationStatusPending},
			{Status: models.InvitationStatusPending},
			{Status: models.InvitationStatusAccepted},
			{Status: models.InvitationStatusDeclined},
		}

		for i := range invitations {
			invitations[i].CenterID = center.ID
			invitations[i].TeacherID = teacher.ID
			invitations[i].InvitedBy = 1
			invitations[i].Token = fmt.Sprintf("COUNT-TOKEN-%d", i)
			invitations[i].InviteType = models.InvitationTypeTalentPool
			invitations[i].ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
			invitations[i].CreatedAt = time.Now()
			if err := db.WithContext(ctx).Create(&invitations[i]).Error; err != nil {
				t.Fatalf("建立邀請 %d 失敗: %v", i, err)
			}
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.CenterInvitation{})
		}()

		repo := repositories.NewCenterInvitationRepository(appInstance)

		pending, accepted, declined, err := repo.CountByCenter(ctx, center.ID)
		if err != nil {
			t.Fatalf("統計邀請失敗: %v", err)
		}

		if pending != 2 {
			t.Errorf("預期待處理數量為 2，實際為 %d", pending)
		}

		if accepted != 1 {
			t.Errorf("預期已接受數量為 1，實際為 %d", accepted)
		}

		if declined != 1 {
			t.Errorf("預期已拒絕數量為 1，實際為 %d", declined)
		}
	})
}
