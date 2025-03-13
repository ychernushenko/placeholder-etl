package placeholder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ExtractPosts retrieves posts from the database starting from a specific ID and saves them to disk
func ExtractPosts(DB *sql.DB, startID, limit int, prefix string) (int, error) {
	query := `
        SELECT ingest_id, api_id, ingest_timestamp, api_user_id, api_title, api_body 
        FROM posts 
        WHERE ingest_id >= $1 
        ORDER BY ingest_id 
        LIMIT $2`
	rows, err := DB.Query(query, startID, limit)
	if err != nil {
		logger.Printf("Failed to query posts: %v\n", err)
		return 0, fmt.Errorf("failed to query posts: %v", err)
	}
	defer rows.Close()

	var posts []RawPost
	var latestID int
	for rows.Next() {
		var post RawPost
		err := rows.Scan(&post.IngestID, &post.APIID, &post.IngestTimestamp, &post.APIUserID, &post.APITitle, &post.APIBody)
		if err != nil {
			logger.Printf("Failed to scan post: %v\n", err)
			return 0, fmt.Errorf("failed to scan post: %v", err)
		}
		posts = append(posts, post)
		latestID = post.IngestID
	}

	if len(posts) == 0 {
		logger.Println("No new posts to extract.")
		return 0, nil
	}

	err = saveToFile(posts, prefix, "raw_posts")
	if err != nil {
		logger.Printf("Failed to save posts to file: %v\n", err)
		return 0, fmt.Errorf("failed to save posts to file: %v", err)
	}

	logger.Printf("Successfully extracted and saved %d posts.\n", len(posts))
	return latestID, nil
}

// TransformPosts reads raw posts from disk, transforms them, and saves the transformed posts back to disk
func TransformPosts(filename string, prefix string) ([]ProcessedPost, error) {
	file, err := os.Open(filename)
	if err != nil {
		logger.Printf("Failed to open file: %v\n", err)
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Printf("Failed to read raw posts from file: %v\n", err)
		return nil, fmt.Errorf("failed to read raw posts from file: %v", err)
	}

	var rawPosts []RawPost
	err = json.Unmarshal(data, &rawPosts)
	if err != nil {
		logger.Printf("Failed to unmarshal raw posts: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal raw posts: %v", err)
	}

	var processedPosts []ProcessedPost
	for _, rawPost := range rawPosts {
		processedPost := ProcessedPost{
			UserID:          rawPost.APIUserID,
			ID:              rawPost.APIID,
			Title:           rawPost.APITitle,
			Body:            rawPost.APIBody,
			IngestTimestamp: rawPost.IngestTimestamp,
		}
		processedPosts = append(processedPosts, processedPost)
	}

	err = saveToFile(processedPosts, prefix, "processed_posts")
	if err != nil {
		logger.Printf("Failed to save processed posts to file: %v\n", err)
		return nil, fmt.Errorf("failed to save processed posts to file: %v", err)
	}

	IncrementETLPosts(uint64(len(processedPosts)))
	logger.Printf("Successfully transformed and saved %d posts.\n", len(processedPosts))

	return processedPosts, nil
}

// saveToFile saves data to a file with a specified prefix and type
func saveToFile(data interface{}, prefix, fileType string) error {
	// Ensure the directory exists
	err := os.MkdirAll(prefix, os.ModePerm)
	if err != nil {
		logger.Printf("Failed to create directory: %v\n", err)
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := filepath.Join(prefix, fmt.Sprintf("%s_%d.json", fileType, time.Now().Unix()))
	file, err := os.Create(filename)
	if err != nil {
		logger.Printf("Failed to create file: %v\n", err)
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Printf("Failed to marshal data: %v\n", err)
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		logger.Printf("Failed to write data to file: %v\n", err)
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}
