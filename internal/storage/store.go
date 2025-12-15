package storage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pavanprakash21/totp-manager-go/internal/crypto"
)

// Store manages encrypted TOTP service storage
type Store struct {
	path       string
	passphrase string
	*Storage
}

// Create creates a new encrypted storage file
func Create(path, passphrase string) (*Store, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate salt for key derivation
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	store := &Store{
		path:       path,
		passphrase: passphrase,
		Storage: &Storage{
			Version:  1,
			Services: []Service{},
			Salt:     salt,
		},
	}

	return store, nil
}

// Load loads and decrypts an existing storage file
func Load(path, passphrase string) (*Store, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	// Parse file format:
	// [4 bytes: Version]
	// [16 bytes: Salt]
	// [12 bytes: Nonce]
	// [N bytes: Encrypted JSON + Auth Tag]

	if len(data) < 4+16+12+16 {
		return nil, fmt.Errorf("invalid storage file: too short")
	}

	// Read version
	version := binary.LittleEndian.Uint32(data[0:4])
	if version != 1 {
		return nil, fmt.Errorf("unsupported storage version: %d", version)
	}

	// Read salt and nonce
	salt := data[4:20]
	nonce := data[20:32]
	ciphertext := data[32:]

	// Derive key from passphrase
	key, err := crypto.DeriveKey(passphrase, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Decrypt
	plaintext, err := crypto.Decrypt(ciphertext, key, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt storage (wrong passphrase?): %w", err)
	}

	// Unmarshal JSON
	var storage Storage
	if err := json.Unmarshal(plaintext, &storage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal storage: %w", err)
	}

	storage.Salt = salt
	storage.Nonce = nonce

	store := &Store{
		path:       path,
		passphrase: passphrase,
		Storage:    &storage,
	}

	return store, nil
}

// Save encrypts and saves storage to disk (atomic write)
func (s *Store) Save() error {
	// Derive key from passphrase
	key, err := crypto.DeriveKey(s.passphrase, s.Salt)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	// Marshal storage to JSON
	jsonData, err := json.Marshal(s.Storage)
	if err != nil {
		return fmt.Errorf("failed to marshal storage: %w", err)
	}

	// Encrypt
	ciphertext, nonce, err := crypto.Encrypt(jsonData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt storage: %w", err)
	}

	// Build file content
	// [4 bytes: Version] [16 bytes: Salt] [12 bytes: Nonce] [N bytes: Ciphertext + Auth Tag]
	fileData := make([]byte, 4+16+12+len(ciphertext))
	binary.LittleEndian.PutUint32(fileData[0:4], uint32(s.Version))
	copy(fileData[4:20], s.Salt)
	copy(fileData[20:32], nonce)
	copy(fileData[32:], ciphertext)

	// Atomic write: write to temp file, then rename
	tmpPath := s.path + ".tmp"

	// Write temp file with 0600 permissions
	if err := os.WriteFile(tmpPath, fileData, 0600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Rename temp file to actual file (atomic on Unix)
	if err := os.Rename(tmpPath, s.path); err != nil {
		os.Remove(tmpPath) // Clean up temp file on error
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	// Update nonce in memory
	s.Nonce = nonce

	return nil
}

// ChangePassphrase re-encrypts storage with a new passphrase
func (s *Store) ChangePassphrase(newPassphrase string) error {
	// Generate new salt
	newSalt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate new salt: %w", err)
	}

	// Update passphrase and salt
	s.passphrase = newPassphrase
	s.Salt = newSalt

	// Save with new passphrase (atomic)
	return s.Save()
}

// GetDefaultStoragePath returns the default storage path
func GetDefaultStoragePath() (string, error) {
	// Use XDG_CONFIG_HOME or ~/.config
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	storageDir := filepath.Join(configDir, "totp-manager")
	storagePath := filepath.Join(storageDir, "secrets.enc")

	return storagePath, nil
}
