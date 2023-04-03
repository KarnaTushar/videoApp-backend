package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Database contains a pointer to the database object
type Database struct {
	*sqlx.DB
}

// CreateDB is used to initialize a new database connection
func CreateDB(dbURL string) (*Database, error) {
	db, err := sqlx.Connect("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("DB connection success")

	return &Database{db}, nil
}
