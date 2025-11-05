package mq

import (
	"fmt"
)

type Consumer struct {
	rabbit *RabbitMQ
	router *Router
}

func NewConsumer(r *RabbitMQ, router *Router) *Consumer {
	return &Consumer{rabbit: r, router: router}
}

func (c *Consumer) Consume(queue string) error {
	ch, err := c.rabbit.GetChannel(queue)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for m := range msgs {
			if err := c.router.Route(m.Body); err != nil {
				fmt.Println("Route error:", err)
			}
		}
	}()

	return nil
}
