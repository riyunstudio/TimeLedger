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
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global"
	"timeLedger/global/errInfos"
	mockRedis "timeLedger/testing/redis"

	"github.com/gin-gonic/gin"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupAdminUserTestApp() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
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

func TestAdminUserController_CRUD(t *testing.T) {
	appInstance, _, cleanup := setupAdminUserTestApp()
	defer cleanup()

	ctx := context.Background()

	center := models.Center{
		Name:      "Test Center",
		PlanLevel: "STARTER",
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, err := centerRepo.Create(ctx, center)
	if err != nil {
		t.Fatalf("Failed to create center: %v", err)
	}

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)

	adminUser := models.AdminUser{
		Email:        "test@adminuser.com",
		PasswordHash: "hashed_password",
		Name:         "Test Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	createdAdminUser, err := adminUserRepo.Create(ctx, adminUser)
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	adminUserController := controllers.NewAdminUserController(appInstance)

	t.Run("GetAdminUsers_Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdminUser.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users", nil)

		adminUserController.GetAdminUsers(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["code"] != float64(0) {
			t.Errorf("Expected code 0, got %v", response["code"])
		}
	})

	t.Run("CreateAdminUser_Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdminUser.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		reqBody := controllers.CreateAdminUserRequest{
			Email:    "newadmin@test.com",
			Password: "newpassword123",
			Name:     "New Admin",
			Role:     "ADMIN",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminUserController.CreateAdminUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("CreateAdminUser_InvalidRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdminUser.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		adminUserController.CreateAdminUser(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("UpdateAdminUser_Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdminUser.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{
			{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)},
			{Key: "admin_id", Value: fmt.Sprintf("%d", createdAdminUser.ID)},
		}

		reqBody := controllers.UpdateAdminUserRequest{
			Name: "Updated Admin Name",
			Role: "ADMIN",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PUT", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users/"+fmt.Sprintf("%d", createdAdminUser.ID), bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminUserController.UpdateAdminUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("DeleteAdminUser_Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdminUser.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{
			{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)},
			{Key: "admin_id", Value: fmt.Sprintf("%d", createdAdminUser.ID)},
		}

		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users/"+fmt.Sprintf("%d", createdAdminUser.ID), nil)

		adminUserController.DeleteAdminUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GetAdminUsers_InvalidCenterID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdminUser.ID)
		c.Set(global.CenterIDKey, 999)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/999/users", nil)

		adminUserController.GetAdminUsers(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}
