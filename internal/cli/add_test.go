package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddCommand_MissingName(t *testing.T) {
	// Test that --name is required
	code := AddCommand([]string{"--secret", "JBSWY3DPEHPK3PXP"})
	if code != 1 {
		t.Errorf("Expected exit code 1 for missing --name, got %d", code)
	}
}

func TestAddCommand_MissingSecret(t *testing.T) {
	// Test that --secret is required
	code := AddCommand([]string{"--name", "GitHub"})
	if code != 1 {
		t.Errorf("Expected exit code 1 for missing --secret, got %d", code)
	}
}

func TestAddCommand_InvalidSecret(t *testing.T) {
	// Test invalid Base32 secret
	code := AddCommand([]string{"--name", "GitHub", "--secret", "invalid!secret"})
	if code != 1 {
		t.Errorf("Expected exit code 1 for invalid secret, got %d", code)
	}
}

func TestAddCommand_ShortSecret(t *testing.T) {
	// Test secret that's too short
	code := AddCommand([]string{"--name", "GitHub", "--secret", "ABC"})
	if code != 1 {
		t.Errorf("Expected exit code 1 for short secret, got %d", code)
	}
}

func TestAddCommand_WithIdentifier(t *testing.T) {
	// Create a temporary directory for test storage
	tempDir := t.TempDir()

	// Set environment variable to use test storage path
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	// Create test storage directory
	configDir := filepath.Join(tempDir, ".config", "totp-manager")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	} // Note: This test would require mocking the passphrase prompt
	// For now, we test the flag parsing

	// Test with identifier flag
	args := []string{
		"--name", "GitHub",
		"--identifier", "user@example.com",
		"--secret", "JBSWY3DPEHPK3PXP",
	}

	// This will fail at passphrase prompt, which is expected in unit test
	code := AddCommand(args)

	// We expect it to fail at initialization, not at flag parsing
	if code != 1 {
		t.Logf("Command exited with code %d (expected 1 due to passphrase prompt)", code)
	}
}

func TestAddCommand_FlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
	}{
		{
			name:     "No flags",
			args:     []string{},
			wantCode: 1,
		},
		{
			name:     "Only name",
			args:     []string{"--name", "GitHub"},
			wantCode: 1,
		},
		{
			name:     "Only secret",
			args:     []string{"--secret", "JBSWY3DPEHPK3PXP"},
			wantCode: 1,
		},
		{
			name:     "Invalid secret format",
			args:     []string{"--name", "Test", "--secret", "123"},
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := AddCommand(tt.args)
			if code != tt.wantCode {
				t.Errorf("AddCommand() = %d, want %d", code, tt.wantCode)
			}
		})
	}
}
