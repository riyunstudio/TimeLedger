package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/global/errInfos"

	"github.com/google/uuid"
)

// ICSCalendarService 處理 iCalendar 格式的課表匯出
type ICSCalendarService struct {
	app *app.App
}

// NewICSCalendarService 建立 ICSCalendarService 實例
func NewICSCalendarService(app *app.App) *ICSCalendarService {
	return &ICSCalendarService{
		app: app,
	}
}

// ScheduleEvent 代表單一課表事件
type ScheduleEvent struct {
	ID          string
	Summary     string    // 課程名稱
	Description string    // 課程描述
	Location    string    // 教室地點
	StartTime   time.Time // 開始時間
	EndTime     time.Time // 結束時間
	TeacherName string    // 老師名稱
	CenterName  string    // 中心名稱
	AllDay      bool      // 是否為全天事件
}

// ICSConfig ICS 匯出配置
type ICSConfig struct {
	TeacherID   uint      // 老師 ID
	CenterID    uint      // 中心 ID
	StartDate   time.Time // 開始日期
	EndDate     time.Time // 結束日期
	CenterName  string    // 中心名稱
	TeacherName string    // 老師名稱
	Events      []ScheduleEvent
}

// GenerateICS 產生 iCalendar 格式的課表資料
func (s *ICSCalendarService) GenerateICS(ctx context.Context, config *ICSConfig) ([]byte, error) {
	var buf bytes.Buffer

	// 寫入 VCALENDAR 開頭
	calUUID := uuid.New().String()
	buf.WriteString("BEGIN:VCALENDAR\r\n")
	buf.WriteString("VERSION:2.0\r\n")
	buf.WriteString("PRODID:-//TimeLedger//Schedule//EN\r\n")
	buf.WriteString(fmt.Sprintf("UID:%s\r\n", calUUID))
	buf.WriteString("CALSCALE:GREGORIAN\r\n")
	buf.WriteString("METHOD:PUBLISH\r\n")
	buf.WriteString(fmt.Sprintf("X-WR-CALNAME:%s 課表\r\n", config.CenterName))

	// 寫入每個事件
	for _, event := range config.Events {
		eventICS, err := s.eventToICS(&event)
		if err != nil {
			continue
		}
		buf.WriteString(eventICS)
	}

	// 寫入 VCALENDAR 結尾
	buf.WriteString("END:VCALENDAR\r\n")

	return buf.Bytes(), nil
}

// eventToICS 將 ScheduleEvent 轉換為 VEVENT 區塊
func (s *ICSCalendarService) eventToICS(event *ScheduleEvent) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("BEGIN:VEVENT\r\n")

	// 事件 UID
	if event.ID != "" {
		buf.WriteString(fmt.Sprintf("UID:%s\r\n", event.ID))
	} else {
		buf.WriteString(fmt.Sprintf("UID:%s\r\n", uuid.New().String()))
	}

	// 時間戳記
	buf.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", time.Now().UTC().Format("20060102T150405Z")))

	// 開始時間
	if event.AllDay {
		// 全天事件格式：YYYYMMDD
		startDate := event.StartTime.Format("20060102")
		buf.WriteString(fmt.Sprintf("DTSTART;VALUE=DATE:%s\r\n", startDate))
	} else {
		// 一般事件格式：YYYYMMDDTHHMMSSZ
		startTime := event.StartTime.UTC().Format("20060102T150405Z")
		endTime := event.EndTime.UTC().Format("20060102T150405Z")
		buf.WriteString(fmt.Sprintf("DTSTART:%s\r\n", startTime))
		buf.WriteString(fmt.Sprintf("DTEND:%s\r\n", endTime))
	}

	// 標題
	summary := event.Summary
	if event.TeacherName != "" {
		summary = fmt.Sprintf("%s (%s)", event.Summary, event.TeacherName)
	}
	buf.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", s.escapeICSText(summary)))

	// 描述
	if event.Description != "" {
		buf.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", s.escapeICSText(event.Description)))
	}

	// 地點
	if event.Location != "" {
		buf.WriteString(fmt.Sprintf("LOCATION:%s\r\n", s.escapeICSText(event.Location)))
	}

	// 狀態
	buf.WriteString("STATUS:CONFIRMED\r\n")

	// 寫入 VEVENT 結尾
	buf.WriteString("END:VEVENT\r\n")

	return buf.String(), nil
}

// escapeICSText 逸出 ICS 文字中的特殊字元
func (s *ICSCalendarService) escapeICSText(text string) string {
	// 需要逸出的字元：\ , ; \n
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, ";", "\\;")
	text = strings.ReplaceAll(text, ",", "\\,")
	// 將換行替換為 \n
	text = strings.ReplaceAll(text, "\n", "\\n")
	return text
}

// ParseICS 解析 ICS 檔案
func (s *ICSCalendarService) ParseICS(reader io.Reader) ([]ScheduleEvent, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read ICS content: %w", err)
	}

	// 簡單解析 - 提取 VEVENT 區塊
	events := make([]ScheduleEvent, 0)
	contentStr := string(content)

	// 檢查是否為有效的 ICS 格式
	if !strings.Contains(contentStr, "BEGIN:VCALENDAR") {
		return nil, fmt.Errorf("invalid ICS format: missing VCALENDAR")
	}

	// 這裡可以實作更完整的 ICS 解析邏輯
	// 對於基本需求，我們先返回空陣列

	return events, nil
}

// GenerateSubscriptionToken 產生訂閱 token
func (s *ICSCalendarService) GenerateSubscriptionToken(teacherID uint, centerID uint) (string, error) {
	token := uuid.New().String()
	return token, nil
}

// GenerateSubscriptionURL 產生課表訂閱連結
func (s *ICSCalendarService) GenerateSubscriptionURL(token string) string {
	baseURL := ""
	if s.app.Env != nil {
		baseURL = s.app.Env.FrontendBaseURL
	}
	if baseURL == "" {
		baseURL = "https://timeledger.app"
	}
	return fmt.Sprintf("%s/api/v1/calendar/subscribe/%s.ics", baseURL, token)
}

// ValidateSubscriptionToken 驗證訂閱 token 是否有效
func (s *ICSCalendarService) ValidateSubscriptionToken(ctx context.Context, token string) (uint, uint, bool, error) {
	// 簡單實作：檢查 token 格式
	if len(token) < 10 {
		return 0, 0, false, fmt.Errorf("invalid token format")
	}

	return 0, 0, false, nil
}

// GetTeacherScheduleForICS 取得老師的課表資料用於 ICS 匯出
func (s *ICSCalendarService) GetTeacherScheduleForICS(ctx context.Context, teacherID, centerID uint, startDate, endDate time.Time) ([]ScheduleEvent, error) {
	scheduleQueryService := NewScheduleQueryService(s.app)
	schedule, err := scheduleQueryService.GetTeacherSchedule(ctx, teacherID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher schedule: %w", err)
	}

	events := make([]ScheduleEvent, 0)
	for _, item := range schedule {
		if item.Type == "CENTER_SESSION" {
			// 解析日期和時間
			dateStr := item.Date
			startTimeStr := item.StartTime
			endTimeStr := item.EndTime

			var startAt, endAt time.Time
			if dateStr != "" && startTimeStr != "" {
				startAt, _ = time.Parse("2006-01-02 15:04", dateStr+" "+startTimeStr)
			}
			if dateStr != "" && endTimeStr != "" {
				endAt, _ = time.Parse("2006-01-02 15:04", dateStr+" "+endTimeStr)
			}

			event := ScheduleEvent{
				ID:          item.ID,
				Summary:     item.Title,
				Description: fmt.Sprintf("課程：%s\\n時間：%s - %s", item.Title, item.StartTime, item.EndTime),
				Location:    "",
				StartTime:   startAt,
				EndTime:     endAt,
				TeacherName: "",
				CenterName:  item.CenterName,
			}
			events = append(events, event)
		}
	}

	return events, nil
}

// ExportTeacherScheduleToICS 匯出老師課表為 ICS 格式
func (s *ICSCalendarService) ExportTeacherScheduleToICS(ctx context.Context, teacherID, centerID uint, startDate, endDate time.Time) ([]byte, error) {
	// 取得課表資料
	events, err := s.GetTeacherScheduleForICS(ctx, teacherID, centerID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 產生 ICS
	config := &ICSConfig{
		TeacherID:  teacherID,
		CenterID:   centerID,
		StartDate:  startDate,
		EndDate:    endDate,
		CenterName: "",
		Events:     events,
	}

	return s.GenerateICS(ctx, config)
}

// SubscriptionInfo 訂閱資訊
type SubscriptionInfo struct {
	Token       string    `json:"token"`
	URL         string    `json:"url"`
	TeacherID   uint      `json:"teacher_id"`
	CenterID    uint      `json:"center_id"`
	CenterName  string    `json:"center_name"`
	TeacherName string    `json:"teacher_name"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateCalendarSubscription 建立課表訂閱
func (s *ICSCalendarService) CreateCalendarSubscription(ctx context.Context, teacherID, centerID uint) (*SubscriptionInfo, *errInfos.Res, error) {
	// 產生訂閱 token
	token, err := s.GenerateSubscriptionToken(teacherID, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SYSTEM_ERROR), err
	}

	// 產生訂閱連結
	subscriptionURL := s.GenerateSubscriptionURL(token)

	now := time.Now()
	expiresAt := now.AddDate(1, 0, 0) // 一年後過期

	info := &SubscriptionInfo{
		Token:       token,
		URL:         subscriptionURL,
		TeacherID:   teacherID,
		CenterID:    centerID,
		CenterName:  "",
		TeacherName: "",
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	}

	return info, nil, nil
}

// Unsubscribe 取消訂閱
func (s *ICSCalendarService) Unsubscribe(ctx context.Context, token string) error {
	return nil
}

// ParseSubscriptionToken 解析訂閱 URL 中的 token
func (s *ICSCalendarService) ParseSubscriptionToken(subscriptionURL string) (string, error) {
	u, err := url.Parse(subscriptionURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// 提取 path 中的 token
	path := u.Path
	path = strings.TrimPrefix(path, "/api/v1/calendar/subscribe/")
	path = strings.TrimSuffix(path, ".ics")

	if path == "" {
		return "", fmt.Errorf("invalid subscription URL format")
	}

	return path, nil
}
