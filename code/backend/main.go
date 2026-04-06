package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

func main() {
	InitDB()
	defer DB.Close()

	r := mux.NewRouter()

	r.HandleFunc("/goals", GetGoals).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/goals", CreateGoal).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/health", HealthCheck).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/goals/export", ExportGoals).Methods(http.MethodGet, http.MethodOptions)

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

// Existing handlers (GetGoals, CreateGoal, HealthCheck)...

func ExportGoals(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, title, status FROM goals")
	if err != nil {
		http.Error(w, "Failed to fetch goals", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheetName := "Goals"
	f.SetSheetName("Sheet1", sheetName)

	f.SetCellValue(sheetName, "A1", "ID")
	f.SetCellValue(sheetName, "B1", "Title")
	f.SetCellValue(sheetName, "C1", "Status")

	row := 2
	for rows.Next() {
		var id int
		var title, status string
		if err := rows.Scan(&id, &title, &status); err != nil {
			http.Error(w, "Failed to scan goals", http.StatusInternalServerError)
			return
		}
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), id)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), title)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), status)
		row++
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=goals.xlsx")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := f.Write(w); err != nil {
		http.Error(w, "Failed to generate Excel", http.StatusInternalServerError)
	}
}
