package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type RequestLog struct {
	ID        int       `json:"id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

var db *sql.DB

func initDB() {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to PostgreSQL")
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec(
		"INSERT INTO requests (method, path, timestamp) VALUES ($1, $2, $3)",
		r.Method,
		r.URL.Path,
		time.Now(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, method, path, timestamp FROM requests ORDER BY id DESC LIMIT 10")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []RequestLog
	for rows.Next() {
		var l RequestLog
		err := rows.Scan(&l.ID, &l.Method, &l.Path, &l.Timestamp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logs = append(logs, l)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func main() {
	initDB()
	http.HandleFunc("/", logRequest)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
