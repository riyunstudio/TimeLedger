package services

import (
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
)

// IsCrossDayTime 檢查是否為跨日時間（結束時間早於開始時間）
func IsCrossDayTime(startTime, endTime string) bool {
	return endTime < startTime
}

// ParseTimeToMinutes 將 HH:MM 格式轉換為當天分鐘數
func ParseTimeToMinutes(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return 0
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	return hour*60 + minute
}

// TimesOverlapCrossDay 處理跨日時間重疊檢測
// 對於跨日課程（如 23:00-02:00），需要檢查：
// 1. 當天 23:00 以後是否有其他課程
// 2. 隔天 02:00 以前是否有其他課程
func TimesOverlapCrossDay(start1, end1 string, isCrossDay1 bool, start2, end2 string, isCrossDay2 bool) bool {
	// 如果兩個都是普通課程（非跨日），使用簡單的時間比較
	if !isCrossDay1 && !isCrossDay2 {
		return start1 < end2 && end1 > start2
	}

	// 將時間轉換為分鐘數
	start1Min := ParseTimeToMinutes(start1)
	end1Min := ParseTimeToMinutes(end1)
	start2Min := ParseTimeToMinutes(start2)
	end2Min := ParseTimeToMinutes(end2)

	// 處理跨日課程（結束時間加 24*60 分鐘）
	if isCrossDay1 {
		end1Min += 24 * 60
	}
	if isCrossDay2 {
		end2Min += 24 * 60
	}

	// 檢查重疊
	return start1Min < end2Min && end1Min > start2Min
}

// GetNextWeekday 取得下一天的星期幾（1-7，週一到週日）
func GetNextWeekday(currentWeekday int) int {
	next := currentWeekday + 1
	if next > 7 {
		next = 1
	}
	return next
}

// GetPreviousWeekday 取得上一天的星期幾（1-7，週一到週日）
func GetPreviousWeekday(currentWeekday int) int {
	prev := currentWeekday - 1
	if prev < 1 {
		prev = 7
	}
	return prev
}

// FormatTimeForStorage 格式化時間用於存儲（HH:MM 格式）
func FormatTimeForStorage(t time.Time) string {
	return t.Format("15:04")
}

// ParseTimeFromStorage 從存儲格式解析時間
func ParseTimeFromStorage(timeStr string) (hour, minute int) {
	parts := strings.Split(timeStr, ":")
	if len(parts) >= 2 {
		hour, _ = strconv.Atoi(parts[0])
		minute, _ = strconv.Atoi(parts[1])
	}
	return hour, minute
}

// GetWeekdayFromTime 從時間取得星期幾（1-7，週一到週日）
func GetWeekdayFromTime(t time.Time) int {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // 週日轉換為 7
	}
	return weekday
}

// CreateCrossDayDateRange 為跨日課程建立有效的日期範圍
// 跨日課程的有效範圍需要考慮到課程會延續到隔天
func CreateCrossDayDateRange(startDate time.Time, weeks int) (start, end time.Time) {
	loc := app.GetTaiwanLocation()
	start = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	end = start.AddDate(0, 0, weeks*7-1)
	return start, end
}

// IsTimeInCrossDayRange 檢查時間是否在跨日課程的範圍內
// 跨日課程的時間範圍是：當天 20:00 到隔天 08:00
func IsTimeInCrossDayRange(timeStr string) bool {
	hour, _ := ParseTimeFromStorage(timeStr)
	// 跨日課程通常在晚上開始（20:00-23:59）或凌晨（00:00-08:00）
	return hour >= 20 || hour < 8
}

// GetCrossDayOverlapInfo 取得跨日課程的重疊詳細資訊
func GetCrossDayOverlapInfo(start1, end1 string, isCrossDay1 bool, start2, end2 string, isCrossDay2 bool) map[string]interface{} {
	result := make(map[string]interface{})

	if TimesOverlapCrossDay(start1, end1, isCrossDay1, start2, end2, isCrossDay2) {
		result["overlaps"] = true

		// 計算重疊的分鐘數
		start1Min := ParseTimeToMinutes(start1)
		end1Min := ParseTimeToMinutes(end1)
		start2Min := ParseTimeToMinutes(start2)
		end2Min := ParseTimeToMinutes(end2)

		if isCrossDay1 {
			end1Min += 24 * 60
		}
		if isCrossDay2 {
			end2Min += 24 * 60
		}

		overlapStart := max(start1Min, start2Min)
		overlapEnd := min(end1Min, end2Min)
		overlapMinutes := overlapEnd - overlapStart

		result["overlap_minutes"] = overlapMinutes

		// 判斷重疊發生的時段
		if overlapStart < 24*60 {
			result["overlap_day"] = "current"
			result["overlap_time"] = FormatMinutesToTime(overlapStart)
		} else {
			result["overlap_day"] = "next"
			result["overlap_time"] = FormatMinutesToTime(overlapStart - 24*60)
		}
	} else {
		result["overlaps"] = false
	}

	return result
}

// FormatMinutesToTime 將分鐘數轉換為 HH:MM 格式
func FormatMinutesToTime(minutes int) string {
	hour := minutes / 60
	min := minutes % 60
	return strconv.Itoa(hour) + ":" + strconv.Itoa(min)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
