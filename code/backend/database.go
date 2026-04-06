package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Database struct {
		Type       string `json:"type"`
		SQLitePath string `json:"sqlite_path"`
	} `json:"database"`
}

var DB *sql.DB

func InitDB() {
	// Load config
	configFile := filepath.Join(".", "config.json")
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Failed to open config.json:", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal("Failed to parse config.json:", err)
	}

	switch config.Database.Type {
	case "postgres":
		// Use PostgreSQL (existing logic)
		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)

		for i := 0; i < 10; i++ {
			DB, err = sql.Open("postgres", connStr)
			if err == nil {
				err = DB.Ping()
				if err == nil {
					break
				}
			}
			log.Println("Waiting for PostgreSQL...")
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			log.Fatal("PostgreSQL connection failed:", err)
		}

	case "sqlite":
		// Use SQLite with configurable path
		dbPath := config.Database.SQLitePath
		if !filepath.IsAbs(dbPath) {
			// Make relative paths relative to backend folder
			dbPath = filepath.Join(".", dbPath)
		}
		// Ensure directory exists
		os.MkdirAll(filepath.Dir(dbPath), 0755)

		DB, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("SQLite connection failed:", err)
		}
		if err = DB.Ping(); err != nil {
			log.Fatal("SQLite ping failed:", err)
		}

	default:
		log.Fatal("Unsupported database type:", config.Database.Type)
	}

	// Create table (works for both DBs with minor adjustments)
	createTable := `
    CREATE TABLE IF NOT EXISTS goals (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        status TEXT DEFAULT 'active'
    );`
	if config.Database.Type == "sqlite" {
		createTable = `
        CREATE TABLE IF NOT EXISTS goals (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            status TEXT DEFAULT 'active'
        );`
	}

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}

	log.Printf("Database connected (%s) and table ready.", config.Database.Type)
}
