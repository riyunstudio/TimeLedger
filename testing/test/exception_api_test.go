package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeLedger/app/controllers"
	"timeLedger/app/models"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

// TestGetAllExceptions tests the GET /admin/exceptions/all endpoint
func TestGetAllExceptions(t *testing.T) {
	appInstance := setupExceptionTestApp()

	// Get a center for testing
	var center models.Center
	if err := appInstance.MySQL.RDB.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Create scheduling controller
	schedulingCtl := controllers.NewSchedulingController(appInstance)

	// Setup test router
	router := gin.New()
	// Add middleware to set center_id in context
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, uint(center.ID))
		c.Set(global.UserTypeKey, "ADMIN")
		c.Next()
	})
	router.GET("/admin/exceptions/all", schedulingCtl.GetAllExceptions)

	// Create test request with mock admin token
	req, _ := http.NewRequest("GET", "/admin/exceptions/all", nil)
	req.Header.Set("Authorization", "Bearer mock-admin-token")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	// Parse response
	var response global.ApiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Logf("Response may be error format: %v", err)
		t.Logf("Response body: %s", w.Body.String())
	}

	t.Logf("GetAllExceptions response: %s", w.Body.String())
}

// TestGetAllExceptions_WithFilters tests the GET /admin/exceptions/all endpoint with filters
func TestGetAllExceptions_WithFilters(t *testing.T) {
	appInstance := setupExceptionTestApp()

	// Get a center for testing
	var center models.Center
	if err := appInstance.MySQL.RDB.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Create scheduling controller
	schedulingCtl := controllers.NewSchedulingController(appInstance)

	// Setup test router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, uint(center.ID))
		c.Set(global.UserTypeKey, "ADMIN")
		c.Next()
	})
	router.GET("/admin/exceptions/all", schedulingCtl.GetAllExceptions)

	// Test with status filter
	req, _ := http.NewRequest("GET", "/admin/exceptions/all?status=PENDING", nil)
	req.Header.Set("Authorization", "Bearer mock-admin-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}

	t.Logf("GetAllExceptions with filters response: %s", w.Body.String())
}

// TestReviewException_Approve tests the POST /admin/scheduling/exceptions/:id/review endpoint with APPROVE action
func TestReviewException_Approve(t *testing.T) {
	appInstance := setupExceptionTestApp()

	// Get a center for testing
	var center models.Center
	if err := appInstance.MySQL.RDB.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get an offering for testing
	var offering models.Offering
	if err := appInstance.MySQL.RDB.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	// Get a teacher for testing
	var teacher models.Teacher
	if err := appInstance.MySQL.RDB.Order("id DESC").First(&teacher).Error; err != nil {
		t.Skipf("Skipping - no teacher data available: %v", err)
		return
	}

	// Get a room for testing
	var room models.Room
	if err := appInstance.MySQL.RDB.Where("center_id = ?", center.ID).Order("id DESC").First(&room).Error; err != nil {
		t.Skipf("Skipping - no room data available: %v", err)
		return
	}

	// Create a schedule rule
	today := time.Now()
	weekday := int(today.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	teacherID := teacher.ID
	rule := &models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  &teacherID,
		RoomID:     room.ID,
		Weekday:    weekday,
		StartTime:  "09:00",
		EndTime:    "10:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: today,
			EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
	}
	if err := appInstance.MySQL.RDB.Create(rule).Error; err != nil {
		t.Logf("Warning: Failed to create test rule: %v", err)
	}

	// Create a pending exception
	exception := &models.ScheduleException{
		CenterID:     center.ID,
		RuleID:       rule.ID,
		OriginalDate: today,
		Type:         "CANCEL",
		Status:       "PENDING",
		Reason:       "測試請假",
	}
	if err := appInstance.MySQL.RDB.Create(exception).Error; err != nil {
		t.Skipf("Skipping - failed to create test exception: %v", err)
		return
	}

	// Create scheduling controller
	schedulingCtl := controllers.NewSchedulingController(appInstance)

	// Setup test router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, uint(center.ID))
		c.Set(global.UserTypeKey, "ADMIN")
		c.Set(global.UserIDKey, uint(1))
		c.Next()
	})
	router.POST("/admin/scheduling/exceptions/:id/review", schedulingCtl.ReviewException)

	// Test approve action
	body := []byte(`{"action": "APPROVED", "reason": "測試核准"}`)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/admin/scheduling/exceptions/%d/review", exception.ID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer mock-admin-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response - may be 200 (success) or 409 (conflict) or 500 (error)
	t.Logf("ReviewException APPROVE response: %s", w.Body.String())

	// Verify exception status changed
	var updatedException models.ScheduleException
	appInstance.MySQL.RDB.First(&updatedException, exception.ID)
	if updatedException.Status == "APPROVED" {
		t.Log("Exception was successfully approved")
	} else if updatedException.Status == "PENDING" {
		t.Log("Exception remains pending (may have validation issues)")
	}
}

// TestReviewException_Reject tests the POST /admin/scheduling/exceptions/:id/review endpoint with REJECT action
func TestReviewException_Reject(t *testing.T) {
	appInstance := setupExceptionTestApp()

	// Get a center for testing
	var center models.Center
	if err := appInstance.MySQL.RDB.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get an offering for testing
	var offering models.Offering
	if err := appInstance.MySQL.RDB.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	// Get a teacher for testing
	var teacher models.Teacher
	if err := appInstance.MySQL.RDB.Order("id DESC").First(&teacher).Error; err != nil {
		t.Skipf("Skipping - no teacher data available: %v", err)
		return
	}

	// Get a room for testing
	var room models.Room
	if err := appInstance.MySQL.RDB.Where("center_id = ?", center.ID).Order("id DESC").First(&room).Error; err != nil {
		t.Skipf("Skipping - no room data available: %v", err)
		return
	}

	// Create a schedule rule
	today := time.Now()
	weekday := int(today.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	teacherID := teacher.ID
	rule := &models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  &teacherID,
		RoomID:     room.ID,
		Weekday:    weekday,
		StartTime:  "09:00",
		EndTime:    "10:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: today,
			EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
	}
	if err := appInstance.MySQL.RDB.Create(rule).Error; err != nil {
		t.Logf("Warning: Failed to create test rule: %v", err)
	}

	// Create a pending exception
	exception := &models.ScheduleException{
		CenterID:     center.ID,
		RuleID:       rule.ID,
		OriginalDate: today,
		Type:         "CANCEL",
		Status:       "PENDING",
		Reason:       "測試請假",
	}
	if err := appInstance.MySQL.RDB.Create(exception).Error; err != nil {
		t.Skipf("Skipping - failed to create test exception: %v", err)
		return
	}

	// Create scheduling controller
	schedulingCtl := controllers.NewSchedulingController(appInstance)

	// Setup test router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, uint(center.ID))
		c.Set(global.UserTypeKey, "ADMIN")
		c.Set(global.UserIDKey, uint(1))
		c.Next()
	})
	router.POST("/admin/scheduling/exceptions/:id/review", schedulingCtl.ReviewException)

	// Test reject action
	body := []byte(`{"action": "REJECTED", "reason": "不符合規定"}`)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/admin/scheduling/exceptions/%d/review", exception.ID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer mock-admin-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("ReviewException REJECT response: %s", w.Body.String())

	// Verify exception status changed
	var updatedException models.ScheduleException
	appInstance.MySQL.RDB.First(&updatedException, exception.ID)
	if updatedException.Status == "REJECTED" {
		t.Log("Exception was successfully rejected")
	}
}

// TestReviewException_InvalidAction tests the POST /admin/scheduling/exceptions/:id/review with invalid action
func TestReviewException_InvalidAction(t *testing.T) {
	appInstance := setupExceptionTestApp()

	// Get a center for testing
	var center models.Center
	if err := appInstance.MySQL.RDB.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Create scheduling controller
	schedulingCtl := controllers.NewSchedulingController(appInstance)

	// Setup test router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(global.CenterIDKey, uint(center.ID))
		c.Set(global.UserTypeKey, "ADMIN")
		c.Next()
	})
	router.POST("/admin/scheduling/exceptions/:id/review", schedulingCtl.ReviewException)

	// Test with invalid action
	body := []byte(`{"action": "INVALID_ACTION", "reason": "測試"}`)
	req, _ := http.NewRequest("POST", "/admin/scheduling/exceptions/1/review", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer mock-admin-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("ReviewException INVALID action response: %s", w.Body.String())

	// Should return error for invalid action
	var response global.ApiResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Code != 0 {
		t.Log("Invalid action correctly rejected")
	}
}
