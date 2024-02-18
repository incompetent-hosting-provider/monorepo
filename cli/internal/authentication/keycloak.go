package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	resp, err := http.PostForm(c.getTokenURL(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to send token exchange request: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to exchange code for token: %d - %s", resp.StatusCode, resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body during token exchange: %w", err)
	}

	var tokens *SessionTokens
	err = json.Unmarshal(respBody, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response during token exchange: %w", err)
	}

	return tokens, nil
}

// Performs a token refresh request against the Keycloak Server and return the new tokens.
// If the refresh token has expired nil will be returned.
func (c KeycloakConfig) performTokenRefreshRequest(refreshToken RefreshToken) (*SessionTokens, error) {
	form := url.Values{}
	form.Set("client_id", c.ClientId)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", string(refreshToken))

	resp, err := http.PostForm(c.getTokenURL(), form)
	if err != nil {
		return nil, fmt.Errorf("failed to send token refresh request: %w", err)
	}

	if resp.StatusCode == 401 || resp.StatusCode == 400 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to refresh token: %d - %s", resp.StatusCode, resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body during token refresh: %w", err)
	}

	var tokens *SessionTokens
	err = json.Unmarshal(respBody, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response during token refresh: %w", err)
	}

	return tokens, nil
}
