package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	mockRedis "timeLedger/testing/redis"

	"github.com/gin-gonic/gin"
)

func setupCrossDayTestApp() (*app.App, *gorm.DB) {
	gin.SetMode(gin.TestMode)

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
	return appInstance, mysqlDB
}

// TestCrossDayTimeFunctions 測試跨日時間處理函數
func TestCrossDayTimeFunctions(t *testing.T) {
	testApp, _ := setupCrossDayTestApp()
	_ = testApp

	ctx := context.Background()
	_ = ctx

	t.Run("IsCrossDayTime", func(t *testing.T) {
		// 正常課程（非跨日）
		if repositories.IsCrossDayTime("09:00", "10:00") {
			t.Error("09:00-10:00 should not be cross-day")
		}

		// 跨日課程
		if !repositories.IsCrossDayTime("23:00", "02:00") {
			t.Error("23:00-02:00 should be cross-day")
		}

		// 邊界情況
		if repositories.IsCrossDayTime("00:00", "23:59") {
			t.Error("00:00-23:59 should not be cross-day")
		}
	})

	t.Run("ParseTimeToMinutes", func(t *testing.T) {
		tests := []struct {
			timeStr  string
			expected int
		}{
			{"00:00", 0},
			{"01:00", 60},
			{"09:00", 540},
			{"12:00", 720},
			{"23:59", 1439},
		}

		for _, tt := range tests {
			result := repositories.ParseTimeToMinutes(tt.timeStr)
			if result != tt.expected {
				t.Errorf("ParseTimeToMinutes(%s) = %d, want %d", tt.timeStr, result, tt.expected)
			}
		}
	})

	t.Run("TimesOverlapCrossDay_NormalCourses", func(t *testing.T) {
		// 兩個普通課程重疊
		if !repositories.TimesOverlapCrossDay("09:00", "10:00", false, "09:30", "10:30", false) {
			t.Error("09:00-10:00 and 09:30-10:30 should overlap")
		}

		// 兩個普通課程不重疊
		if repositories.TimesOverlapCrossDay("09:00", "10:00", false, "10:00", "11:00", false) {
			t.Error("09:00-10:00 and 10:00-11:00 should not overlap")
		}

		// 邊界接續
		if repositories.TimesOverlapCrossDay("09:00", "10:00", false, "10:00", "11:00", false) {
			t.Error("09:00-10:00 and 10:00-11:00 should not overlap (adjacent)")
		}
	})

	t.Run("TimesOverlapCrossDay_CrossDayCourse", func(t *testing.T) {
		// 跨日課程 23:00-02:00 與當天課程 21:00-23:30 重疊
		if !repositories.TimesOverlapCrossDay("23:00", "02:00", true, "21:00", "23:30", false) {
			t.Error("23:00-02:00 and 21:00-23:30 should overlap")
		}

		// 跨日課程 23:00-02:00 與當天課程 22:00-23:00 不重疊（相鄰，無重疊）
		if repositories.TimesOverlapCrossDay("23:00", "02:00", true, "22:00", "23:00", false) {
			t.Error("23:00-02:00 and 22:00-23:00 should not overlap (adjacent)")
		}

		// 跨日課程 23:00-02:00 與當天課程 20:00-21:00 不重疊
		if repositories.TimesOverlapCrossDay("23:00", "02:00", true, "20:00", "21:00", false) {
			t.Error("23:00-02:00 and 20:00-21:00 should not overlap")
		}

		// 跨日課程 23:00-02:00 與當天晚間課程 23:30-23:59 重疊
		if !repositories.TimesOverlapCrossDay("23:00", "02:00", true, "23:30", "23:59", false) {
			t.Error("23:00-02:00 and 23:30-23:59 should overlap")
		}

		// 測試跨日課程與跨日課程（同一天）
		// 22:00-03:00 跨日課程與 23:00-02:00 跨日課程重疊
		if !repositories.TimesOverlapCrossDay("23:00", "02:00", true, "22:00", "03:00", true) {
			t.Error("Cross-day 23:00-02:00 and cross-day 22:00-03:00 should overlap")
		}

		// 跨日課程 23:00-02:00 與跨日課程 03:00-06:00 不重疊
		if repositories.TimesOverlapCrossDay("23:00", "02:00", true, "03:00", "06:00", true) {
			t.Error("Cross-day 23:00-02:00 and cross-day 03:00-06:00 should not overlap")
		}
	})

	t.Run("TimesOverlapCrossDay_BothCrossDay", func(t *testing.T) {
		// 兩個跨日課程重疊
		if !repositories.TimesOverlapCrossDay("23:00", "02:00", true, "00:00", "03:00", true) {
			t.Error("Two cross-day courses 23:00-02:00 and 00:00-03:00 should overlap")
		}

		// 兩個跨日課程不重疊
		if repositories.TimesOverlapCrossDay("23:00", "02:00", true, "03:00", "06:00", true) {
			t.Error("Cross-day courses 23:00-02:00 and 03:00-06:00 should not overlap")
		}

		// 另一個不重疊案例：22:00-01:00 與 02:00-05:00
		if repositories.TimesOverlapCrossDay("22:00", "01:00", true, "02:00", "05:00", true) {
			t.Error("Cross-day courses 22:00-01:00 and 02:00-05:00 should not overlap")
		}
	})

	t.Run("TimesOverlapCrossDayWithNextDay", func(t *testing.T) {
		// 測試 TimesOverlapCrossDayWithNextDay 函數
		// 跨日課程 23:00-02:00 與隔天凌晨課程 01:00-03:00 重疊
		if !repositories.TimesOverlapCrossDayWithNextDay("23:00", "02:00", true, "01:00", "03:00") {
			t.Error("23:00-02:00 and 01:00-03:00 (next day) should overlap")
		}

		// 跨日課程 23:00-02:00 與隔天凌晨課程 03:00-04:00 不重疊
		if repositories.TimesOverlapCrossDayWithNextDay("23:00", "02:00", true, "03:00", "04:00") {
			t.Error("23:00-02:00 and 03:00-04:00 (next day) should not overlap")
		}

		// 1分鐘重疊
		if !repositories.TimesOverlapCrossDayWithNextDay("23:00", "02:00", true, "01:00", "01:01") {
			t.Error("23:00-02:00 and 01:00-01:01 (next day) should overlap (1 minute)")
		}

		// 1分鐘間隔
		if repositories.TimesOverlapCrossDayWithNextDay("23:00", "02:00", true, "02:01", "03:00") {
			t.Error("23:00-02:00 and 02:01-03:00 (next day) should not overlap (1 minute gap)")
		}
	})
}

// TestCrossDayScheduleRule 測試跨日課程規則
func TestCrossDayScheduleRule(t *testing.T) {
	testApp, _ := setupCrossDayTestApp()
	ctx := context.Background()

	// 建立測試用的中心
	centerRepo := repositories.NewCenterRepository(testApp)
	center, err := centerRepo.GetByID(ctx, 1)
	if err != nil {
		t.Skip("No center data available, skipping test")
		return
	}

	// 建立測試用的房間
	roomRepo := repositories.NewRoomRepository(testApp)
	rooms, err := roomRepo.ListByCenterID(ctx, center.ID)
	if err != nil || len(rooms) == 0 {
		t.Skip("No rooms available, skipping test")
		return
	}
	room := rooms[0]

	// 建立測試用的老師
	teacherRepo := repositories.NewTeacherRepository(testApp)
	teachers, err := teacherRepo.List(ctx)
	if err != nil || len(teachers) == 0 {
		t.Skip("No teachers available, skipping test")
		return
	}
	teacher := teachers[0]

	t.Run("CreateCrossDayScheduleRule", func(t *testing.T) {
		// 建立跨日課程規則
		rule := models.ScheduleRule{
			CenterID:   center.ID,
			OfferingID: 1,
			TeacherID:  &teacher.ID,
			RoomID:     room.ID,
			Name:       "跨日瑜伽課程",
			Weekday:    1, // 週一
			StartTime:  "23:00",
			EndTime:    "02:00",
			Duration:   180, // 3小時
			IsCrossDay: true,
			EffectiveRange: models.DateRange{
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 3, 0),
			},
		}

		ruleRepo := repositories.NewScheduleRuleRepository(testApp)
		createdRule, err := ruleRepo.Create(ctx, rule)
		if err != nil {
			t.Fatalf("Failed to create cross-day schedule rule: %v", err)
		}

		// 驗證 IsCrossDay 欄位
		if !createdRule.IsCrossDay {
			t.Error("Created rule should have IsCrossDay = true")
		}

		// 清理測試資料
		defer ruleRepo.Delete(ctx, createdRule.ID)
	})

	t.Run("CheckOverlap_CrossDayWithNormal", func(t *testing.T) {
		ruleRepo := repositories.NewScheduleRuleRepository(testApp)

		// 建立普通課程（21:00-22:00）
		normalRule := models.ScheduleRule{
			CenterID:   center.ID,
			OfferingID: 1,
			TeacherID:  &teacher.ID,
			RoomID:     room.ID,
			Name:       "普通晚間課程",
			Weekday:    1,
			StartTime:  "21:00",
			EndTime:    "22:00",
			Duration:   60,
			EffectiveRange: models.DateRange{
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 3, 0),
			},
		}

		createdNormalRule, err := ruleRepo.Create(ctx, normalRule)
		if err != nil {
			t.Fatalf("Failed to create normal schedule rule: %v", err)
		}
		defer ruleRepo.Delete(ctx, createdNormalRule.ID)

		// 建立跨日課程（23:00-02:00）
		crossDayRule := models.ScheduleRule{
			CenterID:   center.ID,
			OfferingID: 1,
			TeacherID:  &teacher.ID,
			RoomID:     room.ID,
			Name:       "跨日課程",
			Weekday:    1,
			StartTime:  "23:00",
			EndTime:    "02:00",
			Duration:   180,
			IsCrossDay: true,
			EffectiveRange: models.DateRange{
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 3, 0),
			},
		}

		createdCrossDayRule, err := ruleRepo.Create(ctx, crossDayRule)
		if err != nil {
			t.Fatalf("Failed to create cross-day schedule rule: %v", err)
		}
		defer ruleRepo.Delete(ctx, createdCrossDayRule.ID)

		// 檢查跨日課程是否與普通課程衝突
		conflicts, _, err := ruleRepo.CheckOverlap(ctx, center.ID, room.ID, &teacher.ID, 1, "23:00", "02:00", &createdCrossDayRule.ID, time.Now())
		if err != nil {
			t.Fatalf("CheckOverlap failed: %v", err)
		}

		// 應該檢測到衝突
		if len(conflicts) == 0 {
			t.Error("Cross-day course should conflict with overlapping normal course")
		}
	})

	t.Run("CheckOverlap_CrossDayNoOverlap", func(t *testing.T) {
		ruleRepo := repositories.NewScheduleRuleRepository(testApp)

		// 使用 weekday=6（週六），因為這個 weekday 沒有現有課程
		testWeekday := 6

		// 建立普通課程（19:00-20:00）- 不會與跨日課程衝突
		normalRule := models.ScheduleRule{
			CenterID:   center.ID,
			OfferingID: 1,
			TeacherID:  &teacher.ID,
			RoomID:     room.ID,
			Name:       "普通晚間課程",
			Weekday:    testWeekday,
			StartTime:  "19:00",
			EndTime:    "20:00",
			Duration:   60,
			EffectiveRange: models.DateRange{
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 3, 0),
			},
		}

		createdNormalRule, err := ruleRepo.Create(ctx, normalRule)
		if err != nil {
			t.Fatalf("Failed to create normal schedule rule: %v", err)
		}
		defer ruleRepo.Delete(ctx, createdNormalRule.ID)

		// 建立跨日課程（23:00-02:00）
		crossDayRule := models.ScheduleRule{
			CenterID:   center.ID,
			OfferingID: 1,
			TeacherID:  &teacher.ID,
			RoomID:     room.ID,
			Name:       "跨日課程",
			Weekday:    testWeekday,
			StartTime:  "23:00",
			EndTime:    "02:00",
			Duration:   180,
			IsCrossDay: true,
			EffectiveRange: models.DateRange{
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 3, 0),
			},
		}

		createdCrossDayRule, err := ruleRepo.Create(ctx, crossDayRule)
		if err != nil {
			t.Fatalf("Failed to create cross-day schedule rule: %v", err)
		}
		defer ruleRepo.Delete(ctx, createdCrossDayRule.ID)

		// 檢查跨日課程是否與普通課程衝突
		conflicts, _, err := ruleRepo.CheckOverlap(ctx, center.ID, room.ID, &teacher.ID, testWeekday, "23:00", "02:00", &createdCrossDayRule.ID, time.Now())
		if err != nil {
			t.Fatalf("CheckOverlap failed: %v", err)
		}

		// 不應該檢測到衝突
		if len(conflicts) > 0 {
			t.Error("Cross-day course should not conflict with non-overlapping normal course")
		}
	})
}
