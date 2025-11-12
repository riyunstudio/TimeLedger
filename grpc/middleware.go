package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Grpc) InitMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx = s.setRequestTime(ctx)
	s.wg.Add(1)
	defer func() {
		s.wg.Done()
	}()

	return handler(ctx, req)
}

func (s *Grpc) RecoverMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			// 整理 Panic內容
			e := s.app.Tools.PanicParser(r)
			// 紀錄 TraceLog
			s.writePanicLog(ctx, info.FullMethod, req, e)
			// 回傳錯誤
			err = status.Errorf(codes.Unknown, "Panic error, Err: %v", e.Panic)
		}
	}()

	return handler(ctx, req)
}

func (s *Grpc) MainMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// 避免空請求或 HTTP/2 preface
	if info.FullMethod == "" || req == nil {
		return handler(ctx, req)
	}

	// 將 metadata 注入到 ctx
	ctx = s.injectMetadataToContext(ctx)

	var resp any
	var err error

	done := make(chan struct{})
	go func() {
		resp, err = handler(ctx, req)
		close(done)
	}()

	select {
	case <-done:
		// 執行完成
	case <-ctx.Done():
		// Context 被取消，例如超時或 client 取消
		err = status.Errorf(codes.DeadlineExceeded, "context canceled during execution")
		resp = nil
	}

	// 紀錄 TraceLog
	s.writeApiLog(ctx, info.FullMethod, req, resp, err)

	return resp, err
}
