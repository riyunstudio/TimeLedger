package logs

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

type CurlLog struct {
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

func CurlLogInit() *CurlLog {
	return &CurlLog{
		topic:   "Curl_Log",
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

func (l *CurlLog) SetUrl(url string) *CurlLog {
	l.url = url
	return l
}

func (l *CurlLog) SetMethod(method string) *CurlLog {
	l.method = method
	return l
}

func (l *CurlLog) SetArgs(args any) *CurlLog {
	l.args = args
	return l
}

func (l *CurlLog) SetHeaders(headers any) *CurlLog {
	l.headers = headers
	return l
}

func (l *CurlLog) SetDomain(domain string) *CurlLog {
	l.domain = domain
	return l
}

func (l *CurlLog) SetClientIP(ip string) *CurlLog {
	l.clientIP = ip
	return l
}

func (l *CurlLog) SetResponse(response any) *CurlLog {
	l.response = response
	return l
}

func (l *CurlLog) SetExtraInfo(extra any) *CurlLog {
	l.extraInfo = extra
	return l
}

func (l *CurlLog) SetRunTime(runTime float64) *CurlLog {
	l.runTime = runTime
	return l
}

func (l *CurlLog) SetTraceID(traceID string) *CurlLog {
	l.traceID = traceID
	return l
}

func (l *CurlLog) GetTraceID() string {
	return l.traceID
}

func (l *CurlLog) fields() logrus.Fields {
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

func (l *CurlLog) PrintCurl(msg string, err error) {
	if err != nil {
		l.err = err.Error()
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(l.fields()).Error(msg)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(l.fields()).Info(msg)
	}
}

func (l *CurlLog) PrintPanic(err tools.Panic) {
	l.err = err.StackTrace
	logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(l.fields()).Error(err.Panic)
}
