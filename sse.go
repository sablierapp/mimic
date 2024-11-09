package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Flush the headers
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Start sending periodic events to the client every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Send events every second
	for {
		select {
		case <-ticker.C:
			// Send an event with the current timestamp
			fmt.Fprintf(w, "data: Current time: %s\n\n", time.Now().Format(time.RFC3339))
			flusher.Flush()
		case <-r.Context().Done():
			log.Printf("Request canceled: %v", r.Context().Err())
			return
		}
	}
}
