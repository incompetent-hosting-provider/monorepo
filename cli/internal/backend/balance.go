package backend

import (
	"cli/internal/authentication"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)


type GetBalanceResponse struct {
	Balance int `json:"balance,omitempty"`
}

// Gets the current balance of the user associated to the tokens provided.
func GetBalance(tokens authentication.SessionTokens) (int, error) {
	balance, err := getBalance(tokens.AccessToken)
	if err != nil && errors.Is(err, ErrNotAuthenticated) {
		newTokens, err := authentication.RefreshTokens()
		if err != nil || newTokens == nil {
			return 0, fmt.Errorf("failed to get user balance: %w: %w", ErrNotAuthenticated, err)
		}

		return getBalance(newTokens.AccessToken)
	}

	return balance, err
}

func getBalance(token authentication.AccessToken) (int, error) {
	req, err := getAuthenticatedRequest("GET", "/payment", token, nil)
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
