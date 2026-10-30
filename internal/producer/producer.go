package producer

import (
	"encoding/json"
	"net/http"

	"message-queue/internal/queue"
)

// perfoms enqueue
func Enqueue(w http.ResponseWriter, r *http.Request, q *queue.Queue) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	// Read the JSON message
	var msg queue.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the message
	if msg.ID <= 0 || msg.Message == "" {
		http.Error(w, "Invalid message", http.StatusBadRequest)
		return
	}

	// Add message to Redis
	err = queue.Enqueue(r.Context(), q, msg)
	if err != nil {
		http.Error(w, "Failed to enqueue message", http.StatusInternalServerError)
		return
	}

	// // Print messages still waiting in Redis
	// err = queue.PrintQueue(r.Context(), q)
	// if err != nil {
	// 	http.Error(w, "Failed to read queue", http.StatusInternalServerError)
	// 	return
	// }

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(map[string]any{
		"status": "Message accepted",
		"id":     msg.ID,
	})
}
