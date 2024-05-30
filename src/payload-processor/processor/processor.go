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

	"github.com/Baitinq/fs-tracer-backend/lib"
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

		user_id, err := getHeaderValue(m.Headers, "user_id")
		if err != nil {
			log.Fatal("failed to get user_id from headers:", err)
		}

		// TODO: Remove after testing
		if string(m.Value) == "" {
			m.Value = []byte(fmt.Sprintf(`[{
				"absolute_path": "/home/user/file.txt",
				"contents": "Hello, World!",
				"timestamp": "%s"
			}]`, time.Now().Format(time.RFC3339)))
		}

		err = p.handleMessage(ctx, m, string(user_id))
		if err != nil {
			log.Println("failed to handle message:", err)
			p.handleError(ctx, m, err)
			return
		}
		p.kafka_reader.CommitMessages(ctx, m)
	}
}

func (p Processor) handleMessage(ctx context.Context, m kafka.Message, user_id string) error {
	fmt.Printf("(%s): message at paritition %d: offset %d: %s = %s, user_id = %s\n", time.Now().String(), m.Partition, m.Offset, string(m.Key), string(m.Value), user_id)

	var files []lib.File
	err := json.Unmarshal(m.Value, &files)
	if err != nil {
		return err
	}

	err = p.db.InsertFiles(ctx, files, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (p Processor) handleError(ctx context.Context, m kafka.Message, err error) {
	switch err {
	// TODO: If its a recoverable error, don't commit.
	default:
		p.kafka_reader.CommitMessages(ctx, m)
	}
}

func getHeaderValue(headers []kafka.Header, key string) ([]byte, error) {
	for _, header := range headers {
		if header.Key == key {
			return header.Value, nil
		}
	}
	return []byte{}, fmt.Errorf("Header %s not found", key)
}
