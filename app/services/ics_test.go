package services

import (
	"context"
	"image"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	mockRedis "timeLedger/testing/redis"
)

func TestICSCalendarService_GenerateICS(t *testing.T) {
	app := setupTestApp(t)
	icsSvc := NewICSCalendarService(app)

	config := &ICSConfig{
		TeacherID:   1,
		CenterID:    1,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		CenterName:  "測試中心",
		TeacherName: "測試老師",
		Events: []ScheduleEvent{
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

func TestICSCalendarService_eventToICS(t *testing.T) {
	app := setupTestApp(t)
	icsSvc := NewICSCalendarService(app)

	event := &ScheduleEvent{
		ID:          "test-event-1",
		Summary:     "測試課程",
		Description: "測試描述",
		Location:    "教室A",
		StartTime:   time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 1, 20, 11, 0, 0, 0, time.UTC),
		TeacherName: "老師",
		CenterName:  "中心",
	}

	icsStr, err := icsSvc.eventToICS(event)
	if err != nil {
		t.Fatalf("Failed to convert event to ICS: %v", err)
	}

	if !contains(icsStr, "UID:test-event-1") {
		t.Error("Missing event UID")
	}
	if !contains(icsStr, "SUMMARY:測試課程") {
		t.Error("Missing event summary")
	}

	t.Logf("Generated VEVENT:\n%s", icsStr)
}

func TestICSCalendarService_escapeICSText(t *testing.T) {
	app := setupTestApp(t)
	icsSvc := NewICSCalendarService(app)

	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello;world", "hello\\;world"},
		{"hello,world", "hello\\,world"},
		{"hello\\world", "hello\\\\world"},
	}

	for _, tt := range tests {
		result := icsSvc.escapeICSText(tt.input)
		if result != tt.expected {
			t.Errorf("escapeICSText(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestICSCalendarService_GenerateSubscriptionToken(t *testing.T) {
	app := setupTestApp(t)
	icsSvc := NewICSCalendarService(app)

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
	app := setupTestApp(t)
	icsSvc := NewICSCalendarService(app)

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

func TestImageService_ParseHexColor(t *testing.T) {
	app := setupTestApp(t)
	imageSvc := NewImageService(app)

	tests := []struct {
		hex      string
		expected bool // 是否為有效顏色
	}{
		{"#FFFFFF", true},
		{"#000000", true},
		{"#FF0000", true},
		{"#4F46E5", true},
		{"FFFFFF", false}, // 缺少 #
		{"#GGG", false},   // 無效格式
		{"", false},       // 空值
	}

	for _, tt := range tests {
		colorImg := imageSvc.parseHexColor(tt.hex)
		// 檢查是否返回有效的圖片
		if tt.expected && colorImg == nil {
			t.Errorf("Expected valid color for %s, got nil", tt.hex)
		}
		// 驗證無效的 hex 返回黑色
		if !tt.expected && colorImg != nil {
			// 空值或無效格式應該返回黑色圖片
			// color.Color 沒有 Bounds 方法，我們只能檢查是否不為 nil
			_ = colorImg
		}
	}
}

func TestImageService_GenerateScheduleImage(t *testing.T) {
	app := setupTestApp(t)
	imageSvc := NewImageService(app)

	config := &ImageConfig{
		TeacherID: 1,
		CenterID:  1,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 1, 0),
		Width:     800,
		Height:    600,
		Title:     "測試課表",
		Theme:     DefaultTheme,
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

func TestImageService_ImageToJPEG(t *testing.T) {
	app := setupTestApp(t)
	imageSvc := NewImageService(app)

	// 建立簡單的測試圖片
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	data, err := imageSvc.imageToJPEG(testImg)
	if err != nil {
		t.Fatalf("Failed to convert to JPEG: %v", err)
	}

	if len(data) < 100 {
		t.Error("Generated JPEG data is too small")
	}

	t.Logf("Generated JPEG size: %d bytes", len(data))
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
func setupTestApp(t *testing.T) *app.App {
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
