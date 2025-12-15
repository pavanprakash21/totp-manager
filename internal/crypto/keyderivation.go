package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters (memory-hard KDF)
	saltLength = 16        // 16 bytes (128 bits)
	keyLength  = 32        // 32 bytes (256 bits for AES-256)
	time       = 4         // Number of iterations
	memory     = 64 * 1024 // 64 MB memory
	threads    = 4         // Number of parallel threads
)

// DeriveKey derives a 256-bit encryption key from a passphrase using Argon2id
// Parameters: 64MB memory, 4 iterations, 4 threads
func DeriveKey(passphrase string, salt []byte) ([]byte, error) {
	// Validate salt length
	if len(salt) < saltLength {
		return nil, fmt.Errorf("salt too short: need %d bytes, got %d", saltLength, len(salt))
	}

	// Derive key using Argon2id (memory-hard KDF resistant to GPU attacks)
	key := argon2.IDKey(
		[]byte(passphrase),
		salt,
		time,
		memory,
		threads,
		keyLength,
	)

	return key, nil
}

// GenerateSalt generates a cryptographically secure random salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %w", err)
	}
	return salt, nil
}
