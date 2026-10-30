package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// three redis list
const queueKey = "message"
const processingKey = "message:processing"
const deadKey = "message:dead"

type Queue struct {
	rdb *redis.Client
}

// New creates a Queue that uses the given Redis client.
func New(rdb *redis.Client) *Queue {
	return &Queue{rdb: rdb}
}

// Enqueue adds a message to the left side of the list.
func Enqueue(ctx context.Context, q *Queue, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return q.rdb.LPush(ctx, queueKey, data).Err()
}

// Dequeue waits for a message and MOVES it to the processing list
// (it is not deleted yet). It returns the message and its raw text,
// which we need later for Ack and Dead.
func Dequeue(ctx context.Context, q *Queue) (Message, string, error) {
	raw, err := q.rdb.BLMove(ctx, queueKey, processingKey, "RIGHT", "LEFT", 0).Result()
	if err != nil {
		return Message{}, "", err
	}

	var msg Message
	err = json.Unmarshal([]byte(raw), &msg)
	if err != nil {
		// Broken message: put it in the dead list so it doesn't get stuck.
		Dead(ctx, q, raw)
		Ack(ctx, q, raw)
		return Message{}, "", err
	}

	return msg, raw, nil
}

// Ack removes a finished message from the processing list.
func Ack(ctx context.Context, q *Queue, raw string) error {
	return q.rdb.LRem(ctx, processingKey, 1, raw).Err()
}

// Dead stores a message that failed every attempt.
func Dead(ctx context.Context, q *Queue, raw string) error {
	return q.rdb.LPush(ctx, deadKey, raw).Err()
}

func PrintQueue(ctx context.Context, q *Queue) error {
	messages, err := q.rdb.LRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	fmt.Println("-----Current Redis Queue--------")
	for i, data := range messages {
		var msg Message
		err := json.Unmarshal([]byte(data), &msg)
		if err != nil {
			fmt.Printf("Index %d: invalid message: %v\n", i, err)
			continue
		}
		fmt.Printf("Index : %d | ID : %d | Message: %s\n", i, msg.ID, msg.Message)
	}
	fmt.Println("---------------------------------")

	return nil
}
