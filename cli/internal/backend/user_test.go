package backend

import (
	"cli/internal/authentication"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserInfo_RequestParams(t *testing.T) {
	// ARRANGE
	tokens := authentication.SessionTokens{
		AccessToken: authentication.AccessToken("accessToken"),
		RefreshToken: authentication.RefreshToken("refreshToken"),
	}

	// Create a test server to mock the HTTP response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ASSERT
		assert.Equal(t, "/user", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer "+ string(tokens.AccessToken), r.Header.Get("Authorization"))

		// RESPOND
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	testClient := BackendClient{
		baseURL: server.URL,
		client: server.Client(),
	}

	// ACT
	_, _ =  testClient.GetUserInfo(tokens)
}

func TestGetUserInfo_Success(t *testing.T) {
	// ARRANGE
	tokens := authentication.SessionTokens{
		AccessToken: authentication.AccessToken("accessToken"),
		RefreshToken: authentication.RefreshToken("refreshToken"),
	}

	expectedUserInfo := UserInfo{
		Balance: 100,
		Email: "john.doe@foobar.de",
	}

	// Create a test server to mock the HTTP response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ASSERT
		assert.Equal(t, "/user", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer "+ string(tokens.AccessToken), r.Header.Get("Authorization"))

		// RESPOND
		w.WriteHeader(http.StatusOK)
		responseJson, _ := json.Marshal(expectedUserInfo)
		w.Write(responseJson)
	}))
	defer server.Close()

	testClient := BackendClient{
		baseURL: server.URL,
		client: server.Client(),
	}

	// ACT
	userInfo, err :=  testClient.GetUserInfo(tokens)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, &expectedUserInfo, userInfo)
}

func TestGetUserInfo_InternalServerError(t *testing.T) {
		// ARRANGE
		tokens := authentication.SessionTokens{
			AccessToken: authentication.AccessToken("accessToken"),
			RefreshToken: authentication.RefreshToken("refreshToken"),
		}

		mockErrorResponseCode := http.StatusInternalServerError
	
		// Create a test server to mock the HTTP response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(mockErrorResponseCode)
		}))
		defer server.Close()

		testClient := BackendClient{
			baseURL: server.URL,
			client: server.Client(),
		}

		// ACT
		userInfo, err :=  testClient.GetUserInfo(tokens)

		// ASSERT
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), fmt.Sprint(mockErrorResponseCode))
		}
	
		assert.Zero(t, userInfo)
}


func TestGetUserInfo_Unauthenticated(t *testing.T) {
	// ARRANGE
	tokens := authentication.SessionTokens{
		AccessToken: authentication.AccessToken("accessToken"),
		RefreshToken: authentication.RefreshToken("refreshToken"),
	}

	mockErrorResponseCode := http.StatusUnauthorized

	// Create a test server to mock the HTTP response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(mockErrorResponseCode)
	}))
	defer server.Close()

	testClient := BackendClient{
		baseURL: server.URL,
		client: server.Client(),
	}

	// ACT
	userInfo, err := testClient.GetUserInfo(tokens)

	// ASSERT
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNotAuthenticated)
	}

	assert.Zero(t, userInfo)
}