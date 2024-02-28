package backend

import (
	"cli/internal/authentication"
	"testing"
)

// TestGetAuthenticatedRequest tests that
// getAuthenticatedRequest returns a request with the correct URL.
func TestGetAuthenticatedRequestURL(t *testing.T) {
	// ARRANGE
	accessToken := authentication.AccessToken("dummyToken")
	testReqPath := "/test"

	// ACT

	req, err := DefaultBackendClient.buildAuthenticatedRequest("GET", testReqPath, accessToken, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// ASSERT
	expectedURL := DefaultBackendClient.baseURL + testReqPath
	if req.URL.String() != expectedURL {
		t.Errorf("unexpected URL, got: %s, want: %s", req.URL.String(), expectedURL)
	}
}

// TestGetAuthenticatedRequest tests that
// getAuthenticatedRequest returns a request with the correct Headers.
func TestGetAuthenticatedRequestHeader(t *testing.T) {
	// ARRANGE
	accessToken := authentication.AccessToken("dummyToken")
	reqPath := "/test"

	// ACT
	req, err := DefaultBackendClient.buildAuthenticatedRequest("GET", reqPath, accessToken, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// ASSERT
	expectedAuthHeader := "Bearer " + string(accessToken)
	if req.Header.Get("Authorization") != expectedAuthHeader {
		t.Errorf("unexpected Authorization header, got: %s, want: %s", req.Header.Get("Authorization"), expectedAuthHeader)
	}
}