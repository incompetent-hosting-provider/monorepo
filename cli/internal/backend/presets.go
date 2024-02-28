package backend

import (
	"cli/internal/authentication"
	"cli/internal/models"
	"encoding/json"
	"fmt"
	"io"
)

type GetInstancePresetsResponse struct {
	Presets []models.InstancePreset `json:"presets,omitempty"`
}

// Gets the presets for the available instances.
func (client BackendClient) GetInstancePresets(token authentication.AccessToken, authRetry bool) ([]models.InstancePreset, error) {
	req, err := client.buildAuthenticatedRequest("GET", "/service/available-presets", token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get instance presets request: %w", err)
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance presets: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		if authRetry {
			newTokens, err := authentication.RefreshTokens()
			if err != nil || newTokens == nil {
				return nil, fmt.Errorf("failed to get instance presets: %w: %w", ErrNotAuthenticated, err)
			}
			return client.GetInstancePresets(newTokens.AccessToken, false)
		} else {
			return nil, fmt.Errorf("failed to get instance presets: %w: %w", ErrNotAuthenticated, err)
		}
	default:
		return nil, fmt.Errorf("failed to get instance presets: %s", resp.Status)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read instance presets response: %w", err)
	}

	var presets *GetInstancePresetsResponse
	err = json.Unmarshal(body, &presets)
	if err != nil {
		return nil, fmt.Errorf("failed to parse instance presets: %w", err)
	}

	return presets.Presets, nil
}