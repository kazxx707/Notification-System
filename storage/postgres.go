package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DB wraps the database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
