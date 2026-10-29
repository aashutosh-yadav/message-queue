# message-broker

A small message queue I built in Go to learn how producers, brokers and consumers work together.

You send a message to an HTTP endpoint, it gets pushed into a Redis list, and a few goroutines pick messages up from the other end and process them in the background.

```
client -> POST /messages -> Redis list -> consumer goroutines
```

The endpoint replies `202 Accepted` as soon as the message is in Redis, it doesn't wait for the processing to finish.

## How to run it

You need Go and Redis.

Start Redis :

```bash
redis-server
```

Then run the server:

```bash
go mod tidy
go run ./cmd/server
```

It listens on `localhost:8080`. Send a message like this:

```bash
curl -X POST localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{"id":1,"message":"hello"}'
```

You'll see the consumer pick it up in the server logs.

## Project layout

```
cmd/server        starts everything
internal/producer the HTTP handler
internal/queue    talks to Redis
internal/consumer the worker goroutines
```

## Known limitations

This is a learning project, not something to run in production. The main catch is that a message is removed from Redis the moment a consumer takes it, so if the program crashes halfway through processing, that message is lost. Acknowledgements are next on my list.

## Roadmap

I want to build the same queue a few different ways to compare them: channels, a ring buffer, an append-only log, an embedded database, and then Redis (what's here now).
