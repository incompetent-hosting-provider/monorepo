package backend

import (
	"cli/internal/authentication"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


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
