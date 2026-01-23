package test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	dbmysql "timeLedger/database/mysql"
)

// Generate unique identifier for test data
func generateUniqueID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(10000))
}

func setupAuthTestDB(t *testing.T) *gorm.DB {
	db, err := InitializeTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Center{},
		&models.AdminUser{},
		&models.Teacher{},
	); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
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
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)
	adminRepo := repositories.NewAdminUserRepository(testApp)

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

	if response.User.UserType != "ADMIN" {
		t.Errorf("Expected UserType 'ADMIN', got '%s'", response.User.UserType)
	}

	_ = adminRepo
}

func TestAuthController_AdminLogin_InvalidPassword(t *testing.T) {
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

	_, err := authService.AdminLogin(context.Background(), admin.Email, "wrongpassword")
	if err == nil {
		t.Error("Expected error for invalid password")
	}
}

func TestAuthController_AdminLogin_AdminNotFound(t *testing.T) {
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	_ = createTestCenter(t, db)

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

	_, err := authService.AdminLogin(context.Background(), "nonexistent@test.com", "password123")
	if err == nil {
		t.Error("Expected error for nonexistent admin")
	}
}

func TestAuthController_TeacherLineLogin_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	uniqueID := generateUniqueID()
	lineUserID := fmt.Sprintf("U%s", uniqueID)
	teacher := createTestTeacher(t, db, lineUserID)
	_ = teacher.ID

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

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
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

	_, err := authService.TeacherLineLogin(context.Background(), "UNKNOWN_USER", "mock-token")
	if err == nil {
		t.Error("Expected error for nonexistent teacher")
	}
}

func TestAuthController_RefreshToken_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	uniqueID := generateUniqueID()
	lineUserID := fmt.Sprintf("U%s", uniqueID)
	teacher := createTestTeacher(t, db, lineUserID)
	_ = teacher.ID

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

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
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

	_, err := authService.RefreshToken(context.Background(), "invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestAuthController_Logout_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

	err := authService.Logout(context.Background(), "any-token")
	if err != nil {
		t.Errorf("Logout failed: %v", err)
	}
}

func TestAuthController_TokenValidation(t *testing.T) {
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	uniqueID := generateUniqueID()
	lineUserID := fmt.Sprintf("U%s", uniqueID)
	teacher := createTestTeacher(t, db, lineUserID)

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

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
	db := setupAuthTestDB(t)
	defer CloseDB(db)

	center := createTestCenter(t, db)
	admin := createTestAdminUser(t, db, center.ID, "password123")

	testApp := &app.App{
		MySQL: &dbmysql.DB{WDB: db, RDB: db},
	}

	authService := services.NewAuthService(testApp)

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

	if claims.UserType != "ADMIN" {
		t.Errorf("Expected UserType 'ADMIN', got '%s'", claims.UserType)
	}

	if claims.CenterID != admin.CenterID {
		t.Errorf("Expected CenterID %d, got %d", admin.CenterID, claims.CenterID)
	}
}
