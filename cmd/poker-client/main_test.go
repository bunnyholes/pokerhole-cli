package main

import (
	"os"
	"testing"
)

// TestGetServerURL_Default tests default server URL
func TestGetServerURL_Default(t *testing.T) {
	// Clear environment variable
	os.Unsetenv("POKERHOLE_SERVER")

	url := getServerURL()

	expected := "ws://localhost:8080/ws/game"
	if url != expected {
		t.Errorf("Expected default URL %s, got %s", expected, url)
	}
}

// TestGetServerURL_FromEnv tests server URL from environment variable
func TestGetServerURL_FromEnv(t *testing.T) {
	// Set environment variable
	customURL := "ws://custom-server:9999/ws/game"
	os.Setenv("POKERHOLE_SERVER", customURL)
	defer os.Unsetenv("POKERHOLE_SERVER")

	url := getServerURL()

	if url != customURL {
		t.Errorf("Expected custom URL %s, got %s", customURL, url)
	}
}

// TestGetServerURL_EnvOverridesDefault tests environment variable takes precedence
func TestGetServerURL_EnvOverridesDefault(t *testing.T) {
	// Set custom URL
	customURL := "ws://prod-server:443/ws/game"
	os.Setenv("POKERHOLE_SERVER", customURL)
	defer os.Unsetenv("POKERHOLE_SERVER")

	url := getServerURL()

	// Should use environment variable, not default
	if url != customURL {
		t.Errorf("Environment variable should override default, got %s", url)
	}

	// Should NOT be default
	defaultURL := "ws://localhost:8080/ws/game"
	if url == defaultURL {
		t.Error("Should not use default URL when environment variable is set")
	}
}
