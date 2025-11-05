package mq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	r *RabbitMQ
}

func NewProducer(r *RabbitMQ) *Producer {
	return &Producer{r: r}
}

func (p *Producer) Publish(queue string, payload amqp091.Publishing) error {
	ch, err := p.r.getChannel(queue)
	if err != nil {
		return fmt.Errorf("RabbitMQ producer publish get channel error, Queue: %s, Err: %v", queue, err)
	}

	if payload.ContentType == "" {
		payload.ContentType = "application/json"
	}

	return ch.Publish(
		"",    // exchange
		queue, // routing key
		false,
		false,
		payload,
	)
}
