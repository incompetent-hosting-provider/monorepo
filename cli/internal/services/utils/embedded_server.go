package utils

import "fmt"

// A configuration for the embedded server.
type EmbeddedServerConfig struct {
	// The Port of the embedded server.
	Port int `json:"port"`

	// The Path of the callback endpoint.
	// Should start with a slash.
	CallbackPath string `json:"login_callback_path"`
}

// Returns the URI of the embedded server.
// E.g. localhost:3000
func (c EmbeddedServerConfig) GetServerURI() string {
	return fmt.Sprintf("localhost:%v", c.Port)
}

// Returns the URL of the Authentication Callback Endpoint.
func (c EmbeddedServerConfig) GetCallbackURL() string {
	return fmt.Sprintf("http://%s%s", c.GetServerURI(), c.CallbackPath)
}