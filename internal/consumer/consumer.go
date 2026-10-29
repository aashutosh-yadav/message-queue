package consumer

import (
	"context"
	"log"
	"time"

	"message-queue/internal/queue"
)

// Run consumes messages from the queue until ctx is cancelled.
func Run(ctx context.Context, q *queue.Queue, id int) {
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

		log.Printf("worker %d: processing message %d: %s", id, msg.ID, msg.Message)

		// Simulate a slow task.
		time.Sleep(2 * time.Second)

		log.Printf("worker %d: finished processing %d", id, msg.ID)
		// Perform your actual task here.
	}
}
