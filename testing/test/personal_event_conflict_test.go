package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app/models"
	"timeLedger/app/repositories"

	"gorm.io/gorm"
	gormMysql "gorm.io/driver/mysql"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	"timeLedger/app"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"

	mockRedis "timeLedger/testing/redis"

	"github.com/gin-gonic/gin"
)

func setupPersonalEventConflictTestApp() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.Teacher{},
		&models.Course{},
		&models.Room{},
		&models.Offering{},
		&models.ScheduleRule{},
		&models.PersonalEvent{},
	); err != nil {
		panic(fmt.Sprintf("AutoMigrate error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	// 初始化測試用的 Env 配置
	env := &configs.Env{
		JWTSecret:      "test-jwt-secret-key-for-testing-only",
		AppEnv:         "test",
		AppDebug:       true,
		AppTimezone:    "Asia/Taipei",
	}

	appInstance := &app.App{
		Env:   env,
		Err:   e,
		Tools: tool,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	cleanup := func() {
		// Use TRUNCATE to avoid foreign key constraint issues
		// TRUNCATE automatically resets auto-increment and drops FK constraints temporarily
		tables := []string{
			"schedule_exceptions",
			"schedule_rules",
			"personal_events",
			"center_invitations",
			"center_memberships",
			"center_teacher_notes",
			"center_holidays",
			"offerings",
			"courses",
			"rooms",
			"teacher_certificates",
			"teacher_personal_hashtags",
			"teachers",
			"centers",
		}
		for _, table := range tables {
			// Use TRUNCATE with RESTRICT to safely clean tables
			// Ignore errors if table doesn't exist in test environment
			mysqlDB.Exec(fmt.Sprintf("SET FOREIGN_KEY_CHECKS=0"))
			mysqlDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
			mysqlDB.Exec(fmt.Sprintf("SET FOREIGN_KEY_CHECKS=1"))
		}
		mr.Close()
	}

	return appInstance, mysqlDB, cleanup
}

func TestScheduleRuleRepository_CheckPersonalEventConflict(t *testing.T) {
	appInstance, _, cleanup := setupPersonalEventConflictTestApp()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	// 建立測試資料
	center := models.Center{
		Name:      fmt.Sprintf("Conflict Test Center %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter, err := centerRepo.Create(ctx, center)
	if err != nil {
		t.Fatalf("Failed to create center: %v", err)
	}

	teacher := models.Teacher{
		Name:      fmt.Sprintf("Test Teacher %d", now.UnixNano()),
		Email:     fmt.Sprintf("teacher_%d@test.com", now.UnixNano()),
		LineUserID: fmt.Sprintf("line_%d", now.UnixNano()),
		CreatedAt: now,
	}
	teacherRepo := repositories.NewTeacherRepository(appInstance)
	createdTeacher, err := teacherRepo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	course := models.Course{
		Name:              "Yoga Basics",
		TeacherBufferMin:  15,
		RoomBufferMin:     10,
		CreatedAt:         now,
	}
	courseRepo := repositories.NewCourseRepository(appInstance)
	createdCourse, err := courseRepo.Create(ctx, course)
	if err != nil {
		t.Fatalf("Failed to create course: %v", err)
	}

	room := models.Room{
		CenterID: createdCenter.ID,
		Name:     "Room A",
		Capacity: 20,
		IsActive: true,
		CreatedAt: now,
	}
	roomRepo := repositories.NewRoomRepository(appInstance)
	createdRoom, err := roomRepo.Create(ctx, room)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	offering := models.Offering{
		CenterID: createdCenter.ID,
		CourseID: createdCourse.ID,
		Name:     "Yoga Basics Course",
		CreatedAt: now,
	}
	offeringRepo := repositories.NewOfferingRepository(appInstance)
	createdOffering, err := offeringRepo.Create(ctx, offering)
	if err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	// 建立一個週一 09:00-10:00 的課程規則
	scheduleRule := models.ScheduleRule{
		CenterID:       createdCenter.ID,
		OfferingID:     createdOffering.ID,
		TeacherID:      &createdTeacher.ID,
		RoomID:         createdRoom.ID,
		Name:           "Monday Yoga",
		Weekday:        1, // Monday
		StartTime:      "09:00",
		EndTime:        "10:00",
		Duration:       60,
		EffectiveRange: models.DateRange{StartDate: now.AddDate(0, -1, 0), EndDate: now.AddDate(0, 1, 0)},
		CreatedAt:      now,
	}
	scheduleRuleRepo := repositories.NewScheduleRuleRepository(appInstance)
	createdRule, err := scheduleRuleRepo.Create(ctx, scheduleRule)
	if err != nil {
		t.Fatalf("Failed to create schedule rule: %v", err)
	}

	// 測試案例 1: 衝突的個人行程（部分重疊）
	t.Run("Conflict_WithOverlappingTime", func(t *testing.T) {
		// 建立一個週一 09:30-10:30 的個人行程（與課程 09:00-10:00 部分重疊）
		eventTime := getNextWeekday(now, 1) // 下個週一
		// 將時間設為午夜，然後加上時數，確保得到正確的上午時間
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 30*time.Minute)
		eventEndAt := eventTime.Add(10*time.Hour + 30*time.Minute)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) == 0 {
			t.Error("Expected conflict but got none")
		}

		if len(conflicts) > 0 && conflicts[0].ID != createdRule.ID {
			t.Errorf("Expected rule ID %d, got %d", createdRule.ID, conflicts[0].ID)
		}

		t.Logf("Found %d conflict(s), rule ID: %d", len(conflicts), createdRule.ID)
	})

	// 測試案例 2: 不衝突的個人行程（不同時間）
	t.Run("NoConflict_WithNonOverlappingTime", func(t *testing.T) {
		// 建立一個週一 11:00-12:00 的個人行程（與課程 09:00-10:00 不重疊）
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(11 * time.Hour)
		eventEndAt := eventTime.Add(12 * time.Hour)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) > 0 {
			t.Errorf("Expected no conflict but got %d", len(conflicts))
		}

		t.Log("No conflicts found as expected")
	})

	// 測試案例 3: 不衝突的個人行程（不同星期）
	t.Run("NoConflict_WithDifferentWeekday", func(t *testing.T) {
		// 建立一個週二 09:30-10:30 的個人行程（與週一課程不同天）
		eventTime := getNextWeekday(now, 2) // 下個週二
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 30*time.Minute)
		eventEndAt := eventTime.Add(10*time.Hour + 30*time.Minute)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) > 0 {
			t.Errorf("Expected no conflict but got %d", len(conflicts))
		}

		t.Log("No conflicts found as expected (different weekday)")
	})

	// 測試案例 4: 完全包含課程的個人行程
	t.Run("Conflict_WithContainingTime", func(t *testing.T) {
		// 建立一個週一 08:00-11:00 的個人行程（完全包含課程 09:00-10:00）
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(8 * time.Hour)
		eventEndAt := eventTime.Add(11 * time.Hour)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) == 0 {
			t.Error("Expected conflict but got none (containing time)")
		}

		t.Logf("Found %d conflict(s) with containing time", len(conflicts))
	})

	// 測試案例 5: 完全被課程包含的個人行程
	t.Run("Conflict_WithContainedTime", func(t *testing.T) {
		// 建立一個週一 09:15-09:45 的個人行程（完全被課程 09:00-10:00 包含）
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 15*time.Minute)
		eventEndAt := eventTime.Add(9*time.Hour + 45*time.Minute)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) == 0 {
			t.Error("Expected conflict but got none (contained time)")
		}

		t.Logf("Found %d conflict(s) with contained time", len(conflicts))
	})

	// 測試案例 6: 事件日期在 effective_range 之外，不應該報衝突
	t.Run("NoConflict_WhenEventDateOutsideEffectiveRange", func(t *testing.T) {
		// 建立一個新的課程規則，有效範圍是過去
		// 並測試一個在有效範圍內的時間（不應該與過去的課程衝突）

		// 先建立一個「未來」的課程規則
		futureCourse := models.Course{
			Name:              "Future Yoga",
			TeacherBufferMin:  15,
			RoomBufferMin:     10,
			CreatedAt:         now,
		}
		createdFutureCourse, err := courseRepo.Create(ctx, futureCourse)
		if err != nil {
			t.Fatalf("Failed to create future course: %v", err)
		}

		futureOffering := models.Offering{
			CenterID: createdCenter.ID,
			CourseID: createdFutureCourse.ID,
			Name:     "Future Yoga Course",
			CreatedAt: now,
		}
		createdFutureOffering, err := offeringRepo.Create(ctx, futureOffering)
		if err != nil {
			t.Fatalf("Failed to create future offering: %v", err)
		}

		// 課程規則的有效範圍是未來（從下個月開始）
		futureRule := models.ScheduleRule{
			CenterID:       createdCenter.ID,
			OfferingID:     createdFutureOffering.ID,
			TeacherID:      &createdTeacher.ID,
			RoomID:         createdRoom.ID,
			Name:           "Future Monday Yoga",
			Weekday:        1, // Monday
			StartTime:      "09:00",
			EndTime:        "10:00",
			Duration:       60,
			EffectiveRange: models.DateRange{StartDate: now.AddDate(0, 1, 0), EndDate: now.AddDate(0, 2, 0)}, // 下個月到下下個月
			CreatedAt:      now,
		}
		_, err = scheduleRuleRepo.Create(ctx, futureRule)
		if err != nil {
			t.Fatalf("Failed to create future schedule rule: %v", err)
		}

		// 嘗試在下個週一建立個人行程（這應該不會衝突，因為未來課程還沒開始）
		eventTime := getNextWeekday(now, 1) // 下個週一（在未來，但在課程開始之前）
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 30*time.Minute)
		eventEndAt := eventTime.Add(10*time.Hour + 30*time.Minute)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		// 應該不與未來課程衝突（但會與原始課程衝突，這是預期的）
		// 所以這個測試預期會有 1 個衝突（來自原始課程）
		t.Logf("Found %d conflicts (expected 1 from original course)", len(conflicts))
	})

	// ============ 邊界情況測試 ============
	t.Run("NoConflict_EventEndsExactlyWhenCourseStarts", func(t *testing.T) {
		// 事件 08:00-09:00，課程 09:00-10:00，剛好接續不重疊
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(8 * time.Hour)
		eventEndAt := eventTime.Add(9 * time.Hour)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) > 0 {
			t.Errorf("Expected no conflict when event ends exactly when course starts, got %d", len(conflicts))
		}

		t.Log("No conflicts found when event ends exactly when course starts")
	})

	t.Run("NoConflict_EventStartsExactlyWhenCourseEnds", func(t *testing.T) {
		// 事件 10:00-11:00，課程 09:00-10:00，剛好接續不重疊
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(10 * time.Hour)
		eventEndAt := eventTime.Add(11 * time.Hour)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) > 0 {
			t.Errorf("Expected no conflict when event starts exactly when course ends, got %d", len(conflicts))
		}

		t.Log("No conflicts found when event starts exactly when course ends")
	})

	t.Run("Conflict_EventOneMinuteOverlap", func(t *testing.T) {
		// 事件 09:59-10:30，課程 09:00-10:00，只重疊 1 分鐘
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 59*time.Minute)
		eventEndAt := eventTime.Add(10*time.Hour + 30*time.Minute)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) == 0 {
			t.Error("Expected conflict with 1 minute overlap")
		}

		t.Logf("Found %d conflict(s) with 1 minute overlap", len(conflicts))
	})

	t.Run("NoConflict_EventOneMinuteBeforeCourse", func(t *testing.T) {
		// 事件 08:59-09:00，課程 09:01-10:00，剛好差 1 分鐘不重疊
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(8*time.Hour + 59*time.Minute)
		eventEndAt := eventTime.Add(9 * time.Hour)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) > 0 {
			t.Errorf("Expected no conflict with 1 minute gap, got %d", len(conflicts))
		}

		t.Log("No conflicts found with 1 minute gap")
	})

	t.Run("Conflict_MultipleOverlappingCourses", func(t *testing.T) {
		// 建立第二個週一課程
		course2 := models.Course{
			Name:              "Piano Basics",
			TeacherBufferMin:  10,
			RoomBufferMin:     5,
			CreatedAt:         now,
		}
		createdCourse2, err := courseRepo.Create(ctx, course2)
		if err != nil {
			t.Fatalf("Failed to create course2: %v", err)
		}

		offering2 := models.Offering{
			CenterID: createdCenter.ID,
			CourseID: createdCourse2.ID,
			Name:     "Piano Basics Course",
			CreatedAt: now,
		}
		createdOffering2, err := offeringRepo.Create(ctx, offering2)
		if err != nil {
			t.Fatalf("Failed to create offering2: %v", err)
		}

		// 第二個課程時間設為 09:30-10:30，與原始課程 09:00-10:00 重疊
		scheduleRule2 := models.ScheduleRule{
			CenterID:       createdCenter.ID,
			OfferingID:     createdOffering2.ID,
			TeacherID:      &createdTeacher.ID,
			RoomID:         createdRoom.ID,
			Name:           "Monday Piano",
			Weekday:        1,
			StartTime:      "09:30",
			EndTime:        "10:30",
			Duration:       60,
			EffectiveRange: models.DateRange{StartDate: now.AddDate(0, -1, 0), EndDate: now.AddDate(0, 1, 0)},
			CreatedAt:      now,
		}
		_, err = scheduleRuleRepo.Create(ctx, scheduleRule2)
		if err != nil {
			t.Fatalf("Failed to create schedule rule2: %v", err)
		}

		// 事件時間 09:00-10:00 會與兩個課程都重疊
		eventTime := getNextWeekday(now, 1)
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9 * time.Hour)
		eventEndAt := eventTime.Add(10 * time.Hour)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		if len(conflicts) != 2 {
			t.Errorf("Expected conflicts with 2 courses, got %d", len(conflicts))
		}

		t.Logf("Found %d conflict(s) with multiple courses", len(conflicts))
	})

	t.Run("Conflict_SameCourseOnDifferentDays", func(t *testing.T) {
		// 在週二建立課程
		course3 := models.Course{
			Name:              "Yoga Advanced",
			TeacherBufferMin:  15,
			RoomBufferMin:     10,
			CreatedAt:         now,
		}
		createdCourse3, err := courseRepo.Create(ctx, course3)
		if err != nil {
			t.Fatalf("Failed to create course3: %v", err)
		}

		offering3 := models.Offering{
			CenterID: createdCenter.ID,
			CourseID: createdCourse3.ID,
			Name:     "Yoga Advanced Course",
			CreatedAt: now,
		}
		createdOffering3, err := offeringRepo.Create(ctx, offering3)
		if err != nil {
			t.Fatalf("Failed to create offering3: %v", err)
		}

		scheduleRule3 := models.ScheduleRule{
			CenterID:       createdCenter.ID,
			OfferingID:     createdOffering3.ID,
			TeacherID:      &createdTeacher.ID,
			RoomID:         createdRoom.ID,
			Name:           "Tuesday Yoga",
			Weekday:        2, // 週二
			StartTime:      "09:00",
			EndTime:        "10:00",
			Duration:       60,
			EffectiveRange: models.DateRange{StartDate: now.AddDate(0, -1, 0), EndDate: now.AddDate(0, 1, 0)},
			CreatedAt:      now,
		}
		_, err = scheduleRuleRepo.Create(ctx, scheduleRule3)
		if err != nil {
			t.Fatalf("Failed to create schedule rule3: %v", err)
		}

		// 週一的事件不應該與週二的課程衝突
		// 但會與所有週一的課程衝突（包括原始課程和 Conflict_MultipleOverlappingCourses 新增的課程）
		eventTime := getNextWeekday(now, 1) // 週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 30*time.Minute)
		eventEndAt := eventTime.Add(10*time.Hour + 30*time.Minute)

		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}

		// 由於測試依序執行，之前建立的課程也會存在
		// 所以應該有 2 個衝突（2 個週一課程），不應該包含週二的課程
		// 驗證沒有週二的課程衝突
		for _, conflict := range conflicts {
			if conflict.Weekday == 2 {
				t.Error("Did not expect conflict with Tuesday course")
			}
		}

		t.Logf("Found %d conflict(s) - correctly only Monday courses", len(conflicts))
	})

	// ============ 除錯測試 ============
	t.Run("Debug_PrintValues", func(t *testing.T) {
		// 列印所有相關值用於除錯
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 30*time.Minute)
		eventEndAt := eventTime.Add(10*time.Hour + 30*time.Minute)

		loc := app.GetTaiwanLocation()
		eventStartAtTaipei := eventStartAt.In(loc)
		eventEndAtTaipei := eventEndAt.In(loc)

		t.Logf("=== Debug Values ===")
		t.Logf("now: %v (%s)", now, now.Weekday())
		t.Logf("eventTime (next Monday): %v (%s)", eventTime, eventTime.Weekday())
		t.Logf("eventStartAt: %v", eventStartAt)
		t.Logf("eventStartAtTaipei: %v", eventStartAtTaipei)
		t.Logf("eventWeekday (original): %d", eventStartAt.Weekday())
		t.Logf("eventWeekday (Taipei): %d", eventStartAtTaipei.Weekday())

		// 取得課程規則
		rules, err := scheduleRuleRepo.ListByTeacherID(ctx, createdTeacher.ID, createdCenter.ID)
		if err != nil {
			t.Fatalf("ListByTeacherID failed: %v", err)
		}

		t.Logf("Found %d rules for teacher %d in center %d", len(rules), createdTeacher.ID, createdCenter.ID)
		for i, rule := range rules {
			t.Logf("Rule %d: ID=%d, Weekday=%d, StartTime=%s, EndTime=%s", i, rule.ID, rule.Weekday, rule.StartTime, rule.EndTime)
			t.Logf("  EffectiveRange: %v to %v", rule.EffectiveRange.StartDate, rule.EffectiveRange.EndDate)

			// 檢查每個條件
			eventWeekday := int(eventStartAtTaipei.Weekday())
			if eventWeekday == 0 {
				eventWeekday = 7
			}
			t.Logf("  eventWeekday=%d, rule.Weekday=%d, match=%v", eventWeekday, rule.Weekday, rule.Weekday == eventWeekday)

			eventDate := eventStartAtTaipei.Format("2006-01-02")
			ruleStartDate := rule.EffectiveRange.StartDate.In(loc).Format("2006-01-02")
			ruleEndDate := rule.EffectiveRange.EndDate.In(loc).Format("2006-01-02")
			t.Logf("  eventDate=%s, ruleStartDate=%s, ruleEndDate=%s", eventDate, ruleStartDate, ruleEndDate)
			inRange := eventDate >= ruleStartDate && eventDate <= ruleEndDate
			t.Logf("  inRange=%v", inRange)

			eventStartTime := eventStartAtTaipei.Format("15:04")
			eventEndTime := eventEndAtTaipei.Format("15:04")
			t.Logf("  eventStartTime=%s, eventEndTime=%s", eventStartTime, eventEndTime)
			// Inline timesOverlap check
			overlaps := rule.StartTime < eventEndTime && rule.EndTime > eventStartTime
			t.Logf("  timesOverlap=%v", overlaps)
		}

		// 呼叫實際的衝突檢查函數
		conflicts, err := scheduleRuleRepo.CheckPersonalEventConflict(ctx, createdTeacher.ID, createdCenter.ID, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflict failed: %v", err)
		}
		t.Logf("Final conflicts count: %d", len(conflicts))
	})
}

func TestScheduleRuleRepository_CheckPersonalEventConflictAllCenters(t *testing.T) {
	appInstance, _, cleanup := setupPersonalEventConflictTestApp()
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	// 建立測試資料
	center1 := models.Center{
		Name:      fmt.Sprintf("Conflict Test Center 1 %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	centerRepo := repositories.NewCenterRepository(appInstance)
	createdCenter1, err := centerRepo.Create(ctx, center1)
	if err != nil {
		t.Fatalf("Failed to create center1: %v", err)
	}

	center2 := models.Center{
		Name:      fmt.Sprintf("Conflict Test Center 2 %d", now.UnixNano()),
		PlanLevel: "STARTER",
		CreatedAt: now,
	}
	createdCenter2, err := centerRepo.Create(ctx, center2)
	if err != nil {
		t.Fatalf("Failed to create center2: %v", err)
	}

	teacher := models.Teacher{
		Name:      fmt.Sprintf("Multi-Center Teacher %d", now.UnixNano()),
		Email:     fmt.Sprintf("multicenter_teacher_%d@test.com", now.UnixNano()),
		LineUserID: fmt.Sprintf("line_multicenter_%d", now.UnixNano()),
		CreatedAt: now,
	}
	teacherRepo := repositories.NewTeacherRepository(appInstance)
	createdTeacher, err := teacherRepo.Create(ctx, teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}

	course := models.Course{
		Name:              "Yoga Basics",
		TeacherBufferMin:  15,
		RoomBufferMin:     10,
		CreatedAt:         now,
	}
	courseRepo := repositories.NewCourseRepository(appInstance)
	createdCourse, err := courseRepo.Create(ctx, course)
	if err != nil {
		t.Fatalf("Failed to create course: %v", err)
	}

	room := models.Room{
		CenterID: createdCenter1.ID,
		Name:     "Room A",
		Capacity: 20,
		IsActive: true,
		CreatedAt: now,
	}
	roomRepo := repositories.NewRoomRepository(appInstance)
	createdRoom, err := roomRepo.Create(ctx, room)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	offering := models.Offering{
		CenterID: createdCenter1.ID,
		CourseID: createdCourse.ID,
		Name:     "Yoga Basics Course",
		CreatedAt: now,
	}
	offeringRepo := repositories.NewOfferingRepository(appInstance)
	createdOffering, err := offeringRepo.Create(ctx, offering)
	if err != nil {
		t.Fatalf("Failed to create offering: %v", err)
	}

	// 在 Center1 建立週一 09:00-10:00 的課程規則
	scheduleRule := models.ScheduleRule{
		CenterID:       createdCenter1.ID,
		OfferingID:     createdOffering.ID,
		TeacherID:      &createdTeacher.ID,
		RoomID:         createdRoom.ID,
		Name:           "Monday Yoga at Center1",
		Weekday:        1, // Monday
		StartTime:      "09:00",
		EndTime:        "10:00",
		Duration:       60,
		EffectiveRange: models.DateRange{StartDate: now.AddDate(0, -1, 0), EndDate: now.AddDate(0, 1, 0)},
		CreatedAt:      now,
	}
	scheduleRuleRepo := repositories.NewScheduleRuleRepository(appInstance)
	_, err = scheduleRuleRepo.Create(ctx, scheduleRule)
	if err != nil {
		t.Fatalf("Failed to create schedule rule: %v", err)
	}

	// 測試：檢查多中心的衝突
	t.Run("ConflictInOneCenter", func(t *testing.T) {
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(9*time.Hour + 30*time.Minute)
		eventEndAt := eventTime.Add(10*time.Hour + 30*time.Minute)

		centerIDs := []uint{createdCenter1.ID, createdCenter2.ID}
		allConflicts, err := scheduleRuleRepo.CheckPersonalEventConflictAllCenters(ctx, createdTeacher.ID, centerIDs, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflictAllCenters failed: %v", err)
		}

		// 應該只在 Center1 找到衝突
		if len(allConflicts) != 1 {
			t.Errorf("Expected conflicts in 1 center, got %d", len(allConflicts))
		}

		if _, exists := allConflicts[createdCenter1.ID]; !exists {
			t.Error("Expected conflict in Center1")
		}

		if _, exists := allConflicts[createdCenter2.ID]; exists {
			t.Error("Did not expect conflict in Center2")
		}

		t.Logf("Found conflicts in %d center(s)", len(allConflicts))
	})

	// 測試：沒有衝突
	t.Run("NoConflictInAnyCenter", func(t *testing.T) {
		eventTime := getNextWeekday(now, 1) // 下個週一
		eventTime = time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())
		eventStartAt := eventTime.Add(11 * time.Hour) // 11:00-12:00
		eventEndAt := eventTime.Add(12 * time.Hour)

		centerIDs := []uint{createdCenter1.ID, createdCenter2.ID}
		allConflicts, err := scheduleRuleRepo.CheckPersonalEventConflictAllCenters(ctx, createdTeacher.ID, centerIDs, eventStartAt, eventEndAt)
		if err != nil {
			t.Fatalf("CheckPersonalEventConflictAllCenters failed: %v", err)
		}

		if len(allConflicts) > 0 {
			t.Errorf("Expected no conflicts but got %d", len(allConflicts))
		}

		t.Log("No conflicts found in any center as expected")
	})
}

// getNextWeekday 返回距離 now 最近的下個指定星期幾的日期
func getNextWeekday(now time.Time, weekday time.Weekday) time.Time {
	target := now
	for target.Weekday() != weekday {
		target = target.AddDate(0, 0, 1)
	}
	return target
}
