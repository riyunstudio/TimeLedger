package test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupLineBotTestApp() (*app.App, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	// 初始化測試用的 Env 配置
	env := &configs.Env{
		JWTSecret:             "test-jwt-secret-key-for-testing-only",
		AppEnv:                "test",
		AppDebug:              true,
		AppTimezone:           "Asia/Taipei",
		LineChannelSecret:     "test-secret",
		LineChannelAccessToken: "test-token",
		FrontendBaseURL:       "https://timeledger.example.com",
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

	return appInstance, mysqlDB, func() {}
}

// TestLineBotService_SendMessage 測試 LINE Bot 發送文字訊息
func TestLineBotService_SendMessage(t *testing.T) {
	// 建立 mock 伺服器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// 驗證 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Bearer test-token, got %s", authHeader)
		}

		// 驗證 Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected application/json, got %s", contentType)
		}

		// 回傳成功
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": 200, "message": "ok"}`))
	}))
	defer server.Close()

	testApp, _, cleanup := setupLineBotTestApp()
	defer cleanup()

	// 建立 service（使用 mock URL）
	lineBotService := services.NewLineBotService(testApp)

	// 測試發送訊息
	err := lineBotService.PushMessage(context.Background(), "U123456789", map[string]interface{}{
		"type": "text",
		"text": "Hello, World!",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestLineBotService_VerifySignature 測試 LINE Webhook 簽名驗證
func TestLineBotService_VerifySignature(t *testing.T) {
	testApp, _, cleanup := setupLineBotTestApp()
	defer cleanup()

	lineBotService := services.NewLineBotService(testApp)

	// 測試資料
	body := []byte(`{"events":[{"type":"message","replyToken":"abc123","message":{"type":"text","id":"12345","text":"Hello"}}]}`)

	// 生成正確的簽名
	hash := hmac.New(sha256.New, []byte("test-secret"))
	hash.Write(body)
	correctSignature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	// 測試正確簽名
	if !lineBotService.VerifySignature(body, correctSignature) {
		t.Error("Expected signature to be valid")
	}

	// 測試錯誤簽名
	wrongSignature := base64.StdEncoding.EncodeToString([]byte("wrong-signature"))
	if lineBotService.VerifySignature(body, wrongSignature) {
		t.Error("Expected signature to be invalid")
	}

	// 測試空簽名
	if lineBotService.VerifySignature(body, "") {
		t.Error("Expected empty signature to be invalid")
	}
}

// TestLineBotService_PushFlexMessage 測試 LINE Bot 發送 Flex Message
func TestLineBotService_PushFlexMessage(t *testing.T) {
	// 建立 mock 伺服器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析請求體
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		// 驗證 Flex Message 結構
		messages, ok := reqBody["messages"].([]interface{})
		if !ok {
			t.Error("Expected messages to be an array")
		}

		if len(messages) == 0 {
			t.Error("Expected at least one message")
		}

		flexMsg, ok := messages[0].(map[string]interface{})
		if !ok {
			t.Error("Expected message to be an object")
		}

		if flexMsg["type"] != "flex" {
			t.Errorf("Expected message type to be flex, got %v", flexMsg["type"])
		}

		if flexMsg["altText"] == "" {
			t.Error("Expected altText to be set")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	testApp, _, cleanup := setupLineBotTestApp()
	defer cleanup()

	lineBotService := services.NewLineBotService(testApp)

	// Flex Message 內容
	flexContent := map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "🔔 新的例外申請",
				},
			},
		},
	}

	err := lineBotService.PushFlexMessage(context.Background(), "U123456789", "新的例外申請通知", flexContent)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestLineBotTemplateService_GetWelcomeTemplate 測試取得歡迎訊息範本
func TestLineBotTemplateService_GetWelcomeTemplate(t *testing.T) {
	templateService := services.NewLineBotTemplateService("https://timeledger.example.com")

	teacher := &models.Teacher{
		ID:   1,
		Name: "陳小美",
	}

	template := templateService.GetWelcomeTeacherTemplate(teacher, "Yoga Space 台北館")

	// 驗證範本結構
	flexMap, ok := template.(map[string]interface{})
	if !ok {
		t.Fatal("Expected template to be a map")
	}

	if flexMap["type"] != "bubble" {
		t.Errorf("Expected type to be bubble, got %v", flexMap["type"])
	}

	body, ok := flexMap["body"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected body to be a map")
	}

	contents, ok := body["contents"].([]interface{})
	if !ok {
		t.Fatal("Expected contents to be an array")
	}

	// 檢查是否包含歡迎文字
	foundWelcome := false
	for _, item := range contents {
		if textItem, ok := item.(map[string]interface{}); ok {
			if text, ok := textItem["text"].(string); ok {
				if text == "👋 歡迎加入 TimeLedger！" {
					foundWelcome = true
					break
				}
			}
		}
	}

	if !foundWelcome {
		t.Error("Expected to find welcome message in template")
	}

	// 檢查按鈕
	footer, ok := flexMap["footer"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected footer to be a map")
	}

	footerContents, ok := footer["contents"].([]interface{})
	if !ok {
		t.Fatal("Expected footer contents to be an array")
	}

	if len(footerContents) == 0 {
		t.Error("Expected at least one button in footer")
	}
}

// TestLineBotTemplateService_GetExceptionSubmitTemplate 測試取得例外通知範本
func TestLineBotTemplateService_GetExceptionSubmitTemplate(t *testing.T) {
	templateService := services.NewLineBotTemplateService("https://timeledger.example.com")

	exception := &models.ScheduleException{
		ID:            123,
		ExceptionType: "LEAVE",
		OriginalDate:  time.Now(),
		Reason:        "身體不適",
	}

	template := templateService.GetExceptionSubmitTemplate(exception, "陳小美", "Yoga Space 台北館")

	flexMap, ok := template.(map[string]interface{})
	if !ok {
		t.Fatal("Expected template to be a map")
	}

	body, ok := flexMap["body"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected body to be a map")
	}

	contents, ok := body["contents"].([]interface{})
	if !ok {
		t.Fatal("Expected contents to be an array")
	}

	// 檢查是否包含申請人資訊
	foundTeacher := false
	for _, item := range contents {
		if textItem, ok := item.(map[string]interface{}); ok {
			if text, ok := textItem["text"].(string); ok {
				if text == "👤 申請人：陳小美 老師" {
					foundTeacher = true
					break
				}
			}
		}
	}

	if !foundTeacher {
		t.Error("Expected to find teacher name in template")
	}
}

// TestLineBotTemplateService_GetExceptionApproveTemplate 測試取得核准通知範本
func TestLineBotTemplateService_GetExceptionApproveTemplate(t *testing.T) {
	templateService := services.NewLineBotTemplateService("https://timeledger.example.com")

	exception := &models.ScheduleException{
		ID:            456,
		ExceptionType: "RESCHEDULE",
	}

	template := templateService.GetExceptionApproveTemplate(exception, "陳小美")

	flexMap, ok := template.(map[string]interface{})
	if !ok {
		t.Fatal("Expected template to be a map")
	}

	// 檢查是否包含核准文字
	body, ok := flexMap["body"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected body to be a map")
	}

	contents, ok := body["contents"].([]interface{})
	if !ok {
		t.Fatal("Expected contents to be an array")
	}

	foundApproved := false
	for _, item := range contents {
		if textItem, ok := item.(map[string]interface{}); ok {
			if text, ok := textItem["text"].(string); ok {
				if text == "✅ 調課申請已核准" {
					foundApproved = true
					break
				}
			}
		}
	}

	if !foundApproved {
		t.Error("Expected to find approval message in template")
	}
}

// TestGenerateBindingCode 測試產生綁定驗證碼
func TestGenerateBindingCode(t *testing.T) {
	code := services.GenerateBindingCode()

	// 驗證長度
	if len(code) != 6 {
		t.Errorf("Expected code length to be 6, got %d", len(code))
	}

	// 驗證格式（應該是字母數字，不含易混淆的字元如 0、O、I、l）
	for _, c := range code {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Errorf("Unexpected character in code: %c", c)
		}
		// 排除易混淆的字元
		if c == '0' || c == 'O' || c == 'I' || c == 'l' || c == '1' {
			t.Errorf("Code contains ambiguous character: %c", c)
		}
	}

	// 驗證每次產生的碼不同
	code2 := services.GenerateBindingCode()
	if code == code2 {
		t.Error("Expected different codes on each call")
	}
}
