package main

import (
	"context"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
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

	q := &Queue{rdb: rdb}

	// Start the background consumer.
	go worker(q)

	// Register the POST endpoint.
	http.HandleFunc("/messages", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		enqueueHandler(w, r, q)
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
