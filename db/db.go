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

type Acronym struct {
	ID          int    `json:"id"`
	ShortForm   string `json:"short_form"`
	LongForm    string `json:"long_form"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Label struct {
	ID          int    `json:"id"`
	Label       string `json:"label"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// initSchema creates the necessary tables if they don't exist
func initSchema() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS acronyms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_form TEXT NOT NULL,
		long_form TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		
		UNIQUE(short_form, long_form)
	);
	CREATE TABLE IF NOT EXISTS acronym_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		acronym_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

		FOREIGN KEY (acronym_id) REFERENCES acronyms(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS note_joiner (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		acronym_id INTEGER,
		note_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		
		FOREIGN KEY (acronym_id) REFERENCES acronyms(id) ON DELETE SET NULL,
		FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		note TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

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