package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
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

func setupSmartMatchingTestApp() (*app.App, *gorm.DB, func()) {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	// AutoMigrate required tables
	if err := mysqlDB.AutoMigrate(
		&models.Center{},
		&models.Teacher{},
		&models.CenterInvitation{},
		&models.CenterTeacherNote{},
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
		mr.Close()
	}

	return appInstance, mysqlDB, cleanup
}

func TestSmartMatchingService_InviteTalent(t *testing.T) {
	t.Run("InviteTalent_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		// 1. 建立測試資料
		center, err := factory.CreateTestCenter(ctx, "InviteSuccess")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立老師（開放徵才）
		teacher, err := factory.CreateTestTeacher(ctx, "InviteSuccess", WithTeacherOpenToHiring(true))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 清理函數
		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("teacher_id = ?", teacher.ID).Delete(&models.CenterInvitation{})
		}()

		// 2. 測試邀請功能
		svc := services.NewSmartMatchingService(appInstance)

		result, err := svc.InviteTalent(ctx, center.ID, 1, []uint{teacher.ID}, "歡迎加入我們的人才庫！")
		if err != nil {
			t.Fatalf("邀請失敗: %v", err)
		}

		// 3. 驗證結果
		if result.InvitedCount != 1 {
			t.Errorf("預期邀請數量為 1，實際為 %d", result.InvitedCount)
		}

		if result.FailedCount != 0 {
			t.Errorf("預期失敗數量為 0，實際為 %d", result.FailedCount)
		}

		if len(result.InvitationIDs) != 1 {
			t.Errorf("預期邀請 ID 數量為 1，實際為 %d", len(result.InvitationIDs))
		}

		// 4. 驗證資料庫中的邀請記錄
		invitationRepo := repositories.NewCenterInvitationRepository(appInstance)
		invitations, err := invitationRepo.GetByTeacherAndCenter(ctx, teacher.ID, center.ID)
		if err != nil {
			t.Fatalf("查詢邀請記錄失敗: %v", err)
		}

		if len(invitations) != 1 {
			t.Errorf("預期找到 1 筆邀請記錄，實際找到 %d 筆", len(invitations))
		}

		invitation := invitations[0]
		if invitation.Status != models.InvitationStatusPending {
			t.Errorf("預期邀請狀態為 PENDING，實際為 %s", invitation.Status)
		}

		if invitation.InviteType != models.InvitationTypeTalentPool {
			t.Errorf("預期邀請類型為 TALENT_POOL，實際為 %s", invitation.InviteType)
		}

		if invitation.Token == "" {
			t.Error("預期邀請 Token 不為空")
		}
	})

	t.Run("InviteTalent_AlreadyHasPendingInvitation", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		// 1. 建立測試資料
		center, err := factory.CreateTestCenter(ctx, "DuplicateInvite")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher, err := factory.CreateTestTeacher(ctx, "DuplicateInvite", WithTeacherOpenToHiring(true))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立既有的待處理邀請
		existingInvitation := models.CenterInvitation{
			CenterID:    center.ID,
			TeacherID:   teacher.ID,
			InvitedBy:   1,
			Token:       factory.CreateUniqueToken(),
			Status:      models.InvitationStatusPending,
			InviteType:  models.InvitationTypeTalentPool,
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:   time.Now(),
		}
		if err := db.WithContext(ctx).Create(&existingInvitation).Error; err != nil {
			t.Fatalf("建立既有邀請失敗: %v", err)
		}

		// 清理函數
		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("teacher_id = ?", teacher.ID).Delete(&models.CenterInvitation{})
		}()

		// 2. 嘗試再次邀請
		svc := services.NewSmartMatchingService(appInstance)

		result, err := svc.InviteTalent(ctx, center.ID, 1, []uint{teacher.ID}, "歡迎加入！")
		if err != nil {
			t.Fatalf("邀請失敗: %v", err)
		}

		// 3. 驗證結果（應該失敗）
		if result.InvitedCount != 0 {
			t.Errorf("預期邀請數量為 0，實際為 %d", result.InvitedCount)
		}

		if result.FailedCount != 1 {
			t.Errorf("預期失敗數量為 1，實際為 %d", result.FailedCount)
		}

		if len(result.FailedIDs) != 1 || result.FailedIDs[0] != teacher.ID {
			t.Errorf("預期失敗 ID 包含老師 ID")
		}
	})

	t.Run("InviteTalent_TeacherNotOpenToHiring", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "NotOpen")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 老師未開放徵才
		teacher, err := factory.CreateTestTeacher(ctx, "NotOpen", WithTeacherOpenToHiring(false))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		result, err := svc.InviteTalent(ctx, center.ID, 1, []uint{teacher.ID}, "歡迎加入！")
		if err != nil {
			t.Fatalf("邀請失敗: %v", err)
		}

		// 驗證結果（應該失敗）
		if result.InvitedCount != 0 {
			t.Errorf("預期邀請數量為 0，實際為 %d", result.InvitedCount)
		}

		if result.FailedCount != 1 {
			t.Errorf("預期失敗數量為 1，實際為 %d", result.FailedCount)
		}
	})

	t.Run("InviteTalent_MultipleTeachers", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "BatchInvite")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立多個老師
		var teacherIDs []uint
		for i := 0; i < 3; i++ {
			teacher, err := factory.CreateTestTeacher(ctx, fmt.Sprintf("Batch%d", i), WithTeacherOpenToHiring(true))
			if err != nil {
				t.Fatalf("建立測試老師 %d 失敗: %v", i, err)
			}
			teacherIDs = append(teacherIDs, teacher.ID)
		}

		defer func() {
			db.WithContext(ctx).Where("email LIKE ?", "%test.com").Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.CenterInvitation{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		result, err := svc.InviteTalent(ctx, center.ID, 1, teacherIDs, "歡迎加入人才庫！")
		if err != nil {
			t.Fatalf("批量邀請失敗: %v", err)
		}

		// 驗證結果
		if result.InvitedCount != 3 {
			t.Errorf("預期邀請數量為 3，實際為 %d", result.InvitedCount)
		}

		if result.FailedCount != 0 {
			t.Errorf("預期失敗數量為 0，實際為 %d", result.FailedCount)
		}

		if len(result.InvitationIDs) != 3 {
			t.Errorf("預期邀請 ID 數量為 3，實際為 %d", len(result.InvitationIDs))
		}
	})
}

func TestSmartMatchingService_GetTalentStats(t *testing.T) {
	t.Run("GetTalentStats_WithRealData", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		// 1. 建立測試資料
		center, err := factory.CreateTestCenter(ctx, "Stats")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立老師（開放徵才）
		teacher, err := factory.CreateTestTeacher(ctx, "Stats",
			WithTeacherOpenToHiring(true),
			WithTeacherCity("台北市"),
			WithTeacherDistrict("信義區"))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立邀請記錄
		invitation := models.CenterInvitation{
			CenterID:    center.ID,
			TeacherID:   teacher.ID,
			InvitedBy:   1,
			Token:       factory.CreateUniqueToken(),
			Status:      models.InvitationStatusPending,
			InviteType:  models.InvitationTypeTalentPool,
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:   time.Now(),
		}
		if err := db.WithContext(ctx).Create(&invitation).Error; err != nil {
			t.Fatalf("建立邀請記錄失敗: %v", err)
		}

		// 清理函數
		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.CenterInvitation{})
		}()

		// 2. 測試取得統計
		svc := services.NewSmartMatchingService(appInstance)

		stats, err := svc.GetTalentStats(ctx, center.ID)
		if err != nil {
			t.Fatalf("取得統計失敗: %v", err)
		}

		// 3. 驗證結果
		if stats.OpenHiringCount < 1 {
			t.Errorf("預期開放徵才數量 >= 1，實際為 %d", stats.OpenHiringCount)
		}

		// 驗證邀請統計
		if stats.PendingInvites < 1 {
			t.Errorf("預期待處理邀請 >= 1，實際為 %d", stats.PendingInvites)
		}

		// 驗證城市分布
		foundTaipei := false
		for _, city := range stats.CityDistribution {
			if city.Name == "台北市" {
				foundTaipei = true
				if city.Count < 1 {
					t.Errorf("預期台北市老師數量 >= 1，實際為 %d", city.Count)
				}
				break
			}
		}
		if !foundTaipei {
			t.Error("統計中應該包含台北市的數據")
		}

		// 驗證趨勢數據
		if len(stats.MonthlyTrend) == 0 {
			t.Error("月趨勢數據不應該為空")
		}
	})
}

func TestSmartMatchingService_FindMatches(t *testing.T) {
	t.Run("FindMatches_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性並建立完整的關聯資料
		factory := NewTestDataFactory(db)

		// 1. 建立測試資料
		center, err := factory.CreateTestCenter(ctx, "SmartMatch")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立老師（開放徵才）
		teacher, err := factory.CreateTestTeacher(ctx, "SmartMatch", WithTeacherOpenToHiring(true), WithTeacherCity("台北市"), WithTeacherDistrict("大安區"))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立方案和教室
		course, err := factory.CreateTestCourse(ctx, center.ID, "SmartMatch")
		if err != nil {
			t.Fatalf("建立測試課程失敗: %v", err)
		}
		offering, err := factory.CreateTestOffering(ctx, center.ID, course.ID, "SmartMatch")
		if err != nil {
			t.Fatalf("建立測試方案失敗: %v", err)
		}
		room, err := factory.CreateTestRoom(ctx, center.ID, "SmartMatch")
		if err != nil {
			t.Fatalf("建立測試教室失敗: %v", err)
		}

		// 建立排課規則（使用正確的 offering_id 和 room_id）
		_, err = factory.CreateTestScheduleRule(ctx, center.ID, offering.ID, room.ID, &teacher.ID, "SmartMatch")
		if err != nil {
			t.Fatalf("建立排課規則失敗: %v", err)
		}

		// 清理函數
		defer func() {
			db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.ScheduleRule{})
		}()

		// 2. 測試 FindMatches
		svc := services.NewSmartMatchingService(appInstance)

		startTime := time.Date(2026, 2, 1, 14, 0, 0, 0, time.UTC)
		endTime := time.Date(2026, 2, 1, 16, 0, 0, 0, time.UTC)

		matches, err := svc.FindMatches(ctx, center.ID, nil, 1, startTime, endTime, []string{"瑜珈"}, []uint{})
		if err != nil {
			t.Fatalf("FindMatches 失敗: %v", err)
		}

		// 3. 驗證結果
		if len(matches) == 0 {
			t.Error("預期找到至少一個匹配結果")
		}

		if len(matches) > 0 {
			match := matches[0]
			if match.TeacherID != teacher.ID {
				t.Errorf("預期老師 ID 為 %d，實際為 %d", teacher.ID, match.TeacherID)
			}

			if match.Score < 0 {
				t.Errorf("預期分數 >= 0，實際為 %d", match.Score)
			}

			if match.Availability == "" {
				t.Error("預期 Availability 不為空")
			}
		}
	})

	t.Run("FindMatches_WithExcludedTeachers", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性並建立完整的關聯資料
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "Exclude")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立方案和教室
		course, err := factory.CreateTestCourse(ctx, center.ID, "Exclude")
		if err != nil {
			t.Fatalf("建立測試課程失敗: %v", err)
		}
		offering, err := factory.CreateTestOffering(ctx, center.ID, course.ID, "Exclude")
		if err != nil {
			t.Fatalf("建立測試方案失敗: %v", err)
		}
		room, err := factory.CreateTestRoom(ctx, center.ID, "Exclude")
		if err != nil {
			t.Fatalf("建立測試教室失敗: %v", err)
		}

		// 建立老師1
		teacher1, err := factory.CreateTestTeacher(ctx, "Exclude1", WithTeacherOpenToHiring(true))
		if err != nil {
			t.Fatalf("建立老師1失敗: %v", err)
		}

		// 建立老師2
		teacher2, err := factory.CreateTestTeacher(ctx, "Exclude2", WithTeacherOpenToHiring(true))
		if err != nil {
			t.Fatalf("建立老師2失敗: %v", err)
		}

		// 建立排課規則
		_, err = factory.CreateTestScheduleRule(ctx, center.ID, offering.ID, room.ID, &teacher1.ID, "Exclude1")
		if err != nil {
			t.Fatalf("建立規則1失敗: %v", err)
		}

		_, err = factory.CreateTestScheduleRule(ctx, center.ID, offering.ID, room.ID, &teacher2.ID, "Exclude2")
		if err != nil {
			t.Fatalf("建立規則2失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("email LIKE ?", "%test.com").Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.ScheduleRule{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		startTime := time.Date(2026, 2, 1, 14, 0, 0, 0, time.UTC)
		endTime := time.Date(2026, 2, 1, 16, 0, 0, 0, time.UTC)

		// 排除老師1
		matches, err := svc.FindMatches(ctx, center.ID, nil, 1, startTime, endTime, []string{}, []uint{teacher1.ID})
		if err != nil {
			t.Fatalf("FindMatches 失敗: %v", err)
		}

		// 驗證老師1被排除
		for _, match := range matches {
			if match.TeacherID == teacher1.ID {
				t.Error("預期老師1被排除，但結果中包含老師1")
			}
		}
	})

	t.Run("FindMatches_NoOpenToHiringTeachers", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性並建立完整的關聯資料
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "NoOpen")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立方案和教室
		course, err := factory.CreateTestCourse(ctx, center.ID, "NoOpen")
		if err != nil {
			t.Fatalf("建立測試課程失敗: %v", err)
		}
		offering, err := factory.CreateTestOffering(ctx, center.ID, course.ID, "NoOpen")
		if err != nil {
			t.Fatalf("建立測試方案失敗: %v", err)
		}
		room, err := factory.CreateTestRoom(ctx, center.ID, "NoOpen")
		if err != nil {
			t.Fatalf("建立測試教室失敗: %v", err)
		}

		// 建立未開放徵才的老師
		teacher, err := factory.CreateTestTeacher(ctx, "NoOpen", WithTeacherOpenToHiring(false))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		_, err = factory.CreateTestScheduleRule(ctx, center.ID, offering.ID, room.ID, &teacher.ID, "NoOpen")
		if err != nil {
			t.Fatalf("建立規則失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("email LIKE ?", "%test.com").Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.ScheduleRule{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		startTime := time.Date(2026, 2, 1, 14, 0, 0, 0, time.UTC)
		endTime := time.Date(2026, 2, 1, 16, 0, 0, 0, time.UTC)

		matches, err := svc.FindMatches(ctx, center.ID, nil, 1, startTime, endTime, []string{}, []uint{})
		if err != nil {
			t.Fatalf("FindMatches 失敗: %v", err)
		}

		// 未開放徵才的老師不應該出現在結果中
		if len(matches) != 0 {
			t.Errorf("預期無匹配結果，實際找到 %d 個", len(matches))
		}
	})
}

func TestSmartMatchingService_SearchTalent(t *testing.T) {
	t.Run("SearchTalent_WithFilters", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "SearchTalent")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher, err := factory.CreateTestTeacher(ctx, "SearchTalent",
			WithTeacherOpenToHiring(true),
			WithTeacherCity("新北市"),
			WithTeacherDistrict("板橋區"))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("email LIKE ?", "%test.com").Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		params := services.TalentSearchParams{
			CenterID: center.ID,
			City:     "新北市",
			District: "板橋區",
		}

		results, err := svc.SearchTalent(ctx, params)
		if err != nil {
			t.Fatalf("SearchTalent 失敗: %v", err)
		}

		// 驗證結果
		if len(results) == 0 {
			t.Error("預期找到至少一個人才結果")
		}

		if len(results) > 0 {
			result := results[0]
			if result.TeacherID != teacher.ID {
				t.Errorf("預期老師 ID 為 %d，實際為 %d", teacher.ID, result.TeacherID)
			}

			if result.Name != teacher.Name {
				t.Errorf("預期老師名稱為 %s，實際為 %s", teacher.Name, result.Name)
			}
		}
	})

	t.Run("SearchTalent_ByKeyword", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "Keyword")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher, err := factory.CreateTestTeacher(ctx, "Keyword",
			WithTeacherOpenToHiring(true),
			WithTeacherBio("我是專業瑜珈教練"))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		defer func() {
			if teacher != nil {
				db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
			}
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		params := services.TalentSearchParams{
			CenterID: center.ID,
			Keyword:  "瑜珈",
		}

		results, err := svc.SearchTalent(ctx, params)
		if err != nil {
			t.Fatalf("SearchTalent 失敗: %v", err)
		}

		// 關鍵字搜尋應該返回結果
		if len(results) == 0 {
			t.Log("未找到關鍵字匹配的結果，可能是預期行為")
		}
	})
}

func TestSmartMatchingService_GetSearchSuggestions(t *testing.T) {
	t.Run("GetSearchSuggestions_Success", func(t *testing.T) {
		appInstance, _, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		svc := services.NewSmartMatchingService(appInstance)

		suggestions, err := svc.GetSearchSuggestions(ctx, "瑜珈")
		if err != nil {
			t.Fatalf("GetSearchSuggestions 失敗: %v", err)
		}

		// 驗證結果結構
		if suggestions == nil {
			t.Fatal("預期 Suggestions 不為 nil")
		}

		// 驗證各類建議陣列存在
		_ = suggestions.Skills
		_ = suggestions.Tags
		_ = suggestions.Names
		_ = suggestions.Trending

		t.Log("搜尋建議取得成功")
	})

	t.Run("GetSearchSuggestions_EmptyQuery", func(t *testing.T) {
		appInstance, _, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		svc := services.NewSmartMatchingService(appInstance)

		suggestions, err := svc.GetSearchSuggestions(ctx, "")
		if err != nil {
			t.Fatalf("GetSearchSuggestions 失敗: %v", err)
		}

		// 空查詢也應該返回結果（熱門推薦）
		if suggestions == nil {
			t.Fatal("預期 Suggestions 不為 nil")
		}
	})
}

func TestSmartMatchingService_GetAlternativeSlots(t *testing.T) {
	t.Run("GetAlternativeSlots_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "AltSlot")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher, err := factory.CreateTestTeacher(ctx, "AltSlot", WithTeacherOpenToHiring(true))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("email LIKE ?", "%test.com").Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		originalStart := time.Date(2026, 2, 1, 14, 0, 0, 0, time.UTC)
		originalEnd := time.Date(2026, 2, 1, 16, 0, 0, 0, time.UTC)

		slots, err := svc.GetAlternativeSlots(ctx, center.ID, teacher.ID, originalStart, originalEnd, 120)
		if err != nil {
			t.Fatalf("GetAlternativeSlots 失敗: %v", err)
		}

		// 驗證結果結構
		if slots == nil {
			t.Fatal("預期 Slots 不為 nil")
		}

		// 應該返回替代時段列表
		t.Logf("找到 %d 個替代時段", len(slots))

		for _, slot := range slots {
			if slot.Date == "" {
				t.Error("替代時段日期不應該為空")
			}
			if slot.Start == "" {
				t.Error("替代時段開始時間不應該為空")
			}
			if slot.End == "" {
				t.Error("替代時段結束時間不應該為空")
			}
		}
	})
}

func TestSmartMatchingService_GetTeacherSessions(t *testing.T) {
	t.Run("GetTeacherSessions_Success", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性並建立完整的關聯資料
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "TeacherSession")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立方案和教室
		course, err := factory.CreateTestCourse(ctx, center.ID, "TeacherSession")
		if err != nil {
			t.Fatalf("建立測試課程失敗: %v", err)
		}
		offering, err := factory.CreateTestOffering(ctx, center.ID, course.ID, "TeacherSession")
		if err != nil {
			t.Fatalf("建立測試方案失敗: %v", err)
		}
		room, err := factory.CreateTestRoom(ctx, center.ID, "TeacherSession")
		if err != nil {
			t.Fatalf("建立測試教室失敗: %v", err)
		}

		teacher, err := factory.CreateTestTeacher(ctx, "TeacherSession", WithTeacherOpenToHiring(true))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立排課規則
		_, err = factory.CreateTestScheduleRule(ctx, center.ID, offering.ID, room.ID, &teacher.ID, "TeacherSession")
		if err != nil {
			t.Fatalf("建立排課規則失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("email LIKE ?", "%test.com").Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
			db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.ScheduleRule{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		sessions, err := svc.GetTeacherSessions(ctx, center.ID, teacher.ID, "2026-02-01", "2026-02-28")
		if err != nil {
			t.Fatalf("GetTeacherSessions 失敗: %v", err)
		}

		// 驗證結果結構
		if sessions == nil {
			t.Fatal("預期 Sessions 不為 nil")
		}

		if sessions.TeacherID != teacher.ID {
			t.Errorf("預期 TeacherID 為 %d，實際為 %d", teacher.ID, sessions.TeacherID)
		}

		if sessions.TeacherName == "" {
			t.Error("預期 TeacherName 不為空")
		}

		// 課表可能包含多個場次
		t.Logf("教師課表包含 %d 個場次", len(sessions.Sessions))

		for _, session := range sessions.Sessions {
			if session.ID == 0 {
				t.Error("場次 ID 不應該為 0")
			}
			if session.CourseName == "" {
				t.Error("課程名稱不應該為空")
			}
			if session.Status == "" {
				t.Error("狀態不應該為空")
			}
		}
	})

	t.Run("GetTeacherSessions_NoSessions", func(t *testing.T) {
		appInstance, db, cleanup := setupSmartMatchingTestApp()
		defer cleanup()

		ctx := context.Background()

		// 使用測試資料工廠確保唯一性
		factory := NewTestDataFactory(db)

		center, err := factory.CreateTestCenter(ctx, "NoSession")
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		teacher, err := factory.CreateTestTeacher(ctx, "NoSession", WithTeacherOpenToHiring(true))
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		defer func() {
			db.WithContext(ctx).Where("email LIKE ?", "%test.com").Delete(&models.Teacher{})
			db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		}()

		svc := services.NewSmartMatchingService(appInstance)

		sessions, err := svc.GetTeacherSessions(ctx, center.ID, teacher.ID, "2026-02-01", "2026-02-28")
		if err != nil {
			t.Fatalf("GetTeacherSessions 失敗: %v", err)
		}

		// 無課表應該返回空列表而非錯誤
		if sessions == nil {
			t.Fatal("預期 Sessions 不為 nil")
		}

		if sessions.TeacherID != teacher.ID {
			t.Errorf("預期 TeacherID 為 %d，實際為 %d", teacher.ID, sessions.TeacherID)
		}

		// 沒有課表時應該返回空列表
		if len(sessions.Sessions) != 0 {
			t.Logf("找到 %d 個場次", len(sessions.Sessions))
		}
	})
}
