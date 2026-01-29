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

	response, err := authService.AdminLogin(context.Background(), admin.Email, "password123")
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

	_, err := authService.AdminLogin(context.Background(), admin.Email, "wrongpassword")
	if err == nil {
		t.Error("Expected error for invalid password")
	}
}

func TestAuthController_AdminLogin_AdminNotFound(t *testing.T) {
	appInstance, _, cleanup := setupAuthTestApp()
	defer cleanup()

	_ = createTestCenter(t, appInstance.MySQL.RDB)

	authService := services.NewAuthService(appInstance)

	_, err := authService.AdminLogin(context.Background(), "nonexistent@test.com", "password123")
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

	response, err := authService.TeacherLineLogin(context.Background(), lineUserID, "mock-token")
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

	_, err := authService.TeacherLineLogin(context.Background(), "UNKNOWN_USER", "mock-token")
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

	loginResponse, err := authService.TeacherLineLogin(context.Background(), lineUserID, "mock-token")
	if err != nil {
		t.Fatalf("TeacherLineLogin failed: %v", err)
	}

	refreshResponse, err := authService.RefreshToken(context.Background(), loginResponse.Token)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
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

	loginResponse, err := authService.TeacherLineLogin(context.Background(), lineUserID, "mock-token")
	if err != nil {
		t.Fatalf("TeacherLineLogin failed: %v", err)
	}

	claims, err := authService.ValidateToken(loginResponse.Token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
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

	loginResponse, err := authService.AdminLogin(context.Background(), admin.Email, "password123")
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
