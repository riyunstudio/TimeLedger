package mq

import (
	"akali/app"
	rabbitmq "akali/global/rabbitMQ"
	"akali/libs/logs"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	r       *RabbitMQ
	handler *Handler
}

func NewConsumer(r *RabbitMQ, app *app.App, queues []string) *Consumer {
	c := &Consumer{
		r: r,
		handler: &Handler{
			app: app,
			r:   r,
		},
	}

	for _, queue := range queues {
		if err := c.consume(queue); err != nil {
			panic(fmt.Errorf("RabbitMQ new consumer %s error: %w", queue, err))
		}
	}

	return c
}

func (c *Consumer) consume(queue string) error {
	ch, err := c.r.getChannel(queue)
	if err != nil {
		return fmt.Errorf("RabbitMQ consumer consume get channel error, Queue: %s, Err: %v", queue, err)
	}

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("RabbitMQ consumer consume error, Queue: %s, Err: %v", queue, err)
	}

	go func() {
		for msg := range msgs {
			c.r.wg.Add(1)
			go func(m amqp091.Delivery) {
				defer c.r.wg.Done()

				var err error
				defer func() {
					if exception := recover(); exception != nil {
						panicErr := c.r.app.Tools.PanicParser(exception)
						err = fmt.Errorf("%s\n%s", panicErr.Panic, panicErr.StackTrace)
					}
					if err != nil {
						// 初始化 TraceLog
						traceLog := logs.TraceLogInit()
						traceLog.SetTopic("RabbitMQ")
						traceLog.SetMethod(m.Type)
						traceLog.SetArgs(string(m.Body))
						traceLog.SetTraceID(m.MessageId)
						traceLog.PrintErr(err)
					}
				}()

				switch queue {
				case rabbitmq.N_NORMAL:
					switch msg.Type {
					case rabbitmq.T_DEMO:
						c.handler.handleNormal(msg)
					default:
						err = fmt.Errorf("RabbitMQ unknown message type [%s]", msg.Type)
					}
				}
			}(msg)
		}
	}()

	return nil
}
