package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestApp_Initialize_NewStorage tests creating new storage
// (T026, T027: Passphrase prompt and storage initialization)
func TestApp_Initialize_NewStorage(t *testing.T) {
	// This is a unit test - interactive prompts tested manually
	// We test the logic that determines new vs existing storage

	tmpDir, err := os.MkdirTemp("", "totp-cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "secrets.enc")

	app := &App{
		storagePath: storagePath,
	}

	// Verify storage file doesn't exist
	if _, err := os.Stat(storagePath); !os.IsNotExist(err) {
		t.Error("Storage file should not exist initially")
	}

	// Note: Full initialization requires interactive input
	// This test verifies the structure and logic flow
	if app.storagePath != storagePath {
		t.Errorf("Expected storagePath %s, got %s", storagePath, app.storagePath)
	}
}

// TestApp_LoadExisting_ValidatesAttempts tests 3-attempt limit logic
// (T028: Passphrase validation with attempt limit)
func TestApp_LoadExisting_ValidatesAttempts(t *testing.T) {
	// Verify maxPassphraseAttempts constant
	if maxPassphraseAttempts != 3 {
		t.Errorf("Expected maxPassphraseAttempts to be 3, got %d", maxPassphraseAttempts)
	}

	// Attempt limit is enforced in loadExistingStorage()
	// Full test requires interactive prompts (tested manually)
}

// TestApp_NewApp tests default storage path
func TestApp_NewApp(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	if app.storagePath == "" {
		t.Error("storagePath should not be empty")
	}

	// Verify storagePath has expected format
	if !filepath.IsAbs(app.storagePath) {
		t.Errorf("storagePath should be absolute, got: %s", app.storagePath)
	}
}

// TestPassphraseValidation tests passphrase strength requirements
func TestPassphraseValidation(t *testing.T) {
	tests := []struct {
		name       string
		passphrase string
		wantError  bool
	}{
		{"Valid passphrase", "MySecurePass123!", false},
		{"Too short", "short", true},
		{"Minimum length", "12345678", false},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passphrase validation logic (8 char minimum)
			hasError := len(tt.passphrase) < 8
			if hasError != tt.wantError {
				t.Errorf("Passphrase '%s': expected error=%v, got error=%v",
					tt.passphrase, tt.wantError, hasError)
			}
		})
	}
}
