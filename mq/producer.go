package mq

import (
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	rabbit *RabbitMQ
}

func NewProducer(r *RabbitMQ) *Producer {
	return &Producer{rabbit: r}
}

func (p *Producer) Publish(queue string, payload amqp091.Publishing) error {
	ch, err := p.rabbit.GetChannel(queue)
	if err != nil {
		return err
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
