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

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestTeacherMergeService_MergeTeacher 教師合併服務測試
// 測試場景：建立一個佔位老師與幾筆排課規則，執行合併後，驗證排課規則的 teacher_id 是否正確變更，且舊紀錄已刪除（軟刪除）
func TestTeacherMergeService_MergeTeacher(t *testing.T) {
	// 初始化測試資料庫
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("跳過測試 - 資料庫連線失敗: %v", err)
		return
	}
	defer CloseDB(db)

	ctx := context.Background()

	// 設定測試資料
	center, course, offering, room, sourceTeacher, targetTeacher, cleanup, err := setupMergeTestData(ctx, db)
	if err != nil {
		t.Skipf("跳過測試 - 無法建立測試資料: %v", err)
		return
	}
	defer cleanup()
	_ = course // course 在 setupMergeTestData 內部建立，此處不使用

	// 建立 App 實例
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// 建立合併服務
	mergeService := services.NewTeacherMergeService(appInstance)

	// 建立 3 筆測試排課規則
	baseTime := time.Now()

	rule1 := &models.ScheduleRule{
		CenterID:       center.ID,
		OfferingID:     offering.ID,
		TeacherID:      &sourceTeacher.ID,
		RoomID:         room.ID,
		Name:           "合併測試規則1",
		Weekday:        1,
		StartTime:      "14:00",
		EndTime:        "15:00",
		Duration:       60,
		IsCrossDay:     false,
		EffectiveRange: models.DateRange{
			StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		CreatedAt: baseTime,
		UpdatedAt: baseTime,
	}
	if err := db.WithContext(ctx).Create(rule1).Error; err != nil {
		t.Fatalf("建立排課規則1失敗: %v", err)
	}

	rule2 := &models.ScheduleRule{
		CenterID:       center.ID,
		OfferingID:     offering.ID,
		TeacherID:      &sourceTeacher.ID,
		RoomID:         room.ID,
		Name:           "合併測試規則2",
		Weekday:        2,
		StartTime:      "15:00",
		EndTime:        "16:00",
		Duration:       60,
		IsCrossDay:     false,
		EffectiveRange: models.DateRange{
			StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		CreatedAt: baseTime,
		UpdatedAt: baseTime,
	}
	if err := db.WithContext(ctx).Create(rule2).Error; err != nil {
		t.Fatalf("建立排課規則2失敗: %v", err)
	}

	rule3 := &models.ScheduleRule{
		CenterID:       center.ID,
		OfferingID:     offering.ID,
		TeacherID:      &sourceTeacher.ID,
		RoomID:         room.ID,
		Name:           "合併測試規則3",
		Weekday:        3,
		StartTime:      "16:00",
		EndTime:        "17:00",
		Duration:       60,
		IsCrossDay:     false,
		EffectiveRange: models.DateRange{
			StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		CreatedAt: baseTime,
		UpdatedAt: baseTime,
	}
	if err := db.WithContext(ctx).Create(rule3).Error; err != nil {
		t.Fatalf("建立排課規則3失敗: %v", err)
	}

	t.Logf("來源教師 ID: %d", sourceTeacher.ID)
	t.Logf("目標教師 ID: %d", targetTeacher.ID)
	t.Logf("建立測試排課規則: ID=%d, %d, %d", rule1.ID, rule2.ID, rule3.ID)

	// 記錄合併前來源教師的排課規則數量
	var initialSourceRules int64
	db.Model(&models.ScheduleRule{}).
		Where("teacher_id = ? AND center_id = ?", sourceTeacher.ID, center.ID).
		Count(&initialSourceRules)
	preMergeSourceRules := int(initialSourceRules)

	// ============================================
	// 執行合併
	// ============================================

	err = mergeService.MergeTeacher(ctx, sourceTeacher.ID, targetTeacher.ID, center.ID)
	assert.NoError(t, err, "合併不應失敗")
	t.Logf("合併執行成功")

	// ============================================
	// 驗證合併後的狀態
	// ============================================

	t.Run("驗證排課規則已遷移至目標教師", func(t *testing.T) {
		// 檢查目標教師的排課規則數量
		var targetRulesCount int64
		db.WithContext(ctx).Model(&models.ScheduleRule{}).
			Where("teacher_id = ? AND center_id = ?", targetTeacher.ID, center.ID).
			Count(&targetRulesCount)

		// 目標教師原有的 + 來源教師的 = 合併後總數
		var initialTargetRules int64
		db.Model(&models.ScheduleRule{}).
			Where("teacher_id = ? AND center_id = ?", targetTeacher.ID, center.ID).
			Count(&initialTargetRules)

		expectedTotal := preMergeSourceRules + int(initialTargetRules)
		assert.Equal(t, expectedTotal, int(targetRulesCount), "合併後目標教師排課規則數量應正確")
		t.Logf("合併後目標教師排課規則數量: %d (原有 %d + 新增 %d)", targetRulesCount, initialTargetRules, preMergeSourceRules)

		// 檢查來源教師的排課規則數量（應為 0，因為已遷移）
		var sourceRulesCount int64
		db.WithContext(ctx).Model(&models.ScheduleRule{}).
			Where("teacher_id = ? AND center_id = ?", sourceTeacher.ID, center.ID).
			Count(&sourceRulesCount)
		assert.Equal(t, int64(0), sourceRulesCount, "合併後來源教師應有 0 筆排課規則")
		t.Logf("合併後來源教師排課規則數量: %d", sourceRulesCount)

		// 驗證具體排課規則的 teacher_id 已變更
		var rules []models.ScheduleRule
		db.WithContext(ctx).Where("id IN (?, ?, ?)", rule1.ID, rule2.ID, rule3.ID).Find(&rules)

		for i, rule := range rules {
			assert.NotNil(t, rule.TeacherID, "排課規則 %d 的 TeacherID 不應為 nil", i+1)
			assert.Equal(t, targetTeacher.ID, *rule.TeacherID, "排課規則 %d 的 TeacherID 應變更為目標教師 ID", i+1)
			t.Logf("排課規則 %d (ID=%d) teacher_id 已正確變更為 %d ✓", i+1, rule.ID, *rule.TeacherID)
		}
	})

	t.Run("驗證來源教師已軟刪除", func(t *testing.T) {
		// 檢查來源教師是否被軟刪除
		var sourceTeacherCheck models.Teacher
		err := db.WithContext(ctx).Where("id = ?", sourceTeacher.ID).First(&sourceTeacherCheck).Error
		assert.Error(t, err, "來源教師應該找不到（已被軟刪除）")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "應該是 RecordNotFound 錯誤")

		// 使用 Unscoped 應該能找到
		var sourceTeacherUnscoped models.Teacher
		err = db.WithContext(ctx).Unscoped().Where("id = ?", sourceTeacher.ID).First(&sourceTeacherUnscoped).Error
		assert.NoError(t, err, "使用 Unscoped 應該能找到來源教師")
		assert.NotNil(t, sourceTeacherUnscoped.DeletedAt, "來源教師應有刪除時間")
		t.Logf("來源教師已軟刪除，刪除時間: %v", sourceTeacherUnscoped.DeletedAt)
	})

	// ============================================
	// 清理測試資料
	// ============================================

	t.Run("清理測試資料", func(t *testing.T) {
		// 刪除測試排課規則
		db.Where("id IN (?, ?, ?)", rule1.ID, rule2.ID, rule3.ID).Delete(&models.ScheduleRule{})
		// 恢復來源教師（移除軟刪除）
		db.Unscoped().Model(&sourceTeacher).Update("deleted_at", nil)
		t.Logf("測試資料已清理")
	})
}

// setupMergeTestData 設定合併測試所需的資料
func setupMergeTestData(ctx context.Context, db *gorm.DB) (*models.Center, *models.Course, *models.Offering, *models.Room, *models.Teacher, *models.Teacher, func(), error) {
	// 建立測試中心
	centerName := fmt.Sprintf("合併測試中心 - %d", time.Now().UnixNano())
	center := &models.Center{
		Name:      centerName,
		PlanLevel: "STARTER",
		CreatedAt: time.Now(),
	}
	if err := db.WithContext(ctx).Create(center).Error; err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 建立課程
	courseName := fmt.Sprintf("合併測試課程 - %d", time.Now().UnixNano())
	course := &models.Course{
		CenterID:         center.ID,
		Name:             courseName,
		DefaultDuration:  60,
		ColorHex:         "#3498db",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        time.Now(),
	}
	if err := db.WithContext(ctx).Create(course).Error; err != nil {
		db.Where("id = ?", center.ID).Delete(&models.Center{})
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 建立方案
	offeringName := fmt.Sprintf("合併測試方案 - %d", time.Now().UnixNano())
	offering := &models.Offering{
		CenterID:            center.ID,
		CourseID:            course.ID,
		Name:                offeringName,
		AllowBufferOverride: false,
		IsActive:            true,
		CreatedAt:           time.Now(),
	}
	if err := db.WithContext(ctx).Create(offering).Error; err != nil {
		db.Where("center_id = ?", center.ID).Delete(&models.Center{})
		db.Where("center_id = ?", center.ID).Delete(&models.Course{})
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 建立教室
	roomName := fmt.Sprintf("合併測試教室 - %d", time.Now().UnixNano())
	room := &models.Room{
		CenterID:  center.ID,
		Name:      roomName,
		Capacity:  10,
		CreatedAt: time.Now(),
	}
	if err := db.WithContext(ctx).Create(room).Error; err != nil {
		db.Where("center_id = ?", center.ID).Delete(&models.Center{})
		db.Where("center_id = ?", center.ID).Delete(&models.Course{})
		db.Where("center_id = ?", center.ID).Delete(&models.Offering{})
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 建立來源教師
	sourceLineID := fmt.Sprintf("line-source-%d", time.Now().UnixNano())
	sourceTeacher := &models.Teacher{
		LineUserID: sourceLineID,
		Name:       fmt.Sprintf("來源教師 - %d", time.Now().UnixNano()),
		Email:      fmt.Sprintf("source.%d@test.com", time.Now().UnixNano()),
		CreatedAt:  time.Now(),
	}
	if err := db.WithContext(ctx).Create(sourceTeacher).Error; err != nil {
		db.Where("center_id = ?", center.ID).Delete(&models.Center{})
		db.Where("center_id = ?", center.ID).Delete(&models.Course{})
		db.Where("center_id = ?", center.ID).Delete(&models.Offering{})
		db.Where("center_id = ?", center.ID).Delete(&models.Room{})
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 建立來源教師會籍
	db.WithContext(ctx).Create(&models.CenterMembership{
		CenterID:  center.ID,
		TeacherID: sourceTeacher.ID,
		Role:      "TEACHER",
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
	})

	// 建立目標教師
	targetLineID := fmt.Sprintf("line-target-%d", time.Now().UnixNano())
	targetTeacher := &models.Teacher{
		LineUserID: targetLineID,
		Name:       fmt.Sprintf("目標教師 - %d", time.Now().UnixNano()),
		Email:      fmt.Sprintf("target.%d@test.com", time.Now().UnixNano()),
		CreatedAt:  time.Now(),
	}
	if err := db.WithContext(ctx).Create(targetTeacher).Error; err != nil {
		db.Where("center_id = ?", center.ID).Delete(&models.Center{})
		db.Where("center_id = ?", center.ID).Delete(&models.Course{})
		db.Where("center_id = ?", center.ID).Delete(&models.Offering{})
		db.Where("center_id = ?", center.ID).Delete(&models.Room{})
		db.Where("id = ?", sourceTeacher.ID).Delete(&models.Teacher{})
		db.Where("center_id = ? AND teacher_id = ?", center.ID, sourceTeacher.ID).Delete(&models.CenterMembership{})
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 建立目標教師會籍
	db.WithContext(ctx).Create(&models.CenterMembership{
		CenterID:  center.ID,
		TeacherID: targetTeacher.ID,
		Role:      "TEACHER",
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
	})

	// 清理函數
	cleanup := func() {
		db.Where("center_id = ? AND teacher_id IN (?, ?)", center.ID, sourceTeacher.ID, targetTeacher.ID).Delete(&models.CenterMembership{})
		db.Where("center_id = ?", center.ID).Delete(&models.Center{})
		db.Where("center_id = ?", center.ID).Delete(&models.Course{})
		db.Where("center_id = ?", center.ID).Delete(&models.Offering{})
		db.Where("center_id = ?", center.ID).Delete(&models.Room{})
		db.Where("id IN (?, ?)", sourceTeacher.ID, targetTeacher.ID).Delete(&models.Teacher{})
	}

	return center, course, offering, room, sourceTeacher, targetTeacher, cleanup, nil
}

// TestTeacherMergeService_ValidationErrors 教師合併服務 - 驗證錯誤測試
func TestTeacherMergeService_ValidationErrors(t *testing.T) {
	// 初始化測試資料庫
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("跳過測試 - 資料庫連線失敗: %v", err)
		return
	}
	defer CloseDB(db)

	// 建立 App 實例
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// 建立合併服務
	mergeService := services.NewTeacherMergeService(appInstance)
	ctx := context.Background()

	t.Run("來源與目標教師相同時應失敗", func(t *testing.T) {
		err := mergeService.MergeTeacher(ctx, 1, 1, 1)
		assert.Error(t, err, "來源與目標相同時應失敗")
		assert.Contains(t, err.Error(), "相同", "錯誤訊息應包含 '相同'")
		t.Logf("預期錯誤: %v", err)
	})

	t.Run("來源教師 ID 為零時應失敗", func(t *testing.T) {
		err := mergeService.MergeTeacher(ctx, 0, 1, 1)
		assert.Error(t, err, "來源教師 ID 為零時應失敗")
		assert.Contains(t, err.Error(), "不能為零", "錯誤訊息應包含 '不能為零'")
		t.Logf("預期錯誤: %v", err)
	})

	t.Run("目標教師 ID 為零時應失敗", func(t *testing.T) {
		err := mergeService.MergeTeacher(ctx, 1, 0, 1)
		assert.Error(t, err, "目標教師 ID 為零時應失敗")
		assert.Contains(t, err.Error(), "不能為零", "錯誤訊息應包含 '不能為零'")
		t.Logf("預期錯誤: %v", err)
	})
}

// TestTeacherMergeService_NonExistentTeachers 教師合併服務 - 不存在教師測試
func TestTeacherMergeService_NonExistentTeachers(t *testing.T) {
	// 初始化測試資料庫
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("跳過測試 - 資料庫連線失敗: %v", err)
		return
	}
	defer CloseDB(db)

	// 建立 App 實例
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// 建立合併服務
	mergeService := services.NewTeacherMergeService(appInstance)
	ctx := context.Background()

	// 找一个存在的教师 ID
	var existingTeacher models.Teacher
	findErr := db.WithContext(ctx).Order("id DESC").First(&existingTeacher).Error
	if findErr != nil {
		t.Skipf("跳過測試 - 無可用教師資料: %v", findErr)
		return
	}

	nonExistentID := uint(99999)
	for nonExistentID > existingTeacher.ID {
		nonExistentID--
	}

	t.Run("來源教師不存在時應失敗", func(t *testing.T) {
		err := mergeService.MergeTeacher(ctx, nonExistentID, existingTeacher.ID, 1)
		assert.Error(t, err, "來源教師不存在時應失敗")
		assert.Contains(t, err.Error(), "不存在", "錯誤訊息應包含 '不存在'")
		t.Logf("預期錯誤: %v", err)
	})

	t.Run("目標教師不存在時應失敗", func(t *testing.T) {
		err := mergeService.MergeTeacher(ctx, existingTeacher.ID, nonExistentID, 1)
		assert.Error(t, err, "目標教師不存在時應失敗")
		assert.Contains(t, err.Error(), "不存在", "錯誤訊息應包含 '不存在'")
		t.Logf("預期錯誤: %v", err)
	})
}

// TestTeacherMergeService_MultipleScheduleRules 教師合併服務 - 多筆排課規則遷移測試
func TestTeacherMergeService_MultipleScheduleRules(t *testing.T) {
	// 初始化測試資料庫
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("跳過測試 - 資料庫連線失敗: %v", err)
		return
	}
	defer CloseDB(db)

	ctx := context.Background()

	// 設定測試資料
	center, course, offering, room, sourceTeacher, targetTeacher, cleanup, err := setupMergeTestData(ctx, db)
	if err != nil {
		t.Skipf("跳過測試 - 無法建立測試資料: %v", err)
		return
	}
	defer cleanup()
	_ = course // course 在 setupMergeTestData 內部建立，此處不使用

	// 建立 App 實例
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// 建立合併服務
	mergeService := services.NewTeacherMergeService(appInstance)

	// ============================================
	// 建立 5 筆不同時間的排課規則
	// ============================================

	baseTime := time.Now()
	ruleIDs := make([]uint, 5)

	for i := 0; i < 5; i++ {
		rule := &models.ScheduleRule{
			CenterID:       center.ID,
			OfferingID:     offering.ID,
			TeacherID:      &sourceTeacher.ID,
			RoomID:         room.ID,
			Name:           fmt.Sprintf("多規則測試%c", 'A'+i),
			Weekday:        1 + i,
			StartTime:      fmt.Sprintf("%02d:00", 9+i),
			EndTime:        fmt.Sprintf("%02d:00", 10+i),
			Duration:       60,
			IsCrossDay:     false,
			EffectiveRange: models.DateRange{
				StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			CreatedAt: baseTime,
			UpdatedAt: baseTime,
		}
		if err := db.WithContext(ctx).Create(rule).Error; err != nil {
			t.Fatalf("建立排課規則 %d 失敗: %v", i+1, err)
		}
		ruleIDs[i] = rule.ID
		t.Logf("建立排課規則 %d: ID=%d, Time=%02d:00-%02d:00", i+1, rule.ID, 9+i, 10+i)
	}

	// ============================================
	// 執行合併
	// ============================================

	err = mergeService.MergeTeacher(ctx, sourceTeacher.ID, targetTeacher.ID, center.ID)
	assert.NoError(t, err, "合併不應失敗")

	// ============================================
	// 驗證所有排課規則都已遷移
	// ============================================

	t.Run("驗證所有排課規則都已遷移", func(t *testing.T) {
		var rules []models.ScheduleRule
		db.WithContext(ctx).Where("id IN ?", ruleIDs).Find(&rules)

		assert.Equal(t, 5, len(rules), "應找到 5 筆排課規則")

		for i, rule := range rules {
			assert.NotNil(t, rule.TeacherID, "排課規則 %d TeacherID 不應為 nil", i+1)
			assert.Equal(t, targetTeacher.ID, *rule.TeacherID, "排課規則 %d TeacherID 應為目標教師 ID", i+1)
			t.Logf("排課規則 %d (ID=%d) teacher_id=%d ✓", i+1, rule.ID, *rule.TeacherID)
		}
	})

	// 驗證來源教師沒有任何排課規則
	t.Run("驗證來源教師沒有任何排課規則", func(t *testing.T) {
		var count int64
		db.WithContext(ctx).Model(&models.ScheduleRule{}).
			Where("teacher_id = ?", sourceTeacher.ID).
			Count(&count)
		assert.Equal(t, int64(0), count, "來源教師應沒有排課規則")
		t.Logf("來源教師排課規則數量: %d ✓", count)
	})

	// ============================================
	// 清理測試資料
	// ============================================

	t.Run("清理測試資料", func(t *testing.T) {
		db.Where("id IN ?", ruleIDs).Delete(&models.ScheduleRule{})
		db.Unscoped().Model(&sourceTeacher).Update("deleted_at", nil)
		t.Logf("測試資料已清理")
	})
}

// TestTeacherMergeService_MigratePersonalEvents 教師合併服務 - 私人行程遷移測試
func TestTeacherMergeService_MigratePersonalEvents(t *testing.T) {
	// 初始化測試資料庫
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("跳過測試 - 資料庫連線失敗: %v", err)
		return
	}
	defer CloseDB(db)

	ctx := context.Background()

	// 設定測試資料
	center, course, offering, room, sourceTeacher, targetTeacher, cleanup, err := setupMergeTestData(ctx, db)
	if err != nil {
		t.Skipf("跳過測試 - 無法建立測試資料: %v", err)
		return
	}
	defer cleanup()
	_ = course
	_ = offering
	_ = room

	// 建立 App 實例
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// 建立合併服務
	mergeService := services.NewTeacherMergeService(appInstance)

	loc := time.UTC

	// 建立私人行程
	event1 := &models.PersonalEvent{
		TeacherID: sourceTeacher.ID,
		Title:     "醫療門診",
		IsAllDay:  false,
		StartAt:   time.Date(2026, 2, 1, 9, 0, 0, 0, loc),
		EndAt:     time.Date(2026, 2, 1, 10, 30, 0, 0, loc),
		Note:      "私人行程",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.WithContext(ctx).Create(event1).Error; err != nil {
		t.Fatalf("建立私人行程1失敗: %v", err)
	}

	event2 := &models.PersonalEvent{
		TeacherID: sourceTeacher.ID,
		Title:     "家族聚餐",
		IsAllDay:  true,
		StartAt:   time.Date(2026, 2, 15, 0, 0, 0, 0, loc),
		EndAt:     time.Date(2026, 2, 15, 23, 59, 59, 0, loc),
		Note:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.WithContext(ctx).Create(event2).Error; err != nil {
		t.Fatalf("建立私人行程2失敗: %v", err)
	}

	t.Logf("建立來源教師私人行程: %d, %d", event1.ID, event2.ID)

	// ============================================
	// 執行合併
	// ============================================

	err = mergeService.MergeTeacher(ctx, sourceTeacher.ID, targetTeacher.ID, center.ID)
	assert.NoError(t, err, "合併不應失敗")

	// ============================================
	// 驗證私人行程已遷移
	// ============================================

	t.Run("驗證私人行程已遷移至目標教師", func(t *testing.T) {
		var events []models.PersonalEvent
		db.WithContext(ctx).Where("teacher_id = ?", targetTeacher.ID).Find(&events)

		assert.Equal(t, 2, len(events), "目標教師應有 2 筆私人行程")
		t.Logf("目標教師私人行程數量: %d", len(events))

		// 驗證行程標題
		titles := make(map[string]bool)
		for _, event := range events {
			titles[event.Title] = true
		}
		assert.True(t, titles["醫療門診"], "應包含醫療門診行程")
		assert.True(t, titles["家族聚餐"], "應包含家族聚餐行程")
	})

	// 驗證來源教師沒有私人行程
	t.Run("驗證來源教師沒有私人行程", func(t *testing.T) {
		var count int64
		db.WithContext(ctx).Model(&models.PersonalEvent{}).
			Where("teacher_id = ?", sourceTeacher.ID).
			Count(&count)
		assert.Equal(t, int64(0), count, "來源教師應沒有私人行程")
	})

	// ============================================
	// 清理測試資料
	// ============================================

	t.Run("清理測試資料", func(t *testing.T) {
		db.Where("id IN ?", []uint{event1.ID, event2.ID}).Delete(&models.PersonalEvent{})
		db.Unscoped().Model(&sourceTeacher).Update("deleted_at", nil)
		t.Logf("測試資料已清理")
	})
}

// TestTeacherMergeService_MigrateTeacherSkills 教師合併服務 - 技能遷移測試
func TestTeacherMergeService_MigrateTeacherSkills(t *testing.T) {
	// 初始化測試資料庫
	db, err := InitializeTestDB()
	if err != nil {
		t.Skipf("跳過測試 - 資料庫連線失敗: %v", err)
		return
	}
	defer CloseDB(db)

	ctx := context.Background()

	// 設定測試資料
	center, course, offering, room, sourceTeacher, targetTeacher, cleanup, err := setupMergeTestData(ctx, db)
	if err != nil {
		t.Skipf("跳過測試 - 無法建立測試資料: %v", err)
		return
	}
	defer cleanup()
	_ = course
	_ = offering
	_ = room

	// 建立 App 實例
	appInstance := &app.App{
		MySQL: &mysql.DB{WDB: db, RDB: db},
	}

	// 建立合併服務
	mergeService := services.NewTeacherMergeService(appInstance)

	// 建立來源教師的技能
	skill1 := &models.TeacherSkill{
		TeacherID: sourceTeacher.ID,
		Category:  "運動",
		SkillName: "瑜珈",
		Level:     "初級",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.WithContext(ctx).Create(skill1)

	skill2 := &models.TeacherSkill{
		TeacherID: sourceTeacher.ID,
		Category:  "運動",
		SkillName: "皮拉提斯",
		Level:     "中級",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.WithContext(ctx).Create(skill2)

	t.Logf("建立來源教師技能: %s(%s), %s(%s)", skill1.SkillName, skill1.Level, skill2.SkillName, skill2.Level)

	// ============================================
	// 執行合併
	// ============================================

	err = mergeService.MergeTeacher(ctx, sourceTeacher.ID, targetTeacher.ID, center.ID)
	assert.NoError(t, err, "合併不應失敗")

	// ============================================
	// 驗證技能已遷移
	// ============================================

	t.Run("驗證技能已遷移至目標教師", func(t *testing.T) {
		var skills []models.TeacherSkill
		db.WithContext(ctx).Where("teacher_id = ?", targetTeacher.ID).Find(&skills)

		assert.Equal(t, 2, len(skills), "目標教師應有 2 筆技能")
		t.Logf("目標教師技能數量: %d", len(skills))

		// 驗證技能名稱
		skillNames := make(map[string]bool)
		for _, skill := range skills {
			skillNames[skill.SkillName] = true
		}
		assert.True(t, skillNames["瑜珈"], "應包含瑜珈技能")
		assert.True(t, skillNames["皮拉提斯"], "應包含皮拉提斯技能")
	})

	// 驗證來源教師沒有技能
	t.Run("驗證來源教師沒有技能", func(t *testing.T) {
		var count int64
		db.WithContext(ctx).Model(&models.TeacherSkill{}).
			Where("teacher_id = ?", sourceTeacher.ID).
			Count(&count)
		assert.Equal(t, int64(0), count, "來源教師應沒有技能")
	})

	// ============================================
	// 清理測試資料
	// ============================================

	t.Run("清理測試資料", func(t *testing.T) {
		db.Where("teacher_id = ?", []uint{sourceTeacher.ID, targetTeacher.ID}).
			Where("skill_name IN ?", []string{"瑜珈", "皮拉提斯"}).
			Delete(&models.TeacherSkill{})
		db.Unscoped().Model(&sourceTeacher).Update("deleted_at", nil)
		t.Logf("測試資料已清理")
	})
}
