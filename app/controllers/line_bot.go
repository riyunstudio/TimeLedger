package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/services"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

// LineBotController LINE Bot Webhook Controller
type LineBotController struct {
	app             *app.App
	logger          *services.ServiceLogger
	lineBotService  services.LineBotService
	qrCodeService   *services.QRCodeService
	adminService    *services.AdminUserService
	templateService services.LineBotTemplateService
}

// NewLineBotController 建立 LINE Bot Controller
func NewLineBotController(app *app.App) *LineBotController {
	return &LineBotController{
		app:             app,
		logger:          services.NewServiceLogger("LineBotController"),
		lineBotService:  services.NewLineBotService(app),
		qrCodeService:   services.NewQRCodeService(),
		adminService:    services.NewAdminUserService(app),
		templateService: services.NewLineBotTemplateService(app.Env.FrontendBaseURL),
	}
}

// LINEWebhookRequest LINE Webhook 請求結構
type LINEWebhookRequest struct {
	Destination string             `json:"destination"`
	Events      []LINEWebhookEvent `json:"events"`
}

// LINEWebhookEvent LINE Webhook 事件
type LINEWebhookEvent struct {
	Type       string           `json:"type"`
	Mode       string           `json:"mode"`
	Timestamp  int64            `json:"timestamp"`
	Source     LINEEventSource  `json:"source"`
	ReplyToken string           `json:"replyToken,omitempty"`
	Message    LINEEventMessage `json:"message,omitempty"`
}

// LINEEventSource 事件來源
type LINEEventSource struct {
	Type   string `json:"type"`
	UserID string `json:"userId,omitempty"`
}

// LINEEventMessage 事件訊息
type LINEEventMessage struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Text       string `json:"text,omitempty"`
	QuoteToken string `json:"quoteToken,omitempty"`
}

// HandleWebhook 處理 LINE Webhook
func (c *LineBotController) HandleWebhook(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		c.logger.Error("failed to read webhook body", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// 驗證簽名
	signature := ctx.GetHeader("X-Line-Signature")
	if !c.lineBotService.VerifySignature(body, signature) {
		c.logger.Warn("invalid LINE signature")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// 解析請求
	var webhookReq LINEWebhookRequest
	if err := json.Unmarshal(body, &webhookReq); err != nil {
		c.logger.Error("failed to parse webhook request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
		return
	}

	// 處理每個事件 - 每個事件由獨立的 goroutine 處理
	for _, event := range webhookReq.Events {
		go c.handleEvent(ctx, &event)
	}

	// 立即返回 200 OK，goroutine 會在背景繼續處理
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleEvent 處理單個事件
// 使用 gin.Context 來獲取 replyToken 等資訊
// 在內部為資料庫操作創建不會被取消的 context
func (c *LineBotController) handleEvent(gctx *gin.Context, event *LINEWebhookEvent) {
	// 為資料庫操作建立不會被取消的上下文
	// 避免 HTTP 請求結束後 goroutine 中的資料庫操作被取消
	dbCtx := context.WithoutCancel(gctx.Request.Context())

	switch event.Type {
	case "message":
		c.handleMessageEvent(dbCtx, event)
	case "follow":
		c.handleFollowEvent(dbCtx, event)
	case "unfollow":
		c.handleUnfollowEvent(dbCtx, event)
	default:
		c.logger.Debug("unhandled event type", "event_type", event.Type)
	}
}

// handleMessageEvent 處理訊息事件
func (c *LineBotController) handleMessageEvent(ctx context.Context, event *LINEWebhookEvent) {
	if event.Message.Type != "text" {
		return
	}

	text := event.Message.Text
	userID := event.Source.UserID

	// 【監控日誌】記錄用戶互動
	identity, _ := c.lineBotService.GetCombinedIdentity(userID)
	c.logger.Info("line_webhook_message",
		"user_id", userID,
		"primary_role", identity.PrimaryRole,
		"message_type", "text",
		"message_preview", truncateString(text, 50),
	)

	// 處理驗證碼（6位數大寫字母數字）
	// 注意：這裡使用 context 傳遞，但 processBindingCode 內部會調用 lineBotService.ReplyMessage
	// 由於我們已經在 goroutine 中，回覆消息不會受到 HTTP 請求結束的影響
	if len(text) == 6 && isValidBindingCode(text) {
		c.processBindingCode(ctx, text, userID, event.ReplyToken)
		return
	}

	// 處理關鍵字
	switch text {
	case "綁定", "bind", "Bind":
		c.sendBindingInstructions(ctx, event.ReplyToken)
	case "幫助", "幫我", "help", "Help":
		c.sendHelpMessage(ctx, event.ReplyToken)
	case "狀態", "status", "Status":
		c.sendStatusMessage(ctx, event.ReplyToken, userID)
	case "解除綁定", "unbind", "Unbind":
		c.sendUnbindInstructions(ctx, event.ReplyToken)
	case "了解更多", "更多", "more", "More":
		c.sendMoreInfoMessage(ctx, event.ReplyToken)
	case "稍後綁定", "稍後再說":
		c.sendAckMessage(ctx, event.ReplyToken)
	case "課表", "我的課表", "今日課表", "schedule", "Schedule":
		c.sendScheduleMessage(ctx, event.ReplyToken, userID)
	case "明天課表", "明日課表":
		c.sendScheduleMessage(ctx, event.ReplyToken, userID, true)
	default:
		c.sendDefaultResponse(ctx, event.ReplyToken)
	}
}

// handleFollowEvent 處理加入好友事件
func (c *LineBotController) handleFollowEvent(ctx context.Context, event *LINEWebhookEvent) {
	userID := event.Source.UserID

	// 【監控日誌】記錄用戶關注
	identity, _ := c.lineBotService.GetCombinedIdentity(userID)
	c.logger.Info("line_webhook_follow",
		"user_id", userID,
		"primary_role", identity.PrimaryRole,
		"is_bound_admin", identity.PrimaryRole == "ADMIN",
		"is_bound_teacher", identity.PrimaryRole == "TEACHER",
	)

	// 1. 檢查是否為已綁定的管理員
	adminStatus, _, _ := c.adminService.GetLINEBindingStatusByLineUserID(ctx, userID)
	if adminStatus != nil && adminStatus.IsBound {
		// 已綁定的管理員
		centerName := "TimeLedger"

		welcomeFlex := c.templateService.GetWelcomeAdminTemplate(&models.AdminUser{
			LineUserID: userID,
			Role:       adminStatus.Role,
		}, centerName)

		err := c.lineBotService.ReplyFlexMessage(ctx, event.ReplyToken, "歡迎回來！", welcomeFlex)
		if err == nil {
			return // 成功發送 Flex Message
		}
		c.logger.Warn("failed to send admin welcome flex, using text", "error", err)
	}

	// 2. 檢查是否為老師（通過 LINE User ID）
	// 老師的歡迎訊息
	welcomeFlex := c.templateService.GetWelcomeTeacherTemplate(&models.Teacher{
		LineUserID: userID,
	}, "TimeLedger")

	err := c.lineBotService.ReplyFlexMessage(ctx, event.ReplyToken, "歡迎加入 TimeLedger！", welcomeFlex)
	if err == nil {
		return // 成功發送老師歡迎訊息
	}

	// 3. 如果 Flex Message 失敗，發送通用文字訊息
	c.logger.Error("failed to send welcome flex message", "error", err)
	welcomeMessage := map[string]interface{}{
		"type": "text",
		"text": "👋 您好！歡迎加入 TimeLedger！\n\n" +
			"TimeLedger 是教師中心化多據點排課平台，\n" +
			"讓您可以輕鬆管理課表、提交例外申請。\n\n" +
			"如需使用，請透過 LIFF 頁面登入。",
	}
	c.lineBotService.ReplyMessage(ctx, event.ReplyToken, welcomeMessage)
}

// handleUnfollowEvent 處理封鎖/取消好友事件
func (c *LineBotController) handleUnfollowEvent(ctx context.Context, event *LINEWebhookEvent) {
	userID := event.Source.UserID

	// 【監控日誌】記錄用戶取消關注
	c.logger.Info("line_webhook_unfollow",
		"user_id", userID,
	)
}

// processBindingCode 處理綁定驗證碼
func (c *LineBotController) processBindingCode(ctx context.Context, code string, userID string, replyToken string) {
	adminID, eInfo, err := c.adminService.VerifyLINEBinding(ctx, code, userID)
	if err != nil {
		c.logger.Error("failed to verify binding code", "error", err)
		errorMsg := "❌ 綁定失敗，驗證碼錯誤或已過期。"
		if eInfo != nil {
			if eInfo.Code == 90004 {
				errorMsg = "❌ 驗證碼已過期，請至後台重新產生。"
			}
		}
		c.lineBotService.ReplyMessage(ctx, replyToken, map[string]interface{}{
			"type": "text",
			"text": errorMsg,
		})
		return
	}

	// 綁定成功
	c.lineBotService.ReplyMessage(ctx, replyToken, map[string]interface{}{
		"type": "text",
		"text": "✅ 綁定成功！\n\n" +
			"您將會收到：\n" +
			"🔔 老師提交例外申請的通知\n" +
			"🔔 審核結果通知\n\n" +
			"如需調整通知設定，請至後台「設定」→「通知設定」。",
	})

	// 發送歡迎訊息（異步）
	go func() {
		welcomeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := c.adminService.SendWelcomeMessageIfNeeded(welcomeCtx, adminID); err != nil {
			c.logger.Error("failed to send welcome message after binding", "admin_id", adminID, "error", err)
		}
	}()
}

// sendBindingInstructions 發送綁定說明
func (c *LineBotController) sendBindingInstructions(ctx context.Context, replyToken string) {
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
	c.lineBotService.ReplyMessage(ctx, replyToken, message)
}

// sendHelpMessage 發送幫助訊息
func (c *LineBotController) sendHelpMessage(ctx context.Context, replyToken string) {
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
	c.lineBotService.ReplyMessage(ctx, replyToken, message)
}

// sendStatusMessage 發送狀態訊息
func (c *LineBotController) sendStatusMessage(ctx context.Context, replyToken string, userID string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "📊 狀態查詢：\n\n" +
			"您的 LINE 帳號已與 TimeLedger 綁定。\n\n" +
			"如需調整設定，請至管理後台。",
	}
	c.lineBotService.ReplyMessage(ctx, replyToken, message)
}

// sendUnbindInstructions 發送解除綁定說明
func (c *LineBotController) sendUnbindInstructions(ctx context.Context, replyToken string) {
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
	c.lineBotService.ReplyMessage(ctx, replyToken, message)
}

// sendMoreInfoMessage 發送更多資訊
func (c *LineBotController) sendMoreInfoMessage(ctx context.Context, replyToken string) {
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
	c.lineBotService.ReplyMessage(ctx, replyToken, message)
}

// sendAckMessage 發送確認訊息
func (c *LineBotController) sendAckMessage(ctx context.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "ℹ️ 了解！\n\n" +
			"您可以稍後再進行綁定。\n" +
			"當您準備好時，輸入「綁定」即可開始流程。",
	}
	c.lineBotService.ReplyMessage(ctx, replyToken, message)
}

// sendDefaultResponse 發送預設回應
func (c *LineBotController) sendDefaultResponse(ctx context.Context, replyToken string) {
	message := map[string]interface{}{
		"type": "text",
		"text": "🤔 我不太理解您的意思。\n\n" +
			"輸入「幫助」查看可用指令。",
	}
	c.lineBotService.ReplyMessage(ctx, replyToken, message)
}

// sendScheduleMessage 發送課表訊息
// isTomorrow: true 表示查詢明天課表，false 表示查詢今天課表
func (c *LineBotController) sendScheduleMessage(ctx context.Context, replyToken string, userID string, isTomorrow ...bool) {
	// 計算目標日期
	targetDate := time.Now()
	if len(isTomorrow) > 0 && isTomorrow[0] {
		targetDate = targetDate.AddDate(0, 0, 1)
	}

	// 取得當日課表
	agendaItems, err := c.lineBotService.GetAggregatedAgenda(userID, &targetDate)
	if err != nil {
		c.logger.Error("failed to get aggregated agenda", "error", err, "user_id", userID)
		errorMsg := map[string]interface{}{
			"type": "text",
			"text": "❌ 取得課表失敗，請稍後再試。\n\n" +
				"如有問題，請聯繫系統管理員。",
		}
		c.lineBotService.ReplyMessage(ctx, replyToken, errorMsg)
		return
	}

	// 取得用戶名稱（用於標題顯示）
	userName := ""
	identity, err := c.lineBotService.GetCombinedIdentity(userID)
	if err == nil {
		if identity.TeacherProfile != nil {
			userName = identity.TeacherProfile.Name
		} else if len(identity.AdminProfiles) > 0 {
			userName = identity.AdminProfiles[0].Name
		}
	}
	if userName == "" {
		userName = "您"
	}

	// 產生 Flex Message
	flexContent := c.templateService.GenerateAgendaFlex(agendaItems, targetDate, userName)

	// 發送 Flex Message
	err = c.lineBotService.ReplyFlexMessage(ctx, replyToken, "今日課表", flexContent)
	if err != nil {
		c.logger.Error("failed to send schedule flex message", "error", err)
		// Flex Message 失敗時，發送文字訊息
		if len(agendaItems) == 0 {
			weekdayStr := getWeekdayChinese(targetDate)
			dateStr := targetDate.Format("1月2日")
			fallbackMsg := map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("📅 %s (%s)\n\n"+

					"目前沒有課表。\n\n"+
					"💡 您可以透過 LIFF 頁面查看完整課表。", dateStr, weekdayStr),
			}
			c.lineBotService.ReplyMessage(ctx, replyToken, fallbackMsg)
		} else {
			fallbackMsg := c.buildScheduleFallbackMessage(agendaItems, targetDate)
			c.lineBotService.ReplyMessage(ctx, replyToken, fallbackMsg)
		}
	}
}

// buildScheduleFallbackMessage 建立課表文字回覆（當 Flex Message 失敗時使用）
func (c *LineBotController) buildScheduleFallbackMessage(agendaItems []services.AgendaItem, targetDate time.Time) map[string]interface{} {
	dateStr := targetDate.Format("1月2日")
	weekdayStr := getWeekdayChinese(targetDate)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📅 %s (%s)\n\n", dateStr, weekdayStr))

	// 分組顯示
	var centerItems, personalItems []services.AgendaItem
	for _, item := range agendaItems {
		if item.SourceType == services.AgendaSourceTypeCenter {
			centerItems = append(centerItems, item)
		} else {
			personalItems = append(personalItems, item)
		}
	}

	// 顯示中心課表
	if len(centerItems) > 0 {
		sb.WriteString("🏢 中心課表\n")
		for _, item := range centerItems {
			sb.WriteString(fmt.Sprintf("  %s │ %s (%s)\n", item.Time, item.Title, item.SourceName))
		}
		if len(personalItems) > 0 {
			sb.WriteString("\n")
		}
	}

	// 顯示個人行程
	if len(personalItems) > 0 {
		sb.WriteString("📌 個人行程\n")
		for _, item := range personalItems {
			sb.WriteString(fmt.Sprintf("  %s │ %s\n", item.Time, item.Title))
		}
	}

	sb.WriteString("\n💡 輸入「課表」查看明日課表")

	return map[string]interface{}{
		"type": "text",
		"text": sb.String(),
	}
}

// getWeekdayChinese 取得星期幾的中文名稱
func getWeekdayChinese(date time.Time) string {
	weekdays := []string{"週日", "週一", "週二", "週三", "週四", "週五", "週六"}
	return weekdays[date.Weekday()]
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

// truncateString 截斷字串（用於日誌顯示）
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	return s[:maxLen-3] + "..."
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
	_, exists := ctx.Get(global.UserIDKey)
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
		c.logger.Error("failed to generate LINE binding QR code", "error", err)
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
	_, exists := ctx.Get(global.UserIDKey)
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
	c.logger.Debug("generating verification QR code",
		"line_id", lineOfficialAccountID,
		"verification_code", code,
	)

	// 如果環境變數沒有設定，回傳預設的 LINE ID
	if lineOfficialAccountID == "" {
		lineOfficialAccountID = "timeledger"
	}

	// 產生 QR Code
	qrBytes, err := c.qrCodeService.GenerateVerificationCodeQR(lineOfficialAccountID, code)
	if err != nil {
		c.logger.Error("failed to generate verification code QR code", "error", err)
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: err.Error(),
		})
		return
	}

	c.logger.Debug("QR code generated successfully", "size", len(qrBytes))

	// 將 QR Code 轉換為 base64 字串返回，避免二進制流被代理截斷
	base64Image := base64.StdEncoding.EncodeToString(qrBytes)

	// 回傳 JSON，包含完整的 data URL
	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "success",
		Datas: map[string]string{
			"image": "data:image/png;base64," + base64Image,
		},
	})
}
