package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// OpenDB opens a postgres connection using environment variables:
//
//	DB_HOST (default: localhost)
//	DB_PORT (default: 5432)
//	DB_USER
//	DB_PASSWORD
//	DB_NAME
//
// It returns a *sql.DB ready to use.
func OpenDB() (*sql.DB, error) {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || name == "" {
		return nil, fmt.Errorf("DB_USER, DB_PASSWORD and DB_NAME must be set")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
