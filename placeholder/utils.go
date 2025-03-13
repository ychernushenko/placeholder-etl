package placeholder

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Check if environment variables are set
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatalf("Missing required environment variables: DB_HOST=%s DB_USER=%s DB_PASSWORD=%s DB_NAME=%s", dbHost, dbUser, dbPassword, dbName)
	}

	// PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	fmt.Println("Connecting with DSN:", dsn) // Log DSN for debugging

	var err error
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Retry connection logic
	for i := 0; i < 5; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		if i == 4 {
			log.Fatalf("Database not reachable after 5 attempts: %v", err)
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	fmt.Println("Connected to PostgreSQL database!")

	// Ensure table exists
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS posts (
        ingest_id SERIAL PRIMARY KEY,
		api_id INT NOT NULL,
		ingest_timestamp TIMESTAMPTZ NOT NULL,
        api_user_id INT NOT NULL,
        api_title TEXT NOT NULL,
        api_body TEXT NOT NULL
    );`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Database table ensured.")

	return DB
}
