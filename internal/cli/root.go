package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/pavanprakash21/totp-manager-go/internal/storage"
	"golang.org/x/term"
)

const maxPassphraseAttempts = 3

// App represents the CLI application
type App struct {
	store       *storage.Store
	storagePath string
}

// NewApp creates a new CLI application instance
func NewApp() (*App, error) {
	path, err := storage.GetDefaultStoragePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage path: %w", err)
	}
	return &App{
		storagePath: path,
	}, nil
}

// Initialize loads or creates the encrypted storage
// (T026, T027, T028: Passphrase prompt, storage init, validation)
func (a *App) Initialize() error {
	// Check if storage file exists
	if _, err := os.Stat(a.storagePath); os.IsNotExist(err) {
		// First time setup: create new storage
		return a.createNewStorage()
	}

	// Load existing storage with passphrase attempts
	return a.loadExistingStorage()
}

// createNewStorage creates a new encrypted storage with passphrase confirmation
// (T026: Passphrase prompt with confirmation)
func (a *App) createNewStorage() error {
	fmt.Println("Welcome to TOTP Manager!")
	fmt.Println("No storage found. Let's create a new one.")
	fmt.Println()

	// Get new passphrase with confirmation
	passphrase, err := a.promptNewPassphrase()
	if err != nil {
		return fmt.Errorf("passphrase setup failed: %w", err)
	}

	// Create storage (T027: Storage initialization)
	store, err := storage.Create(a.storagePath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// Save storage to disk (creates file with 0600 permissions - T031)
	if err := store.Save(); err != nil {
		return fmt.Errorf("failed to save storage: %w", err)
	}

	a.store = store

	// Log success (T030: Security event logging)
	fmt.Println("✓ Storage created successfully")
	fmt.Printf("✓ Storage location: %s\n", a.storagePath)
	fmt.Printf("✓ File permissions: 0600 (owner read/write only)\n")
	fmt.Println()

	return nil
}

// loadExistingStorage loads existing storage with 3-attempt limit
// (T028: Passphrase validation with 3-attempt limit)
func (a *App) loadExistingStorage() error {
	var lastErr error

	// Allow up to 3 attempts
	for attempt := 1; attempt <= maxPassphraseAttempts; attempt++ {
		passphrase, err := a.promptPassphrase(attempt)
		if err != nil {
			return fmt.Errorf("passphrase input failed: %w", err)
		}

		// Try to load storage
		store, err := storage.Load(a.storagePath, passphrase)
		if err == nil {
			a.store = store
			return nil
		}

		lastErr = err

		// T029: Error handling with clear messages
		if attempt < maxPassphraseAttempts {
			fmt.Printf("✗ Incorrect passphrase (attempt %d/%d)\n", attempt, maxPassphraseAttempts)
			fmt.Println()
		}
	}

	// T029: Failed after 3 attempts
	fmt.Printf("✗ Failed to unlock storage after %d attempts\n", maxPassphraseAttempts)
	fmt.Println("For security reasons, the application will now exit.")
	fmt.Println()
	// T030: Log security event (no passphrase logged)
	fmt.Fprintf(os.Stderr, "SECURITY: Failed authentication attempts for storage: %s\n", a.storagePath)

	return fmt.Errorf("authentication failed: %w", lastErr)
}

// promptNewPassphrase prompts for a new passphrase with confirmation
func (a *App) promptNewPassphrase() (string, error) {
	fmt.Print("Enter new passphrase: ")
	passphrase1, err := readPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase: %w", err)
	}
	fmt.Println()

	// Validate passphrase strength
	if len(passphrase1) < 8 {
		return "", fmt.Errorf("passphrase must be at least 8 characters")
	}

	fmt.Print("Confirm passphrase: ")
	passphrase2, err := readPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read confirmation: %w", err)
	}
	fmt.Println()

	if passphrase1 != passphrase2 {
		return "", fmt.Errorf("passphrases do not match")
	}

	return passphrase1, nil
}

// promptPassphrase prompts for passphrase (for existing storage)
func (a *App) promptPassphrase(attempt int) (string, error) {
	if attempt == 1 {
		fmt.Println("Enter passphrase to unlock storage:")
	}

	fmt.Print("Passphrase: ")
	passphrase, err := readPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase: %w", err)
	}
	fmt.Println()

	return passphrase, nil
}

// readPassword reads a password from stdin without echoing
func readPassword() (string, error) {
	// Try to read from terminal (supports masking)
	if term.IsTerminal(int(syscall.Stdin)) {
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(bytePassword)), nil
	}

	// Fallback for non-terminal input (e.g., tests)
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}

// GetStore returns the initialized storage store
func (a *App) GetStore() *storage.Store {
	return a.store
}
