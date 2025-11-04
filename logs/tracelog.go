package logs

import (
	"akali/global"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type TraceLog struct {
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

func TraceLogInit() *TraceLog {
	return &TraceLog{
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

func (tl *TraceLog) SetTopic(topic string) {
	tl.topic = topic
}

func (tl *TraceLog) SetUrl(url string) {
	tl.url = url
}

func (tl *TraceLog) SetMethod(method string) {
	tl.method = method
}

func (tl *TraceLog) SetArgs(args any) {
	tl.args = args
}

func (tl *TraceLog) SetHeaders(headers any) {
	tl.headers = headers
}

func (tl *TraceLog) SetDomain(domain string) {
	tl.domain = domain
}

func (tl *TraceLog) SetClientIP(ip string) {
	tl.clientIP = ip
}

func (tl *TraceLog) SetResponse(response any) {
	tl.response = response
}

func (tl *TraceLog) SetExtraInfo(extra any) {
	tl.extraInfo = extra
}

func (tl *TraceLog) SetRunTime(runTime float64) {
	tl.runTime = runTime
}

func (tl *TraceLog) SetTraceID(traceID string) {
	tl.traceID = traceID
}

func (tl *TraceLog) GetTraceID() string {
	return tl.traceID
}

func getlogfields(tl *TraceLog) logrus.Fields {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	logField := logrus.Fields{}

	// 固定資訊
	logField["service"] = tl.service
	logField["podName"] = tl.podName
	logField["url"] = tl.url
	logField["method"] = tl.method
	logField["headers"] = tl.headers
	logField["runTime"] = tl.runTime
	logField["traceID"] = tl.traceID

	// 非固定資訊
	if tl.topic != "" {
		logField["topic"] = tl.topic
	}
	if tl.args != nil {
		logField["args"] = tl.args
	}
	if tl.domain != "" {
		logField["domain"] = tl.domain
	}
	if tl.clientIP != "" {
		logField["clientIP"] = tl.clientIP
	}
	if tl.err != "" {
		logField["error"] = tl.err
	}
	if tl.response != nil {
		logField["response"] = tl.response
	}
	if tl.extraInfo != nil {
		logField["extraInfo"] = tl.extraInfo
	}

	return logField
}

// Panic 打印專用
func (tl *TraceLog) PrintPanic(err global.Panic) {
	tl.err = err.StackTrace
	lf := getlogfields(tl)
	logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(lf).Error(err.Panic)
}

// Gin Server 打印專用
func (tl *TraceLog) PrintServer(res global.Ret) {
	tl.SetResponse(res)

	switch res.ErrInfo.Code {
	case 0:
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(lf).Info()
	default:
		tl.err = res.Err.Error()
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(lf).Error()
	}
}

// Curl 打印專用
func (tl *TraceLog) PrintCurl(msg string, err error) {
	if err != nil {
		tl.err = err.Error()
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(lf).Error(msg)
	} else {
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(lf).Info(msg)
	}
}

func (tl *TraceLog) PrintSchedule(err error) {
	if err != nil {
		tl.err = err.Error()
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(lf).Error()
	} else {
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(lf).Info()
	}
}

// gRPC Server 打印專用
func (tl *TraceLog) PrintGrpc(res any, err error) {
	tl.SetResponse(res)

	switch err {
	case nil:
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(lf).Info()
	default:
		tl.err = err.Error()
		lf := getlogfields(tl)
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(lf).Error()
	}
}
