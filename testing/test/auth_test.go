package test

import (
	"context"
	"fmt"
	"math/rand"
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
	"timeLedger/global/logger"
	mockRedis "timeLedger/testing/redis"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
)

// Generate unique identifier for test data
func generateUniqueID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(10000))
}

func setupAuthTestApp() (*app.App, *gorm.DB, func()) {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.AdminUser{},
		&models.Teacher{},
		&models.AdminLoginHistory{},
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

	// 初始化測試用的 logger（使用空配置，避免寫入檔案）
	logConfig := &logger.Config{
		Level:      "debug",
		Format:     "console",
		OutputPath: "",
	}
	_, _ = logger.Initialize(logConfig)

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

func createTestCenter(t *testing.T, db *gorm.DB) *models.Center {
	uniqueID := generateUniqueID()
	center := &models.Center{
		Name:      fmt.Sprintf("Test Center %s", uniqueID),
		PlanLevel: "PRO",
	}
	if err := db.Create(center).Error; err != nil {
		t.Fatalf("Failed to create test center: %v", err)
	}
	return center
}

func createTestAdminUser(t *testing.T, db *gorm.DB, centerID uint, password string) *models.AdminUser {
	uniqueID := generateUniqueID()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	admin := &models.AdminUser{
		CenterID:     centerID,
		Email:        fmt.Sprintf("admin_%s@test.com", uniqueID),
		PasswordHash: string(hashedPassword),
		Name:         fmt.Sprintf("Test Admin %s", uniqueID),
		Role:         "OWNER",
		Status:       "ACTIVE",
	}
	if err := db.Create(admin).Error; err != nil {
		t.Fatalf("Failed to create test admin: %v", err)
	}
	return admin
}

func createTestTeacher(t *testing.T, db *gorm.DB, lineUserID string) *models.Teacher {
	uniqueID := generateUniqueID()
	teacher := &models.Teacher{
		LineUserID: lineUserID,
		Name:       fmt.Sprintf("Test Teacher %s", uniqueID),
		Email:      fmt.Sprintf("teacher_%s@test.com", uniqueID),
		Bio:        "Test bio",
		City:       "台北市",
		District:   "大安區",
	}
	if err := db.Create(teacher).Error; err != nil {
		t.Fatalf("Failed to create test teacher: %v", err)
	}
	return teacher
}

func TestAuthController_AdminLogin_Success(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	authService := services.NewAuthService(appInstance)
	adminRepo := repositories.NewAdminUserRepository(appInstance)

	response, err := authService.AdminLogin(context.Background(), admin.Email, "password123", "127.0.0.1", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("AdminLogin failed: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}

	if response.User.Email != admin.Email {
		t.Errorf("Expected email '%s', got '%s'", admin.Email, response.User.Email)
	}

	if response.User.UserType != "OWNER" {
		t.Errorf("Expected UserType 'OWNER', got '%s'", response.User.UserType)
	}

	_ = adminRepo
}

func TestAuthController_AdminLogin_InvalidPassword(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	authService := services.NewAuthService(appInstance)

	_, err := authService.AdminLogin(context.Background(), admin.Email, "wrongpassword", "127.0.0.1", "Mozilla/5.0")
	if err == nil {
		t.Error("Expected error for invalid password")
	}
}

func TestAuthController_AdminLogin_AdminNotFound(t *testing.T) {
	appInstance, _, cleanup := setupAuthTestApp()
	defer cleanup()

	_ = createTestCenter(t, appInstance.MySQL.RDB)

	authService := services.NewAuthService(appInstance)

	_, err := authService.AdminLogin(context.Background(), "nonexistent@test.com", "password123", "127.0.0.1", "Mozilla/5.0")
	if err == nil {
		t.Error("Expected error for nonexistent admin")
	}
}

func TestAuthController_TeacherLineLogin_Success(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	uniqueID := generateUniqueID()
	lineUserID := fmt.Sprintf("U%s", uniqueID)
	teacher := createTestTeacher(t, db, lineUserID)
	_ = teacher.ID

	authService := services.NewAuthService(appInstance)

	response, err, _ := authService.TeacherLineLogin(context.Background(), lineUserID, "mock-token")
	if err != nil {
		t.Fatalf("TeacherLineLogin failed: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}

	if response.User.UserType != "TEACHER" {
		t.Errorf("Expected UserType 'TEACHER', got '%s'", response.User.UserType)
	}
}

func TestAuthController_TeacherLineLogin_TeacherNotFound(t *testing.T) {
	appInstance, _, cleanup := setupAuthTestApp()
	defer cleanup()

	authService := services.NewAuthService(appInstance)

	_, err, _ := authService.TeacherLineLogin(context.Background(), "UNKNOWN_USER", "mock-token")
	if err == nil {
		t.Error("Expected error for nonexistent teacher")
	}
}

func TestAuthController_RefreshToken_Success(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	uniqueID := generateUniqueID()
	lineUserID := fmt.Sprintf("U%s", uniqueID)
	teacher := createTestTeacher(t, db, lineUserID)
	_ = teacher.ID

	authService := services.NewAuthService(appInstance)

	loginResponse, err, _ := authService.TeacherLineLogin(context.Background(), lineUserID, "mock-token")
	if err != nil {
		t.Fatalf("TeacherLineLogin failed: %v", err)
	}

	refreshResponse, refreshErr := authService.RefreshToken(context.Background(), loginResponse.Token)
	if refreshErr != nil {
		t.Fatalf("RefreshToken failed: %v", refreshErr)
	}

	if refreshResponse.Token == "" {
		t.Error("Expected token in refresh response")
	}
}

func TestAuthController_RefreshToken_InvalidToken(t *testing.T) {
	appInstance, _, cleanup := setupAuthTestApp()
	defer cleanup()

	authService := services.NewAuthService(appInstance)

	_, err := authService.RefreshToken(context.Background(), "invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestAuthController_Logout_Success(t *testing.T) {
	appInstance, _, cleanup := setupAuthTestApp()
	defer cleanup()

	authService := services.NewAuthService(appInstance)

	err := authService.Logout(context.Background(), "any-token")
	if err != nil {
		t.Errorf("Logout failed: %v", err)
	}
}

func TestAuthController_TokenValidation(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	uniqueID := generateUniqueID()
	lineUserID := fmt.Sprintf("U%s", uniqueID)
	teacher := createTestTeacher(t, db, lineUserID)

	authService := services.NewAuthService(appInstance)

	loginResponse, err, _ := authService.TeacherLineLogin(context.Background(), lineUserID, "mock-token")
	if err != nil {
		t.Fatalf("TeacherLineLogin failed: %v", err)
	}

	claims, tokenErr := authService.ValidateToken(loginResponse.Token)
	if tokenErr != nil {
		t.Fatalf("ValidateToken failed: %v", tokenErr)
	}

	if claims.UserID != teacher.ID {
		t.Errorf("Expected UserID %d, got %d", teacher.ID, claims.UserID)
	}

	if claims.UserType != "TEACHER" {
		t.Errorf("Expected UserType 'TEACHER', got '%s'", claims.UserType)
	}
}

func TestAuthController_AdminTokenValidation(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	authService := services.NewAuthService(appInstance)

	// 等待非同步寫入完成
	time.Sleep(100 * time.Millisecond)

	loginResponse, err := authService.AdminLogin(context.Background(), admin.Email, "password123", "127.0.0.1", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("AdminLogin failed: %v", err)
	}

	claims, err := authService.ValidateToken(loginResponse.Token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != admin.ID {
		t.Errorf("Expected UserID %d, got %d", admin.ID, claims.UserID)
	}

	if claims.UserType != "OWNER" {
		t.Errorf("Expected UserType 'OWNER', got '%s'", claims.UserType)
	}

	if claims.CenterID != admin.CenterID {
		t.Errorf("Expected CenterID %d, got %d", admin.CenterID, claims.CenterID)
	}
}

func TestAuthService_AdminLogin_SuccessRecordsHistory(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	authService := services.NewAuthService(appInstance)
	loginHistoryRepo := repositories.NewAdminLoginHistoryRepository(appInstance)

	// 等待非同步寫入完成
	time.Sleep(100 * time.Millisecond)

	// 執行成功登入
	loginResponse, err := authService.AdminLogin(
		context.Background(),
		admin.Email,
		"password123",
		"192.168.1.100",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	)
	if err != nil {
		t.Fatalf("AdminLogin failed: %v", err)
	}

	if loginResponse.Token == "" {
		t.Error("Expected token in response")
	}

	// 驗證登入紀錄已建立
	time.Sleep(200 * time.Millisecond) // 等待非同步寫入

	histories, err := loginHistoryRepo.GetByAdminID(context.Background(), admin.ID, 10)
	if err != nil {
		t.Fatalf("GetByAdminID failed: %v", err)
	}

	if len(histories) == 0 {
		t.Error("Expected at least one login history record")
		return
	}

	// 驗證最新的紀錄是成功的
	lastRecord := histories[len(histories)-1]
	if lastRecord.Status != models.LoginStatusSuccess {
		t.Errorf("Expected status 'SUCCESS', got '%s'", lastRecord.Status)
	}

	if lastRecord.IPAddress != "192.168.1.100" {
		t.Errorf("Expected IP '192.168.1.100', got '%s'", lastRecord.IPAddress)
	}

	if lastRecord.Email != admin.Email {
		t.Errorf("Expected email '%s', got '%s'", admin.Email, lastRecord.Email)
	}
}

func TestAuthService_AdminLogin_FailedPasswordRecordsHistory(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	authService := services.NewAuthService(appInstance)
	loginHistoryRepo := repositories.NewAdminLoginHistoryRepository(appInstance)

	// 執行失敗登入（密碼錯誤）
	_, err := authService.AdminLogin(
		context.Background(),
		admin.Email,
		"wrongpassword",
		"192.168.1.101",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
	)
	if err == nil {
		t.Error("Expected error for invalid password")
	}

	// 驗證登入紀錄已建立
	time.Sleep(200 * time.Millisecond) // 等待非同步寫入

	histories, err := loginHistoryRepo.GetByAdminID(context.Background(), admin.ID, 10)
	if err != nil {
		t.Fatalf("GetByAdminID failed: %v", err)
	}

	// 找到最新的失敗紀錄
	var lastFailedRecord *models.AdminLoginHistory
	for i := len(histories) - 1; i >= 0; i-- {
		if histories[i].Status == models.LoginStatusFailed {
			lastFailedRecord = &histories[i]
			break
		}
	}

	if lastFailedRecord == nil {
		t.Fatal("Expected at least one failed login history record")
	}

	if lastFailedRecord.Status != models.LoginStatusFailed {
		t.Errorf("Expected status 'FAILED', got '%s'", lastFailedRecord.Status)
	}

	if lastFailedRecord.Reason != "invalid password" {
		t.Errorf("Expected reason 'invalid password', got '%s'", lastFailedRecord.Reason)
	}

	if lastFailedRecord.IPAddress != "192.168.1.101" {
		t.Errorf("Expected IP '192.168.1.101', got '%s'", lastFailedRecord.IPAddress)
	}
}

func TestAuthService_AdminLogin_AdminNotFoundRecordsHistory(t *testing.T) {
	appInstance, _, cleanup := setupAuthTestApp()
	defer cleanup()

	authService := services.NewAuthService(appInstance)
	loginHistoryRepo := repositories.NewAdminLoginHistoryRepository(appInstance)

	// 執行失敗登入（帳號不存在）
	_, err := authService.AdminLogin(
		context.Background(),
		"nonexistent@example.com",
		"password123",
		"192.168.1.102",
		"Mozilla/5.0 (Linux; Android 10)",
	)
	if err == nil {
		t.Error("Expected error for non-existent admin")
	}

	// 驗證登入紀錄已建立
	time.Sleep(200 * time.Millisecond) // 等待非同步寫入

	histories, err := loginHistoryRepo.GetByAdminID(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("GetByAdminID failed: %v", err)
	}

	// 找到失敗紀錄（admin_id 為 0）
	var lastFailedRecord *models.AdminLoginHistory
	for i := len(histories) - 1; i >= 0; i-- {
		if histories[i].Email == "nonexistent@example.com" {
			lastFailedRecord = &histories[i]
			break
		}
	}

	if lastFailedRecord == nil {
		t.Fatal("Expected login history record for non-existent admin")
	}

	if lastFailedRecord.Status != models.LoginStatusFailed {
		t.Errorf("Expected status 'FAILED', got '%s'", lastFailedRecord.Status)
	}

	if lastFailedRecord.AdminID != 0 {
		t.Errorf("Expected AdminID 0, got %d", lastFailedRecord.AdminID)
	}

	if lastFailedRecord.Reason != "admin not found" {
		t.Errorf("Expected reason 'admin not found', got '%s'", lastFailedRecord.Reason)
	}
}

func TestAdminLoginHistoryRepository_GetRecentFailedLogins(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	loginHistoryRepo := repositories.NewAdminLoginHistoryRepository(appInstance)

	// 建立一些失敗的登入紀錄
	failedHistories := []models.AdminLoginHistory{
		{
			AdminID:   admin.ID,
			Email:     admin.Email,
			Status:    models.LoginStatusFailed,
			IPAddress: "10.0.0.1",
			Reason:    "wrong password",
			CreatedAt: time.Now().Add(-1 * time.Hour), // 1小時前
		},
		{
			AdminID:   admin.ID,
			Email:     admin.Email,
			Status:    models.LoginStatusFailed,
			IPAddress: "10.0.0.2",
			Reason:    "wrong password",
			CreatedAt: time.Now().Add(-30 * time.Minute), // 30分鐘前
		},
	}

	for _, history := range failedHistories {
		if err := db.Create(&history).Error; err != nil {
			t.Fatalf("Failed to create test login history: %v", err)
		}
	}

	// 查詢最近30分鐘內的失敗登入
	since := time.Now().Add(-45 * time.Minute)
	recentLogins, err := loginHistoryRepo.GetRecentFailedLogins(context.Background(), admin.ID, since)
	if err != nil {
		t.Fatalf("GetRecentFailedLogins failed: %v", err)
	}

	// 應該只有1筆記錄（30分鐘前那筆）
	if len(recentLogins) != 1 {
		t.Errorf("Expected 1 recent failed login, got %d", len(recentLogins))
	}

	if len(recentLogins) > 0 && recentLogins[0].IPAddress != "10.0.0.2" {
		t.Errorf("Expected IP '10.0.0.2', got '%s'", recentLogins[0].IPAddress)
	}
}

func TestAdminLoginHistoryRepository_CountFailedLoginsSince(t *testing.T) {
	appInstance, db, cleanup := setupAuthTestApp()
	defer cleanup()

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	loginHistoryRepo := repositories.NewAdminLoginHistoryRepository(appInstance)

	// 建立失敗的登入紀錄
	failedHistories := []models.AdminLoginHistory{
		{
			AdminID:   admin.ID,
			Email:     admin.Email,
			Status:    models.LoginStatusFailed,
			IPAddress: "10.0.0.1",
			Reason:    "wrong password",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			AdminID:   admin.ID,
			Email:     admin.Email,
			Status:    models.LoginStatusFailed,
			IPAddress: "10.0.0.2",
			Reason:    "wrong password",
			CreatedAt: time.Now().Add(-10 * time.Minute),
		},
	}

	for _, history := range failedHistories {
		if err := db.Create(&history).Error; err != nil {
			t.Fatalf("Failed to create test login history: %v", err)
		}
	}

	// 計算最近30分鐘內的失敗登入次數
	since := time.Now().Add(-30 * time.Minute)
	count := loginHistoryRepo.CountFailedLoginsSince(context.Background(), admin.ID, since)

	if count != 1 {
		t.Errorf("Expected 1 failed login in last 30 minutes, got %d", count)
	}
}
