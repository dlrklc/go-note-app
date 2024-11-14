package db

import (
	"database/sql"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func Init() {
	var err error

	config, err := loadConfig("db/config.txt")

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	connStr := config
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

func loadConfig(filePath string) (string, error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer configFile.Close()

	configData, err := io.ReadAll(configFile)
	if err != nil {
		return "", err
	}

	return string(configData), nil
}
