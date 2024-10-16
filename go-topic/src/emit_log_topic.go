package main

import (
	"context"
	"os"
	"rehmanm/go-topic/lib"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://admin:guest@localhost:5672/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	lib.FailOnError(err, "Faild to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	body := lib.BodyFrom(os.Args)

	err = ch.PublishWithContext(
		ctx,
		"logs_topic",
		lib.SeverityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	lib.FailOnError(err, "Failed to Send Message")

	lib.Message("Message Sent " + body)
	lib.Message("Sending Message to " + lib.SeverityFrom(os.Args))

}
