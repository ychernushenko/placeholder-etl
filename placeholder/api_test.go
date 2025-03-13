package placeholder

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestMain runs before any tests are executed
func TestMain(m *testing.M) {
	// Initialize logger
	InitLogger()

	// Run the tests
	m.Run()
}

// TestGetAPIPosts tests the GetAPIPosts function
func TestGetAPIPosts(t *testing.T) {
	// Create a test server with a mock API response
	mockResponse := `[{"userId": 1, "id": 1, "title": "Test Title", "body": "Test Body"}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Call the GetAPIPosts function with the test server URL
	posts, err := GetAPIPosts(server.URL + "/posts")

	// Assert no error and correct response
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, 1, posts[0].UserID)
	assert.Equal(t, 1, posts[0].ID)
	assert.Equal(t, "Test Title", posts[0].Title)
	assert.Equal(t, "Test Body", posts[0].Body)
}

// TestIngestPosts tests the IngestPosts function
func TestIngestPosts(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Define the mock posts
	posts := []APIPost{
		{UserID: 1, ID: 1, Title: "Test Title", Body: "Test Body"},
	}

	// Define the expected query and result
	mock.ExpectExec("INSERT INTO posts").
		WithArgs(posts[0].ID, sqlmock.AnyArg(), posts[0].UserID, posts[0].Title, posts[0].Body).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the IngestPosts function
	err = IngestPosts(db, posts)

	// Assert no error and correct query execution
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
