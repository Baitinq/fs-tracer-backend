package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Baitinq/fs-tracer-backend/src/rest-api/handler"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	kafka_password, ok := os.LookupEnv("KAFKA_PASSWORD")
	if !ok {
		log.Fatal("KAFKA_PASSWORD not set")
	}

	kafka_writer := &kafka.Writer{
		Addr: kafka.TCP("kafka.default.svc.cluster.local:9092"),
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: "user1",
				Password: kafka_password,
			},
		},
		Topic:    "topic-A",
		Balancer: &kafka.LeastBytes{},
		// Async:                  true, //TODO: Creat the topic beforehand, if not this doesnt work
		AllowAutoTopicCreation: true,
	}

	handler := handler.NewHandler(kafka_writer)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello folks!")
	})
	mux.Handle("/payload", handler)

	http.ListenAndServe(":8080", mux)
}
