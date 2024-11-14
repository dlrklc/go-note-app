package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func Init() {
	var err error
	connStr := "user=postgres password=postgres dbname=note_db sslmode=disable" //todo: make private
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Check if the database is reachable
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	log.Println("Database connected successfully!")
}

// CloseDB closes the database connection
func Close() {
	if err := DB.Close(); err != nil {
		log.Fatal("Error closing the database connection: ", err)
	}
}
