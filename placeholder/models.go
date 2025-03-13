package placeholder

import "time"

// Post represents the structure of a post from the API
type APIPost struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// RawPost represents the structure of a post read from the PostgreSQL database
type RawPost struct {
	IngestID        int       `json:"ingest_id"`
	APIID           int       `json:"api_id"`
	IngestTimestamp time.Time `json:"ingest_timestamp"`
	APIUserID       int       `json:"api_user_id"`
	APITitle        string    `json:"api_title"`
	APIBody         string    `json:"api_body"`
}

// ProcessedPost represents the structure of a processed post
type ProcessedPost struct {
	UserID          int       `json:"user_id"`
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Body            string    `json:"body"`
	IngestTimestamp time.Time `json:"ingest_timestamp"`
}
