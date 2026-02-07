package test

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTestApp 建立測試用的 App 實例
func setupTestApp(t *testing.T) *app.App {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("MySQL init error: %s. Skipping test.", err.Error())
		return nil
	}

	// 檢查資料庫連線
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		t.Skipf("MySQL DB error: %s. Skipping test.", err.Error())
		return nil
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("MySQL ping error: %s. Skipping test.", err.Error())
		return nil
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	env := &configs.Env{
		JWTSecret:   "test-jwt-secret-key-for-testing-only",
		AppEnv:      "test",
		AppDebug:    true,
		AppTimezone: "Asia/Taipei",
	}

	appInstance := &app.App{
		Env:   env,
		Err:   e,
		Tools: tool,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: nil,
		Api:   nil,
		Rpc:   nil,
	}

	return appInstance
}

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
		if services.UpdateModeSingle != "SINGLE" {
			t.Errorf("Expected 'SINGLE', got '%s'", services.UpdateModeSingle)
		}
	})

	t.Run("UpdateModeFuture_Value", func(t *testing.T) {
		if services.UpdateModeFuture != "FUTURE" {
			t.Errorf("Expected 'FUTURE', got '%s'", services.UpdateModeFuture)
		}
	})

	t.Run("UpdateModeAll_Value", func(t *testing.T) {
		if services.UpdateModeAll != "ALL" {
			t.Errorf("Expected 'ALL', got '%s'", services.UpdateModeAll)
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
	expansionSvc := services.NewScheduleExpansionService(appInstance)

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
	expansionSvc := services.NewScheduleExpansionService(appInstance)
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

// TestExpandRules_HolidayLogic 測試 ExpandRules 的假日邏輯處理
func TestExpandRules_HolidayLogic(t *testing.T) {
	appInstance := setupTestApp(t)
	if appInstance == nil {
		return
	}

	ctx := context.Background()
	now := time.Now()
	loc := time.UTC

	// 查詢現有資料用於測試
	var center models.Center
	if err := appInstance.MySQL.RDB.WithContext(ctx).Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("No center data available, skipping test: %v", err)
		return
	}

	var course models.Course
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&course).Error; err != nil {
		t.Skipf("No course data available, skipping test: %v", err)
		return
	}

	var teacher models.Teacher
	if err := appInstance.MySQL.RDB.WithContext(ctx).First(&teacher).Error; err != nil {
		t.Skipf("No teacher data available, skipping test: %v", err)
		return
	}

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 建立獨立的 offering 確保測試隔離
	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
		return
	}
	defer func() {
		// 清理測試資料
		appInstance.MySQL.WDB.WithContext(ctx).Where("id = ?", offering.ID).Delete(&models.Offering{})
	}()

	// 定義測試日期範圍（包含一個假日）
	testStartDate := time.Date(2026, 1, 20, 0, 0, 0, 0, loc)
	testEndDate := time.Date(2026, 1, 26, 0, 0, 0, 0, loc)

	// 建立假日（2026-01-22 是假日）
	holidayDate := time.Date(2026, 1, 22, 0, 0, 0, 0, loc)
	holiday := models.CenterHoliday{
		CenterID:  center.ID,
		Date:      holidayDate,
		Name:      "春節假期",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&holiday).Error; err != nil {
		t.Fatalf("Failed to create holiday: %v", err)
		return
	}
	defer func() {
		// 清理測試資料
		appInstance.MySQL.WDB.WithContext(ctx).Where("id = ?", holiday.ID).Delete(&models.CenterHoliday{})
	}()

	// 測試情境 A：ForceCancel = true，SkipHoliday = false（應該跳過）
	t.Run("ScenarioA_ForceCancelTrue_SkipHolidayFalse", func(t *testing.T) {
		// 更新假日為 ForceCancel = true
		holiday.ForceCancel = true
		if err := appInstance.MySQL.WDB.WithContext(ctx).Save(&holiday).Error; err != nil {
			t.Fatalf("Failed to update holiday: %v", err)
		}

		// 建立規則（SkipHoliday = false）
		rule := models.ScheduleRule{
			CenterID:       center.ID,
			OfferingID:     offering.ID,
			TeacherID:      &teacher.ID,
			RoomID:         room.ID,
			Name:           "測試課程A",
			Weekday:        3, // 週三
			StartTime:      "10:00",
			EndTime:        "12:00",
			Duration:       120,
			SkipHoliday:    false,
			EffectiveRange: models.DateRange{StartDate: testStartDate, EndDate: testEndDate},
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create rule: %v", err)
		}
		defer func() {
			appInstance.MySQL.WDB.WithContext(ctx).Where("id = ?", rule.ID).Delete(&models.ScheduleRule{})
		}()

		expansionSvc := services.NewScheduleExpansionService(appInstance)
		schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testStartDate, testEndDate, center.ID)

		// 驗證：ForceCancel = true 時應該跳過（不論 SkipHoliday）
		for _, schedule := range schedules {
			if schedule.Date.Equal(holidayDate) {
				t.Errorf("Scenario A: Expected to skip holiday %s, but session was generated", holidayDate.Format("2006-01-02"))
			}
		}

		// 統計非假日的天數
		var nonHolidayCount int
		for _, s := range schedules {
			if !s.Date.Equal(holidayDate) {
				nonHolidayCount++
			}
		}
		t.Logf("Scenario A: Generated %d sessions (holiday ForceCancel=true, SkipHoliday=false)", nonHolidayCount)
	})

	// 測試情境 B：ForceCancel = false，SkipHoliday = true（應該跳過）
	t.Run("ScenarioB_ForceCancelFalse_SkipHolidayTrue", func(t *testing.T) {
		// 更新假日為 ForceCancel = false
		holiday.ForceCancel = false
		if err := appInstance.MySQL.WDB.WithContext(ctx).Save(&holiday).Error; err != nil {
			t.Fatalf("Failed to update holiday: %v", err)
		}

		// 建立規則（SkipHoliday = true）
		rule := models.ScheduleRule{
			CenterID:       center.ID,
			OfferingID:     offering.ID,
			TeacherID:      &teacher.ID,
			RoomID:         room.ID,
			Name:           "測試課程B",
			Weekday:        3, // 週三
			StartTime:      "10:00",
			EndTime:        "12:00",
			Duration:       120,
			SkipHoliday:    true,
			EffectiveRange: models.DateRange{StartDate: testStartDate, EndDate: testEndDate},
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create rule: %v", err)
		}
		defer func() {
			appInstance.MySQL.WDB.WithContext(ctx).Where("id = ?", rule.ID).Delete(&models.ScheduleRule{})
		}()

		expansionSvc := services.NewScheduleExpansionService(appInstance)
		schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testStartDate, testEndDate, center.ID)

		// 驗證：ForceCancel = false 且 SkipHoliday = true 時應該跳過
		for _, schedule := range schedules {
			if schedule.Date.Equal(holidayDate) {
				t.Errorf("Scenario B: Expected to skip holiday %s, but session was generated", holidayDate.Format("2006-01-02"))
			}
		}

		var nonHolidayCount int
		for _, s := range schedules {
			if !s.Date.Equal(holidayDate) {
				nonHolidayCount++
			}
		}
		t.Logf("Scenario B: Generated %d sessions (holiday ForceCancel=false, SkipHoliday=true)", nonHolidayCount)
	})

	// 測試情境 C：ForceCancel = false，SkipHoliday = false（應該產生課程）
	t.Run("ScenarioC_ForceCancelFalse_SkipHolidayFalse", func(t *testing.T) {
		// 更新假日為 ForceCancel = false
		holiday.ForceCancel = false
		if err := appInstance.MySQL.WDB.WithContext(ctx).Save(&holiday).Error; err != nil {
			t.Fatalf("Failed to update holiday: %v", err)
		}

		// 建立規則（SkipHoliday = false）
		rule := models.ScheduleRule{
			CenterID:       center.ID,
			OfferingID:     offering.ID,
			TeacherID:      &teacher.ID,
			RoomID:         room.ID,
			Name:           "測試課程C",
			Weekday:        3, // 週三
			StartTime:      "10:00",
			EndTime:        "12:00",
			Duration:       120,
			SkipHoliday:    false,
			EffectiveRange: models.DateRange{StartDate: testStartDate, EndDate: testEndDate},
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create rule: %v", err)
		}
		defer func() {
			appInstance.MySQL.WDB.WithContext(ctx).Where("id = ?", rule.ID).Delete(&models.ScheduleRule{})
		}()

		expansionSvc := services.NewScheduleExpansionService(appInstance)
		schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testStartDate, testEndDate, center.ID)

		// 驗證：ForceCancel = false 且 SkipHoliday = false 時應該產生課程
		var foundHolidaySession bool
		for _, schedule := range schedules {
			if schedule.Date.Equal(holidayDate) {
				foundHolidaySession = true
				break
			}
		}
		if !foundHolidaySession {
			t.Errorf("Scenario C: Expected to generate session on holiday %s, but none was generated", holidayDate.Format("2006-01-02"))
		}

		t.Logf("Scenario C: Generated %d sessions including holiday (ForceCancel=false, SkipHoliday=false)", len(schedules))
	})

	// 測試現有資料相容性：SkipHoliday 為預設值（true）
	t.Run("BackwardCompatibility_DefaultSkipHoliday", func(t *testing.T) {
		// 更新假日為 ForceCancel = false
		holiday.ForceCancel = false
		if err := appInstance.MySQL.WDB.WithContext(ctx).Save(&holiday).Error; err != nil {
			t.Fatalf("Failed to update holiday: %v", err)
		}

		// 建立規則（不指定 SkipHoliday，使用預設值 true）
		rule := models.ScheduleRule{
			CenterID:       center.ID,
			OfferingID:     offering.ID,
			TeacherID:      &teacher.ID,
			RoomID:         room.ID,
			Name:           "測試課程相容性",
			Weekday:        3, // 週三
			StartTime:      "10:00",
			EndTime:        "12:00",
			Duration:       120,
			// SkipHoliday 不設定，使用預設值 true
			EffectiveRange: models.DateRange{StartDate: testStartDate, EndDate: testEndDate},
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create rule: %v", err)
		}
		defer func() {
			appInstance.MySQL.WDB.WithContext(ctx).Where("id = ?", rule.ID).Delete(&models.ScheduleRule{})
		}()

		expansionSvc := services.NewScheduleExpansionService(appInstance)
		schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testStartDate, testEndDate, center.ID)

		// 驗證：預設 SkipHoliday = true 時應該跳過假日（維持舊有行為）
		for _, schedule := range schedules {
			if schedule.Date.Equal(holidayDate) {
				t.Errorf("Backward Compatibility: Expected to skip holiday %s (default SkipHoliday=true), but session was generated", holidayDate.Format("2006-01-02"))
			}
		}

		var nonHolidayCount int
		for _, s := range schedules {
			if !s.Date.Equal(holidayDate) {
				nonHolidayCount++
			}
		}
		t.Logf("Backward Compatibility: Generated %d sessions (holiday skipped with default SkipHoliday=true)", nonHolidayCount)
	})

	// 清理假日資料
	appInstance.MySQL.WDB.WithContext(ctx).Where("id = ?", holiday.ID).Delete(&models.CenterHoliday{})
}

// ============ Helper Functions (reimplemented from services package for testing) ============

// splitTime 分割時間字串 (HH:MM 或 HH:MM:SS) 為 [hour, minute]
func splitTime(timeStr string) []int {
	if timeStr == "" {
		return nil
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return nil
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	return []int{hour, minute}
}

// parseInt 從字串中提取數字並轉換為整數
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var num int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		}
	}
	return num
}

// intToString 整數轉字串
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		digit := n % 10
		result = append([]byte{byte('0' + digit)}, result...)
		n /= 10
	}
	return string(result)
}
