package main

import (
	"akali/app"
	"akali/app/console"
	ginSrv "akali/app/servers"
	"akali/mq"
	rpcSrv "akali/rpc/servers"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rabbitmq/amqp091-go"
)

// @title Akali 阿卡莉 模板框架_1
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

	queue := mq.NewRabbitMQ(app)

	body, _ := json.Marshal(mq.User{
		ID:   1,
		Name: "Test1",
	})

	// 發送測試消息
	queue.Producer.Publish("normal", amqp091.Publishing{
		Type: "",
		Body: body,
	})
	// queue.Producer.Publish("delay", &mq.User{
	// 	Type: "normal",
	// 	ID:   2,
	// 	Name: "Test2",
	// })

	// 啟動 API server
	srv := ginSrv.Initialize(app)
	srv.Start()

	rpcServer := rpcSrv.Initialize(app)
	rpcServer.Start()

	// 優雅退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go func() {
		<-signals
		// 優雅退出要結束的程式寫在這 Ex.關閉連線
		cancel()
		srv.Stop()
		rpcServer.Stop()
		scheduler.Stop()
		fmt.Println("Exit.")
		exit <- struct{}{}
	}()
	<-ctx.Done()
	<-exit
}
