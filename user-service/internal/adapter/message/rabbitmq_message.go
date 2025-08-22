package message

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
	"user-service/config"
)

func PublishMessage(email, message, notif_type string) error {
	conn, err := config.NewConfig().NewRabbitMQ()
	if err != nil {
		log.Fatalf("[PublishMessage-1] Failed to connect to RabbitMQ: %v", err)
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[PublishMessage-2] Failed to open a channel: %v", err)
	}
	defer ch.Close()
	queue, err := ch.QueueDeclare(notif_type, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[PublishMessage-3] Failed to declare a queue: %v", err)
		return err
	}

	notification := map[string]string{
		"email":   email,
		"message": message,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		log.Fatalf("[PublishMessage-4] Failed to marshal notification: %v", err)
		return err
	}

	return ch.Publish(
		"",
		queue.Name, // exchange
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     amqp.Table{"x-delay": 1000}, // Delay in milliseconds

		},
	)
}
