package main

import (
	"context"
	"log"
	"os"
	"rehmanm/go-rabbitmq/rpc/lib"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	n := lib.BodyFrom(os.Args)

	log.Printf(" [x] Requesting fib(%d)", n)

	res, err := fibonacciRPC(n)

	lib.FailOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)

}

func fibonacciRPC(n int) (res int, err error) {
	conn, ch := lib.GetRabbitMqConnectionandChannel()

	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	lib.FailOnError(err, "Failed to declare a queue")

	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	lib.FailOnError(err, "Failed to register a consumer")

	corrId := lib.RandomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(lib.IntToString(n)),
		})

	lib.FailOnError(err, "Failed to publish a message")

	for d := range msg {
		if corrId == d.CorrelationId {
			res = lib.StringToInt(string(d.Body))
			lib.FailOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}
