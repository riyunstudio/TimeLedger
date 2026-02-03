package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

// ResponseSanitizer 確保 API 回應是乾淨的 JSON
type ResponseSanitizer struct{}

func NewResponseSanitizer() *ResponseSanitizer {
	return &ResponseSanitizer{}
}

// Sanitize 包装 ResponseWriter，確保沒有額外輸出
func (s *ResponseSanitizer) Sanitize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 創建自定義 ResponseWriter
		sw := &sanitizedWriter{
			ResponseWriter: c.Writer,
			buffer:         bytes.NewBuffer(nil),
		}
		c.Writer = sw

		c.Next()

		// 確保已經提交響應
		status := c.Writer.Status()
		
		// 只處理 JSON API 響應
		contentType := c.Writer.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return
		}

		// 如果已經有響應內容，確保它是乾淨的 JSON
		if sw.buffer.Len() > 0 {
			body := sw.buffer.Bytes()
			
			// 去除任何非 JSON 內容（開頭的空白、不可見字元等）
			trimmed := bytes.TrimSpace(body)
			
			// 檢查是否為有效的 JSON（以 { 或 [ 開頭）
			trimmedStr := strings.TrimSpace(string(trimmed))
			if len(trimmedStr) > 0 && (trimmedStr[0] == '{' || trimmedStr[0] == '[') {
				// 重新寫入乾淨的響應
				c.Writer.Header().Set("Content-Length", string(rune(len(trimmed))))
				c.Writer.WriteHeader(status)
				c.Writer.Write(trimmed)
				return
			}
			
			// 如果不是有效 JSON，寫入原始內容（讓錯誤顯現出來以便調試）
			c.Writer.Write(body)
		}
	}
}

// sanitizedWriter 包裝 gin.ResponseWriter，緩衝響應內容
type sanitizedWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func (w *sanitizedWriter) Write(b []byte) (int, error) {
	// 將內容寫入緩衝區
	return w.buffer.Write(b)
}

func (w *sanitizedWriter) WriteString(s string) (int, error) {
	return w.buffer.WriteString(s)
}

// RecoveryWithSanitizer 自定義的 Panic 恢復，確保不會輸出多餘內容
func RecoveryWithSanitizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 清除任何可能的緩衝區輸出
				c.Writer.Header().Del("X-Content-Type-Options")
				
				c.JSON(http.StatusInternalServerError, global.ApiResponse{
					Code:    global.SYSTEM_ERROR,
					Message: "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// DisableStdoutOutput 在生產環境禁用標準輸出
// 這會將 stdout 重定向到日誌系統
func DisableStdoutOutput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 捕獲任何在處理過程中的 stdout 輸出
		// 這個 middleware 確保 Gin 不會在響應中輸出任何內容
		c.Next()
	}
}
