package services

import (
	"context"
	"testing"
	"time"

	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

// TestScheduleService_HelperFunctions 測試輔助函數
func TestScheduleService_HelperFunctions(t *testing.T) {
	t.Run("splitTime_ValidFormat", func(t *testing.T) {
		result := splitTime("14:30")
		if len(result) != 2 {
			t.Fatalf("Expected 2 parts, got %d", len(result))
		}
		if result[0] != 14 {
			t.Errorf("Expected hour 14, got %d", result[0])
		}
		if result[1] != 30 {
			t.Errorf("Expected minute 30, got %d", result[1])
		}
	})

	t.Run("splitTime_EmptyString", func(t *testing.T) {
		result := splitTime("")
		if result != nil {
			t.Error("Expected nil for empty string")
		}
	})

	t.Run("splitTime_SinglePart", func(t *testing.T) {
		result := splitTime("14")
		if result != nil {
			t.Error("Expected nil for single part")
		}
	})

	t.Run("parseInt_ValidNumber", func(t *testing.T) {
		result := parseInt("123")
		if result != 123 {
			t.Errorf("Expected 123, got %d", result)
		}
	})

	t.Run("parseInt_EmptyString", func(t *testing.T) {
		result := parseInt("")
		if result != 0 {
			t.Errorf("Expected 0 for empty string, got %d", result)
		}
	})

	t.Run("parseInt_ComplexNumber", func(t *testing.T) {
		// parseInt 解析數字字串，不處理特殊字元
		// "09:00" 會被解析為 900（跳過冒號後解析數字）
		result := parseInt("09:00")
		if result != 900 {
			t.Errorf("Expected 900 for '09:00' (parsing numeric parts), got %d", result)
		}
	})

	t.Run("intToString_Zero", func(t *testing.T) {
		result := intToString(0)
		if result != "0" {
			t.Errorf("Expected '0', got '%s'", result)
		}
	})

	t.Run("intToString_MultiDigit", func(t *testing.T) {
		result := intToString(12345)
		if result != "12345" {
			t.Errorf("Expected '12345', got '%s'", result)
		}
	})
}

// TestScheduleService_DeadlineCalculation 測試截止日計算
func TestScheduleService_DeadlineCalculation(t *testing.T) {
	t.Run("CalculateDeadline", func(t *testing.T) {
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		expected := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)

		if !deadline.Equal(expected) {
			t.Errorf("Expected deadline %v, got %v", expected, deadline)
		}
	})

	t.Run("DaysRemaining_Positive", func(t *testing.T) {
		// 使用未來日期測試
		now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
		exceptionDate := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		daysRemaining := int(deadline.Sub(now).Hours() / 24)

		if daysRemaining < 0 {
			t.Errorf("Expected positive days remaining, got %d", daysRemaining)
		}
	})

	t.Run("DaysRemaining_Negative", func(t *testing.T) {
		now := time.Date(2026, 1, 20, 12, 0, 0, 0, time.UTC)
		exceptionDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		leadDays := 14

		deadline := exceptionDate.AddDate(0, 0, -leadDays)
		daysRemaining := int(deadline.Sub(now).Hours() / 24)

		if daysRemaining >= 0 {
			t.Errorf("Expected negative days remaining (past deadline), got %d", daysRemaining)
		}
	})
}

// TestScheduleService_UpdateModeConstants 測試更新模式常量
func TestScheduleService_UpdateModeConstants(t *testing.T) {
	t.Run("UpdateModeSingle_Value", func(t *testing.T) {
		if UpdateModeSingle != "SINGLE" {
			t.Errorf("Expected 'SINGLE', got '%s'", UpdateModeSingle)
		}
	})

	t.Run("UpdateModeFuture_Value", func(t *testing.T) {
		if UpdateModeFuture != "FUTURE" {
			t.Errorf("Expected 'FUTURE', got '%s'", UpdateModeFuture)
		}
	})

	t.Run("UpdateModeAll_Value", func(t *testing.T) {
		if UpdateModeAll != "ALL" {
			t.Errorf("Expected 'ALL', got '%s'", UpdateModeAll)
		}
	})
}

// TestScheduleRuleModel_EffectiveRange 測試規則有效範圍
func TestScheduleRuleModel_EffectiveRange(t *testing.T) {
	t.Run("WithinEffectiveRange", func(t *testing.T) {
		ruleStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		ruleEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
		checkDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

		isWithin := !checkDate.Before(ruleStart) && !checkDate.After(ruleEnd)
		if !isWithin {
			t.Error("Date should be within effective range")
		}
	})

	t.Run("BeforeEffectiveRange", func(t *testing.T) {
		ruleStart := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		ruleEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
		checkDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

		isWithin := !checkDate.Before(ruleStart) && !checkDate.After(ruleEnd)
		if isWithin {
			t.Error("Date should NOT be within effective range")
		}
	})

	t.Run("AfterEffectiveRange", func(t *testing.T) {
		ruleStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		ruleEnd := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		checkDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		isWithin := !checkDate.Before(ruleStart) && !checkDate.After(ruleEnd)
		if isWithin {
			t.Error("Date should NOT be within effective range")
		}
	})

	t.Run("OpenEndedRange", func(t *testing.T) {
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

// TestExpandRulesBatchFetchExceptions 測試 ExpandRules 的批次例外查詢優化
func TestExpandRulesBatchFetchExceptions(t *testing.T) {
	appInstance := setupTestApp(t)
	if appInstance == nil {
		return
	}

	ctx := context.Background()
	expansionSvc := NewScheduleExpansionService(appInstance)

	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	rules, err := ruleRepo.ListByCenterID(ctx, 1)
	if err != nil || len(rules) == 0 {
		t.Skip("No schedule rules available for testing")
		return
	}

	startDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)

	schedules := expansionSvc.ExpandRules(ctx, rules, startDate, endDate, 1)

	if len(schedules) == 0 {
		t.Log("No schedules generated in date range")
	}

	for i, schedule := range schedules {
		if schedule.RuleID == 0 {
			t.Errorf("Schedule %d: RuleID should not be zero", i)
		}
		if schedule.Date.IsZero() {
			t.Errorf("Schedule %d: Date should not be zero", i)
		}
		if schedule.StartTime == "" {
			t.Errorf("Schedule %d: StartTime should not be empty", i)
		}
		if schedule.EndTime == "" {
			t.Errorf("Schedule %d: EndTime should not be empty", i)
		}
	}

	t.Logf("Generated %d schedules using batch exception fetch", len(schedules))
}

// TestScheduleExceptionRepositoryBatchFetch 測試 Repository 的批次查詢方法
func TestScheduleExceptionRepositoryBatchFetch(t *testing.T) {
	appInstance := setupTestApp(t)
	if appInstance == nil {
		return
	}

	ctx := context.Background()
	exceptionRepo := repositories.NewScheduleExceptionRepository(appInstance)

	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	rules, err := ruleRepo.ListByCenterID(ctx, 1)
	if err != nil || len(rules) == 0 {
		t.Skip("No schedule rules available for testing")
		return
	}

	ruleIDs := make([]uint, 0, len(rules))
	for _, rule := range rules {
		ruleIDs = append(ruleIDs, rule.ID)
	}

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	exceptionsMap, err := exceptionRepo.GetByRuleIDsAndDateRange(ctx, ruleIDs, startDate, endDate)
	if err != nil {
		t.Fatalf("Batch query failed: %v", err)
	}

	if exceptionsMap == nil {
		t.Fatal("Exceptions map should not be nil")
	}

	for _, ruleID := range ruleIDs {
		if _, ok := exceptionsMap[ruleID]; !ok {
			t.Logf("Rule %d has no exceptions in date range (this is expected)", ruleID)
			if exceptionsMap[ruleID] == nil {
				exceptionsMap[ruleID] = make(map[string][]models.ScheduleException)
			}
		}
	}

	totalExceptions := 0
	for _, dateMap := range exceptionsMap {
		for _, exceptions := range dateMap {
			totalExceptions += len(exceptions)
		}
	}

	t.Logf("Found %d total exceptions for %d rules in date range", totalExceptions, len(ruleIDs))
}

// TestExpandRulesWithExceptions 測試 ExpandRules 處理例外狀態的正確性
func TestExpandRulesWithExceptions(t *testing.T) {
	appInstance := setupTestApp(t)
	if appInstance == nil {
		return
	}

	ctx := context.Background()
	expansionSvc := NewScheduleExpansionService(appInstance)
	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)

	rules, err := ruleRepo.ListByCenterID(ctx, 1)
	if err != nil || len(rules) == 0 {
		t.Skip("No schedule rules available for testing")
		return
	}

	startDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)

	schedules := expansionSvc.ExpandRules(ctx, rules, startDate, endDate, 1)

	for _, schedule := range schedules {
		if (schedule.ExceptionInfo != nil) != schedule.HasException {
			t.Errorf("Schedule %d: HasException (%v) should match ExceptionInfo != nil (%v)",
				schedule.RuleID, schedule.HasException, schedule.ExceptionInfo != nil)
		}
	}

	t.Logf("Tested %d schedules for exception handling", len(schedules))
}
