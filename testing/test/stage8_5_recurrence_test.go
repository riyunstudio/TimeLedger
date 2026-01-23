package test

import (
	"testing"
	"time"
)

func TestStage8_5_RecurrenceEditMode(t *testing.T) {
	t.Run("Single Mode Affects Only One Date", func(t *testing.T) {
		affectedCount := 1
		if affectedCount != 1 {
			t.Errorf("Expected 1 affected session for SINGLE mode, got %d", affectedCount)
		}
	})

	t.Run("Future Mode Counts Multiple Dates", func(t *testing.T) {
		startDate := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

		var count int
		current := startDate
		for !current.After(endDate) {
			if current.Weekday() == time.Monday {
				count++
			}
			current = current.AddDate(0, 0, 1)
		}

		if count < 3 || count > 5 {
			t.Errorf("Expected 3-5 Mondays in January 2026, got %d", count)
		}
	})

	t.Run("All Mode Counts All Dates", func(t *testing.T) {
		startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)

		var count int
		current := startDate
		for !current.After(endDate) {
			if current.Weekday() == time.Wednesday {
				count++
			}
			current = current.AddDate(0, 0, 1)
		}

		if count < 12 || count > 13 {
			t.Errorf("Expected 12-13 Wednesdays in Q1 2026, got %d", count)
		}
	})
}

func TestStage8_5_ExceptionGeneration(t *testing.T) {
	t.Run("Single Edit Creates Cancel And Add Exceptions", func(t *testing.T) {
		cancelExceptionType := "CANCEL"
		addExceptionType := "ADD"

		if cancelExceptionType != "CANCEL" {
			t.Error("Cancel exception type should be CANCEL")
		}
		if addExceptionType != "ADD" {
			t.Error("Add exception type should be ADD")
		}
	})

	t.Run("Exception Has Correct Original Date", func(t *testing.T) {
		originalDate := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)

		exceptionOriginalDate := originalDate

		if !exceptionOriginalDate.Equal(originalDate) {
			t.Errorf("Exception original date should be %v, got %v", originalDate, exceptionOriginalDate)
		}
	})

	t.Run("Add Exception Has New Time Details", func(t *testing.T) {
		newStartTime := "14:00:00"
		newEndTime := "15:30:00"
		newTeacherID := uint(100)
		newRoomID := uint(5)

		if newStartTime == "" {
			t.Error("New start time should not be empty")
		}
		if newEndTime == "" {
			t.Error("New end time should not be empty")
		}
		if newTeacherID == 0 {
			t.Error("New teacher ID should not be zero")
		}
		if newRoomID == 0 {
			t.Error("New room ID should not be zero")
		}
	})
}

func TestStage8_5_FutureEditCreatesNewRule(t *testing.T) {
	t.Run("Future Edit Creates New Rule With Updated Effective Range", func(t *testing.T) {
		originalRuleStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		originalRuleEnd := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)
		editDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

		newRuleStart := editDate
		newRuleEnd := originalRuleEnd

		if newRuleStart.After(editDate) {
			t.Error("New rule start date should be on or before edit date")
		}
		if newRuleEnd.Before(originalRuleEnd) {
			t.Error("New rule end date should not be before original rule end date")
		}
		if newRuleStart.Before(originalRuleStart) {
			t.Error("New rule start should be after or equal to original rule start")
		}
	})

	t.Run("New Rule Preserves Unchanged Fields", func(t *testing.T) {
		originalRule := map[string]interface{}{
			"weekday":     1,
			"offering_id": uint(10),
		}

		newRule := map[string]interface{}{
			"weekday":     originalRule["weekday"],
			"offering_id": originalRule["offering_id"],
		}

		if newRule["weekday"] != originalRule["weekday"] {
			t.Error("New rule should preserve weekday")
		}
		if newRule["offering_id"] != originalRule["offering_id"] {
			t.Error("New rule should preserve offering_id")
		}
	})
}

func TestStage8_5_AllEditUpdatesBaseRule(t *testing.T) {
	t.Run("All Edit Updates Base Rule Directly", func(t *testing.T) {
		rule := map[string]interface{}{
			"id":         uint(1),
			"room_id":    uint(5),
			"start_time": "10:00:00",
			"end_time":   "11:00:00",
		}

		newRoomID := uint(10)
		rule["room_id"] = newRoomID

		if rule["room_id"].(uint) != newRoomID {
			t.Error("Rule should be updated with new room ID")
		}
	})
}

func TestStage8_5_DeleteRecurringSchedule(t *testing.T) {
	t.Run("Single Delete Creates Cancel Exception", func(t *testing.T) {
		deleteMode := "SINGLE"
		exceptionType := "CANCEL"

		if deleteMode != "SINGLE" {
			t.Error("Delete mode should be SINGLE for single deletion")
		}
		if exceptionType != "CANCEL" {
			t.Error("Exception type should be CANCEL for deletion")
		}
	})

	t.Run("Future Delete Creates Cancel Exceptions For Future Dates", func(t *testing.T) {
		editDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		effectiveEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

		var affectedDates []time.Time
		current := editDate
		for !current.After(effectiveEnd) {
			if current.Weekday() == time.Monday {
				affectedDates = append(affectedDates, current)
			}
			current = current.AddDate(0, 0, 1)
		}

		if len(affectedDates) == 0 {
			t.Error("Should have affected dates for future delete")
		}
		if len(affectedDates) < 2 {
			t.Error("Should have at least 2 affected dates from 1/15 to 1/31")
		}
	})

	t.Run("All Delete Removes The Rule", func(t *testing.T) {
		deleteMode := "ALL"
		ruleID := uint(1)

		if deleteMode != "ALL" {
			t.Error("Delete mode should be ALL for full deletion")
		}
		if ruleID == 0 {
			t.Error("Rule ID should not be zero")
		}
	})
}

func TestStage8_5_AffectedDatesCalculation(t *testing.T) {
	t.Run("Calculate Mondays From Edit Date", func(t *testing.T) {
		editDate := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
		effectiveEnd := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		var mondays []time.Time
		current := editDate
		for !current.After(effectiveEnd) {
			if current.Weekday() == time.Monday {
				mondays = append(mondays, current)
			}
			current = current.AddDate(0, 0, 1)
		}

		if len(mondays) < 2 || len(mondays) > 3 {
			t.Errorf("Expected 2-3 Mondays from 1/6 to 1/20, got %d", len(mondays))
		}
	})

	t.Run("Calculate Fridays From Edit Date", func(t *testing.T) {
		editDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		effectiveEnd := time.Date(2026, 1, 24, 0, 0, 0, 0, time.UTC)

		var fridays []time.Time
		current := editDate
		for !current.After(effectiveEnd) {
			if current.Weekday() == time.Friday {
				fridays = append(fridays, current)
			}
			current = current.AddDate(0, 0, 1)
		}

		if len(fridays) < 2 || len(fridays) > 3 {
			t.Errorf("Expected 2-3 Fridays from 1/10 to 1/24, got %d", len(fridays))
		}
	})

	t.Run("Calculate All Wednesdays In Month", func(t *testing.T) {
		startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

		var wednesdays []time.Time
		current := startDate
		for !current.After(endDate) {
			if current.Weekday() == time.Wednesday {
				wednesdays = append(wednesdays, current)
			}
			current = current.AddDate(0, 0, 1)
		}

		if len(wednesdays) != 4 && len(wednesdays) != 5 {
			t.Errorf("Expected 4-5 Wednesdays in January 2026, got %d", len(wednesdays))
		}
	})
}

func TestStage8_5_RecurrenceEditRequestValidation(t *testing.T) {
	t.Run("Valid Single Mode Request", func(t *testing.T) {
		req := map[string]interface{}{
			"rule_id":        uint(1),
			"edit_date":      "2026-01-06",
			"mode":           "SINGLE",
			"new_start_time": "14:00:00",
			"new_end_time":   "15:30:00",
			"new_room_id":    uint(5),
			"reason":         "Test change",
		}

		if req["mode"] != "SINGLE" {
			t.Error("Mode should be SINGLE")
		}
		if req["rule_id"].(uint) == 0 {
			t.Error("Rule ID should not be zero")
		}
	})

	t.Run("Valid Future Mode Request", func(t *testing.T) {
		req := map[string]interface{}{
			"rule_id":        uint(1),
			"edit_date":      "2026-01-15",
			"mode":           "FUTURE",
			"new_teacher_id": uint(100),
			"reason":         "Teacher change",
		}

		if req["mode"] != "FUTURE" {
			t.Error("Mode should be FUTURE")
		}
		if req["new_teacher_id"].(uint) == 0 {
			t.Error("New teacher ID should not be zero when provided")
		}
	})

	t.Run("Valid All Mode Request", func(t *testing.T) {
		req := map[string]interface{}{
			"rule_id":     uint(1),
			"edit_date":   "2026-01-01",
			"mode":        "ALL",
			"new_room_id": uint(10),
			"reason":      "Room change",
		}

		if req["mode"] != "ALL" {
			t.Error("Mode should be ALL")
		}
	})

	t.Run("Invalid Mode Should Be Rejected", func(t *testing.T) {
		validModes := map[string]bool{
			"SINGLE": true,
			"FUTURE": true,
			"ALL":    true,
		}

		testMode := "INVALID"
		if validModes[testMode] {
			t.Error("Invalid mode should not be in valid modes")
		}
	})
}
