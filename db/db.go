package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	return sql.Open("sqlite3", "./storage/db/cooper.db")
}
