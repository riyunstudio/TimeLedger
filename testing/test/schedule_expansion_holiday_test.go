package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"
	mockRedis "timeLedger/testing/redis"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupExpansionTestApp() (*app.App, *gorm.DB, func()) {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	// AutoMigrate required tables
	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.Course{},
		&models.Offering{},
		&models.Room{},
		&models.Teacher{},
		&models.ScheduleRule{},
		&models.CenterHoliday{},
	); err != nil {
		panic(fmt.Sprintf("AutoMigrate error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
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
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	cleanup := func() {
		mr.Close()
	}

	return appInstance, mysqlDB, cleanup
}

func TestScheduleExpansionService_HolidayFiltering(t *testing.T) {
	t.Run("HolidaySessionsShouldBeFiltered", func(t *testing.T) {
		appInstance, db, cleanup := setupExpansionTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		// 1. 建立測試中心
		center, err := factory.CreateTestCenter(ctx, "HolidayTest")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 2. 建立測試課程和方案
		course, err := factory.CreateTestCourse(ctx, center.ID, "HolidayTest")
		if err != nil {
			t.Fatalf("建立測試課程失敗: %v", err)
		}

		offering, err := factory.CreateTestOffering(ctx, center.ID, course.ID, "HolidayTest")
		if err != nil {
			t.Fatalf("建立測試方案失敗: %v", err)
		}

		// 3. 建立測試教室
		room, err := factory.CreateTestRoom(ctx, center.ID, "HolidayTest")
		if err != nil {
			t.Fatalf("建立測試教室失敗: %v", err)
		}

		// 4. 建立測試老師
		teacher, err := factory.CreateTestTeacher(ctx, "HolidayTest")
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 5. 建立排課規則（週三上課）
		// 2026-01-21 是週三
		rule := models.ScheduleRule{
			CenterID:   center.ID,
			OfferingID: offering.ID,
			TeacherID:  &teacher.ID,
			RoomID:     room.ID,
			Weekday:    3, // 週三
			StartTime:  "10:00:00",
			EndTime:    "12:00:00",
			Duration:   120,
			EffectiveRange: models.DateRange{
				StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&rule).Error; err != nil {
			t.Fatalf("建立排課規則失敗: %v", err)
		}

		// 6. 建立假日（2026年1月21日是週三）
		holiday := models.CenterHoliday{
			CenterID:  center.ID,
			Date:      time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC),
			Name:      "測試假日",
			CreatedAt: time.Now(),
		}
		if err := db.WithContext(ctx).Create(&holiday).Error; err != nil {
			t.Fatalf("建立假日失敗: %v", err)
		}

		// 清理函數
		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("id = ?", rule.ID).Delete(&models.ScheduleRule{})
			db.WithContext(ctx).Where("id = ?", holiday.ID).Delete(&models.CenterHoliday{})
		}()

		// 7. 測試 ExpandRules
		// 測試範圍：2026年1月13日至2月10日（涵蓋多個週三）
		// 週三日期：1/14, 1/21(假日), 1/28, 2/4, 2/11
		startDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)

		expansionSvc := services.NewScheduleExpansionService(appInstance)

		rules := []models.ScheduleRule{rule}
		schedules := expansionSvc.ExpandRules(ctx, rules, startDate, endDate, center.ID)

		// 驗證結果
		// 預期有 4 個週三會排課：1/14, 1/28, 2/4, 2/11
		// 1/21 是假日，應該被過濾掉
		expectedCount := 4
		if len(schedules) != expectedCount {
			t.Errorf("預期有 %d 個課表項目，實際有 %d 個", expectedCount, len(schedules))
		}

		// 驗證 1/21（假日）不在結果中
		for _, schedule := range schedules {
			if schedule.Date.Format("2006-01-02") == "2026-01-21" {
				t.Error("2026-01-21 是假日，不應該出現在課表結果中")
			}
		}

		// 驗證其他週三在結果中
		expectedDates := []string{"2026-01-14", "2026-01-28", "2026-02-04", "2026-02-11"}
		for _, expectedDate := range expectedDates {
			found := false
			for _, schedule := range schedules {
				if schedule.Date.Format("2006-01-02") == expectedDate {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s 應該在課表結果中", expectedDate)
			}
		}

		t.Logf("假日過濾測試通過，共 %d 個課表項目", len(schedules))
		for _, s := range schedules {
			t.Logf("  - %s: %s %s", s.Date.Format("2006-01-02"), s.StartTime, s.EndTime)
		}
	})
}
