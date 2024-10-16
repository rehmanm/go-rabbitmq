package main

import (
	"context"
	"log"
	"rehmanm/go-rabbitmq/rpc/lib"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	fibCache := make(map[int]int)

	conn, ch := lib.GetRabbitMqConnectionandChannel()

	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)

	lib.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,
		0,
		false)

	lib.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)

	lib.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for d := range msgs {
			n := lib.StringToInt(string(d.Body))
			response := 0
			if val, ok := fibCache[n]; ok {
				log.Printf(" [.] fib(%d) from Cache = %d", n, val)
				response = val
			} else {
				response = lib.Fib(n)
				fibCache[n] = response
			}

			log.Printf(" [.] fib(%d) = %d", n, response)

			err = ch.PublishWithContext(ctx,
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(lib.IntToString(response)),
				})

			lib.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever

}
