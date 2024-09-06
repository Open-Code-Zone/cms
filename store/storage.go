package store

import (
	"database/sql"
	"log"
)

// NewSQLiteStorage opens a connection to the SQLite database
func NewSQLiteStorage(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
