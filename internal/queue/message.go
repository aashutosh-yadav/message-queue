package queue

// Message is a single item in the queue.
type Message struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}
