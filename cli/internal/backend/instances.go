package backend

import (
	"bytes"
	"cli/internal/authentication"
	"cli/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetInstancesResponse struct {
	Instances []models.Instance `json:"instances,omitempty"`
}

// Gets the instance with the given id belonging to the user associate to the tokens.
func (client BackendClient) GetUserInstances(token authentication.AccessToken, authRetry bool) ([]models.Instance, error) {
	req, err := client.buildAuthenticatedRequest("GET", "/instances", token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get instances request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user instances: %w", err)
	}

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		if authRetry {
			newTokens, err := authentication.RefreshTokens()
			if err != nil || newTokens == nil {
				return nil, fmt.Errorf("failed to get user instances: %w: %w", ErrNotAuthenticated, err)
			}
			return client.GetUserInstances(newTokens.AccessToken, false)
		} else {
			return nil, fmt.Errorf("failed to get user instances: %w", ErrNotAuthenticated)
		}
	default:
		return nil, fmt.Errorf("failed to get user instances: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body of user instances request: %w", err)
	}

	var instancesResponse *GetInstancesResponse
	err = json.Unmarshal(body, &instancesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body of user instances request: %w", err)
	}

	return instancesResponse.Instances, nil
}

// Gets the instance with the given id belonging to the user associate to the tokens.
func (client BackendClient) GetUserInstance(token authentication.AccessToken, instanceID string, authRetry bool) (*models.InstanceDetail, error) {
	req, err := client.buildAuthenticatedRequest("GET", fmt.Sprintf("/instances/%s", instanceID), token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get instance request with id %s: %w", instanceID, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance with id %s: %w", instanceID, err)
	}

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		if authRetry {
			newTokens, err := authentication.RefreshTokens()
			if err != nil || newTokens == nil {
				return nil, fmt.Errorf("failed to get user instance with id %s: %w: %w", instanceID, ErrNotAuthenticated, err)
			}
	
			return client.GetUserInstance(newTokens.AccessToken, instanceID, false)
		} else {
			return nil, fmt.Errorf("failed to get instance with id %s: %w", instanceID, ErrNotAuthenticated)
		}
	default:
		return nil, fmt.Errorf("failed to get instance with id %s: %s", instanceID, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body of user instance request with id %s: %w", instanceID, err)
	}

	var instanceDetail *models.InstanceDetail
	err = json.Unmarshal(body, &instanceDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body of user instance request with id %s: %w", instanceID, err)
	}

	return instanceDetail, nil
}


type CreatePresetInstanceRequest struct {
	PresetID string `json:"preset,omitempty"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type CreateInstanceResponse struct {
	InstanceID string `json:"instance_id,omitempty"`
}

// Creates a custom instance for the user associated to the tokens provided.
// Returns the id of the created instance.
func (client BackendClient) CreatePresetInstance(token authentication.AccessToken, request CreatePresetInstanceRequest, authRetry bool) (string, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal create preset instance request with id %s: %w", request.PresetID, err)
	}

	req, err := client.buildAuthenticatedRequest("POST", "/instances/preset/", token, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create create preset instance request with id %s: %w", request.PresetID, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create preset instance with id %s: %w", request.PresetID, err)
	}

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		if authRetry {
			newTokens, err := authentication.RefreshTokens()
			if err != nil || newTokens == nil {
				return "", fmt.Errorf("failed to create preset instance with id %s: %w: %w", request.PresetID, ErrNotAuthenticated, err)
			}
	
			return client.CreatePresetInstance(newTokens.AccessToken, request, false)
		} else {
			return "", fmt.Errorf("failed to create preset instance with id %s: %w", request.PresetID, ErrNotAuthenticated)
		}
	default:
		return "", fmt.Errorf("failed to create preset instance with id %s: %s", request.PresetID, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to read response body of create preset instance request with id %s: %w", request.PresetID, err)
	}

	var createInstanceResponse *CreateInstanceResponse
	err = json.Unmarshal(body, &createInstanceResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body of create preset instance request with id %s: %w", request.PresetID, err)
	}

	return createInstanceResponse.InstanceID, nil
}

type CreateCustomInstanceRequest struct {
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Image models.ContainerImage `json:"image,omitempty"`
	EnvVars map[string]string `json:"env-vars,omitempty"`
	Ports []int `json:"ports,omitempty"`
}


// Creates a custom instance for the user associated to the tokens provided.
// Returns the id of the created instance.
func (client BackendClient) CreateCustomInstance(token authentication.AccessToken, r CreateCustomInstanceRequest, authRetry bool) (string, error) {
	reqBody, err := json.Marshal(r)
	if err != nil {
		return "",fmt.Errorf("failed to marshal create custom instance request: %w", err)
	}

	req, err := client.buildAuthenticatedRequest("POST", "/instances/custom", token, bytes.NewBuffer(reqBody))
	if err != nil {
		return "",fmt.Errorf("failed to create create custom instance request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "",fmt.Errorf("failed to create custom instance: %w", err)
	}

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		if authRetry {
			newTokens, err := authentication.RefreshTokens()
			if err != nil || newTokens == nil {
				return "", fmt.Errorf("failed to create custom instance: %w: %w", ErrNotAuthenticated, err)
			}
	
			return client.CreateCustomInstance(newTokens.AccessToken, r, false)
		} else {
			return "",fmt.Errorf("failed to create custom instance: %w", ErrNotAuthenticated)
		}
	default:
		return "",fmt.Errorf("failed to create custom instance: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "",fmt.Errorf("failed to read response body of create custom instance request: %w", err)
	}

	var createInstanceResponse *CreateInstanceResponse
	err = json.Unmarshal(body, &createInstanceResponse)
	if err != nil {
		return "",fmt.Errorf("failed to parse response body of create custom instance request: %w", err)
	}

	return createInstanceResponse.InstanceID, nil
}


// Deletes the instance with the given id belonging to the user associate to the tokens.
func (client BackendClient) DeleteInstance(token authentication.AccessToken, instanceID string, authRetry bool) error {
	req, err := client.buildAuthenticatedRequest("DELETE", fmt.Sprintf("/instances/%s", instanceID), token, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete instance request with id %s: %w", instanceID, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete instance with id %s: %w", instanceID, err)
	}

	switch resp.StatusCode {
	case 200:
		break
	case 401:
		if authRetry {
			newTokens, err := authentication.RefreshTokens()
			if err != nil || newTokens == nil {
				return fmt.Errorf("failed to delete instance with id %s: %w: %w", instanceID, ErrNotAuthenticated, err)
			}
	
			return client.DeleteInstance(newTokens.AccessToken, instanceID, false)
		} else {
			return fmt.Errorf("failed to delete instance with id %s: %w", instanceID, ErrNotAuthenticated)
		}
	default:
		return fmt.Errorf("failed to delete instance with id %s: %s", instanceID, resp.Status)
	}

	return nil
}