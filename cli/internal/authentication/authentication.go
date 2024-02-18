package authentication

import (
	"fmt"
	"net"
	"net/http"
)

func GetCurrentAuthentication() (*SessionTokens, error) {
	tokens, err := readTokens()
	if err != nil {
		// Unable to read tokens, clear them and
		// assume that the user was not logged in.
		if clearErr := clearTokens(); clearErr != nil {
			return nil, fmt.Errorf(
				"failed to read tokens to get current authentication: %w."+
					"Tried to clear tokens but failed aswell: %w", err, clearErr)
		}

		return nil, nil
	}

	return tokens, nil
}

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

	http.Serve(server, mux)
}

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

	http.Serve(server, mux)
}
