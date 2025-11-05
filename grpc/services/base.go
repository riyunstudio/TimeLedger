package services

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaseService struct{}

// RunWithTimeout 封裝 API 超時處理
func RunWithTimeout[T any](ctx context.Context, timeout time.Duration, fn func(ctx context.Context, do func(func() error) error) (T, error)) (T, error) {
	var zero T
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultCh := make(chan struct {
		data T
		err  error
	}, 1)

	// do 函式用來包裹耗時操作，自動檢查 context
	do := func(op func() error) error {
		select {
		case <-ctx.Done():
			return status.Error(codes.DeadlineExceeded, "context canceled")
		default:
			return op()
		}
	}

	go func() {
		data, err := fn(ctx, do)
		resultCh <- struct {
			data T
			err  error
		}{data, err}
	}()

	select {
	case r := <-resultCh:
		return r.data, r.err
	case <-ctx.Done():
		return zero, status.Error(codes.DeadlineExceeded, "request timeout")
	}
}
