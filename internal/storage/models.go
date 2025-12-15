package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/pavanprakash21/totp-manager-go/internal/totp"
)

// Service represents a single TOTP service configuration
type Service struct {
	// Name is the user-visible identifier (e.g., "GitHub", "AWS")
	Name string `json:"name"`

	// Identifier is an optional additional identifier (e.g., email, username)
	Identifier string `json:"identifier,omitempty"`

	// Secret is the Base32-encoded shared secret
	Secret string `json:"secret"`

	// CreatedAt is the timestamp when service was added
	CreatedAt time.Time `json:"created_at"`

	// LastUsed is updated when TOTP code is copied
	LastUsed *time.Time `json:"last_used,omitempty"`
}

// Validate validates the Service struct
func (s *Service) Validate() error {
	// Validate name
	if err := ValidateServiceName(s.Name); err != nil {
		return err
	}

	// Validate secret
	if err := totp.ValidateSecret(s.Secret); err != nil {
		return fmt.Errorf("invalid secret: %w", err)
	}

	return nil
}

// Storage encapsulates encrypted service data and metadata
type Storage struct {
	// Version for future format migrations (current: 1)
	Version int `json:"version"`

	// Services is the list of configured TOTP services
	Services []Service `json:"services"`

	// Salt for Argon2id key derivation (stored separately in file)
	Salt []byte `json:"-"`

	// Nonce for AES-GCM encryption (stored separately in file)
	Nonce []byte `json:"-"`
}

// AddService adds a new service to storage
func (s *Storage) AddService(service Service) error {
	// Validate service
	if err := service.Validate(); err != nil {
		return err
	}

	// Check for duplicate name (case-insensitive)
	for _, existing := range s.Services {
		if strings.EqualFold(existing.Name, service.Name) {
			return fmt.Errorf("service '%s' already exists", service.Name)
		}
	}

	// Add service
	s.Services = append(s.Services, service)
	return nil
}

// GetService retrieves a service by name (case-insensitive)
func (s *Storage) GetService(name string) (*Service, error) {
	for i := range s.Services {
		if strings.EqualFold(s.Services[i].Name, name) {
			return &s.Services[i], nil
		}
	}
	return nil, fmt.Errorf("service '%s' not found", name)
}

// UpdateLastUsed updates the LastUsed timestamp for a service
func (s *Storage) UpdateLastUsed(name string) error {
	for i := range s.Services {
		if strings.EqualFold(s.Services[i].Name, name) {
			now := time.Now()
			s.Services[i].LastUsed = &now
			return nil
		}
	}
	return fmt.Errorf("service '%s' not found", name)
}

// ValidateServiceName validates a service name
func ValidateServiceName(name string) error {
	// Trim whitespace for validation
	trimmed := strings.TrimSpace(name)

	// Check empty
	if trimmed == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Check length (1-50 characters)
	if len(trimmed) > 50 {
		return fmt.Errorf("service name too long: max 50 characters, got %d", len(trimmed))
	}

	// Check for control characters and path separators
	for _, c := range trimmed {
		if c < 32 || c == 127 { // Control characters
			return fmt.Errorf("service name contains control character")
		}
		if c == '/' || c == '\\' { // Path separators
			return fmt.Errorf("service name cannot contain path separators")
		}
	}

	return nil
}
