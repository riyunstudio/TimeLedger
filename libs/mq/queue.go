package mq

import (
	"akali/app"
	"fmt"
	"log"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

var MQ *RabbitMQ

type RabbitMQ struct {
	app *app.App

	conn   *amqp091.Connection
	chPool map[string]*amqp091.Channel
	mutex  sync.Mutex

	Consumer *Consumer
	Producer *Producer

	wg sync.WaitGroup
}

func Initialize(app *app.App) *RabbitMQ {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("RabbitMQ panic error, Err: %v", err))
		}
	}()

	conn, err := amqp091.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			app.Env.RabbitMqAccount,
			app.Env.RabbitMqPassword,
			app.Env.RabbitMqHost,
			app.Env.RabbitMqPort,
		),
	)
	if err != nil {
		panic(fmt.Errorf("RabbitMQ connect error: %w", err))
	}

	r := &RabbitMQ{
		app:    app,
		conn:   conn,
		chPool: make(map[string]*amqp091.Channel),
		wg:     sync.WaitGroup{},
	}

	// 要監聽的 Queue
	queues := []string{}

	// 初始化所有隊列
	for _, cfg := range QueuesConfig {
		ch, err := conn.Channel()
		if err != nil {
			panic(fmt.Errorf("RabbitMQ create channel for %s failed: %w", cfg.Name, err))
		}

		if err := ch.Qos(cfg.Qos.PrefetchCount, cfg.Qos.PrefetchSize, cfg.Qos.Global); err != nil {
			panic(fmt.Errorf("RabbitMQ Qos set error for %s: %w", cfg.Name, err))
		}

		_, err = ch.QueueDeclare(cfg.Name, true, false, false, false, cfg.Args)
		if err != nil {
			panic(fmt.Errorf("RabbitMQ queue declare error for %s: %w", cfg.Name, err))
		}

		r.chPool[cfg.Name] = ch

		if cfg.Start {
			queues = append(queues, cfg.Name)
		}
	}

	r.Producer = NewProducer(r)
	r.Consumer = NewConsumer(r, app, queues)

	// 寫入全局
	MQ = r

	return r
}

func (r *RabbitMQ) getChannel(queue string) (*amqp091.Channel, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ch, ok := r.chPool[queue]
	if !ok {
		return nil, fmt.Errorf("queue channel not found: %s", queue)
	}
	return ch, nil
}

func (r *RabbitMQ) Stop() {
	log.Println("RabbitMQ stopping...")

	// 等待所有正在處理的消息完成
	r.wg.Wait()
	for _, ch := range r.chPool {
		if err := ch.Close(); err != nil {
			fmt.Printf("RabbitMQ failed to close channel: %v\n", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			fmt.Printf("RabbitMQ failed to close connection: %v\n", err)
		}
	}

	log.Println("RabbitMQ Stopped")
}
