package placeholder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetAPIPosts retrieves posts from the API
func GetAPIPosts(url string) ([]APIPost, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Printf("Failed to ingest posts: %v\n", err)
		return nil, fmt.Errorf("failed to ingest posts: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Failed to read response body: %v\n", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var posts []APIPost
	err = json.Unmarshal(body, &posts)
	if err != nil {
		logger.Printf("Failed to unmarshal JSON: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	logger.Println("Successfully retrieved posts from API")
	return posts, nil
}

// IngestPosts inserts posts into the PostgreSQL database
func IngestPosts(DB *sql.DB, posts []APIPost) error {
	timestamp := time.Now().UTC()
	for _, post := range posts {
		_, err := DB.Exec(`
            INSERT INTO posts (api_id, ingest_timestamp, api_user_id, api_title, api_body) 
            VALUES ($1, $2, $3, $4, $5)`,
			post.ID, timestamp, post.UserID, post.Title, post.Body)
		if err != nil {
			logger.Printf("Failed to insert post: %v\n", err)
			IncrementInjestErrors()
			return fmt.Errorf("failed to insert post: %v", err)
		}
	}
	logger.Println("Data successfully stored in database.")
	IncrementPostsInjested(uint64(len(posts)))
	return nil
}
