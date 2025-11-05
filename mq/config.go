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
		Name:  QNAME_NORMAL,
		Qos:   Qos{PrefetchCount: 1, PrefetchSize: 0, Global: false},
		Args:  nil,
		Start: true,
	},
	{
		Name: QNAME_DELAY,
		Qos:  Qos{PrefetchCount: 3, PrefetchSize: 0, Global: false},
		Args: amqp091.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": QNAME_NORMAL,
			"x-message-ttl":             int32(10 * 1000), // 延遲10秒
		},
		Start: false, // 不啟動此 Consumer
	},
}

type User struct {
	ID   uint
	Name string
}

const (
	QNAME_NORMAL string = "Normal"
	QNAME_DELAY  string = "Delay"
)

const (
	QTYPE_DEMO string = "Demo"
)
