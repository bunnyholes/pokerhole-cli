package network

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestConnectWithTimeout_Success tests successful connection within timeout
func TestConnectWithTimeout_Success(t *testing.T) {
	// Create test WebSocket server
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade: %v", err)
		}
		defer conn.Close()

		// Read REGISTER message
		_, _, err = conn.ReadMessage()
		if err != nil {
			t.Logf("Read error (expected during test cleanup): %v", err)
		}
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Create client
	client := NewClient(wsURL, "test-uuid", "test-nickname")

	// Connect with 3 second timeout
	err := client.ConnectWithTimeout(3 * time.Second)
	if err != nil {
		t.Fatalf("Expected successful connection, got error: %v", err)
	}

	// Verify connection
	if !client.IsConnected() {
		t.Error("Client should be connected")
	}

	// Cleanup
	client.Close()
}

// TestConnectWithTimeout_Timeout tests connection timeout
func TestConnectWithTimeout_Timeout(t *testing.T) {
	// Use invalid URL that will timeout
	invalidURL := "ws://192.0.2.1:9999/ws"

	// Create client
	client := NewClient(invalidURL, "test-uuid", "test-nickname")

	// Connect with short timeout (100ms)
	start := time.Now()
	err := client.ConnectWithTimeout(100 * time.Millisecond)
	elapsed := time.Since(start)

	// Should fail
	if err == nil {
		t.Fatal("Expected connection to timeout, but it succeeded")
	}

	// Should timeout within reasonable time (< 200ms)
	if elapsed > 200*time.Millisecond {
		t.Errorf("Timeout took too long: %v (expected < 200ms)", elapsed)
	}

	// Verify not connected
	if client.IsConnected() {
		t.Error("Client should not be connected after timeout")
	}
}

// TestConnectWithTimeout_InvalidURL tests connection to invalid URL
func TestConnectWithTimeout_InvalidURL(t *testing.T) {
	// Use invalid URL
	invalidURL := "ws://localhost:99999/invalid"

	// Create client
	client := NewClient(invalidURL, "test-uuid", "test-nickname")

	// Connect with timeout
	err := client.ConnectWithTimeout(1 * time.Second)

	// Should fail
	if err == nil {
		t.Fatal("Expected connection to invalid URL to fail")
	}

	// Verify not connected
	if client.IsConnected() {
		t.Error("Client should not be connected")
	}
}
