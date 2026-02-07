package test

import (
	"context"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	mockRedis "timeLedger/testing/redis"
)

func TestICSCalendarService_GenerateICS(t *testing.T) {
	app := setupICSTestApp(t)
	icsSvc := services.NewICSCalendarService(app)

	config := &services.ICSConfig{
		TeacherID:   1,
		CenterID:    1,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		CenterName:  "測試中心",
		TeacherName: "測試老師",
		Events: []services.ScheduleEvent{
			{
				ID:          "event-1",
				Summary:     "瑜珈課程",
				Description: "基礎瑜珈課程",
				Location:    "教室A",
				StartTime:   time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2026, 1, 20, 11, 0, 0, 0, time.UTC),
				TeacherName: "測試老師",
				CenterName:  "測試中心",
			},
			{
				ID:          "event-2",
				Summary:     "舞蹈課程",
				Description: "進階舞蹈課程",
				Location:    "教室B",
				StartTime:   time.Date(2026, 1, 21, 14, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2026, 1, 21, 15, 30, 0, 0, time.UTC),
				TeacherName: "測試老師",
				CenterName:  "測試中心",
			},
		},
	}

	data, err := icsSvc.GenerateICS(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to generate ICS: %v", err)
	}

	// 驗證產出的 ICS 格式
	icsStr := string(data)
	if !contains(icsStr, "BEGIN:VCALENDAR") {
		t.Error("Missing VCALENDAR header")
	}
	if !contains(icsStr, "END:VCALENDAR") {
		t.Error("Missing VCALENDAR footer")
	}
	if !contains(icsStr, "BEGIN:VEVENT") {
		t.Error("Missing VEVENT")
	}
	if !contains(icsStr, "END:VEVENT") {
		t.Error("Missing VEVENT footer")
	}
	if !contains(icsStr, "瑜珈課程") {
		t.Error("Missing event summary")
	}

	t.Logf("Generated ICS:\n%s", icsStr)
}

func TestICSCalendarService_GenerateSubscriptionToken(t *testing.T) {
	app := setupICSTestApp(t)
	icsSvc := services.NewICSCalendarService(app)

	token, err := icsSvc.GenerateSubscriptionToken(1, 1)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
	if len(token) < 10 {
		t.Error("Token should be at least 10 characters")
	}

	t.Logf("Generated token: %s", token)
}

func TestICSCalendarService_GenerateSubscriptionURL(t *testing.T) {
	app := setupICSTestApp(t)
	icsSvc := services.NewICSCalendarService(app)

	// 測試預設 baseURL
	url := icsSvc.GenerateSubscriptionURL("test-token-123")
	expectedBase := "https://timeledger.app"
	if !contains(url, expectedBase) {
		t.Errorf("Expected URL to contain %s, got %s", expectedBase, url)
	}
	if !contains(url, "test-token-123.ics") {
		t.Errorf("Expected URL to contain token, got %s", url)
	}

	t.Logf("Generated subscription URL: %s", url)
}

func TestImageService_GenerateScheduleImage(t *testing.T) {
	app := setupICSTestApp(t)
	imageSvc := services.NewImageService(app)

	config := &services.ImageConfig{
		TeacherID: 1,
		CenterID:  1,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 1, 0),
		Width:     800,
		Height:    600,
		Title:     "測試課表",
		Theme:     services.DefaultTheme,
	}

	data, err := imageSvc.GenerateScheduleImage(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	// 驗證產出的 JPEG 格式
	if len(data) < 100 {
		t.Error("Generated image data is too small")
	}

	// JPEG 檔案通常以 FFD8 開頭
	if data[0] != 0xFF || data[1] != 0xD8 {
		t.Error("Generated image does not appear to be JPEG format")
	}

	t.Logf("Generated image size: %d bytes", len(data))
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// setupTestApp 建立測試用的 App 實例
func setupICSTestApp(t *testing.T) *app.App {
	t.Helper()

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("MySQL not available, skipping test: %v", err)
		return nil
	}

	rdb, _, err := mockRedis.Initialize()
	if err != nil {
		t.Skipf("Redis not available, skipping test: %v", err)
		return nil
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

	return appInstance
}
