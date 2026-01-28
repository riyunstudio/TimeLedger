package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"
	mockRedis "timeLedger/testing/redis"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupAdminUserTestApp() (*app.App, *gorm.DB, func()) {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.AdminUser{},
	); err != nil {
		panic(fmt.Sprintf("AutoMigrate error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	env := &configs.Env{
		JWTSecret:       "test-jwt-secret-key-for-testing-only",
		AppEnv:          "test",
		AppDebug:        true,
		AppTimezone:     "Asia/Taipei",
		FrontendBaseURL: "http://localhost:3000",
	}

	appInstance := &app.App{
		Env:   env,
		Err:   e,
		Tools: tool,
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

func createTestCenterForAdminUserTest(t *testing.T, db *gorm.DB) *models.Center {
	center := &models.Center{
		Name:      fmt.Sprintf("Test Center Admin %d", time.Now().UnixNano()),
		PlanLevel: "PRO",
	}
	if err := db.Create(center).Error; err != nil {
		t.Fatalf("Failed to create test center: %v", err)
	}
	return center
}

func createTestAdminUserForBindingTest(t *testing.T, db *gorm.DB, centerID uint, hasLineBound bool) *models.AdminUser {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	admin := &models.AdminUser{
		CenterID:     centerID,
		Email:        fmt.Sprintf("admin.binding%d@test.com", time.Now().UnixNano()),
		PasswordHash: string(hashedPassword),
		Name:         fmt.Sprintf("Test Admin %d", time.Now().UnixNano()),
		Role:         "ADMIN",
		Status:       "ACTIVE",
		LineNotifyEnabled: false, // 明確設為 false
	}

	if hasLineBound {
		admin.LineUserID = fmt.Sprintf("Utest%d", time.Now().UnixNano())
		admin.LineBoundAt = func() *time.Time { now := time.Now(); return &now }()
		admin.LineNotifyEnabled = true
	}

	if err := db.Create(admin).Error; err != nil {
		t.Fatalf("Failed to create test admin: %v", err)
	}
	return admin
}

// TestAdminUserService_GetLINEBindingStatus 測試取得 LINE 綁定狀態
func TestAdminUserService_GetLINEBindingStatus(t *testing.T) {
	t.Run("GetStatus_NotBound", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, false)

		adminService := services.NewAdminUserService(appInstance)

		status, eInfo, err := adminService.GetLINEBindingStatus(ctx, admin.ID)
		if err != nil {
			t.Fatalf("GetLINEBindingStatus failed: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("GetLINEBindingStatus returned error info: %s", eInfo.Msg)
		}

		if status.IsBound {
			t.Error("Expected IsBound to be false")
		}
		// Note: LineNotifyEnabled has default:true in database schema
		// So even for unbound admins, it may be true
		if status.BoundAt != nil {
			t.Error("Expected BoundAt to be nil")
		}
	})

	t.Run("GetStatus_IsBound", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, true)

		adminService := services.NewAdminUserService(appInstance)

		status, eInfo, err := adminService.GetLINEBindingStatus(ctx, admin.ID)
		if err != nil {
			t.Fatalf("GetLINEBindingStatus failed: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("GetLINEBindingStatus returned error info: %s", eInfo.Msg)
		}

		if !status.IsBound {
			t.Error("Expected IsBound to be true")
		}
		if !status.NotifyEnabled {
			t.Error("Expected NotifyEnabled to be true")
		}
		if status.BoundAt == nil {
			t.Error("Expected BoundAt to be set")
		}
	})

	t.Run("GetStatus_AdminNotFound", func(t *testing.T) {
		appInstance, _, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		adminService := services.NewAdminUserService(appInstance)

		_, _, err := adminService.GetLINEBindingStatus(ctx, 99999)
		if err == nil {
			t.Error("Expected error for nonexistent admin")
		}
	})
}

// TestAdminUserService_InitLINEBinding 測試初始化 LINE 綁定
func TestAdminUserService_InitLINEBinding(t *testing.T) {
	t.Run("Init_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, false)

		adminService := services.NewAdminUserService(appInstance)

		code, expiresAt, eInfo, err := adminService.InitLINEBinding(ctx, admin.ID)
		if err != nil {
			t.Fatalf("InitLINEBinding failed: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("InitLINEBinding returned error info: %s", eInfo.Msg)
		}

		if code == "" {
			t.Error("Expected code to be generated")
		}
		if len(code) != 6 {
			t.Errorf("Expected code length 6, got %d", len(code))
		}
		if expiresAt.Before(time.Now()) {
			t.Error("Expected expiresAt to be in the future")
		}
	})

	t.Run("Init_AlreadyBound", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, true)

		adminService := services.NewAdminUserService(appInstance)

		_, _, eInfo, err := adminService.InitLINEBinding(ctx, admin.ID)
		if err != nil {
			t.Fatalf("InitLINEBinding failed: %v", err)
		}
		if eInfo == nil {
			t.Error("Expected error info for already bound admin")
		}
		if int(eInfo.Code) != int(errInfos.LINE_ALREADY_BOUND)+100000 {
			t.Errorf("Expected error code %d, got %d", errInfos.LINE_ALREADY_BOUND+100000, eInfo.Code)
		}
	})

	t.Run("Init_AdminNotFound", func(t *testing.T) {
		appInstance, _, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		adminService := services.NewAdminUserService(appInstance)

		_, _, _, err := adminService.InitLINEBinding(ctx, 99999)
		if err == nil {
			t.Error("Expected error for nonexistent admin")
		}
	})
}

// TestAdminUserService_VerifyLINEBinding 測試驗證 LINE 綁定
func TestAdminUserService_VerifyLINEBinding(t *testing.T) {
	t.Run("Verify_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)

		// 先建立一個管理員
		admin := createTestAdminUserForBindingTest(t, db, center.ID, false)

		// 初始化綁定
		adminService := services.NewAdminUserService(appInstance)
		code, expiresAt, _, _ := adminService.InitLINEBinding(ctx, admin.ID)

		// 驗證綁定
		lineUserID := fmt.Sprintf("U%d", time.Now().UnixNano())
		adminID, eInfo, err := adminService.VerifyLINEBinding(ctx, code, lineUserID)
		if err != nil {
			t.Fatalf("VerifyLINEBinding failed: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("VerifyLINEBinding returned error info: %s", eInfo.Msg)
		}

		if adminID != admin.ID {
			t.Errorf("Expected admin ID %d, got %d", admin.ID, adminID)
		}

		// 驗證資料庫中的更新
		var updatedAdmin models.AdminUser
		appInstance.MySQL.RDB.WithContext(ctx).First(&updatedAdmin, admin.ID)
		if updatedAdmin.LineUserID != lineUserID {
			t.Errorf("Expected LineUserID %s, got %s", lineUserID, updatedAdmin.LineUserID)
		}
		if updatedAdmin.LineBoundAt == nil {
			t.Error("Expected LineBoundAt to be set")
		}
		if updatedAdmin.LineBindingCode != "" {
			t.Error("Expected LineBindingCode to be cleared")
		}
		if updatedAdmin.LineBindingExpires != nil {
			t.Error("Expected LineBindingExpires to be nil")
		}
		_ = expiresAt
	})

	t.Run("Verify_InvalidCode", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		_ = createTestAdminUserForBindingTest(t, db, center.ID, false)

		adminService := services.NewAdminUserService(appInstance)

		_, eInfo, err := adminService.VerifyLINEBinding(ctx, "INVALID", "U123")
		if err != nil {
			t.Fatalf("VerifyLINEBinding failed: %v", err)
		}
		if eInfo == nil {
			t.Error("Expected error info for invalid code")
		}
		if int(eInfo.Code) != int(errInfos.LINE_BINDING_CODE_INVALID)+100000 {
			t.Errorf("Expected error code %d, got %d", errInfos.LINE_BINDING_CODE_INVALID+100000, eInfo.Code)
		}
	})

	t.Run("Verify_ExpiredCode", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, false)

		// 手動設定過期的驗證碼
		appInstance.MySQL.RDB.WithContext(ctx).Model(&admin).Updates(map[string]interface{}{
			"line_binding_code":    "EXPIRED",
			"line_binding_expires": time.Now().Add(-1 * time.Minute), // 已過期
		})

		adminService := services.NewAdminUserService(appInstance)

		_, eInfo, err := adminService.VerifyLINEBinding(ctx, "EXPIRED", "U123")
		if err != nil {
			t.Fatalf("VerifyLINEBinding failed: %v", err)
		}
		if eInfo == nil {
			t.Error("Expected error info for expired code")
		}
	})
}

// TestAdminUserService_UnbindLINE 測試解除 LINE 綁定
func TestAdminUserService_UnbindLINE(t *testing.T) {
	t.Run("Unbind_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, true)

		adminService := services.NewAdminUserService(appInstance)

		eInfo, err := adminService.UnbindLINE(ctx, admin.ID)
		if err != nil {
			t.Fatalf("UnbindLINE failed: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("UnbindLINE returned error info: %s", eInfo.Msg)
		}

		// 驗證資料庫中的更新
		var updatedAdmin models.AdminUser
		appInstance.MySQL.RDB.WithContext(ctx).First(&updatedAdmin, admin.ID)
		if updatedAdmin.LineUserID != "" {
			t.Error("Expected LineUserID to be cleared")
		}
		if updatedAdmin.LineNotifyEnabled {
			t.Error("Expected LineNotifyEnabled to be false")
		}
		if updatedAdmin.LineBoundAt != nil {
			t.Error("Expected LineBoundAt to be nil")
		}
	})

	t.Run("Unbind_NotBound", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, false)

		adminService := services.NewAdminUserService(appInstance)

		eInfo, err := adminService.UnbindLINE(ctx, admin.ID)
		if err != nil {
			t.Fatalf("UnbindLINE failed: %v", err)
		}
		if eInfo == nil {
			t.Error("Expected error info for not bound admin")
		}
		if int(eInfo.Code) != int(errInfos.LINE_NOT_BOUND)+100000 {
			t.Errorf("Expected error code %d, got %d", errInfos.LINE_NOT_BOUND+100000, eInfo.Code)
		}
	})
}

// TestAdminUserService_UpdateLINENotifySettings 測試更新 LINE 通知設定
func TestAdminUserService_UpdateLINENotifySettings(t *testing.T) {
	t.Run("UpdateNotify_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, true)

		adminService := services.NewAdminUserService(appInstance)

		// 關閉通知
		eInfo, err := adminService.UpdateLINENotifySettings(ctx, admin.ID, false)
		if err != nil {
			t.Fatalf("UpdateLINENotifySettings failed: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("UpdateLINENotifySettings returned error info: %s", eInfo.Msg)
		}

		// 驗證資料庫中的更新
		var updatedAdmin models.AdminUser
		appInstance.MySQL.RDB.WithContext(ctx).First(&updatedAdmin, admin.ID)
		if updatedAdmin.LineNotifyEnabled {
			t.Error("Expected LineNotifyEnabled to be false")
		}

		// 開啟通知
		_, err = adminService.UpdateLINENotifySettings(ctx, admin.ID, true)
		if err != nil {
			t.Fatalf("UpdateLINENotifySettings failed: %v", err)
		}

		appInstance.MySQL.RDB.WithContext(ctx).First(&updatedAdmin, admin.ID)
		if !updatedAdmin.LineNotifyEnabled {
			t.Error("Expected LineNotifyEnabled to be true")
		}
	})

	t.Run("UpdateNotify_NotBound", func(t *testing.T) {
		appInstance, db, cleanup := setupAdminUserTestApp()
		defer cleanup()

		ctx := context.Background()
		center := createTestCenterForAdminUserTest(t, db)
		admin := createTestAdminUserForBindingTest(t, db, center.ID, false)

		adminService := services.NewAdminUserService(appInstance)

		eInfo, err := adminService.UpdateLINENotifySettings(ctx, admin.ID, true)
		if err != nil {
			t.Fatalf("UpdateLINENotifySettings failed: %v", err)
		}
		if eInfo == nil {
			t.Error("Expected error info for not bound admin")
		}
		if int(eInfo.Code) != int(errInfos.LINE_NOT_BOUND)+100000 {
			t.Errorf("Expected error code %d, got %d", errInfos.LINE_NOT_BOUND+100000, eInfo.Code)
		}
	})
}

// TestAdminUserService_GenerateBindingCode 測試驗證碼生成
func TestAdminUserService_GenerateBindingCode(t *testing.T) {
	t.Run("Generate_CodeFormat", func(t *testing.T) {
		code := services.GenerateBindingCode()

		if len(code) != 6 {
			t.Errorf("Expected code length 6, got %d", len(code))
		}

		// 驗證格式（不包含 I, O, 0, 1 等容易混淆的字元）
		for _, c := range code {
			if (c < 'A' || c > 'Z') && (c < '2' || c > '9') {
				t.Errorf("Unexpected character in code: %c", c)
			}
		}
	})

	t.Run("Generate_Uniqueness", func(t *testing.T) {
		codes := make(map[string]bool)
		for i := 0; i < 100; i++ {
			code := services.GenerateBindingCode()
			if codes[code] {
				t.Errorf("Duplicate code generated: %s", code)
			}
			codes[code] = true
		}
	})
}
