package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"timeLedger/app"
	"timeLedger/app/console"
	"timeLedger/app/servers"
	_ "timeLedger/docs"

	_ "github.com/joho/godotenv/autoload"
)

// @title timeLedger 阿卡莉 模板框架
// @version 1.0
// @description API 維護文件
// @schemes http
func main() {
	// recover 防止因服務 panic 直接關閉
	defer func() {
		if err := recover(); err != nil {
			e := fmt.Sprintf("[Main panic] %v", err)
			fmt.Println(e)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())

	// 初始化 App
	app := app.Initialize()

	// 排程
	scheduler := console.Initialize(app)
	scheduler.Start()

	// MQ
	// rabbitMQ := mq.Initialize(app)

	// Websocket server (optional)
	// wsEnabled := os.Getenv("WS_ENABLED")
	// if wsEnabled != "false" {
	// 	wsServer := ws.InitializeServer(app)
	// 	wsServer.Start()

	// 	wsClient := ws.InitializeClient(fmt.Sprintf("ws://%s:%s/ws?uuid=%s", app.Env.ServerHost, app.Env.WsServerPort, app.Env.ServerName))
	// 	if err := wsClient.Start(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	// 啟動 API server
	gin := servers.Initialize(app)
	gin.Start()

	// gRPC server (optional)
	// grpcEnabled := os.Getenv("GRPC_ENABLED")
	// if grpcEnabled != "false" {
	// 	grpcSrv := grpc.Initialize(app)
	// 	grpcSrv.Start()
	// }

	// 優雅退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go func() {
		<-signals
		cancel()
		scheduler.Stop()
		// rabbitMQ.Stop()
		gin.Stop()
		fmt.Println("Exit.")
		exit <- struct{}{}
	}()
	<-ctx.Done()
	<-exit
}
