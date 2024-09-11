package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	bin, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, err
	}
	return sql.Open("postgres", fmt.Sprintf("postgres://postgres:%s@db:5432/example?sslmode=disable", string(bin)))
}

type RunningStatistic struct {
	ID        int       `json:"id"`
	Date      string    `json:"date"`
	Distance  string    `json:"distance"`
	Time      string    `json:"time"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Statistics []RunningStatistic `json:"statistics"`
}

func runningStatsHandler(w http.ResponseWriter, r *http.Request) {
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
	var stats []RunningStatistic
	for rows.Next() {
		var stat RunningStatistic
		if err = rows.Scan(&stat.ID, &stat.Date, &stat.Distance, &stat.Time, &stat.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stats = append(stats, stat)
	}

	// Create a response object
	response := Response{Statistics: stats}

	// Encode the response object as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createRunningStat(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var stat RunningStatistic
	if err := json.NewDecoder(r.Body).Decode(&stat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int
	err = db.QueryRow("INSERT INTO running_statistics (date, distance, time) VALUES ($1, $2, $3) RETURNING id",
		stat.Date, stat.Distance, stat.Time).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func deleteRunningStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM running_statistics WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Running Statistic not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getRunningStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var stat RunningStatistic
	err = db.QueryRow("SELECT id, date, distance, time, created_at FROM running_statistics WHERE id = $1", id).
		Scan(&stat.ID, &stat.Date, &stat.Distance, &stat.Time, &stat.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Running Statistic not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(stat)
}

func main() {
	log.Print("Prepare db...")
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	log.Print("Listening on :8000")
	r := mux.NewRouter()
	r.HandleFunc("/running-statistics", runningStatsHandler).Methods("GET")
	r.HandleFunc("/running-statistics", createRunningStat).Methods("POST")
	r.HandleFunc("/running-statistics/{id}", deleteRunningStat).Methods("DELETE")
	r.HandleFunc("/running-statistics/{id}", getRunningStat).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
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

	if _, err := db.Exec("DROP TABLE IF EXISTS running_statistics"); err != nil {
		return err
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

	return nil
}
