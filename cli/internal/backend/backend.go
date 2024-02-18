package backend

import (
	"cli/internal/authentication"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	backendBaseURL = "http://localhost:8081"
)

var (
	ErrNotAuthenticated = fmt.Errorf("not authenticated")
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

type GetBalanceResponse struct {
	Balance int `json:"balance,omitempty"`
}

func GetBalance(tokens authentication.SessionTokens) (int, error) {
	req, err := getAuthenticatedRequest("GET", "/balance", tokens.AccessToken)
	if err != nil {
		return 0, fmt.Errorf("failed to create get balance request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get user balance: %w", err)
	}

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		return 0, fmt.Errorf("failed to get user balance: %w", ErrNotAuthenticated)
	default:
		return 0, fmt.Errorf("failed to get user balance: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to read response body of user balance request: %w", err)
	}

	var balanceResponse *GetBalanceResponse
	err = json.Unmarshal(body, &balanceResponse)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response body of user balance request: %w", err)
	}

	return balanceResponse.Balance, nil
}

func getAuthenticatedRequest(method string, path string, accessToken authentication.AccessToken) (*http.Request, error) {
	req, err := http.NewRequest("GET", backendBaseURL+"/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer " + string(accessToken))
	return req, nil
}
