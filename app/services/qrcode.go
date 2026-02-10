package services

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"image/png"

	"github.com/skip2/go-qrcode"
)

// QRCodeService QR Code 生成服務
type QRCodeService struct{}

// NewQRCodeService 建立 QRCodeService
func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

// GeneratePNG 生成 PNG 格式的 QR Code
func (s *QRCodeService) GeneratePNG(content string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256 // 預設大小
	}

	qrc, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// 使用 Image() 方法取得圖像，然後編碼為 PNG
	img := qrc.Image(size)

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code as PNG: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateLINEBindingQR 生成 LINE 官方帳號加好友 QR Code
// QR Code 內容為 LINE Deep Link，掃描後直接開啟官方帳號聊天視窗
func (s *QRCodeService) GenerateLINEBindingQR(lineOfficialAccountID string) ([]byte, error) {
	// LINE Deep Link 格式：line://nv/addContacts/{LINE_ID}
	// 這種方式可以讓用戶直接開啟 LINE 並跳轉到官方帳號的聊天視窗
	lineDeepLink := fmt.Sprintf("line://nv/addContacts/%s", lineOfficialAccountID)

	return s.GeneratePNG(lineDeepLink, 256)
}

// GenerateBindingURLQR 生成綁定 URL 的 QR Code
// 產生一個包含綁定頁面 URL 的 QR Code，用戶掃描後可直接開啟綁定頁面
func (s *QRCodeService) GenerateBindingURLQR(bindingURL string) ([]byte, error) {
	return s.GeneratePNG(bindingURL, 256)
}

// GenerateVerificationCodeQR 生成包含驗證碼的 LINE 綁定 QR Code
// 這種方式會產生一個 LINE Message API 的 URL，掃描後會開啟 LINE
// 並自動帶入驗證碼文字，用戶可以直接傳送
func (s *QRCodeService) GenerateVerificationCodeQR(lineOfficialAccountID, verificationCode string) ([]byte, error) {
	// 檢查 LINE ID 是否有效
	if lineOfficialAccountID == "" {
		return nil, fmt.Errorf("LINE官方帳號ID未設定，請聯繫系統管理員設定環境變數LINE_OFFICIAL_ACCOUNT_ID")
	}

	// 確保 ID 包含 @（LINE 官方帳號規範）
	formattedID := lineOfficialAccountID
	if !strings.HasPrefix(formattedID, "@") {
		formattedID = "@" + formattedID
	}

	// 使用更可靠的 LINE oaMessage URL 格式
	// https://line.me/R/oaMessage/{LINE_ID}/?{message}
	// 這種方式在行動裝置上能更穩定地開啟聊天視窗並填入內容
	url := fmt.Sprintf(
		"https://line.me/R/oaMessage/%s/?%s",
		formattedID,
		verificationCode,
	)

	return s.GeneratePNG(url, 256)
}

// GetLineOfficialAccountID 取得環境變數中的 LINE 官方帳號 ID
func (s *QRCodeService) GetLineOfficialAccountID() string {
	lineID := os.Getenv("LINE_OFFICIAL_ACCOUNT_ID")
	if lineID == "" {
		return ""
	}
	// 不再主動移除 @，交給 Generate 函數處理格式
	return lineID
}
