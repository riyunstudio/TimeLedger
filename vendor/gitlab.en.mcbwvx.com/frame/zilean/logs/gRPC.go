package logs

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

type GrpcLog struct {
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

func GrpcLogInit() *GrpcLog {
	return &GrpcLog{
		topic:   "gRPC_Log",
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

func (l *GrpcLog) SetUrl(url string) *GrpcLog {
	l.url = url
	return l
}

func (l *GrpcLog) SetMethod(method string) *GrpcLog {
	l.method = method
	return l
}

func (l *GrpcLog) SetArgs(args any) *GrpcLog {
	l.args = args
	return l
}

func (l *GrpcLog) SetHeaders(headers any) *GrpcLog {
	l.headers = headers
	return l
}

func (l *GrpcLog) SetDomain(domain string) *GrpcLog {
	l.domain = domain
	return l
}

func (l *GrpcLog) SetClientIP(ip string) *GrpcLog {
	l.clientIP = ip
	return l
}

func (l *GrpcLog) SetResponse(response any) *GrpcLog {
	l.response = response
	return l
}

func (l *GrpcLog) SetExtraInfo(extra any) *GrpcLog {
	l.extraInfo = extra
	return l
}

func (l *GrpcLog) SetRunTime(runTime float64) *GrpcLog {
	l.runTime = runTime
	return l
}

func (l *GrpcLog) SetTraceID(traceID string) *GrpcLog {
	l.traceID = traceID
	return l
}

func (l *GrpcLog) GetTraceID() string {
	return l.traceID
}

func (l *GrpcLog) fields() logrus.Fields {
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

func (l *GrpcLog) PrintGrpc(res any, err error) {
	l.SetResponse(res)

	switch err {
	case nil:
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(l.fields()).Info()
	default:
		l.err = err.Error()
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(l.fields()).Error()
	}
}

func (l *GrpcLog) PrintPanic(err tools.Panic) {
	l.err = err.StackTrace
	logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(l.fields()).Error(err.Panic)
}
