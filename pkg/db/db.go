package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	HOME, _ = os.UserHomeDir()
	DB_PATH = HOME + "/.askGPT.db"
)

type MyDB struct {
	DB *sql.DB
}

// Init initializes the database connection. Returns a pointer to the database object and an error.
func Init() (sdb *MyDB, err error) {
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return nil, err
	}

	//Create a new SDB object with the opened database
	mydb := &MyDB{
		DB: db,
	}

	// Create the table if it doesn't exist
	_, err = mydb.DB.Exec(`CREATE TABLE IF NOT EXISTS responses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		id_gpt TEXT NOT NULL,
		input TEXT NOT NULL,
		answer TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return nil, err
	}

	return mydb, nil
}

// SaveStatement saves the response from the GPT model. Receives the ID, the input and the answer.
// Returns an error.
func (m *MyDB) SaveStatement(id, input, ans string) (err error) {
	// Save the response with ID as identifier in the database
	stmt, err := m.DB.Prepare(`INSERT INTO responses(
		id_gpt,
		input,
		answer) VALUES(?, ?, ?);`)
	if err != nil {
		e := fmt.Sprintf("Error preparing the statement: '%s' '%s' '%s'", id, input, ans)

		return errors.New(e)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, input, ans)
	if err != nil {
		e := fmt.Sprintf("Error executing the statement: '%s' '%s' '%s'", id, input, ans)

		return errors.New(e)
	}

	return nil
}
