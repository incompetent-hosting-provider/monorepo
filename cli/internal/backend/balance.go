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

// Gets the current balance of the user associated to the tokens provided.
func (client BackendClient) GetBalance(token authentication.AccessToken, retryAuth bool) (int, error) {
	req, err := client.buildAuthenticatedRequest("GET", "/payment", token, nil)
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
		if retryAuth {
			newTokens, err := authentication.RefreshTokens()
		if err != nil || newTokens == nil {
			return 0, fmt.Errorf("failed to get user balance: %w: %w", ErrNotAuthenticated, err)
		}

		return client.GetBalance(newTokens.AccessToken, false)
		} else {

			return 0, fmt.Errorf("failed to get user balance: %w", ErrNotAuthenticated)
		}
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
