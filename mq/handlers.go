package mq

import (
	"akali/app"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	app *app.App
}

func (h *Handler) normal(msg amqp091.Delivery) error {
	switch msg.Type {
	case QTYPE_DEMO:
		var u User
		if err := json.Unmarshal(msg.Body, &u); err != nil {
			return err
		}
		fmt.Printf("[Normal.%s] Received user: ID=%d, Name=%s, %v\n", msg.Type, u.ID, u.Name, h.app.Tools.Locat())
	default:
		return fmt.Errorf("RabbitMQ unknown message type [%s]", msg.Type)
	}

	return nil
}
