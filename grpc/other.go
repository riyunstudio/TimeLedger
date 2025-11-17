package grpc

import (
	"akali/global"
	"akali/libs/logs"
	"context"
	"time"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	"google.golang.org/grpc/metadata"
)

func (s *Grpc) setRequestTime(ctx context.Context) context.Context {
	start := time.Now()
	ctx = context.WithValue(ctx, global.RequestTimeKey, start)

	return ctx
}

func (s *Grpc) getRequestRunTime(ctx context.Context) float64 {
	if startTime, ok := ctx.Value(global.RequestTimeKey).(time.Time); ok {
		return float64(time.Since(startTime).Milliseconds()) / 1000
	}

	return 0
}

func (s *Grpc) injectMetadataToContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx // 沒有 metadata，直接返回原始 ctx
	}

	// 完整 MD
	metaMap := make(map[string]string)
	type ctxKey string

	// 將每個 metadata key-value 放進 ctx
	for k, vals := range md {
		if len(vals) > 0 {
			ctx = context.WithValue(ctx, ctxKey(k), vals[0])
			metaMap[k] = vals[0]
		}
	}

	// 完整 metadata 放進 ctx
	ctx = context.WithValue(ctx, global.MeteDataKey, metaMap)

	// 最後才檢查MD有沒有給tid, 沒給不新增到MD, 直接新增到ctx提供後續流程使用
	if _, ok := metaMap[string(global.TraceIDKey)]; !ok {
		tid, err := s.app.Tools.NewTraceId()
		if err != nil {
			tid = ""
		}
		ctx = context.WithValue(ctx, global.TraceIDKey, tid)
	}

	return ctx
}

func (s *Grpc) writeApiLog(ctx context.Context, method string, req any, resp any, err error) {
	// 初始化 TraceLog
	traceLog := logs.TraceLogInit()
	traceLog.SetTopic("server_gRPC")

	headers, _ := ctx.Value(global.MeteDataKey).(map[string]string)
	traceLog.SetHeaders(headers)
	traceLog.SetArgs(req)
	traceLog.SetUrl(method)
	traceLog.SetMethod(method)

	if clientIP, ok := headers[string(global.XForwardForKey)]; ok {
		traceLog.SetClientIP(clientIP)
	}
	if tid, ok := ctx.Value(global.TraceIDKey).(string); ok {
		traceLog.SetTraceID(tid)
	}
	traceLog.SetRunTime(s.getRequestRunTime(ctx))
	traceLog.PrintGrpc(resp, err)
}

func (s *Grpc) writePanicLog(ctx context.Context, method string, req any, err tools.Panic) {
	// 初始化 TraceLog
	traceLog := logs.TraceLogInit()
	traceLog.SetTopic("server_gRPC")

	headers, _ := ctx.Value(global.MeteDataKey).(map[string]string)
	traceLog.SetHeaders(headers)
	traceLog.SetArgs(req)
	traceLog.SetUrl(method)
	traceLog.SetMethod(method)

	if clientIP, ok := headers[string(global.XForwardForKey)]; ok {
		traceLog.SetClientIP(clientIP)
	}
	if tid, ok := ctx.Value(global.TraceIDKey).(string); ok {
		traceLog.SetTraceID(tid)
	}
	traceLog.SetRunTime(s.getRequestRunTime(ctx))
	traceLog.PrintPanic(err)
}
