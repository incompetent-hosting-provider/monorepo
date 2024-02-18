package backend

import (
	"cli/internal/authentication"
	"fmt"
	"net/http"
)

const (
	backendBaseURL = "http://localhost:8081"
)

var (
	ErrNotAuthenticated = fmt.Errorf("not authenticated")
)

func getAuthenticatedRequest(method string, path string, accessToken authentication.AccessToken) (*http.Request, error) {
	req, err := http.NewRequest("GET", backendBaseURL+"/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer " + string(accessToken))
	return req, nil
}
