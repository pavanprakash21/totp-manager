package tui

import (
	"testing"
	"time"

	"github.com/pavanprakash21/totp-manager-go/internal/storage"
)

// TestNewModel tests creating a new TUI model
func TestNewModel(t *testing.T) {
	// Create a mock store with test data
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{
					Name:      "GitHub",
					Secret:    "JBSWY3DPEHPK3PXP",
					CreatedAt: time.Now(),
				},
			},
		},
	}

	model := NewModel(store)

	if model.store == nil {
		t.Error("Model store should not be nil")
	}

	if len(model.services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(model.services))
	}

	if model.cursor != 0 {
		t.Errorf("Expected cursor at 0, got %d", model.cursor)
	}

	if model.totpCodes == nil {
		t.Error("totpCodes map should be initialized")
	}
}

// TestCalculateRemainingSeconds tests the countdown calculation
func TestCalculateRemainingSeconds(t *testing.T) {
	remaining := calculateRemainingSeconds()

	if remaining < 1 || remaining > 30 {
		t.Errorf("Expected remaining time between 1-30 seconds, got %d", remaining)
	}
}

// TestEmptyModel tests model with no services
// (T032: Test for empty list with instructions)
func TestEmptyModel(t *testing.T) {
	store := &storage.Storage{
		Version:  1,
		Services: []storage.Service{},
	}

	model := NewModel(&storage.Store{Storage: store})

	if len(model.services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(model.services))
	}

	// View should contain empty state message
	view := model.View()
	if view == "" {
		t.Error("View should not be empty for empty model")
	}
}
