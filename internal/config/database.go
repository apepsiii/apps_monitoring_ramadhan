package config

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "./amaliah.db"
	}

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}
