package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/V-Lam/High-Throughput-Event-Broker/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

// structure for incoming message format
type EventPayload struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

// buffered channel
var eventQueue chan EventPayload

func main() {

	ctx := context.Background()

	connStr := "postgres://admin:secretpassword@localhost:5432/event_broker?sslmode=disable"

	pool, err := database.InitPool(ctx, connStr)
	if err != nil {
		slog.Error("unable to initialize connection pool: %w", err)
		os.Exit(1)
	}

	slog.Info("Server successfully booted and connected to the backend cluster", "pool_status", "active")

	eventQueue = make(chan EventPayload, 10000)

	go startWorker(ctx, pool) // launch background go routine

	// Keep the main thread alive temporarily
	// select {}

	http.HandleFunc("/publish", handlePublish)

	slog.Info("HTTP ingestion server listening on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("unable to start HTTP server on port 8080", slog.Any("err", err))
		os.Exit(1)
	}

}

func startWorker(ctx context.Context, pool *pgxpool.Pool) {
	slog.Info("Background event worker pool spinning up...")

	for currentPayload := range eventQueue {
		slog.Info("event successfully popped off the queue")
		_ = currentPayload
	}
}

func handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
