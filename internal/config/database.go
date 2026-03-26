package config

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var lock = &sync.Mutex{}

// DatabaseConfig holds the database configuration parameters.
type DatabaseConnection struct {
	DB *sql.DB
}

var singleInstance *DatabaseConnection

func NewDatabaseConfig() *DatabaseConnection {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		db, _ := openDB()
		singleInstance = &DatabaseConnection{DB: db}
	}
	return singleInstance
}

// OpenDB opens a postgres connection using environment variables:
//
//	DB_HOST (default: localhost)
//	DB_PORT (default: 5432)
//	DB_USER
//	DB_PASSWORD
//	DB_NAME
//
// It returns a *sql.DB ready to use.
func openDBPostGres() (*sql.DB, error) {
	//host := getenv("DB_HOST", "host.docker.internal")
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "postgres")
	pass := getenv("DB_PASSWORD", "postgres")
	name := getenv("DB_NAME", "estoque")

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

func openDB() (*sql.DB, error) {
	//host := getenv("MYSQL_HOST", "host.docker.internal")
	host := getenv("MYSQL_HOST", "localhost")
	port := getenv("MYSQL_PORT", "3306")
	user := getenv("MYSQL_USER", "root")
	pass := getenv("MYSQL_ROOT_PASSWORD", "root")
	name := getenv("MYSQL_DATABASE_KENJIX", "estoque")

	if user == "" || name == "" {
		return nil, fmt.Errorf("MYSQL_PORT and MYSQL_DATABASE_KENJIX must be set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user,
		pass,
		host,
		port,
		name,
	)

	db, err := sql.Open("mysql", dsn)
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
