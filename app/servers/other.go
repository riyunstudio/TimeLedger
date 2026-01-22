package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	"gitlab.en.mcbwvx.com/frame/zilean/logs"
)

func (s *Server) setRequestTime(ctx *gin.Context) context.Context {
	start := time.Now()
	ctx.Set(global.RequestTimeKey, start)

	return ctx
}

func (s *Server) getRequestRunTime(ctx *gin.Context) float64 {
	if val, exists := ctx.Get(global.RequestTimeKey); exists {
		if t, ok := val.(time.Time); ok {
			return float64(time.Since(t).Milliseconds()) / 1000
		}
	}
	return 0
}

func (s *Server) setTraceID(ctx *gin.Context) {
	tid := ctx.GetHeader(string(global.TraceIDKey))
	if tid == "" {
		tid, _ = s.app.Tools.NewTraceId()
	}

	ctx.Set(global.TraceIDKey, tid)
}

func (s *Server) getTraceID(ctx *gin.Context) string {
	if val, exists := ctx.Get(global.TraceIDKey); exists {
		if tid, ok := val.(string); ok {
			return tid
		}
	}
	return ""
}

func (s *Server) getRequestHeaders(ctx *gin.Context) map[string]string {
	headers := make(map[string]string)

	// 要過濾的Header欄位
	noKeylist := []string{"Cookie"}

	for k, v := range ctx.Request.Header {
		// 過濾
		skip := false
		for _, b := range noKeylist {
			if strings.EqualFold(k, b) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		if len(v) > 1 {
			headers[k] = strings.Join(v, ",")
		} else {
			headers[k] = v[0]
		}
	}
	return headers
}

func (s *Server) getQueryParams(ctx *gin.Context) map[string]string {
	params := make(map[string]string)
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 1 {
			params[k] = strings.Join(v, ",")
		} else {
			params[k] = v[0]
		}
	}
	return params
}

func (s *Server) getBodyParams(ctx *gin.Context) map[string]any {
	var body map[string]any
	if val, exists := ctx.Get(global.ArgsBodyKey); exists {
		if bytesData, ok := val.([]byte); ok {
			_ = json.Unmarshal(bytesData, &body)
		}
	}
	return body
}

func (s *Server) setBodyParams(ctx *gin.Context) {
	bodyBytes, _ := io.ReadAll(ctx.Request.Body)
	// 儲存到 context
	ctx.Set(global.ArgsBodyKey, bodyBytes)
	// 重置 Body，讓 controller 可以正常讀
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}

func (s *Server) getRet(ctx *gin.Context) global.Ret {
	if val, exists := ctx.Get(global.RetKey); exists {
		if retInfo, ok := val.(global.Ret); ok {
			return retInfo
		}
	}
	return global.Ret{}
}

func (s *Server) writeApiLog(ctx *gin.Context) {
	traceLog := logs.GinLogInit()
	traceLog.SetHeaders(s.getRequestHeaders(ctx))

	switch ctx.Request.Method {
	case http.MethodGet, http.MethodDelete:
		traceLog.SetArgs(s.getQueryParams(ctx))
	case http.MethodPost, http.MethodPut:
		traceLog.SetArgs(s.getBodyParams(ctx))
	}

	traceLog.SetUrl(ctx.Request.URL.String())
	traceLog.SetMethod(ctx.Request.Method)
	traceLog.SetDomain(ctx.Request.Host)
	traceLog.SetClientIP(ctx.ClientIP())
	traceLog.SetTraceID(s.getTraceID(ctx))
	traceLog.SetRunTime(s.getRequestRunTime(ctx))

	// 只記錄，不存取 response
	traceLog.PrintServer(0, nil, nil)
}

func (s *Server) writePanicLog(ctx *gin.Context, err tools.Panic) {
	// 初始化 TraceLog
	traceLog := logs.GinLogInit()
	traceLog.SetHeaders(s.getRequestHeaders(ctx))

	switch ctx.Request.Method {
	case http.MethodGet, http.MethodDelete:
		traceLog.SetArgs(s.getQueryParams(ctx))
	case http.MethodPost, http.MethodPut:
		traceLog.SetArgs(s.getBodyParams(ctx))
	}

	traceLog.SetUrl(ctx.Request.URL.String())
	traceLog.SetMethod(ctx.Request.Method)
	traceLog.SetDomain(ctx.Request.Host)
	traceLog.SetClientIP(ctx.ClientIP())
	traceLog.SetTraceID(s.getTraceID(ctx))
	traceLog.SetRunTime(s.getRequestRunTime(ctx))
	traceLog.PrintPanic(err)
}
