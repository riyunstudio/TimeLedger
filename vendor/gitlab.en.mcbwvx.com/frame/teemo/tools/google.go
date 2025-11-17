package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

// 新增 2FA Key
func (tl *Tools) Generate2FASecret() string {
	return gotp.RandomSecret(16)
}

// 產生 2FA QRcode base64
func (tl *Tools) Generate2FAQrcode(account, secret string) (qr string, err error) {
	// 創建 TOTP 對象
	totp := gotp.NewDefaultTOTP(secret)

	// 生成 otpauth URL (讓 Google Authenticator 掃描)
	otpAuthURL := totp.ProvisioningUri(account, "Account")

	// 生成 QR Code 並保存到 buffer
	var qrCode *qrcode.QRCode
	qrCode, err = qrcode.New(otpAuthURL, qrcode.Medium)
	if err != nil {
		return qr, fmt.Errorf("產生 QRcode 錯誤, Err: %s", err.Error())
	}

	// 將 QR Code 寫入 buffer
	var buf bytes.Buffer
	if err = qrCode.Write(256, &buf); err != nil {
		return qr, fmt.Errorf("轉換 QRcode base64 錯誤, Err: %s", err.Error())
	}

	// 將 QR Code 轉換為 Base64 編碼
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// 驗證 2FA Key
func (tl *Tools) Verify2FASecret(secret, code string) bool {
	totp := gotp.NewDefaultTOTP(secret)
	return totp.Verify(code, time.Now().Unix())
}
