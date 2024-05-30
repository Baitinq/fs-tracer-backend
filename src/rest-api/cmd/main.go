package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Baitinq/fs-tracer-backend/src/rest-api/handler"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	kafka_password, ok := os.LookupEnv("KAFKA_PASSWORD")
	if !ok {
		log.Fatal("KAFKA_PASSWORD not set")
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: "user1",
			Password: kafka_password,
		},
	}
	// Create topic
	_, err := dialer.DialLeader(context.Background(), "tcp", "kafka.default.svc.cluster.local:9092", "topic-A", 0)
	if err != nil {
		log.Fatal(err)
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
		// Async:                  true,
		AllowAutoTopicCreation: true,
	}

	db_password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("DB_PASSWORD not set")
	}
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://postgres.slpoocycjgqsuoedhkbn:%s@aws-0-eu-central-1.pooler.supabase.com:5432/postgres", db_password))
	if err != nil {
		log.Fatal("cannot initalize db client", err)
	}

	handler := handler.NewHandler(db, kafka_writer)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		fmt.Fprint(w, "Hello folks!")
	})
	mux.Handle("/api/v1/file/", handler)

	http.ListenAndServe(":8080", mux)
}
