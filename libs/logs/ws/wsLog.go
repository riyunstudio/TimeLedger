package ws

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type WsLog struct {
	service   string
	podName   string
	topic     string // 主題或操作名稱
	event     string // 事件名稱
	uuid      string // 客戶 ID
	clientIP  string // 客戶端 IP
	url       string // 路由
	err       string // 錯誤訊息
	extraInfo any    // 額外資訊
}

func WsLogInit() *WsLog {
	return &WsLog{
		topic:   "Websocket",
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

const TOPIC_SRV = "Websocket_Server"
const TOPIC_CLI = "Websocket_Client"

const EVENT_SRV_UPGRADE_ERR = "UpgradeError"
const EVENT_SRV_SET_DEADLINE_ERR = "SetReadDeadline"
const EVENT_SRV_READ_ERR = "ReadError"
const EVENT_SRV_BROADCAST_ERR = "BroadcastError"
const EVENT_SRV_CLIENT_CONN = "ClientConnect"
const EVENT_SRV_CLIENT_DIS_CONN = "ClientDisconnect"
const EVENT_SRV_HEART_TIMEOUT = "HeartbeatTimeout"
const EVENT_SRV_SHUTDOWN_ERR = "ServerShutdownError"

const EVENT_CLI_READ_ERR = "ReadError"
const EVENT_CLI_PING_ERR = "PingError"

func (wl *WsLog) SetTopic(topic string) *WsLog {
	wl.topic = topic
	return wl
}

func (wl *WsLog) SetEvent(event string) *WsLog {
	wl.event = event
	return wl
}

func (wl *WsLog) SetUuid(uuid string) *WsLog {
	wl.uuid = uuid
	return wl
}

func (wl *WsLog) SetClientIP(ip string) *WsLog {
	wl.clientIP = ip
	return wl
}

func (wl *WsLog) SetUrl(url string) *WsLog {
	wl.url = url
	return wl
}

func (wl *WsLog) SetError(err error) *WsLog {
	if err != nil {
		wl.err = err.Error()
	}
	return wl
}

func (wl *WsLog) SetExtraInfo(info any) *WsLog {
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
	if wl.topic != "" {
		f["topic"] = wl.topic
	}
	if wl.event != "" {
		f["event"] = wl.event
	}
	if wl.uuid != "" {
		f["uuid"] = wl.uuid
	}
	if wl.clientIP != "" {
		f["clientIP"] = wl.clientIP
	}
	if wl.url != "" {
		f["url"] = wl.url
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
