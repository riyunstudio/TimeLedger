package grpc

import (
	"fmt"
	"log"
	"net"
	"sync"
	"timeLedger/app"

	"google.golang.org/grpc"
)

type Grpc struct {
	app *app.App

	srv *grpc.Server   // Shutdown時直接使用
	wg  sync.WaitGroup // For 優雅退出
}

func Initialize(app *app.App) *Grpc {
	return &Grpc{
		app: app,
	}
}

func (s *Grpc) registerServices() {
}

func (s *Grpc) Start() {
	s.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.InitMiddleware,    // wg
			s.RecoverMiddleware, // panic 保護
			s.MainMiddleware,    // 主程式
		),
	)
	s.registerServices()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(fmt.Errorf("gRPC server start error, Err: failed to listen: %v", err))
	}

	go func() {
		if err := s.srv.Serve(lis); err != nil {
			panic(fmt.Errorf("gRPC server failed to serve: %v", err))
		}
	}()

	log.Println("gRPC server started")
}

func (s *Grpc) Stop() {
	log.Println("gRPC server stopping...")

	// 1. 停止 gRPC 伺服器的相關操作
	s.srv.GracefulStop()

	// 2. 等待所有請求完成
	s.wg.Wait()
	log.Println("gRPC server graceful Stop complete")
}
