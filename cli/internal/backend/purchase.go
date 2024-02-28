package backend

import (
	"bytes"
	"cli/internal/authentication"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PurchaseResponse struct {
	Balance int `json:"balance"`
}

// Purchase credits for the user associated to the tokens provided.
func Purchase(tokens authentication.SessionTokens, amount int) (int, error) {
	newAmount, err := makePurchaseRequest(tokens, amount)
	if err != nil && errors.Is(err, ErrNotAuthenticated) {
		newTokens, err := authentication.RefreshTokens()
		if err != nil || newTokens != nil{
			return 0, fmt.Errorf("failed to purchase credits: %w: %w", ErrNotAuthenticated, err)
		}

		return makePurchaseRequest(*newTokens, amount)
	} else if err != nil {
		return 0, err
	}

	return newAmount, nil
}

func makePurchaseRequest(tokens authentication.SessionTokens, amount int) (int, error){
	reqJson, err := json.Marshal(map[string]any{
		"amount": amount,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create purchase request: %w", err)
	}

	req, err := getAuthenticatedRequest("POST", "/payment", tokens.AccessToken, bytes.NewBuffer(reqJson))
	if err != nil {
		return 0, fmt.Errorf("failed to create purchase request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send purchase request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
	case 201:
		break
	case 401:
		return 0, fmt.Errorf("failed to purchase credits: %w", ErrNotAuthenticated)
	default:
		return 0, fmt.Errorf("failed to purchase credits: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body of purchase request: %w", err)
	}

	var purchaseResponse PurchaseResponse
	err = json.Unmarshal(body, &purchaseResponse)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response body of purchase request: %w", err)
	}
	
	return purchaseResponse.Balance, nil
}