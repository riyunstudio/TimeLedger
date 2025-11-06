package mq

import (
	"akali/app"
	rabbitmq "akali/global/rabbitMQ"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	app *app.App
	r   *RabbitMQ
}

func (h *Handler) handleNormal(msg amqp091.Delivery) error {
	var u rabbitmq.NormalDemoUser
	if err := json.Unmarshal(msg.Body, &u); err != nil {
		return err
	}
	fmt.Printf("[Normal.%s] Received user: ID=%d, Name=%s, %v\n", msg.Type, u.ID, u.Name, h.app.Tools.Locat())

	return nil
}
