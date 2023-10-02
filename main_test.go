package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// This should be added to your main_test.go
func TestMain(m *testing.M) {
	// Start your server here
	go main()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Run the tests
	os.Exit(m.Run())
}

func TestKeyGeneration(t *testing.T) {
	err := initKeys()
	assert.NoError(t, err, "Key generation should not return an error")
}

func TestJWKSHandler(t *testing.T) {
	// Create a request to the JWKS endpoint
	req, err := http.NewRequest("GET", "/jwks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a test router and serve the request
	r := setupRouter()
	r.ServeHTTP(rr, req)

	// Check the status code and response body
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
}

func TestAuthHandlerValidCredentials(t *testing.T) {
	// Create a request to the auth endpoint with valid credentials
	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = url.Values{
		"username": {"userABC"},
		"password": {"password123"},
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a test router and serve the request
	r := setupRouter()
	r.ServeHTTP(rr, req)

	// Check the status code and response body
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	// Add additional assertions as needed to check the JWT response.
}

func TestAuthHandlerInvalidCredentials(t *testing.T) {
	// Create a request to the auth endpoint with invalid credentials
	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = url.Values{
		"username": {"invalidUser"},
		"password": {"invalidPassword"},
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a test router and serve the request
	r := setupRouter()
	r.ServeHTTP(rr, req)

	// Check the status code (should be Unauthorized)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
	// Add additional assertions as needed for error response.
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/jwks", JWKSHandler)
	// Register other routes if needed
	return r
}

func TestServerIsUp(t *testing.T) {
	// Define the base URL of your server
	baseURL := "http://localhost:8080"

	// Send a GET request to the base URL
	resp, err := http.Get(baseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")
}
