package placeholder

import (
	"database/sql"
	"fmt"
	"time"
)

// IngestAPIData runs the ETL job every 30 seconds
func IngestAPIData(DB *sql.DB, secRate int) {
	ticker := time.NewTicker(time.Duration(secRate) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Starting ETL process...")

		// Extract
		posts, err := GetAPIPosts()
		if err != nil {
			fmt.Printf("Error during extraction: %v\n", err)
			continue
		}

		// Ingest into DB
		err = IngestPosts(DB, posts)
		if err != nil {
			fmt.Printf("Error during data ingest: %v\n", err)
			continue
		}

		fmt.Println("ETL process completed.")
	}
}

// ETLtoDatalake periodically reads data from the database, saves it to disk, transforms it, and saves the transformed data to disk
func ETLtoDatalake(DB *sql.DB, secRate int, startID, limit int, rawPrefix, processedPrefix string) {
	ticker := time.NewTicker(time.Duration(secRate) * time.Second)
	defer ticker.Stop()

	var latestID = startID

	for range ticker.C {
		fmt.Println("Starting ETL cycle...")

		// Extract posts from the database
		newLatestID, err := ExtractPosts(DB, latestID, limit, rawPrefix)
		if err != nil {
			fmt.Printf("Error during extraction: %v\n", err)
			continue
		}

		if newLatestID == 0 {
			fmt.Println("No new posts to extract.")
			continue
		}

		latestID = newLatestID

		// Transform the extracted posts
		rawFilename := fmt.Sprintf("%sraw_posts_%d.json", rawPrefix, time.Now().Unix())
		_, err = TransformPosts(rawFilename, processedPrefix)
		if err != nil {
			fmt.Printf("Error during transformation: %v\n", err)
			continue
		}

		fmt.Println("ETL cycle completed successfully.")
	}
}
