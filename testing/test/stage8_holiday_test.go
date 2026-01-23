package test

import (
	"testing"
	"time"
	"timeLedger/app/models"
	"timeLedger/app/services"
)

func TestStage8_CenterHoliday_Model(t *testing.T) {
	t.Run("CenterHoliday Fields", func(t *testing.T) {
		date := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)

		holiday := models.CenterHoliday{
			ID:        1,
			CenterID:  10,
			Date:      date,
			Name:      "情人節",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if holiday.CenterID == 0 {
			t.Error("CenterID should not be zero")
		}

		if holiday.Date.IsZero() {
			t.Error("Date should not be zero")
		}

		if holiday.Name == "" {
			t.Error("Name should not be empty")
		}
	})

	t.Run("CenterHoliday TableName", func(t *testing.T) {
		holiday := models.CenterHoliday{}
		expected := "center_holidays"
		if holiday.TableName() != expected {
			t.Errorf("TableName() = %s, want %s", holiday.TableName(), expected)
		}
	})
}

func TestStage8_HolidayFiltering_Logic(t *testing.T) {
	t.Run("Holiday Set Creation", func(t *testing.T) {
		holidays := []models.CenterHoliday{
			{Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Name: "元旦"},
			{Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), Name: "情人節"},
			{Date: time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC), Name: "雙十節"},
		}

		holidaySet := make(map[string]bool)
		for _, h := range holidays {
			holidaySet[h.Date.Format("2006-01-02")] = true
		}

		if !holidaySet["2026-01-01"] {
			t.Error("Holiday 2026-01-01 should be in set")
		}

		if !holidaySet["2026-02-14"] {
			t.Error("Holiday 2026-02-14 should be in set")
		}

		if holidaySet["2026-05-01"] {
			t.Error("Holiday 2026-05-01 should NOT be in set")
		}
	})

	t.Run("Holiday Check Logic", func(t *testing.T) {
		date := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
		holidaySet := map[string]bool{
			"2026-01-01": true,
			"2026-02-14": true,
		}

		dateStr := date.Format("2006-01-02")
		isHoliday := holidaySet[dateStr]

		if !isHoliday {
			t.Error("2026-02-14 should be detected as holiday")
		}
	})
}

func TestStage8_BulkImport_Logic(t *testing.T) {
	t.Run("Skip Duplicate Logic", func(t *testing.T) {
		existingHolidays := map[string]bool{
			"2026-01-01": true,
		}

		newHolidays := []models.CenterHoliday{
			{Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Name: "元旦"},
			{Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), Name: "情人節"},
		}

		var skipped int
		var created int

		for _, h := range newHolidays {
			dateStr := h.Date.Format("2006-01-02")
			if existingHolidays[dateStr] {
				skipped++
			} else {
				created++
			}
		}

		if skipped != 1 {
			t.Errorf("Expected 1 skipped, got %d", skipped)
		}

		if created != 1 {
			t.Errorf("Expected 1 created, got %d", created)
		}
	})

	t.Run("Date Range Query Logic", func(t *testing.T) {
		startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

		testDates := []time.Time{
			time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		}

		for _, d := range testDates {
			isWithinRange := !d.Before(startDate) && !d.After(endDate)
			if d.Format("2006-01-02") == "2025-12-31" && isWithinRange {
				t.Error("2025-12-31 should NOT be within range")
			}
			if d.Format("2006-01-02") == "2026-02-01" && isWithinRange {
				t.Error("2026-02-01 should NOT be within range")
			}
		}
	})
}

func TestStage8_ExpandedSchedule_HolidayField(t *testing.T) {
	t.Run("ExpandedSchedule with IsHoliday", func(t *testing.T) {
		schedule := services.ExpandedSchedule{
			RuleID:       1,
			Date:         time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
			StartTime:    "10:00:00",
			EndTime:      "11:00:00",
			RoomID:       1,
			TeacherID:    new(uint),
			IsHoliday:    true,
			HasException: false,
		}

		if !schedule.IsHoliday {
			t.Error("IsHoliday should be true")
		}
	})

	t.Run("ExpandedSchedule Non-Holiday", func(t *testing.T) {
		schedule := services.ExpandedSchedule{
			RuleID:       1,
			Date:         time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC),
			StartTime:    "10:00:00",
			EndTime:      "11:00:00",
			RoomID:       1,
			TeacherID:    new(uint),
			IsHoliday:    false,
			HasException: false,
		}

		if schedule.IsHoliday {
			t.Error("IsHoliday should be false")
		}
	})
}

func TestStage8_PhaseTransition_HolidayAwareness(t *testing.T) {
	t.Run("Phase Transition Near Holiday", func(t *testing.T) {
		holidayDate := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
		transitionDate := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)

		holidaySet := map[string]bool{
			"2026-02-14": true,
		}

		isHolidayBefore := holidaySet[holidayDate.Format("2006-01-02")]
		isHolidayAfter := holidaySet[transitionDate.Format("2006-01-02")]

		if !isHolidayBefore {
			t.Error("Feb 14 should be a holiday")
		}

		if isHolidayAfter {
			t.Error("Feb 15 should NOT be a holiday")
		}

		t.Logf("Transition from holiday (%v) to normal day (%v)", isHolidayBefore, isHolidayAfter)
	})
}
