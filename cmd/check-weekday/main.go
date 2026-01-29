package main

import (
	"fmt"
	"time"
)

func main() {
	// 檢查 2026年1月的日期對應
	dates := []time.Time{
		time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC), // 週一
		time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC), // 週二
		time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC), // 週三
	}

	weekdayNames := []string{"週日", "週一", "週二", "週三", "週四", "週五", "週六"}

	for _, d := range dates {
		weekday := int(d.Weekday())
		if weekday == 0 {
			weekday = 7 // 將週日改為 7
		}
		fmt.Printf("%s: Go Weekday=%d, Our System=%d (%s)\n",
			d.Format("2006-01-02"),
			int(d.Weekday()),
			weekday,
			weekdayNames[d.Weekday()])
	}
}
