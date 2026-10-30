package consumer

import (
	"context"
	"errors"
	"log"
	"time"

	"message-queue/internal/queue"
)

const maxAttempts = 3

// process is the task. It fails on purpose so you can see retries work.
// Replace it with your real work later.
func process(msg queue.Message, attempt int) error {
	if msg.ID%3 == 0 {
		return errors.New("always fails")
	}
	if msg.ID%4 == 0 && attempt == 1 {
		return errors.New("fails once")
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func Run(ctx context.Context, q *queue.Queue, id int) {
	for {
		msg, raw, err := queue.Dequeue(ctx, q)
		if err != nil {
			log.Printf("worker %d: dequeue error: %v", id, err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("worker %d: processing message %d", id, msg.ID)

		// Try the task up to 3 times.
		success := false
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			err := process(msg, attempt)
			if err == nil {
				success = true
				break
			}
			log.Printf("worker %d: message %d failed (attempt %d of %d): %v",
				id, msg.ID, attempt, maxAttempts, err)
			time.Sleep(time.Second)
		}

		if success {
			log.Printf("worker %d: message %d done", id, msg.ID)
		} else {
			log.Printf("worker %d: message %d gave up, moving to dead list", id, msg.ID)
			queue.Dead(ctx, q, raw)
		}

		// Either way we are finished with it, so remove it from processing.
		err = queue.Ack(ctx, q, raw)
		if err != nil {
			log.Printf("worker %d: ack error: %v", id, err)
		}
	}
}
