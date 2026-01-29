package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// LineBotService LINE Messaging API 服務
type LineBotService interface {
	// 基礎訊息發送
	PushMessage(ctx context.Context, userID string, message interface{}) error
	ReplyMessage(ctx context.Context, replyToken string, message interface{}) error
	Multicast(ctx context.Context, userIDs []string, message interface{}) error

	// Flex Message 發送（便捷方法）
	PushFlexMessage(ctx context.Context, userID string, altText string, flexContent interface{}) error
	ReplyFlexMessage(ctx context.Context, replyToken string, altText string, flexContent interface{}) error

	// 驗證
	VerifySignature(body []byte, signature string) bool

	// 範本發送
	SendWelcomeTeacher(ctx context.Context, teacher *models.Teacher, centerName string) error
	SendWelcomeAdmin(ctx context.Context, admin *models.AdminUser, centerName string) error
	SendExceptionNotification(ctx context.Context, admin *models.AdminUser, exception *models.ScheduleException, teacherName string) error
	SendInvitationAcceptedNotification(ctx context.Context, admins []*models.AdminUser, teacher *models.Teacher, centerName string, role string) error
}

// LineBotServiceImpl LINE Messaging API 服務實現
type LineBotServiceImpl struct {
	app             *app.App
	channelSecret   string
	channelToken    string
	apiURL          string
	profileURL      string
	replyURL        string
	multicastURL    string
	client          *http.Client
	templateService LineBotTemplateService
}

// NewLineBotService 建立 LINE Bot Service
func NewLineBotService(app *app.App) LineBotService {
	return &LineBotServiceImpl{
		app:             app,
		channelSecret:   app.Env.LineChannelSecret,
		channelToken:    app.Env.LineChannelAccessToken,
		apiURL:          "https://api.line.me/v2/bot/message/push",
		profileURL:      "https://api.line.me/v2/bot/profile",
		replyURL:        "https://api.line.me/v2/bot/message/reply",
		multicastURL:    "https://api.line.me/v2/bot/message/multicast",
		templateService: NewLineBotTemplateService(app.Env.FrontendBaseURL),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LineBotRequest LINE Bot API 請求結構
type LineBotRequest struct {
	Messages []interface{} `json:"messages"`
}

// LineBotReplyRequest LINE Bot Reply API 請求結構
type LineBotReplyRequest struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

// LineBotMulticastRequest LINE Bot Multicast API 請求結構
type LineBotMulticastRequest struct {
	To       []string      `json:"to"`
	Messages []interface{} `json:"messages"`
}

// LineBotErrorResponse LINE Bot API 錯誤回應
type LineBotErrorResponse struct {
	Message string `json:"message"`
	Details []struct {
		Property string `json:"property"`
		Message  string `json:"message"`
	} `json:"details"`
}

// PushMessage 發送推播訊息給單一用戶
func (s *LineBotServiceImpl) PushMessage(ctx context.Context, userID string, message interface{}) error {
	messages := []interface{}{message}
	return s.pushToEndpoint(ctx, s.apiURL, map[string]interface{}{
		"to":       userID,
		"messages": messages,
	})
}

// ReplyMessage 回復訊息（用於 Webhook 回覆）
func (s *LineBotServiceImpl) ReplyMessage(ctx context.Context, replyToken string, message interface{}) error {
	messages := []interface{}{message}
	reqBody := LineBotReplyRequest{
		ReplyToken: replyToken,
		Messages:   messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	return s.sendRequest(ctx, s.replyURL, jsonBody)
}

// Multicast 群發訊息給多個用戶
func (s *LineBotServiceImpl) Multicast(ctx context.Context, userIDs []string, message interface{}) error {
	if len(userIDs) == 0 {
		return nil
	}

	messages := []interface{}{message}
	reqBody := LineBotMulticastRequest{
		To:       userIDs,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	return s.sendRequest(ctx, s.multicastURL, jsonBody)
}

// PushFlexMessage 發送 Flex Message 給單一用戶
func (s *LineBotServiceImpl) PushFlexMessage(ctx context.Context, userID string, altText string, flexContent interface{}) error {
	flexMessage := s.createFlexMessage(altText, flexContent)
	return s.PushMessage(ctx, userID, flexMessage)
}

// ReplyFlexMessage 回復 Flex Message
func (s *LineBotServiceImpl) ReplyFlexMessage(ctx context.Context, replyToken string, altText string, flexContent interface{}) error {
	flexMessage := s.createFlexMessage(altText, flexContent)
	return s.ReplyMessage(ctx, replyToken, flexMessage)
}

// createFlexMessage 建立 Flex Message 結構
func (s *LineBotServiceImpl) createFlexMessage(altText string, contents interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "flex",
		"altText":  altText,
		"contents": contents,
	}
}

// pushToEndpoint 發送到指定端點
func (s *LineBotServiceImpl) pushToEndpoint(ctx context.Context, endpoint string, body map[string]interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	return s.sendRequest(ctx, endpoint, jsonBody)
}

// sendRequest 發送 HTTP 請求到 LINE API
func (s *LineBotServiceImpl) sendRequest(ctx context.Context, endpoint string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.channelToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		var errResp LineBotErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return fmt.Errorf("LINE API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("LINE API error: %s", errResp.Message)
	}

	return nil
}

// VerifySignature 驗證 LINE Webhook 簽名
func (s *LineBotServiceImpl) VerifySignature(body []byte, signature string) bool {
	expectedSignature := s.generateSignature(body)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// generateSignature 生成簽名
func (s *LineBotServiceImpl) generateSignature(body []byte) string {
	hash := hmac.New(sha256.New, []byte(s.channelSecret))
	hash.Write(body)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// SendWelcomeTeacher 發送老師歡迎訊息
func (s *LineBotServiceImpl) SendWelcomeTeacher(ctx context.Context, teacher *models.Teacher, centerName string) error {
	if teacher.LineUserID == "" {
		return nil // 沒有 LINE ID，無法發送
	}

	flexMessage := s.templateService.GetWelcomeTeacherTemplate(teacher, centerName)
	return s.PushFlexMessage(ctx, teacher.LineUserID, "歡迎加入 TimeLedger！", flexMessage)
}

// SendWelcomeAdmin 發送管理員歡迎訊息
func (s *LineBotServiceImpl) SendWelcomeAdmin(ctx context.Context, admin *models.AdminUser, centerName string) error {
	if admin.LineUserID == "" {
		return nil // 沒有 LINE ID，無法發送
	}

	flexMessage := s.templateService.GetWelcomeAdminTemplate(admin, centerName)
	return s.PushFlexMessage(ctx, admin.LineUserID, "歡迎使用 TimeLedger！", flexMessage)
}

// SendExceptionNotification 發送例外申請通知給管理員
func (s *LineBotServiceImpl) SendExceptionNotification(ctx context.Context, admin *models.AdminUser, exception *models.ScheduleException, teacherName string) error {
	if admin.LineUserID == "" || !admin.LineNotifyEnabled {
		return nil // 沒有綁定或通知已關閉
	}

	flexMessage := s.templateService.GetExceptionSubmitTemplate(exception, teacherName, "")
	return s.PushFlexMessage(ctx, admin.LineUserID, "新的例外申請通知", flexMessage)
}

// SendInvitationAcceptedNotification 發送邀請接受通知給管理員
func (s *LineBotServiceImpl) SendInvitationAcceptedNotification(ctx context.Context, admins []*models.AdminUser, teacher *models.Teacher, centerName string, role string) error {
	if len(admins) == 0 {
		return nil // 沒有管理員可通知
	}

	flexMessage := s.templateService.GetInvitationAcceptedTemplate(teacher, centerName, role)

	// 收集所有已綁定 LINE 的管理員
	var lineUserIDs []string
	for _, admin := range admins {
		if admin.LineUserID != "" && admin.LineNotifyEnabled {
			lineUserIDs = append(lineUserIDs, admin.LineUserID)
		}
	}

	if len(lineUserIDs) == 0 {
		return nil // 沒有已綁定的管理員
	}

	// 群發通知給所有已綁定的管理員
	return s.Multicast(ctx, lineUserIDs, flexMessage)
}
