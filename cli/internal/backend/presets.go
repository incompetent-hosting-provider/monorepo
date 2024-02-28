package backend

import (
	"cli/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Gets the presets for the available instances.
func (client BackendClient) GetInstancePresets() ([]models.InstancePreset, error) {
	resp, err := http.Get(client.baseURL + "/service/available-presets")
	if err != nil {
		return nil, fmt.Errorf("failed to get instance presets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get instance presets: %s", resp.Status)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read instance presets response: %w", err)
	}

	var presets []models.InstancePreset
	err = json.Unmarshal(body, &presets)
	if err != nil {
		return nil, fmt.Errorf("failed to parse instance presets: %w", err)
	}

	return presets, nil
}