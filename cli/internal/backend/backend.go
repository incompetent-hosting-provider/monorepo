package backend

import (
	"cli/internal/authentication"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrNotAuthenticated = fmt.Errorf("not authenticated")
)

type BackendClient struct {
	baseURL string
	client *http.Client
}

var DefaultBackendClient = &BackendClient{
	baseURL: "http://localhost:8081",
	client: &http.Client{},
}

func (client BackendClient) buildAuthenticatedRequest(method string, path string, accessToken authentication.AccessToken, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, client.baseURL + path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer " + string(accessToken))
	return req, nil
}
