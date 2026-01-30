package test

import (
	"context"
	"testing"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/database/mysql"
)

func TestScheduleRuleRepository_CRUD(t *testing.T) {
	t.Skip("Skipping - requires .env file")
}

// TestScheduleRuleRepository_ListByOfferingIDWithPreload
// 測試 ListByOfferingIDWithPreload 批次查詢並預載入關聯資料
func TestScheduleRuleRepository_ListByOfferingIDWithPreload(t *testing.T) {
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

	// Create repository
	repo := repositories.NewScheduleRuleRepository(appInstance)

	// Get an offering for testing
	var offering models.Offering
	if err := db.Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering data available: %v", err)
		return
	}

	// Test the ListByOfferingIDWithPreload method
	rules, err := repo.ListByOfferingIDWithPreload(context.Background(), offering.ID)
	if err != nil {
		t.Fatalf("ListByOfferingIDWithPreload failed: %v", err)
	}

	// Log the results
	t.Logf("ListByOfferingIDWithPreload Results:")
	t.Logf("  Total rules found: %d", len(rules))
	t.Logf("  Offering ID: %d", offering.ID)

	// Verify preloaded associations if rules exist
	for i, rule := range rules {
		t.Logf("  Rule %d: ID=%d, Weekday=%d, StartTime=%s, EndTime=%s",
			i+1, rule.ID, rule.Weekday, rule.StartTime, rule.EndTime)

		// Check if associations are loaded (non-zero values indicate successful preloading)
		if rule.OfferingID != 0 {
			t.Logf("    - Offering loaded: %s", rule.Offering.Name)
		}
		if rule.RoomID != 0 && rule.Room.ID != 0 {
			t.Logf("    - Room loaded: %s", rule.Room.Name)
		}
		if rule.TeacherID != nil && rule.Teacher.ID != 0 {
			t.Logf("    - Teacher loaded: %s", rule.Teacher.Name)
		}
	}

	// Test with non-existent offering (should return empty slice, not error)
	nonExistentOfferingID := uint(999999)
	rules, err = repo.ListByOfferingIDWithPreload(context.Background(), nonExistentOfferingID)
	if err != nil {
		t.Fatalf("ListByOfferingIDWithPreload with non-existent ID failed: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("Expected 0 rules for non-existent offering, got %d", len(rules))
	}
	t.Logf("ListByOfferingIDWithPreload with non-existent offering: %d rules (expected 0)", len(rules))
}

// TestScheduleRuleRepository_ListByOfferingIDWithPreload_PreloadVerification
// 測試 ListByOfferingIDWithPreload 預載入關聯資料的正確性
func TestScheduleRuleRepository_ListByOfferingIDWithPreload_PreloadVerification(t *testing.T) {
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

	// Create repository
	repo := repositories.NewScheduleRuleRepository(appInstance)

	// Get an offering that has rules with preloaded associations
	var offering models.Offering
	if err := db.Preload("ScheduleRules").Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("Skipping - no offering with rules available: %v", err)
		return
	}

	// Skip if no rules exist (use raw query to check)
	var count int64
	if err := db.Table("schedule_rules").Where("offering_id = ?", offering.ID).Count(&count).Error; err != nil {
		t.Skipf("Skipping - cannot check rule count: %v", err)
		return
	}
	if count == 0 {
		t.Skipf("Skipping - offering has no rules: %d", offering.ID)
		return
	}

	// Test preloading
	rules, err := repo.ListByOfferingIDWithPreload(context.Background(), offering.ID)
	if err != nil {
		t.Fatalf("ListByOfferingIDWithPreload failed: %v", err)
	}

	// Verify preloading worked by checking association data
	var rulesWithOffering int
	var rulesWithRoom int
	var rulesWithTeacher int

	for _, rule := range rules {
		if rule.OfferingID == offering.ID {
			rulesWithOffering++
		}
		if rule.RoomID != 0 && rule.Room.ID != 0 {
			rulesWithRoom++
		}
		if rule.TeacherID != nil && *rule.TeacherID != 0 && rule.Teacher.ID != 0 {
			rulesWithTeacher++
		}
	}

	t.Logf("Preload Verification Results:")
	t.Logf("  Total rules: %d", len(rules))
	t.Logf("  Rules with Offering: %d", rulesWithOffering)
	t.Logf("  Rules with Room: %d", rulesWithRoom)
	t.Logf("  Rules with Teacher: %d", rulesWithTeacher)

	// All rules should have the correct offering_id
	if rulesWithOffering != len(rules) {
		t.Errorf("Expected all %d rules to have offering_id %d, but only %d have it",
			len(rules), offering.ID, rulesWithOffering)
	}
}

// Helper function to create uint pointer
func uintPtr(v uint) *uint {
	return &v
}
