package services

import (
	"bytes"
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
	"image/png"
)

// QRCodeService QR Code 生成服務
type QRCodeService struct{}

// NewQRCodeService 建立 QRCodeService
func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

// GeneratePNG 生成 PNG 格式的 QR Code
// @param content QR Code 內容（URL 或文字）
// @param size QR Code 大小（像素）
// @return []byte PNG 圖片資料
// @return error 錯誤資訊
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
// @param lineOfficialAccountID LINE 官方帳號 ID（如 @timeledger）
// @return []byte PNG 圖片資料
// @return error 錯誤資訊
func (s *QRCodeService) GenerateLINEBindingQR(lineOfficialAccountID string) ([]byte, error) {
	// LINE Deep Link 格式：line://nv/addContacts/{LINE_ID}
	// 這種方式可以讓用戶直接開啟 LINE 並跳轉到官方帳號的聊天視窗
	lineDeepLink := fmt.Sprintf("line://nv/addContacts/%s", lineOfficialAccountID)

	return s.GeneratePNG(lineDeepLink, 256)
}

// GenerateBindingURLQR 生成綁定 URL 的 QR Code
// 產生一個包含綁定頁面 URL 的 QR Code，用戶掃描後可直接開啟綁定頁面
// @param bindingURL 綁定頁面的完整 URL
// @return []byte PNG 圖片資料
// @return error 錯誤資訊
func (s *QRCodeService) GenerateBindingURLQR(bindingURL string) ([]byte, error) {
	return s.GeneratePNG(bindingURL, 256)
}

// GenerateVerificationCodeQR 生成包含驗證碼的 LINE 綁定 QR Code
// 這種方式會產生一個 LINE Message API 的 URL，掃描後會開啟 LINE
// 並自動帶入驗證碼文字，用戶可以直接傳送
// @param lineOfficialAccountID LINE 官方帳號 ID
// @param verificationCode 6 位數驗證碼
// @return []byte PNG 圖片資料
// @return error 錯誤資訊
func (s *QRCodeService) GenerateVerificationCodeQR(lineOfficialAccountID, verificationCode string) ([]byte, error) {
	// 使用 LINE Message API 的 URL 格式
	// https://line.me/R/ti/p/{LINE_ID}?text={message}
	// 這種方式可以開啟 LINE 並自動填入訊息文字
	url := fmt.Sprintf(
		"https://line.me/R/ti/p/%s?text=%s",
		lineOfficialAccountID,
		verificationCode,
	)

	return s.GeneratePNG(url, 256)
}

// GetLineOfficialAccountID 取得環境變數中的 LINE 官方帳號 ID
// @return string LINE 官方帳號 ID
func (s *QRCodeService) GetLineOfficialAccountID() string {
	lineID := os.Getenv("LINE_OFFICIAL_ACCOUNT_ID")
	if lineID == "" {
		// 如果環境變數沒有設定，回傳空字串
		// 前端應該要有預設值
		return ""
	}
	// 移除 @ 符號（如果有的話）
	if len(lineID) > 0 && lineID[0] == '@' {
		lineID = lineID[1:]
	}
	return lineID
}
