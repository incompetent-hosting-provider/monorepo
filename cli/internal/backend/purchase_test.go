package backend

import (
	"cli/internal/authentication"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestMakePurchaseRequest_RequestParams(t *testing.T) {
	// ARRANGE
	testToken := authentication.AccessToken("dummy-token")
	testPurchaseAmount := 10

	// Create a mock server
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/payment", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "Bearer "+ string(testToken), r.Header.Get("Authorization"))

			reqBody, err := io.ReadAll(r.Body);
			assert.NoError(t, err)
			
			var reqBodyMap map[string]any
			err = json.Unmarshal(reqBody, &reqBodyMap)
			assert.NoError(t, err)
			assert.EqualValues(t, testPurchaseAmount, reqBodyMap["amount"])

			w.WriteHeader(200)
		}),
	)
	defer mockServer.Close()

	testClient := BackendClient{
		baseURL: mockServer.URL,
		client: mockServer.Client(),
	}

	// ACT
	_, _ = testClient.PurchaseCredits(testToken, testPurchaseAmount, true)
}

func TestMakePurchaseRequest_Success(t *testing.T) {
	// ARRANGE
	testToken := authentication.AccessToken("valid-token")
	testPurchaseAmount := 10
	expectedNewBalance := 100

	// Create a mock server that returns an 200 response
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"balance": %d}`, expectedNewBalance)
		}),
	)
	defer mockServer.Close()
	
	testClient := BackendClient{
		baseURL: mockServer.URL,
		client: mockServer.Client(),
	}


	// ACT
	balance, err := testClient.PurchaseCredits(testToken, testPurchaseAmount, true)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, 100, balance)
}

func TestMakePurchaseRequest_InternalServerError(t *testing.T) {
	// ARRANGE
	testPurchaseAmount := 10
	testToken := authentication.AccessToken("valid-token")

	// Create a mock server that returns an 500 response
	mockErrorResponseCode := http.StatusInternalServerError
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(mockErrorResponseCode)
			fmt.Fprint(w, "internal server error")
		}),
	)
	defer mockServer.Close()
	
	testClient := BackendClient{
		baseURL: mockServer.URL,
		client: mockServer.Client(),
	}

	// ACT
	balance, err := testClient.PurchaseCredits(testToken, testPurchaseAmount, true)
	
	// ASSERT
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), strconv.Itoa(mockErrorResponseCode))
	}

	assert.Zero(t, balance)
}

func TestMakePurchaseRequest_Unauthenticated(t *testing.T) {
	// ARRANGE
	testPurchaseAmount := 10
	testToken := authentication.AccessToken("invalid-token")

	// Create a mock server that returns an 500 response
	mockErrorResponseCode := http.StatusUnauthorized
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(mockErrorResponseCode)
			fmt.Fprint(w, "unauthorized")
		}),
	)
	defer mockServer.Close()

	testClient := BackendClient{
		baseURL: mockServer.URL,
		client: mockServer.Client(),
	}

	// ACT
	balance, err := testClient.PurchaseCredits(testToken, testPurchaseAmount, true)
	
	// ASSERT
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNotAuthenticated)
	}

	
	assert.Zero(t, balance)
}