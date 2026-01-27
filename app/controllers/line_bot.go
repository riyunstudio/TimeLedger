package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

// LineBotController LINE Bot Webhook Controller
type LineBotController struct {
	app            *app.App
	lineBotService services.LineBotService
	qrCodeService  *services.QRCodeService
	adminService   *services.AdminUserService
}

// NewLineBotController 建立 LINE Bot Controller
func NewLineBotController(app *app.App) *LineBotController {
	return &LineBotController{
		app:           app,
		lineBotService: services.NewLineBotService(app),
		qrCodeService: services.NewQRCodeService(),
		adminService:  services.NewAdminUserService(app),
	}
}

// LINEWebhookRequest LINE Webhook 請求結構
type LINEWebhookRequest struct {
	Destination string             `json:"destination"`
	Events      []LINEWebhookEvent `json:"events"`
}

// LINEWebhookEvent LINE Webhook 事件
type LINEWebhookEvent struct {
	Type       string          `json:"type"`
	Mode       string          `json:"mode"`
	Timestamp  int64           `json:"timestamp"`
	Source     LINEEventSource `json:"source"`
	ReplyToken string          `json:"replyToken,omitempty"`
	Message    LINEEventMessage `json:"message,omitempty"`
}

// LINEEventSource 事件來源
type LINEEventSource struct {
	Type   string `json:"type"`
	UserID string `json:"userId,omitempty"`
}

// LINEEventMessage 事件訊息
type LINEEventMessage struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Text        string `json:"text,omitempty"`
	QuoteToken  string `json:"quoteToken,omitempty"`
}

// HandleWebhook 處理 LINE Webhook
func (c *LineBotController) HandleWebhook(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read webhook body: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// 驗證簽名
	signature := ctx.GetHeader("X-Line-Signature")
	if !c.lineBotService.VerifySignature(body, signature) {
		fmt.Printf("[WARN] Invalid LINE signature\n")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// 解析請求
	var webhookReq LINEWebhookRequest
	if err := json.Unmarshal(body, &webhookReq); err != nil {
		fmt.Printf("[ERROR] Failed to parse webhook request: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
		return
	}

	// 處理每個事件
	for _, event := range webhookReq.Events {
		go c.handleEvent(ctx, &event)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleEvent 處理單個事件
func (c *LineBotController) handleEvent(gctx *gin.Context, event *LINEWebhookEvent) {
	switch event.Type {
	case "message":
		c.handleMessageEvent(gctx, event)
	case "follow":
		c.handleFollowEvent(gctx, event)
	case "unfollow":
		c.handleUnfollowEvent(gctx, event)
	default:
		fmt.Printf("[DEBUG] Unhandled event type: %s\n", event.Type)
	}
}

// handleMessageEvent 處理訊息事件
func (c *LineBotController) handleMessageEvent(gctx *gin.Context, event *LINEWebhookEvent) {
	if event.Message.Type != "text" {
		return
	}

	text := event.Message.Text
	userID := event.Source.UserID

	// 處理驗證碼（6位數大寫字母數字）
	if len(text) == 6 && isValidBindingCode(text) {
		c.processBindingCode(gctx, text, userID, event.ReplyToken)
		return
	}

	// 處理關鍵字
	switch text {
	case "綁定", "bind", "Bind":
		c.sendBindingInstructions(gctx, event.ReplyToken)
	case "幫助", "幫我", "help", "Help":
		c.sendHelpMessage(gctx, event.ReplyToken)
	case "狀態", "status", "Status":
		c.sendStatusMessage(gctx, event.ReplyToken, userID)
	case "解除綁定", "unbind", "Unbind":
		c.sendUnbindInstructions(gctx, event.ReplyToken)
	case "了解更多", "更多", "more", "More":
		c.sendMoreInfoMessage(gctx, event.ReplyToken)
	case "稍後綁定", "稍後再說":
		c.sendAckMessage(gctx, event.ReplyToken)
	default:
		c.sendDefaultResponse(gctx, event.ReplyToken)
	}
}

// handleFollowEvent 處理加入好友事件
func (c *LineBotController) handleFollowEvent(gctx *gin.Context, event *LINEWebhookEvent) {
	userID := event.Source.UserID
	fmt.Printf("[INFO] User followed: %s\n", userID)

	welcomeMessage := map[string]interface{}{
		"type": "text",
		"text": "👋 您好！歡迎加入 TimeLedger！\n\n" +
			"如果您是管理員，請登入後台進行 LINE 綁定，即可收到即時例外通知。\n\n" +
			"輸入「綁定」開始綁定流程。",
	}

	c.lineBotService.ReplyMessage(gctx.Request.Context(), event.ReplyToken, welcomeMessage)
}

// handleUnfollowEvent 處理封鎖/取消好友事件
func (c *LineBotController) handleUnfollowEvent(gctx *gin.Context, event *LINEWebhookEvent) {
	userID := event.Source.UserID
	fmt.Printf("[INFO] User unfollowed: %s\n", userID)
}

// processBindingCode 處理綁定驗證碼
func (c *LineBotController) processBindingCode(gctx *gin.Context, code string, userID string, replyToken string) {
	_, eInfo, err := c.adminService.VerifyLINEBinding(gctx.Request.Context(), code, userID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to verify binding code: %v\n", err)
		errorMsg := "❌ 綁定失敗，驗證碼錯誤或已過期。"
		if eInfo != nil {
			if eInfo.Code == 90004 {
				errorMsg = "❌ 驗證碼已過期，請至後台重新產生。"
			}
		}
		c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, map[string]interface{}{
			"type": "text",
			"text": errorMsg,
		})
		return
	}

	// 綁定成功
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, map[string]interface{}{
		"type": "text",
		"text": "✅ 綁定成功！\n\n" +
			"您將會收到：\n" +
			"🔔 老師提交例外申請的通知\n" +
			"🔔 審核結果通知\n\n" +
			"如需調整通知設定，請至後台「設定」→「通知設定」。",
	})
}

// sendBindingInstructions 發送綁定說明
func (c *LineBotController) sendBindingInstructions(gctx *gin.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "🔗 綁定步驟：\n\n" +
			"1. 登入 TimeLedger 管理後台\n" +
			"2. 點擊右上角「設定」\n" +
			"3. 點擊「LINE 通知」\n" +
			"4. 點擊「開始綁定」\n" +
			"5. 掃描 QR Code 或輸入顯示的驗證碼\n\n" +
			"如有問題，請聯繫系統管理員。",
	}
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, message)
}

// sendHelpMessage 發送幫助訊息
func (c *LineBotController) sendHelpMessage(gctx *gin.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "❓ TimeLedger 指令說明：\n\n" +
			"📌 綁定相關：\n" +
			"• 「綁定」- 開始 LINE 綁定流程\n" +
			"• 「解除綁定」- 解除 LINE 綁定\n\n" +
			"📌 查詢相關：\n" +
			"• 「狀態」- 查看綁定狀態\n" +
			"• 「幫助」- 顯示此說明訊息\n\n" +
			"如有問題，請聯繫系統管理員。",
	}
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, message)
}

// sendStatusMessage 發送狀態訊息
func (c *LineBotController) sendStatusMessage(gctx *gin.Context, replyToken string, userID string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "📊 狀態查詢：\n\n" +
			"您的 LINE 帳號已與 TimeLedger 綁定。\n\n" +
			"如需調整設定，請至管理後台。",
	}
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, message)
}

// sendUnbindInstructions 發送解除綁定說明
func (c *LineBotController) sendUnbindInstructions(gctx *gin.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "🔓 解除綁定：\n\n" +
			"請至 TimeLedger 管理後台：\n" +
			"1. 點擊右上角「設定」\n" +
			"2. 點擊「LINE 通知」\n" +
			"3. 點擊「解除綁定」\n" +
			"4. 確認解除綁定\n\n" +
			"⚠️ 解除綁定後將無法收到即時通知。",
	}
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, message)
}

// sendMoreInfoMessage 發送更多資訊
func (c *LineBotController) sendMoreInfoMessage(gctx *gin.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "ℹ️ TimeLedger 介紹：\n\n" +
			"TimeLedger 是教師中心化多據點排課平台，\n" +
			"讓您可以：\n\n" +
			"• 接收老師的例外申請通知\n" +
			"• 即時處理請假、調課等申請\n" +
			"• 透過手機 LINE 隨時掌握動態\n\n" +
			"如有問題，請聯繫系統管理員。",
	}
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, message)
}

// sendAckMessage 發送確認訊息
func (c *LineBotController) sendAckMessage(gctx *gin.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "ℹ️ 了解！\n\n" +
			"您可以稍後再進行綁定。\n" +
			"當您準備好時，輸入「綁定」即可開始流程。",
	}
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, message)
}

// sendDefaultResponse 發送預設回應
func (c *LineBotController) sendDefaultResponse(gctx *gin.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "🤔 我不太理解您的意思。\n\n" +
			"輸入「幫助」查看可用指令。",
	}
	c.lineBotService.ReplyMessage(gctx.Request.Context(), replyToken, message)
}

// isValidBindingCode 檢查是否為有效的綁定驗證碼格式
func isValidBindingCode(code string) bool {
	if len(code) != 6 {
		return false
	}
	for _, c := range code {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

// HealthCheck 健康檢查
func (c *LineBotController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "line-bot",
		"timestamp": time.Now().Unix(),
	})
}

// GenerateLINEBindingQR 產生 LINE 官方帳號加好友 QR Code
// @Summary 產生 LINE 官方帳號 QR Code
// @Description 產生可掃描開啟 LINE 官方帳號聊天的 QR Code
// @Tags LINE
// @Produce octet-stream
// @Security BearerAuth
// @Success 200 {file} binary "PNG image"
// @Failure 401 {object} global.ApiResponse
// @Failure 500 {object} global.ApiResponse
// @Router /api/v1/admin/me/line/qrcode [GET]
func (c *LineBotController) GenerateLINEBindingQR(ctx *gin.Context) {
	// 驗證管理員身份（從 JWT middleware 取得）
	_, exists := ctx.Get(string(global.UserIDKey))
	if !exists {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Unauthorized",
		})
		return
	}

	// 取得 LINE 官方帳號 ID
	lineOfficialAccountID := c.qrCodeService.GetLineOfficialAccountID()

	// 如果環境變數沒有設定，回傳預設的 LINE ID
	if lineOfficialAccountID == "" {
		lineOfficialAccountID = "timeledger"
	}

	// 產生 QR Code
	qrBytes, err := c.qrCodeService.GenerateLINEBindingQR(lineOfficialAccountID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to generate LINE binding QR code: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: "系統錯誤",
		})
		return
	}

	// 輸出 PNG 圖片
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", "inline; filename=line-binding-qr.png")
	ctx.Data(http.StatusOK, "image/png", qrBytes)
}

// GenerateVerificationCodeQR 產生包含驗證碼的 LINE 綁定 QR Code
// @Summary 產生包含驗證碼的 QR Code
// @Description 產生掃描後會自動帶入驗證碼文字的 QR Code
// @Tags LINE
// @Produce octet-stream
// @Security BearerAuth
// @Param code query string true "6位數驗證碼"
// @Success 200 {file} binary "PNG image"
// @Failure 400 {object} global.ApiResponse
// @Failure 401 {object} global.ApiResponse
// @Failure 500 {object} global.ApiResponse
// @Router /api/v1/admin/me/line/qrcode-with-code [GET]
func (c *LineBotController) GenerateVerificationCodeQR(ctx *gin.Context) {
	// 驗證管理員身份（從 JWT middleware 取得）
	_, exists := ctx.Get(string(global.UserIDKey))
	if !exists {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Unauthorized",
		})
		return
	}

	// 取得驗證碼
	code := ctx.Query("code")
	if len(code) != 6 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "驗證碼必須是6位數",
		})
		return
	}

	// 取得 LINE 官方帳號 ID
	lineOfficialAccountID := c.qrCodeService.GetLineOfficialAccountID()

	// 如果環境變數沒有設定，回傳預設的 LINE ID
	if lineOfficialAccountID == "" {
		lineOfficialAccountID = "timeledger"
	}

	// 產生 QR Code
	qrBytes, err := c.qrCodeService.GenerateVerificationCodeQR(lineOfficialAccountID, code)
	if err != nil {
		fmt.Printf("[ERROR] Failed to generate verification code QR code: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: "系統錯誤",
		})
		return
	}

	// 輸出 PNG 圖片
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", "inline; filename=line-verification-qr.png")
	ctx.Data(http.StatusOK, "image/png", qrBytes)
}
