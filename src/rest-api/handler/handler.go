package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	ch        *amqp.Channel
	queueName string
}

func NewHandler(ch *amqp.Channel, queueName string) *Handler {
	return &Handler{
		ch:        ch,
		queueName: queueName,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		log.Fatal(err)
	}

	body := fmt.Sprint("Hello World!", r.RemoteAddr, string(bytes))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	go h.ch.PublishWithContext(ctx, "", h.queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	fmt.Fprint(w, "Hello, World!", string(bytes))
	log.Println("Request received", r.RemoteAddr, string(bytes))
}
