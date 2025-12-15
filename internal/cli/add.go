package cli

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pavanprakash21/totp-manager-go/internal/storage"
	"github.com/pavanprakash21/totp-manager-go/internal/totp"
)

// AddCommand handles adding a new TOTP service
// (T059-T066: CLI add command implementation)
func AddCommand(args []string) int {
	// T059: Parse CLI flags for --name and --secret
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	name := fs.String("name", "", "Service name (required)")
	identifier := fs.String("identifier", "", "Optional identifier (e.g., email, username)")
	secret := fs.String("secret", "", "Base32 TOTP secret (required)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return 1 // T065: Exit code 1 for errors
	}

	// Validate required flags
	if *name == "" {
		fmt.Fprintln(os.Stderr, "Error: --name is required")
		fmt.Fprintln(os.Stderr, "Usage: totp add --name SERVICE_NAME --secret BASE32_SECRET")
		return 1
	}

	if *secret == "" {
		fmt.Fprintln(os.Stderr, "Error: --secret is required")
		fmt.Fprintln(os.Stderr, "Usage: totp add --name SERVICE_NAME --secret BASE32_SECRET")
		return 1
	}

	// T062: Validate Base32 secret
	if err := totp.ValidateSecret(*secret); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid TOTP secret: %v\n", err)
		fmt.Fprintln(os.Stderr, "Secret must be valid Base32 (A-Z, 2-7) and at least 16 characters")
		return 1
	}

	// Initialize app and load storage
	app, err := NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// T060: Load storage (prompts for passphrase if exists, creates if not)
	if err := app.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// T061: Check for duplicate name
	if _, err := app.store.GetService(*name); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Service '%s' already exists\n", *name)
		fmt.Fprintln(os.Stderr, "Use a different name or remove the existing service first")
		return 1
	}

	// Create new service
	service := storage.Service{
		Name:       *name,
		Identifier: *identifier,
		Secret:     *secret,
		CreatedAt:  time.Now(),
	}

	// Add service to storage
	if err := app.store.AddService(service); err != nil {
		fmt.Fprintf(os.Stderr, "Error adding service: %v\n", err)
		return 1
	}

	// T063: Save storage (re-encrypts with updated data)
	if err := app.store.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving storage: %v\n", err)
		return 1
	}

	// T064: Success message to stdout
	fmt.Printf("✓ Service '%s' added successfully\n", *name)
	fmt.Println("✓ Storage updated and encrypted")

	return 0 // T065: Exit code 0 for success
}
