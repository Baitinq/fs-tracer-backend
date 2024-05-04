package processor

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type Processor struct {
	kafka_reader *kafka.Reader
	concurrency  int
}

func NewProcessor(kafka_reader *kafka.Reader, concurrency int) Processor {
	log.Println("Created processor with concurrency: ", concurrency)
	return Processor{
		kafka_reader: kafka_reader,
		concurrency:  concurrency,
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

	wg := sync.WaitGroup{}
	wg.Add(p.concurrency)
	for i := 0; i < p.concurrency; i++ {
		go func() {
			defer wg.Done()
			p.process(ctx, cancel)
		}()
	}

	wg.Wait()

	if err := p.kafka_reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func (p Processor) process(ctx context.Context, cancel context.CancelFunc) {
	for {
		m, err := p.kafka_reader.FetchMessage(ctx)
		if err != nil {
			log.Println("failed to fetch message:", err)
			cancel()
			break
		}
		err = p.handleMessage(m)
		if err != nil {
			log.Println("failed to handle message:", err)
			continue
		}
		p.kafka_reader.CommitMessages(ctx, m)
	}
}

func (p Processor) handleMessage(m kafka.Message) error {
	fmt.Printf("(%s): message at offset %d: %s = %s\n", time.Now().String(), m.Offset, string(m.Key), string(m.Value))
	return nil
}
