package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"

	"github.com/google/uuid"
)

// ImageService 處理課表圖片生成
type ImageService struct {
	BaseService
	app              *app.App
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	teacherRepo      *repositories.TeacherRepository
	centerRepo       *repositories.CenterRepository
}

// NewImageService 建立 ImageService 實例
func NewImageService(app *app.App) *ImageService {
	return &ImageService{
		app:              app,
		scheduleRuleRepo: repositories.NewScheduleRuleRepository(app),
		teacherRepo:      repositories.NewTeacherRepository(app),
		centerRepo:       repositories.NewCenterRepository(app),
	}
}

// ImageConfig 圖片生成配置
type ImageConfig struct {
	TeacherID       uint       // 老師 ID
	CenterID        uint       // 中心 ID
	StartDate       time.Time  // 開始日期
	EndDate         time.Time  // 結束日期
	Width           int        // 圖片寬度 (預設 800)
	Height          int        // 圖片高度 (預設 600)
	BackgroundImage string     // 自訂背景圖 URL 或本地路徑
	Theme           ImageTheme // 主題配色
	Title           string     // 標題文字
	IncludePersonal bool       // 是否包含個人行程
}

// ImageTheme 主題配色
type ImageTheme struct {
	Background       string // 背景色
	HeaderBackground string // 標題背景色
	HeaderText       string // 標題文字色
	SessionBg        string // 課程區塊背景色
	SessionText      string // 課程文字色
	TimeText         string // 時間文字色
	GridLine         string // 網格線顏色
	PersonalBg       string // 個人行程背景色
	PersonalText     string // 個人行程文字色
}

// DefaultTheme 預設主題
var DefaultTheme = ImageTheme{
	Background:       "#FFFFFF",
	HeaderBackground: "#4F46E5",
	HeaderText:       "#FFFFFF",
	SessionBg:        "#EEF2FF",
	SessionText:      "#1E1B4B",
	TimeText:         "#6B7280",
	GridLine:         "#E5E7EB",
	PersonalBg:       "#FEF3C7",
	PersonalText:     "#92400E",
}

// ScheduleImageData 課表圖片資料
type ScheduleImageData struct {
	Date        string // 日期 (YYYY-MM-DD)
	Weekday     string // 星期幾
	StartTime   string // 開始時間
	EndTime     string // 結束時間
	CourseName  string // 課程名稱
	TeacherName string // 老師名稱
	RoomName    string // 教室名稱
	IsPersonal  bool   // 是否為個人行程
	Title       string // 行程標題
	Color       string // 自訂顏色
}

// GenerateScheduleImage 生成課表圖片
func (s *ImageService) GenerateScheduleImage(ctx context.Context, config *ImageConfig) ([]byte, error) {
	// 設定預設值
	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}
	if config.Title == "" {
		config.Title = "課表"
	}
	if config.Theme.Background == "" {
		config.Theme = DefaultTheme
	}

	// 建立圖片
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	// 繪製背景
	backgroundColor := s.parseHexColor(config.Theme.Background)
	draw.Draw(img, img.Bounds(), s.colorToImage(backgroundColor), image.Point{}, draw.Src)

	// 取得課表資料
	scheduleData, err := s.getScheduleData(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule data: %w", err)
	}

	// 如果有自訂背景圖，疊加到背景上
	if config.BackgroundImage != "" {
		if err := s.overlayBackgroundImage(img, config.BackgroundImage); err != nil {
			// 靜默忽略背景圖疊加錯誤
			_ = err
		}
	}

	// 繪製標題區域
	s.drawHeader(img, config)

	// 繪製課表內容
	s.drawSchedule(img, scheduleData, config)

	// 轉換為 JPEG
	return s.imageToJPEG(img)
}

// getScheduleData 取得課表資料
func (s *ImageService) getScheduleData(ctx context.Context, config *ImageConfig) ([]ScheduleImageData, error) {
	scheduleQueryService := NewScheduleQueryService(s.app)
	schedule, err := scheduleQueryService.GetTeacherSchedule(ctx, config.TeacherID, config.StartDate, config.EndDate)
	if err != nil {
		return nil, err
	}

	data := make([]ScheduleImageData, 0)
	for _, item := range schedule {
		scheduleItem := ScheduleImageData{
			Date:        item.Date,
			Weekday:     s.formatWeekday(item.Date),
			StartTime:   item.StartTime,
			EndTime:     item.EndTime,
			CourseName:  item.Title,
			TeacherName: "", // TeacherScheduleItem 沒有 TeacherName 欄位
			RoomName:    "", // TeacherScheduleItem 沒有 RoomName 欄位，需要另外查詢
			IsPersonal:  item.Type == "personal_event",
			Title:       item.Title,
		}

		if item.Type == "personal_event" {
			scheduleItem.Color = config.Theme.PersonalBg
		} else {
			scheduleItem.Color = config.Theme.SessionBg
		}

		data = append(data, scheduleItem)
	}

	return data, nil
}

// formatWeekday 從日期字串取得星期幾
func (s *ImageService) formatWeekday(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ""
	}
	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	return weekdays[date.Weekday()]
}

// drawHeader 繪製標題區域
func (s *ImageService) drawHeader(img *image.RGBA, config *ImageConfig) {
	headerHeight := 80

	// 繪製標題背景
	headerColor := s.parseHexColor(config.Theme.HeaderBackground)
	headerRect := image.Rect(0, 0, config.Width, headerHeight)
	draw.Draw(img, headerRect, s.colorToImage(headerColor), image.Point{}, draw.Src)

	// 繪製標題文字
	titleColor := s.colorToImage(s.parseHexColor(config.Theme.HeaderText))
	titlePoint := image.Point{
		X: config.Width / 2,
		Y: headerHeight/2 + 10,
	}
	s.drawText(img, config.Title, titlePoint, titleColor, 32, true)

	// 繪製日期範圍
	dateRange := fmt.Sprintf("%s ~ %s",
		config.StartDate.Format("2006/01/02"),
		config.EndDate.Format("2006/01/02"))
	dateColor := s.colorToImage(s.parseHexColor("#FFFFFF"))
	datePoint := image.Point{
		X: config.Width / 2,
		Y: headerHeight/2 + 40,
	}
	s.drawText(img, dateRange, datePoint, dateColor, 14, true)
}

// drawSchedule 繪製課表內容
func (s *ImageService) drawSchedule(img *image.RGBA, scheduleData []ScheduleImageData, config *ImageConfig) {
	headerHeight := 80
	padding := 20
	itemHeight := 60
	startY := headerHeight + padding

	// 繪製網格線
	gridColor := s.parseHexColor(config.Theme.GridLine)
	for y := startY; y < config.Height; y += itemHeight {
		lineRect := image.Rect(padding, y, config.Width-padding, y+1)
		draw.Draw(img, lineRect, s.colorToImage(gridColor), image.Point{}, draw.Src)
	}

	// 繪製每個課表項目
	for i, item := range scheduleData {
		if i >= (config.Height-headerHeight-padding)/itemHeight-1 {
			break // 超過可容納數量
		}

		y := startY + i*itemHeight

		// 繪製項目背景
		itemBgColor := s.parseHexColor(item.Color)
		itemRect := image.Rect(padding, y+5, config.Width-padding, y+itemHeight-5)
		draw.Draw(img, itemRect, s.colorToImage(itemBgColor), image.Point{}, draw.Src)

		// 繪製日期和星期
		dateColorImg := s.colorToImage(s.parseHexColor(config.Theme.TimeText))
		datePoint := image.Point{
			X: padding + 15,
			Y: y + itemHeight/2 + 5,
		}
		s.drawText(img, fmt.Sprintf("%s (%s)", item.Date, item.Weekday), datePoint, dateColorImg, 12, false)

		// 繪製時間
		timeText := fmt.Sprintf("%s - %s", item.StartTime, item.EndTime)
		timePoint := image.Point{
			X: padding + 150,
			Y: y + itemHeight/2 + 5,
		}
		s.drawText(img, timeText, timePoint, dateColorImg, 12, false)

		// 繪製課程名稱
		nameColorImg := s.colorToImage(s.parseHexColor(config.Theme.SessionText))
		namePoint := image.Point{
			X: padding + 280,
			Y: y + itemHeight/2 + 5,
		}
		s.drawText(img, item.CourseName, namePoint, nameColorImg, 14, false)

		// 繪製老師名稱
		if item.TeacherName != "" {
			teacherColorImg := s.colorToImage(s.parseHexColor(config.Theme.TimeText))
			teacherPoint := image.Point{
				X: padding + 500,
				Y: y + itemHeight/2 + 5,
			}
			s.drawText(img, "老師: "+item.TeacherName, teacherPoint, teacherColorImg, 12, false)
		}

		// 繪製教室名稱
		if item.RoomName != "" {
			roomColorImg := s.colorToImage(s.parseHexColor(config.Theme.TimeText))
			roomPoint := image.Point{
				X: padding + 650,
				Y: y + itemHeight/2 + 5,
			}
			s.drawText(img, item.RoomName, roomPoint, roomColorImg, 12, false)
		}
	}
}

// drawText 繪製文字（簡化版本，實際應使用 font/gofont）
func (s *ImageService) drawText(img *image.RGBA, text string, pt image.Point, c image.Image, size int, center bool) {
	// 注意：此為簡化版本，實際需要載入字型檔案
	// 文字渲染需要使用 golang.org/x/image/font 或類似套件
	// 這裡僅作預留位置

	_ = text
	_ = pt
	_ = c
	_ = size
	_ = center
}

// parseHexColor 解析 HEX 顏色為 color.Color
func (s *ImageService) parseHexColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.Black
	}

	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// colorToImage 將 color.Color 轉換為 image.Image (用於 draw.Draw)
func (s *ImageService) colorToImage(c color.Color) image.Image {
	return image.NewUniform(c)
}

// overlayBackgroundImage 疊加背景圖
func (s *ImageService) overlayBackgroundImage(img *image.RGBA, backgroundPath string) error {
	var bgImg image.Image
	var err error

	// 支援 URL 和本地路徑
	if strings.HasPrefix(backgroundPath, "http://") || strings.HasPrefix(backgroundPath, "https://") {
		bgImg, err = s.loadImageFromURL(backgroundPath)
	} else {
		bgImg, err = s.loadImageFromFile(backgroundPath)
	}

	if err != nil {
		return err
	}

	// 調整背景圖大小以符合目標圖片
	bgImg = s.resizeImage(bgImg, img.Bounds().Size())

	// 疊加背景圖（使用透明度混合）
	draw.Draw(img, img.Bounds(), bgImg, image.Point{}, draw.Over)

	return nil
}

// loadImageFromURL 從 URL 載入圖片
func (s *ImageService) loadImageFromURL(urlStr string) (image.Image, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
	}

	return s.decodeImage(resp.Body)
}

// loadImageFromFile 從檔案載入圖片
func (s *ImageService) loadImageFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return s.decodeImage(file)
}

// decodeImage 解碼圖片
func (s *ImageService) decodeImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

// resizeImage 調整圖片大小
func (s *ImageService) resizeImage(img image.Image, size image.Point) image.Image {
	// 計算縮放比例
	srcBounds := img.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	scaleX := float64(size.X) / float64(srcWidth)
	scaleY := float64(size.Y) / float64(srcHeight)
	scale := math.Min(scaleX, scaleY)

	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)

	// 建立新圖片
	newImg := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))

	// 計算置中位置
	offsetX := (size.X - newWidth) / 2
	offsetY := (size.Y - newHeight) / 2

	// 繪製縮放後的圖片
	draw.Draw(newImg, image.Rect(offsetX, offsetY, offsetX+newWidth, offsetY+newHeight), img, srcBounds.Min, draw.Over)

	return newImg
}

// imageToJPEG 將圖片轉換為 JPEG 格式
func (s *ImageService) imageToJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
	if err != nil {
		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
	}
	return buf.Bytes(), nil
}

// imageToPNG 將圖片轉換為 PNG 格式
func (s *ImageService) imageToPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

// UploadBackgroundImage 上傳自訂背景圖
func (s *ImageService) UploadBackgroundImage(ctx context.Context, teacherID uint, file io.Reader, filename string) (string, *errInfos.Res, error) {
	// 產生唯一檔案名稱
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	newFilename := fmt.Sprintf("backgrounds/%d/%s%s", teacherID, uuid.New().String(), ext)

	// 儲存到本地或雲端儲存
	// 這裡使用本地儲存作為範例
	uploadDir := filepath.Join("uploads", "backgrounds", fmt.Sprintf("%d", teacherID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", s.app.Err.New(errInfos.SQL_ERROR), err
	}

	filePath := filepath.Join(uploadDir, fmt.Sprintf("%s%s", uuid.New().String(), ext))
	dst, err := os.Create(filePath)
	if err != nil {
		return "", s.app.Err.New(errInfos.SQL_ERROR), err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 回傳相對路徑
	return newFilename, nil, nil
}

// DeleteBackgroundImage 刪除自訂背景圖
func (s *ImageService) DeleteBackgroundImage(ctx context.Context, teacherID uint, imagePath string) error {
	filePath := filepath.Join("uploads", imagePath)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete background image: %w", err)
	}
	return nil
}

// GetTeacherBackgroundImages 取得老師的所有自訂背景圖
func (s *ImageService) GetTeacherBackgroundImages(ctx context.Context, teacherID uint) ([]string, error) {
	uploadDir := filepath.Join("uploads", "backgrounds", fmt.Sprintf("%d", teacherID))

	var images []string
	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel("uploads", path)
			images = append(images, relPath)
		}
		return nil
	})

	return images, nil
}

// ExportTeacherScheduleToImage 匯出老師課表為圖片
func (s *ImageService) ExportTeacherScheduleToImage(ctx context.Context, teacherID, centerID uint, startDate, endDate time.Time, backgroundImage string) ([]byte, *errInfos.Res, error) {
	// 取得中心名稱
	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	centerName := center.Name

	config := &ImageConfig{
		TeacherID:       teacherID,
		CenterID:        centerID,
		StartDate:       startDate,
		EndDate:         endDate,
		BackgroundImage: backgroundImage,
		Title:           fmt.Sprintf("%s - 課表", centerName),
		IncludePersonal: true,
	}

	imgData, err := s.GenerateScheduleImage(ctx, config)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SYSTEM_ERROR), err
	}

	return imgData, nil, nil
}
