package cli

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

// ChangePassphraseCommand handles changing the storage passphrase
func ChangePassphraseCommand(args []string) int {
	// Create app and initialize with current passphrase
	app, err := NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Load existing storage (prompts for current passphrase)
	fmt.Println("Changing storage passphrase...")
	if err := app.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Prompt for new passphrase with confirmation
	newPassphrase, err := promptNewPassphrase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Change passphrase (re-encrypts the file)
	if err := app.store.ChangePassphrase(newPassphrase); err != nil {
		fmt.Fprintf(os.Stderr, "Error changing passphrase: %v\n", err)
		return 1
	}

	fmt.Println("✓ Passphrase changed successfully!")
	fmt.Println("  The storage file has been re-encrypted with the new passphrase.")
	return 0
}

// promptNewPassphrase prompts for a new passphrase with confirmation
func promptNewPassphrase() (string, error) {
	// Get new passphrase
	fmt.Print("Enter new passphrase: ")
	newPass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read new passphrase: %w", err)
	}

	if len(newPass) == 0 {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	// Confirm new passphrase
	fmt.Print("Confirm new passphrase: ")
	confirmPass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read confirmation: %w", err)
	}

	if string(newPass) != string(confirmPass) {
		return "", fmt.Errorf("passphrases do not match")
	}

	return string(newPass), nil
}
