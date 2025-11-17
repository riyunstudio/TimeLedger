package tools

import (
	"time"
)

// 取得當前時區
func (tl *Tools) Locat() *time.Location {
	return tl.loc
}

// 取得當前時間字串 日期時間 YYYY-mm-dd HH:ii:ss
func (tl *Tools) NowDateTime() string {
	return time.Now().In(tl.loc).Format("2006-01-02 15:04:05")
}

// 取得當前時間字串 日期 YYYY-mm-dd
func (tl *Tools) NowDate() string {
	return time.Now().In(tl.loc).Format("2006-01-02")
}

// 取得當前時間字串 時間 HH:ii:ss
func (tl *Tools) NowTime() string {
	return time.Now().In(tl.loc).Format("15:04:05")
}

// 取得當前時間字串 年 YYYY
func (tl *Tools) NowYear() string {
	return time.Now().In(tl.loc).Format("2006")
}

// 取得當前時間字串 月 mm
func (tl *Tools) NowMonth() string {
	return time.Now().In(tl.loc).Format("01")
}

// 取得當前時間字串 日 dd
func (tl *Tools) NowDay() string {
	return time.Now().In(tl.loc).Format("02")
}

// 取得當前時間字串 時 HH
func (tl *Tools) NowHour() string {
	return time.Now().In(tl.loc).Format("15")
}

// 取得當前時間字串 分 ii
func (tl *Tools) NowMin() string {
	return time.Now().In(tl.loc).Format("04")
}

// 取得當前時間字串 秒 ss
func (tl *Tools) NowSec() string {
	return time.Now().In(tl.loc).Format("05")
}

// 取得當前時間字串 自訂格式
func (tl *Tools) CustomNowFormat(format string) string {
	return time.Now().In(tl.loc).Format(format)
}

// 取得當前時間戳
func (tl *Tools) NowUnix() int64 {
	return time.Now().In(tl.loc).Unix()
}

// 取得當前時間
func (tl *Tools) Now() time.Time {
	return time.Now().In(tl.loc)
}

func (tl *Tools) NowUnixTrimmedToMinute(sec int) int64 {
	now := time.Now().In(tl.loc)
	trimmed := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		sec, // 秒設為0
		0,   // 奈秒設為0
		tl.loc,
	)
	return trimmed.Unix()
}

func (tl *Tools) TodayTime(dayOffset ...int) time.Time {
	offset := 0
	if len(dayOffset) > 0 {
		offset = dayOffset[0]
	}

	now := time.Now().In(tl.loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tl.loc)

	// 加上天數偏移
	return today.AddDate(0, 0, offset)
}

// 取得指定時間戳轉當日整天區間時間字串 日期時間 YYYY-mm-dd HH:ii:ss
func (tl *Tools) RangeDateTime(unix int64) []string {
	day := time.Unix(unix, 0).In(tl.loc)

	return []string{
		time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, tl.loc).Format("2006-01-02 15:04:05"),
		time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, tl.loc).Format("2006-01-02 15:04:05"),
	}
}

// 取得指定時間戳轉當日整天區間時間戳
func (tl *Tools) RangeUnix(unix int64) []int64 {
	day := time.Unix(unix, 0).In(tl.loc)

	return []int64{
		time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, tl.loc).Unix(),
		time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, tl.loc).Unix(),
	}
}

// 取得指定時間字串轉時間戳
func (tl *Tools) DateTimeToUnix(DateTime string) int64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", DateTime, tl.loc)

	if err != nil {
		return 0
	}
	return t.Unix()
}

// 取得指定時間字串轉時間戳
func (tl *Tools) DateToUnix(Date string) int64 {
	t, err := time.ParseInLocation("2006-01-02", Date, tl.loc)

	if err != nil {
		return 0
	}
	return t.Unix()
}

// 取得指定時間字串轉時間戳
func (tl *Tools) TimeToUnix(Time string) int64 {
	t, err := time.ParseInLocation("15:04:05", Time, tl.loc)

	if err != nil {
		return 0
	}
	return t.Unix()
}

// 指定時間戳轉時間字串 日期時間 YYYY-mm-dd HH:ii:ss
func (tl *Tools) UnixToDateTime(unix int64) string {
	return time.Unix(unix, 0).In(tl.loc).Format("2006-01-02 15:04:05")
}

// 指定時間戳轉時間字串 日期 YYYY-mm-dd
func (tl *Tools) UnixToDate(unix int64) string {
	return time.Unix(unix, 0).In(tl.loc).Format("2006-01-02")
}

// 指定時間戳轉時間字串 時間 HH:ii:ss
func (tl *Tools) UnixToTime(unix int64) string {
	return time.Unix(unix, 0).In(tl.loc).Format("15:04:05")
}

// 指定時間戳轉時間字串 自訂格式
func (tl *Tools) CustUnixToDate(format string, unix int64) string {
	return time.Unix(unix, 0).In(tl.loc).Format(format)
}

// 指定時間戳增減時間, d 給正的是增加, 負的是減少
func (tl *Tools) UnixAdd(unix int64, d time.Duration) int64 {
	return time.Unix(unix, 0).In(tl.loc).Add(d).Unix()
}

// 時間戳(毫秒)轉時間戳
func (tl *Tools) UnixMilliToUnix(unix int64) int64 {
	return time.UnixMilli(unix).In(tl.loc).Unix()
}

// 解析字串時區 + | - 符號, 回傳對應的時間
func (tl *Tools) ParserPrefixToTime(prefex string) int {
	switch prefex[:1] {
	case "-":
		return tl.StrToInt(prefex[1:]) * -1
	default:
		return tl.StrToInt(prefex[1:])
	}
}

// ? --------------- 以下為 指定時區調用方法 -----------------

// 取得指定時區
func (tl *Tools) LoadLocat(loc string) (*time.Location, error) {
	return time.LoadLocation(loc)
}

// 取得當前時間戳 自訂時區
func (tl *Tools) NowUnixByLocat(loc *time.Location) int64 {
	return time.Now().In(loc).Unix()
}

// 取得當前時間戳 自訂時區 (毫秒)
func (tl *Tools) NowUnixMilliByLocat(loc *time.Location) int64 {
	return time.Now().In(loc).UnixMilli()
}

// 時間戳轉時間戳(毫秒) 自訂時區
func (tl *Tools) UnixToUnixMilliByLocat(loc *time.Location, unix int64) int64 {
	return time.Unix(unix, 0).In(loc).UnixMilli()
}

// 時間戳轉時間字串, 自訂格式、時區
// noSec: true時會回傳00的秒數
func (tl *Tools) CustUnixToDateByLocat(loc *time.Location, format string, unix int64, noSec bool) string {
	if noSec {
		return time.Unix(unix, 0).Truncate(time.Minute).In(loc).Format(format)
	}
	return time.Unix(unix, 0).In(loc).Format(format)
}

// 取得指定時間字串轉時間戳
func (tl *Tools) TimeToUnixByLocat(loc *time.Location, dateTime string) int64 {
	// 取得當前時間，設置今天的日期
	now := time.Now().In(loc)

	// 解析時間字串，格式為 "15:04:05"，並使用傳入的時區
	parsedTime, err := time.ParseInLocation("15:04:05", dateTime, loc)
	if err != nil {
		return 0
	}

	// 將解析後的時間與今天的日期結合
	finalTime := time.Date(now.Year(), now.Month(), now.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, loc)

	return finalTime.Unix()
}
