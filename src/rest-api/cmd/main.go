package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Baitinq/fs-tracer-backend/src/rest-api/handler"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitmq_password, ok := os.LookupEnv("RABBITMQ_PASSWORD")
	if !ok {
		panic("RABBITMQ_PASSWORD not set")
	}
	log.Println("RabbitMQ password", rabbitmq_password)
	conn, err := amqp.Dial(fmt.Sprintf("amqp://user:%s@rabbitmq:5672/", rabbitmq_password))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
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
		panic(err)
	}
	handler := handler.NewHandler(ch, q.Name)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello folks!")
	})
	mux.Handle("/payload", handler)

	http.ListenAndServe(":8080", mux)
}
