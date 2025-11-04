package tools

import "fmt"

// Telegram 發送警報訊息
func (tl *Tools) TgSchedulerAlert(topic string, err error, tid string) string {
	return fmt.Sprintf("[Scheduler] %s\n錯誤內容: %v\nTID: %s", topic, err.Error(), tid)
}
