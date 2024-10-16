package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://admin:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	body := "Hello World!"
	defer cancel()

	for i := 0; i < 10; i++ {
		err = ch.PublishWithContext(ctx,
			"",
			q.Name,
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		failOnError(err, "Failed to publish a message")

	}

	log.Printf(" [x] Sent %s", body)

}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
