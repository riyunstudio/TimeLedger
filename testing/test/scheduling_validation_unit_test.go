package test

import (
	"context"
	"testing"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/database/mysql"

	"github.com/stretchr/testify/assert"
)

// TestScheduleValidationService_Complete 排課驗證服務完整測試
func TestScheduleValidationService_Complete(t *testing.T) {
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

	// 建立驗證服務
	validationService := services.NewScheduleValidationService(appInstance)

	// 取得測試資料
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("跳過測試 - 無可用中心資料: %v", err)
		return
	}

	var teacher models.Teacher
	if err := db.Order("id ASC").First(&teacher).Error; err != nil {
		t.Skipf("跳過測試 - 無可用老師資料: %v", err)
		return
	}

	var offering models.Offering
	if err := db.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("跳過測試 - 無可用課程資料: %v", err)
		return
	}

	teacherID := teacher.ID
	loc := app.GetTaiwanLocation()

	// ============================================
	// CheckOverlap 測試案例
	// ============================================

	t.Run("CheckOverlap_NoConflict", func(t *testing.T) {
		// 測試時段：無衝突
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc) // Monday 14:00
		endTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc)   // Monday 15:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("無衝突測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("CheckOverlap_TeacherOverlap", func(t *testing.T) {
		// 測試時段：老師時間重疊
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc) // Monday 08:00
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)   // Monday 09:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("老師重疊測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
		for _, conflict := range result.Conflicts {
			t.Logf("  衝突: Type=%s, Message=%s", conflict.Type, conflict.Message)
		}
	})

	t.Run("CheckOverlap_RoomOverlap", func(t *testing.T) {
		// 測試時段：教室重疊（nil teacherID）
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			nil,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("教室重疊測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("CheckOverlap_SundayWeekday", func(t *testing.T) {
		// 測試時段：週日轉換（Go weekday=0, DB weekday=7）
		startTime := time.Date(2026, 1, 26, 14, 0, 0, 0, loc) // Sunday 14:00
		endTime := time.Date(2026, 1, 26, 15, 0, 0, 0, loc)   // Sunday 15:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			0,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("週日轉換測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	// ============================================
	// CheckTeacherBuffer 測試案例
	// ============================================

	t.Run("CheckTeacherBuffer_Sufficient", func(t *testing.T) {
		// 測試緩衝：足夠的老師緩衝時間
		prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)
		nextStartTime := time.Date(2026, 1, 27, 11, 30, 0, 0, loc)

		result, err := validationService.CheckTeacherBuffer(
			context.Background(),
			center.ID,
			teacherID,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("老師足夠緩衝測試結果: Valid=%v", result.Valid)
	})

	t.Run("CheckTeacherBuffer_Insufficient", func(t *testing.T) {
		// 測試緩衝：不足的老師緩衝時間
		prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)
		nextStartTime := time.Date(2026, 1, 27, 10, 30, 0, 0, loc)

		result, err := validationService.CheckTeacherBuffer(
			context.Background(),
			center.ID,
			teacherID,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("老師不足緩衝測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
		for _, conflict := range result.Conflicts {
			t.Logf("  衝突: Type=%s, Message=%s, CanOverride=%v",
				conflict.Type, conflict.Message, conflict.CanOverride)
		}
	})

	// ============================================
	// CheckRoomBuffer 測試案例
	// ============================================

	t.Run("CheckRoomBuffer_Sufficient", func(t *testing.T) {
		// 測試緩衝：足夠的教室緩衝時間
		prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)
		nextStartTime := time.Date(2026, 1, 27, 10, 30, 0, 0, loc)

		result, err := validationService.CheckRoomBuffer(
			context.Background(),
			center.ID,
			1,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("教室足夠緩衝測試結果: Valid=%v", result.Valid)
	})

	t.Run("CheckRoomBuffer_Insufficient", func(t *testing.T) {
		// 測試緩衝：不足的教室緩衝時間
		prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)
		nextStartTime := time.Date(2026, 1, 27, 10, 15, 0, 0, loc)

		result, err := validationService.CheckRoomBuffer(
			context.Background(),
			center.ID,
			1,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("教室不足緩衝測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	// ============================================
	// ValidateFull 測試案例
	// ============================================

	t.Run("ValidateFull_NoConflicts", func(t *testing.T) {
		// 測試完整驗證：無衝突
		startTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 16, 0, 0, 0, loc)

		result, err := validationService.ValidateFull(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			offering.ID,
			startTime,
			endTime,
			nil,
			false,
			nil,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("完整驗證無衝突測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("ValidateFull_WithOverlap", func(t *testing.T) {
		// 測試完整驗證：有重疊
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)

		result, err := validationService.ValidateFull(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			offering.ID,
			startTime,
			endTime,
			nil,
			false,
			nil,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("完整驗證重疊測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))

		for _, conflict := range result.Conflicts {
			if conflict.Type == "TEACHER_OVERLAP" || conflict.Type == "ROOM_OVERLAP" {
				assert.False(t, conflict.CanOverride, "重疊衝突不應可覆寫")
			}
		}
	})

	t.Run("ValidateFull_WithBufferOverride", func(t *testing.T) {
		// 測試完整驗證：允許緩衝覆寫
		startTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)

		result, err := validationService.ValidateFull(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			offering.ID,
			startTime,
			endTime,
			nil,
			true,
			nil,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("完整驗證緩衝覆寫測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))

		for _, conflict := range result.Conflicts {
			if conflict.Type == "TEACHER_BUFFER" || conflict.Type == "ROOM_BUFFER" {
				assert.True(t, conflict.CanOverride, "緩衝衝突應可覆寫")
			}
		}
	})

	t.Run("ValidateFull_NilTeacher", func(t *testing.T) {
		// 測試完整驗證：無老師ID（只檢查教室）
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc)

		result, err := validationService.ValidateFull(
			context.Background(),
			center.ID,
			nil,
			1,
			offering.ID,
			startTime,
			endTime,
			nil,
			false,
			nil,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("完整驗證無老師ID測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))

		for _, conflict := range result.Conflicts {
			assert.NotEqual(t, "TEACHER_OVERLAP", conflict.Type)
			assert.NotEqual(t, "TEACHER_BUFFER", conflict.Type)
		}
	})

	// ============================================
	// 邊界情況測試案例
	// ============================================

	t.Run("Boundary_BackToBackCourses", func(t *testing.T) {
		// 邊界：背靠背課程（時間剛好銜接）
		prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)
		nextStartTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)

		result1, err1 := validationService.CheckTeacherBuffer(
			context.Background(),
			center.ID,
			teacherID,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)
		assert.NoError(t, err1)
		t.Logf("背靠背老師緩衝: Valid=%v", result1.Valid)

		result2, err2 := validationService.CheckRoomBuffer(
			context.Background(),
			center.ID,
			1,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)
		assert.NoError(t, err2)
		t.Logf("背靠背教室緩衝: Valid=%v", result2.Valid)
	})

	t.Run("Boundary_CrossDayCourse", func(t *testing.T) {
		// 邊界：跨日課程（跨越午夜）
		startTime := time.Date(2026, 1, 27, 23, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 28, 2, 0, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			5,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("跨日課程測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_EarlyMorning", func(t *testing.T) {
		// 邊界：凌晨課程（00:00）
		startTime := time.Date(2026, 1, 27, 0, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 1, 0, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("凌晨課程測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_LateNight", func(t *testing.T) {
		// 邊界：深夜課程（23:59）
		startTime := time.Date(2026, 1, 27, 23, 30, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 23, 59, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("深夜課程測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_SaturdayNormalDay", func(t *testing.T) {
		// 邊界：週六 - 通常為正常上課日（非假日）
		// 週六不在 center_holiday 中定義時，可正常排課
		startTime := time.Date(2026, 1, 31, 10, 0, 0, 0, loc) // Saturday 10:00
		endTime := time.Date(2026, 1, 31, 11, 0, 0, 0, loc)   // Saturday 11:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			6, // Saturday in database (Monday=1, Saturday=6)
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("週六正常上課日測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_SundayHoliday", func(t *testing.T) {
		// 邊界：週日 - 通常為例假日
		// 週日是否停課取決於 center_holiday 表設定
		// 若 2026/2/1 在 center_holiday 中，則為假日；否則可能正常上課
		startTime := time.Date(2026, 2, 1, 10, 0, 0, 0, loc) // Sunday 10:00
		endTime := time.Date(2026, 2, 1, 11, 0, 0, 0, loc)   // Sunday 11:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			0, // Sunday = 0 in Go (will be converted to 7 for database)
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("週日邊界測試結果: Valid=%v, Conflicts=%d (若 2026/2/1 在 center_holiday 中則為假日)", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_ExactMinute", func(t *testing.T) {
		// 邊界：精確分鐘（10:05, 10:10）
		startTime := time.Date(2026, 1, 27, 10, 5, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 10, 10, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("精確分鐘測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_HourBoundary", func(t *testing.T) {
		// 邊界：小時邊界（09:00, 10:00）
		startTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("小時邊界測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_LongDuration", func(t *testing.T) {
		// 邊界：長時段課程（3小時）
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 17, 0, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("長時段課程測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_ShortDuration", func(t *testing.T) {
		// 邊界：短時段課程（15分鐘）
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 14, 15, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("短時段課程測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_SingleMinute", func(t *testing.T) {
		// 邊界：單分鐘課程
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 14, 1, 0, 0, loc)

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("單分鐘課程測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})
}

// TestValidationResult_Structure 測試 ValidationResult 結構體
func TestValidationResult_Structure(t *testing.T) {
	result := services.ValidationResult{
		Valid: false,
		Conflicts: []services.ValidationConflict{
			{
				Type:             "TEACHER_OVERLAP",
				Message:          "老師在該時段已有課程安排",
				CanOverride:      false,
				ConflictSource:   "RULE",
				ConflictSourceID: 1,
				Details:          "rule_id:1, offering_id:1",
			},
			{
				Type:            "TEACHER_BUFFER",
				Message:         "老師轉場時間不足，需間隔 15 分鐘，實際間隔 5 分鐘",
				CanOverride:     true,
				RequiredMinutes: 15,
				DiffMinutes:     10,
				ConflictSource:  "SESSION",
			},
		},
	}

	assert.False(t, result.Valid)
	assert.Len(t, result.Conflicts, 2)

	// 檢查第一個衝突（重疊 - 不可覆寫）
	assert.Equal(t, "TEACHER_OVERLAP", result.Conflicts[0].Type)
	assert.False(t, result.Conflicts[0].CanOverride)
	assert.Equal(t, "RULE", result.Conflicts[0].ConflictSource)

	// 檢查第二個衝突（緩衝 - 可覆寫）
	assert.Equal(t, "TEACHER_BUFFER", result.Conflicts[1].Type)
	assert.True(t, result.Conflicts[1].CanOverride)
	assert.Equal(t, 15, result.Conflicts[1].RequiredMinutes)
	assert.Equal(t, 10, result.Conflicts[1].DiffMinutes)
}

// TestScheduleValidationService_AdditionalCases 補齊測試案例
func TestScheduleValidationService_AdditionalCases(t *testing.T) {
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

	// 建立驗證服務
	validationService := services.NewScheduleValidationService(appInstance)

	// 取得測試資料
	var center models.Center
	if err := db.Order("id DESC").First(&center).Error; err != nil {
		t.Skipf("跳過測試 - 無可用中心資料: %v", err)
		return
	}

	var teacher models.Teacher
	if err := db.Where("user_type = ?", "TEACHER").Order("id ASC").First(&teacher).Error; err != nil {
		t.Skipf("跳過測試 - 無可用老師資料: %v", err)
		return
	}

	var offering models.Offering
	if err := db.Where("center_id = ?", center.ID).Order("id DESC").First(&offering).Error; err != nil {
		t.Skipf("跳過測試 - 無可用課程資料: %v", err)
		return
	}

	teacherID := teacher.ID
	loc := app.GetTaiwanLocation()

	// ============================================
	// CheckOverlap 補齊測試案例
	// ============================================

	t.Run("CheckOverlap_BothTeacherAndRoomOverlap", func(t *testing.T) {
		// 測試時段：老師和教室同時重疊
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc) // Monday 08:00
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)   // Monday 09:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("老師和教室重疊測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))

		// 應該有兩個衝突：TEACHER_OVERLAP 和 ROOM_OVERLAP
		hasTeacherOverlap := false
		hasRoomOverlap := false
		for _, conflict := range result.Conflicts {
			if conflict.Type == "TEACHER_OVERLAP" {
				hasTeacherOverlap = true
			}
			if conflict.Type == "ROOM_OVERLAP" {
				hasRoomOverlap = true
			}
		}
		assert.True(t, hasTeacherOverlap, "應該有老師重疊衝突")
		assert.True(t, hasRoomOverlap, "應該有教室重疊衝突")
	})

	t.Run("CheckOverlap_WithExcludeRuleID", func(t *testing.T) {
		// 測試時段：排除特定規則
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc) // Monday 14:00
		endTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc)   // Monday 15:00

		// 排除不存在的規則 ID（應該正常運行）
		excludeID := uint(99999)
		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			&excludeID,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("排除規則測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("CheckOverlap_EmptyTeacherIDs", func(t *testing.T) {
		// 測試時段：空的 teacherIDs 列表
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc)
		endTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc)

		var emptyTeacherID *uint = nil
		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			emptyTeacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("空老師ID測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	// ============================================
	// CheckTeacherBuffer 補齊測試案例
	// ============================================

	t.Run("CheckTeacherBuffer_ZeroBuffer", func(t *testing.T) {
		// 測試緩衝：緩衝時間為 0（應該直接通過）
		prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)
		nextStartTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc) // 剛好銜接

		// 這個測試需要一個緩衝時間為 0 的課程
		// 我們使用現有的 offering，但如果它的緩衝不是 0，可能會失敗
		result, err := validationService.CheckTeacherBuffer(
			context.Background(),
			center.ID,
			teacherID,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("零緩衝測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	// ============================================
	// CheckRoomBuffer 補齊測試案例
	// ============================================

	t.Run("CheckRoomBuffer_ZeroBuffer", func(t *testing.T) {
		// 測試緩衝：緩衝時間為 0
		prevEndTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)
		nextStartTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)

		result, err := validationService.CheckRoomBuffer(
			context.Background(),
			center.ID,
			1,
			prevEndTime,
			nextStartTime,
			offering.ID,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("教室零緩衝測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	// ============================================
	// ValidateFull 補齊測試案例
	// ============================================

	t.Run("ValidateFull_NoPreviousSession", func(t *testing.T) {
		// 測試完整驗證：沒有前一堂課（早上第一堂課）
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc) // Monday 08:00
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)   // Monday 09:00

		result, err := validationService.ValidateFull(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			offering.ID,
			startTime,
			endTime,
			nil,
			false,
			nil,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("無前一堂課測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("ValidateFull_NoPreviousRoomSession", func(t *testing.T) {
		// 測試完整驗證：教室沒有前一堂課
		startTime := time.Date(2026, 1, 27, 14, 0, 0, 0, loc) // Monday 14:00
		endTime := time.Date(2026, 1, 27, 15, 0, 0, 0, loc)   // Monday 15:00

		// 使用一個在測試資料中沒有使用的教室 ID
		result, err := validationService.ValidateFull(
			context.Background(),
			center.ID,
			&teacherID,
			999, // 不存在的教室 ID
			offering.ID,
			startTime,
			endTime,
			nil,
			false,
			nil,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("教室無前一堂課測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("ValidateFull_BothOverlapAndBuffer", func(t *testing.T) {
		// 測試完整驗證：同時有重疊和緩衝問題
		// 08:00-09:00 有課程，測試 08:30-09:30
		startTime := time.Date(2026, 1, 27, 8, 30, 0, 0, loc) // Monday 08:30
		endTime := time.Date(2026, 1, 27, 9, 30, 0, 0, loc)   // Monday 09:30

		result, err := validationService.ValidateFull(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			offering.ID,
			startTime,
			endTime,
			nil,
			false,
			nil,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("同時有重疊和緩衝測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))

		// 重疊衝突不應該可覆寫
		for _, conflict := range result.Conflicts {
			if conflict.Type == "TEACHER_OVERLAP" || conflict.Type == "ROOM_OVERLAP" {
				assert.False(t, conflict.CanOverride, "重疊衝突不應可覆寫")
			}
		}
	})

	// ============================================
	// 邊界情況補齊測試
	// ============================================

	t.Run("Boundary_ExactOverlap", func(t *testing.T) {
		// 邊界：完全重疊（同一時間區間）
		// 假設 08:00-09:00 有課程
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc) // Monday 08:00
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)   // Monday 09:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("完全重疊測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
		assert.False(t, result.Valid, "完全重疊應該無效")
	})

	t.Run("Boundary_AdjacentButNotOverlapping", func(t *testing.T) {
		// 邊界：相鄰但不重疊（09:00-10:00 與 08:00-09:00）
		startTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc) // Monday 09:00
		endTime := time.Date(2026, 1, 27, 10, 0, 0, 0, loc)  // Monday 10:00

		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("相鄰不重疊測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))
	})

	t.Run("Boundary_DifferentRooms", func(t *testing.T) {
		// 邊界：不同教室不應該有 room overlap
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc) // Monday 08:00
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)   // Monday 09:00

		// 使用不同的教室 ID（假設 room_id=2）
		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&teacherID,
			2, // 不同的教室
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("不同教室測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))

		// 應該沒有 room overlap
		for _, conflict := range result.Conflicts {
			assert.NotEqual(t, "ROOM_OVERLAP", conflict.Type, "不同教室不應該有 room overlap")
		}
	})

	t.Run("Boundary_DifferentTeachers", func(t *testing.T) {
		// 邊界：不同老師不應該有 teacher overlap
		startTime := time.Date(2026, 1, 27, 8, 0, 0, 0, loc) // Monday 08:00
		endTime := time.Date(2026, 1, 27, 9, 0, 0, 0, loc)   // Monday 09:00

		// 使用不同的老師 ID（假設 teacher_id=99999）
		var differentTeacherID uint = 99999
		result, err := validationService.CheckOverlap(
			context.Background(),
			center.ID,
			&differentTeacherID,
			1,
			startTime,
			endTime,
			1,
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		t.Logf("不同老師測試結果: Valid=%v, Conflicts=%d", result.Valid, len(result.Conflicts))

		// 應該沒有 teacher overlap
		for _, conflict := range result.Conflicts {
			assert.NotEqual(t, "TEACHER_OVERLAP", conflict.Type, "不同老師不應該有 teacher overlap")
		}
	})
}
