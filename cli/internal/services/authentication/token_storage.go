package authentication

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	// The name of the file where the session and refresh token are stored.
	fileName = "session"
)


// Saves the session and refresh token to the config directory.
func saveTokens(tokens Authentication) error {
	data, err := json.Marshal(tokens)
	if(err != nil) {
		return err
	}

	file, err := getConfigPath()
	if(err != nil) {
		return err
	}

	err = os.WriteFile(file, data, 0640)
	if(err != nil) {
		return err
	}

	return nil
}

// Reads the session and refresh token from the config directory.
func readTokens() (Authentication, error) {
	file, err := getConfigPath()
	if(err != nil) {
		return Authentication{}, err
	}

	data, err := os.ReadFile(file)
	if(err != nil) {
		return Authentication{}, err
	}

	var tokens Authentication
	err = json.Unmarshal(data, &tokens)
	if(err != nil) {
		return Authentication{}, err
	}

	return tokens, nil
}

// Deletes the session and refresh token from the config directory.
func clearToken() error {
	file, err := getConfigPath()
	if(err != nil) {
		return err
	}

	err = os.Remove(file)
	if(err != nil) {
		return err
	}

	return nil
}

// Returns the directory where the tokens are stored.
// If the directory does not exist, it will be created.
func getConfigDirectory() (string, error) {
	configDir,err := os.UserConfigDir()
	if(err != nil) {
		return "", err
	}

	configDirPath := filepath.Join(configDir, "ihp")

	// Create the directory if it does not exist
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755)
		if err != nil {
			return "", err
		}
	}

	return configDirPath, nil
}

// Returns the path to the file where the tokens are stored.
// If the directory containing the file does not exist, it will be created.
func getConfigPath() (string, error) {
	configDir,err := getConfigDirectory()
	if(err != nil) {
		return "", err
	}

	return filepath.Join(configDir, fileName), nil
}