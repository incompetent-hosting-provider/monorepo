package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// The Keycloak Configuration
var DefaultKeycloakConfig = KeycloakConfig{
	URL:      "http://localhost:8080",
	Realm:    "ihp-realm",
	ClientId: "ihp-cli",
}

// Configuration for the Keycloak Server.
type KeycloakConfig struct {
	URL      string `json:"url"`
	Realm    string `json:"realm"`
	ClientId string `json:"client_id"`
}

// Returns the Authentication URI of the Keycloak Server.
func (c KeycloakConfig) GetLoginURL(redirect string) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_mode=query&response_type=code&scope=openid",
		c.URL,
		c.Realm,
		c.ClientId,
		redirect,
	)
}

// Returns the Registration URI of the Keycloak Server.
func (c KeycloakConfig) GetRegisterURL(redirect string) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/registrations?client_id=%s&redirect_uri=%s&response_mode=query&response_type=code&scope=openid",
		c.URL,
		c.Realm,
		c.ClientId,
		redirect,
	)
}

// Returns the Registration URI of the Keycloak Server.
func (c KeycloakConfig) GetLogoutURL(redirect string) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout?client_id=%s&post_logout_redirect_uri=%s&response_mode=query&response_type=code&scope=openid",
		c.URL,
		c.Realm,
		c.ClientId,
		redirect,
	)
}

// Returns the Token URI (used for Token Exchange) of the Keycloak Server.
func (c KeycloakConfig) getTokenURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.URL, c.Realm)
}

// Performs the Token Exchange Request against the Keycloak Server.
func (c KeycloakConfig) performTokenExchangeRequest(authenticationCode string, redirect string) (*SessionTokens, error) {
	body := url.Values{}
	body.Set("client_id", c.ClientId)
	body.Set("grant_type", "authorization_code")
	body.Set("code", authenticationCode)
	body.Set("redirect_uri", redirect)

	req, err := http.NewRequest("POST", c.getTokenURL(), strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send token exchange request: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to exchange code for token: %d - %s", resp.StatusCode, resp.Status)
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokens *SessionTokens
	err = json.Unmarshal(respBody, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return tokens, nil
}
