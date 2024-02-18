package backend

import (
	"cli/internal/authentication"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserInfo struct {
	Balance int    `json:"balance,omitempty"`
	Email   string `json:"email,omitempty"`
}

func GetUserInfo(tokens authentication.SessionTokens) (*UserInfo, error) {
	req, err := getAuthenticatedRequest("GET", "/user", tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to build request for user info: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		return nil, fmt.Errorf("failed to get user info: %w", ErrNotAuthenticated)
	default:
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body of user info request: %w", err)
	}

	var userInfo *UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body of user info request: %w", err)
	}

	return userInfo, nil
}
