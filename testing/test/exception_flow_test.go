package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	mockRedis "timeLedger/testing/redis"
)

func setupExceptionTestApp() *app.App {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	appInstance := &app.App{
		Env:   nil,
		Err:   e,
		Tools: tool,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	_ = mr
	return appInstance
}

// TestExpandRulesWithPendingException 測試 ExpandRules 對於待審核例外的處理
func TestExpandRulesWithPendingException(t *testing.T) {
	appInstance := setupExceptionTestApp()
	ctx := context.Background()
	now := time.Now()

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

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	// 使用 UTC 時間避免時區問題
	testStartDateUTC := now.AddDate(0, 0, 1)
	testEndDateUTC := now.AddDate(0, 3, 0)

	rule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  nil,
		RoomID:     room.ID,
		Name:       fmt.Sprintf("Test Pending %d", now.UnixNano()),
		Weekday:    3,
		StartTime:  "10:00",
		EndTime:    "11:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: testStartDateUTC,
			EndDate:   testEndDateUTC,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	// 計算符合規則 weekday 的目標日期（使用本地時間，然後轉 UTC）
	targetDateUTC := getNextWeekday(testStartDateUTC, time.Weekday(3))
	targetDateUTCUTC := time.Date(targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(), 0, 0, 0, 0, time.UTC)

	t.Logf("=== Test: ExpandRules with PENDING Exception ===")
	t.Logf("規則: ID=%d, Weekday=%d", rule.ID, rule.Weekday)
	t.Logf("目標日期: %s (週%d)", targetDateUTCUTC.Format("2006-01-02"), targetDateUTCUTC.Weekday())

	exception := models.ScheduleException{
		CenterID:     center.ID,
		RuleID:       rule.ID,
		OriginalDate: targetDateUTCUTC,
		Type:         "CANCEL",
		Status:       "PENDING",
		Reason:       fmt.Sprintf("Test pending exception %d", now.UnixNano()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&exception).Error; err != nil {
		t.Fatalf("Failed to create exception: %v", err)
	}
	t.Logf("建立 PENDING 例外: ID=%d", exception.ID)

	t.Logf("--- 驗證 ExpandRules 返回 HasException=true ---")

	expansionSvc := services.NewScheduleExpansionService(appInstance)
	schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, targetDateUTC, targetDateUTC, center.ID)

	if len(schedules) == 0 {
		t.Fatalf("ExpandRules 應該為例外日期產生課程但未產生")
	}

	schedule := schedules[0]
	t.Logf("Schedule: Date=%s, HasException=%v", schedule.Date.Format("2006-01-02"), schedule.HasException)

	if !schedule.HasException {
		t.Errorf("待審核例外的課程應該標記 HasException=true，但得到 false")
	}

	if schedule.ExceptionInfo == nil {
		t.Errorf("待審核例外的課程應該有 ExceptionInfo，但得到 nil")
	} else {
		t.Logf("ExceptionInfo: ID=%d, Type=%s, Status=%s",
			schedule.ExceptionInfo.ID, schedule.ExceptionInfo.Type, schedule.ExceptionInfo.Status)
	}

	t.Logf("=== PENDING Exception 測試通過 ===")
}

// TestExceptionFlowE2ECancel 端到端測試：教師申請停課 -> 管理員核准
func TestExceptionFlowE2ECancel(t *testing.T) {
	appInstance := setupExceptionTestApp()
	ctx := context.Background()
	now := time.Now()

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

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 使用固定日期測試：2026-01-29 是週四 (Weekday=4)
	// 這個日期在現在之後，確保有效
	testDate := time.Date(2026, 1, 29, 0, 0, 0, 0, time.UTC)
	testEndDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	rule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  nil,
		RoomID:     room.ID,
		Name:       fmt.Sprintf("E2E Cancel Test %d", now.UnixNano()),
		Weekday:    4, // 週四
		StartTime:  "09:00",
		EndTime:    "10:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: testDate,
			EndDate:   testEndDate,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	t.Logf("=== E2E Cancel Exception Flow Test ===")
	t.Logf("步驟 1: 建立排課規則 - Rule ID=%d (週四 09:00-10:00, 日期範圍: %s - %s)",
		rule.ID, testDate.Format("2006-01-02"), testEndDate.Format("2006-01-02"))

	// Step 2: 教師提交 CANCEL 例外申請
	exception := models.ScheduleException{
		CenterID:     center.ID,
		RuleID:       rule.ID,
		OriginalDate: testDate,
		Type:         "CANCEL",
		Status:       "PENDING",
		Reason:       fmt.Sprintf("E2E Test: Teacher sick leave %d", now.UnixNano()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&exception).Error; err != nil {
		t.Fatalf("Failed to create exception: %v", err)
	}
	t.Logf("步驟 2: 教師提交 CANCEL 例外 - Exception ID=%d, Status=PENDING", exception.ID)

	// Step 3: 驗證 PENDING 狀態下課程仍顯示，但標記 HasException
	expansionSvc := services.NewScheduleExpansionService(appInstance)

	// 除錯：先驗證規則可以產生課程
	t.Logf("除錯: 查詢日期=%s (Weekday=%d)", testDate.Format("2006-01-02"), testDate.Weekday())

	schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testDate, testDate, center.ID)
	t.Logf("除錯: ExpandRules 返回 %d 個 schedules", len(schedules))

	if len(schedules) == 0 {
		// 嘗試擴大查詢範圍
		allSchedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testDate, testEndDate, center.ID)
		t.Logf("除錯: 擴大範圍返回 %d 個 schedules", len(allSchedules))
		t.Fatalf("PENDING 例外應該產生課程但未產生")
	}

	schedule := schedules[0]
	t.Logf("步驟 3: PENDING 狀態 - Schedule HasException=%v", schedule.HasException)

	if !schedule.HasException {
		t.Errorf("PENDING 例外應該標記 HasException=true")
	}
	if schedule.ExceptionInfo == nil || schedule.ExceptionInfo.Status != "PENDING" {
		t.Errorf("PENDING 例外的 ExceptionInfo 應該有正確的 Status")
	}

	// Step 4: 管理員核准例外
	exception.Status = "APPROVED"
	exception.UpdatedAt = time.Now()
	if err := appInstance.MySQL.WDB.WithContext(ctx).Save(&exception).Error; err != nil {
		t.Fatalf("Failed to approve exception: %v", err)
	}
	t.Logf("步驟 4: 管理員核准例外 - Exception ID=%d, Status=APPROVED", exception.ID)

	// Step 5: 驗證 APPROVED CANCEL 讓 ExpandRules 跳過該日期
	schedules = expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testDate, testDate, center.ID)
	t.Logf("步驟 5: APPROVED 狀態 - ExpandRules 返回 %d 個 sessions", len(schedules))

	if len(schedules) > 0 {
		t.Errorf("已核准的 CANCEL 應該讓 ExpandRules 跳過該日期，但返回 %d 個 sessions", len(schedules))
	}

	t.Logf("=== E2E Cancel Exception Flow Test PASSED ===")
	t.Logf("流程驗證: 教師申請 CANCEL -> 管理員核准 -> 課程正確跳過")
}

// TestExceptionFlowE2EReschedule 端到端測試：教師申請改期 -> 管理員審核
func TestExceptionFlowE2EReschedule(t *testing.T) {
	appInstance := setupExceptionTestApp()
	ctx := context.Background()
	now := time.Now()

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

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 使用固定日期測試：2026-01-30 是週五 (Weekday=5)
	// 使用本地時區（台北 UTC+8）來避免 MySQL 時區轉換問題
	loc, _ := time.LoadLocation("Asia/Taipei")
	testDate := time.Date(2026, 1, 30, 0, 0, 0, 0, loc)
	testEndDate := time.Date(2026, 3, 1, 0, 0, 0, 0, loc)
	newStartAt := time.Date(2026, 1, 30, 16, 0, 0, 0, loc)
	newEndAt := time.Date(2026, 1, 30, 17, 0, 0, 0, loc)

	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	rule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  nil,
		RoomID:     room.ID,
		Name:       fmt.Sprintf("E2E Reschedule Test %d", now.UnixNano()),
		Weekday:    5, // 週五
		StartTime:  "14:00",
		EndTime:    "15:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: testDate,
			EndDate:   testEndDate,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	t.Logf("=== E2E Reschedule Exception Flow Test ===")
	t.Logf("步驟 1: 建立排課規則 - Rule ID=%d (週五 14:00-15:00)", rule.ID)

	// Step 2: 教師提交 RESCHEDULE 例外申請
	exception := models.ScheduleException{
		CenterID:     center.ID,
		RuleID:       rule.ID,
		OriginalDate: testDate,
		Type:         "RESCHEDULE",
		Status:       "PENDING",
		NewStartAt:   &newStartAt,
		NewEndAt:     &newEndAt,
		Reason:       fmt.Sprintf("E2E Test: Room change %d", now.UnixNano()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&exception).Error; err != nil {
		t.Fatalf("Failed to create exception: %v", err)
	}
	t.Logf("步驟 2: 教師提交 RESCHEDULE 例外 - Exception ID=%d (新時段 16:00-17:00)", exception.ID)

	// Step 3: 驗證 PENDING 狀態下課程仍顯示原時段
	expansionSvc := services.NewScheduleExpansionService(appInstance)
	schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testDate, testDate, center.ID)

	if len(schedules) == 0 {
		t.Fatalf("PENDING RESCHEDULE 應該產生課程但未產生")
	}
	schedule := schedules[0]
	t.Logf("步驟 3: PENDING 狀態 - Schedule StartTime=%s (原時段)", schedule.StartTime)

	if !schedule.HasException {
		t.Errorf("PENDING RESCHEDULE 應該標記 HasException=true")
	}
	if schedule.ExceptionInfo == nil || schedule.ExceptionInfo.Type != "RESCHEDULE" {
		t.Errorf("PENDING RESCHEDULE 應該有正確的 ExceptionInfo")
	}

	// Step 4: 管理員核准例外（使用服務層正確呼叫）
	exceptionSvc := services.NewScheduleExceptionService(appInstance)
	// 使用 adminID=1 進行審核
	if err := exceptionSvc.ReviewException(ctx, exception.ID, 1, "APPROVED", false, "Test approval"); err != nil {
		t.Fatalf("Failed to approve exception via service: %v", err)
	}
	t.Logf("步驟 4: 管理員核准例外 - Exception ID=%d", exception.ID)

	// Step 5: 驗證 RESCHEDULE 創建了新規則段
	var updatedRules []models.ScheduleRule
	if err := appInstance.MySQL.RDB.WithContext(ctx).
		Where("offering_id = ?", offering.ID).
		Order("id ASC").
		Find(&updatedRules).Error; err != nil {
		t.Fatalf("Failed to query updated rules: %v", err)
	}

	t.Logf("步驟 5: RESCHEDULE 後規則數量: %d (預期: 2)", len(updatedRules))

	if len(updatedRules) < 2 {
		t.Errorf("RESCHEDULE 應該創建新規則段，但只找到 %d 個規則", len(updatedRules))
	} else {
		newRule := updatedRules[1]
		t.Logf("新規則: ID=%d, StartTime=%s, EndTime=%s, Weekday=%d",
			newRule.ID, newRule.StartTime, newRule.EndTime, newRule.Weekday)

		if newRule.StartTime != "16:00" || newRule.EndTime != "17:00" {
			t.Errorf("新規則應該有正確的時間 16:00-17:00，但得到 %s-%s",
				newRule.StartTime, newRule.EndTime)
		}
	}

	t.Logf("=== E2E Reschedule Exception Flow Test PASSED ===")
	t.Logf("流程驗證: 教師申請改期 -> 管理員核准 -> 系統創建新規則段")
}

// TestExceptionFlowE2ERevoke 端到端測試：教師撤回例外申請
func TestExceptionFlowE2ERevoke(t *testing.T) {
	appInstance := setupExceptionTestApp()
	ctx := context.Background()
	now := time.Now()

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

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 使用固定日期測試：2026-01-31 是週六 (Weekday=6)
	loc, _ := time.LoadLocation("Asia/Taipei")
	testDate := time.Date(2026, 1, 31, 0, 0, 0, 0, loc)
	testEndDate := time.Date(2026, 3, 1, 0, 0, 0, 0, loc)

	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	rule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  nil,
		RoomID:     room.ID,
		Name:       fmt.Sprintf("E2E Revoke Test %d", now.UnixNano()),
		Weekday:    6, // 週六
		StartTime:  "10:00",
		EndTime:    "11:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: testDate,
			EndDate:   testEndDate,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	t.Logf("=== E2E Revoke Exception Flow Test ===")
	t.Logf("步驟 1: 建立排課規則 - Rule ID=%d (週六 10:00-11:00)", rule.ID)

	// Step 2: 教師提交 CANCEL 例外申請
	exception := models.ScheduleException{
		CenterID:     center.ID,
		RuleID:       rule.ID,
		OriginalDate: testDate,
		Type:         "CANCEL",
		Status:       "PENDING",
		Reason:       fmt.Sprintf("E2E Test: Teacher revokes %d", now.UnixNano()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&exception).Error; err != nil {
		t.Fatalf("Failed to create exception: %v", err)
	}
	t.Logf("步驟 2: 教師提交 CANCEL 例外 - Exception ID=%d, Status=PENDING", exception.ID)

	// Step 3: 教師撤回例外
	exceptionSvc := services.NewScheduleExceptionService(appInstance)
	// 使用 teacherID=1 進行撤回
	if err := exceptionSvc.RevokeException(ctx, exception.ID, 1); err != nil {
		t.Fatalf("Failed to revoke exception via service: %v", err)
	}
	t.Logf("步驟 3: 教師撤回例外 - Exception ID=%d", exception.ID)

	// Step 4: 驗證例外狀態變為 REVOKED
	var updatedException models.ScheduleException
	if err := appInstance.MySQL.RDB.WithContext(ctx).First(&updatedException, exception.ID).Error; err != nil {
		t.Fatalf("Failed to query updated exception: %v", err)
	}

	if updatedException.Status != "REVOKED" {
		t.Errorf("例外狀態應該是 REVOKED，但得到 %s", updatedException.Status)
	} else {
		t.Logf("步驟 4: 例外狀態已更新為 REVOKED")
	}

	// Step 5: 驗證 REVOKED 例外不影響課程產生
	expansionSvc := services.NewScheduleExpansionService(appInstance)
	schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testDate, testDate, center.ID)

	if len(schedules) == 0 {
		t.Fatalf("已撤回的例外不應該影響課程產生，但 ExpandRules 返回 0 個 sessions")
	}

	schedule := schedules[0]
	if schedule.HasException {
		t.Errorf("已撤回的例外不應該讓課程標記 HasException，但得到 true")
	}

	t.Logf("=== E2E Revoke Exception Flow Test PASSED ===")
	t.Logf("流程驗證: 教師申請 -> 教師撤回 -> 課程正常產生")
}

// TestExceptionFlowE2EReject 端到端測試：管理員拒絕例外申請
func TestExceptionFlowE2EReject(t *testing.T) {
	appInstance := setupExceptionTestApp()
	ctx := context.Background()
	now := time.Now()

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

	var room models.Room
	if err := appInstance.MySQL.RDB.WithContext(ctx).Where("center_id = ?", center.ID).First(&room).Error; err != nil {
		t.Skipf("No room data available, skipping test: %v", err)
		return
	}

	// 使用固定日期測試：2026-02-01 是週日 (Weekday=0)
	loc, _ := time.LoadLocation("Asia/Taipei")
	testDate := time.Date(2026, 2, 1, 0, 0, 0, 0, loc)
	testEndDate := time.Date(2026, 3, 1, 0, 0, 0, 0, loc)

	offering := models.Offering{
		CenterID:   center.ID,
		CourseID:   course.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&offering).Error; err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	rule := models.ScheduleRule{
		CenterID:   center.ID,
		OfferingID: offering.ID,
		TeacherID:  nil,
		RoomID:     room.ID,
		Name:       fmt.Sprintf("E2E Reject Test %d", now.UnixNano()),
		Weekday:    7, // 週日 (系統使用 7 表示週日)
		StartTime:  "14:00",
		EndTime:    "15:00",
		Duration:   60,
		EffectiveRange: models.DateRange{
			StartDate: testDate,
			EndDate:   testEndDate,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	t.Logf("=== E2E Reject Exception Flow Test ===")
	t.Logf("步驟 1: 建立排課規則 - Rule ID=%d (週日 14:00-15:00)", rule.ID)

	// Step 2: 教師提交 CANCEL 例外申請
	exception := models.ScheduleException{
		CenterID:     center.ID,
		RuleID:       rule.ID,
		OriginalDate: testDate,
		Type:         "CANCEL",
		Status:       "PENDING",
		Reason:       fmt.Sprintf("E2E Test: Admin rejects %d", now.UnixNano()),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := appInstance.MySQL.WDB.WithContext(ctx).Create(&exception).Error; err != nil {
		t.Fatalf("Failed to create exception: %v", err)
	}
	t.Logf("步驟 2: 教師提交 CANCEL 例外 - Exception ID=%d, Status=PENDING", exception.ID)

	// Step 3: 管理員拒絕例外
	exceptionSvc := services.NewScheduleExceptionService(appInstance)
	// 使用 adminID=1 進行拒絕
	if err := exceptionSvc.ReviewException(ctx, exception.ID, 1, "REJECTED", false, "Reason not sufficient"); err != nil {
		t.Fatalf("Failed to reject exception via service: %v", err)
	}
	t.Logf("步驟 3: 管理員拒絕例外 - Exception ID=%d", exception.ID)

	// Step 4: 驗證例外狀態變為 REJECTED
	var updatedException models.ScheduleException
	if err := appInstance.MySQL.RDB.WithContext(ctx).First(&updatedException, exception.ID).Error; err != nil {
		t.Fatalf("Failed to query updated exception: %v", err)
	}

	if updatedException.Status != "REJECTED" {
		t.Errorf("例外狀態應該是 REJECTED，但得到 %s", updatedException.Status)
	} else {
		t.Logf("步驟 4: 例外狀態已更新為 REJECTED")
	}

	// Step 5: 驗證 REJECTED 例外不影響課程產生
	expansionSvc := services.NewScheduleExpansionService(appInstance)
	schedules := expansionSvc.ExpandRules(ctx, []models.ScheduleRule{rule}, testDate, testDate, center.ID)

	if len(schedules) == 0 {
		t.Fatalf("已拒絕的例外不應該影響課程產生，但 ExpandRules 返回 0 個 sessions")
	}

	schedule := schedules[0]
	if schedule.HasException {
		t.Errorf("已拒絕的例外不應該讓課程標記 HasException，但得到 true")
	}

	t.Logf("=== E2E Reject Exception Flow Test PASSED ===")
	t.Logf("流程驗證: 教師申請 -> 管理員拒絕 -> 課程正常產生")
}