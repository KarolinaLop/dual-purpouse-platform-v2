package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// DB stores and retrieves our model data.
var DB *sql.DB

func init() {
	var err error
	if DB, err = sql.Open("sqlite3", "./data/dp-platform.db"); err != nil {
		log.Fatal(err)
	}
}
