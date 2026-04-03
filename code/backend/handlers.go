package main

import (
	"encoding/json"
	"net/http"
)

type Goal struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func CreateGoal(w http.ResponseWriter, r *http.Request) {
	var goal Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if goal.Title == "" {
		http.Error(w, `{"error":"Title required"}`, http.StatusBadRequest)
		return
	}

	err := DB.QueryRow(
		"INSERT INTO goals(title, status) VALUES($1, $2) RETURNING id",
		goal.Title, "active",
	).Scan(&goal.ID)

	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(goal)
}

func GetGoals(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, title FROM goals WHERE status='active' ORDER BY id DESC")
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var g Goal
		rows.Scan(&g.ID, &g.Title)
		goals = append(goals, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goals)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "eOffice API",
	})
}
