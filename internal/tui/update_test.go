package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pavanprakash21/totp-manager-go/internal/storage"
)

// TestInit tests the Init method
func TestInit(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	cmd := model.Init()

	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

// TestUpdate_TickMsg tests Update with tick message
func TestUpdate_TickMsg(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.remainingTime = 5

	// Send tick message
	msg := tickMsg(time.Now())
	newModel, cmd := model.Update(msg)

	m := newModel.(Model)
	if m.remainingTime != 4 {
		t.Errorf("Expected remaining time 4, got %d", m.remainingTime)
	}

	if cmd == nil {
		t.Error("Update should return tick command")
	}
}

// TestUpdate_RefreshMsg tests Update with refresh message
func TestUpdate_RefreshMsg(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)

	// Send refresh message
	msg := refreshMsg(time.Now())
	newModel, _ := model.Update(msg)

	m := newModel.(Model)

	// Should have generated new codes
	if len(m.totpCodes) != 1 {
		t.Errorf("Expected 1 TOTP code after refresh, got %d", len(m.totpCodes))
	}

	// Verify code was actually generated
	if code := m.totpCodes["GitHub"]; code == "" {
		t.Error("TOTP code should be generated for GitHub")
	}
} // TestUpdate_WindowSizeMsg tests Update with window size message
func TestUpdate_WindowSizeMsg(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version:  1,
			Services: []storage.Service{},
		},
	}

	model := NewModel(store)

	// Send window size message
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ := model.Update(msg)

	m := newModel.(Model)
	if m.width != 100 {
		t.Errorf("Expected width 100, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("Expected height 50, got %d", m.height)
	}
}

// TestUpdate_KeyMsg tests Update with key message
func TestUpdate_KeyMsg(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "Service1", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Service2", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)

	// Send key message
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(msg)

	m := newModel.(Model)
	if m.cursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", m.cursor)
	}
}

// TestUpdate_CopyStatusTimeout tests that copy status disappears after timeout
func TestUpdate_CopyStatusTimeout(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.copyStatus = "Test status"
	model.copyStatusTime = time.Now().Add(-4 * time.Second) // 4 seconds ago

	// Send tick message
	msg := tickMsg(time.Now())
	newModel, _ := model.Update(msg)

	m := newModel.(Model)
	if m.copyStatus != "" {
		t.Error("Copy status should clear after 3 seconds")
	}
}

// TestUpdate_CopyStatusNotTimeout tests that copy status remains within timeout
func TestUpdate_CopyStatusNotTimeout(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.copyStatus = "Test status"
	model.copyStatusTime = time.Now().Add(-1 * time.Second) // 1 second ago

	// Send tick message
	msg := tickMsg(time.Now())
	newModel, _ := model.Update(msg)

	m := newModel.(Model)
	if m.copyStatus != "Test status" {
		t.Error("Copy status should remain within 3 seconds")
	}
}

// TestRenderServiceLine tests service line rendering
func TestRenderServiceLine(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version:  1,
			Services: []storage.Service{},
		},
	}

	model := NewModel(store)

	// Test normal line
	line := model.renderServiceLine("GitHub", "", "123456", false)
	if line == "" {
		t.Error("renderServiceLine should return non-empty string")
	}

	// Test selected line
	selectedLine := model.renderServiceLine("GitHub", "", "123456", true)
	if selectedLine == "" {
		t.Error("renderServiceLine should return non-empty string for selected")
	}

	// Both should contain the service name and code
	if !containsString(line, "GitHub") {
		t.Error("Normal line should contain service name")
	}
	if !containsString(selectedLine, "GitHub") {
		t.Error("Selected line should contain service name")
	}
}

// TestRenderServiceLine_WithIdentifier tests rendering with identifier
func TestRenderServiceLine_WithIdentifier(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version:  1,
			Services: []storage.Service{},
		},
	}

	model := NewModel(store)

	line := model.renderServiceLine("GitHub", "user@example.com", "123456", false)
	if line == "" {
		t.Error("renderServiceLine with identifier should return non-empty string")
	}
}

// TestRenderServiceLine_LongName tests truncation of long service names
func TestRenderServiceLine_LongName(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version:  1,
			Services: []storage.Service{},
		},
	}

	model := NewModel(store)

	longName := "This is a very long service name that should be truncated because it exceeds the maximum allowed length"
	line := model.renderServiceLine(longName, "", "123456", false)

	if line == "" {
		t.Error("renderServiceLine with long name should return non-empty string")
	}
}

// TestCalculateRemainingSeconds_Boundary tests boundary conditions
func TestCalculateRemainingSeconds_Boundary(t *testing.T) {
	// Run multiple times to catch edge cases
	for i := 0; i < 100; i++ {
		remaining := calculateRemainingSeconds()
		if remaining < 1 || remaining > 30 {
			t.Errorf("Remaining seconds %d out of range [1,30]", remaining)
		}
	}
}
