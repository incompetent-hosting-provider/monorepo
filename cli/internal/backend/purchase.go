package backend

import (
	"bytes"
	"cli/internal/authentication"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)
type PurchaseCreditsRequest struct {
	Amount int `json:"amount"`
}

type PurchaseCreditsResponse struct {
	Balance int `json:"balance"`
}

// PurchaseCredits credits for the user associated to the tokens provided.
func (client BackendClient) PurchaseCredits(token authentication.AccessToken, amount int, retryAuth bool) (int, error){
	reqBody, err := json.Marshal(PurchaseCreditsRequest{Amount: amount})
	if err != nil {
		return 0, fmt.Errorf("failed to create purchase request: %w", err)
	}

	req, err := client.buildAuthenticatedRequest("POST", "/payment", token, bytes.NewBuffer(reqBody))
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
		if retryAuth {
			newTokens, err := authentication.RefreshTokens()
			if err != nil || newTokens == nil{
				return 0, fmt.Errorf("failed to purchase credits: %w: %w", ErrNotAuthenticated, err)
			}
	
			return client.PurchaseCredits(newTokens.AccessToken, amount, false)
		} else {
			return 0, fmt.Errorf("failed to purchase credits: %w", ErrNotAuthenticated)
		}
	default:
		return 0, fmt.Errorf("failed to purchase credits: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body of purchase request: %w", err)
	}

	var purchaseResponse PurchaseCreditsResponse
	err = json.Unmarshal(body, &purchaseResponse)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response body of purchase request: %w", err)
	}
	
	return purchaseResponse.Balance, nil
}