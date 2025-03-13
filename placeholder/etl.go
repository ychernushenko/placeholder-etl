package placeholder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
		return 0, fmt.Errorf("failed to query posts: %v", err)
	}
	defer rows.Close()

	var posts []RawPost
	var latestID int
	for rows.Next() {
		var post RawPost
		err := rows.Scan(&post.IngestID, &post.APIID, &post.IngestTimestamp, &post.APIUserID, &post.APITitle, &post.APIBody)
		if err != nil {
			return 0, fmt.Errorf("failed to scan post: %v", err)
		}
		posts = append(posts, post)
		latestID = post.IngestID
	}

	if len(posts) == 0 {
		return 0, nil
	}

	data, err := json.Marshal(posts)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal posts: %v", err)
	}

	filename := fmt.Sprintf("%s_raw_posts_%d.json", prefix, time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return 0, fmt.Errorf("failed to write posts to file: %v", err)
	}

	return latestID, nil
}

// TransformPosts reads raw posts from disk, transforms them, and saves the transformed posts back to disk
func TransformPosts(filename string, prefix string) ([]ProcessedPost, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read raw posts from file: %v", err)
	}

	var rawPosts []RawPost
	err = json.Unmarshal(data, &rawPosts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw posts: %v", err)
	}

	var processedPosts []ProcessedPost
	for _, rawPost := range rawPosts {
		processedPost := ProcessedPost{
			UserID: rawPost.APIUserID,
			ID:     rawPost.APIID,
			Title:  rawPost.APITitle,
			Body:   rawPost.APIBody,
		}
		processedPosts = append(processedPosts, processedPost)
	}

	data, err = json.Marshal(processedPosts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal processed posts: %v", err)
	}

	outputFilename := fmt.Sprintf("%s_processed_posts_%d.json", prefix, time.Now().Unix())
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write processed posts to file: %v", err)
	}

	return processedPosts, nil
}
