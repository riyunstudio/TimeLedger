package services

import (
	"os"
	"testing"
)

// TestQRCodeService_GeneratePNG 測試 GeneratePNG 基本功能
func TestQRCodeService_GeneratePNG(t *testing.T) {
	service := NewQRCodeService()

	tests := []struct {
		name    string
		content string
		size    int
		wantErr bool
	}{
		{
			name:    "產生基本 URL QR Code",
			content: "https://example.com",
			size:    256,
			wantErr: false,
		},
		{
			name:    "產生簡短文字 QR Code",
			content: "Hello World",
			size:    128,
			wantErr: false,
		},
		{
			name:    "產生空內容 QR Code（應該失敗）",
			content: "",
			size:    256,
			wantErr: true,
		},
		{
			name:    "使用預設大小",
			content: "https://timeledger.app",
			size:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := service.GeneratePNG(tt.content, tt.size)

			if tt.wantErr {
				if err == nil {
					t.Error("GeneratePNG 應該返回錯誤，但沒有")
				}
				return
			}

			if err != nil {
				t.Errorf("GeneratePNG 發生錯誤: %v", err)
				return
			}

			// 驗證 PNG 資料
			if len(data) == 0 {
				t.Error("GeneratePNG 返回空的資料")
				return
			}

			// 驗證 PNG 檔案頭（PNG 檔案以 8 字節的 PNG signature 開頭）
			expectedPNGHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
			for i, b := range expectedPNGHeader {
				if data[i] != b {
					t.Errorf("PNG 檔案頭不正確，索引 %d: 期望 %x，得到 %x", i, b, data[i])
					return
				}
			}
		})
	}
}

// TestQRCodeService_GenerateLINEBindingQR 測試產生 LINE 官方帳號 QR Code
func TestQRCodeService_GenerateLINEBindingQR(t *testing.T) {
	service := NewQRCodeService()

	tests := []struct {
		name               string
		lineOfficialAccountID string
		wantErr            bool
	}{
		{
			name:               "產生標準 LINE ID QR Code",
			lineOfficialAccountID: "@timeledger",
			wantErr:            false,
		},
		{
			name:               "產生無 @ 符號的 LINE ID QR Code",
			lineOfficialAccountID: "timeledger",
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := service.GenerateLINEBindingQR(tt.lineOfficialAccountID)

			if tt.wantErr {
				if err == nil {
					t.Error("GenerateLINEBindingQR 應該返回錯誤，但沒有")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateLINEBindingQR 發生錯誤: %v", err)
				return
			}

			// 驗證 PNG 資料
			if len(data) == 0 {
				t.Error("GenerateLINEBindingQR 返回空的資料")
				return
			}

			// 驗證 PNG 檔案頭
			expectedPNGHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
			for i, b := range expectedPNGHeader {
				if data[i] != b {
					t.Errorf("PNG 檔案頭不正確，索引 %d: 期望 %x，得到 %x", i, b, data[i])
					return
				}
			}
		})
	}
}

// TestQRCodeService_GenerateVerificationCodeQR 測試產生含驗證碼的 QR Code
func TestQRCodeService_GenerateVerificationCodeQR(t *testing.T) {
	service := NewQRCodeService()

	tests := []struct {
		name               string
		lineOfficialAccountID string
		verificationCode   string
		wantErr            bool
	}{
		{
			name:               "產生含驗證碼的 QR Code",
			lineOfficialAccountID: "@timeledger",
			verificationCode:   "ABC123",
			wantErr:            false,
		},
		{
			name:               "產生含數字驗證碼的 QR Code",
			lineOfficialAccountID: "timeledger",
			verificationCode:   "123456",
			wantErr:            false,
		},
		{
			name:               "產生含混合驗證碼的 QR Code",
			lineOfficialAccountID: "@test",
			verificationCode:   "X9Y2Z1",
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := service.GenerateVerificationCodeQR(tt.lineOfficialAccountID, tt.verificationCode)

			if tt.wantErr {
				if err == nil {
					t.Error("GenerateVerificationCodeQR 應該返回錯誤，但沒有")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateVerificationCodeQR 發生錯誤: %v", err)
				return
			}

			// 驗證 PNG 資料
			if len(data) == 0 {
				t.Error("GenerateVerificationCodeQR 返回空的資料")
				return
			}

			// 驗證 PNG 檔案頭
			expectedPNGHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
			for i, b := range expectedPNGHeader {
				if data[i] != b {
					t.Errorf("PNG 檔案頭不正確，索引 %d: 期望 %x，得到 %x", i, b, data[i])
					return
				}
			}
		})
	}
}

// TestQRCodeService_GetLineOfficialAccountID 測試取得 LINE 官方帳號 ID
func TestQRCodeService_GetLineOfficialAccountID(t *testing.T) {
	service := NewQRCodeService()

	// 測試環境變數未設定的情況
	originalValue := os.Getenv("LINE_OFFICIAL_ACCOUNT_ID")
	os.Unsetenv("LINE_OFFICIAL_ACCOUNT_ID")
	defer func() {
		if originalValue != "" {
			os.Setenv("LINE_OFFICIAL_ACCOUNT_ID", originalValue)
		}
	}()

	t.Run("環境變數未設定時回傳空字串", func(t *testing.T) {
		lineID := service.GetLineOfficialAccountID()
		if lineID != "" {
			t.Errorf("期望空字串，得到 %s", lineID)
		}
	})

	// 測試環境變數設定為帶 @ 的情況
	os.Setenv("LINE_OFFICIAL_ACCOUNT_ID", "@timeledger")
	t.Run("環境變數帶 @ 符號時應該移除 @", func(t *testing.T) {
		lineID := service.GetLineOfficialAccountID()
		if lineID == "timeledger" {
			t.Log("成功移除 @ 符號")
		} else if lineID == "@timeledger" {
			t.Error("沒有移除 @ 符號")
		}
	})

	// 測試環境變數設定為不帶 @ 的情況
	os.Setenv("LINE_OFFICIAL_ACCOUNT_ID", "test123")
	t.Run("環境變數不帶 @ 符號時保持不變", func(t *testing.T) {
		lineID := service.GetLineOfficialAccountID()
		if lineID != "test123" {
			t.Errorf("期望 test123，得到 %s", lineID)
		}
	})
}

// TestQRCodeService_GenerateBindingURLQR 測試產生綁定 URL QR Code
func TestQRCodeService_GenerateBindingURLQR(t *testing.T) {
	service := NewQRCodeService()

	tests := []struct {
		name      string
		bindingURL string
		wantErr   bool
	}{
		{
			name:      "產生綁定 URL QR Code",
			bindingURL: "https://timeledger.app/admin/line-bind",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := service.GenerateBindingURLQR(tt.bindingURL)

			if tt.wantErr {
				if err == nil {
					t.Error("GenerateBindingURLQR 應該返回錯誤，但沒有")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateBindingURLQR 發生錯誤: %v", err)
				return
			}

			// 驗證 PNG 資料
			if len(data) == 0 {
				t.Error("GenerateBindingURLQR 返回空的資料")
				return
			}

			// 驗證 PNG 檔案頭
			expectedPNGHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
			for i, b := range expectedPNGHeader {
				if data[i] != b {
					t.Errorf("PNG 檔案頭不正確，索引 %d: 期望 %x，得到 %x", i, b, data[i])
					return
				}
			}
		})
	}
}
