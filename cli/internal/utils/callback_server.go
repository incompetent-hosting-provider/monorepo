package utils

import (
	"fmt"
	"net"
)

func GetCallbackServer() (*net.TCPListener, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start a callback server: %w", err)
	}

	return listener.(*net.TCPListener), nil
}
