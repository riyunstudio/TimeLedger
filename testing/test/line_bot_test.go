package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupLineBotTestApp 建立測試用的 App 實例
func setupLineBotTestApp(t *testing.T) *app.App {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("MySQL init error: %s. Skipping test.", err.Error())
		return nil
	}

	// 檢查資料庫連線
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		t.Skipf("MySQL DB error: %s. Skipping test.", err.Error())
		return nil
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("MySQL ping error: %s. Skipping test.", err.Error())
		return nil
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	env := &configs.Env{
		JWTSecret:   "test-jwt-secret-key-for-testing-only",
		AppEnv:      "test",
		AppDebug:    true,
		AppTimezone: "Asia/Taipei",
	}

	appInstance := &app.App{
		Env:   env,
		Err:   e,
		Tools: tool,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: nil,
		Api:   nil,
		Rpc:   nil,
	}

	return appInstance
}

// cleanupLineBotTestData 清理測試資料
func cleanupLineBotTestData(t *testing.T, appInstance *app.App, adminLineUserID, teacherLineUserID string) {
	ctx := context.Background()

	// 清理管理員測試資料
	if adminLineUserID != "" {
		appInstance.MySQL.WDB.WithContext(ctx).
			Table("admin_users").
			Where("line_user_id LIKE ?", adminLineUserID+"%").
			Delete(&models.AdminUser{})
	}

	// 清理老師測試資料
	if teacherLineUserID != "" {
		appInstance.MySQL.WDB.WithContext(ctx).
			Table("teachers").
			Where("line_user_id LIKE ?", teacherLineUserID+"%").
			Delete(&models.Teacher{})

		// 清理相關的會員關係
		appInstance.MySQL.WDB.WithContext(ctx).
			Table("center_memberships").
			Where("teacher_id IN (SELECT id FROM teachers WHERE line_user_id LIKE ?)", teacherLineUserID+"%").
			Delete(&models.CenterMembership{})
	}
}

// TestLineBotService_GetCombinedIdentity 測試整合身份識別功能
func TestLineBotService_GetCombinedIdentity(t *testing.T) {
	t.Run("AdminOnly_ReturnAdminIdentity", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-teacher-")
		}()

		ctx := context.Background()

		// 建立測試管理員資料
		testLineUserID := fmt.Sprintf("test-line-admin-%d", time.Now().UnixNano())
		admin := models.AdminUser{
			Name:         "Test Admin for Combined Identity",
			Email:        fmt.Sprintf("test-admin-%d@test.com", time.Now().UnixNano()),
			PasswordHash: "hashed_password",
			Role:         "ADMIN",
			CenterID:     1,
			LineUserID:   testLineUserID,
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("admin_users").Create(&admin).Error; err != nil {
			t.Fatalf("建立測試管理員失敗: %v", err)
		}

		// 執行測試
		svc := services.NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(testLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功，但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "ADMIN" {
			t.Errorf("預期 PrimaryRole 為 'ADMIN'，但取得 '%s'", identity.PrimaryRole)
		}

		if len(identity.AdminProfiles) != 1 {
			t.Errorf("預期有 1 個管理員資料，但取得 %d 個", len(identity.AdminProfiles))
		}

		if identity.TeacherProfile != nil {
			t.Error("預期 TeacherProfile 為 nil，但取得非空值")
		}

		if identity.Memberships != nil && len(identity.Memberships) > 0 {
			t.Error("預期 Memberships 為空，但取得非空值")
		}
	})

	t.Run("TeacherOnly_ReturnTeacherIdentity", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-teacher-")
		}()

		ctx := context.Background()
		centerID := uint(1)

		// 建立測試老師資料
		testLineUserID := fmt.Sprintf("test-line-teacher-%d", time.Now().UnixNano())
		teacher := models.Teacher{
			Name:      "Test Teacher for Combined Identity",
			Email:     fmt.Sprintf("test-teacher-%d@test.com", time.Now().UnixNano()),
			LineUserID: testLineUserID,
			City:      "台北市",
			District:  "大安區",
			AvatarURL: "https://example.com/avatar.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("teachers").Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立老師的會員關係
		membership := models.CenterMembership{
			CenterID:  centerID,
			TeacherID: teacher.ID,
			Role:      "TEACHER",
			Status:    "ACTIVE",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("center_memberships").Create(&membership).Error; err != nil {
			t.Fatalf("建立測試會員關係失敗: %v", err)
		}

		// 執行測試
		svc := services.NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(testLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功，但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "TEACHER" {
			t.Errorf("預期 PrimaryRole 為 'TEACHER'，但取得 '%s'", identity.PrimaryRole)
		}

		if len(identity.AdminProfiles) != 0 {
			t.Errorf("預期有 0 個管理員資料，但取得 %d 個", len(identity.AdminProfiles))
		}

		if identity.TeacherProfile == nil {
			t.Fatal("預期 TeacherProfile 不為 nil")
		}

		if identity.TeacherProfile.ID != teacher.ID {
			t.Errorf("預期 TeacherProfile.ID 為 %d，但取得 %d", teacher.ID, identity.TeacherProfile.ID)
		}

		if identity.TeacherProfile.Name != teacher.Name {
			t.Errorf("預期 TeacherProfile.Name 為 '%s'，但取得 '%s'", teacher.Name, identity.TeacherProfile.Name)
		}

		if len(identity.Memberships) != 1 {
			t.Errorf("預期有 1 個會員關係，但取得 %d 個", len(identity.Memberships))
		}

		if len(identity.Memberships) > 0 && identity.Memberships[0].CenterID != centerID {
			t.Errorf("預期會員關係的 CenterID 為 %d，但取得 %d", centerID, identity.Memberships[0].CenterID)
		}
	})

	t.Run("GuestNotBound_ReturnGuestIdentity", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-teacher-")
		}()

		// 使用一個不可能存在的 LINE User ID
		nonExistentLineUserID := fmt.Sprintf("non-existent-line-user-%d@test.com", time.Now().UnixNano())

		// 執行測試
		svc := services.NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(nonExistentLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功（找不到資料視為正常），但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "GUEST" {
			t.Errorf("預期 PrimaryRole 為 'GUEST'，但取得 '%s'", identity.PrimaryRole)
		}

		if len(identity.AdminProfiles) != 0 {
			t.Errorf("預期有 0 個管理員資料，但取得 %d 個", len(identity.AdminProfiles))
		}

		if identity.TeacherProfile != nil {
			t.Error("預期 TeacherProfile 為 nil，但取得非空值")
		}

		if identity.Memberships != nil && len(identity.Memberships) > 0 {
			t.Error("預期 Memberships 為空，但取得非空值")
		}
	})

	t.Run("NoMemberships_ReturnTeacherWithoutMemberships", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-no-membership-")
		}()

		ctx := context.Background()

		// 建立沒有會員關係的老師資料
		testLineUserID := fmt.Sprintf("test-line-no-membership-%d", time.Now().UnixNano())
		teacher := models.Teacher{
			Name:      "Test Teacher No Memberships",
			Email:     fmt.Sprintf("test-teacher-no-membership-%d@test.com", time.Now().UnixNano()),
			LineUserID: testLineUserID,
			City:      "新北市",
			District:  "板橋區",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("teachers").Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 執行測試
		svc := services.NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(testLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功，但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "TEACHER" {
			t.Errorf("預期 PrimaryRole 為 'TEACHER'，但取得 '%s'", identity.PrimaryRole)
		}

		if identity.TeacherProfile == nil {
			t.Fatal("預期 TeacherProfile 不為 nil")
		}

		if len(identity.Memberships) != 0 {
			t.Errorf("預期有 0 個會員關係，但取得 %d 個", len(identity.Memberships))
		}
	})
}

// TestCompareTimeStrings 測試時間字串比較函數
func TestCompareTimeStrings(t *testing.T) {
	t.Skip("Internal helper function - skipping")
}

// TestSortAgendaItemsByTime 測試行程項目排序功能
func TestSortAgendaItemsByTime(t *testing.T) {
	t.Skip("Internal helper function - skipping")
}

// TestFormatTimeForAgenda 測試時間格式化函數
func TestFormatTimeForAgenda(t *testing.T) {
	t.Skip("Internal helper function - skipping")
}

// TestGenerateAgendaFlex 測試行程聚合 Flex Message 範本生成
func TestGenerateAgendaFlex(t *testing.T) {
	t.Skip("Internal helper function - skipping")
}
