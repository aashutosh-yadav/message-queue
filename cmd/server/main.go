package main

import (
	"context"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"

	"message-queue/internal/consumer"
	"message-queue/internal/producer"
	"message-queue/internal/queue"
)

func main() {
	ctx := context.Background()

	// Connect to Redis.
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Verify the connection.
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Cannot connect to Redis: %v", err)
	}

	q := queue.New(rdb)

	// Start the background consumer.
	const numWorker = 5

	for i := 1; i <= numWorker; i++ {
		go consumer.Run(ctx, q, i)
	}

	// Register the POST endpoint.
	http.HandleFunc("/messages", handler.Enqueue(q))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
