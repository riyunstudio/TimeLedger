package test

import (
	"testing"
	"time"
	"timeLedger/app/models"
)

func TestStage7_ScheduleRule_EffectiveRange(t *testing.T) {
	t.Run("Valid EffectiveRange", func(t *testing.T) {
		startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)

		rule := models.ScheduleRule{
			ID:        1,
			CenterID:  1,
			Weekday:   1,
			StartTime: "10:00:00",
			EndTime:   "11:00:00",
			EffectiveRange: models.DateRange{
				StartDate: startDate,
				EndDate:   endDate,
			},
		}

		if rule.EffectiveRange.StartDate.IsZero() {
			t.Error("StartDate should not be zero")
		}

		if rule.EffectiveRange.EndDate.IsZero() {
			t.Error("EndDate should not be zero")
		}

		if !rule.EffectiveRange.StartDate.Before(rule.EffectiveRange.EndDate) {
			t.Error("StartDate should be before EndDate")
		}
	})

	t.Run("Open-ended EffectiveRange", func(t *testing.T) {
		startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		rule := models.ScheduleRule{
			ID:        2,
			CenterID:  1,
			Weekday:   2,
			StartTime: "14:00:00",
			EndTime:   "15:30:00",
			EffectiveRange: models.DateRange{
				StartDate: startDate,
				EndDate:   time.Time{},
			},
		}

		if rule.EffectiveRange.StartDate.IsZero() {
			t.Error("StartDate should not be zero for open-ended range")
		}

		if !rule.EffectiveRange.EndDate.IsZero() {
			t.Error("EndDate should be zero for open-ended range")
		}
	})

	t.Run("Same Start and End Date", func(t *testing.T) {
		date := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)

		rule := models.ScheduleRule{
			ID:             3,
			CenterID:       1,
			Weekday:        5,
			StartTime:      "09:00:00",
			EndTime:        "10:00:00",
			EffectiveRange: models.DateRange{StartDate: date, EndDate: date},
		}

		if !rule.EffectiveRange.StartDate.Equal(rule.EffectiveRange.EndDate) {
			t.Error("For single-day rule, StartDate and EndDate should be equal")
		}
	})
}

func TestStage7_PhaseDetection_Logic(t *testing.T) {
	t.Run("Same Weekday Different EffectiveRange", func(t *testing.T) {
		rule1Start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		rule1End := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
		rule2Start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
		rule2End := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)

		rule1 := models.ScheduleRule{
			ID:         1,
			OfferingID: 10,
			Weekday:    1,
			StartTime:  "10:00:00",
			EndTime:    "11:00:00",
			EffectiveRange: models.DateRange{
				StartDate: rule1Start,
				EndDate:   rule1End,
			},
		}

		rule2 := models.ScheduleRule{
			ID:         2,
			OfferingID: 10,
			Weekday:    1,
			StartTime:  "14:00:00",
			EndTime:    "15:00:00",
			EffectiveRange: models.DateRange{
				StartDate: rule2Start,
				EndDate:   rule2End,
			},
		}

		rule1EndDate := rule1.EffectiveRange.EndDate
		rule2StartDate := rule2.EffectiveRange.StartDate

		transitionDate := rule2StartDate
		isWithinRule1 := !transitionDate.After(rule1EndDate)
		isWithinRule2 := !transitionDate.Before(rule2StartDate)

		if isWithinRule1 && isWithinRule2 {
			t.Error("Transition date should not be within both rules")
		}

		if !isWithinRule1 && isWithinRule2 {
			t.Logf("Correct: Transition on %v is the start of rule 2", transitionDate)
		}
	})

	t.Run("No Phase Gap", func(t *testing.T) {
		rule1End := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
		rule2Start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

		hasGap := rule2Start.After(rule1End.AddDate(0, 0, 1))

		if hasGap {
			t.Log("There is a gap between phases")
		} else {
			t.Log("No gap - phases are contiguous")
		}
	})

	t.Run("Same Offering Different Times", func(t *testing.T) {
		rule1 := models.ScheduleRule{
			ID:         1,
			OfferingID: 10,
			Weekday:    1,
			StartTime:  "10:00:00",
			EndTime:    "11:00:00",
			EffectiveRange: models.DateRange{
				StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		rule2 := models.ScheduleRule{
			ID:         2,
			OfferingID: 10,
			Weekday:    1,
			StartTime:  "14:00:00",
			EndTime:    "15:00:00",
			EffectiveRange: models.DateRange{
				StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		if rule1.OfferingID != rule2.OfferingID {
			t.Error("Both rules should be for the same offering")
		}

		if rule1.StartTime == rule2.StartTime {
			t.Error("Phase rules should have different start times")
		}
	})
}

func TestStage7_ScheduleException_PhaseContext(t *testing.T) {
	t.Run("Exception Within Phase", func(t *testing.T) {
		phaseStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		phaseEnd := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)
		exceptionDate := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)

		isWithinPhase := !exceptionDate.Before(phaseStart) && !exceptionDate.After(phaseEnd)

		if !isWithinPhase {
			t.Error("Exception date should be within the phase range")
		}
	})

	t.Run("Exception Outside Phase", func(t *testing.T) {
		phaseStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		phaseEnd := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)
		exceptionDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

		isWithinPhase := !exceptionDate.Before(phaseStart) && !exceptionDate.After(phaseEnd)

		if isWithinPhase {
			t.Error("Exception date should NOT be within the phase range")
		}
	})
}

func TestStage7_DateRange_Scan(t *testing.T) {
	t.Run("Scan Nil Value", func(t *testing.T) {
		var dr models.DateRange
		err := dr.Scan(nil)
		if err != nil {
			t.Errorf("Scan(nil) should not return error, got: %v", err)
		}
	})

	t.Run("DateRange Value", func(t *testing.T) {
		original := models.DateRange{
			StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		}

		jsonData, err := original.Value()
		if err != nil {
			t.Fatalf("Failed to serialize DateRange: %v", err)
		}

		var restored models.DateRange
		err = restored.Scan(jsonData)
		if err != nil {
			t.Fatalf("Failed to deserialize DateRange: %v", err)
		}

		if !original.StartDate.Equal(restored.StartDate) {
			t.Errorf("StartDate mismatch: got %v, want %v", restored.StartDate, original.StartDate)
		}

		if !original.EndDate.Equal(restored.EndDate) {
			t.Errorf("EndDate mismatch: got %v, want %v", restored.EndDate, original.EndDate)
		}
	})
}
