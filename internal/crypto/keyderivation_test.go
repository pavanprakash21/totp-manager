package crypto

import (
	"bytes"
	"testing"
)

// TestDeriveKey tests Argon2id key derivation
func TestDeriveKey(t *testing.T) {
	passphrase := "test-passphrase-123"
	salt := []byte("1234567890123456") // 16 bytes

	key, err := DeriveKey(passphrase, salt)
	if err != nil {
		t.Fatalf("DeriveKey() error = %v", err)
	}

	// Verify key is 32 bytes (256 bits for AES-256)
	if len(key) != 32 {
		t.Errorf("DeriveKey() key length = %d, want 32", len(key))
	}
}

// TestDeriveKey_Consistency tests that same inputs produce same key
func TestDeriveKey_Consistency(t *testing.T) {
	passphrase := "test-passphrase"
	salt := []byte("1234567890123456")

	key1, err1 := DeriveKey(passphrase, salt)
	if err1 != nil {
		t.Fatalf("First DeriveKey() error = %v", err1)
	}

	key2, err2 := DeriveKey(passphrase, salt)
	if err2 != nil {
		t.Fatalf("Second DeriveKey() error = %v", err2)
	}

	if !bytes.Equal(key1, key2) {
		t.Error("DeriveKey() produced different keys for same inputs")
	}
}

// TestDeriveKey_DifferentPassphrases tests that different passphrases produce different keys
func TestDeriveKey_DifferentPassphrases(t *testing.T) {
	salt := []byte("1234567890123456")

	key1, err1 := DeriveKey("passphrase1", salt)
	if err1 != nil {
		t.Fatalf("DeriveKey() error = %v", err1)
	}

	key2, err2 := DeriveKey("passphrase2", salt)
	if err2 != nil {
		t.Fatalf("DeriveKey() error = %v", err2)
	}

	if bytes.Equal(key1, key2) {
		t.Error("DeriveKey() produced same key for different passphrases")
	}
}

// TestDeriveKey_DifferentSalts tests that different salts produce different keys
func TestDeriveKey_DifferentSalts(t *testing.T) {
	passphrase := "test-passphrase"

	salt1 := []byte("1111111111111111")
	key1, err1 := DeriveKey(passphrase, salt1)
	if err1 != nil {
		t.Fatalf("DeriveKey() error = %v", err1)
	}

	salt2 := []byte("2222222222222222")
	key2, err2 := DeriveKey(passphrase, salt2)
	if err2 != nil {
		t.Fatalf("DeriveKey() error = %v", err2)
	}

	if bytes.Equal(key1, key2) {
		t.Error("DeriveKey() produced same key for different salts")
	}
}

// TestDeriveKey_InvalidSalt tests error handling for invalid salt
func TestDeriveKey_InvalidSalt(t *testing.T) {
	tests := []struct {
		name string
		salt []byte
	}{
		{
			name: "Empty salt",
			salt: []byte{},
		},
		{
			name: "Salt too short",
			salt: []byte("short"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeriveKey("passphrase", tt.salt)
			if err == nil {
				t.Error("DeriveKey() expected error for invalid salt, got nil")
			}
		})
	}
}

// TestGenerateSalt tests salt generation
func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error = %v", err)
	}

	// Verify salt is 16 bytes
	if len(salt) != 16 {
		t.Errorf("GenerateSalt() salt length = %d, want 16", len(salt))
	}
}

// TestGenerateSalt_Uniqueness tests that generated salts are unique
func TestGenerateSalt_Uniqueness(t *testing.T) {
	salt1, err1 := GenerateSalt()
	if err1 != nil {
		t.Fatalf("First GenerateSalt() error = %v", err1)
	}

	salt2, err2 := GenerateSalt()
	if err2 != nil {
		t.Fatalf("Second GenerateSalt() error = %v", err2)
	}

	if bytes.Equal(salt1, salt2) {
		t.Error("GenerateSalt() produced identical salts (should be cryptographically random)")
	}
}

// TestDeriveKey_EmptyPassphrase tests key derivation with empty passphrase
func TestDeriveKey_EmptyPassphrase(t *testing.T) {
	salt := []byte("1234567890123456")

	key, err := DeriveKey("", salt)
	if err != nil {
		t.Fatalf("DeriveKey() error = %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected 32-byte key even for empty passphrase, got %d", len(key))
	}
}

// TestDeriveKey_LongPassphrase tests key derivation with very long passphrase
func TestDeriveKey_LongPassphrase(t *testing.T) {
	salt := []byte("1234567890123456")
	longPassphrase := ""
	for i := 0; i < 1000; i++ {
		longPassphrase += "a"
	}

	key, err := DeriveKey(longPassphrase, salt)
	if err != nil {
		t.Fatalf("DeriveKey() error = %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected 32-byte key for long passphrase, got %d", len(key))
	}
}

// TestDeriveKey_UnicodePassphrase tests key derivation with unicode characters
func TestDeriveKey_UnicodePassphrase(t *testing.T) {
	salt := []byte("1234567890123456")
	unicodePassphrase := "パスワード🔐"

	key, err := DeriveKey(unicodePassphrase, salt)
	if err != nil {
		t.Fatalf("DeriveKey() error = %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected 32-byte key for unicode passphrase, got %d", len(key))
	}
}

// TestGenerateSalt_Multiple tests generating multiple salts
func TestGenerateSalt_Multiple(t *testing.T) {
	salts := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		salt, err := GenerateSalt()
		if err != nil {
			t.Fatalf("GenerateSalt() error = %v", err)
		}
		salts[i] = salt
	}

	// Check all salts are unique
	for i := 0; i < len(salts); i++ {
		for j := i + 1; j < len(salts); j++ {
			if bytes.Equal(salts[i], salts[j]) {
				t.Errorf("Salts %d and %d are identical", i, j)
			}
		}
	}
}

// BenchmarkDeriveKey benchmarks key derivation performance
func BenchmarkDeriveKey(b *testing.B) {
	passphrase := "test-passphrase"
	salt := []byte("1234567890123456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DeriveKey(passphrase, salt)
	}
}

// BenchmarkGenerateSalt benchmarks salt generation performance
func BenchmarkGenerateSalt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateSalt()
	}
}
