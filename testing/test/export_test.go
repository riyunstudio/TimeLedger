package test

import (
	"testing"
	"time"

	"timeLedger/app/models"
)

func TestExportService_CourseCreation(t *testing.T) {
	db, err := InitializeTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer CloseDB(db)

	if err := db.AutoMigrate(&models.Center{}, &models.Course{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	center := &models.Center{
		Name:      "Test Center",
		PlanLevel: "PRO",
	}
	db.Create(center)

	course := &models.Course{
		CenterID:        center.ID,
		Name:            "Test Course",
		DefaultDuration: 60,
		ColorHex:        "#FF5733",
		RoomBufferMin:   10,
		TeacherBufferMin: 5,
		IsActive:        true,
	}
	db.Create(course)

	if course.ID == 0 {
		t.Error("Expected course ID to be set")
	}

	if course.Name != "Test Course" {
		t.Errorf("Expected course name 'Test Course', got '%s'", course.Name)
	}
}

func TestExportService_DateRange(t *testing.T) {
	dateRange := models.DateRange{
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	if dateRange.StartDate.IsZero() {
		t.Error("Expected StartDate to be set")
	}

	if dateRange.EndDate.IsZero() {
		t.Error("Expected EndDate to be set")
	}

	if dateRange.StartDate.After(dateRange.EndDate) {
		t.Error("Expected StartDate to be before EndDate")
	}
}
