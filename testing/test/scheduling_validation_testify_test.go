package test

import (
	"context"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/database/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScheduleRuleRepository 模擬 ScheduleRuleRepository
type MockScheduleRuleRepository struct {
	mock.Mock
}

func (m *MockScheduleRuleRepository) GetByID(ctx context.Context, id uint) (models.ScheduleRule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.ScheduleRule), args.Error(1)
}

func (m *MockScheduleRuleRepository) FindByDateRange(ctx context.Context, centerID uint, offeringID uint, startDate, endDate time.Time, weekdays []int) ([]models.ScheduleRule, error) {
	args := m.Called(ctx, centerID, offeringID, startDate, endDate, weekdays)
	return args.Get(0).([]models.ScheduleRule), args.Error(1)
}

func (m *MockScheduleRuleRepository) FindOverlapping(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, weekday int, excludeRuleID *uint) ([]models.ScheduleRule, error) {
	args := m.Called(ctx, centerID, teacherID, roomID, startTime, endTime, weekday, excludeRuleID)
	return args.Get(0).([]models.ScheduleRule), args.Error(1)
}

// MockCourseRepository 模擬 CourseRepository
type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id uint) (models.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Course), args.Error(1)
}

// MockRoomRepository 模擬 RoomRepository
type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) GetByID(ctx context.Context, id uint) (models.Room, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Room), args.Error(1)
}

// MockTeacherRepository 模擬 TeacherRepository
type MockTeacherRepository struct {
	mock.Mock
}

func (m *MockTeacherRepository) GetByID(ctx context.Context, id uint) (models.Teacher, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Teacher), args.Error(1)
}

// TestScheduleValidation_CheckOverlap_NoConflict
// 測試時段檢查：無衝突情況
func TestScheduleValidation_CheckOverlap_NoConflict(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	teacherID := uint(1)
	roomID := uint(1)
	loc := app.GetTaiwanLocation()

	// Test with a time that should have no conflict
	// Using a time range that doesn't overlap with existing schedules
	startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc) // Monday 14:00
	endTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc)   // Monday 15:00
	weekday := 1 // Monday

	result, err := validationService.CheckOverlap(
		context.Background(),
		center.ID,
		&teacherID,
		roomID,
		startTime,
		endTime,
		weekday,
		nil,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid, "Expected no overlap conflicts")
	assert.Empty(t, result.Conflicts, "Expected no conflicts to be reported")
}

// TestScheduleValidation_CheckOverlap_TeacherOverlap
// 測試時段檢查：老師時間重疊
func TestScheduleValidation_CheckOverlap_TeacherOverlap(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get a teacher who has existing schedule
	var teacher models.User
	if err := db.Where("user_type = ?", "TEACHER").Order("id ASC").First(&teacher).Error; err != nil {
		t.Skipf("Skipping - no teacher data available: %v", err)
		return
	}

	teacherID := teacher.ID
	roomID := uint(1)
	loc := app.GetTaiwanLocation()

	// Try to create a rule at 08:00-09:00 on Monday
	// This should trigger an overlap conflict if there's already a class at that time
	startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc)  // Monday 08:00
	endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)    // Monday 09:00
	weekday := 1                                           // Monday

	result, err := validationService.CheckOverlap(
		context.Background(),
		center.ID,
		&teacherID,
		roomID,
		startTime,
		endTime,
		weekday,
		nil,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Log the result for debugging
	t.Logf("Overlap check result: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	for _, conflict := range result.Conflicts {
		t.Logf("  Conflict: Type=%s, Message=%s", conflict.Type, conflict.Message)
	}
}

// TestScheduleValidation_CheckOverlap_RoomOverlap
// 測試時段檢查：教室重疊
func TestScheduleValidation_CheckOverlap_RoomOverlap(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	teacherID := uintPtr(1) // Using nil teacherID to focus on room overlap
	roomID := uint(1)
	loc := app.GetTaiwanLocation()

	// Test with overlapping room
	startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc) // Monday 08:00
	endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)   // Monday 09:00
	weekday := 1                                          // Monday

	result, err := validationService.CheckOverlap(
		context.Background(),
		center.ID,
		teacherID,
		roomID,
		startTime,
		endTime,
		weekday,
		nil,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	t.Logf("Room overlap check result: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	for _, conflict := range result.Conflicts {
		t.Logf("  Conflict: Type=%s, Message=%s", conflict.Type, conflict.Message)
		if conflict.Type == "ROOM_OVERLAP" {
			t.Logf("  Room overlap detected correctly")
		}
	}
}

// TestScheduleValidation_CheckTeacherBuffer_Sufficient
// 測試老師緩衝時間：足夠的緩衝
func TestScheduleValidation_CheckTeacherBuffer_Sufficient(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get an offering to get course info
	var offering models.Offering
	if err := db.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	teacherID := uint(1)
	loc := app.GetTaiwanLocation()

	// Test with sufficient gap between sessions
	prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc) // Previous class ends at 10:00
	nextStartTime := time.Date(2026, 1, 27, 11, 30, 0, 0, loc) // Next class starts at 11:30 (90 min gap)

	result, err := validationService.CheckTeacherBuffer(
		context.Background(),
		center.ID,
		teacherID,
		prevEndTime,
		nextStartTime,
		offering.ID,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid, "Expected sufficient buffer time")
	assert.Empty(t, result.Conflicts, "Expected no buffer conflicts")

	t.Logf("Teacher buffer check: Valid=%v (sufficient buffer)", result.Valid)
}

// TestScheduleValidation_CheckTeacherBuffer_Insufficient
// 測試老師緩衝時間：不足的緩衝
func TestScheduleValidation_CheckTeacherBuffer_Insufficient(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get an offering to get course info
	var offering models.Offering
	if err := db.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	teacherID := uint(1)
	loc := app.GetTaiwanLocation()

	// Test with insufficient gap between sessions
	prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc) // Previous class ends at 10:00
	nextStartTime := time.Date(2026, 1, 27, 10, 30, 0, 0, loc) // Next class starts at 10:30 (30 min gap)

	result, err := validationService.CheckTeacherBuffer(
		context.Background(),
		center.ID,
		teacherID,
		prevEndTime,
		nextStartTime,
		offering.ID,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Log the result
	t.Logf("Teacher buffer check: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	for _, conflict := range result.Conflicts {
		t.Logf("  Conflict: Type=%s, Message=%s, CanOverride=%v",
			conflict.Type, conflict.Message, conflict.CanOverride)
	}

	// Buffer conflicts should be overridable
	if !result.Valid {
		for _, conflict := range result.Conflicts {
			assert.Equal(t, "TEACHER_BUFFER", conflict.Type, "Expected TEACHER_BUFFER conflict type")
			assert.True(t, conflict.CanOverride, "Buffer conflicts should be overridable")
		}
	}
}

// TestScheduleValidation_CheckRoomBuffer
// 測試教室緩衝時間
func TestScheduleValidation_CheckRoomBuffer(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get an offering to get course info
	var offering models.Offering
	if err := db.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	roomID := uint(1)
	loc := app.GetTaiwanLocation()

	// Test with insufficient room buffer
	prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)  // Previous class ends at 10:00
	nextStartTime := time.Date(2026, 1, 27, 10, 15, 0, 0, loc) // Next class starts at 10:15 (15 min gap)

	result, err := validationService.CheckRoomBuffer(
		context.Background(),
		center.ID,
		roomID,
		prevEndTime,
		nextStartTime,
		offering.ID,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	t.Logf("Room buffer check: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	for _, conflict := range result.Conflicts {
		t.Logf("  Conflict: Type=%s, Message=%s", conflict.Type, conflict.Message)
	}
}

// TestScheduleValidation_ValidateFull
// 測試完整驗證流程
func TestScheduleValidation_ValidateFull(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get an offering for testing
	var offering models.Offering
	if err := db.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	teacherID := uintPtr(1)
	roomID := uint(1)
	loc := app.GetTaiwanLocation()

	// Test validation with a valid time slot (no conflicts)
	startTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc) // Monday 15:00
	endTime := time.Date(2026, 1, 27, 16, 0, 0, 0, loc)   // Monday 16:00

	result, err := validationService.ValidateFull(
		context.Background(),
		center.ID,
		teacherID,
		roomID,
		offering.ID,
		startTime,
		endTime,
		nil,   // no exclude rule
		false, // no buffer override
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	t.Logf("Full validation result: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	for _, conflict := range result.Conflicts {
		t.Logf("  Conflict: Type=%s, Message=%s, CanOverride=%v",
			conflict.Type, conflict.Message, conflict.CanOverride)
	}

	// Overlap conflicts should NOT be overridable
	for _, conflict := range result.Conflicts {
		if conflict.Type == "TEACHER_OVERLAP" || conflict.Type == "ROOM_OVERLAP" {
			assert.False(t, conflict.CanOverride, "Overlap conflicts should NOT be overridable")
		}
	}
}

// TestScheduleValidation_ValidateFull_WithBufferOverride
// 測試完整驗證流程（含緩衝覆寫）
func TestScheduleValidation_ValidateFull_WithBufferOverride(t *testing.T) {
	// Skip if no database connection
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("Skipping - database connection failed: %v", err)
		return
	}
	defer CloseDB(db)

	// Create app instance
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// Create validation service
	validationService := services.NewScheduleValidationService(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	// Get an offering for testing
	var offering models.Offering
	if err := db.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	teacherID := uintPtr(1)
	roomID := uint(1)
	loc := app.GetTaiwanLocation()

	// Test with buffer override enabled
	// Using a time that might have buffer conflicts
	startTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)  // Monday 09:00
	endTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)   // Monday 10:00

	result, err := validationService.ValidateFull(
		context.Background(),
		center.ID,
		teacherID,
		roomID,
		offering.ID,
		startTime,
		endTime,
		nil,  // no exclude rule
		true, // allow buffer override
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	t.Logf("Full validation with override: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	for _, conflict := range result.Conflicts {
		t.Logf("  Conflict: Type=%s, Message=%s, CanOverride=%v",
			conflict.Type, conflict.Message, conflict.CanOverride)
	}

	// With buffer override, buffer conflicts should not invalidate the result
	// (though the result may still show conflicts, they should be marked as overridable)
	for _, conflict := range result.Conflicts {
		if conflict.Type == "TEACHER_BUFFER" || conflict.Type == "ROOM_BUFFER" {
			assert.True(t, conflict.CanOverride, "Buffer conflicts should be overridable")
		}
	}
}

// TestValidationResult_Structure
// 測試 ValidationResult 結構體
func TestValidationResult_Structure(t *testing.T) {
	result := services.ValidationResult{
		Valid: false,
		Conflicts: []services.ValidationConflict{
			{
				Type:            "TEACHER_OVERLAP",
				Message:         "老師在該時段已有課程安排",
				CanOverride:     false,
				ConflictSource:  "RULE",
				ConflictSourceID: 1,
				Details:         "rule_id:1, offering_id:1",
			},
			{
				Type:            "TEACHER_BUFFER",
				Message:         "老師轉場時間不足，需間隔 15 分鐘，實際間隔 5 分鐘",
				CanOverride:     true,
				RequiredMinutes: 15,
				DiffMinutes:     10,
				ConflictSource:  "SESSION",
			},
		},
	}

	assert.False(t, result.Valid)
	assert.Len(t, result.Conflicts, 2)

	// Check first conflict (overlap - not overridable)
	assert.Equal(t, "TEACHER_OVERLAP", result.Conflicts[0].Type)
	assert.False(t, result.Conflicts[0].CanOverride)
	assert.Equal(t, "RULE", result.Conflicts[0].ConflictSource)

	// Check second conflict (buffer - overridable)
	assert.Equal(t, "TEACHER_BUFFER", result.Conflicts[1].Type)
	assert.True(t, result.Conflicts[1].CanOverride)
	assert.Equal(t, 15, result.Conflicts[1].RequiredMinutes)
	assert.Equal(t, 10, result.Conflicts[1].DiffMinutes)
}

// TestPaginationParams_Validate
// 測試分頁參數驗證
func TestPaginationParams_Validate(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		sortBy   string
		sortOrder string
		expectedPage     int
		expectedLimit    int
		expectedSortBy   string
		expectedSortOrder string
	}{
		{
			name:           "Normal pagination",
			page:           2,
			limit:          25,
			sortBy:         "created_at",
			sortOrder:      "ASC",
			expectedPage:    2,
			expectedLimit:   25,
			expectedSortBy:  "created_at",
			expectedSortOrder: "ASC",
		},
		{
			name:           "Negative page defaults to 1",
			page:           -1,
			limit:          10,
			sortBy:         "",
			sortOrder:      "DESC",
			expectedPage:    1,
			expectedLimit:   10,
			expectedSortBy:  "id",
			expectedSortOrder: "DESC",
		},
		{
			name:           "Zero limit defaults to 20",
			page:           1,
			limit:          0,
			sortBy:         "",
			sortOrder:      "ASC",
			expectedPage:    1,
			expectedLimit:   20,
			expectedSortBy:  "id",
			expectedSortOrder: "ASC",
		},
		{
			name:           "Limit over 100 caps at 100",
			page:           1,
			limit:          200,
			sortBy:         "",
			sortOrder:      "DESC",
			expectedPage:    1,
			expectedLimit:   100,
			expectedSortBy:  "id",
			expectedSortOrder: "DESC",
		},
		{
			name:           "Invalid sort order defaults to DESC",
			page:           1,
			limit:          10,
			sortBy:         "",
			sortOrder:      "INVALID",
			expectedPage:    1,
			expectedLimit:   10,
			expectedSortBy:  "id",
			expectedSortOrder: "DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &services.PaginationParams{
				Page:      tt.page,
				Limit:     tt.limit,
				SortBy:    tt.sortBy,
				SortOrder: tt.sortOrder,
			}

			params.Validate()

			assert.Equal(t, tt.expectedPage, params.Page)
			assert.Equal(t, tt.expectedLimit, params.Limit)
			assert.Equal(t, tt.expectedSortBy, params.SortBy)
			assert.Equal(t, tt.expectedSortOrder, params.SortOrder)
		})
	}
}

// TestPaginationParams_GetOffset
// 測試分頁偏移量計算
func TestPaginationParams_GetOffset(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		limit        int
		expectedOffset int
	}{
		{
			name:           "First page",
			page:           1,
			limit:          20,
			expectedOffset: 0,
		},
		{
			name:           "Second page",
			page:           2,
			limit:          20,
			expectedOffset: 20,
		},
		{
			name:           "Third page with 10 items per page",
			page:           3,
			limit:          10,
			expectedOffset: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &services.PaginationParams{
				Page:  tt.page,
				Limit: tt.limit,
			}
			assert.Equal(t, tt.expectedOffset, params.GetOffset())
		})
	}
}

// TestFilterBuilder
// 測試過濾建構器
func TestFilterBuilder(t *testing.T) {
	t.Run("AddEq", func(t *testing.T) {
		fb := services.NewFilterBuilder()
		fb.AddEq("status", "active")

		conditions := fb.Build()
		assert.Len(t, conditions, 1)
		assert.Equal(t, "status", conditions[0].Field)
		assert.Equal(t, "eq", conditions[0].Operator)
		assert.Equal(t, "active", conditions[0].Value)
	})

	t.Run("AddIn", func(t *testing.T) {
		fb := services.NewFilterBuilder()
		fb.AddIn("status", []interface{}{"active", "pending"})

		conditions := fb.Build()
		assert.Len(t, conditions, 1)
		assert.Equal(t, "in", conditions[0].Operator)
	})

	t.Run("AddBetween", func(t *testing.T) {
		fb := services.NewFilterBuilder()
		fb.AddBetween("created_at", "2026-01-01", "2026-12-31")

		conditions := fb.Build()
		assert.Len(t, conditions, 1)
		assert.Equal(t, "between", conditions[0].Operator)
	})

	t.Run("AddCenterScope", func(t *testing.T) {
		fb := services.NewFilterBuilder()
		fb.AddCenterScope(123)

		conditions := fb.Build()
		assert.Len(t, conditions, 1)
		assert.Equal(t, "center_id", conditions[0].Field)
		assert.Equal(t, uint(123), conditions[0].Value)
	})

	t.Run("MethodChaining", func(t *testing.T) {
		fb := services.NewFilterBuilder().
			AddEq("status", "active").
			AddEq("center_id", uint(1)).
			AddLike("name", "test")

		conditions := fb.Build()
		assert.Len(t, conditions, 3)
		assert.False(t, fb.IsEmpty())
	})

	t.Run("IsEmpty", func(t *testing.T) {
		fb := services.NewFilterBuilder()
		assert.True(t, fb.IsEmpty())

		fb.AddEq("id", 1)
		assert.False(t, fb.IsEmpty())
	})
}
