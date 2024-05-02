package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Processor struct {
	kafka_reader *kafka.Reader
}

func NewProcessor(kafka_reader *kafka.Reader) Processor {
	log.Println("Created processor")
	return Processor{
		kafka_reader: kafka_reader,
	}
}

func (p Processor) ProcessMessages() {
	for {
		m, err := p.kafka_reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := p.kafka_reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
