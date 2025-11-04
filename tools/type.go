package tools

import (
	"strconv"
)

// 型態轉換、檢查工具

// string 轉 int
func (tl *Tools) StrToInt(str string) int {
	num, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}
	return num
}

// string 轉 int32
func (tl *Tools) StrToInt32(str string) int32 {
	num, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}
	return int32(num)
}

// string 轉 int64
func (tl *Tools) StrToInt64(str string) int64 {
	num, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}
	return int64(num)
}

// string 轉 boolean
func (tl *Tools) StrToBool(str string) bool {
	bool, _ := strconv.ParseBool(str)
	return bool
}

// string 轉 float64
func (tl *Tools) StrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0
	}
	return f
}

// int 轉 string
func (tl *Tools) IntToStr(i int) string {
	return strconv.Itoa(i)
}

// int 轉 int32
func (tl *Tools) IntToInt32(i int) int32 {
	return int32(i)
}

// int 轉 int64
func (tl *Tools) IntToInt64(i int) int64 {
	return int64(i)
}

// int 轉 float32
func (tl *Tools) IntTofloat32(i int) float32 {
	return float32(i)
}

// int 轉 float64
func (tl *Tools) IntTofloat64(i int) float64 {
	return float64(i)
}

// int32 轉 string
func (tl *Tools) Int32ToStr(i int32) string {
	return strconv.Itoa(int(i))
}

// int64 轉 string
func (tl *Tools) Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

// int64 轉 int32
func (tl *Tools) Int64ToInt32(i int64) int32 {
	return int32(i)
}

// int32 轉 int64
func (tl *Tools) Int32ToInt64(i int32) int64 {
	return int64(i)
}

// float64 轉 string
func (tl *Tools) Float64ToStr(value float64, decimals int) string {
	return strconv.FormatFloat(value, 'f', decimals, 64)
}

// 檢查 []string 中是否包含特定 string
func (tl *Tools) InStrArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 檢查 []int 中是否包含特定 int
func (tl *Tools) InIntArray(value int, array []int) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 檢查 []int32 中是否包含特定 int32
func (tl *Tools) InInt32Array(value int32, array []int32) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 檢查 []int64 中是否包含特定 int64
func (tl *Tools) InInt64Array(value int64, array []int64) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 陣列刪除指定值
func (tl *Tools) UnsetStrArray(value string, arr []string) []string {
	var result []string

	for _, v := range arr {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
