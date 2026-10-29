package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const queueKey = "message"

var ctx = context.Background()

type Message struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

type Queue struct {
	rdb *redis.Client
}

func (q *Queue) Enqueue(ctx context.Context, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// add to the left side of the redis list
	return q.rdb.LPush(ctx, queueKey, data).Err()
}

func (q *Queue) Dequeue(ctx context.Context) (Message, error) {
	// wait until a message is availabe , then remove
	// it from the right side of the list
	result, err := q.rdb.BRPop(ctx, 0, queueKey).Result()
	if err != nil {
		return Message{}, err
	}

	var msg Message
	if err := json.Unmarshal([]byte(result[1]), &msg); err != nil {
		return Message{}, err
	}

	return msg, nil

}

func (q *Queue) PrintQueue(ctx context.Context) error {
	message, err := q.rdb.LRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	fmt.Println("-----Current Redis Queue--------")
	for i, data := range message {
		var msg Message
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			fmt.Printf("Index %d: invalid message: %v\n", i, err)
			continue
		}

		fmt.Printf("Index : %d | ID : %d | Message: %s\n",
			i, msg.ID, msg.Message)
	}
	fmt.Println("---------------------------------")

	return nil
}
