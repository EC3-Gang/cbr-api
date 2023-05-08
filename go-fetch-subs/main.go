package main

import (
	"encoding/json"
	"fmt"
	"github.com/EC3-Gang/cbr-api/redis"
	"log"
	"net/http"
)

func main() {
	client := redis.NewClient("localhost", 6379)

	go redis.PeriodicallyUpdate(client)

	http.HandleFunc("/attempts", func(w http.ResponseWriter, r *http.Request) {
		problemID := r.URL.Query().Get("problem")
		if problemID == "" {
			http.Error(w, "Missing problem ID parameter", http.StatusBadRequest)
			return
		}

		log.Println("[*] Got request for problem", problemID)

		attempts := redis.GetAttemptsFromCache(client, problemID)

		// Encode attempts as JSON and write to response writer
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(attempts); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode attempts: %v", err), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/allAttempts", func(w http.ResponseWriter, r *http.Request) {
		allAttempts := redis.GetAllAttemptsFromCache(client)

		// Encode attempts as JSON and write to response writer
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(allAttempts); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode attempts: %v", err), http.StatusInternalServerError)
			return
		}
	})

	println("[*] Go server started on port 3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
