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

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"timeLedger/app"
	"timeLedger/app/controllers"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global"
	"timeLedger/global/errInfos"
	mockRedis "timeLedger/testing/redis"

	"github.com/gin-gonic/gin"
)

// setupBufferOverrideTestApp 設定整合測試環境
func setupBufferOverrideTestApp() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	// 初始化自訂驗證器（解決 time_format 等問題）
	requests.InitValidators()

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	// AutoMigrate required tables
	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.Teacher{},
		&models.Room{},
		&models.Course{},
		&models.Offering{},
		&models.ScheduleRule{},
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
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	cleanup := func() {
		mr.Close()
	}

	return appInstance, mysqlDB, cleanup
}

// createTestData 建立測試所需的資料
func createTestData(appInstance *app.App) (center models.Center, teacher models.Teacher, room models.Room, course models.Course, offering models.Offering, err error) {
	ctx := context.Background()

	// 建立 Center
	center = models.Center{
		Name:      fmt.Sprintf("Buffer Override Test Center %d", time.Now().UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: time.Now(),
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	center, err = centerRepo.Create(ctx, center)
	if err != nil {
		return
	}

	// 建立 Teacher
	teacher = models.Teacher{
		Name:       fmt.Sprintf("Test Teacher %d", time.Now().UnixNano()),
		Email:      fmt.Sprintf("teacher%d@test.com", time.Now().UnixNano()),
		LineUserID: fmt.Sprintf("LINE_%d", time.Now().UnixNano()),
		CreatedAt:  time.Now(),
	}
	teacherRepo := repositories.NewTeacherRepository(appInstance)
	teacher, err = teacherRepo.Create(ctx, teacher)
	if err != nil {
		return
	}

	// 建立 Room
	room = models.Room{
		CenterID:  center.ID,
		Name:      fmt.Sprintf("Test Room %d", time.Now().UnixNano()),
		Capacity:  20,
		CreatedAt: time.Now(),
	}
	roomRepo := repositories.NewRoomRepository(appInstance)
	room, err = roomRepo.Create(ctx, room)
	if err != nil {
		return
	}

	// 建立 Course
	course = models.Course{
		CenterID:         center.ID,
		Name:             fmt.Sprintf("Test Course %d", time.Now().UnixNano()),
		DefaultDuration:  60,
		ColorHex:         "#3498db",
		TeacherBufferMin: 15,
		RoomBufferMin:    10,
		IsActive:         true,
		CreatedAt:        time.Now(),
	}
	courseRepo := repositories.NewCourseRepository(appInstance)
	course, err = courseRepo.Create(ctx, course)
	if err != nil {
		return
	}

	// 建立 Offering
	offering = models.Offering{
		CenterID:  center.ID,
		CourseID:  course.ID,
		Name:      fmt.Sprintf("Test Offering %d", time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	offeringRepo := repositories.NewOfferingRepository(appInstance)
	offering, err = offeringRepo.Create(ctx, offering)
	if err != nil {
		return
	}

	return
}

// TestIntegration_ApplyTemplate_BufferOverride
// 測試模板套用的 Buffer Override 功能
func TestIntegration_ApplyTemplate_BufferOverride(t *testing.T) {
	appInstance, _, cleanup := setupBufferOverrideTestApp()
	defer cleanup()

	// 建立測試資料
	center, teacher, room, _, offering, err := createTestData(appInstance)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	ctx := context.Background()

	// 建立現有的排課規則（08:00-09:00）造成 Buffer 衝突
	existingRule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  &teacher.ID,
		RoomID:     room.ID,
		Weekday:    1, // 週一
		StartTime:  "08:00",
		EndTime:    "09:00",
		EffectiveRange: models.DateRange{
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 6, 0),
		},
		CreatedAt: time.Now(),
	}
	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	_, err = ruleRepo.Create(ctx, existingRule)
	if err != nil {
		t.Fatalf("Failed to create existing rule: %v", err)
	}

	// 建立模板
	template := models.TimetableTemplate{
		CenterID:  center.ID,
		Name:      fmt.Sprintf("Test Template %d", time.Now().UnixNano()),
		RowType:   "TIME",
		CreatedAt: time.Now(),
	}
	templateRepo := repositories.NewTimetableTemplateRepository(appInstance)
	createdTemplate, err := templateRepo.Create(ctx, template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// 建立模板格子（09:00-10:00，與現有規則沒有緩衝時間）
	cell := models.TimetableCell{
		TemplateID: createdTemplate.ID,
		RowNo:      1,
		ColNo:      1,
		StartTime:  "09:00",
		EndTime:    "10:00",
		RoomID:     &room.ID,
		TeacherID:  &teacher.ID,
	}
	cellRepo := repositories.NewTimetableCellRepository(appInstance)
	_, err = cellRepo.Create(ctx, cell)
	if err != nil {
		t.Fatalf("Failed to create cell: %v", err)
	}

	t.Run("ApplyTemplate_BufferConflict_WithoutOverride", func(t *testing.T) {
		templateController := controllers.NewTimetableTemplateController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, uint(1))
		c.Set(global.CenterIDKey, center.ID)
		c.Params = gin.Params{
			{Key: "id", Value: fmt.Sprintf("%d", center.ID)},
			{Key: "templateId", Value: fmt.Sprintf("%d", createdTemplate.ID)},
		}

		// 嘗試套用模板（09:00-10:00，與 08:00-09:00 沒有緩衝時間）
		reqBody := map[string]interface{}{
			"offering_id":     offering.ID,
			"start_date":      time.Now().Format("2006-01-02"),
			"end_date":        time.Now().AddDate(0, 3, 0).Format("2006-01-02"),
			"weekdays":        []int{1},
			"override_buffer": false,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", center.ID)+"/templates/"+fmt.Sprintf("%d", createdTemplate.ID)+"/apply", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		templateController.ApplyTemplate(c)

		// 應該返回 400（有衝突）
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
			return
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		t.Logf("Buffer conflict detected: %s", response["message"])
	})

	t.Run("ApplyTemplate_OverrideBuffer_Success", func(t *testing.T) {
		templateController := controllers.NewTimetableTemplateController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, uint(1))
		c.Set(global.CenterIDKey, center.ID)
		c.Params = gin.Params{
			{Key: "id", Value: fmt.Sprintf("%d", center.ID)},
			{Key: "templateId", Value: fmt.Sprintf("%d", createdTemplate.ID)},
		}

		// 使用 override_buffer = true
		reqBody := map[string]interface{}{
			"offering_id":     offering.ID,
			"start_date":      time.Now().Format("2006-01-02"),
			"end_date":        time.Now().AddDate(0, 3, 0).Format("2006-01-02"),
			"weekdays":        []int{2}, // 週二，避免與現有規則衝突
			"override_buffer": true,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/centers/"+fmt.Sprintf("%d", center.ID)+"/templates/"+fmt.Sprintf("%d", createdTemplate.ID)+"/apply", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		templateController.ApplyTemplate(c)

		// 應該成功
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
			return
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		t.Logf("Template applied with override: %s", response["message"])
	})
}

// TestIntegration_CreateRule_BufferOverride
// 測試建立規則的 Buffer Override 功能
func TestIntegration_CreateRule_BufferOverride(t *testing.T) {
	appInstance, _, cleanup := setupBufferOverrideTestApp()
	defer cleanup()

	// 建立測試資料
	center, teacher, room, _, offering, err := createTestData(appInstance)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	ctx := context.Background()

	// 建立現有的排課規則（08:00-09:00）
	existingRule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  &teacher.ID,
		RoomID:     room.ID,
		Weekday:    1, // 週一
		StartTime:  "08:00",
		EndTime:    "09:00",
		EffectiveRange: models.DateRange{
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 6, 0),
		},
		CreatedAt: time.Now(),
	}
	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	_, err = ruleRepo.Create(ctx, existingRule)
	if err != nil {
		t.Fatalf("Failed to create existing rule: %v", err)
	}

	t.Run("CreateRule_BufferConflict_WithoutOverride", func(t *testing.T) {
		schedulingController := controllers.NewSchedulingController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, uint(1))
		c.Set(global.CenterIDKey, center.ID)

		// 嘗試在 09:00-10:00 建立規則（沒有緩衝時間）
		reqBody := requests.CreateRuleRequest{
			Name:           fmt.Sprintf("Test Rule %d", time.Now().UnixNano()),
			OfferingID:     offering.ID,
			TeacherID:      &teacher.ID,
			RoomID:         room.ID,
			StartTime:      "09:00",
			EndTime:        "10:00",
			Duration:       60,
			Weekdays:       []int{1},
			StartDate:      time.Now().Format("2006-01-02"),
			OverrideBuffer: false,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/rules", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CreateRule(c)

		// 應該返回 500（Insufficient buffer time 錯誤）
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d. Body: %s", w.Code, w.Body.String())
			return
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		t.Logf("Buffer conflict detected: %s", response["message"])
	})

	t.Run("CreateRule_OverrideBuffer_Success", func(t *testing.T) {
		schedulingController := controllers.NewSchedulingController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, uint(1))
		c.Set(global.CenterIDKey, center.ID)

		// 使用 override_buffer = true
		reqBody := requests.CreateRuleRequest{
			Name:           fmt.Sprintf("Override Test Rule %d", time.Now().UnixNano()),
			OfferingID:     offering.ID,
			TeacherID:      &teacher.ID,
			RoomID:         room.ID,
			StartTime:      "09:00",
			EndTime:        "10:00",
			Duration:       60,
			Weekdays:       []int{2}, // 週二，避免與現有規則衝突
			StartDate:      time.Now().Format("2006-01-02"),
			OverrideBuffer: true,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/rules", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CreateRule(c)

		// 應該成功
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
			return
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		t.Logf("Rule created with override: %s", response["message"])
	})
}

// TestIntegration_OverlapConflict_CannotOverride
// 測試重疊衝突不可被覆蓋
func TestIntegration_OverlapConflict_CannotOverride(t *testing.T) {
	appInstance, _, cleanup := setupBufferOverrideTestApp()
	defer cleanup()

	// 建立測試資料
	center, teacher, room, _, offering, err := createTestData(appInstance)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	ctx := context.Background()

	// 建立現有的排課規則（08:00-10:00）
	existingRule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  &teacher.ID,
		RoomID:     room.ID,
		Weekday:    1, // 週一
		StartTime:  "08:00",
		EndTime:    "10:00",
		EffectiveRange: models.DateRange{
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 6, 0),
		},
		CreatedAt: time.Now(),
	}
	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	_, err = ruleRepo.Create(ctx, existingRule)
	if err != nil {
		t.Fatalf("Failed to create existing rule: %v", err)
	}

	t.Run("CreateRule_OverlapConflict_EvenWithOverride", func(t *testing.T) {
		schedulingController := controllers.NewSchedulingController(appInstance)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(global.UserIDKey, uint(1))
		c.Set(global.CenterIDKey, center.ID)

		// 嘗試在完全相同的時間建立規則（完全重疊）
		reqBody := requests.CreateRuleRequest{
			Name:           fmt.Sprintf("Overlap Rule %d", time.Now().UnixNano()),
			OfferingID:     offering.ID,
			TeacherID:      &teacher.ID,
			RoomID:         room.ID,
			StartTime:      "09:00", // 完全在現有規則的時間範圍內
			EndTime:        "10:00",
			Duration:       60,
			Weekdays:       []int{1},
			StartDate:      time.Now().Format("2006-01-02"),
			OverrideBuffer: true, // 即使允許覆蓋
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/api/v1/admin/rules", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		schedulingController.CreateRule(c)

		// 應該被拒絕（Overlap 不可覆蓋）
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
			return
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		t.Logf("Overlap conflict rejected: %s", response["message"])

		// 驗證返回的是 Overlap 衝突代碼 (40002)，不是 Buffer 衝突代碼 (40003)
		if response["code"] != float64(40002) {
			t.Errorf("Expected code 40002 for overlap conflict, got %v", response["code"])
		}
	})
}
