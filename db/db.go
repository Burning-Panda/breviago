package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
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
	UUID        uuid.UUID `json:"uuid"`
	ShortForm   string `json:"short_form"`
	LongForm    string `json:"long_form"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type AcronymCategory struct {
	ID          int    `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	AcronymID   int `json:"acronym_id"`
	CategoryID  int `json:"category_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
// 

type Category struct {
	ID          int    `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Label struct {
	ID          int    `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Label       string `json:"label"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// initSchema creates the necessary tables if they don't exist
func initSchema() {
	// Read the init.sql file
	initSQL, err := os.ReadFile("./db/init.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(string(initSQL))
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

// InsertAcronym inserts a single acronym into the database
func InsertAcronym(acronym *Acronym) error {
	database := GetDB()
	
	// Generate a new UUID if not provided
	if acronym.UUID == uuid.Nil {
		acronym.UUID = uuid.New()
	}
	
	res, err := database.Exec(
		"INSERT INTO acronyms (uuid, short_form, long_form, description) VALUES (?, ?, ?, ?)",
		acronym.UUID,
		acronym.ShortForm,
		acronym.LongForm,
		acronym.Description,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	return database.QueryRow(
		"SELECT id, uuid, short_form, long_form, description, created_at, updated_at FROM acronyms WHERE id = ?", 
		id,
	).Scan(
		&acronym.ID,
		&acronym.UUID,
		&acronym.ShortForm,
		&acronym.LongForm,
		&acronym.Description,
		&acronym.CreatedAt,
		&acronym.UpdatedAt,
	)
}

// InsertAcronyms inserts multiple acronyms into the database
func InsertAcronyms(acronyms []Acronym) ([]Acronym, error) {
	database := GetDB()
	
	// Start a transaction
	tx, err := database.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Prepare the insert statement
	stmt, err := tx.Prepare(
		"INSERT INTO acronyms (uuid, short_form, long_form, description) VALUES (?, ?, ?, ?)",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Insert each acronym
	for i := range acronyms {
		// Generate a new UUID if not provided
		if acronyms[i].UUID == uuid.Nil {
			acronyms[i].UUID = uuid.New()
		}

		res, err := stmt.Exec(
			acronyms[i].UUID,
			acronyms[i].ShortForm,
			acronyms[i].LongForm,
			acronyms[i].Description,
		)
		if err != nil {
			return nil, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		// Update the acronym with the new ID and timestamps
		err = tx.QueryRow(
			"SELECT id, uuid, short_form, long_form, description, created_at, updated_at FROM acronyms WHERE id = ?",
			id,
		).Scan(
			&acronyms[i].ID,
			&acronyms[i].UUID,
			&acronyms[i].ShortForm,
			&acronyms[i].LongForm,
			&acronyms[i].Description,
			&acronyms[i].CreatedAt,
			&acronyms[i].UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return acronyms, nil
} 

func SearchAcronyms(query string) ([]Acronym, error) {
	database := GetDB()

	rows, err := database.Query(
		"SELECT id, uuid, short_form, long_form, description, created_at, updated_at FROM acronyms WHERE short_form LIKE ? OR long_form LIKE ? OR description LIKE ?",
		fmt.Sprintf("%%%s%%", query),
		fmt.Sprintf("%%%s%%", query),
		fmt.Sprintf("%%%s%%", query),
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acronyms []Acronym
	for rows.Next() {
		var acronym Acronym
		err := rows.Scan(
			&acronym.ID,
			&acronym.UUID,
			&acronym.ShortForm,
			&acronym.LongForm,
			&acronym.Description,
			&acronym.CreatedAt,
			&acronym.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		acronyms = append(acronyms, acronym)
	}

	return acronyms, nil
}
