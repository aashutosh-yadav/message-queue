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
