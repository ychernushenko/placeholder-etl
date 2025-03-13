package placeholder

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestExtractPosts tests the ExtractPosts function
func TestExtractPosts(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Define the mock rows
	rows := sqlmock.NewRows([]string{"ingest_id", "api_id", "ingest_timestamp", "api_user_id", "api_title", "api_body"}).
		AddRow(1, 1, time.Date(2025, time.March, 13, 9, 53, 56, 0, time.UTC), 1, "Test Title", "Test Body")

	// Define the expected query and result
	mock.ExpectQuery("SELECT ingest_id, api_id, ingest_timestamp, api_user_id, api_title, api_body FROM posts WHERE ingest_id >= \\$1 ORDER BY ingest_id LIMIT \\$2").
		WithArgs(0, 10).
		WillReturnRows(rows)

	// Ensure the directory exists
	err = os.MkdirAll("./test_data/raw", os.ModePerm)
	assert.NoError(t, err)

	// Call the ExtractPosts function
	latestID, err := ExtractPosts(db, 0, 10, "./test_data/raw")

	// Assert no error and correct response
	assert.NoError(t, err)
	assert.Equal(t, 1, latestID)

	// Verify that the file was created
	files, err := os.ReadDir("./test_data/raw")
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	// Clean up
	os.RemoveAll("./test_data/raw")
}

// TestTransformPosts tests the TransformPosts function
func TestTransformPosts(t *testing.T) {
	// Create a test file with raw posts
	rawPosts := []RawPost{
		{IngestID: 1, APIID: 1, IngestTimestamp: time.Date(2025, time.March, 13, 9, 53, 56, 0, time.UTC), APIUserID: 1, APITitle: "Test Title", APIBody: "Test Body"},
	}
	data, err := json.Marshal(rawPosts)
	assert.NoError(t, err)

	err = os.MkdirAll("./test_data/raw", os.ModePerm)
	assert.NoError(t, err)

	err = os.WriteFile("./test_data/raw/raw_posts_test.json", data, 0644)
	assert.NoError(t, err)

	// Call the TransformPosts function
	processedPosts, err := TransformPosts("./test_data/raw/raw_posts_test.json", "./test_data/processed")

	// Assert no error and correct response
	assert.NoError(t, err)
	assert.Len(t, processedPosts, 1)
	assert.Equal(t, "Test Title", processedPosts[0].Title)

	// Verify that the file was created
	files, err := os.ReadDir("./test_data/processed")
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	// Clean up
	os.RemoveAll("./test_data/raw")
	os.RemoveAll("./test_data/processed")
}
