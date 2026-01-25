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
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"

	mockRedis "timeLedger/testing/redis"

	"github.com/gin-gonic/gin"
)

func setupPersonalEventConflictTestApp() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
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

	appInstance := &app.App{
		Env:   nil,
		Err:   e,
		Tools: tool,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	cleanup := func() {
		mysqlDB.Exec("DELETE FROM schedule_rules")
		mysqlDB.Exec("DELETE FROM personal_events")
		mysqlDB.Exec("DELETE FROM offerings")
		mysqlDB.Exec("DELETE FROM courses")
		mysqlDB.Exec("DELETE FROM rooms")
		mysqlDB.Exec("DELETE FROM teachers")
		mysqlDB.Exec("DELETE FROM centers")
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

	// 測試案例 1: 衝突的個人行程（同一時間）
	t.Run("Conflict_WithOverlappingTime", func(t *testing.T) {
		// 建立一個週一 09:30-10:30 的個人行程（與課程 09:00-10:00 重疊）
		eventTime := getNextWeekday(now, 1) // 下個週一
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
