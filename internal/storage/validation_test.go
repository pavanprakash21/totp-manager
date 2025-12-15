package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestStore_AddServiceDuplicate tests adding duplicate service
func TestStore_AddServiceDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test.enc")
	passphrase := "test-pass"

	store, err := Create(storePath, passphrase)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	service := Service{
		Name:      "GitHub",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}

	// Add first time
	err = store.AddService(service)
	if err != nil {
		t.Fatalf("First AddService() error = %v", err)
	}

	// Try to add duplicate
	err = store.AddService(service)
	if err == nil {
		t.Error("Expected error adding duplicate service")
	}
}

// TestStore_Load_CorruptedFile tests loading corrupted file
func TestStore_Load_CorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "corrupted.enc")

	// Write invalid data
	err := os.WriteFile(storePath, []byte("corrupted data"), 0600)
	if err != nil {
		t.Fatalf("Failed to write corrupted file: %v", err)
	}

	// Try to load
	_, err = Load(storePath, "password")
	if err == nil {
		t.Error("Expected error loading corrupted file")
	}
}

// TestStore_Load_WrongPassphrase tests loading with wrong passphrase
func TestStore_Load_WrongPassphrase(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test.enc")
	correctPass := "correct"
	wrongPass := "wrong"

	// Create and save with correct password
	store, err := Create(storePath, correctPass)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Try to load with wrong password
	_, err = Load(storePath, wrongPass)
	if err == nil {
		t.Error("Expected error loading with wrong passphrase")
	}
}

// TestStore_Save_InvalidPath tests saving to invalid path
func TestStore_Save_InvalidPath(t *testing.T) {
	store := &Store{
		Storage: &Storage{
			Version:  1,
			Services: []Service{},
		},
		passphrase: "test",
		path:       "/invalid/path/that/does/not/exist/file.enc",
	}

	err := store.Save()
	if err == nil {
		t.Error("Expected error saving to invalid path")
	}
}

// TestStore_ChangePassphrase tests changing passphrase
func TestStore_ChangePassphrase(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "test.enc")
	oldPass := "old-password"
	newPass := "new-password"

	// Create with old password
	store, err := Create(storePath, oldPass)
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

	err = store.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Change passphrase
	err = store.ChangePassphrase(newPass)
	if err != nil {
		t.Fatalf("ChangePassphrase() error = %v", err)
	}

	// Load with new password
	loaded, err := Load(storePath, newPass)
	if err != nil {
		t.Fatalf("Load() with new password error = %v", err)
	}

	if len(loaded.Services) != 1 {
		t.Errorf("Expected 1 service after password change, got %d", len(loaded.Services))
	}

	// Old password should not work
	_, err = Load(storePath, oldPass)
	if err == nil {
		t.Error("Old password should not work after change")
	}
}
