package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type Processor struct {
	kafka_reader *kafka.Reader
	db           DB
	concurrency  int
}

func NewProcessor(kafka_reader *kafka.Reader, db *sqlx.DB, concurrency int) Processor {
	log.Println("Created processor with concurrency: ", concurrency)
	return Processor{
		kafka_reader: kafka_reader,
		db:           NewDB(db),
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

		// TODO: Remove after testing
		if string(m.Value) == "" {
			m.Value = []byte(fmt.Sprintf(`{
				"user_id": "%s",
				"absolute_path": "/home/user/file.txt",
				"contents": "Hello, World!",
				"timestamp": "%s"
			}`, uuid.New(), time.Now().Format(time.RFC3339)))
		}

		err = p.handleMessage(ctx, m)
		if err != nil {
			log.Println("failed to handle message:", err)
			continue
		}
		p.kafka_reader.CommitMessages(ctx, m)
	}
}

func (p Processor) handleMessage(ctx context.Context, m kafka.Message) error {
	fmt.Printf("(%s): message at paritition %d: offset %d: %s = %s\n", time.Now().String(), m.Partition, m.Offset, string(m.Key), string(m.Value))

	var file File
	err := json.Unmarshal(m.Value, &file)
	if err != nil {
		return err
	}

	err = p.db.InsertFile(ctx, file)
	if err != nil {
		return err
	}
	return nil
}
