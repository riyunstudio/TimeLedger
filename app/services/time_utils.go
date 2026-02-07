package services

import (
	"strconv"
	"strings"
)

// splitTime 分割時間字串 (HH:MM 或 HH:MM:SS) 為 [hour, minute]
func splitTime(timeStr string) []int {
	if timeStr == "" {
		return nil
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return nil
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	return []int{hour, minute}
}

// compareTimeStrings 比較兩個時間字串 HH:MM 格式
// 回傳 -1 表示 t1 < t2, 0 表示相等, 1 表示 t1 > t2
func compareTimeStrings(t1, t2 string) int {
	m1 := timeStringToMinutes(t1)
	m2 := timeStringToMinutes(t2)

	if m1 < m2 {
		return -1
	} else if m1 > m2 {
		return 1
	}
	return 0
}

// timeStringToMinutes 將 HH:MM 格式轉換為總分鐘數
func timeStringToMinutes(s string) int {
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return 0
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	return hour*60 + minute
}
