package logs

import (
	"os"
	"time"
	"timeLedger/global/errInfos"

	"github.com/sirupsen/logrus"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

type GinLog struct {
	// 必要資訊
	service string
	podName string

	// 選填資訊
	topic     string
	url       string
	method    string
	args      any
	headers   any
	domain    string
	clientIP  string
	err       string
	response  any
	extraInfo any
	runTime   float64
	traceID   string
}

func GinLogInit() *GinLog {
	return &GinLog{
		topic:   "Gin_Log",
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

func (l *GinLog) SetUrl(url string) *GinLog {
	l.url = url
	return l
}

func (l *GinLog) SetMethod(method string) *GinLog {
	l.method = method
	return l
}

func (l *GinLog) SetArgs(args any) *GinLog {
	l.args = args
	return l
}

func (l *GinLog) SetHeaders(headers any) *GinLog {
	l.headers = headers
	return l
}

func (l *GinLog) SetDomain(domain string) *GinLog {
	l.domain = domain
	return l
}

func (l *GinLog) SetClientIP(ip string) *GinLog {
	l.clientIP = ip
	return l
}

func (l *GinLog) SetResponse(response any) *GinLog {
	l.response = response
	return l
}

func (l *GinLog) SetExtraInfo(extra any) *GinLog {
	l.extraInfo = extra
	return l
}

func (l *GinLog) SetRunTime(runTime float64) *GinLog {
	l.runTime = runTime
	return l
}

func (l *GinLog) SetTraceID(traceID string) *GinLog {
	l.traceID = traceID
	return l
}

func (l *GinLog) GetTraceID() string {
	return l.traceID
}

func (l *GinLog) fields() logrus.Fields {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	logField := logrus.Fields{}

	if l.service != "" {
		logField["service"] = l.service
	}
	if l.podName != "" {
		logField["podName"] = l.podName
	}
	if l.url != "" {
		logField["url"] = l.url
	}
	if l.method != "" {
		logField["method"] = l.method
	}
	if l.headers != "" {
		logField["headers"] = l.headers
	}
	if l.runTime != 0 {
		logField["runTime"] = l.runTime
	}
	if l.traceID != "" {
		logField["traceID"] = l.traceID
	}
	if l.topic != "" {
		logField["topic"] = l.topic
	}
	if l.args != nil {
		logField["args"] = l.args
	}
	if l.domain != "" {
		logField["domain"] = l.domain
	}
	if l.clientIP != "" {
		logField["clientIP"] = l.clientIP
	}
	if l.err != "" {
		logField["error"] = l.err
	}
	if l.response != nil {
		logField["response"] = l.response
	}
	if l.extraInfo != nil {
		logField["extraInfo"] = l.extraInfo
	}

	return logField
}

// Gin Server 打印專用
func (l *GinLog) PrintServer(code errInfos.ErrCode, res any, err error) {
	l.SetResponse(res)

	switch code {
	case 0:
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(l.fields()).Info()
	default:
		if err != nil {
			l.err = err.Error()
		}
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(l.fields()).Error()
	}
}

func (l *GinLog) PrintPanic(err tools.Panic) {
	l.err = err.StackTrace
	logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(l.fields()).Error(err.Panic)
}
