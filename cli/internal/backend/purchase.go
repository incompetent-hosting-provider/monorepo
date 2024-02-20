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

type PurchaseRequest struct {
	Amount int `json:"amount"`
}

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
	req, err := getAuthenticatedRequest("POST", "/purchase", tokens.AccessToken)
	if err != nil {
		return 0, fmt.Errorf("failed to create purchase request: %w", err)
	}

	reqBody, err := json.Marshal(PurchaseRequest{Amount: amount})
	if err != nil {
		return 0, fmt.Errorf("failed to encode purchase request: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewReader(reqBody))
	req.ContentLength = int64(len(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send purchase request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
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