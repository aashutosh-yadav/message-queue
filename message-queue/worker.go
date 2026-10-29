package main

import (
	"context"
	"log"
	"time"
)

func worker(q *Queue) {
	ctx := context.Background()

	for {
		msg, err := q.Dequeue(ctx)
		if err != nil {
			log.Printf("Failed to dequeue: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Processing message %d: %s",
			msg.ID, msg.Message)

		// Perform your actual task here.
	}
}
