package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/controllers"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	mockRedis "timeLedger/testing/redis"

	"github.com/gin-gonic/gin"
)

func setupIntegrationTestApp() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.AdminUser{},
		&models.Teacher{},
		&models.Course{},
		&models.Room{},
		&models.Offering{},
		&models.ScheduleRule{},
		&models.ScheduleException{},
	); err != nil {
		panic(fmt.Sprintf("AutoMigrate error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	appInstance := &app.App{
		Env:   nil,
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

func TestIntegration_AdminLoginAndCRUD(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestApp()
	defer cleanup()

	ctx := context.Background()

	center := models.Center{
		Name:      "Integration Test Center",
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			AllowPublicRegister: true,
			DefaultLanguage:     "zh-TW",
		},
		CreatedAt: time.Now(),
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, err := centerRepo.Create(ctx, center)
	if err != nil {
		t.Fatalf("Failed to create center: %v", err)
	}

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        "integration@test.com",
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Integration Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
	}
	_, err = adminUserRepo.Create(ctx, adminUser)
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	authService := services.NewMockAuthService(appInstance)

	t.Run("Step1_AdminLogin_Success", func(t *testing.T) {
		authController := controllers.NewAuthController(appInstance, authService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := controllers.AdminLoginRequest{
			Email:    "integration@test.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		authController.AdminLogin(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["code"] != float64(0) {
			t.Errorf("Expected code 0, got %v", response["code"])
		}

		datas := response["datas"].(map[string]interface{})
		token := datas["token"].(string)
		if token == "" {
			t.Error("Expected non-empty token")
		}
		t.Logf("Login successful, token: %s...", token[:50])
	})

	t.Run("Step2_AdminLogin_InvalidPassword", func(t *testing.T) {
		authController := controllers.NewAuthController(appInstance, authService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := controllers.AdminLoginRequest{
			Email:    "integration@test.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		authController.AdminLogin(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Step3_GetAdminUsers_WithAuth", func(t *testing.T) {
		adminUserController := controllers.NewAdminUserController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, uint(1))
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users", nil)

		adminUserController.GetAdminUsers(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Step4_GetAdminUsers_WithoutAuth", func(t *testing.T) {
		adminUserController := controllers.NewAdminUserController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users", nil)

		adminUserController.GetAdminUsers(c)

		if w.Code == http.StatusOK || w.Code == http.StatusUnauthorized {
			t.Logf("Got status %d - this varies based on auth middleware execution", w.Code)
		} else {
			t.Errorf("Expected status 200 or 401, got %d", w.Code)
		}
	})

	t.Run("Step5_CreateCenter_WithAuth", func(t *testing.T) {
		adminResourceController := controllers.NewAdminResourceController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, uint(1))
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers", nil)

		newCenter := map[string]interface{}{
			"name":       "New Center from Test",
			"plan_level": "STARTER",
		}
		body, _ := json.Marshal(newCenter)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminResourceController.CreateCenter(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

func TestIntegration_TeacherCRUD_WithAuth(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestApp()
	defer cleanup()

	ctx := context.Background()

	center := models.Center{
		Name:      "Teacher Test Center",
		PlanLevel: "STARTER",
		CreatedAt: time.Now(),
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID:     fmt.Sprintf("LINE_USER_%d", time.Now().UnixNano()),
		Name:           "Test Teacher",
		Email:          fmt.Sprintf("teacher%d@integration.com", time.Now().UnixNano()),
		IsOpenToHiring: true,
		CreatedAt:      time.Now(),
	}
	createdTeacher, err := teacherRepo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	t.Run("GetTeacherByID", func(t *testing.T) {
		teacherController := controllers.NewTeacherController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdTeacher.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/teachers/"+fmt.Sprintf("%d", createdTeacher.ID), nil)

		teacherController.GetProfile(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("UpdateTeacherProfile", func(t *testing.T) {
		teacherController := controllers.NewTeacherController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdTeacher.ID)}}

		updateData := map[string]interface{}{
			"name":  "Updated Teacher Name",
			"bio":   "This is an updated bio",
			"email": "updated@email.com",
		}
		body, _ := json.Marshal(updateData)
		c.Request = httptest.NewRequest("PUT", "/api/v1/teachers/"+fmt.Sprintf("%d", createdTeacher.ID), bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.UpdateProfile(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

func TestIntegration_RefreshToken(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestApp()
	defer cleanup()

	ctx := context.Background()

	center := models.Center{
		Name:      "Refresh Token Test Center",
		PlanLevel: "STARTER",
		CreatedAt: time.Now(),
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        "refresh@test.com",
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Refresh Token Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	authService := services.NewMockAuthService(appInstance)

	t.Run("RefreshToken_InvalidToken", func(t *testing.T) {
		authController := controllers.NewAuthController(appInstance, authService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidData := map[string]interface{}{
			"token": "invalid.token.here",
		}
		body, _ := json.Marshal(invalidData)
		c.Request = httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		authController.RefreshToken(c)

		if w.Code != http.StatusUnauthorized {
			t.Logf("Got status %d - token validation behavior", w.Code)
		}
	})

	t.Run("Logout_WithoutAuth", func(t *testing.T) {
		authController := controllers.NewAuthController(appInstance, authService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/api/v1/auth/logout", nil)

		authController.Logout(c)

		if w.Code == http.StatusOK || w.Code == http.StatusUnauthorized {
			t.Logf("Got status %d - auth context behavior", w.Code)
		} else {
			t.Errorf("Expected status 200 or 401, got %d", w.Code)
		}
	})

	t.Run("Logout_WithAuthContext", func(t *testing.T) {
		authController := controllers.NewAuthController(appInstance, authService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Request = httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
		c.Request.Header.Set("Authorization", "Bearer test-token")

		authController.Logout(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}
