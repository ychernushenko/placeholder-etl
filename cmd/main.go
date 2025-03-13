package main

import (
	"net/http"
	"placeholder-etl/placeholder"
)

func main() {
	// Initialize logger
	placeholder.InitLogger()

	// Initialize database
	DB := placeholder.InitDB()
	defer DB.Close()

	// Run IngestAPIData and ETLtoDatalake in parallel as goroutines
	go placeholder.IngestAPIData(DB, 30)
	go placeholder.ETLtoDatalake(DB, 10, 0, 50, "./data/raw/", "./data/processed/")

	// Expose health and metrics endpoints
	http.HandleFunc("/health", placeholder.HealthCheckHandler)
	http.HandleFunc("/metrics", placeholder.MetricsHandler)

	// Start HTTP server
	http.ListenAndServe(":8080", nil)
}
