package lib

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return Fib(n-1) + Fib(n-2)
	}
}

func GetRabbitMqConnectionandChannel() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://admin:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	return conn, ch
}

func StringToInt(s string) int {
	data, err := strconv.Atoi(s)
	FailOnError(err, "Failed to convert string to int")
	return data
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func BodyFrom(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	FailOnError(err, "Failed to convert arg to integer")
	return n
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}
