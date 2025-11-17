package logs

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
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
		topic:   "Websocket_Log",
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

const WS_TOPIC_SRV = "Websocket_Server_Log"
const WS_TOPIC_CLI = "Websocket_Client_Log"

const WS_EVENT_SRV_UPGRADE_ERR = "UpgradeError"
const WS_EVENT_SRV_SET_DEADLINE_ERR = "SetReadDeadline"
const WS_EVENT_SRV_READ_ERR = "ReadError"
const WS_EVENT_SRV_BROADCAST_ERR = "BroadcastError"
const WS_EVENT_SRV_CLIENT_CONN = "ClientConnect"
const WS_EVENT_SRV_CLIENT_DIS_CONN = "ClientDisconnect"
const WS_EVENT_SRV_HEART_TIMEOUT = "HeartbeatTimeout"
const WS_EVENT_SRV_SHUTDOWN_ERR = "ServerShutdownError"

const WS_EVENT_CLI_READ_ERR = "ReadError"
const WS_EVENT_CLI_PING_ERR = "PingError"

func (l *WsLog) SetTopic(topic string) *WsLog {
	l.topic = topic
	return l
}

func (l *WsLog) SetEvent(event string) *WsLog {
	l.event = event
	return l
}

func (l *WsLog) SetUuid(uuid string) *WsLog {
	l.uuid = uuid
	return l
}

func (l *WsLog) SetClientIP(ip string) *WsLog {
	l.clientIP = ip
	return l
}

func (l *WsLog) SetUrl(url string) *WsLog {
	l.url = url
	return l
}

func (l *WsLog) SetError(err error) *WsLog {
	if err != nil {
		l.err = err.Error()
	}
	return l
}

func (l *WsLog) SetExtraInfo(info any) *WsLog {
	l.extraInfo = info
	return l
}

// 生成 logrus fields
func (l *WsLog) fields() logrus.Fields {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	f := logrus.Fields{}
	if l.service != "" {
		f["service"] = l.service
	}
	if l.podName != "" {
		f["podName"] = l.podName
	}
	if l.topic != "" {
		f["topic"] = l.topic
	}
	if l.event != "" {
		f["event"] = l.event
	}
	if l.uuid != "" {
		f["uuid"] = l.uuid
	}
	if l.clientIP != "" {
		f["clientIP"] = l.clientIP
	}
	if l.url != "" {
		f["url"] = l.url
	}
	if l.err != "" {
		f["error"] = l.err
	}
	if l.extraInfo != nil {
		f["extraInfo"] = l.extraInfo
	}
	return f
}

// Info log
func (l *WsLog) PrintInfo(msg string) {
	logrus.WithFields(l.fields()).Info(msg)
}

// Error log
func (l *WsLog) PrintError(msg string) {
	logrus.WithFields(l.fields()).Error(msg)
}

func (l *WsLog) PrintPanic(err tools.Panic) {
	l.err = err.StackTrace
	logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(l.fields()).Error(err.Panic)
}
