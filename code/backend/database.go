package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error

	// Retry until DB is ready
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}
		log.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS goals (
		id SERIAL PRIMARY KEY,
		title TEXT,
		status TEXT
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}

	log.Println("Database connected and table ready.")
}
