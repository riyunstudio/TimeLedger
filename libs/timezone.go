package libs

import (
	"sync"
	"time"
)

// 中央時區管理
var (
	taipeiLoc *time.Location
	locOnce   sync.Once
	locErr    error
)

// LoadTaiwanLocation 載入台灣時區（中央化管理，避免重複載入）
func LoadTaiwanLocation() (*time.Location, error) {
	locOnce.Do(func() {
		taipeiLoc, locErr = time.LoadLocation("Asia/Taipei")
	})
	return taipeiLoc, locErr
}

// GetTaiwanLocation 取得台灣時區（可快取）
func GetTaiwanLocation() *time.Location {
	loc, _ := LoadTaiwanLocation()
	return loc
}

// NowInTaiwan 取得台灣時間的現在時間
func NowInTaiwan() time.Time {
	return time.Now().In(GetTaiwanLocation())
}

// TodayInTaiwan 取得台灣時間的今天日期
func TodayInTaiwan() time.Time {
	now := NowInTaiwan()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, GetTaiwanLocation())
}

// TimeToTaiwan 將 time.Time 轉換為台灣時區
// 如果時間本身已是台灣時區則不變，否則轉換
func TimeToTaiwan(t time.Time) time.Time {
	return t.In(GetTaiwanLocation())
}

// TimeToTaiwanOrNow 如果時間為零值，返回台灣時區的現在時間；否則轉換為台灣時區
func TimeToTaiwanOrNow(t time.Time) time.Time {
	if t.IsZero() {
		return NowInTaiwan()
	}
	return TimeToTaiwan(t)
}

// TimeToTaiwanOrNil 將 *time.Time 轉換為台灣時區
// 如果指標為 nil 或時間為零值，返回 nil
func TimeToTaiwanOrNil(t *time.Time) *time.Time {
	if t == nil || t.IsZero() {
		return nil
	}
	result := TimeToTaiwan(*t)
	return &result
}

// TimeToTaiwanOrNowIfNil 如果 *time.Time 為 nil，返回台灣時區的現在時間；否則轉換為台灣時區
func TimeToTaiwanOrNowIfNil(t *time.Time) *time.Time {
	if t == nil {
		now := NowInTaiwan()
		return &now
	}
	result := TimeToTaiwan(*t)
	return &result
}

// ParseDateInTaiwan 解析 YYYY-MM-DD 日期字串為台灣時區的時間
func ParseDateInTaiwan(dateStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", dateStr, GetTaiwanLocation())
}

// ParseDateTimeInTaiwan 解析日期時間字串為台灣時區的時間
func ParseDateTimeInTaiwan(dateTimeStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", dateTimeStr, GetTaiwanLocation())
}
