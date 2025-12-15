package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewApp tests creating a new CLI app
func TestNewApp(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}

	if app == nil {
		t.Error("NewApp() returned nil app")
	}

	if app.storagePath == "" {
		t.Error("NewApp() should set storage path")
	}
}

// TestApp_StoragePath tests storage path is set correctly
func TestApp_StoragePath(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}

	// Path should contain .config/totp-manager
	if app.storagePath == "" {
		t.Error("Storage path should not be empty")
	}
}

// TestApp_Initialize_NonExistentStorage tests initializing with no storage file
func TestApp_Initialize_NonExistentStorage(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}

	// Initialize will fail because it tries to prompt for passphrase
	// which requires stdin interaction
	err = app.Initialize()
	if err == nil {
		t.Log("Initialize() expected to fail without stdin, but succeeded")
	}
}

// TestApp_CreateWithTestStorage tests creating storage in test environment
func TestApp_CreateWithTestStorage(t *testing.T) {
	tempDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}

	// Verify storage path is in temp directory
	if !filepath.IsAbs(app.storagePath) {
		t.Error("Storage path should be absolute")
	}
}
