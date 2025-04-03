package models

import (
	"database/sql"
)

type Table interface {
	InitTable(db *sql.DB) error
}

var db *sql.DB

// SetDB sets the database connection for the models package
func SetDB(database *sql.DB) {
	db = database
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// InitDB initializes all database tables required by the models
func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./models.db")
	if err != nil {
		return err
	}

	tables := []Table{
		&User{},
		&RemoteComputer{},
		&Credential{},
		&ComputerCredentialMapping{},
	}
	for _, table := range tables {
		err = table.InitTable(db)
		if err != nil {
			return err
		}
	}
	return nil
}
