package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func GetDBConnection() (*sql.DB, error) {
	// Replace with your own connection details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_AURORA_USERNAME"),
		os.Getenv("DB_AURORA_PASSWORD"),
		os.Getenv("DB_AURORA_HOST"),
		os.Getenv("DB_AURORA_PORT"),
		os.Getenv("DB_AURORA_NAME"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error opening database connection:", err)
		return nil, err
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return nil, err
	}

	fmt.Println("Successfully connected to Aurora database!")
	return db, nil
}
