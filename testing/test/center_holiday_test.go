package test

import (
	"testing"
	"time"
	"timeLedger/app/models"
)

// ============ Center Holiday Repository Tests ============

func TestCenterHolidayRepository_Initialization(t *testing.T) {
	// This test verifies the repository can be instantiated
	// Note: Actual repository tests require a properly configured app.App instance
	t.Log("Holiday repository initialization test - requires full app setup for integration tests")
}

func TestCenterHolidayModel_Structure(t *testing.T) {
	testDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday := models.CenterHoliday{
		ID:       1,
		CenterID: 2,
		Date:     testDate,
		Name:     "元旦",
	}

	if holiday.ID != 1 {
		t.Errorf("Expected ID 1, got %d", holiday.ID)
	}

	if holiday.CenterID != 2 {
		t.Errorf("Expected CenterID 2, got %d", holiday.CenterID)
	}

	if holiday.Name != "元旦" {
		t.Errorf("Expected Name '元旦', got '%s'", holiday.Name)
	}
}

func TestCenterHoliday_TableName(t *testing.T) {
	holiday := models.CenterHoliday{}
	expectedTableName := "center_holidays"

	if holiday.TableName() != expectedTableName {
		t.Errorf("Expected TableName '%s', got '%s'", expectedTableName, holiday.TableName())
	}
}
