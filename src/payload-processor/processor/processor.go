package processor

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	// go routine for getting signals asynchronously
	go func() {
		sig := <-signals
		log.Println("Got signal: ", sig)
		cancel()
	}()
	for {
		m, err := p.kafka_reader.FetchMessage(ctx)
		if err != nil {
			log.Panic("failed to fetch message:", err)
		}
		fmt.Printf("(%s): message at offset %d: %s = %s\n", time.Now().String(), m.Offset, string(m.Key), string(m.Value))
		p.kafka_reader.CommitMessages(ctx, m)
	}

	if err := p.kafka_reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
