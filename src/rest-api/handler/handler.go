package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type Handler struct {
	kafka_writer *kafka.Writer
}

func NewHandler(kafka_writer *kafka.Writer) Handler {
	return Handler{
		kafka_writer: kafka_writer,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.kafka_writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("key"), //TODO:This routes to a partition. We should probably route by agent UUID TODO: wont this negate having multiple topics
		Value: bytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, "Hello, World!", string(bytes))
	log.Println("Request received :)", r.RemoteAddr, string(bytes))
}
