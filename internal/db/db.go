package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // _ for indirectly using postgres with database/sql
)

func ConnectToPostgres() (*sql.DB, error) {
	// Fetch credentials from environment variables
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Check env variables loaded
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("missing database configuration")
	}

	// Format data into string
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
