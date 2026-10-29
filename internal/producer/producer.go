package handler

import (
	"encoding/json"
	"net/http"

	"message-queue/internal/queue"
)

// Enqueue returns an HTTP handler that accepts a JSON message
// and pushes it onto the queue.
func Enqueue(q *queue.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow only POST requests.
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		// Read the JSON message.
		var msg queue.Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate the message.
		if msg.ID <= 0 || msg.Message == "" {
			http.Error(w, "Invalid message", http.StatusBadRequest)
			return
		}

		// Add message to Redis.
		if err := q.Enqueue(r.Context(), msg); err != nil {
			http.Error(w, "Failed to enqueue message",
				http.StatusInternalServerError)
			return
		}

		// Print messages still waiting in Redis.
		if err := q.PrintQueue(r.Context()); err != nil {
			http.Error(w, "Failed to read queue",
				http.StatusInternalServerError)
			return
		}

		// Send response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		json.NewEncoder(w).Encode(map[string]any{
			"status": "Message accepted",
			"id":     msg.ID,
		})
	}
}
