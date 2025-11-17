package logs

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Log struct {
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

func PanicLogInit() *Log {
	return &Log{
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

func (l *Log) SetTopic(topic string) *Log {
	l.topic = topic
	return l
}

func (l *Log) SetUrl(url string) *Log {
	l.url = url
	return l
}

func (l *Log) SetMethod(method string) *Log {
	l.method = method
	return l
}

func (l *Log) SetArgs(args any) *Log {
	l.args = args
	return l
}

func (l *Log) SetHeaders(headers any) *Log {
	l.headers = headers
	return l
}

func (l *Log) SetDomain(domain string) *Log {
	l.domain = domain
	return l
}

func (l *Log) SetClientIP(ip string) *Log {
	l.clientIP = ip
	return l
}

func (l *Log) SetResponse(response any) *Log {
	l.response = response
	return l
}

func (l *Log) SetExtraInfo(extra any) *Log {
	l.extraInfo = extra
	return l
}

func (l *Log) SetRunTime(runTime float64) *Log {
	l.runTime = runTime
	return l
}

func (l *Log) SetTraceID(traceID string) *Log {
	l.traceID = traceID
	return l
}

func (l *Log) GetTraceID() string {
	return l.traceID
}

func (l *Log) fields() logrus.Fields {
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

func (l *Log) PrintErr(err error) {
	l.err = err.Error()
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.WithFields(l.fields()).Error()
}
