package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	log.Println(" [x] Sent", body)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello folks!")
	})
	mux.HandleFunc("/payload", handleRequest)

	http.ListenAndServe(":8080", mux)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, "Hello, World!", string(bytes))
	log.Println("Request received", r.RemoteAddr, string(bytes))
}
