package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	database *sql.DB
	once     sync.Once
)

// GetDB returns the database instance, creating it if necessary
func GetDB() *sql.DB {
	once.Do(func() {
		db, err := sql.Open("sqlite3", "./acronyms.db")
		if err != nil {
			log.Fatal(err)
		}

		// Test the connection
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		database = db

		// Initialize the database schema
		initSchema()
	})

	return database
}

// initSchema creates the necessary tables if they don't exist
func initSchema() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS acronyms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_form TEXT NOT NULL UNIQUE,
		long_form TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := database.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Is WAL mode enabled?
	walMode := database.QueryRow("PRAGMA journal_mode")
	var mode string
	err = walMode.Scan(&mode)
	if err != nil {
		log.Fatal(err)
	}
	
	if mode != "WAL" {
		log.Println("WAL mode is not enabled, enabling it")
		database.Exec("PRAGMA journal_mode = WAL")
	}
}

// CloseDB closes the database connection
func CloseDB() error {
	if database != nil {
		return database.Close()
	}
	return nil
} 