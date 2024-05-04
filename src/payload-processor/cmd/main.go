package main

import (
	"log"
	"os"
	"time"

	"github.com/Baitinq/fs-tracer-backend/src/payload-processor/processor"
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
	processor := processor.NewProcessor(kafka_reader, 4)
	processor.ProcessMessages()
}
