package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	"github.com/gin-gonic/gin"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupNotificationTestApp() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("MySQL init error: " + err.Error())
	}

	// 初始化 Redis mock
	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic("Redis init error: " + err.Error())
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	// 初始化測試用的 Env 配置
	env := &configs.Env{
		JWTSecret:              "test-jwt-secret-key-for-testing-only",
		AppEnv:                 "test",
		AppDebug:               true,
		AppTimezone:            "Asia/Taipei",
		LineChannelSecret:      "test-secret",
		LineChannelAccessToken: "test-token",
		FrontendBaseURL:        "https://timeledger.example.com",
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

// TestNotificationQueueService_PushNotification 測試將通知加入 Redis 佇列
func TestNotificationQueueService_PushNotification(t *testing.T) {
	appInstance, _, cleanup := setupNotificationTestApp()
	defer cleanup()

	// 建立 notification queue service
	queueService := services.NewNotificationQueueService(appInstance)

	// 建立測試用佇列項目
	queueItem := &models.NotificationQueue{
		Type:          models.NotificationTypeExceptionSubmit,
		RecipientID:   1,
		RecipientType: "ADMIN",
		Payload:       `{"type":"text","text":"測試訊息"}`,
		Status:        models.NotificationStatusPending,
		ScheduledAt:   time.Now(),
	}

	err := queueService.PushNotification(context.Background(), queueItem)
	if err != nil {
		t.Errorf("Failed to push notification to queue: %v", err)
	}
}

// TestNotificationQueueService_ProcessQueue 測試處理通知佇列
// 注意：此測試會無限阻塞，因為 BRPop 不支援 context 取消
// ProcessQueue 設計為背景 worker 持續執行，不適合單元測試
func TestNotificationQueueService_ProcessQueue(t *testing.T) {
	t.Skip("ProcessQueue 使用 BRPop 會無限阻塞，設計為背景 worker 執行，不適合單元測試")
}

// TestNotificationQueueService_NotifyExceptionSubmitted 測試例外申請通知
func TestNotificationQueueService_NotifyExceptionSubmitted(t *testing.T) {
	appInstance, _, cleanup := setupNotificationTestApp()
	defer cleanup()

	// 建立 notification queue service
	queueService := services.NewNotificationQueueService(appInstance)

	// 測試用例外資料
	exception := &models.ScheduleException{
		ID:            0, // 會被資料庫產生
		CenterID:      1,
		RuleID:        1,
		OriginalDate:  time.Now(),
		ExceptionType: "LEAVE",
		Status:        "PENDING",
		Reason:        "測試請假",
	}

	// 測試通知函數（不會真的發送 LINE 訊息，因爲佇列是空的）
	err := queueService.NotifyExceptionSubmitted(context.Background(), exception, "陳小美", "Yoga Space")
	if err != nil {
		t.Errorf("NotifyExceptionSubmitted should not error: %v", err)
	}
}

// TestNotificationQueueService_NotifyWelcomeAdmin 測試管理員歡迎訊息
func TestNotificationQueueService_NotifyWelcomeAdmin(t *testing.T) {
	appInstance, _, cleanup := setupNotificationTestApp()
	defer cleanup()

	// 建立 notification queue service
	queueService := services.NewNotificationQueueService(appInstance)

	// 測試用管理員（未綁定 LINE）
	admin := &models.AdminUser{
		ID:         1,
		LineUserID: "", // 未綁定
		Name:       "測試管理員",
	}

	// 應該不會報錯，但因爲未綁定所以不會實際發送
	err := queueService.NotifyWelcomeAdmin(context.Background(), admin, "測試中心")
	if err != nil {
		t.Errorf("NotifyWelcomeAdmin should not error for unbound admin: %v", err)
	}
}

// TestLineWebhookHandler 測試 LINE Webhook 處理
func TestLineWebhookHandler(t *testing.T) {
	// 建立 mock 伺服器來接收 LINE Webhook
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// 驗證 Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected application/json, got %s", contentType)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	// 驗證簽名函數
	env := &configs.Env{
		LineChannelSecret: "test-secret",
	}
	testApp := &app.App{
		Env: env,
	}

	lineBotService := services.NewLineBotService(testApp)

	// 測試無效的簽名
	body := []byte(`{"events":[{"type":"message"}]}`)
	if lineBotService.VerifySignature(body, "invalid-signature") {
		t.Error("Expected invalid signature to fail verification")
	}
}

// TestFlexMessageJSON 測試 Flex Message JSON 格式正確性
func TestFlexMessageJSON(t *testing.T) {
	templateService := services.NewLineBotTemplateService("https://timeledger.example.com")

	teacher := &models.Teacher{
		ID:   1,
		Name: "陳小美",
	}

	template := templateService.GetWelcomeTeacherTemplate(teacher, "Yoga Space")

	// 序列化爲 JSON
	jsonBytes, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal template to JSON: %v", err)
	}

	// 反序列化回來
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// 驗證必要欄位
	if parsed["type"] != "bubble" {
		t.Errorf("Expected type to be bubble, got: %v", parsed["type"])
	}

	// 驗證有 body 區塊
	body, ok := parsed["body"].(map[string]interface{})
	if !ok {
		t.Error("Expected body to be a map")
	} else {
		// 驗證 body 有 contents
		if body["contents"] == nil {
			t.Error("Expected body to have contents")
		}
	}

	// 驗證有 hero 或 body 區塊
	if parsed["hero"] == nil && parsed["body"] == nil {
		t.Error("Expected template to have hero or body section")
	}
}
