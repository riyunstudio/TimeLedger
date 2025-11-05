package mq

import (
	"akali/app"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	app *app.App

	conn   *amqp091.Connection
	chPool map[string]*amqp091.Channel
	mutex  sync.Mutex

	Consumer *Consumer
	Producer *Producer
}

func NewRabbitMQ(app *app.App) *RabbitMQ {
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
	}

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
	}

	// Consumer 路由
	router := NewRouter()
	router.RegisterBulk(GetRegistrations())

	// Consumer
	r.Consumer = NewConsumer(r, router)
	for _, q := range []string{"normal"} {
		if err := r.Consumer.Consume(q); err != nil {
			panic(fmt.Errorf("RabbitMQ new consume error, Err: %v", err))
		}
	}

	// Producer
	r.Producer = NewProducer(r)

	return r
}

func (r *RabbitMQ) GetChannel(queue string) (*amqp091.Channel, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ch, ok := r.chPool[queue]
	if !ok {
		return nil, fmt.Errorf("queue channel not found: %s", queue)
	}
	return ch, nil
}

func (r *RabbitMQ) Close() error {
	for _, ch := range r.chPool {
		ch.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}
