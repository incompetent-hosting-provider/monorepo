package utils

import (
	"net"
	"testing"
)

// Tests that the function returns a non-nil listener
func TestGetCallbackServerNonNill(t *testing.T) {
	// ACT
	listener, err := GetCallbackServer()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	// ASSERT
	if listener == nil {
		t.Error("expected listener to not be nil")
	}
	defer listener.Close()
}

// Tests that the function returns a non-nil listener with localhost address
func TestGetCallbackServerLocalhost(t *testing.T) {
	// ARRANGE
	expectedHost := "127.0.0.1"

	// ACT
	listener, err := GetCallbackServer()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	// ASSERT
	if listener == nil {
		t.Error("expected listener to not be nil")
	}
	
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	if addr.IP.String() != expectedHost {
		t.Errorf("unexpected IP address, got: %s, want: %s", addr.IP.String(), expectedHost)
	}
}

// Tests that the function returns a non-nil listener with port != 0
func TestGetCallbackServerNonZeroPort(t *testing.T) {
	// ACT
	listener, err := GetCallbackServer()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	// ASSERT
	if listener == nil {
		t.Error("expected listener to not be nil")
	}
	defer listener.Close()
	
	addr := listener.Addr().(*net.TCPAddr)
	if addr.Port == 0 {
		t.Error("expected non-zero port number")
	}
}

// Tests that the function returns a non-nil listener with TCP protocol
func TestGetCallbackServerTCP(t *testing.T) {
	// ARRANGE
	expectedNetwork := "tcp"

	// ACT
	listener, err := GetCallbackServer()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	// ASSERT
	if listener == nil {
		t.Error("expected listener to not be nil")
	}
	defer listener.Close()
	
	addr := listener.Addr().(*net.TCPAddr)
	if addr.Network() != expectedNetwork {
		t.Errorf("unexpected network, got: %s, want: %s", addr.Network(), expectedNetwork)
	}
}