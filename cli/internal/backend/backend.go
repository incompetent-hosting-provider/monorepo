package backend

import (
	"cli/internal/authentication"
	"fmt"
	"io"
	"net/http"
)

// ERRORS
var (
	ErrNotAuthenticated = fmt.Errorf("not authenticated")
)

var baseURL = "http://localhost:8081"

func getAuthenticatedRequest(method string, path string, accessToken authentication.AccessToken, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, baseURL + path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer " + string(accessToken))
	return req, nil
}
