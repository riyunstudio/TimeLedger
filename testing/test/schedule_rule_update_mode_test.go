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

func setupTestAppWithRealDB() (*app.App, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
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

	_ = mr
	return appInstance, mysqlDB
}

// TestScheduleRuleUpdateMode_Single 測試 SINGLE 模式
func TestScheduleRuleUpdateMode_Single(t *testing.T) {
	appInstance, _ := setupTestAppWithRealDB()

	ctx := context.Background()
	now := time.Now()

	// 查詢現有資料
	var center models.Center
	if err := appInstance.MySQL.RDB.WithContext(ctx).Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("No center data available, skipping test: %v", err)
		return
	}

	var course models.Course
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&course).Error; err != nil {
		t.Skipf("No course data available, skipping test: %v", err)
		return
	}

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 建立獨立的 offering 確保測試隔離
	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	// 建立測試規則
	testStartDate := now.AddDate(0, 1, 0)
	testEndDate := now.AddDate(0, 3, 0)

	rule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  nil,
		RoomID:     room.ID,
		Name:       fmt.Sprintf("Test Piano %d", now.UnixNano()),
		Weekday:    1,
		StartTime:  "09:00",
		EndTime:    "10:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: testStartDate,
			EndDate:   testEndDate,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	t.Logf("=== Test: SINGLE Mode ===")
	t.Logf("原始規則: ID=%d, Name=%s, Weekday=%d, StartTime=%s, StartDate=%s, EndDate=%s",
		rule.ID, rule.Name, rule.Weekday, rule.StartTime,
		rule.EffectiveRange.StartDate.Format("2006-01-02"),
		rule.EffectiveRange.EndDate.Format("2006-01-02"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, center.ID)
		c.Set(global.UserIDKey, uint(1))
		c.Next()
	})

	schedulingCtrl := controllers.NewSchedulingController(appInstance)
	router.PUT("/admin/rules/:ruleId", schedulingCtrl.UpdateRule)

	// 測試 SINGLE 模式
	newStartDate := now.AddDate(0, 1, 7).Format("2006-01-02")
	newEndDate := now.AddDate(0, 4, 0).Format("2006-01-02")

	t.Logf("預期更新: StartDate=%s, EndDate=%s, Mode=SINGLE", newStartDate, newEndDate)

	body := map[string]interface{}{
		"name":        rule.Name,
		"offering_id": offering.ID,
		"room_id":     room.ID,
		"start_time":  "10:00",
		"end_time":    "11:00",
		"duration":    60,
		"weekdays":    []int{1},
		"start_date":  newStartDate,
		"end_date":    newEndDate,
		"update_mode": "SINGLE",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/rules/%d", rule.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response global.ApiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	rules := response.Datas.([]interface{})
	expectedCount := 1
	actualCount := len(rules)

	t.Logf("預期規則數量: %d", expectedCount)
	t.Logf("實際規則數量: %d", actualCount)

	if actualCount != expectedCount {
		t.Errorf("規則數量不符: 預期=%d, 實際=%d", expectedCount, actualCount)
	}

	// 驗證規則內容
	if actualCount > 0 {
		ruleData := rules[0].(map[string]interface{})
		effectiveRange := ruleData["effective_range"].(map[string]interface{})
		actualStartDate := effectiveRange["start_date"].(string)
		actualEndDate := effectiveRange["end_date"].(string)

		expectedStartDate := newStartDate + "T00:00:00Z"
		expectedEndDate := newEndDate + "T00:00:00Z"

		t.Logf("預期 StartDate: %s", expectedStartDate)
		t.Logf("實際 StartDate:   %s", actualStartDate)
		t.Logf("預期 EndDate:   %s", expectedEndDate)
		t.Logf("實際 EndDate:     %s", actualEndDate)

		if actualStartDate != expectedStartDate {
			t.Errorf("StartDate 不符: 預期=%s, 實際=%s", expectedStartDate, actualStartDate)
		}
		if actualEndDate != expectedEndDate {
			t.Errorf("EndDate 不符: 預期=%s, 實際=%s", expectedEndDate, actualEndDate)
		}
	}

	t.Logf("=== SINGLE Mode 測試結果: %s ===", mapResultToText(actualCount == expectedCount))
}

// TestScheduleRuleUpdateMode_Future 測試 FUTURE 模式
func TestScheduleRuleUpdateMode_Future(t *testing.T) {
	appInstance, _ := setupTestAppWithRealDB()

	ctx := context.Background()
	now := time.Now()

	var center models.Center
	if err := appInstance.MySQL.RDB.WithContext(ctx).Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("No center data available, skipping test: %v", err)
		return
	}

	var course models.Course
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&course).Error; err != nil {
		t.Skipf("No course data available, skipping test: %v", err)
		return
	}

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 建立獨立的 offering 確保測試隔離
	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	testStartDate := now.AddDate(0, 1, 0)
	testEndDate := now.AddDate(0, 3, 0)

	rule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  nil,
		RoomID:     room.ID,
		Name:       fmt.Sprintf("Test Yoga %d", now.UnixNano()),
		Weekday:    2,
		StartTime:  "14:00",
		EndTime:    "15:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: testStartDate,
			EndDate:   testEndDate,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	t.Logf("=== Test: FUTURE Mode ===")
	t.Logf("原始規則: ID=%d, Name=%s, Weekday=%d, StartDate=%s, EndDate=%s",
		rule.ID, rule.Name, rule.Weekday,
		rule.EffectiveRange.StartDate.Format("2006-01-02"),
		rule.EffectiveRange.EndDate.Format("2006-01-02"))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, center.ID)
		c.Set(global.UserIDKey, uint(1))
		c.Next()
	})

	schedulingCtrl := controllers.NewSchedulingController(appInstance)
	router.PUT("/admin/rules/:ruleId", schedulingCtrl.UpdateRule)

	newStartDate := now.AddDate(0, 1, 14).Format("2006-01-02")
	newEndDate := now.AddDate(0, 6, 0).Format("2006-01-02")

	t.Logf("預期更新: 新規則從 %s 開始，原規則截斷到前一天", newStartDate)

	body := map[string]interface{}{
		"name":        rule.Name,
		"offering_id": offering.ID,
		"room_id":     room.ID,
		"start_time":  "15:00",
		"end_time":    "16:00",
		"duration":    60,
		"weekdays":    []int{2},
		"start_date":  newStartDate,
		"end_date":    newEndDate,
		"update_mode": "FUTURE",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/rules/%d", rule.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response global.ApiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	rules := response.Datas.([]interface{})
	// FUTURE 模式應該產生：1條截斷的舊規則 + 1條新規則 = 2條
	expectedCount := 2
	actualCount := len(rules)

	t.Logf("預期規則數量: %d (1條截斷 + 1條新規則)", expectedCount)
	t.Logf("實際規則數量: %d", actualCount)

	if actualCount != expectedCount {
		t.Errorf("規則數量不符: 預期=%d, 實際=%d", expectedCount, actualCount)
	}

	// 輸出各規則詳情
	for i, ruleData := range rules {
		er := ruleData.(map[string]interface{})["effective_range"].(map[string]interface{})
		t.Logf("規則 %d: StartDate=%s, EndDate=%s", i+1, er["start_date"], er["end_date"])
	}

	t.Logf("=== FUTURE Mode 測試結果: %s ===", mapResultToText(actualCount == expectedCount))
}

// TestScheduleRuleUpdateMode_All 測試 ALL 模式
func TestScheduleRuleUpdateMode_All(t *testing.T) {
	appInstance, _ := setupTestAppWithRealDB()

	ctx := context.Background()
	now := time.Now()

	var center models.Center
	if err := appInstance.MySQL.RDB.WithContext(ctx).Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("No center data available, skipping test: %v", err)
		return
	}

	var course models.Course
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&course).Error; err != nil {
		t.Skipf("No course data available, skipping test: %v", err)
		return
	}

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 建立獨立的 offering 確保測試隔離
	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	testStartDate := now.AddDate(0, 1, 0)
	testEndDate := now.AddDate(0, 3, 0)
	testName := fmt.Sprintf("Test Dance %d", now.UnixNano())

	// 建立兩條相關規則（相同 offering, 相同 time, 不同 weekday）
	rules := []models.ScheduleRule{
		{
			CenterID:   center.ID,
			OfferingID: offering.ID,
			RoomID:     room.ID,
			Name:       testName,
			Weekday:    1,
			StartTime:  "10:00",
			EndTime:    "11:00",
			Duration:   60,
			EffectiveRange: models.DateRange{
				StartDate: testStartDate,
				EndDate:   testEndDate,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			CenterID:   center.ID,
			OfferingID: offering.ID,
			RoomID:     room.ID,
			Name:       testName,
			Weekday:    3,
			StartTime:  "10:00",
			EndTime:    "11:00",
			Duration:   60,
			EffectiveRange: models.DateRange{
				StartDate: testStartDate,
				EndDate:   testEndDate,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for i := range rules {
		if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rules[i]).Error; err != nil {
			t.Fatalf("Failed to create test rule %d: %v", i, err)
		}
	}

	t.Logf("=== Test: ALL Mode ===")
	t.Logf("原始規則: %d 條相關規則", len(rules))
	for i, r := range rules {
		t.Logf("  規則%d: ID=%d, Weekday=%d, StartDate=%s, EndDate=%s",
			i+1, r.ID, r.Weekday,
			r.EffectiveRange.StartDate.Format("2006-01-02"),
			r.EffectiveRange.EndDate.Format("2006-01-02"))
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, center.ID)
		c.Set(global.UserIDKey, uint(1))
		c.Next()
	})

	schedulingCtrl := controllers.NewSchedulingController(appInstance)
	router.PUT("/admin/rules/:ruleId", schedulingCtrl.UpdateRule)

	newStartDate := now.AddDate(0, 1, 14).Format("2006-01-02")
	newEndDate := now.AddDate(0, 4, 0).Format("2006-01-02")

	t.Logf("預期更新: 所有 %d 條規則的日期範圍都更新為 %s ~ %s", len(rules), newStartDate, newEndDate)

	body := map[string]interface{}{
		"name":        testName,
		"offering_id": offering.ID,
		"room_id":     room.ID,
		"start_time":  "11:00",
		"end_time":    "12:00",
		"duration":    60,
		"weekdays":    []int{1, 3},
		"start_date":  newStartDate,
		"end_date":    newEndDate,
		"update_mode": "ALL",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/rules/%d", rules[0].ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response global.ApiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	resultRules := response.Datas.([]interface{})
	// ALL 模式應該更新所有相關規則（這裡是 2 條）
	expectedCount := 2
	actualCount := len(resultRules)

	t.Logf("預期規則數量: %d (所有相關規則)", expectedCount)
	t.Logf("實際規則數量: %d", actualCount)

	if actualCount != expectedCount {
		t.Errorf("規則數量不符: 預期=%d, 實際=%d", expectedCount, actualCount)
	}

	// 驗證所有規則都有新的日期範圍
	expectedStartDate := newStartDate + "T00:00:00Z"
	expectedEndDate := newEndDate + "T00:00:00Z"

	allMatch := true
	for i, ruleData := range resultRules {
		er := ruleData.(map[string]interface{})["effective_range"].(map[string]interface{})
		actualStartDate := er["start_date"].(string)
		actualEndDate := er["end_date"].(string)

		t.Logf("規則 %d: StartDate=%s, EndDate=%s", i+1, actualStartDate, actualEndDate)

		if actualStartDate != expectedStartDate || actualEndDate != expectedEndDate {
			allMatch = false
		}
	}

	if !allMatch {
		t.Errorf("部分規則日期範圍不符合預期")
	}

	t.Logf("=== ALL Mode 測試結果: %s ===", mapResultToText(actualCount == expectedCount && allMatch))
}

func mapResultToText(passed bool) string {
	if passed {
		return "✅ 通過"
	}
	return "❌ 失敗"
}
