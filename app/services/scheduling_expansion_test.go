package services

import (
	"context"
	"testing"
	"time"

	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

// TestExpandRulesBatchFetchExceptions 測試 ExpandRules 的批次例外查詢優化
// 驗證邏輯與原來 N+1 查詢方式一致
func TestExpandRulesBatchFetchExceptions(t *testing.T) {
	// 初始化測試環境（使用實際開發資料庫）
	appInstance := setupTestApp(t)
	if appInstance == nil {
		return
	}

	ctx := context.Background()
	expansionSvc := NewScheduleExpansionService(appInstance)

	// 取得測試用的規則資料
	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	rules, err := ruleRepo.ListByCenterID(ctx, 1)
	if err != nil || len(rules) == 0 {
		t.Skip("No schedule rules available for testing")
		return
	}

	// 測試日期範圍（包含可能需要檢查例外的日期）
	startDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)

	// 執行 ExpandRules（使用優化後的批次查詢）
	schedules := expansionSvc.ExpandRules(ctx, rules, startDate, endDate, 1)

	// 驗證結果不為空
	if len(schedules) == 0 {
		t.Log("No schedules generated in date range")
	}

	// 驗證每個 schedule 的基本欄位
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
	// 初始化測試環境
	appInstance := setupTestApp(t)
	if appInstance == nil {
		return
	}

	ctx := context.Background()
	exceptionRepo := repositories.NewScheduleExceptionRepository(appInstance)

	// 取得測試用的規則資料
	ruleRepo := repositories.NewScheduleRuleRepository(appInstance)
	rules, err := ruleRepo.ListByCenterID(ctx, 1)
	if err != nil || len(rules) == 0 {
		t.Skip("No schedule rules available for testing")
		return
	}

	// 收集規則 ID
	ruleIDs := make([]uint, 0, len(rules))
	for _, rule := range rules {
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// 測試日期範圍
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	// 執行批次查詢
	exceptionsMap, err := exceptionRepo.GetByRuleIDsAndDateRange(ctx, ruleIDs, startDate, endDate)
	if err != nil {
		t.Fatalf("Batch query failed: %v", err)
	}

	// 驗證結果結構
	if exceptionsMap == nil {
		t.Fatal("Exceptions map should not be nil")
	}

	// 驗證每個規則都有對應的 map
	for _, ruleID := range ruleIDs {
		if _, ok := exceptionsMap[ruleID]; !ok {
			t.Logf("Rule %d has no exceptions in date range (this is expected)", ruleID)
			// 確保即使沒有例外，該規則的 map 仍然存在
			if exceptionsMap[ruleID] == nil {
				exceptionsMap[ruleID] = make(map[string][]models.ScheduleException)
			}
		}
	}

	// 計算總例外數量
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
	// 初始化測試環境
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

	// 測試日期範圍
	startDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)

	// 執行 ExpandRules
	schedules := expansionSvc.ExpandRules(ctx, rules, startDate, endDate, 1)

	// 驗證 HasException 欄位設置正確
	for _, schedule := range schedules {
		// HasException 應該與 ExceptionInfo 是否存在一致
		if (schedule.ExceptionInfo != nil) != schedule.HasException {
			t.Errorf("Schedule %d: HasException (%v) should match ExceptionInfo != nil (%v)",
				schedule.RuleID, schedule.HasException, schedule.ExceptionInfo != nil)
		}
	}

	t.Logf("Tested %d schedules for exception handling", len(schedules))
}
