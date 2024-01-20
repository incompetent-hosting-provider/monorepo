package utils

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

// A configuration for the embedded server.
type EmbeddedServerConfig struct {
	// The Port of the embedded server.
	Port int `json:"port"`

	// The Path of the callback endpoint.
	// Should start with a slash.
	CallbackPath string `json:"login_callback_path"`
}

// Returns the URL of the embedded server.
// E.g. http://localhost:3000
func (c EmbeddedServerConfig) ServerURL() string {
	return fmt.Sprintf("localhost:%v", c.Port)
}

// Returns the URL of the Authentication Callback Endpoint.
func (c EmbeddedServerConfig) CallbackURL() string {
	return fmt.Sprintf("%s%s", c.ServerURL(), c.CallbackPath)
}

// Starts the embedded server.
func (c EmbeddedServerConfig) StartServer(
	onRequest func(w http.ResponseWriter, r *http.Request),
) {
	http.HandleFunc(c.CallbackPath, onRequest)
	go func () {http.ListenAndServe(c.ServerURL(), nil)}()
}

// Opens the systems browser with the given URL.
func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}