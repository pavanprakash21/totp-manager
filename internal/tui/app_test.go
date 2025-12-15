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

// TestGenerateAllCodes tests TOTP code generation
func TestGenerateAllCodes(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{
					Name:      "GitHub",
					Secret:    "JBSWY3DPEHPK3PXP",
					CreatedAt: time.Now(),
				},
				{
					Name:      "Google",
					Secret:    "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
					CreatedAt: time.Now(),
				},
			},
		},
	}

	model := NewModel(store)
	model.generateAllCodes()

	if len(model.totpCodes) != 2 {
		t.Errorf("Expected 2 TOTP codes, got %d", len(model.totpCodes))
	}

	githubCode := model.totpCodes["GitHub"]
	if len(githubCode) != 6 {
		t.Errorf("Expected 6-digit code, got %s", githubCode)
	}
}

// TestGenerateAllCodes_InvalidSecret tests error handling for invalid secrets
func TestGenerateAllCodes_InvalidSecret(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{
					Name:      "Invalid",
					Secret:    "INVALID!!!",
					CreatedAt: time.Now(),
				},
			},
		},
	}

	model := NewModel(store)
	model.generateAllCodes()

	code := model.totpCodes["Invalid"]
	if code != "ERROR" {
		t.Errorf("Expected 'ERROR' for invalid secret, got %s", code)
	}
}

// TestFilterServices tests fuzzy search filtering
func TestFilterServices(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "GitLab", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Google", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "AWS", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)

	// Test empty query - should show all
	model.searchQuery = ""
	model.filterServices()
	if len(model.filteredIndices) != 4 {
		t.Errorf("Expected 4 services with empty query, got %d", len(model.filteredIndices))
	}

	// Test "git" - should match GitHub and GitLab
	model.searchQuery = "git"
	model.filterServices()
	if len(model.filteredIndices) != 2 {
		t.Errorf("Expected 2 services matching 'git', got %d", len(model.filteredIndices))
	}

	// Test "goog" - should match Google
	model.searchQuery = "goog"
	model.filterServices()
	if len(model.filteredIndices) != 1 {
		t.Errorf("Expected 1 service matching 'goog', got %d", len(model.filteredIndices))
	}

	// Test "xyz" - should match nothing
	model.searchQuery = "xyz"
	model.filterServices()
	if len(model.filteredIndices) != 0 {
		t.Errorf("Expected 0 services matching 'xyz', got %d", len(model.filteredIndices))
	}
}

// TestFilterServices_WithIdentifier tests filtering with identifier field
func TestFilterServices_WithIdentifier(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Identifier: "user@example.com", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "GitHub", Identifier: "admin@example.com", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)

	// Test "user" - should match first GitHub
	model.searchQuery = "user"
	model.filterServices()
	if len(model.filteredIndices) != 1 {
		t.Errorf("Expected 1 service matching 'user', got %d", len(model.filteredIndices))
	}

	// Test "example" - should match both
	model.searchQuery = "example"
	model.filterServices()
	if len(model.filteredIndices) != 2 {
		t.Errorf("Expected 2 services matching 'example', got %d", len(model.filteredIndices))
	}
}

// TestFuzzyMatch tests the fuzzy matching algorithm
func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		query    string
		expected bool
	}{
		{name: "Exact match", text: "github", query: "github", expected: true},
		{name: "Fuzzy match", text: "github", query: "gh", expected: true},
		{name: "Fuzzy match spread", text: "github", query: "gtb", expected: true},
		{name: "No match", text: "github", query: "xyz", expected: false},
		{name: "Empty query", text: "github", query: "", expected: true},
		{name: "Query longer than text", text: "git", query: "github", expected: false},
		{name: "Case sensitive", text: "GitHub", query: "github", expected: false},
		{name: "Substring", text: "github.com", query: "hub", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fuzzyMatch(tt.text, tt.query)
			if result != tt.expected {
				t.Errorf("fuzzyMatch(%q, %q) = %v, want %v", tt.text, tt.query, result, tt.expected)
			}
		})
	}
}

// TestModelView tests the View rendering
func TestModelView(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.generateAllCodes()

	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	// View should contain service name
	if !containsString(view, "GitHub") {
		t.Error("View should contain service name 'GitHub'")
	}
}

// TestModelView_SearchMode tests View in search mode
func TestModelView_SearchMode(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.searchMode = true
	model.searchQuery = "git"

	view := model.View()
	if !containsString(view, "Search:") {
		t.Error("View should contain 'Search:' in search mode")
	}
}

// TestModelView_FilteredMode tests View with filter active
func TestModelView_FilteredMode(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Google", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.searchMode = false
	model.searchQuery = "git"
	model.filterServices()

	view := model.View()
	if !containsString(view, "Filter:") {
		t.Error("View should contain 'Filter:' when filtered but not in search mode")
	}
}

// TestModelView_NoResults tests View with no search results
func TestModelView_NoResults(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.searchQuery = "nonexistent"
	model.filterServices()

	view := model.View()
	if !containsString(view, "No matching services") {
		t.Error("View should contain 'No matching services' when no results")
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
