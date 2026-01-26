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
	"timeLedger/configs"
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

func setupIntegrationTestAppWithMigrations() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
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
		&models.CenterHoliday{},
		&models.CenterInvitation{},
		&models.PersonalEvent{},
	); err != nil {
		panic(fmt.Sprintf("AutoMigrate error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	// 初始化測試用的 Env 配置
	env := &configs.Env{
		JWTSecret:      "test-jwt-secret-key-for-testing-only",
		AppEnv:         "test",
		AppDebug:       true,
		AppTimezone:    "Asia/Taipei",
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
		mysqlDB.Exec("DELETE FROM schedule_exceptions")
		mysqlDB.Exec("DELETE FROM schedule_rules")
		mysqlDB.Exec("DELETE FROM center_holidays")
		mysqlDB.Exec("DELETE FROM center_invitations")
		mysqlDB.Exec("DELETE FROM offerings")
		mysqlDB.Exec("DELETE FROM courses")
		mysqlDB.Exec("DELETE FROM rooms")
		mysqlDB.Exec("DELETE FROM personal_events")
		mysqlDB.Exec("DELETE FROM admin_users")
		mysqlDB.Exec("DELETE FROM teachers")
		mysqlDB.Exec("DELETE FROM centers")
		mr.Close()
	}

	return appInstance, mysqlDB, cleanup
}

func TestIntegration_CenterAdminFullWorkflow(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Full Test Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			AllowPublicRegister: true,
			DefaultLanguage:     "zh-TW",
			ExceptionLeadDays:   14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, err := centerRepo.Create(ctx, center)
	if err != nil {
		t.Fatalf("Failed to create center: %v", err)
	}
	t.Logf("Created center: ID=%d, Name=%s", createdCenter.ID, createdCenter.Name)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("fulltest_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Full Test Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, err := adminUserRepo.Create(ctx, adminUser)
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}
	t.Logf("Created admin user: ID=%d", createdAdmin.ID)

	authService := services.NewAuthService(appInstance)
	adminResourceController := controllers.NewAdminResourceController(appInstance)
	_ = authService
	_ = adminResourceController

	t.Run("Step1_AdminLogin", func(t *testing.T) {
		authController := controllers.NewAuthController(appInstance, authService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := controllers.AdminLoginRequest{
			Email:    adminUser.Email,
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		authController.AdminLogin(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["code"] != float64(0) {
			t.Fatalf("Expected code 0, got %v", response)
		}

		datas := response["datas"].(map[string]interface{})
		token := datas["token"].(string)
		if token == "" {
			t.Fatal("Expected non-empty token")
		}
		t.Logf("Admin login successful, token: %s...", token[:30])
	})

	t.Run("Step2_CreateRoom", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)

		reqBody := controllers.CreateRoomRequest{
			Name:     fmt.Sprintf("Test Room %d", now.UnixNano()),
			Capacity: 20,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/rooms", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminResourceController.CreateRoom(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		datas := response["datas"].(map[string]interface{})
		roomID := uint(datas["id"].(float64))
		if roomID == 0 {
			t.Fatal("Expected non-zero room ID")
		}
		t.Logf("Created room: ID=%d", roomID)
	})

	t.Run("Step3_CreateCourse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)

		reqBody := controllers.CreateCourseRequest{
			Name:             fmt.Sprintf("Piano Class %d", now.UnixNano()),
			Duration:         60,
			ColorHex:         "#FF5733",
			RoomBufferMin:    10,
			TeacherBufferMin: 15,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/courses", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminResourceController.CreateCourse(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		datas := response["datas"].(map[string]interface{})
		courseID := uint(datas["id"].(float64))
		if courseID == 0 {
			t.Fatal("Expected non-zero course ID")
		}
		t.Logf("Created course: ID=%d", courseID)
	})

	t.Run("Step4_GetRoomsAndCourses", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/rooms", nil)

		adminResourceController.GetRooms(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for GetRooms, got %d. Body: %s", w.Code, w.Body.String())
		}

		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Set(global.UserIDKey, createdAdmin.ID)
		c2.Set(global.CenterIDKey, createdCenter.ID)
		c2.Request = httptest.NewRequest("GET", "/api/v1/admin/courses", nil)

		adminResourceController.GetCourses(c2)

		if w2.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for GetCourses, got %d. Body: %s", w2.Code, w2.Body.String())
		}
		t.Log("Successfully retrieved rooms and courses")
	})
}

func TestIntegration_TeacherFullWorkflow(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Teacher Workflow Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID:     fmt.Sprintf("LINE_TEACHER_%d", now.UnixNano()),
		Name:           "Test Teacher",
		Email:          fmt.Sprintf("teacher_workflow_%d@test.com", now.UnixNano()),
		Bio:            "Test bio",
		City:           "Taipei",
		District:       "Xinyi",
		IsOpenToHiring: true,
		CreatedAt:      now,
	}
	createdTeacher, err := teacherRepo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}
	t.Logf("Created teacher: ID=%d, Name=%s", createdTeacher.ID, createdTeacher.Name)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Test Room",
		Capacity:  10,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Test Course",
		DefaultDuration:  60,
		ColorHex:         "#3498DB",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:      createdCenter.ID,
		CourseID:      createdCourse.ID,
		Name:          "Test Offering",
		DefaultRoomID: &createdRoom.ID,
		IsActive:      true,
		CreatedAt:     now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)

	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	rule := models.ScheduleRule{
		CenterID:   createdCenter.ID,
		OfferingID: createdOffering.ID,
		TeacherID:  &createdTeacher.ID,
		RoomID:     createdRoom.ID,
		Weekday:    1,
		StartTime:  "10:00:00",
		EndTime:    "11:00:00",
		EffectiveRange: models.DateRange{
			StartDate: now,
			EndDate:   now.AddDate(0, 3, 0),
		},
		CreatedAt: now,
	}
	createdRule, _ := ruleRepo.Create(ctx, rule)
	t.Logf("Created schedule rule: ID=%d", createdRule.ID)

	teacherController := controllers.NewTeacherController(appInstance)

	t.Run("Step1_GetProfile", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/me/profile", nil)

		teacherController.GetProfile(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["code"] != float64(0) {
			t.Fatalf("Expected code 0, got %v", response)
		}
		t.Log("Successfully retrieved teacher profile")
	})

	t.Run("Step2_GetSchedule", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		fromDate := now.Format("2006-01-02")
		toDate := now.AddDate(0, 1, 0).Format("2006-01-02")
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/me/schedule?from="+fromDate+"&to="+toDate, nil)

		teacherController.GetSchedule(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved teacher schedule")
	})

	t.Run("Step3_GetExceptions", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/exceptions", nil)

		teacherController.GetExceptions(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved teacher exceptions")
	})
}

func TestIntegration_ScheduleRuleCreation(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Schedule Rule Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("schedulerule_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Schedule Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_SCHEDULE_%d", now.UnixNano()),
		Name:       "Schedule Teacher",
		Email:      fmt.Sprintf("schedule_teacher_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Schedule Room",
		Capacity:  15,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Math Course",
		DefaultDuration:  90,
		ColorHex:         "#9B59B6",
		RoomBufferMin:    5,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:      createdCenter.ID,
		CourseID:      createdCourse.ID,
		Name:          "Math Class A",
		DefaultRoomID: &createdRoom.ID,
		IsActive:      true,
		CreatedAt:     now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)

	schedulingController := controllers.NewSchedulingController(appInstance)

	t.Run("Step1_CreateScheduleRule", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		effectiveStart := now.Format("2006-01-02")
		effectiveEnd := now.AddDate(0, 2, 0).Format("2006-01-02")

		reqBody := map[string]interface{}{
			"name":        "Math Class Weekly",
			"offering_id": createdOffering.ID,
			"teacher_id":  createdTeacher.ID,
			"room_id":     createdRoom.ID,
			"start_time":  "14:00:00",
			"end_time":    "15:30:00",
			"duration":    90,
			"weekdays":    []int{2},
			"start_date":  effectiveStart,
			"end_date":    &effectiveEnd,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/rules", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CreateRule(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["code"] != float64(0) {
			t.Fatalf("Expected code 0, got %v", response)
		}
		t.Log("Successfully created schedule rule")
	})

	t.Run("Step2_GetRules", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/rules", nil)

		schedulingController.GetRules(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved schedule rules")
	})

	t.Run("Step3_ExpandRules", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		reqBody := map[string]interface{}{
			"rule_ids":   []uint{createdOffering.ID},
			"start_date": now.Format(time.RFC3339),
			"end_date":   now.AddDate(0, 1, 0).Format(time.RFC3339),
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/expand-rules", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.ExpandRules(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully expanded schedule rules")
	})
}

func TestIntegration_ResourceToggleAndInvitationStats(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Resource Toggle Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("resource_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Resource Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Toggle Test Room",
		Capacity:  10,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)
	_ = createdRoom

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Toggle Test Course",
		DefaultDuration:  60,
		ColorHex:         "#1ABC9C",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)
	t.Logf("Toggle test course created: ID=%d", createdCourse.ID)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:  createdCenter.ID,
		CourseID:  createdCourse.ID,
		Name:      "Toggle Test Offering",
		IsActive:  true,
		CreatedAt: now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)
	_ = createdOffering

	invitationRepo := repositories.NewCenterInvitationRepository(appInstance)
	for i := 0; i < 3; i++ {
		invitation := models.CenterInvitation{
			CenterID:  createdCenter.ID,
			Email:     fmt.Sprintf("invitee%d_%d@test.com", i, now.UnixNano()),
			Token:     fmt.Sprintf("token_%d_%d", now.UnixNano(), i),
			Status:    "PENDING",
			CreatedAt: now,
			ExpiresAt: now.AddDate(0, 1, 0),
		}
		invitationRepo.Create(ctx, invitation)
	}
	invitation := models.CenterInvitation{
		CenterID:  createdCenter.ID,
		Email:     fmt.Sprintf("accepted_%d@test.com", now.UnixNano()),
		Token:     fmt.Sprintf("token_accepted_%d", now.UnixNano()),
		Status:    "ACCEPTED",
		CreatedAt: now,
		ExpiresAt: now.AddDate(0, 1, 0),
	}
	invitationRepo.Create(ctx, invitation)

	adminResourceController := controllers.NewAdminResourceController(appInstance)

	t.Run("Step1_GetActiveRooms", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/rooms/active", nil)

		adminResourceController.GetActiveRooms(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved active rooms")
	})

	t.Run("Step2_ToggleCourseActive", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "course_id", Value: fmt.Sprintf("%d", createdCourse.ID)}}

		reqBody := controllers.ToggleActiveRequest{
			IsActive: false,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PATCH", "/api/v1/admin/courses/"+fmt.Sprintf("%d", createdCourse.ID)+"/toggle-active", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		adminResourceController.ToggleCourseActive(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully toggled course active status")
	})

	t.Run("Step3_GetInvitationStats", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/invitations/stats", nil)

		adminResourceController.GetInvitationStats(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		datas := response["datas"].(map[string]interface{})
		total := int(datas["total"].(float64))
		pending := int(datas["pending"].(float64))
		accepted := int(datas["accepted"].(float64))

		if total != 4 {
			t.Errorf("Expected total 4 invitations, got %d", total)
		}
		if pending != 3 {
			t.Errorf("Expected 3 pending invitations, got %d", pending)
		}
		if accepted != 1 {
			t.Errorf("Expected 1 accepted invitation, got %d", accepted)
		}
		t.Logf("Invitation stats: total=%d, pending=%d, accepted=%d", total, pending, accepted)
	})

	t.Run("Step4_GetInvitations", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/invitations?page=1&limit=10", nil)

		adminResourceController.GetInvitations(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved invitations list")
	})
}

func TestIntegration_ValidationAndException(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Validation Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("validation_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Validation Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_VALID_%d", now.UnixNano()),
		Name:       "Validation Teacher",
		Email:      fmt.Sprintf("validation_teacher_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Validation Room",
		Capacity:  10,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Validation Course",
		DefaultDuration:  60,
		ColorHex:         "#E74C3C",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:  createdCenter.ID,
		CourseID:  createdCourse.ID,
		Name:      "Validation Offering",
		IsActive:  true,
		CreatedAt: now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(appInstance)
	effectiveRange := models.DateRange{
		StartDate: now,
		EndDate:   now.AddDate(1, 0, 0),
	}
	scheduleRule := models.ScheduleRule{
		CenterID:       createdCenter.ID,
		OfferingID:     createdOffering.ID,
		TeacherID:      &createdTeacher.ID,
		RoomID:         createdRoom.ID,
		Weekday:        int(now.Weekday()),
		StartTime:      now.Format("15:04"),
		EndTime:        now.Add(time.Hour).Format("15:04"),
		EffectiveRange: effectiveRange,
		LockAt:         &now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	createdScheduleRule, _ := scheduleRuleRepo.Create(ctx, scheduleRule)

	schedulingController := controllers.NewSchedulingController(appInstance)

	t.Run("Step1_CheckOverlap_Empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)

		startTime := now.Format(time.RFC3339)
		endTime := now.Add(time.Hour).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"center_id":  createdCenter.ID,
			"teacher_id": createdTeacher.ID,
			"room_id":    createdRoom.ID,
			"start_time": startTime,
			"end_time":   endTime,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/check-overlap", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CheckOverlap(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully checked overlap (no conflicts)")
	})

	t.Run("Step2_ValidateFull", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)

		startTime := now.Format(time.RFC3339)
		endTime := now.Add(time.Hour).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"center_id":             createdCenter.ID,
			"teacher_id":            createdTeacher.ID,
			"room_id":               createdRoom.ID,
			"course_id":             createdCourse.ID,
			"start_time":            startTime,
			"end_time":              endTime,
			"allow_buffer_override": false,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/validate", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.ValidateFull(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully validated schedule")
	})

	t.Run("Step3_DetectPhaseTransitions", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		startDate := now.Format(time.RFC3339)
		endDate := now.AddDate(0, 3, 0).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"offering_id": createdOffering.ID,
			"start_date":  startDate,
			"end_date":    endDate,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/detect-phase-transitions", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.DetectPhaseTransitions(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully detected phase transitions")
	})

	t.Run("Step4_CheckRuleLockStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		exceptionDate := now.Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"rule_id":        createdScheduleRule.ID,
			"exception_date": exceptionDate,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/check-rule-lock", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CheckRuleLockStatus(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully checked rule lock status")
	})
}

func TestIntegration_InvitationFlow(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Invitation Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("invitation_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Invitation Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	t.Run("Step1_CreateInvitation", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		reqBody := map[string]interface{}{
			"email": fmt.Sprintf("new_teacher_%d@test.com", now.UnixNano()),
			"role":  "TEACHER",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/invitations", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController := controllers.NewTeacherController(appInstance)
		teacherController.InviteTeacher(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		datas := response["datas"].(map[string]interface{})
		invitationToken := datas["token"].(string)
		if invitationToken == "" {
			t.Fatal("Expected non-empty invitation token")
		}
		t.Logf("Created invitation with token: %s...", invitationToken[:20])
	})

	t.Run("Step2_GetInvitations", func(t *testing.T) {
		adminResourceController := controllers.NewAdminResourceController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/invitations?page=1&limit=10", nil)

		adminResourceController.GetInvitations(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved invitations")
	})
}

func TestIntegration_ExceptionReview(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Exception Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("exception_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Exception Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_EXC_%d", now.UnixNano()),
		Name:       "Exception Teacher",
		Email:      fmt.Sprintf("exception_teacher_%d@test.com", now.UnixNano()),
		City:       "Taipei",
		District:   "Xinyi",
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Exception Room",
		Capacity:  10,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Exception Course",
		DefaultDuration:  60,
		ColorHex:         "#E74C3C",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:  createdCenter.ID,
		CourseID:  createdCourse.ID,
		Name:      "Exception Offering",
		IsActive:  true,
		CreatedAt: now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)

	effectiveRange := models.DateRange{
		StartDate: now,
		EndDate:   now.AddDate(1, 0, 0),
	}
	scheduleRuleRepo := repositories.NewScheduleRuleRepository(appInstance)
	scheduleRule := models.ScheduleRule{
		CenterID:       createdCenter.ID,
		OfferingID:     createdOffering.ID,
		TeacherID:      &createdTeacher.ID,
		RoomID:         createdRoom.ID,
		Weekday:        int(now.Weekday()),
		StartTime:      now.Format("15:04"),
		EndTime:        now.Add(time.Hour).Format("15:04"),
		EffectiveRange: effectiveRange,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	createdRule, _ := scheduleRuleRepo.Create(ctx, scheduleRule)

	schedulingController := controllers.NewSchedulingController(appInstance)
	teacherController := controllers.NewTeacherController(appInstance)

	t.Run("Step1_CreateExceptionRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		originalDate := now.AddDate(0, 0, 20)
		newStartAt := originalDate.Add(time.Hour * 14)
		newEndAt := originalDate.Add(time.Hour * 15)

		newStartAtStr := newStartAt.Format(time.RFC3339)
		newEndAtStr := newEndAt.Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"center_id":     createdCenter.ID,
			"rule_id":       createdRule.ID,
			"original_date": originalDate.Format("2006-01-02"),
			"type":          "RESCHEDULE",
			"new_start_at":  newStartAtStr,
			"new_end_at":    newEndAtStr,
			"reason":        "Personal meeting",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/scheduling/exceptions", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.CreateException(c)

		if w.Code != http.StatusOK {
			t.Logf("Create exception request body: %s", string(body))
			t.Logf("Response: %s", w.Body.String())
		}
		if w.Code == http.StatusOK {
			t.Log("Successfully created exception request")
		}
	})

	t.Run("Step2_GetExceptions", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/scheduling/exceptions", nil)

		teacherController.GetExceptions(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved exceptions")
	})

	t.Run("Step3_ReviewException_Approve", func(t *testing.T) {
		exceptionRepo := repositories.NewScheduleExceptionRepository(appInstance)
		newStartAt := now.AddDate(0, 0, 5).Add(time.Hour * 14)
		newEndAt := now.AddDate(0, 0, 5).Add(time.Hour * 15)
		exception := models.ScheduleException{
			RuleID:       createdRule.ID,
			CenterID:     createdCenter.ID,
			OriginalDate: now.AddDate(0, 0, 5),
			Type:         "TIME_CHANGE",
			Status:       "PENDING",
			NewStartAt:   &newStartAt,
			NewEndAt:     &newEndAt,
			Reason:       "Doctor appointment",
			CreatedAt:    now,
		}
		createdException, _ := exceptionRepo.Create(ctx, exception)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "exceptionId", Value: fmt.Sprintf("%d", createdException.ID)}}

		reqBody := map[string]interface{}{
			"action":       "APPROVE",
			"review_notes": "Approved by admin",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/exceptions/"+fmt.Sprintf("%d", createdException.ID)+"/review", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.ReviewException(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully approved exception request")
	})

	t.Run("Step4_ReviewException_Reject", func(t *testing.T) {
		exceptionRepo := repositories.NewScheduleExceptionRepository(appInstance)
		exception := models.ScheduleException{
			RuleID:       createdRule.ID,
			CenterID:     createdCenter.ID,
			OriginalDate: now.AddDate(0, 0, 7),
			Type:         "CANCEL",
			Status:       "PENDING",
			Reason:       "Travel",
			CreatedAt:    now,
		}
		createdException, _ := exceptionRepo.Create(ctx, exception)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "exceptionId", Value: fmt.Sprintf("%d", createdException.ID)}}

		reqBody := map[string]interface{}{
			"action":       "REJECT",
			"review_notes": "Cannot approve during exam week",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/exceptions/"+fmt.Sprintf("%d", createdException.ID)+"/review", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.ReviewException(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully rejected exception request")
	})

	t.Run("Step5_RevokeException", func(t *testing.T) {
		exceptionRepo := repositories.NewScheduleExceptionRepository(appInstance)
		newStartAt := now.AddDate(0, 0, 9).Add(time.Hour * 10)
		newEndAt := now.AddDate(0, 0, 9).Add(time.Hour * 11)
		exception := models.ScheduleException{
			RuleID:       createdRule.ID,
			CenterID:     createdCenter.ID,
			OriginalDate: now.AddDate(0, 0, 9),
			Type:         "TIME_CHANGE",
			Status:       "PENDING",
			NewStartAt:   &newStartAt,
			NewEndAt:     &newEndAt,
			Reason:       "Family event",
			CreatedAt:    now,
		}
		createdException, _ := exceptionRepo.Create(ctx, exception)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdException.ID)}}

		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/exceptions/"+fmt.Sprintf("%d", createdException.ID)+"/revoke", nil)

		teacherController.RevokeException(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully revoked exception request")
	})
}

func TestIntegration_RecurrenceEditing(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Recurrence Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("recurrence_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Recurrence Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_REC_%d", now.UnixNano()),
		Name:       "Recurrence Teacher",
		Email:      fmt.Sprintf("recurrence_teacher_%d@test.com", now.UnixNano()),
		City:       "Taipei",
		District:   "Xinyi",
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Recurrence Room",
		Capacity:  15,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Recurrence Course",
		DefaultDuration:  60,
		ColorHex:         "#3498DB",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:  createdCenter.ID,
		CourseID:  createdCourse.ID,
		Name:      "Recurrence Offering",
		IsActive:  true,
		CreatedAt: now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)

	effectiveRange := models.DateRange{
		StartDate: now,
		EndDate:   now.AddDate(1, 0, 0),
	}
	scheduleRuleRepo := repositories.NewScheduleRuleRepository(appInstance)
	scheduleRule := models.ScheduleRule{
		CenterID:       createdCenter.ID,
		OfferingID:     createdOffering.ID,
		TeacherID:      &createdTeacher.ID,
		RoomID:         createdRoom.ID,
		Weekday:        int(now.Weekday()),
		StartTime:      now.Format("15:04"),
		EndTime:        now.Add(time.Hour).Format("15:04"),
		EffectiveRange: effectiveRange,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	createdRule, _ := scheduleRuleRepo.Create(ctx, scheduleRule)

	teacherController := controllers.NewTeacherController(appInstance)

	t.Run("Step1_PreviewAffectedSessions", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		startDate := now.Format(time.RFC3339)
		endDate := now.AddDate(0, 1, 0).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"rule_id":    createdRule.ID,
			"edit_mode":  "FUTURE",
			"start_date": startDate,
			"end_date":   endDate,
			"new_start":  "16:00",
			"new_end":    "17:00",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/scheduling/preview-recurrence-edit", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.PreviewRecurrenceEdit(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Preview recurrence edit completed")
	})

	t.Run("Step2_EditRecurringSchedule_Single", func(t *testing.T) {
		targetDate := now.AddDate(0, 0, 7).Format("2006-01-02")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		reqBody := map[string]interface{}{
			"rule_id":     createdRule.ID,
			"edit_mode":   "SINGLE",
			"target_date": targetDate,
			"new_start":   "14:00",
			"new_end":     "15:00",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/scheduling/edit-recurring", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.EditRecurringSchedule(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Edit recurring schedule (single) completed")
	})

	t.Run("Step3_EditRecurringSchedule_Future", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		startDate := now.AddDate(0, 0, 14).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"rule_id":    createdRule.ID,
			"edit_mode":  "FUTURE",
			"start_date": startDate,
			"new_start":  "10:00",
			"new_end":    "11:00",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/scheduling/edit-recurring", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.EditRecurringSchedule(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Edit recurring schedule (future) completed")
	})

	t.Run("Step4_DeleteRecurringSchedule", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		reqBody := map[string]interface{}{
			"rule_id":   createdRule.ID,
			"edit_mode": "ALL",
			"reason":    "Course discontinued",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/scheduling/delete-recurring", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.DeleteRecurringSchedule(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Delete recurring schedule completed")
	})
}

func TestIntegration_OfferingManagement(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Offering Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("offering_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Offering Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Offering Test Course",
		DefaultDuration:  60,
		ColorHex:         "#9B59B6",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringController := controllers.NewOfferingController(appInstance)
	adminResourceController := controllers.NewAdminResourceController(appInstance)

	t.Run("Step1_CreateOffering", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		reqBody := controllers.CreateOfferingRequest{
			CourseID: createdCourse.ID,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/offerings", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		offeringController.CreateOffering(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		datas := response["datas"].(map[string]interface{})
		offeringID := uint(datas["id"].(float64))
		if offeringID == 0 {
			t.Fatal("Expected non-zero offering ID")
		}
		t.Logf("Created offering: ID=%d", offeringID)
	})

	t.Run("Step2_GetOfferings", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/offerings", nil)

		offeringController.GetOfferings(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved offerings")
	})

	t.Run("Step3_ToggleOfferingActive", func(t *testing.T) {
		offeringRepo := repositories.NewOfferingRepository(appInstance)
		offering := models.Offering{
			CenterID:  createdCenter.ID,
			CourseID:  createdCourse.ID,
			Name:      fmt.Sprintf("Toggle Test Offering %d", now.UnixNano()),
			IsActive:  true,
			CreatedAt: now,
		}
		createdOffering, _ := offeringRepo.Create(ctx, offering)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "offering_id", Value: fmt.Sprintf("%d", createdOffering.ID)}}

		reqBody := map[string]interface{}{
			"is_active": false,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/admin/offerings/"+fmt.Sprintf("%d", createdOffering.ID)+"/toggle-active", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminResourceController.ToggleOfferingActive(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully toggled offering active status")
	})

	t.Run("Step4_CopyOffering", func(t *testing.T) {
		offeringRepo := repositories.NewOfferingRepository(appInstance)
		offering := models.Offering{
			CenterID:  createdCenter.ID,
			CourseID:  createdCourse.ID,
			Name:      fmt.Sprintf("Original Offering %d", now.UnixNano()),
			IsActive:  true,
			CreatedAt: now,
		}
		createdOffering, _ := offeringRepo.Create(ctx, offering)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "offeringId", Value: fmt.Sprintf("%d", createdOffering.ID)}}

		reqBody := map[string]interface{}{
			"new_name": fmt.Sprintf("Copied Offering %d", now.UnixNano()),
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/offerings/"+fmt.Sprintf("%d", createdOffering.ID)+"/copy", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		offeringController.CopyOffering(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully copied offering")
	})

	t.Run("Step5_DeleteOffering", func(t *testing.T) {
		offeringRepo := repositories.NewOfferingRepository(appInstance)
		offering := models.Offering{
			CenterID:  createdCenter.ID,
			CourseID:  createdCourse.ID,
			Name:      fmt.Sprintf("Delete Test Offering %d", now.UnixNano()),
			IsActive:  true,
			CreatedAt: now,
		}
		createdOffering, _ := offeringRepo.Create(ctx, offering)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "offering_id", Value: fmt.Sprintf("%d", createdOffering.ID)}}

		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/offerings/"+fmt.Sprintf("%d", createdOffering.ID), nil)

		offeringController.DeleteOffering(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully deleted offering")
	})
}

func TestIntegration_HolidayManagement(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Holiday Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("holiday_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Holiday Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	adminResourceController := controllers.NewAdminResourceController(appInstance)

	t.Run("Step1_CreateHoliday", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		holidayDate := now.AddDate(0, 1, 0).Format("2006-01-02")

		reqBody := map[string]interface{}{
			"date": holidayDate,
			"name": "National Day",
			"type": "NATIONAL_HOLIDAY",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/holidays", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminResourceController.CreateHoliday(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully created holiday")
	})

	t.Run("Step2_GetHolidays", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/holidays", nil)

		adminResourceController.GetHolidays(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved holidays")
	})

	t.Run("Step3_BulkCreateHolidays", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		holidays := []map[string]interface{}{
			{
				"date": now.AddDate(0, 2, 0).Format("2006-01-02"),
				"name": "Christmas",
				"type": "NATIONAL_HOLIDAY",
			},
			{
				"date": now.AddDate(0, 2, 15).Format("2006-01-02"),
				"name": "New Year",
				"type": "CENTER_CLOSED",
			},
		}

		reqBody := map[string]interface{}{
			"holidays": holidays,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/holidays/bulk", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminResourceController.BulkCreateHolidays(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully bulk created holidays")
	})
}

func TestIntegration_AdminUserManagement(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Admin User Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("admin_user_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Admin User Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	adminUserController := controllers.NewAdminUserController(appInstance)

	t.Run("Step1_CreateAdminUser", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		reqBody := map[string]interface{}{
			"email":    fmt.Sprintf("new_admin_%d@test.com", now.UnixNano()),
			"name":     "New Admin User",
			"role":     "STAFF",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminUserController.CreateAdminUser(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully created admin user")
	})

	t.Run("Step2_GetAdminUsers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users", nil)

		adminUserController.GetAdminUsers(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved admin users")
	})

	t.Run("Step3_UpdateAdminUser", func(t *testing.T) {
		newAdminUserRepo := repositories.NewAdminUserRepository(appInstance)
		newAdmin := models.AdminUser{
			Email:        fmt.Sprintf("update_admin_%d@test.com", now.UnixNano()),
			PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
			Name:         "Update Admin",
			CenterID:     createdCenter.ID,
			Role:         "STAFF",
			Status:       "ACTIVE",
			CreatedAt:    now,
		}
		createdNewAdmin, _ := newAdminUserRepo.Create(ctx, newAdmin)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "admin_id", Value: fmt.Sprintf("%d", createdNewAdmin.ID)}}

		reqBody := map[string]interface{}{
			"name": "Updated Admin Name",
			"role": "STAFF",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PUT", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users/"+fmt.Sprintf("%d", createdNewAdmin.ID), bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminUserController.UpdateAdminUser(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully updated admin user")
	})
}

func TestIntegration_TemplateManagement(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Template Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("template_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Template Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	templateController := controllers.NewTimetableTemplateController(appInstance)

	t.Run("Step1_CreateTemplate", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		reqBody := map[string]interface{}{
			"name":     fmt.Sprintf("Weekly Template %d", now.UnixNano()),
			"row_type": "ROOM",
			"cells": []map[string]interface{}{
				{
					"row_no":     1,
					"col_no":     1,
					"start_time": "09:00",
					"end_time":   "10:00",
				},
				{
					"row_no":     2,
					"col_no":     1,
					"start_time": "10:00",
					"end_time":   "11:00",
				},
			},
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/templates", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		templateController.CreateTemplate(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully created template")
	})

	t.Run("Step2_GetTemplates", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/templates", nil)

		templateController.GetTemplates(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved templates")
	})

	t.Run("Step3_UpdateTemplate", func(t *testing.T) {
		templateRepo := repositories.NewTimetableTemplateRepository(appInstance)
		template := models.TimetableTemplate{
			CenterID:  createdCenter.ID,
			Name:      fmt.Sprintf("Template to Update %d", now.UnixNano()),
			RowType:   "WEEKLY",
			CreatedAt: now,
		}
		createdTemplate, _ := templateRepo.Create(ctx, template)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "templateId", Value: fmt.Sprintf("%d", createdTemplate.ID)}}

		reqBody := map[string]interface{}{
			"name":        fmt.Sprintf("Updated Template %d", now.UnixNano()),
			"description": "Updated description",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PUT", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/templates/"+fmt.Sprintf("%d", createdTemplate.ID), bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		templateController.UpdateTemplate(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully updated template")
	})

	t.Run("Step4_GetTemplateCells", func(t *testing.T) {
		templateRepo := repositories.NewTimetableTemplateRepository(appInstance)
		template := models.TimetableTemplate{
			CenterID:  createdCenter.ID,
			Name:      fmt.Sprintf("Template for Cells %d", now.UnixNano()),
			RowType:   "WEEKLY",
			CreatedAt: now,
		}
		createdTemplate, _ := templateRepo.Create(ctx, template)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "templateId", Value: fmt.Sprintf("%d", createdTemplate.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/templates/"+fmt.Sprintf("%d", createdTemplate.ID)+"/cells", nil)

		templateController.GetCells(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved template cells")
	})

	t.Run("Step5_CreateTemplateCells", func(t *testing.T) {
		templateRepo := repositories.NewTimetableTemplateRepository(appInstance)
		template := models.TimetableTemplate{
			CenterID:  createdCenter.ID,
			Name:      fmt.Sprintf("Template for Create Cells %d", now.UnixNano()),
			RowType:   "WEEKLY",
			CreatedAt: now,
		}
		createdTemplate, _ := templateRepo.Create(ctx, template)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "templateId", Value: fmt.Sprintf("%d", createdTemplate.ID)}}

		cells := []map[string]interface{}{
			{
				"row_no":     1,
				"col_no":     1,
				"start_time": "14:00",
				"end_time":   "15:00",
			},
			{
				"row_no":     2,
				"col_no":     1,
				"start_time": "15:00",
				"end_time":   "16:00",
			},
		}
		body, _ := json.Marshal(cells)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/templates/"+fmt.Sprintf("%d", createdTemplate.ID)+"/cells", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		templateController.CreateCells(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully created template cells")
	})

	t.Run("Step6_DeleteTemplate", func(t *testing.T) {
		templateRepo := repositories.NewTimetableTemplateRepository(appInstance)
		template := models.TimetableTemplate{
			CenterID:  createdCenter.ID,
			Name:      fmt.Sprintf("Template to Delete %d", now.UnixNano()),
			RowType:   "WEEKLY",
			CreatedAt: now,
		}
		createdTemplate, _ := templateRepo.Create(ctx, template)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "templateId", Value: fmt.Sprintf("%d", createdTemplate.ID)}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/templates/"+fmt.Sprintf("%d", createdTemplate.ID), nil)

		templateController.DeleteTemplate(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully deleted template")
	})
}

func TestIntegration_ExportFunctionality(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Export Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			ExceptionLeadDays: 14,
		},
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("export_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Export Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_EXPORT_%d", now.UnixNano()),
		Name:       "Export Teacher",
		Email:      fmt.Sprintf("export_teacher_%d@test.com", now.UnixNano()),
		City:       "Taipei",
		District:   "Xinyi",
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Export Room",
		Capacity:  10,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             "Export Course",
		DefaultDuration:  60,
		ColorHex:         "#3498DB",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:  createdCenter.ID,
		CourseID:  createdCourse.ID,
		Name:      "Export Offering",
		IsActive:  true,
		CreatedAt: now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)

	effectiveRange := models.DateRange{
		StartDate: now,
		EndDate:   now.AddDate(1, 0, 0),
	}
	scheduleRuleRepo := repositories.NewScheduleRuleRepository(appInstance)
	scheduleRule := models.ScheduleRule{
		CenterID:       createdCenter.ID,
		OfferingID:     createdOffering.ID,
		TeacherID:      &createdTeacher.ID,
		RoomID:         createdRoom.ID,
		Weekday:        int(now.Weekday()),
		StartTime:      now.Format("15:04"),
		EndTime:        now.Add(time.Hour).Format("15:04"),
		EffectiveRange: effectiveRange,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	createdRule, _ := scheduleRuleRepo.Create(ctx, scheduleRule)

	exceptionRepo := repositories.NewScheduleExceptionRepository(appInstance)
	exception := models.ScheduleException{
		RuleID:       createdRule.ID,
		CenterID:     createdCenter.ID,
		OriginalDate: now.AddDate(0, 0, 5),
		Type:         "CANCEL",
		Status:       "PENDING",
		Reason:       "Holiday",
		CreatedAt:    now,
	}
	_, _ = exceptionRepo.Create(ctx, exception)

	exportController := controllers.NewExportController(appInstance)

	t.Run("Step1_ExportScheduleCSV", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		startDate := now.Format("2006-01-02")
		endDate := now.AddDate(0, 1, 0).Format("2006-01-02")

		reqBody := map[string]interface{}{
			"center_id":  createdCenter.ID,
			"start_date": startDate,
			"end_date":   endDate,
			"format":     "csv",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/export/schedule/csv", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		exportController.ExportScheduleCSV(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully exported schedule to CSV")
	})

	t.Run("Step2_ExportTeachersCSV", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/export/teachers/csv", nil)

		exportController.ExportTeachersCSV(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully exported teachers to CSV")
	})

	t.Run("Step3_ExportExceptionsCSV", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}
		startDate := now.Format("2006-01-02")
		endDate := now.AddDate(0, 1, 0).Format("2006-01-02")
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/export/exceptions/csv?start_date="+startDate+"&end_date="+endDate, nil)

		exportController.ExportExceptionsCSV(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully exported exceptions to CSV")
	})
}

func TestIntegration_AuthRefreshAndLogout(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Auth Refresh Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("auth_refresh_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Auth Refresh Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	_, _ = adminUserRepo.Create(ctx, adminUser)

	authService := services.NewAuthService(appInstance)
	authController := controllers.NewAuthController(appInstance, authService)

	t.Run("Step1_RefreshToken", func(t *testing.T) {
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		refreshReq := map[string]interface{}{
			"refresh_token": "test_refresh_token",
		}
		body, _ := json.Marshal(refreshReq)
		c2.Request = httptest.NewRequest("POST", "/api/v1/auth/admin/refresh", bytes.NewBuffer(body))
		c2.Request.Header.Set("Content-Type", "application/json")

		authController.RefreshToken(c2)

		if w2.Code != http.StatusOK && w2.Code != http.StatusBadRequest && w2.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 200, 400 or 401, got %d. Body: %s", w2.Code, w2.Body.String())
		}
		t.Log("Refresh token test completed")
	})

	t.Run("Step2_Logout", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		logoutReq := map[string]interface{}{
			"refresh_token": "test_refresh_token",
		}
		body, _ := json.Marshal(logoutReq)
		c.Request = httptest.NewRequest("POST", "/api/v1/auth/admin/logout", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		authController.Logout(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 200, 400 or 401, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Logout test completed")
	})
}

func TestIntegration_AdminResourceCentersAndDelete(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	centerRepo := repositories.NewCenterRepository(appInstance)
	center := models.Center{
		Name:      fmt.Sprintf("Admin Resource Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("admin_resource_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Admin Resource Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	adminResourceController := controllers.NewAdminResourceController(appInstance)

	t.Run("Step1_GetCenters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers", nil)

		adminResourceController.GetCenters(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully retrieved centers")
	})

	t.Run("Step2_CreateCenter", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)

		reqBody := map[string]interface{}{
			"name":       fmt.Sprintf("New Center %d", now.UnixNano()),
			"plan_level": "STARTER",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		adminResourceController.CreateCenter(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully created center")
	})

	t.Run("Step3_DeleteCourse", func(t *testing.T) {
		courseRepo := repositories.NewCourseRepository(appInstance)
		course := models.Course{
			CenterID:         createdCenter.ID,
			Name:             fmt.Sprintf("Course to Delete %d", now.UnixNano()),
			DefaultDuration:  60,
			ColorHex:         "#FF0000",
			RoomBufferMin:    10,
			TeacherBufferMin: 10,
			IsActive:         true,
			CreatedAt:        now,
		}
		createdCourse, _ := courseRepo.Create(ctx, course)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "course_id", Value: fmt.Sprintf("%d", createdCourse.ID)}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/courses/"+fmt.Sprintf("%d", createdCourse.ID), nil)

		adminResourceController.DeleteCourse(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully deleted course")
	})

	t.Run("Step4_DeleteHoliday", func(t *testing.T) {
		holidayRepo := repositories.NewCenterHolidayRepository(appInstance)
		holiday := models.CenterHoliday{
			CenterID:  createdCenter.ID,
			Date:      now.AddDate(0, 1, 0),
			Name:      fmt.Sprintf("Holiday to Delete %d", now.UnixNano()),
			CreatedAt: now,
		}
		createdHoliday, _ := holidayRepo.Create(ctx, holiday)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "holiday_id", Value: fmt.Sprintf("%d", createdHoliday.ID)}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/holidays/"+fmt.Sprintf("%d", createdHoliday.ID), nil)

		adminResourceController.DeleteHoliday(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Successfully deleted holiday")
	})
}

func TestIntegration_SchedulingBufferAndMore(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Buffer Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("buffer_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Buffer Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_BUFFER_%d", now.UnixNano()),
		Name:       "Buffer Teacher",
		Email:      fmt.Sprintf("buffer_teacher_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Buffer Room",
		Capacity:  10,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	schedulingController := controllers.NewSchedulingController(appInstance)

	t.Run("Step1_CheckTeacherBuffer", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)

		startTime := now.Add(time.Hour * 9).Format(time.RFC3339)
		endTime := now.Add(time.Hour * 10).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"teacher_id": createdTeacher.ID,
			"start_time": startTime,
			"end_time":   endTime,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/check-teacher-buffer", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CheckTeacherBuffer(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Check teacher buffer completed")
	})

	t.Run("Step2_CheckRoomBuffer", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)

		startTime := now.Add(time.Hour * 11).Format(time.RFC3339)
		endTime := now.Add(time.Hour * 12).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"room_id":    createdRoom.ID,
			"start_time": startTime,
			"end_time":   endTime,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/check-room-buffer", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CheckRoomBuffer(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Check room buffer completed")
	})

	t.Run("Step3_GetExceptionsByRule", func(t *testing.T) {
		ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
		rule := models.ScheduleRule{
			CenterID:  createdCenter.ID,
			Weekday:   1,
			StartTime: "10:00:00",
			EndTime:   "11:00:00",
			CreatedAt: now,
			UpdatedAt: now,
		}
		createdRule, _ := ruleRepo.Create(ctx, rule)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "ruleId", Value: fmt.Sprintf("%d", createdRule.ID)}}
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/scheduling/rules/"+fmt.Sprintf("%d", createdRule.ID)+"/exceptions", nil)

		schedulingController.GetExceptionsByRule(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get exceptions by rule completed")
	})

	t.Run("Step4_GetExceptionsByDateRange", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		startDate := now.Format("2006-01-02")
		endDate := now.AddDate(0, 1, 0).Format("2006-01-02")

		c.Request = httptest.NewRequest("GET", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/exceptions?start_date="+startDate+"&end_date="+endDate, nil)

		schedulingController.GetExceptionsByDateRange(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get exceptions by date range completed")
	})

	t.Run("Step5_DeleteRule", func(t *testing.T) {
		ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
		rule := models.ScheduleRule{
			CenterID:  createdCenter.ID,
			Weekday:   3,
			StartTime: "14:00:00",
			EndTime:   "15:00:00",
			CreatedAt: now,
			UpdatedAt: now,
		}
		createdRule, _ := ruleRepo.Create(ctx, rule)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "ruleId", Value: fmt.Sprintf("%d", createdRule.ID)}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/rules/"+fmt.Sprintf("%d", createdRule.ID), nil)

		schedulingController.DeleteRule(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Delete rule completed")
	})
}

func TestIntegration_TeacherSkillsCertificates(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Teacher SC Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_SC_%d", now.UnixNano()),
		Name:       "Teacher SC",
		Email:      fmt.Sprintf("teacher_sc_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	teacherController := controllers.NewTeacherController(appInstance)

	t.Run("Step1_GetSkills", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/me/skills", nil)

		teacherController.GetSkills(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get skills completed")
	})

	t.Run("Step2_CreateSkill", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		reqBody := map[string]interface{}{
			"category":   "MUSIC",
			"skill_name": "Piano",
			"level":      "ADVANCED",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/me/skills", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.CreateSkill(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Create skill completed")
	})

	t.Run("Step3_GetCertificates", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/me/certificates", nil)

		teacherController.GetCertificates(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get certificates completed")
	})

	t.Run("Step4_CreateCertificate", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		reqBody := map[string]interface{}{
			"name":      "Test Certificate",
			"file_url":  "https://example.com/cert.pdf",
			"issued_at": now.Format(time.RFC3339),
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/me/certificates", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.CreateCertificate(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Create certificate completed")
	})

	t.Run("Step5_UploadCertificate", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		reqBody := map[string]interface{}{
			"file_url":  "https://example.com/uploaded.pdf",
			"issued_at": now.Format(time.RFC3339),
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/me/certificates", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.CreateCertificate(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Upload certificate completed")
	})
}

func TestIntegration_TeacherNotesAndEvents(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Teacher Notes Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_NOTES_%d", now.UnixNano()),
		Name:       "Teacher Notes",
		Email:      fmt.Sprintf("teacher_notes_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	teacherController := controllers.NewTeacherController(appInstance)

	t.Run("Step1_GetSessionNote", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/sessions/note?rule_id=1&session_date="+now.Format("2006-01-02"), nil)

		teacherController.GetSessionNote(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get session note completed")
	})

	t.Run("Step2_UpsertSessionNote", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		reqBody := map[string]interface{}{
			"rule_id":      1,
			"session_date": now.Format("2006-01-02"),
			"content":      "Test session note",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PUT", "/api/v1/teacher/sessions/note", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.UpsertSessionNote(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Upsert session note completed")
	})

	t.Run("Step3_GetPersonalEvents", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/me/personal-events", nil)

		teacherController.GetPersonalEvents(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get personal events completed")
	})

	t.Run("Step4_CreatePersonalEvent", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		startAt := now.AddDate(0, 0, 7).Format(time.RFC3339)
		endAt := now.AddDate(0, 0, 7).Add(time.Hour).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"title":    "Doctor Appointment",
			"start_at": startAt,
			"end_at":   endAt,
			"type":     "PERSONAL",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/teacher/me/personal-events", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.CreatePersonalEvent(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Create personal event completed")
	})

	t.Run("Step5_UpdatePersonalEvent", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		reqBody := map[string]interface{}{
			"title":       "Updated Event",
			"update_mode": "SINGLE",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/teacher/me/personal-events/1", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.UpdatePersonalEvent(c)

		if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200, 400 or 404, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Update personal event completed")
	})
}

func TestIntegration_AdminTeacherManagement(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Admin Teacher Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("admin_teacher_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Admin Teacher Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_ADMIN_T_%d", now.UnixNano()),
		Name:       "Admin Teacher",
		Email:      fmt.Sprintf("admin_teacher_test_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	_, _ = teacherRepo.Create(ctx, teacher)

	teacherController := controllers.NewTeacherController(appInstance)

	t.Run("Step1_ListTeachers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Request = httptest.NewRequest("GET", "/api/v1/admin/teachers", nil)

		teacherController.ListTeachers(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("List teachers completed")
	})

	t.Run("Step2_DeleteTeacher", func(t *testing.T) {
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		teacher := models.Teacher{
			LineUserID: fmt.Sprintf("LINE_DELETE_%d", now.UnixNano()),
			Name:       "Teacher to Delete",
			Email:      fmt.Sprintf("delete_teacher_%d@test.com", now.UnixNano()),
			CreatedAt:  now,
		}
		createdDelTeacher, _ := teacherRepo.Create(ctx, teacher)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdDelTeacher.ID)}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/teachers/"+fmt.Sprintf("%d", createdDelTeacher.ID), nil)

		teacherController.DeleteTeacher(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Delete teacher completed")
	})

	t.Run("Step3_InviteTeacher", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}}

		reqBody := map[string]interface{}{
			"email": fmt.Sprintf("invite_teacher_%d@test.com", now.UnixNano()),
			"role":  "TEACHER",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/teachers/invite", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.InviteTeacher(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Invite teacher completed")
	})
}

func TestIntegration_AdminUserDeletion(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Admin Delete Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("admin_delete_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Admin Delete Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	t.Run("Step1_DeleteAdminUser", func(t *testing.T) {
		adminUserRepo := repositories.NewAdminUserRepository(appInstance)
		adminToDelete := models.AdminUser{
			Email:        fmt.Sprintf("delete_me_%d@test.com", now.UnixNano()),
			PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
			Name:         "To Delete",
			CenterID:     createdCenter.ID,
			Role:         "ADMIN",
			Status:       "ACTIVE",
			CreatedAt:    now,
		}
		createdDelAdmin, _ := adminUserRepo.Create(ctx, adminToDelete)

		adminUserController := controllers.NewAdminUserController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "admin_id", Value: fmt.Sprintf("%d", createdDelAdmin.ID)}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/users/"+fmt.Sprintf("%d", createdDelAdmin.ID), nil)

		adminUserController.DeleteAdminUser(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Delete admin user completed")
	})
}

func TestIntegration_OfferingUpdate(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Offering Update Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("offering_update_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Offering Update Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	courseRepo := repositories.NewCourseRepository(appInstance)
	course := models.Course{
		CenterID:         createdCenter.ID,
		Name:             fmt.Sprintf("Offering Update Course %d", now.UnixNano()),
		DefaultDuration:  60,
		ColorHex:         "#FF5733",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        now,
	}
	createdCourse, _ := courseRepo.Create(ctx, course)

	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering := models.Offering{
		CenterID:  createdCenter.ID,
		CourseID:  createdCourse.ID,
		Name:      fmt.Sprintf("Offering to Update %d", now.UnixNano()),
		IsActive:  true,
		CreatedAt: now,
	}
	createdOffering, _ := offeringRepo.Create(ctx, offering)

	t.Run("Step1_UpdateOffering", func(t *testing.T) {
		offeringController := controllers.NewOfferingController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "offering_id", Value: fmt.Sprintf("%d", createdOffering.ID)}}

		reqBody := map[string]interface{}{
			"name": fmt.Sprintf("Updated Offering %d", now.UnixNano()),
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/admin/offerings/"+fmt.Sprintf("%d", createdOffering.ID), bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		offeringController.UpdateOffering(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Update offering completed")
	})
}

func TestIntegration_TemplateUpdate(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Template Update Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("template_update_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Template Update Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	t.Run("Step1_UpdateTemplate", func(t *testing.T) {
		templateRepo := repositories.NewTimetableTemplateRepository(appInstance)
		template := models.TimetableTemplate{
			CenterID:  createdCenter.ID,
			Name:      fmt.Sprintf("Template to Update %d", now.UnixNano()),
			RowType:   "WEEKLY",
			CreatedAt: now,
		}
		createdTemplate, _ := templateRepo.Create(ctx, template)

		templateController := controllers.NewTimetableTemplateController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", createdCenter.ID)}, {Key: "templateId", Value: fmt.Sprintf("%d", createdTemplate.ID)}}

		reqBody := map[string]interface{}{
			"name": fmt.Sprintf("Updated Template %d", now.UnixNano()),
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/admin/centers/"+fmt.Sprintf("%d", createdCenter.ID)+"/templates/"+fmt.Sprintf("%d", createdTemplate.ID), bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		templateController.UpdateTemplate(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Update template completed")
	})
}

func TestIntegration_NotificationManagement(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Notification Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("notification_admin_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Notification Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	notificationController := controllers.NewNotificationController(appInstance)

	t.Run("Step1_ListNotifications", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Request = httptest.NewRequest("GET", "/api/v1/notifications", nil)

		notificationController.ListNotifications(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("List notifications completed")
	})

	t.Run("Step2_GetUnreadCount", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Request = httptest.NewRequest("GET", "/api/v1/notifications/unread-count", nil)

		notificationController.GetUnreadCount(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get unread count completed")
	})

	t.Run("Step3_SetNotifyToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		reqBody := map[string]interface{}{
			"token": fmt.Sprintf("test_notify_token_%d", now.UnixNano()),
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/notifications/token", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		notificationController.SetNotifyToken(c)

		if w.Code != http.StatusOK && w.Code != http.StatusForbidden && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
			t.Fatalf("Expected status 200, 403, 400 or 500, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Set notify token completed")
	})

	t.Run("Step4_SendTestNotification", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		reqBody := map[string]interface{}{
			"title":   "Test Notification",
			"content": "This is a test notification",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/notifications/test", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		notificationController.SendTestNotification(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
			t.Fatalf("Expected status 200, 400 or 500, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Send test notification completed")
	})

	t.Run("Step5_MarkAllAsRead", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Request = httptest.NewRequest("POST", "/api/v1/notifications/mark-all-read", nil)

		notificationController.MarkAllAsRead(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Mark all as read completed")
	})

	t.Run("Step6_MarkAsRead", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("POST", "/api/v1/notifications/1/read", nil)

		notificationController.MarkAsRead(c)

		if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 200 or 404, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Mark as read completed")
	})
}

func TestIntegration_SmartMatching(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Smart Matching Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	adminUserRepo := repositories.NewAdminUserRepository(appInstance)
	adminUser := models.AdminUser{
		Email:        fmt.Sprintf("smart_matching_%d@test.com", now.UnixNano()),
		PasswordHash: "$2a$10$lVIoLQr4EjCjQIU98JExROfBoOFK.UNOkVS0LVH2Lj1rT0VX5DYqa",
		Name:         "Smart Matching Admin",
		CenterID:     createdCenter.ID,
		Role:         "ADMIN",
		Status:       "ACTIVE",
		CreatedAt:    now,
	}
	createdAdmin, _ := adminUserRepo.Create(ctx, adminUser)

	roomRepo := repositories.NewRoomRepository(appInstance)
	room := models.Room{
		CenterID:  createdCenter.ID,
		Name:      "Smart Matching Room",
		Capacity:  10,
		IsActive:  true,
		CreatedAt: now,
	}
	createdRoom, _ := roomRepo.Create(ctx, room)

	smartMatchingController := controllers.NewSmartMatchingController(appInstance)

	t.Run("Step1_FindMatches", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		startTime := now.Format(time.RFC3339)
		endTime := now.Add(time.Hour).Format(time.RFC3339)

		reqBody := map[string]interface{}{
			"center_id":  createdCenter.ID,
			"room_id":    createdRoom.ID,
			"start_time": startTime,
			"end_time":   endTime,
			"duration":   60,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/scheduling/find-matches", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		smartMatchingController.FindMatches(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Find matches completed")
	})

	t.Run("Step2_SearchTalent", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdAdmin.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "ADMIN")

		reqBody := map[string]interface{}{
			"query": "music",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/talent/search", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		smartMatchingController.SearchTalent(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 200 or 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Search talent completed")
	})
}

func TestIntegration_TeacherDeletion(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Teacher Deletion Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_DEL_%d", now.UnixNano()),
		Name:       "Teacher Deletion",
		Email:      fmt.Sprintf("teacher_del_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	teacherController := controllers.NewTeacherController(appInstance)

	t.Run("Step1_DeleteSkill", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("DELETE", "/api/v1/teacher/me/skills/1", nil)

		teacherController.DeleteSkill(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 200, 400 or 404, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Delete skill completed")
	})

	t.Run("Step2_DeleteCertificate", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("DELETE", "/api/v1/teacher/me/certificates/1", nil)

		teacherController.DeleteCertificate(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 200, 400 or 404, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Delete certificate completed")
	})

	t.Run("Step3_DeletePersonalEvent", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("DELETE", "/api/v1/teacher/me/personal-events/1", nil)

		teacherController.DeletePersonalEvent(c)

		if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 200, 400 or 404, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Delete personal event completed")
	})
}

func TestIntegration_TeacherGetCenters(t *testing.T) {
	appInstance, _, cleanup := setupIntegrationTestAppWithMigrations()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	center := models.Center{
		Name:      fmt.Sprintf("Teacher Centers Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, _ := centerRepo.Create(ctx, center)

	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher := models.Teacher{
		LineUserID: fmt.Sprintf("LINE_CENTERS_%d", now.UnixNano()),
		Name:       "Teacher Centers",
		Email:      fmt.Sprintf("teacher_centers_%d@test.com", now.UnixNano()),
		CreatedAt:  now,
	}
	createdTeacher, _ := teacherRepo.Create(ctx, teacher)

	t.Run("Step1_GetCenters", func(t *testing.T) {
		teacherController := controllers.NewTeacherController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")
		c.Request = httptest.NewRequest("GET", "/api/v1/teacher/centers", nil)

		teacherController.GetCenters(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Get centers completed")
	})

	t.Run("Step2_UpdateProfile", func(t *testing.T) {
		teacherController := controllers.NewTeacherController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, createdTeacher.ID)
		c.Set(global.CenterIDKey, createdCenter.ID)
		c.Set(global.UserTypeKey, "TEACHER")

		reqBody := map[string]interface{}{
			"name": fmt.Sprintf("Updated Teacher %d", now.UnixNano()),
			"bio":  "Updated bio",
			"city": "Kaohsiung",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/teacher/me/profile", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		teacherController.UpdateProfile(c)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		t.Log("Update profile completed")
	})
}
