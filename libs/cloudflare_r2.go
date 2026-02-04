package libs

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"timeLedger/configs"
)

// R2StorageService Cloudflare R2 儲存服務（純 HTTP 實作）
type R2StorageService struct {
	config     *configs.Env
	httpClient *http.Client
}

// NewR2StorageService 建立 R2 儲存服務
func NewR2StorageService(env *configs.Env) (*R2StorageService, error) {
	if !env.CloudflareR2Enabled {
		return nil, nil
	}

	// 驗證配置
	if env.CloudflareR2AccountID == "" || env.CloudflareR2AccessKey == "" ||
		env.CloudflareR2SecretKey == "" || env.CloudflareR2BucketName == "" {
		return nil, fmt.Errorf("Cloudflare R2 configuration is incomplete")
	}

	return &R2StorageService{
		config: env,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// UploadFile 上傳檔案到 R2
func (s *R2StorageService) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	if s == nil {
		return "", fmt.Errorf("R2 storage service is not initialized")
	}

	// 生成唯一的檔案 key
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(filename)
	uniqueKey := fmt.Sprintf("certificates/%s_%s%s", timestamp, randomString(8), ext)

	// 讀取檔案內容
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// 建立 R2 API URL
	r2URL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s",
		s.config.CloudflareR2AccountID,
		s.config.CloudflareR2BucketName,
		uniqueKey)

	// 建立請求
	req, err := http.NewRequestWithContext(ctx, "PUT", r2URL, bytes.NewReader(fileContent))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 計算 payload hash
	payloadHash := sha256Hex(fileContent)

	// 設定 Header
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-amz-acl", "public-read")
	req.Header.Set("x-amz-content-sha256", payloadHash)

	// 產生 AWS Signature v4
	date := time.Now().UTC()
	amzDate := date.Format("20060102T150405Z")
	shortDate := date.Format("20060102")

	service := "s3"
	region := "auto"
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", shortDate, region, service)

	// 產生簽名
	signedHeaders := "content-type;host;x-amz-acl;x-amz-content-sha256;x-amz-date"

	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s.r2.cloudflarestorage.com\nx-amz-acl:public-read\nx-amz-content-sha256:%s\nx-amz-date:%s\n",
		contentType, s.config.CloudflareR2AccountID, payloadHash, amzDate)

	canonicalRequest := fmt.Sprintf("PUT\n/%s/%s\n\n%s\n%s\n%s",
		s.config.CloudflareR2BucketName, uniqueKey, signedHeaders, canonicalHeaders, payloadHash)

	requestHash := sha256Hex([]byte(canonicalRequest))

	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate, credentialScope, requestHash)

	// 產生簽名金鑰
	kDate := hmacSHA256([]byte("AWS4"+s.config.CloudflareR2SecretKey), []byte(shortDate))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))

	signature := hmacSHA256Hex(kSigning, []byte(stringToSign))

	// 設定 Authorization header
	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.config.CloudflareR2AccessKey, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)

	// 發送請求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("R2 upload failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	// 返回公開 URL 或 R2 URL
	if s.config.CloudflareR2PublicURL != "" {
		return fmt.Sprintf("%s/%s", s.config.CloudflareR2PublicURL, uniqueKey), nil
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s",
		s.config.CloudflareR2AccountID, s.config.CloudflareR2BucketName, uniqueKey), nil
}

// UploadFileFromPath 上傳本地檔案到 R2
func (s *R2StorageService) UploadFileFromPath(ctx context.Context, filePath string) (string, error) {
	if s == nil {
		return "", fmt.Errorf("R2 storage service is not initialized")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	contentType := GetContentType(filename)

	return s.UploadFile(ctx, file, filename, contentType)
}

// DeleteFile 從 R2 刪除檔案
func (s *R2StorageService) DeleteFile(ctx context.Context, fileURL string) error {
	if s == nil {
		return fmt.Errorf("R2 storage service is not initialized")
	}

	// 從 URL 提取 key
	key, err := extractKeyFromURL(fileURL, s.config.CloudflareR2PublicURL, s.config.CloudflareR2BucketName, s.config.CloudflareR2AccountID)
	if err != nil {
		return err
	}

	// 建立 R2 API URL
	r2URL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s",
		s.config.CloudflareR2AccountID, s.config.CloudflareR2BucketName, key)

	// 建立請求
	req, err := http.NewRequestWithContext(ctx, "DELETE", r2URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 產生簽名（簡化版本）
	date := time.Now().UTC()
	amzDate := date.Format("20060102T150405Z")
	shortDate := date.Format("20060102")

	service := "s3"
	region := "auto"
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", shortDate, region, service)

	signedHeaders := "host;x-amz-date"

	canonicalHeaders := fmt.Sprintf("host:%s.r2.cloudflarestorage.com\nx-amz-date:%s\n",
		s.config.CloudflareR2AccountID, amzDate)

	canonicalRequest := fmt.Sprintf("DELETE\n/%s/%s\n\n%s\n%s\nUNSIGNED-PAYLOAD",
		s.config.CloudflareR2BucketName, key, signedHeaders, canonicalHeaders)

	requestHash := sha256Hex([]byte(canonicalRequest))

	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate, credentialScope, requestHash)

	kDate := hmacSHA256([]byte("AWS4"+s.config.CloudflareR2SecretKey), []byte(shortDate))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))

	signature := hmacSHA256Hex(kSigning, []byte(stringToSign))

	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.config.CloudflareR2AccessKey, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("x-amz-date", amzDate)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete from R2: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("R2 delete failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// IsEnabled 檢查 R2 是否已啟用
func (s *R2StorageService) IsEnabled() bool {
	return s != nil && s.config != nil && s.config.CloudflareR2Enabled
}

// GetPublicURL 取得檔案的公開 URL
func (s *R2StorageService) GetPublicURL(filename string) string {
	if s == nil {
		return ""
	}

	if s.config.CloudflareR2PublicURL != "" {
		return fmt.Sprintf("%s/%s", s.config.CloudflareR2PublicURL, filename)
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s",
		s.config.CloudflareR2AccountID, s.config.CloudflareR2BucketName, filename)
}

// 輔助函數

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

func sha256Hex(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func computeSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func hmacSHA256Hex(key []byte, data []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func extractKeyFromURL(fileURL, publicURL, bucketName, accountID string) (string, error) {
	// 如果有 publicURL，從中提取
	if publicURL != "" {
		u, err := url.Parse(fileURL)
		if err == nil {
			return filepath.Base(u.Path), nil
		}
	}

	// 否則從標準 R2 URL 提取
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	return filepath.Base(u.Path), nil
}

func GetContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// LocalStorageService 本地儲存服務（回退方案）
type LocalStorageService struct {
	basePath string
}

// NewLocalStorageService 建立本地儲存服務
func NewLocalStorageService(basePath string) *LocalStorageService {
	return &LocalStorageService{
		basePath: basePath,
	}
}

// UploadFile 上傳檔案到本地
func (s *LocalStorageService) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// 生成唯一的檔案 key
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s_%s%s", timestamp, randomString(8), ext)
	filePath := filepath.Join(s.basePath, uniqueFilename)

	// 確保目錄存在
	if err := os.MkdirAll(s.basePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 建立檔案
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 複製內容
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return fmt.Sprintf("/uploads/certificates/%s", uniqueFilename), nil
}

// DeleteFile 刪除本地檔案
func (s *LocalStorageService) DeleteFile(ctx context.Context, fileURL string) error {
	// 從 URL 提取檔案名稱
	filename := filepath.Base(fileURL)
	filePath := filepath.Join(s.basePath, filename)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// UploadFileFromPath 上傳本地檔案
func (s *LocalStorageService) UploadFileFromPath(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	contentType := GetContentType(filename)

	return s.UploadFile(ctx, file, filename, contentType)
}
