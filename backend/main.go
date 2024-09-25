package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Record represents a running statistic
type Record struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Distance  string `json:"distance"`
	Time      string `json:"time"`
	CreatedAt string `json:"createdAt"`
}

// Response represents a structured JSON response
type Response struct {
	Message string  `json:"message"`
	Record  *Record `json:"record,omitempty"`
}

// In-memory storage
var (
	records []Record
	mu      sync.Mutex
)

func main() {
	http.HandleFunc("/api/records", corsMiddleware(handleRecords))

	// Start the server
	http.ListenAndServe(":8080", nil)
}

// CORS middleware to handle CORS preflight requests
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// handleRecords handles incoming record submissions
func handleRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var record Record
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		records = append(records, record)
		mu.Unlock()

		// Create a response object
		response := Response{
			Message: "Record added successfully",
			Record:  &record,
		}

		// Respond with a JSON object after a successful record addition
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	} else if r.Method == http.MethodGet {
		mu.Lock()
		json.NewEncoder(w).Encode(records)
		mu.Unlock()
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
