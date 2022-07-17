package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	UserID    int64     `db:"user_id" json:"user_id"`
	Name      string    `db:"user_name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Location  string    `db:"location" json:"location"`
	CreatedAt time.Time `db:"created_datetime" json:"created_at"`
}

func serveMux() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/users/{id}", userHandler).Methods("GET")

	return logger(router)
}

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		handler.ServeHTTP(w, r)
		after := time.Now()
		duration := after.Sub(before)
		log.Printf("%s % 4s %s (%s)", r.RemoteAddr, r.Method, r.URL.Path, duration)
	})
}

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

func sendErrorJSON(w http.ResponseWriter, err error, statusCode int) error {
	log.Printf("ERROR: %+v", err)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	return enc.Encode(map[string]string{"error": err.Error()})
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("id: %v", id)

	user := &User{}
	rows, err := db.Queryx("SELECT user_id, user_name, email, location, created_datetime FROM `user` WHERE `user_id` = ?", id)
	if err != nil {
		sendErrorJSON(w, err, 500)
		return
	}

	for rows.Next() {
		if err := rows.StructScan(user); err != nil {
			sendErrorJSON(w, err, 500)
			return
		}
	}

	sendJSON(w, user, 200)
}
