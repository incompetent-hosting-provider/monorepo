package authentication

import (
	"fmt"
	"net"
	"net/http"
)

// Reads the tokens from the storage and returns them.
// If no tokens are present, it will return nil.
//
// If an error occurs while reading the tokens, it tries to clear them
// and assumes the user is not authenticated.
func GetCurrentAuthentication() *SessionTokens {
	tokens, err := readTokens()
	if err != nil {
		// Unable to read tokens, clear them and
		// assume that the user was not logged in.
		_ = clearTokens()
		return nil
	}

	return tokens
}

// Performs a token refresh request and saves the new tokens.
// If the refresh token has expired, the tokens will be cleared and nil will be returned.
func RefreshTokens() (*SessionTokens, error) {
	tokens := GetCurrentAuthentication()
	if tokens == nil {
		return nil, nil
	}

	// Perform the token refresh request
	newTokens, err := DefaultKeycloakConfig.performTokenRefreshRequest(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh tokens: %w", err)
	}

	// If the new tokens are nil and we dont have an error,
	// it means that the refresh token has expired.
	// In this case, we should clear the tokens and return nil.
	if newTokens == nil {
		_ = clearTokens()
		return nil, nil
	}

	// Save the new tokens
	err = saveTokens(*newTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	return newTokens, nil
}

// Starts the token exchange process by listening on the given server
// to receive the callback with an authentication code.
// The result channel will be used to signal the result of the token exchange.
func PerformTokenExchange(server *net.TCPListener, result chan<- error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "ERROR: The authentication code is missing.")
			result <- fmt.Errorf("authentication code missing")
			return
		}

		url := fmt.Sprintf("http://%s", server.Addr().String())
		tokens, err := DefaultKeycloakConfig.performTokenExchangeRequest(code, url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "ERROR: Failed to exchange the authentication code for tokens.")
			result <- fmt.Errorf("failed to exchange the authentication code for tokens: %w", err)
			return
		}

		err = saveTokens(*tokens)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "ERROR: Failed to save the tokens.")
			result <- fmt.Errorf("failed to save tokens after token exchange: %w", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success. Please return back to the CLI.")

		result <- nil
	})

	err := http.Serve(server, mux)

	if err != nil {
		fmt.Printf("Could not start callback server due to an error: %v", err)
	}
}

// Starts the logout process by listening on the given server
// to receive the callback. The result channel will be used to signal the result of the logout.
func PerformLogout(server *net.TCPListener, result chan<- error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := clearTokens(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "ERROR: Failed to clear the tokens.")
			result <- fmt.Errorf("failed to clear tokens after logout: %w", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success. Please return back to the CLI.")

		result <- nil
	})

	err := http.Serve(server, mux)
	if err != nil {
		fmt.Printf("Could not start callback server due to an error: %v", err)
	}
}
