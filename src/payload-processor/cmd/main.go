package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Baitinq/fs-tracer-backend/src/payload-processor/processor"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	kafka_password, ok := os.LookupEnv("KAFKA_PASSWORD")
	if !ok {
		log.Fatal("KAFKA_PASSWORD not set")
	}
	kafka_reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka.default.svc.cluster.local:9092"},
		Dialer: &kafka.Dialer{
			SASLMechanism: plain.Mechanism{
				Username: "user1",
				Password: kafka_password,
			},
			Timeout:   10 * time.Second,
			DualStack: true,
		},
		Topic:    "topic-A",
		GroupID:  "group-A",
		MaxBytes: 10e6, // 10MB
	})

	db_password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("DB_PASSWORD not set")
	}
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://postgres.slpoocycjgqsuoedhkbn:%s@aws-0-eu-central-1.pooler.supabase.com:5432/postgres", db_password))
	if err != nil {
		log.Fatal("cannot initalize db client", err)
	}

	processor := processor.NewProcessor(kafka_reader, db, 4)
	processor.ProcessMessages()
}
