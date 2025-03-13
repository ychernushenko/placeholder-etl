package main

import "placeholder-etl/placeholder"

func main() {
	// Initialize database
	DB := placeholder.InitDB()
	defer DB.Close()

	// Run IngestAPIData and ETLtoDatalake in parallel as goroutines
	go placeholder.IngestAPIData(DB, 30)
	go placeholder.ETLtoDatalake(DB, 10, 0, 50, "./data/raw/", "./data/processed/")

	// Prevent the main function from exiting immediately
	select {}
}
