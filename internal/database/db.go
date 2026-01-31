package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	authToken := os.Getenv("DATABASE_AUTH_TOKEN")

	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	// Construct connection string with auth token
	connStr := dbURL
	if authToken != "" {
		connStr += "?authToken=" + authToken
	}

	var err error
	DB, err = sql.Open("libsql", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
