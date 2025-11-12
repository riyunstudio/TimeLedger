package ws

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type WsLog struct {
	service   string
	podName   string
	event     string // 事件名稱
	clientIP  string // 客戶端 IP
	topic     string // 主題或操作名稱
	err       string // 錯誤訊息
	extraInfo any    // 額外資訊
}

func WsLogInit() *WsLog {
	return &WsLog{
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

func (wl *WsLog) SetEvent(event string) *WsLog {
	wl.event = event
	return wl
}

func (wl *WsLog) SetClientIP(ip string) *WsLog {
	wl.clientIP = ip
	return wl
}

func (wl *WsLog) SetTopic(topic string) *WsLog {
	wl.topic = topic
	return wl
}

func (wl *WsLog) SetError(err error) *WsLog {
	if err != nil {
		wl.err = err.Error()
	}
	return wl
}

func (wl *WsLog) SetExtraInfo(info interface{}) *WsLog {
	wl.extraInfo = info
	return wl
}

// 生成 logrus fields
func (wl *WsLog) fields() logrus.Fields {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	f := logrus.Fields{}
	if wl.service != "" {
		f["service"] = wl.service
	}
	if wl.podName != "" {
		f["podName"] = wl.podName
	}
	if wl.event != "" {
		f["event"] = wl.event
	}
	if wl.clientIP != "" {
		f["clientIP"] = wl.clientIP
	}
	if wl.topic != "" {
		f["topic"] = wl.topic
	}
	if wl.err != "" {
		f["error"] = wl.err
	}
	if wl.extraInfo != nil {
		f["extraInfo"] = wl.extraInfo
	}
	return f
}

// Info log
func (wl *WsLog) PrintInfo(msg string) {
	logrus.WithFields(wl.fields()).Info(msg)
}

// Error log
func (wl *WsLog) PrintError(msg string) {
	logrus.WithFields(wl.fields()).Error(msg)
}
