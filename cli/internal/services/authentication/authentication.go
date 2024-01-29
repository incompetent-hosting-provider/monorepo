package authentication

import (
	"cli/internal/services/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// Authentication tokens.
type Authentication struct {
	SessionToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

const (
	// The name of the file where the session and refresh token are stored.
	fileName = "session"
)

var (
	// The Keycloak Configuration
	// TODO: Read this from a config file or environment variables
	keycloakConfig = KeycloakConfig{
		URL: "http://localhost:8080",
		Realm: "ihp-realm",
		ClientId: "ihp-cli",
	}

	// The Embedded Server Configuration
	// TODO: Read this from a config file or environment variables
	embeddedServerConfig = utils.EmbeddedServerConfig{
		Port: 3000,
		CallbackPath: "/sso-cb",
	}
)

// Starts the login process.
// This will open the browser and redirect the user to keycloak.
// After the user logged in, the session and refresh token will be saved.
// TODO: If the user is already logged in, the login process will be skipped.
func Login() error {
	// Open the browser and redirect the user to the internal server
	if err := utils.OpenBrowser(keycloakConfig.getLoginURL("http://" + embeddedServerConfig.CallbackURL())); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	// Wait for the embedded server to receive the callback
	var asyncError error
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start the embedded server
	embeddedServerConfig.StartServer(func(w http.ResponseWriter, r *http.Request) {
		// Get the autorization code from the query parameters
		code := r.URL.Query().Get("code")

		// Exchange the code for a session and refresh token
	 	auth, err := keycloakConfig.performTokenExchangeRequest(code, "http://" + embeddedServerConfig.CallbackURL())
		if err != nil {
			fmt.Fprintf(w, "Login failed! Return to the CLI to see further details...")
			asyncError = fmt.Errorf("failed to exchange code for token: %w", err)
			wg.Done()
			return
		}

		// Save the authentication
		if err := saveTokens(auth); err != nil {
			asyncError = fmt.Errorf("failed to save tokens: %w", err)
			wg.Done()
			return
		}

		// Inform the user to close the browser window
		fmt.Fprintf(w, "Login successful! You can close this window now.")
		wg.Done()
	})

	// Wait for the embedded server to receive the callback
	wg.Wait()
	return asyncError
}

// Starts the registration process.
// This will open the browser and redirect the user to keycloak.
// After the user registered, the session and refresh token will be saved.
func Register() error {
	// Pretend we getting a session and refresh token from keycloak
	sessionToken:= "sessionToken"
	refreshToken:= "refreshToken"

	tokens := Authentication{	
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	}

	return saveTokens(tokens)
}

// Starts the logout process.
// This will delete the session and refresh token.
func Logout() error {
	return clearToken()
}

// Returns the session token.
func GetSessionToken() (string, error) {
	tokens, err := readTokens()
	if(err != nil) {
		clearToken()
		return "", err
	}

	return tokens.SessionToken, nil
}

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