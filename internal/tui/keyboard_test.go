package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pavanprakash21/totp-manager-go/internal/storage"
)

// TestHandleKeyPress_Navigation tests arrow and vim key navigation
func TestHandleKeyPress_Navigation(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "Service1", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Service2", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Service3", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.cursor = 0

	// Test down key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)
	if m.cursor != 1 {
		t.Errorf("Expected cursor at 1 after down, got %d", m.cursor)
	}

	// Test up key
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)
	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 after up, got %d", m.cursor)
	}

	// Test up at boundary (should stay at 0)
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)
	if m.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", m.cursor)
	}
}

// TestHandleKeyPress_VimKeys tests j/k navigation
func TestHandleKeyPress_VimKeys(t *testing.T) {
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

	// Test j key (down)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)
	if m.cursor != 1 {
		t.Errorf("Expected cursor at 1 after 'j', got %d", m.cursor)
	}

	// Test k key (up)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)
	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 after 'k', got %d", m.cursor)
	}
}

// TestHandleKeyPress_HomeEnd tests home/end navigation
func TestHandleKeyPress_HomeEnd(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "Service1", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Service2", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Service3", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.cursor = 1

	// Test G key (end)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)
	if m.cursor != 2 {
		t.Errorf("Expected cursor at 2 after 'G', got %d", m.cursor)
	}

	// Test g key (home)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)
	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 after 'g', got %d", m.cursor)
	}
}

// TestHandleKeyPress_SearchMode tests search mode toggle
func TestHandleKeyPress_SearchMode(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)

	// Test / key to enter search mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)
	if !m.searchMode {
		t.Error("Expected search mode to be true after '/'")
	}

	// Test ESC to exit search mode
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)
	if m.searchMode {
		t.Error("Expected search mode to be false after ESC")
	}
}

// TestHandleKeyPress_SearchTyping tests typing in search mode
func TestHandleKeyPress_SearchTyping(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "GitLab", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "Google", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.searchMode = true

	// Type "oog" (avoiding g/k which are now navigation keys)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)

	if m.searchQuery != "oog" {
		t.Errorf("Expected search query 'oog', got %q", m.searchQuery)
	}

	// Should filter to Google (oog matches Google)
	if len(m.filteredIndices) != 1 {
		t.Errorf("Expected 1 filtered service, got %d", len(m.filteredIndices))
	}
}

// TestHandleKeyPress_SearchBackspace tests backspace in search mode
func TestHandleKeyPress_SearchBackspace(t *testing.T) {
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

	// Test backspace
	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)

	if m.searchQuery != "gi" {
		t.Errorf("Expected search query 'gi' after backspace, got %q", m.searchQuery)
	}
}

// TestHandleKeyPress_SearchClearFilter tests Ctrl+U to clear search
func TestHandleKeyPress_SearchClearFilter(t *testing.T) {
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
	model.searchMode = true
	model.searchQuery = "git"
	model.filterServices()

	// Test Ctrl+U
	msg := tea.KeyMsg{Type: tea.KeyCtrlU}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)

	if m.searchQuery != "" {
		t.Errorf("Expected empty search query after Ctrl+U, got %q", m.searchQuery)
	}

	if len(m.filteredIndices) != 2 {
		t.Errorf("Expected all services after clearing filter, got %d", len(m.filteredIndices))
	}
}

// TestHandleKeyPress_CopyInSearchMode tests space/enter in search mode
func TestHandleKeyPress_CopyInSearchMode(t *testing.T) {
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
	model.searchMode = true
	model.cursor = 0

	// Test space key in search mode (should copy, not add to search)
	msg := tea.KeyMsg{Type: tea.KeySpace}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)

	// Should have copy status set
	if m.copyStatus == "" {
		t.Error("Expected copy status to be set after space in search mode")
	}
}

// TestHandleKeyPress_NavigationInSearchMode tests arrow key navigation in search mode
func TestHandleKeyPress_NavigationInSearchMode(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
				{Name: "GitLab", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.searchMode = true
	model.searchQuery = ""
	model.cursor = 0

	// Arrow down in search mode - should navigate
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)

	if m.cursor != 1 {
		t.Errorf("Expected cursor at 1 after down arrow, got %d", m.cursor)
	}

	// Arrow up - should navigate
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.handleKeyPress(msg)
	m = newModel.(Model)

	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 after up arrow, got %d", m.cursor)
	}
}

// TestHandleKeyPress_SearchTypingJK tests that j/k/g/G are added to search query
func TestHandleKeyPress_SearchTypingJK(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "Slack", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)
	model.searchMode = true
	model.searchQuery = "slac"

	// Type 'k' in search mode - should add to query
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)

	if m.searchQuery != "slack" {
		t.Errorf("Expected search query 'slack', got %q", m.searchQuery)
	}
}

// TestHandleKeyPress_Quit tests quit keys
func TestHandleKeyPress_Quit(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version: 1,
			Services: []storage.Service{
				{Name: "GitHub", Secret: "JBSWY3DPEHPK3PXP", CreatedAt: time.Now()},
			},
		},
	}

	model := NewModel(store)

	// Test q key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.handleKeyPress(msg)

	if cmd == nil {
		t.Error("Expected quit command")
	}

	// Test Ctrl+C
	msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd = model.handleKeyPress(msg)

	if cmd == nil {
		t.Error("Expected quit command for Ctrl+C")
	}
}

// TestHandleKeyPress_EmptyList tests navigation on empty list
func TestHandleKeyPress_EmptyList(t *testing.T) {
	store := &storage.Store{
		Storage: &storage.Storage{
			Version:  1,
			Services: []storage.Service{},
		},
	}

	model := NewModel(store)

	// Test down key on empty list
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.handleKeyPress(msg)
	m := newModel.(Model)

	if m.cursor != 0 {
		t.Errorf("Expected cursor at 0 on empty list, got %d", m.cursor)
	}
}
