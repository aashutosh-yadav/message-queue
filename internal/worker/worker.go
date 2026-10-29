package worker

import (
	"context"
	"log"
	"time"

	"message-queue/internal/queue"
)

// Run consumes messages from the queue until ctx is cancelled.
func Run(ctx context.Context, q *queue.Queue) {
	for {
		if ctx.Err() != nil {
			return
		}

		msg, err := q.Dequeue(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Failed to dequeue: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Processing message %d: %s", msg.ID, msg.Message)

		// Perform your actual task here.
	}
}
