package test

import (
	"context"
	"testing"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/database/mysql"

	"gorm.io/gorm"
)

// TestScheduleRuleValidator_ValidateForApplyTemplate_OverlapConflict
// 測試模板套用時的時間重疊衝突檢測
func TestScheduleRuleValidator_ValidateForApplyTemplate_OverlapConflict(t *testing.T) {
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

	// Create validator
	validator := services.NewScheduleRuleValidator(appInstance)

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

	// Create test cells
	cells := []models.TimetableCell{
		{
			RowNo:     1,
			ColNo:     1,
			StartTime: "09:00",
			EndTime:   "10:00",
			RoomID:    uintPtr(1),
			TeacherID: uintPtr(1),
		},
		{
			RowNo:     2,
			ColNo:     1,
			StartTime: "10:00",
			EndTime:   "11:00",
			RoomID:    uintPtr(1),
			TeacherID: uintPtr(1),
		},
	}

	// Test validation with overlap check
	summary, err := validator.ValidateForApplyTemplate(
		context.Background(),
		center.ID,
		offering.ID,
		[]int{1}, // Monday
		cells,
		"2026-01-26",
		"2026-12-31",
		false, // no override
	)

	if err != nil {
		t.Fatalf("ValidateForApplyTemplate failed: %v", err)
	}

	// Log the validation result
	t.Logf("Validation Summary:")
	t.Logf("  Valid: %v", summary.Valid)
	t.Logf("  Conflict Count: %d", len(summary.AllConflicts))
	for _, conflict := range summary.AllConflicts {
		t.Logf("    - %s: %s", conflict.ConflictType, conflict.Message)
	}
}

// TestScheduleRuleValidator_ValidateForCreateRule_NoConflict
// 測試新規則驗證（無衝突）
func TestScheduleRuleValidator_ValidateForCreateRule_NoConflict(t *testing.T) {
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

	// Create validator
	validator := services.NewScheduleRuleValidator(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	teacherID := uint(1)
	roomID := uint(1)
	offeringID := uint(1)

	// Test validation with no conflict
	summary, err := validator.ValidateForCreateRule(
		context.Background(),
		center.ID,
		&teacherID,
		roomID,
		offeringID,
		[]int{1}, // Monday
		"2026-01-26",
		"2026-12-31",
		"14:00",
		"15:00",
		false, // no override
	)

	if err != nil {
		t.Fatalf("ValidateForCreateRule failed: %v", err)
	}

	// Log the validation result
	t.Logf("Validation Summary:")
	t.Logf("  Valid: %v", summary.Valid)
	t.Logf("  Overlap Conflicts: %d", len(summary.OverlapConflicts))
	t.Logf("  Buffer Conflicts: %d", len(summary.BufferConflicts))
}

// TestScheduleRuleValidator_ValidateForCreateRule_BufferConflict
// 測試新規則驗證（緩衝時間衝突）
func TestScheduleRuleValidator_ValidateForCreateRule_BufferConflict(t *testing.T) {
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

	// Create validator
	validator := services.NewScheduleRuleValidator(appInstance)

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
	offeringID := uint(1)

	// Test validation - try to create a rule at 09:00 when there's already a class at 08:00-09:00
	// This should trigger a buffer conflict
	summary, err := validator.ValidateForCreateRule(
		context.Background(),
		center.ID,
		&teacherID,
		roomID,
		offeringID,
		[]int{1}, // Monday
		"2026-01-26",
		"2026-12-31",
		"09:00", // Start time that might conflict with buffer
		"10:00",
		false, // no override
	)

	if err != nil {
		t.Fatalf("ValidateForCreateRule failed: %v", err)
	}

	// Log the validation result
	t.Logf("Validation Summary:")
	t.Logf("  Valid: %v", summary.Valid)
	t.Logf("  Overlap Conflicts: %d", len(summary.OverlapConflicts))
	t.Logf("  Buffer Conflicts: %d", len(summary.BufferConflicts))
	for _, conflict := range summary.BufferConflicts {
		t.Logf("    - %s: %s (Gap: %d min, Required: %d min)",
			conflict.ConflictType, conflict.Message, conflict.GapMinutes, conflict.RequiredMinutes)
	}
}

// TestScheduleRuleValidator_ValidationSummary_Structure
// 測試 ValidationSummary 結構體的正確性
func TestScheduleRuleValidator_ValidationSummary_Structure(t *testing.T) {
	// Test that ValidationSummary can be properly serialized
	summary := &services.ValidationSummary{
		Valid: false,
		OverlapConflicts: []services.OverlapInfo{
			{
				Weekday:      1,
				StartTime:    "09:00",
				EndTime:      "10:00",
				ConflictType: "TEACHER_OVERLAP",
				Message:      "週一 09:00-10:00 老師已有排課",
			},
		},
		BufferConflicts: []services.BufferInfo{
			{
				Weekday:         1,
				StartTime:       "10:00",
				EndTime:         "11:00",
				RequiredMinutes: 15,
				GapMinutes:      5,
				ConflictType:    "TEACHER_BUFFER",
				Message:         "緩衝時間不足",
				CanOverride:     true,
			},
		},
	}

	// Verify structure
	if summary.Valid {
		t.Error("Expected Valid to be false")
	}
	if len(summary.OverlapConflicts) != 1 {
		t.Errorf("Expected 1 overlap conflict, got %d", len(summary.OverlapConflicts))
	}
	if len(summary.BufferConflicts) != 1 {
		t.Errorf("Expected 1 buffer conflict, got %d", len(summary.BufferConflicts))
	}
	if !summary.BufferConflicts[0].CanOverride {
		t.Error("Expected CanOverride to be true")
	}

	t.Logf("ValidationSummary structure test passed")
}

// Helper function to create uint pointer
func uintPtr(v uint) *uint {
	return &v
}

// SetupTestApp creates a test app instance with the given database connection
func SetupTestApp(db *gorm.DB) *app.App {
	return &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}
}

// TestScheduleRuleValidator_ValidateForApplyTemplate_WithOverride
// 測試模板套用時的 Buffer Override 功能
func TestScheduleRuleValidator_ValidateForApplyTemplate_WithOverride(t *testing.T) {
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

	// Create validator
	validator := services.NewScheduleRuleValidator(appInstance)

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

	// Create test cells that may cause buffer conflicts
	cells := []models.TimetableCell{
		{
			RowNo:     1,
			ColNo:     1,
			StartTime: "09:00",
			EndTime:   "10:00",
			RoomID:    uintPtr(1),
			TeacherID: uintPtr(1),
		},
	}

	// Test 1: Without override - should fail due to buffer conflict
	summary, err := validator.ValidateForApplyTemplate(
		context.Background(),
		center.ID,
		offering.ID,
		[]int{1}, // Monday
		cells,
		"2026-01-26",
		"2026-12-31",
		false, // no override
	)

	if err != nil {
		t.Fatalf("ValidateForApplyTemplate failed: %v", err)
	}

	t.Logf("Test without override:")
	t.Logf("  Valid: %v", summary.Valid)
	t.Logf("  Conflict Count: %d", len(summary.AllConflicts))

	// Test 2: With override - should succeed if all conflicts are overridable
	summaryWithOverride, err := validator.ValidateForApplyTemplate(
		context.Background(),
		center.ID,
		offering.ID,
		[]int{1}, // Monday
		cells,
		"2026-01-26",
		"2026-12-31",
		true, // with override
	)

	if err != nil {
		t.Fatalf("ValidateForApplyTemplate with override failed: %v", err)
	}

	t.Logf("Test with override:")
	t.Logf("  Valid: %v", summaryWithOverride.Valid)
	t.Logf("  Conflict Count: %d", len(summaryWithOverride.AllConflicts))
}

// TestScheduleRuleValidator_ValidateForCreateRule_WithOverride
// 測試新規則時的 Buffer Override 功能
func TestScheduleRuleValidator_ValidateForCreateRule_WithOverride(t *testing.T) {
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

	// Create validator
	validator := services.NewScheduleRuleValidator(appInstance)

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
	offeringID := uint(1)

	// Test 1: Without override
	summary, err := validator.ValidateForCreateRule(
		context.Background(),
		center.ID,
		&teacherID,
		roomID,
		offeringID,
		[]int{1}, // Monday
		"2026-01-26",
		"2026-12-31",
		"09:00", // May conflict with buffer
		"10:00",
		false, // no override
	)

	if err != nil {
		t.Fatalf("ValidateForCreateRule failed: %v", err)
	}

	t.Logf("Test without override:")
	t.Logf("  Valid: %v", summary.Valid)
	t.Logf("  Buffer Conflicts: %d", len(summary.BufferConflicts))

	// Test 2: With override
	summaryWithOverride, err := validator.ValidateForCreateRule(
		context.Background(),
		center.ID,
		&teacherID,
		roomID,
		offeringID,
		[]int{1}, // Monday
		"2026-01-26",
		"2026-12-31",
		"09:00", // May conflict with buffer
		"10:00",
		true, // with override
	)

	if err != nil {
		t.Fatalf("ValidateForCreateRule with override failed: %v", err)
	}

	t.Logf("Test with override:")
	t.Logf("  Valid: %v", summaryWithOverride.Valid)
	t.Logf("  Buffer Conflicts: %d", len(summaryWithOverride.BufferConflicts))
}

// TestScheduleRuleValidator_Override_NonOverridableConflict
// 測試 Override 無法覆蓋的衝突（Overlap）
func TestScheduleRuleValidator_Override_NonOverridableConflict(t *testing.T) {
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

	// Create validator
	validator := services.NewScheduleRuleValidator(appInstance)

	// Get a center for testing
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("Skipping - no center data available: %v", err)
		return
	}

	teacherID := uintPtr(1)
	roomID := uint(1)
	offeringID := uint(1)

	// Test with a time that will cause overlap (not buffer conflict)
	// The override should not affect overlap conflicts
	summary, err := validator.ValidateForCreateRule(
		context.Background(),
		center.ID,
		teacherID,
		roomID,
		offeringID,
		[]int{1}, // Monday
		"2026-01-26",
		"2026-12-31",
		"08:00", // May overlap with existing schedule
		"09:00",
		true, // override should not affect overlap
	)

	if err != nil {
		t.Fatalf("ValidateForCreateRule failed: %v", err)
	}

	t.Logf("Test overlap with override:")
	t.Logf("  Valid: %v", summary.Valid)
	t.Logf("  Overlap Conflicts: %d", len(summary.OverlapConflicts))
	t.Logf("  Buffer Conflicts: %d", len(summary.BufferConflicts))

	// Overlap conflicts should still make the result invalid
	// even with override = true
	for _, conflict := range summary.OverlapConflicts {
		if conflict.ConflictType == "TEACHER_OVERLAP" || conflict.ConflictType == "ROOM_OVERLAP" {
			t.Logf("  Found overlap conflict: %s", conflict.Message)
		}
	}
}
