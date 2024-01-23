package authentication

import (
	"cli/internal/services/utils"
	"fmt"
	"net/http"
	"sync"
)

// Authentication tokens.
type Authentication struct {
	SessionToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

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
	if err := utils.OpenBrowser(keycloakConfig.getLoginURL(embeddedServerConfig.GetCallbackURL())); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	err := performTokenExchange("Login")
	return err
}

// Starts the registration process.
// This will open the browser and redirect the user to keycloak.
// After the user registered, the session and refresh token will be saved.
func Register() error {
	// Open the browser and redirect the user to the internal server
	if err := utils.OpenBrowser(keycloakConfig.getRegisterURL(embeddedServerConfig.GetCallbackURL())); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	err := performTokenExchange("Registration")
	return err
}


// Performs the token exchange.
// This will start an embedded server and wait for the callback from keycloak.
// After the callback is received, the session and refresh token will be saved.
//
// The action parameter is used to inform the user about the result of the token exchange.
// E.g. "Login successful! You can close this window now." or "Registration successful! You can close this window now."
func performTokenExchange(action string) error {
	// Wait for the embedded server to receive the callback
	var asyncError error
	wg := sync.WaitGroup{}
	wg.Add(1)

	http.HandleFunc(embeddedServerConfig.CallbackPath, func(w http.ResponseWriter, r *http.Request) {
		// Get the autorization code from the query parameters
		code := r.URL.Query().Get("code")

		// Exchange the code for a session and refresh token
	 	auth, err := keycloakConfig.performTokenExchangeRequest(code, embeddedServerConfig.GetCallbackURL())
		if err != nil {
			fmt.Fprintf(w, "%s failed! Return to the CLI to see further details...", action)
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
		fmt.Fprintf(w, "%s successful! You can close this window now.", action)
		wg.Done()
	})

	// Start the embedded server
	go func() {
		http.ListenAndServe(embeddedServerConfig.GetServerURI(), nil)
	}()

	wg.Wait()
	
	return asyncError;
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
