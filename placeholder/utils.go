package placeholder

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

var (
	logger             *log.Logger
	injestErrorsTotal  uint64
	postsInjestedTotal uint64
	postsETLTotal      uint64
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

func InitLogger() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		os.Exit(1)
	}

	// Open log file
	logFile, err := os.OpenFile("logs/etl.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}

	// Create logger
	logger = log.New(logFile, "", log.LstdFlags)
	logger.Println("Logger initialized successfully")
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// MetricsHandler handles metrics requests (Prometheus style)
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := fmt.Sprintf(`
# HELP injest_post_total The total number posts ingested to DB.
# TYPE injest_post_total counter
injest_post_total %d

# HELP injest_errors_total The total number of ingest errors.
# TYPE injest_errors_total counter
injest_errors_total %d

# HELP etl_post_total The total number of ETL posts.
# TYPE etl_post_total counter
etl_post_total %d
`, atomic.LoadUint64(&postsInjestedTotal), atomic.LoadUint64(&injestErrorsTotal), atomic.LoadUint64(&postsETLTotal))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
}

// IncrementInjectErrors increments the ingest errors counter
func IncrementInjestErrors() {
	atomic.AddUint64(&injestErrorsTotal, 1)
}

// IncrementPostsInjested increments the posts injested to database counter
func IncrementPostsInjested(numOfPosts uint64) {
	atomic.AddUint64(&postsInjestedTotal, numOfPosts)
}

// IncrementETLPosts increments the posts processed with ETL to database counter
func IncrementETLPosts(numOfPosts uint64) {
	atomic.AddUint64(&postsETLTotal, numOfPosts)
}
