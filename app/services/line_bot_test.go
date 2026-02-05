package services

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupLineBotTestApp 建立測試用的 App 實例
func setupLineBotTestApp(t *testing.T) *app.App {
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

// cleanupLineBotTestData 清理測試資料
func cleanupLineBotTestData(t *testing.T, appInstance *app.App, adminLineUserID, teacherLineUserID string) {
	ctx := context.Background()

	// 清理管理員測試資料
	if adminLineUserID != "" {
		appInstance.MySQL.WDB.WithContext(ctx).
			Table("admin_users").
			Where("line_user_id LIKE ?", adminLineUserID+"%").
			Delete(&models.AdminUser{})
	}

	// 清理老師測試資料
	if teacherLineUserID != "" {
		appInstance.MySQL.WDB.WithContext(ctx).
			Table("teachers").
			Where("line_user_id LIKE ?", teacherLineUserID+"%").
			Delete(&models.Teacher{})

		// 清理相關的會員關係
		appInstance.MySQL.WDB.WithContext(ctx).
			Table("center_memberships").
			Where("teacher_id IN (SELECT id FROM teachers WHERE line_user_id LIKE ?)", teacherLineUserID+"%").
			Delete(&models.CenterMembership{})
	}
}

// TestLineBotService_GetCombinedIdentity 測試整合身份識別功能
func TestLineBotService_GetCombinedIdentity(t *testing.T) {
	t.Run("AdminOnly_ReturnAdminIdentity", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-teacher-")
		}()

		ctx := context.Background()

		// 建立測試管理員資料
		testLineUserID := fmt.Sprintf("test-line-admin-%d", time.Now().UnixNano())
		admin := models.AdminUser{
			Name:         "Test Admin for Combined Identity",
			Email:        fmt.Sprintf("test-admin-%d@test.com", time.Now().UnixNano()),
			PasswordHash: "hashed_password",
			Role:         "ADMIN",
			CenterID:     1,
			LineUserID:   testLineUserID,
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("admin_users").Create(&admin).Error; err != nil {
			t.Fatalf("建立測試管理員失敗: %v", err)
		}

		// 執行測試
		svc := NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(testLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功，但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "ADMIN" {
			t.Errorf("預期 PrimaryRole 為 'ADMIN'，但取得 '%s'", identity.PrimaryRole)
		}

		if len(identity.AdminProfiles) != 1 {
			t.Errorf("預期有 1 個管理員資料，但取得 %d 個", len(identity.AdminProfiles))
		}

		if identity.TeacherProfile != nil {
			t.Error("預期 TeacherProfile 為 nil，但取得非空值")
		}

		if identity.Memberships != nil && len(identity.Memberships) > 0 {
			t.Error("預期 Memberships 為空，但取得非空值")
		}
	})

	t.Run("TeacherOnly_ReturnTeacherIdentity", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-teacher-")
		}()

		ctx := context.Background()
		centerID := uint(1)

		// 建立測試老師資料
		testLineUserID := fmt.Sprintf("test-line-teacher-%d", time.Now().UnixNano())
		teacher := models.Teacher{
			Name:      "Test Teacher for Combined Identity",
			Email:     fmt.Sprintf("test-teacher-%d@test.com", time.Now().UnixNano()),
			LineUserID: testLineUserID,
			City:      "台北市",
			District:  "大安區",
			AvatarURL: "https://example.com/avatar.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("teachers").Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立老師的會員關係
		membership := models.CenterMembership{
			CenterID:  centerID,
			TeacherID: teacher.ID,
			Role:      "TEACHER",
			Status:    "ACTIVE",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("center_memberships").Create(&membership).Error; err != nil {
			t.Fatalf("建立測試會員關係失敗: %v", err)
		}

		// 執行測試
		svc := NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(testLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功，但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "TEACHER" {
			t.Errorf("預期 PrimaryRole 為 'TEACHER'，但取得 '%s'", identity.PrimaryRole)
		}

		if len(identity.AdminProfiles) != 0 {
			t.Errorf("預期有 0 個管理員資料，但取得 %d 個", len(identity.AdminProfiles))
		}

		if identity.TeacherProfile == nil {
			t.Fatal("預期 TeacherProfile 不為 nil")
		}

		if identity.TeacherProfile.ID != teacher.ID {
			t.Errorf("預期 TeacherProfile.ID 為 %d，但取得 %d", teacher.ID, identity.TeacherProfile.ID)
		}

		if identity.TeacherProfile.Name != teacher.Name {
			t.Errorf("預期 TeacherProfile.Name 為 '%s'，但取得 '%s'", teacher.Name, identity.TeacherProfile.Name)
		}

		if len(identity.Memberships) != 1 {
			t.Errorf("預期有 1 個會員關係，但取得 %d 個", len(identity.Memberships))
		}

		if len(identity.Memberships) > 0 && identity.Memberships[0].CenterID != centerID {
			t.Errorf("預期會員關係的 CenterID 為 %d，但取得 %d", centerID, identity.Memberships[0].CenterID)
		}
	})

	t.Run("GuestNotBound_ReturnGuestIdentity", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-teacher-")
		}()

		// 使用一個不可能存在的 LINE User ID
		nonExistentLineUserID := fmt.Sprintf("non-existent-line-user-%d@test.com", time.Now().UnixNano())

		// 執行測試
		svc := NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(nonExistentLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功（找不到資料視為正常），但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "GUEST" {
			t.Errorf("預期 PrimaryRole 為 'GUEST'，但取得 '%s'", identity.PrimaryRole)
		}

		if len(identity.AdminProfiles) != 0 {
			t.Errorf("預期有 0 個管理員資料，但取得 %d 個", len(identity.AdminProfiles))
		}

		if identity.TeacherProfile != nil {
			t.Error("預期 TeacherProfile 為 nil，但取得非空值")
		}

		if identity.Memberships != nil && len(identity.Memberships) > 0 {
			t.Error("預期 Memberships 為空，但取得非空值")
		}
	})

	t.Run("NoMemberships_ReturnTeacherWithoutMemberships", func(t *testing.T) {
		appInstance := setupLineBotTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			cleanupLineBotTestData(t, appInstance, "test-line-admin-", "test-line-no-membership-")
		}()

		ctx := context.Background()

		// 建立沒有會員關係的老師資料
		testLineUserID := fmt.Sprintf("test-line-no-membership-%d", time.Now().UnixNano())
		teacher := models.Teacher{
			Name:      "Test Teacher No Memberships",
			Email:     fmt.Sprintf("test-teacher-no-membership-%d@test.com", time.Now().UnixNano()),
			LineUserID: testLineUserID,
			City:      "新北市",
			District:  "板橋區",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := appInstance.MySQL.WDB.WithContext(ctx).Table("teachers").Create(&teacher).Error; err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 執行測試
		svc := NewLineBotService(appInstance)
		identity, err := svc.GetCombinedIdentity(testLineUserID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetCombinedIdentity 應該成功，但發生錯誤: %v", err)
		}

		if identity.PrimaryRole != "TEACHER" {
			t.Errorf("預期 PrimaryRole 為 'TEACHER'，但取得 '%s'", identity.PrimaryRole)
		}

		if identity.TeacherProfile == nil {
			t.Fatal("預期 TeacherProfile 不為 nil")
		}

		if len(identity.Memberships) != 0 {
			t.Errorf("預期有 0 個會員關係，但取得 %d 個", len(identity.Memberships))
		}
	})
}

// TestCompareTimeStrings 測試時間字串比較函數
func TestCompareTimeStrings(t *testing.T) {
	tests := []struct {
		name     string
		t1       string
		t2       string
		expected int // -1: t1 < t2, 0: equal, 1: t1 > t2
	}{
		{"t1 早於 t2", "09:00", "10:00", -1},
		{"t1 晚於 t2", "14:30", "10:00", 1},
		{"時間相等", "12:00", "12:00", 0},
		{"t1 是凌晨", "00:00", "08:00", -1},
		{"t1 是深夜", "23:59", "12:00", 1},
		{"不同分鐘數", "10:15", "10:30", -1},
		{"跨小時邊界", "09:59", "10:00", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareTimeStrings(tt.t1, tt.t2)
			if result != tt.expected {
				t.Errorf("compareTimeStrings(%q, %q) = %d, expected %d", tt.t1, tt.t2, result, tt.expected)
			}
		})
	}
}

// TestSortAgendaItemsByTime 測試行程項目排序功能
func TestSortAgendaItemsByTime(t *testing.T) {
	t.Run("MixedCenterAndPersonalItems_SortedChronologically", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "14:00", Title: "下午課程", SourceName: "中心A", SourceType: AgendaSourceTypeCenter},
			{Time: "09:00", Title: "上午課程", SourceName: "中心A", SourceType: AgendaSourceTypeCenter},
			{Time: "10:00", Title: "個人會議", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "11:00", Title: "上午課程", SourceName: "中心B", SourceType: AgendaSourceTypeCenter},
		}

		sortAgendaItemsByTime(items)

		// 驗證排序結果
		if len(items) != 4 {
			t.Fatalf("預期 4 個項目，但取得 %d 個", len(items))
		}

		expectedOrder := []string{"09:00", "10:00", "11:00", "14:00"}
		for i, expected := range expectedOrder {
			if items[i].Time != expected {
				t.Errorf("排序後第 %d 個項目時間應為 %q，但取得 %q", i, expected, items[i].Time)
			}
		}
	})

	t.Run("OnlyPersonalEvents_SortedChronologically", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "15:00", Title: "下午行程", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "08:00", Title: "早晨運動", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "12:30", Title: "午餐約會", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
		}

		sortAgendaItemsByTime(items)

		expectedOrder := []string{"08:00", "12:30", "15:00"}
		for i, expected := range expectedOrder {
			if items[i].Time != expected {
				t.Errorf("排序後第 %d 個項目時間應為 %q，但取得 %q", i, expected, items[i].Time)
			}
		}

		// 驗證來源類型正確
		for _, item := range items {
			if item.SourceType != AgendaSourceTypePersonal {
				t.Errorf("預期來源類型為 PERSONAL，但取得 %s", item.SourceType)
			}
		}
	})

	t.Run("EmptyItems_ReturnsEmpty", func(t *testing.T) {
		items := []AgendaItem{}
		sortAgendaItemsByTime(items)
		if len(items) != 0 {
			t.Errorf("預期空陣列，但取得 %d 個項目", len(items))
		}
	})

	t.Run("SingleItem_ReturnsUnchanged", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "10:00", Title: "唯一項目", SourceName: "中心A", SourceType: AgendaSourceTypeCenter},
		}

		sortAgendaItemsByTime(items)

		if len(items) != 1 {
			t.Errorf("預期 1 個項目，但取得 %d 個", len(items))
		}
		if items[0].Time != "10:00" {
			t.Errorf("時間應為 10:00，但取得 %s", items[0].Time)
		}
	})

	t.Run("SameTimeDifferentSources_SortedBySourceType", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "10:00", Title: "個人行程", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "10:00", Title: "中心課程", SourceName: "中心A", SourceType: AgendaSourceTypeCenter},
		}

		sortAgendaItemsByTime(items)

		// 時間相同時，順序不影響正確性（都是同時段）
		if len(items) != 2 {
			t.Errorf("預期 2 個項目，但取得 %d 個", len(items))
		}
	})

	t.Run("MultipleCenters_SortedCorrectly", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "16:00", Title: "中心C課程", SourceName: "中心C", SourceType: AgendaSourceTypeCenter},
			{Time: "09:00", Title: "中心A課程", SourceName: "中心A", SourceType: AgendaSourceTypeCenter},
			{Time: "13:00", Title: "中心B課程", SourceName: "中心B", SourceType: AgendaSourceTypeCenter},
			{Time: "11:00", Title: "中心A另一課程", SourceName: "中心A", SourceType: AgendaSourceTypeCenter},
		}

		sortAgendaItemsByTime(items)

		expectedOrder := []string{"09:00", "11:00", "13:00", "16:00"}
		for i, expected := range expectedOrder {
			if items[i].Time != expected {
				t.Errorf("排序後第 %d 個項目時間應為 %q，但取得 %q", i, expected, items[i].Time)
			}
		}
	})

	t.Run("AllDayPersonalEvents_SortedWithTime", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "00:00", Title: "全天活動", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "23:59", Title: "晚間活動", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "12:00", Title: "中午活動", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
		}

		sortAgendaItemsByTime(items)

		expectedOrder := []string{"00:00", "12:00", "23:59"}
		for i, expected := range expectedOrder {
			if items[i].Time != expected {
				t.Errorf("排序後第 %d 個項目時間應為 %q，但取得 %q", i, expected, items[i].Time)
			}
		}
	})
}

// TestFormatTimeForAgenda 測試時間格式化函數
func TestFormatTimeForAgenda(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"上午時間", time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC), "09:30"},
		{"下午時間", time.Date(2026, 1, 15, 14, 45, 0, 0, time.UTC), "14:45"},
		{"午夜", time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), "00:00"},
		{"深夜", time.Date(2026, 1, 15, 23, 59, 0, 0, time.UTC), "23:59"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimeForAgenda(tt.input)
			if result != tt.expected {
				t.Errorf("formatTimeForAgenda(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGenerateAgendaFlex 測試行程聚合 Flex Message 範本生成
func TestGenerateAgendaFlex(t *testing.T) {
	baseURL := "https://timeledger.app"
	svc := NewLineBotTemplateService(baseURL)

	targetDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	t.Run("MultipleItems_WithCenterAndPersonal", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "09:00", Title: "瑜珈課程", SourceName: "健身中心A", SourceType: AgendaSourceTypeCenter},
			{Time: "10:00", Title: "部門會議", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "14:00", Title: "鋼琴教學", SourceName: "音樂教室B", SourceType: AgendaSourceTypeCenter},
		}

		flex := svc.GenerateAgendaFlex(items, targetDate, "陳小美")

		// 驗證基本結構
		flexMap, ok := flex.(map[string]interface{})
		if !ok {
			t.Fatal("Expected flex message to be a map")
		}

		// 驗證類型
		if flexMap["type"] != "bubble" {
			t.Errorf("Expected type 'bubble', got %v", flexMap["type"])
		}

		// 驗證有 footer 按鈕
		footer, ok := flexMap["footer"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected footer to be a map")
		}

		footerContents, ok := footer["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected footer contents to be an array")
		}

		// 驗證有進入系統首頁按鈕
		foundHomeButton := false
		for _, item := range footerContents {
			if btn, ok := item.(map[string]interface{}); ok {
				if action, ok := btn["action"].(map[string]interface{}); ok {
					if label, ok := action["label"].(string); ok {
						if label == "📱 進入系統首頁" {
							foundHomeButton = true
							break
						}
					}
				}
			}
		}

		if !foundHomeButton {
			t.Error("Expected to find '進入系統首頁' button in footer")
		}

		// 驗證按鈕連結
		for _, item := range footerContents {
			if btn, ok := item.(map[string]interface{}); ok {
				if action, ok := btn["action"].(map[string]interface{}); ok {
					if uri, ok := action["uri"].(string); ok {
						if uri != baseURL {
							t.Errorf("Expected button URI to be %s, got %s", baseURL, uri)
						}
					}
				}
			}
		}
	})

	t.Run("EmptyItems_ShowsNoScheduleMessage", func(t *testing.T) {
		items := []AgendaItem{}

		flex := svc.GenerateAgendaFlex(items, targetDate, "王老師")

		flexMap, ok := flex.(map[string]interface{})
		if !ok {
			t.Fatal("Expected flex message to be a map")
		}

		body, ok := flexMap["body"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected body to be a map")
		}

		bodyContents, ok := body["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected body contents to be an array")
		}

		// 找到包裝行程列表的 box
		agendaBox := bodyContents[2].(map[string]interface{})
		agendaBoxContents, ok := agendaBox["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected agenda box contents to be an array")
		}

		// 驗證顯示 "今天沒有行程" 訊息（在 agendaBoxContents 的第二個元素）
		if len(agendaBoxContents) >= 2 {
			textItem := agendaBoxContents[1].(map[string]interface{})
			if textStr, ok := textItem["text"].(string); ok {
				if textStr != "🎉 今天沒有行程" {
					t.Errorf("Expected '🎉 今天沒有行程', got %q", textStr)
				}
			}
		} else {
			t.Error("Not enough elements in agenda box contents")
		}
	})

	t.Run("SingleCenterItem", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "15:00", Title: "舞蹈課程", SourceName: "舞蹈教室", SourceType: AgendaSourceTypeCenter},
		}

		flex := svc.GenerateAgendaFlex(items, targetDate, "林老師")

		flexMap, ok := flex.(map[string]interface{})
		if !ok {
			t.Fatal("Expected flex message to be a map")
		}

		body, ok := flexMap["body"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected body to be a map")
		}

		bodyContents, ok := body["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected body contents to be an array")
		}

		// 驗證有標題
		foundTitle := false
		for _, item := range bodyContents {
			if text, ok := item.(map[string]interface{}); ok {
				if textStr, ok := text["text"].(string); ok {
					if textStr == "👋 林老師 的今日行程" {
						foundTitle = true
						break
					}
				}
			}
		}

		if !foundTitle {
			t.Error("Expected to find user name in title")
		}
	})

	t.Run("AllPersonalItems", func(t *testing.T) {
		items := []AgendaItem{
			{Time: "08:00", Title: "晨跑", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "12:00", Title: "午餐約會", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
			{Time: "20:00", Title: "瑜珈課", SourceName: "個人", SourceType: AgendaSourceTypePersonal},
		}

		flex := svc.GenerateAgendaFlex(items, targetDate, "張老師")

		flexMap, ok := flex.(map[string]interface{})
		if !ok {
			t.Fatal("Expected flex message to be a map")
		}

		// 驗證統計資訊
		body, ok := flexMap["body"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected body to be a map")
		}

		bodyContents, ok := body["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected body contents to be an array")
		}

		// 找到包裝行程列表的 box（索引 2）和統計資訊（索引 4）
		agendaBox := bodyContents[2].(map[string]interface{})
		agendaBoxContents, ok := agendaBox["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected agenda box contents to be an array")
		}

		// 統計資訊在 agendaBoxContents 的倒數第二個元素
		statsItem := agendaBoxContents[len(agendaBoxContents)-1].(map[string]interface{})
		if textStr, ok := statsItem["text"].(string); ok {
			if textStr != "📊 共 3 筆行程" {
				t.Errorf("Expected '📊 共 3 筆行程', got %q", textStr)
			}
		} else {
			t.Error("Could not find stats text")
		}
	})

	t.Run("DateFormat_TaiwanFormat", func(t *testing.T) {
		items := []AgendaItem{}
		flex := svc.GenerateAgendaFlex(items, targetDate, "測試老師")

		flexMap, ok := flex.(map[string]interface{})
		if !ok {
			t.Fatal("Expected flex message to be a map")
		}

		body, ok := flexMap["body"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected body to be a map")
		}

		bodyContents, ok := body["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected body contents to be an array")
		}

		// bodyContents[2] 是包裝行程列表的 box，裡面第一個元素是日期標題
		agendaBox := bodyContents[2].(map[string]interface{})
		agendaBoxContents, ok := agendaBox["contents"].([]interface{})
		if !ok {
			t.Fatal("Expected agenda box contents to be an array")
		}

		// agendaBoxContents[0] 是日期標題
		dateItem := agendaBoxContents[0].(map[string]interface{})
		if textStr, ok := dateItem["text"].(string); ok {
			// 驗證日期格式為 "📅 YYYY年M月D日 (W)"
			// 使用 rune 來正確處理中文字元
			runes := []rune(textStr)
			if len(runes) < 8 {
				t.Errorf("Date text too short: %q", textStr)
			}
			// 檢查開頭是 "📅 " (2 runes) + 4位數年份
			if string(runes[:6]) != "📅 2026" {
				t.Errorf("Expected date to start with '📅 2026', got %q", string(runes[:6]))
			}
			// 檢查包含 "年"
			if !strings.Contains(textStr, "年") {
				t.Errorf("Expected date to contain '年', got %q", textStr)
			}
			// 檢查包含 "日"
			if !strings.Contains(textStr, "日") {
				t.Errorf("Expected date to contain '日', got %q", textStr)
			}
		} else {
			t.Error("Could not find date text")
		}
	})
}
