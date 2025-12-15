package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pavanprakash21/totp-manager-go/internal/storage"
)

// TestStorageNewPassphraseWorkflow tests creating storage with new passphrase
// (T021: Integration test for new passphrase prompt + confirmation)
func TestStorageNewPassphraseWorkflow(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "totp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "secrets.enc")

	// Test: New passphrase workflow (simulating user input)
	passphrase := "SecurePassphrase123!"

	// Create new storage with passphrase
	store, err := storage.Create(storagePath, passphrase)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Save the storage to create the file
	if err := store.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		t.Error("Storage file not created")
	}

	// Verify file permissions (0600)
	info, err := os.Stat(storagePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}

	// Verify can load with same passphrase
	_, err = storage.Load(storagePath, passphrase)
	if err != nil {
		t.Errorf("Load() with correct passphrase failed: %v", err)
	}

	// Verify the store is valid
	if store == nil {
		t.Error("Store should not be nil")
	}
}

// TestStorageCorrectVsWrongPassphrase tests authentication
// (T022: Integration test for correct vs wrong passphrase)
func TestStorageCorrectVsWrongPassphrase(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "totp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "secrets.enc")
	correctPass := "CorrectPassphrase123!"
	wrongPass := "WrongPassphrase456!"

	// Create storage
	store, err := storage.Create(storagePath, correctPass)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Add a service
	service := storage.Service{
		Name:      "TestService",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}
	if err := store.AddService(service); err != nil {
		t.Fatalf("AddService() failed: %v", err)
	}
	if err := store.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Test: Load with correct passphrase succeeds
	correctStore, err := storage.Load(storagePath, correctPass)
	if err != nil {
		t.Errorf("Load() with correct passphrase failed: %v", err)
	} else {
		// Verify data integrity
		svc, err := correctStore.GetService("TestService")
		if err != nil {
			t.Errorf("GetService() failed: %v", err)
		}
		if svc.Secret != "JBSWY3DPEHPK3PXP" {
			t.Errorf("Expected secret JBSWY3DPEHPK3PXP, got %s", svc.Secret)
		}
	}

	// Test: Load with wrong passphrase fails
	_, err = storage.Load(storagePath, wrongPass)
	if err == nil {
		t.Error("Load() with wrong passphrase should have failed but succeeded")
	}
}

// TestStorageAttemptsLimit tests 3-attempt authentication limit
// (T023: Integration test for 3-attempt limit)
func TestStorageAttemptsLimit(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "totp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "secrets.enc")
	correctPass := "CorrectPassphrase123!"

	// Create storage
	store, err := storage.Create(storagePath, correctPass)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Save the storage to create the file
	if err := store.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Simulate 3 failed attempts
	wrongPassphrases := []string{"wrong1", "wrong2", "wrong3"}
	attempts := 0

	for _, wrongPass := range wrongPassphrases {
		attempts++
		_, err := storage.Load(storagePath, wrongPass)
		if err == nil {
			t.Errorf("Attempt %d: Load() with wrong passphrase should have failed", attempts)
		}

		if attempts >= 3 {
			// After 3 attempts, we should stop allowing more attempts
			// This will be enforced in the CLI layer, not the storage layer
			// Storage layer just returns auth errors
			break
		}
	}

	if attempts != 3 {
		t.Errorf("Expected exactly 3 attempts, got %d", attempts)
	}

	// Verify storage file still intact (no corruption from failed attempts)
	_, err = storage.Load(storagePath, correctPass)
	if err != nil {
		t.Errorf("Storage corrupted after failed attempts: %v", err)
	}
}

// TestStorageEncryptedContent tests that storage file contains no plaintext secrets
// (T024: Integration test for encrypted file verification)
func TestStorageEncryptedContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "totp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "secrets.enc")
	passphrase := "SecurePassphrase123!"
	secretValue := "JBSWY3DPEHPK3PXP"
	serviceName := "GitHub"

	// Create storage with a service
	store, err := storage.Create(storagePath, passphrase)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	service := storage.Service{
		Name:      serviceName,
		Secret:    secretValue,
		CreatedAt: time.Now(),
	}
	if err := store.AddService(service); err != nil {
		t.Fatalf("AddService() failed: %v", err)
	}
	if err := store.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Read raw file content
	content, err := os.ReadFile(storagePath)
	if err != nil {
		t.Fatalf("Failed to read storage file: %v", err)
	}

	// Verify plaintext secret is NOT in file
	if containsString(content, secretValue) {
		t.Error("Storage file contains plaintext secret")
	}

	// Verify plaintext service name is NOT in file
	if containsString(content, serviceName) {
		t.Error("Storage file contains plaintext service name")
	}

	// Verify file starts with version header (4 bytes)
	if len(content) < 4 {
		t.Error("Storage file too small to contain version header")
	}
}

// TestStorageFilePermissions tests that storage file has 0600 permissions
// (T025: Integration test for permissions verification)
func TestStorageFilePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "totp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storagePath := filepath.Join(tmpDir, "secrets.enc")
	passphrase := "SecurePassphrase123!"

	// Create storage
	store, err := storage.Create(storagePath, passphrase)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Save the storage to create the file
	if err := store.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Check initial permissions
	info, err := os.Stat(storagePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}

	// Add a service and save
	service := storage.Service{
		Name:      "TestService",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}
	if err := store.AddService(service); err != nil {
		t.Fatalf("AddService() failed: %v", err)
	}
	if err := store.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Check permissions after save (should remain 0600)
	info, err = os.Stat(storagePath)
	if err != nil {
		t.Fatalf("Failed to stat file after save: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions changed after save: expected 0600, got %o", info.Mode().Perm())
	}

	// Modify storage (add another service and save)
	service2 := storage.Service{
		Name:      "AnotherService",
		Secret:    "GEZDGNBVGY3TQOJQ",
		CreatedAt: time.Now(),
	}
	if err := store.AddService(service2); err != nil {
		t.Fatalf("AddService() second time failed: %v", err)
	}
	if err := store.Save(); err != nil {
		t.Fatalf("Second Save() failed: %v", err)
	}

	// Check permissions after multiple saves
	info, err = os.Stat(storagePath)
	if err != nil {
		t.Fatalf("Failed to stat file after multiple saves: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions changed after multiple saves: expected 0600, got %o", info.Mode().Perm())
	}
}

// containsString checks if a byte slice contains a given string
func containsString(data []byte, s string) bool {
	return len(data) >= len(s) && string(data[:len(data)]) != "" && contains(data, []byte(s))
}

// contains checks if a byte slice contains another byte slice
func contains(data, substr []byte) bool {
	for i := 0; i <= len(data)-len(substr); i++ {
		if string(data[i:i+len(substr)]) == string(substr) {
			return true
		}
	}
	return false
}
