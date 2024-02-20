package utils

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

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

func GetCallbackServer() (*net.TCPListener, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start a callback server: %w", err)
	}

	return listener.(*net.TCPListener), nil
}
