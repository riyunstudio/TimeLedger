package mq

import "github.com/rabbitmq/amqp091-go"

type Qos struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

type QueueConfig struct {
	Name  string
	Qos   Qos
	Args  amqp091.Table
	Start bool
}

var QueuesConfig = []QueueConfig{
	{
		Name:  "normal",
		Qos:   Qos{PrefetchCount: 1, PrefetchSize: 0, Global: false},
		Args:  nil,
		Start: true,
	},
	{
		Name: "delay",
		Qos:  Qos{PrefetchCount: 3, PrefetchSize: 0, Global: false},
		Args: amqp091.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "normal",
			"x-message-ttl":             int32(10000),
		},
		Start: false, // 不啟動此 Consumer
	},
}

type User struct {
	Type string
	ID   uint
	Name string
}
