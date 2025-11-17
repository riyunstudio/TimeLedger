package logs

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
)

type ScheduleLog struct {
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

func ScheduleLogInit() *ScheduleLog {
	return &ScheduleLog{
		topic:   "Schedule_Log",
		service: os.Getenv("SERVICE_NAME"),
		podName: os.Getenv("HOSTNAME"),
	}
}

func (l *ScheduleLog) SetTopic(topic string) *ScheduleLog {
	l.topic = topic
	return l
}

func (l *ScheduleLog) SetUrl(url string) *ScheduleLog {
	l.url = url
	return l
}

func (l *ScheduleLog) SetMethod(method string) *ScheduleLog {
	l.method = method
	return l
}

func (l *ScheduleLog) SetArgs(args any) *ScheduleLog {
	l.args = args
	return l
}

func (l *ScheduleLog) SetHeaders(headers any) *ScheduleLog {
	l.headers = headers
	return l
}

func (l *ScheduleLog) SetDomain(domain string) *ScheduleLog {
	l.domain = domain
	return l
}

func (l *ScheduleLog) SetClientIP(ip string) *ScheduleLog {
	l.clientIP = ip
	return l
}

func (l *ScheduleLog) SetResponse(response any) *ScheduleLog {
	l.response = response
	return l
}

func (l *ScheduleLog) SetExtraInfo(extra any) *ScheduleLog {
	l.extraInfo = extra
	return l
}

func (l *ScheduleLog) SetRunTime(runTime float64) *ScheduleLog {
	l.runTime = runTime
	return l
}

func (l *ScheduleLog) SetTraceID(traceID string) *ScheduleLog {
	l.traceID = traceID
	return l
}

func (l *ScheduleLog) GetTraceID() string {
	return l.traceID
}

func (l *ScheduleLog) fields() logrus.Fields {
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

func (l *ScheduleLog) PrintSchedule(err error) {
	if err != nil {
		l.err = err.Error()
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.WithFields(l.fields()).Error()
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithFields(l.fields()).Info()
	}
}

func (l *ScheduleLog) PrintPanic(err tools.Panic) {
	l.err = err.StackTrace
	logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(l.fields()).Error(err.Panic)
}
