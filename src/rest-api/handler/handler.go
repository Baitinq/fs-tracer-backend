package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	db           DB
	kafka_writer *kafka.Writer
}

func NewHandler(db *sqlx.DB, kafka_writer *kafka.Writer) Handler {
	return Handler{
		db:           NewDB(db),
		kafka_writer: kafka_writer,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	_, filePath, ok := strings.Cut(r.URL.Path, "/file/")
	if !ok {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	log.Println("File path: ", filePath)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	file, err := h.db.GetLatestFileByPath(ctx, filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "File: ", file)
}

func (h Handler) handlePost(w http.ResponseWriter, r *http.Request) {
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
