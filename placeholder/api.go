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
func GetAPIPosts() ([]APIPost, error) {
	base_url := "https://jsonplaceholder.typicode.com/"
	endpoint_suffix := "posts"
	url := base_url + endpoint_suffix

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to ingest posts: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var posts []APIPost
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

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
			return fmt.Errorf("failed to insert post: %v", err)
		}
	}
	fmt.Println("Data successfully stored in database.")
	return nil
}
