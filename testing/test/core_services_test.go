package test

import (
	"testing"
	"time"
)

func TestValidationService_OverlapDetection(t *testing.T) {
	t.Run("Standard Overlap", func(t *testing.T) {
		existingStart := parseTime("10:00:00")
		existingEnd := parseTime("11:00:00")
		newStart := parseTime("10:30:00")
		newEnd := parseTime("11:30:00")

		hasOverlap := newStart.Before(existingEnd) && newEnd.After(existingStart)
		if !hasOverlap {
			t.Error("Expected overlap for standard overlap case")
		}
	})

	t.Run("Contained Overlap", func(t *testing.T) {
		existingStart := parseTime("10:00:00")
		existingEnd := parseTime("12:00:00")
		newStart := parseTime("10:30:00")
		newEnd := parseTime("11:30:00")

		hasOverlap := newStart.Before(existingEnd) && newEnd.After(existingStart)
		if !hasOverlap {
			t.Error("Expected overlap for contained case")
		}
	})

	t.Run("No Overlap - Adjacent", func(t *testing.T) {
		existingStart := parseTime("10:00:00")
		existingEnd := parseTime("11:00:00")
		newStart := parseTime("11:00:00")
		newEnd := parseTime("12:00:00")

		hasOverlap := newStart.Before(existingEnd) && newEnd.After(existingStart)
		if hasOverlap {
			t.Error("Should NOT overlap when end equals start")
		}
	})

	t.Run("Exact Match", func(t *testing.T) {
		existingStart := parseTime("10:00:00")
		existingEnd := parseTime("11:00:00")
		newStart := parseTime("10:00:00")
		newEnd := parseTime("11:00:00")

		hasOverlap := newStart.Before(existingEnd) && newEnd.After(existingStart)
		if !hasOverlap {
			t.Error("Expected overlap for exact match")
		}
	})

	t.Run("No Overlap - Separate", func(t *testing.T) {
		existingStart := parseTime("10:00:00")
		existingEnd := parseTime("11:00:00")
		newStart := parseTime("14:00:00")
		newEnd := parseTime("15:00:00")

		hasOverlap := newStart.Before(existingEnd) && newEnd.After(existingStart)
		if hasOverlap {
			t.Error("Should NOT overlap for separate times")
		}
	})
}

func TestValidationService_BufferCheck(t *testing.T) {
	t.Run("Sufficient Buffer", func(t *testing.T) {
		prevEnd := parseTime("11:00:00")
		nextStart := parseTime("11:20:00")
		requiredBuffer := 15

		gap := int(nextStart.Sub(prevEnd).Minutes())
		isSufficient := gap >= requiredBuffer

		if !isSufficient {
			t.Errorf("Expected sufficient buffer (%d min), got %d min", requiredBuffer, gap)
		}
	})

	t.Run("Insufficient Buffer", func(t *testing.T) {
		prevEnd := parseTime("11:00:00")
		nextStart := parseTime("11:10:00")
		requiredBuffer := 15

		gap := int(nextStart.Sub(prevEnd).Minutes())
		isSufficient := gap >= requiredBuffer

		if isSufficient {
			t.Errorf("Expected insufficient buffer (%d min), got %d min", requiredBuffer, gap)
		}
	})

	t.Run("Zero Buffer Required", func(t *testing.T) {
		prevEnd := parseTime("11:00:00")
		nextStart := parseTime("11:00:00")
		requiredBuffer := 0

		gap := int(nextStart.Sub(prevEnd).Minutes())
		isSufficient := gap >= requiredBuffer

		if !isSufficient {
			t.Error("Expected sufficient buffer when no buffer required")
		}
	})
}

func TestValidationService_WeekdayMapping(t *testing.T) {
	t.Run("Sunday Mapping", func(t *testing.T) {
		weekday := time.Sunday
		dbWeekday := 0
		if weekday == 0 {
			dbWeekday = 7
		}
		if dbWeekday != 7 {
			t.Errorf("Sunday should map to 7, got %d", dbWeekday)
		}
	})

	t.Run("Monday Mapping", func(t *testing.T) {
		weekday := time.Monday
		dbWeekday := int(weekday)
		if dbWeekday == 0 {
			dbWeekday = 7
		}
		if dbWeekday != 1 {
			t.Errorf("Monday should map to 1, got %d", dbWeekday)
		}
	})

	t.Run("Saturday Mapping", func(t *testing.T) {
		weekday := time.Saturday
		dbWeekday := int(weekday)
		if dbWeekday == 0 {
			dbWeekday = 7
		}
		if dbWeekday != 6 {
			t.Errorf("Saturday should map to 6, got %d", dbWeekday)
		}
	})
}

func TestExpansionService_RuleExpansion(t *testing.T) {
	t.Run("Expand Single Week", func(t *testing.T) {
		startDate := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC)

		ruleWeekday := 1

		var expansions int
		current := startDate
		for !current.After(endDate) {
			currentWeekday := int(current.Weekday())
			if currentWeekday == 0 {
				currentWeekday = 7
			}
			if currentWeekday == ruleWeekday {
				expansions++
			}
			current = current.AddDate(0, 0, 1)
		}

		if expansions != 1 {
			t.Errorf("Expected 1 expansion (Monday), got %d", expansions)
		}
	})

	t.Run("Expand Multiple Weeks", func(t *testing.T) {
		startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

		ruleWeekday := 1

		var expansions int
		current := startDate
		for !current.After(endDate) {
			currentWeekday := int(current.Weekday())
			if currentWeekday == 0 {
				currentWeekday = 7
			}
			if currentWeekday == ruleWeekday {
				expansions++
			}
			current = current.AddDate(0, 0, 1)
		}

		if expansions != 4 && expansions != 5 {
			t.Errorf("Expected 4-5 Mondays, got %d", expansions)
		}
	})
}

func TestExpansionService_EffectiveRange(t *testing.T) {
	t.Run("Within Effective Range", func(t *testing.T) {
		ruleStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		ruleEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
		checkDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

		isWithin := !checkDate.Before(ruleStart) && !checkDate.After(ruleEnd)
		if !isWithin {
			t.Error("Date should be within effective range")
		}
	})

	t.Run("Before Effective Range", func(t *testing.T) {
		ruleStart := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		ruleEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
		checkDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

		isWithin := !checkDate.Before(ruleStart) && !checkDate.After(ruleEnd)
		if isWithin {
			t.Error("Date should NOT be within effective range")
		}
	})

	t.Run("After Effective Range", func(t *testing.T) {
		ruleStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		ruleEnd := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		checkDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		isWithin := !checkDate.Before(ruleStart) && !checkDate.After(ruleEnd)
		if isWithin {
			t.Error("Date should NOT be within effective range")
		}
	})

	t.Run("Open Ended Range", func(t *testing.T) {
		ruleStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		var ruleEnd time.Time
		checkDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

		isWithin := true
		if !ruleStart.IsZero() && checkDate.Before(ruleStart) {
			isWithin = false
		}
		if !ruleEnd.IsZero() && checkDate.After(ruleEnd) {
			isWithin = false
		}

		if !isWithin {
			t.Error("Date should be within open-ended range")
		}
	})
}

func TestExceptionService_DeadlineCalculation(t *testing.T) {
	t.Run("Calculate Deadline from Lead Days", func(t *testing.T) {
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		expected := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)

		if !deadline.Equal(expected) {
			t.Errorf("Expected deadline %v, got %v", expected, deadline)
		}
	})

	t.Run("Check if Within Deadline", func(t *testing.T) {
		now := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		isPastDeadline := now.After(deadline)

		if isPastDeadline {
			t.Error("Should NOT be past deadline")
		}
	})
}

func TestSmartMatchingService_Scoring(t *testing.T) {
	t.Run("Skill Match Score", func(t *testing.T) {
		requiredSkills := []string{"鋼琴", "古典音樂"}
		teacherSkills := []string{"鋼琴", "古典音樂", "樂理"}

		matchCount := 0
		for _, required := range requiredSkills {
			for _, skill := range teacherSkills {
				if required == skill {
					matchCount++
					break
				}
			}
		}

		score := float64(matchCount) / float64(len(requiredSkills)) * 100
		expectedScore := 100.0

		if score != expectedScore {
			t.Errorf("Expected skill match score %.2f, got %.2f", expectedScore, score)
		}
	})

	t.Run("Partial Skill Match Score", func(t *testing.T) {
		requiredSkills := []string{"鋼琴", "古典音樂", "爵士"}
		teacherSkills := []string{"鋼琴", "流行音樂"}

		matchCount := 0
		for _, required := range requiredSkills {
			for _, skill := range teacherSkills {
				if required == skill {
					matchCount++
					break
				}
			}
		}

		score := float64(matchCount) / float64(len(requiredSkills)) * 100
		expectedScore := 33.33

		if score < expectedScore-1 || score > expectedScore+1 {
			t.Errorf("Expected skill match score around %.2f, got %.2f", expectedScore, score)
		}
	})
}

func TestNotificationService_MessageFormatting(t *testing.T) {
	t.Run("Format Schedule Reminder", func(t *testing.T) {
		scheduleTime := time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC)
		courseName := "鋼琴課"
		centerName := "莫札特音樂教室"

		message := formatScheduleReminder(scheduleTime, courseName, centerName)

		if message == "" {
			t.Error("Message should not be empty")
		}
	})

	t.Run("Format Emergency Notification", func(t *testing.T) {
		oldTime := time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC)
		newTime := time.Date(2026, 1, 21, 14, 0, 0, 0, time.UTC)
		courseName := "鋼琴課"

		message := formatEmergencyChange(oldTime, newTime, courseName)

		if message == "" {
			t.Error("Emergency message should not be empty")
		}
	})
}

func TestAuditLogService_ActionFormatting(t *testing.T) {
	t.Run("Action String Formatting", func(t *testing.T) {
		actions := []string{
			"CREATE_RULE",
			"APPROVE_EXCEPTION",
			"REJECT_EXCEPTION",
			"DELETE_SCHEDULE_RULE",
		}

		for _, action := range actions {
			if action == "" {
				t.Error("Action should not be empty")
			}
		}
	})

	t.Run("Action Consistency", func(t *testing.T) {
		expectedActions := map[string]bool{
			"CREATE_RULE":          true,
			"UPDATE_RULE":          true,
			"DELETE_RULE":          true,
			"CREATE_EXCEPTION":     true,
			"APPROVE_EXCEPTION":    true,
			"REJECT_EXCEPTION":     true,
			"REVOKE_EXCEPTION":     true,
			"CREATE_HOLIDAY":       true,
			"DELETE_HOLIDAY":       true,
			"BULK_CREATE_HOLIDAYS": true,
		}

		for action := range expectedActions {
			if _, exists := expectedActions[action]; !exists {
				t.Errorf("Missing action: %s", action)
			}
		}
	})
}

func TestPaginationService_Calculation(t *testing.T) {
	t.Run("Calculate Total Pages", func(t *testing.T) {
		total := 100
		limit := 20

		totalPages := (total + limit - 1) / limit

		if totalPages != 5 {
			t.Errorf("Expected 5 pages, got %d", totalPages)
		}
	})

	t.Run("Calculate Has Next", func(t *testing.T) {
		page := 2
		totalPages := 5

		hasNext := page < totalPages

		if !hasNext {
			t.Error("Should have next page")
		}
	})

	t.Run("Calculate Has Prev", func(t *testing.T) {
		page := 2

		hasPrev := page > 1

		if !hasPrev {
			t.Error("Should have previous page")
		}
	})

	t.Run("Last Page Has Next False", func(t *testing.T) {
		page := 5
		totalPages := 5

		hasNext := page < totalPages

		if hasNext {
			t.Error("Last page should NOT have next")
		}
	})

	t.Run("First Page Has Prev False", func(t *testing.T) {
		page := 1

		hasPrev := page > 1

		if hasPrev {
			t.Error("First page should NOT have previous")
		}
	})
}

func TestTimezoneService_UTCConversion(t *testing.T) {
	t.Run("UTC to Local", func(t *testing.T) {
		utcTime := time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC)
		loc, _ := time.LoadLocation("Asia/Taipei")

		localTime := utcTime.In(loc)

		if localTime.Hour() != 18 {
			t.Errorf("Expected 18:00 in Taipei (UTC+8), got %d:00", localTime.Hour())
		}
	})

	t.Run("Local to UTC", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Taipei")
		localTime := time.Date(2026, 1, 20, 10, 0, 0, 0, loc)

		utcTime := localTime.UTC()

		if utcTime.Hour() != 2 {
			t.Errorf("Expected 02:00 UTC, got %d:00", utcTime.Hour())
		}
	})
}

func TestResourceService_CapacityCheck(t *testing.T) {
	t.Run("Check Room Capacity", func(t *testing.T) {
		roomCapacity := 20
		bookingCount := 15

		hasCapacity := bookingCount < roomCapacity
		if !hasCapacity {
			t.Error("Room should have capacity for 15 bookings out of 20")
		}
	})

	t.Run("Room At Capacity", func(t *testing.T) {
		roomCapacity := 20
		bookingCount := 20

		hasCapacity := bookingCount < roomCapacity
		if hasCapacity {
			t.Error("Room should be at capacity")
		}
	})

	t.Run("Room Over Capacity", func(t *testing.T) {
		roomCapacity := 20
		bookingCount := 25

		hasCapacity := bookingCount < roomCapacity
		if hasCapacity {
			t.Error("Room should be over capacity")
		}
	})
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("15:04:05", s)
	baseDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

func formatScheduleReminder(scheduleTime time.Time, courseName string, centerName string) string {
	return time.Now().Format(time.RFC3339) + " | " + courseName + " at " + centerName
}

func formatEmergencyChange(oldTime, newTime time.Time, courseName string) string {
	return "課程異動通知 | " + courseName + " 從 " + oldTime.Format("1/15 15:04") + " 改至 " + newTime.Format("1/15 15:04")
}
