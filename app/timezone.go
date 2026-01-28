package app

import (
	"sync"
	"time"
)

// 中央時區管理
var (
	taipeiLoc *time.Location
	locOnce  sync.Once
	locErr   error
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
