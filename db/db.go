package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)

	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln("DB connection error:", err)
	}

	if err := runMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func runMigrations() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sequences (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		open_tracking_enabled BOOLEAN NOT NULL DEFAULT FALSE,
		click_tracking_enabled BOOLEAN NOT NULL DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS sequence_steps (
		id SERIAL PRIMARY KEY,
		sequence_id INTEGER NOT NULL REFERENCES sequences(id) ON DELETE CASCADE,
		subject TEXT NOT NULL,
		content TEXT NOT NULL,
		waiting_days INTEGER NOT NULL DEFAULT 0
	);`

	_, err := DB.Exec(schema)
	return err
}
