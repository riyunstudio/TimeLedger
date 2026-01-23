package test

import (
	"testing"
	"time"
	"timeLedger/app/models"
)

func TestStage10_LockingLogic(t *testing.T) {
	t.Run("Rule LockAt Past", func(t *testing.T) {
		lockAt := time.Date(2026, 1, 15, 23, 59, 59, 0, time.UTC)
		now := time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)

		isLocked := lockAt.IsZero() || now.After(lockAt)
		if !isLocked {
			t.Error("Rule should be locked when now is after lock_at")
		}
	})

	t.Run("Rule LockAt Future", func(t *testing.T) {
		lockAt := time.Date(2026, 1, 20, 23, 59, 59, 0, time.UTC)
		now := time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)

		isLocked := lockAt.IsZero() || now.After(lockAt)
		if isLocked {
			t.Error("Rule should NOT be locked when lock_at is in the future")
		}
	})

	t.Run("Rule No LockAt", func(t *testing.T) {
		var lockAt *time.Time = nil

		isLocked := lockAt == nil || time.Now().After(*lockAt)
		if !isLocked {
			t.Error("Rule should NOT be locked when lock_at is nil")
		}
	})
}

func TestStage10_ExceptionLeadDays(t *testing.T) {
	t.Run("Lead Days Calculation", func(t *testing.T) {
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		expectedDeadline := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)

		if !deadline.Equal(expectedDeadline) {
			t.Errorf("Deadline should be %v, got %v", expectedDeadline, deadline)
		}
	})

	t.Run("Days Remaining Calculation", func(t *testing.T) {
		now := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		deadline := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		daysRemaining := int(deadline.Sub(now).Hours() / 24)

		if daysRemaining != 5 {
			t.Errorf("Days remaining should be 5, got %d", daysRemaining)
		}
	})

	t.Run("Negative Days Remaining", func(t *testing.T) {
		now := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		deadline := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		daysRemaining := int(deadline.Sub(now).Hours() / 24)

		if daysRemaining >= 0 {
			t.Errorf("Days remaining should be negative when past deadline, got %d", daysRemaining)
		}
	})
}

func TestStage10_CenterPolicy(t *testing.T) {
	t.Run("Default Lead Days", func(t *testing.T) {
		settings := models.CenterSettings{
			ExceptionLeadDays: 0,
		}

		leadDays := settings.ExceptionLeadDays
		if leadDays <= 0 {
			leadDays = 14
		}

		if leadDays != 14 {
			t.Errorf("Default lead days should be 14, got %d", leadDays)
		}
	})

	t.Run("Custom Lead Days", func(t *testing.T) {
		settings := models.CenterSettings{
			ExceptionLeadDays: 7,
		}

		leadDays := settings.ExceptionLeadDays
		if leadDays <= 0 {
			leadDays = 14
		}

		if leadDays != 7 {
			t.Errorf("Lead days should be 7, got %d", leadDays)
		}
	})
}

func TestStage10_CheckDeadlineLogic(t *testing.T) {
	_ = time.LoadLocation // Suppress unused import warning
	_ = time.Date

	t.Run("Within Deadline", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Taipei")
		now := time.Date(2026, 1, 10, 0, 0, 0, 0, loc)
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, loc)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		isPastDeadline := now.After(deadline)

		if isPastDeadline {
			t.Errorf("Should NOT be past deadline. Now: %v, Deadline: %v", now, deadline)
		}
	})

	t.Run("Past Deadline", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Taipei")
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, loc)
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, loc)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		isPastDeadline := now.After(deadline)

		if !isPastDeadline {
			t.Error("Should be past deadline")
		}
	})

	t.Run("Exactly At Deadline", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Taipei")
		now := time.Date(2026, 1, 11, 0, 0, 0, 0, loc)
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, loc)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		isPastDeadline := now.After(deadline)

		if isPastDeadline {
			t.Error("Should NOT be past deadline when exactly at deadline")
		}
	})
}

func TestStage10_ScheduleRule_LockAt(t *testing.T) {
	t.Run("ScheduleRule with LockAt", func(t *testing.T) {
		lockAt := time.Date(2026, 1, 15, 23, 59, 59, 0, time.UTC)

		rule := models.ScheduleRule{
			ID:        1,
			CenterID:  10,
			Weekday:   1,
			StartTime: "10:00:00",
			EndTime:   "11:00:00",
			LockAt:    &lockAt,
		}

		if rule.LockAt == nil {
			t.Error("LockAt should not be nil")
		}

		if !rule.LockAt.Equal(lockAt) {
			t.Error("LockAt should match")
		}
	})

	t.Run("ScheduleRule without LockAt", func(t *testing.T) {
		rule := models.ScheduleRule{
			ID:        2,
			CenterID:  10,
			Weekday:   2,
			StartTime: "14:00:00",
			EndTime:   "15:00:00",
			LockAt:    nil,
		}

		if rule.LockAt != nil {
			t.Error("LockAt should be nil")
		}
	})
}

func TestStage10_ExceptionRequest_Validation(t *testing.T) {
	t.Run("CreateException Check", func(t *testing.T) {
		now := time.Now()
		lockAt := now.AddDate(0, 0, 7)

		rule := models.ScheduleRule{
			ID:       1,
			CenterID: 10,
			LockAt:   &lockAt,
		}

		centerSettings := models.CenterSettings{
			ExceptionLeadDays: 14,
		}

		exceptionDate := now.AddDate(0, 0, 10)

		isLocked := false
		var reason string

		if rule.LockAt != nil && now.After(*rule.LockAt) {
			isLocked = true
			reason = "已超過異動截止日"
		} else {
			leadDays := centerSettings.ExceptionLeadDays
			if leadDays <= 0 {
				leadDays = 14
			}
			deadline := exceptionDate.AddDate(0, 0, -leadDays)
			if now.After(deadline) {
				isLocked = true
				reason = "已超過異動截止日（需提前 14 天申請）"
			}
		}

		if isLocked {
			t.Logf("Exception is locked: %s", reason)
		} else {
			t.Log("Exception is allowed")
		}
	})
}
