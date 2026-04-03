package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitDB()
	defer DB.Close()

	r := mux.NewRouter()

	// route methods including OPTIONS so middleware gets invoked
	r.HandleFunc("/goals", GetGoals).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/goals", CreateGoal).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/health", HealthCheck).Methods(http.MethodGet, http.MethodOptions)

	// fallback options (if you prefer)
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Use(withCORS)

	log.Println("🚀 Starting http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
