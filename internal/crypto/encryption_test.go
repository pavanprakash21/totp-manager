package crypto

import (
	"bytes"
	"testing"
)

// TestEncryptDecrypt tests encryption and decryption round-trip
func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32) // 256-bit key
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := []byte("This is a secret message for TOTP storage")

	// Encrypt
	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Verify nonce is 12 bytes
	if len(nonce) != 12 {
		t.Errorf("Nonce length = %d, want 12", len(nonce))
	}

	// Verify ciphertext is longer than plaintext (includes auth tag)
	if len(ciphertext) <= len(plaintext) {
		t.Errorf("Ciphertext length %d should be > plaintext length %d", len(ciphertext), len(plaintext))
	}

	// Decrypt
	decrypted, err := Decrypt(ciphertext, key, nonce)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	// Verify decrypted matches original
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted text doesn't match original.\nWant: %s\nGot:  %s", plaintext, decrypted)
	}
}

// TestEncrypt_DifferentNonces tests that each encryption uses a unique nonce
func TestEncrypt_DifferentNonces(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("test message")

	_, nonce1, err1 := Encrypt(plaintext, key)
	if err1 != nil {
		t.Fatalf("First Encrypt() error = %v", err1)
	}

	_, nonce2, err2 := Encrypt(plaintext, key)
	if err2 != nil {
		t.Fatalf("Second Encrypt() error = %v", err2)
	}

	if bytes.Equal(nonce1, nonce2) {
		t.Error("Encrypt() produced identical nonces (should be cryptographically random)")
	}
}

// TestDecrypt_WrongKey tests that decryption fails with wrong key
func TestDecrypt_WrongKey(t *testing.T) {
	correctKey := make([]byte, 32)
	wrongKey := make([]byte, 32)
	for i := range wrongKey {
		wrongKey[i] = byte(i + 1)
	}

	plaintext := []byte("secret message")

	// Encrypt with correct key
	ciphertext, nonce, err := Encrypt(plaintext, correctKey)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Try to decrypt with wrong key
	_, err = Decrypt(ciphertext, wrongKey, nonce)
	if err == nil {
		t.Error("Decrypt() should fail with wrong key, but succeeded")
	}
}

// TestDecrypt_TamperedCiphertext tests that decryption fails if ciphertext is tampered
func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("secret message")

	// Encrypt
	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Tamper with ciphertext
	ciphertext[0] ^= 0xFF

	// Try to decrypt tampered ciphertext
	_, err = Decrypt(ciphertext, key, nonce)
	if err == nil {
		t.Error("Decrypt() should fail with tampered ciphertext (auth tag verification)")
	}
}

// TestEncrypt_InvalidKey tests error handling for invalid key size
func TestEncrypt_InvalidKey(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{
			name: "Empty key",
			key:  []byte{},
		},
		{
			name: "Key too short",
			key:  make([]byte, 16), // 128 bits, need 256
		},
		{
			name: "Key too long",
			key:  make([]byte, 64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := Encrypt([]byte("test"), tt.key)
			if err == nil {
				t.Error("Encrypt() expected error for invalid key size, got nil")
			}
		})
	}
}

// TestDecrypt_InvalidNonce tests error handling for invalid nonce
func TestDecrypt_InvalidNonce(t *testing.T) {
	key := make([]byte, 32)
	ciphertext := []byte("dummy ciphertext")

	tests := []struct {
		name  string
		nonce []byte
	}{
		{
			name:  "Empty nonce",
			nonce: []byte{},
		},
		{
			name:  "Nonce too short",
			nonce: []byte("short"),
		},
		{
			name:  "Nonce too long",
			nonce: make([]byte, 24),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(ciphertext, key, tt.nonce)
			if err == nil {
				t.Error("Decrypt() expected error for invalid nonce, got nil")
			}
		})
	}
}

// TestEncryptDecrypt_EmptyPlaintext tests encryption of empty data
func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte{}

	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Even empty plaintext should produce ciphertext (auth tag)
	if len(ciphertext) == 0 {
		t.Error("Encrypt() produced empty ciphertext for empty plaintext")
	}

	decrypted, err := Decrypt(ciphertext, key, nonce)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted empty plaintext doesn't match")
	}
}

// TestEncryptDecrypt_LargeData tests encryption of large data
func TestEncryptDecrypt_LargeData(t *testing.T) {
	key := make([]byte, 32)

	// Create 1MB of data
	plaintext := make([]byte, 1024*1024)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}

	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key, nonce)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypted large data doesn't match original")
	}
}

// TestEncrypt_InvalidKeySize tests error handling for invalid key size
func TestEncrypt_InvalidKeySize(t *testing.T) {
	invalidKey := make([]byte, 16) // Wrong size (should be 32)
	plaintext := []byte("test")

	_, _, err := Encrypt(plaintext, invalidKey)
	if err == nil {
		t.Error("Encrypt() should fail with invalid key size")
	}
}

// TestDecrypt_InvalidKeySize tests error handling for invalid key size in decryption
func TestDecrypt_InvalidKeySize(t *testing.T) {
	invalidKey := make([]byte, 16) // Wrong size
	ciphertext := []byte("dummy")
	nonce := make([]byte, 12)

	_, err := Decrypt(ciphertext, invalidKey, nonce)
	if err == nil {
		t.Error("Decrypt() should fail with invalid key size")
	}
}

// TestDecrypt_TamperedNonce tests decryption with tampered nonce
func TestDecrypt_TamperedNonce(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("secret message")

	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Tamper with nonce
	nonce[0] ^= 0xFF

	// Try to decrypt
	_, err = Decrypt(ciphertext, key, nonce)
	if err == nil {
		t.Error("Decrypt() should fail with tampered nonce")
	}
}

// TestEncrypt_EmptyPlaintext tests encrypting empty data
func TestEncrypt_EmptyPlaintext(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte{}

	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key, nonce)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypted empty plaintext doesn't match")
	}
}

// BenchmarkEncrypt benchmarks encryption performance
func BenchmarkEncrypt(b *testing.B) {
	key := make([]byte, 32)
	plaintext := make([]byte, 1024) // 1KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = Encrypt(plaintext, key)
	}
}

// BenchmarkDecrypt benchmarks decryption performance
func BenchmarkDecrypt(b *testing.B) {
	key := make([]byte, 32)
	plaintext := make([]byte, 1024)
	ciphertext, nonce, _ := Encrypt(plaintext, key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decrypt(ciphertext, key, nonce)
	}
}
