package backend

import (
	"cli/internal/authentication"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGetBalance_RequestParams(t *testing.T) {
	// ARRANGE
	token := authentication.AccessToken("dummy-token")

	// Create a mock server
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/payment", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Bearer "+string(token), r.Header.Get("Authorization"))

			w.WriteHeader(200)
		}),
	)
	defer mockServer.Close()

	// Replace the default client with the mock server client
	http.DefaultClient = mockServer.Client()
	baseURL = mockServer.URL

	// ACT
	_, _ = getBalance(token)
}

func TestGetBalance_Success(t *testing.T) {
	// ARRANGE
	token := authentication.AccessToken("dummy-token")
	expectedBalance := 100

	// Create a mock server
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"balance": %d}`, expectedBalance)
		}),
	)
	defer mockServer.Close()
	
	// Replace the default client with the mock server client
	http.DefaultClient = mockServer.Client()
	baseURL = mockServer.URL

	// ACT
	balance, err := getBalance(token)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}
func TestGetBalance_InternalServerError(t *testing.T) {
	// ARRANGE
	token := authentication.AccessToken("dummy-token")
	mockErrorResponseCode := http.StatusInternalServerError

	// Create a mock server
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(mockErrorResponseCode)
			fmt.Fprint(w, "internal server error")
		}),
	)
	defer mockServer.Close()
	
	// Replace the default client with the mock server client
	http.DefaultClient = mockServer.Client()
	baseURL = mockServer.URL

	// ACT
	balance, err := getBalance(token)

	// ASSERT
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), strconv.Itoa(mockErrorResponseCode))
	}

	assert.Zero(t, balance)
}

func TestGetBalance_Unauthenticated(t *testing.T) {
	// ARRANGE
	token := authentication.AccessToken("dummy-token")
	mockErrorResponseCode := http.StatusUnauthorized

	// Create a mock server
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(mockErrorResponseCode)
			fmt.Fprint(w, "unauthorized")
		}),
	)
	defer mockServer.Close()
	
	// Replace the default client with the mock server client
	http.DefaultClient = mockServer.Client()
	baseURL = mockServer.URL

	// ACT
	balance, err := getBalance(token)

	// ASSERT
	if assert.Error(t, err) {
		assert.True(t, errors.Is(err, ErrNotAuthenticated))
	}

	assert.Zero(t, balance)
}



