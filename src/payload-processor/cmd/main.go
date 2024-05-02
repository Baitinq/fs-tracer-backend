package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Baitinq/fs-tracer-backend/src/payload-processor/processor"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitmq_password, ok := os.LookupEnv("RABBITMQ_PASSWORD")
	if !ok {
		log.Fatal("RABBITMQ_PASSWORD not set")
	}
	log.Println("RabbitMQ password", rabbitmq_password)
	conn, err := amqp.Dial(fmt.Sprintf("amqp://user:%s@rabbitmq:5672/", rabbitmq_password))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	processor := processor.NewProcessor(ch, q.Name)
	processor.ProcessMessages()
}
