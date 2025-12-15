package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

const (
	nonceSize = 12 // 12 bytes for GCM (96 bits)
)

// Encrypt encrypts plaintext using AES-256-GCM with authenticated encryption
// Returns ciphertext (including auth tag), nonce, and error
func Encrypt(plaintext, key []byte) (ciphertext, nonce []byte, err error) {
	// Validate key size (must be 32 bytes for AES-256)
	if len(key) != 32 {
		return nil, nil, fmt.Errorf("invalid key size: need 32 bytes for AES-256, got %d", len(key))
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce (12 bytes for GCM)
	nonce = make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate
	// GCM automatically appends 16-byte authentication tag
	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM and verifies authentication tag
// Returns plaintext and error (error if authentication fails or decryption fails)
func Decrypt(ciphertext, key, nonce []byte) (plaintext []byte, err error) {
	// Validate key size
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: need 32 bytes for AES-256, got %d", len(key))
	}

	// Validate nonce size
	if len(nonce) != nonceSize {
		return nil, fmt.Errorf("invalid nonce size: need %d bytes, got %d", nonceSize, len(nonce))
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt and verify authentication tag
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong key or tampered data): %w", err)
	}

	return plaintext, nil
}
