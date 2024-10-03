package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	// Load the database password from a secret file
	password, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, err
	}
	return sql.Open("postgres", fmt.Sprintf("postgres://postgres:%s@db:5432/example?sslmode=disable", string(password)))
}

type RunningStatistic struct {
	ID        int       `json:"id"`
	Date      string    `json:"date"`
	Distance  string    `json:"distance"`
	Time      string    `json:"time"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	log.Print("Setting up the database...")
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	r.HandleFunc("/running-statistics", runningStatsHandler).Methods(http.MethodGet)
	r.HandleFunc("/running-statistics", createRunningStatHandler).Methods(http.MethodPost)

	// Apply CORS middleware directly to the router
	corsHandler := corsMiddleware(r)

	log.Print("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", corsHandler))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello from running app")
}

// CORS middleware to handle CORS preflight requests
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
func runningStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	db, err := connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, distance, time, created_at FROM running_statistics")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stats []RunningStatistic
	for rows.Next() {
		var stat RunningStatistic
		if err := rows.Scan(&stat.ID, &stat.Date, &stat.Distance, &stat.Time, &stat.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stats = append(stats, stat)
	}

	response := struct {
		Statistics []RunningStatistic `json:"statistics"`
	}{Statistics: stats}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createRunningStatHandler(w http.ResponseWriter, r *http.Request) {
	var stat RunningStatistic
	// Decode JSON request body and handle any decoding errors
	if err := json.NewDecoder(r.Body).Decode(&stat); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	db, err := connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Validate input values here
	if stat.Date == "" || stat.Distance == "" || stat.Time == "" {
		http.Error(w, "Invalid input: Date, Distance, and Time are required", http.StatusBadRequest)
		return
	}

	// Execute the query to insert the new running statistic
	query := "INSERT INTO running_statistics (date, distance, time, created_at) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	if err := db.QueryRow(query, stat.Date, stat.Distance, stat.Time, time.Now()).Scan(&stat.ID, &stat.CreatedAt); err != nil {
		log.Printf("Failed to insert record: %v", err)
		http.Error(w, "Failed to insert record: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Successful insertion, respond with the new record
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stat)
}
func prepare() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS running_statistics (
        id SERIAL PRIMARY KEY,
        date VARCHAR(10),
        distance VARCHAR(10),
        time VARCHAR(10),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`); err != nil {
		return err
	}

	fmt.Println("Database initialized successfully.")
	return nil
}
