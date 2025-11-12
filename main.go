package main

import (
	"akali/app"
	"akali/app/console"
	"akali/app/servers"
	_ "akali/docs"
	"akali/grpc"
	"akali/libs/mq"
	"akali/libs/ws"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

// @title Akali 阿卡莉 模板框架
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
	rabbitMQ := mq.Initialize(app)
	// body, _ := json.Marshal(rabbitmq.NormalDemoUser{ID: 1, Name: "Test1"})
	// rabbitMQ.Producer.Publish(rabbitmq.N_NORMAL, amqp091.Publishing{Type: rabbitmq.T_DEMO, Body: body, MessageId: "TRACE_ID_XXXXX"})
	// body, _ = json.Marshal(rabbitmq.NormalDemoUser{ID: 2, Name: "Test2"})
	// rabbitMQ.Producer.Publish(rabbitmq.N_DELAY, amqp091.Publishing{Type: rabbitmq.T_DEMO, Body: body, MessageId: "TRACE_ID_XXXXX"})

	// Websocket server
	ws := ws.Initialize(app)
	ws.Start()

	// 啟動 API server
	gin := servers.Initialize(app)
	gin.Start()

	grpc := grpc.Initialize(app)
	grpc.Start()

	// 優雅退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go func() {
		<-signals
		// 優雅退出要結束的程式寫在這 Ex.關閉連線
		cancel()
		scheduler.Stop()
		rabbitMQ.Stop()
		ws.Stop()
		gin.Stop()
		grpc.Stop()
		fmt.Println("Exit.")
		exit <- struct{}{}
	}()
	<-ctx.Done()
	<-exit
}
