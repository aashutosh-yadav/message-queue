# Message Queue Roadmap

## Producer

- HTTP endpoint that accepts messages (done)

## Queue backends (in order)

1. Channel-only queue
2. Slice / ring buffer
3. Append-only log
4. Embedded database
5. Redis list (current)

## Consumer

- Use goroutines (worker pool)

## Delivery guarantees

- Acknowledgement (ack / retry)

How it works
A consumer takes a message with BLMOVE, which moves it from message to message:processing instead of deleting it.
The consumer tries the task up to 3 times, waiting 1s, then 2s between attempts.
On success it acks, which removes the message from message:processing.
After 3 failures it copies the message to message:dead, then acks.
If the program crashes at any point, the message is still in message:processing, so nothing is lost.
