package authentication

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Reads the token from the session file.
// The session tokens can be nil if the user is not authenticated.
func readTokens() (*SessionTokens, error) {
	sessionFilePath, err := getSessionFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get session file path to read tokens: %w", err)
	}

	data, err := os.ReadFile(sessionFilePath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			// No session file exists. This is not really an error but rather
			// an idicator that the user is not authenticated.
			return nil, nil
		default:
			return nil, fmt.Errorf("failed to get session file path to read tokens: %w", err)
		}
	}

	var tokens SessionTokens
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tokens to read tokens: %w", err)
	}

	return &tokens, nil
}

// Writes the token to the session file.
func saveTokens(tokens SessionTokens) error {
	sessionFilePath, err := getSessionFilePath()
	if err != nil {
		return fmt.Errorf("failed to get session file path to save tokens: %w", err)
	}

	data, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to encode the tokens to save them: %w", err)
	}

	err = os.WriteFile(sessionFilePath, data, 0640)
	if err != nil {
		return fmt.Errorf("failed to write the tokens to the file: %w", err)
	}

	return nil
}

// Clears the session file.
func clearTokens() error {
	sessionFilePath, err := getSessionFilePath()
	if err != nil {
		return fmt.Errorf("failed to get session file path to clear tokens: %w", err)
	}

	err = os.Remove(sessionFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove session file to clear tokens: %w", err)
	}

	return nil
}

// Returns the path of the session file
func getSessionFilePath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get session file path: %w", err)
	}

	configDir := filepath.Join(userConfigDir, "ihp")

	// Ensure ihp config dir is existing
	if _, err := os.Stat(configDir); err != nil {
		switch {
		case os.IsNotExist(err):
			if err = os.Mkdir(configDir, 0755); err != nil {
				return "", fmt.Errorf("failed to create IHP config directory at %s: %w", configDir, err)
			}
		default:
			return "", fmt.Errorf("failed to read IHP config directory at %s: %w", configDir, err)
		}
	}

	sessionFilePath := filepath.Join(configDir, "session")
	return sessionFilePath, nil
}
