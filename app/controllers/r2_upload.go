package controllers

import (
	"io"
	"net/http"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/libs"

	"github.com/gin-gonic/gin"
)

// R2TestController 測試 Cloudflare R2 上傳功能
type R2TestController struct {
	app *app.App
	r2  *libs.R2StorageService
}

// NewR2TestController 建立 R2 測試控制器
func NewR2TestController(app *app.App) *R2TestController {
	r2, _ := libs.NewR2StorageService(app.Env)
	return &R2TestController{
		app: app,
		r2:  r2,
	}
}

// UploadTestRequest 上傳測試請求
type UploadTestRequest struct {
	Category string `form:"category" binding:"required"`
}

// UploadTestResponse 上傳測試回應
type UploadTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	FileURL string `json:"file_url,omitempty"`
	FileKey string `json:"file_key,omitempty"`
	Storage string `json:"storage"`
	Error   string `json:"error,omitempty"`
}

// UploadTest 上傳測試 API
// @Summary 測試圖片上傳到 R2
// @Description 測試 Cloudflare R2 或本地儲存上傳功能
// @Tags R2 Test
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上傳的檔案"
// @Param category formData string true "分類 (certificates, backgrounds, avatars)"
// @Success 200 {object} UploadTestResponse
// @Failure 400 {object} UploadTestResponse
// @Router /api/test/upload [post]
func (c *R2TestController) UploadTest(ctx *gin.Context) {
	// 檢查 R2 狀態
	r2Enabled := c.r2 != nil && c.r2.IsEnabled()

	// 取得上傳的檔案
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, UploadTestResponse{
			Success: false,
			Message: "無法取得上傳檔案",
			Storage: getStorageType(r2Enabled),
			Error:   err.Error(),
		})
		return
	}
	defer file.Close()

	// 驗證檔案類型
	contentType := header.Header.Get("Content-Type")
	if !isAllowedContentType(contentType) {
		ctx.JSON(http.StatusBadRequest, UploadTestResponse{
			Success: false,
			Message: "不支援的檔案類型",
			Storage: getStorageType(r2Enabled),
			Error:   "允許的類型: image/jpeg, image/png, application/pdf",
		})
		return
	}

	// 取得分類
	category := ctx.PostForm("category")
	if category == "" {
		category = "test"
	}

	// 生成檔案名稱
	timestamp := time.Now().Format("20060102_150405")
	ext := getFileExtension(header.Filename)
	uniqueFilename := strings.ToLower(category) + "/" + timestamp + "_" + libs.RandomString(8) + ext

	// 根據設定選擇儲存方式
	var fileURL string
	var uploadErr error

	if r2Enabled {
		// 使用 R2 儲存
		fileURL, uploadErr = c.r2.UploadFile(ctx.Request.Context(), file, uniqueFilename, contentType)
		if uploadErr == nil {
			ctx.JSON(http.StatusOK, UploadTestResponse{
				Success: true,
				Message: "成功上傳到 Cloudflare R2",
				FileURL: fileURL,
				FileKey: uniqueFilename,
				Storage: "Cloudflare R2",
			})
			return
		}
		// 如果 R2 上傳失敗，回退到本地儲存
	}

	// 使用本地儲存（回退方案）
	localStorage := libs.NewLocalStorageService("./uploads")
	fileURL, uploadErr = localStorage.UploadFile(ctx.Request.Context(), file, uniqueFilename, contentType)
	if uploadErr != nil {
		ctx.JSON(http.StatusInternalServerError, UploadTestResponse{
			Success: false,
			Message: "上傳失敗",
			Storage: getStorageType(r2Enabled),
			Error:   uploadErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, UploadTestResponse{
		Success: true,
		Message: "成功上傳到本地儲存",
		FileURL: fileURL,
		FileKey: uniqueFilename,
		Storage: "Local Storage",
	})
}

// StatusTest R2 狀態測試 API
// @Summary 測試 R2 連線狀態
// @Description 檢查 Cloudflare R2 是否已正確配置
// @Tags R2 Test
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/test/r2-status [get]
func (c *R2TestController) StatusTest(ctx *gin.Context) {
	r2Enabled := c.r2 != nil && c.r2.IsEnabled()

	status := map[string]interface{}{
		"r2_enabled":   r2Enabled,
		"storage_type": getStorageType(r2Enabled),
		"config": map[string]interface{}{
			"account_id":  maskString(c.app.Env.CloudflareR2AccountID, 8),
			"bucket_name": c.app.Env.CloudflareR2BucketName,
			"public_url":  c.app.Env.CloudflareR2PublicURL != "",
		},
	}

	if r2Enabled {
		status["message"] = "Cloudflare R2 已啟用"
	} else {
		status["message"] = "使用本地儲存（請確認 .env 中的 CLOUDFLARE_R2_ENABLED=true 來啟用 R2）"
	}

	ctx.JSON(http.StatusOK, status)
}

// BatchUploadTest 批量上傳測試 API
// @Summary 測試批量圖片上傳
// @Description 測試多張圖片同時上傳到 R2
// @Tags R2 Test
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/test/upload-batch [post]
func (c *R2TestController) BatchUploadTest(ctx *gin.Context) {
	r2Enabled := c.r2 != nil && c.r2.IsEnabled()

	// 解析 multipart 表單
	err := ctx.Request.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "無法解析表單資料",
		})
		return
	}

	results := make([]map[string]interface{}, 0)
	successCount := 0
	failCount := 0

	// 遍歷所有檔案
	for _, fileHeaders := range ctx.Request.MultipartForm.File {
		for _, header := range fileHeaders {
			result := map[string]interface{}{
				"filename": header.Filename,
				"success":  false,
			}

			// 打開檔案
			file, err := header.Open()
			if err != nil {
				result["error"] = err.Error()
				results = append(results, result)
				failCount++
				continue
			}

			// 讀取內容
			content, err := io.ReadAll(file)
			if err != nil {
				result["error"] = err.Error()
				file.Close()
				results = append(results, result)
				failCount++
				continue
			}
			file.Close()

			// 生成檔案名稱
			timestamp := time.Now().Format("20060102_150405")
			ext := getFileExtension(header.Filename)
			uniqueFilename := "batch/" + timestamp + "_" + libs.RandomString(6) + ext

			// 上傳到 R2
			if r2Enabled {
				reader := strings.NewReader(string(content))
				fileURL, uploadErr := c.r2.UploadFile(ctx.Request.Context(), reader, uniqueFilename, header.Header.Get("Content-Type"))
				if uploadErr == nil {
					result["success"] = true
					result["file_url"] = fileURL
					result["storage"] = "Cloudflare R2"
					successCount++
					results = append(results, result)
					continue
				}
			}

			// 回退到本地儲存
			localStorage := libs.NewLocalStorageService("./uploads")
			reader := strings.NewReader(string(content))
			fileURL, uploadErr := localStorage.UploadFile(ctx.Request.Context(), reader, uniqueFilename, header.Header.Get("Content-Type"))
			if uploadErr == nil {
				result["success"] = true
				result["file_url"] = fileURL
				result["storage"] = "Local Storage"
			} else {
				result["error"] = uploadErr.Error()
				failCount++
			}
			results = append(results, result)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":       failCount == 0,
		"total":         len(results),
		"success_count": successCount,
		"fail_count":    failCount,
		"storage_type":  getStorageType(r2Enabled),
		"results":       results,
	})
}

// 輔助函數

func getStorageType(r2Enabled bool) string {
	if r2Enabled {
		return "Cloudflare R2"
	}
	return "Local Storage"
}

func isAllowedContentType(contentType string) bool {
	allowedTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/pdf",
	}

	for _, allowed := range allowedTypes {
		if strings.HasPrefix(contentType, allowed) {
			return true
		}
	}
	return false
}

func getFileExtension(filename string) string {
	ext := ""
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		ext = strings.ToLower(filename[idx:])
	}
	if ext == "" {
		ext = ".jpg"
	}
	return ext
}

// (已移除 redundant randomString，改用 libs.RandomString)

func maskString(s string, visibleChars int) string {
	if len(s) <= visibleChars {
		return s
	}
	return s[:visibleChars] + "***"
}
