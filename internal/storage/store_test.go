package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestStore_CreateAndLoad tests creating and loading encrypted storage
func TestStore_CreateAndLoad(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-secrets.enc")

	passphrase := "test-passphrase-123"

	// Create new storage
	store, err := Create(storePath, passphrase)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Add a service
	service := Service{
		Name:      "GitHub",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}
	err = store.AddService(service)
	if err != nil {
		t.Fatalf("AddService() error = %v", err)
	}

	// Save
	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		t.Fatal("Storage file was not created")
	}

	// Load storage
	loaded, err := Load(storePath, passphrase)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify service was loaded
	if len(loaded.Services) != 1 {
		t.Errorf("Loaded services count = %d, want 1", len(loaded.Services))
	}

	if loaded.Services[0].Name != "GitHub" {
		t.Errorf("Loaded service name = %s, want GitHub", loaded.Services[0].Name)
	}
}

// TestStore_LoadWithWrongPassphrase tests that loading fails with wrong passphrase
func TestStore_LoadWithWrongPassphrase(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-secrets.enc")

	correctPass := "correct-passphrase"
	wrongPass := "wrong-passphrase"

	// Create storage
	store, err := Create(storePath, correctPass)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Try to load with wrong passphrase
	_, err = Load(storePath, wrongPass)
	if err == nil {
		t.Error("Load() should fail with wrong passphrase, but succeeded")
	}
}

// TestStore_FilePermissions tests that storage file has correct permissions (0600)
func TestStore_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-secrets.enc")

	passphrase := "test-passphrase"

	// Create and save storage
	store, err := Create(storePath, passphrase)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check file permissions
	info, err := os.Stat(storePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	mode := info.Mode()
	expected := os.FileMode(0600)

	if mode.Perm() != expected {
		t.Errorf("File permissions = %o, want %o", mode.Perm(), expected)
	}
}

// TestStore_AtomicWrite tests that writes are atomic (temp file + rename)
func TestStore_AtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-secrets.enc")

	passphrase := "test-passphrase"

	// Create storage
	store, err := Create(storePath, passphrase)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Add service
	service := Service{
		Name:      "GitHub",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}
	err = store.AddService(service)
	if err != nil {
		t.Fatalf("AddService() error = %v", err)
	}

	// Save
	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify no temp files left behind
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".tmp" {
			t.Errorf("Temporary file left behind: %s", entry.Name())
		}
	}
}

// TestStore_EncryptedContent tests that file content is encrypted
func TestStore_EncryptedContent(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-secrets.enc")

	passphrase := "test-passphrase"
	secret := "JBSWY3DPEHPK3PXP"

	// Create storage with service
	store, err := Create(storePath, passphrase)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	service := Service{
		Name:      "GitHub",
		Secret:    secret,
		CreatedAt: time.Now(),
	}
	err = store.AddService(service)
	if err != nil {
		t.Fatalf("AddService() error = %v", err)
	}

	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read raw file content
	content, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	// Verify secret is not in plaintext
	contentStr := string(content)
	if contains(contentStr, secret) {
		t.Error("Secret found in plaintext in encrypted file")
	}

	// Verify service name is not in plaintext
	if contains(contentStr, "GitHub") {
		t.Error("Service name found in plaintext in encrypted file")
	}
}

// TestStore_ReEncrypt tests re-encryption with new passphrase
func TestStore_ReEncrypt(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-secrets.enc")

	oldPass := "old-passphrase"
	newPass := "new-passphrase"

	// Create storage with old passphrase
	store, err := Create(storePath, oldPass)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	service := Service{
		Name:      "GitHub",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}
	err = store.AddService(service)
	if err != nil {
		t.Fatalf("AddService() error = %v", err)
	}

	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load with old passphrase
	store, err = Load(storePath, oldPass)
	if err != nil {
		t.Fatalf("Load() with old passphrase error = %v", err)
	}

	// Re-encrypt with new passphrase
	err = store.ChangePassphrase(newPass)
	if err != nil {
		t.Fatalf("ChangePassphrase() error = %v", err)
	}

	// Try to load with old passphrase (should fail)
	_, err = Load(storePath, oldPass)
	if err == nil {
		t.Error("Load() with old passphrase should fail after re-encryption")
	}

	// Load with new passphrase (should succeed)
	store, err = Load(storePath, newPass)
	if err != nil {
		t.Fatalf("Load() with new passphrase error = %v", err)
	}

	// Verify data integrity
	if len(store.Services) != 1 {
		t.Errorf("Services count = %d, want 1", len(store.Services))
	}
}

// TestStore_MultipleServices tests storage with multiple services
func TestStore_MultipleServices(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test-secrets.enc")

	passphrase := "test-passphrase"

	// Create storage
	store, err := Create(storePath, passphrase)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Add multiple services
	services := []Service{
		{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
		{Name: "AWS", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
		{Name: "Google", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
	}

	for _, svc := range services {
		err = store.AddService(svc)
		if err != nil {
			t.Fatalf("AddService() error = %v", err)
		}
	}

	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load and verify
	loaded, err := Load(storePath, passphrase)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded.Services) != 3 {
		t.Errorf("Loaded services count = %d, want 3", len(loaded.Services))
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findIndex(s, substr) >= 0
}

func findIndex(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
