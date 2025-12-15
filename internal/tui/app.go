package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pavanprakash21/totp-manager-go/internal/storage"
	"github.com/pavanprakash21/totp-manager-go/internal/totp"
)

// Model represents the Bubbletea TUI model
// (T039: Bubbletea Model struct with TUIState fields)
type Model struct {
	store           *storage.Store
	services        []storage.Service
	filteredIndices []int // indices of filtered services
	cursor          int
	totpCodes       map[string]string // service name -> current TOTP code
	remainingTime   int               // seconds remaining until refresh
	lastUpdate      time.Time
	copyStatus      string // Status message for clipboard operations
	copyStatusTime  time.Time
	width           int
	height          int
	searchMode      bool   // whether in search mode
	searchQuery     string // current search query
}

// tickMsg is sent every second for countdown updates
type tickMsg time.Time

// refreshMsg is sent when TOTP codes should refresh
type refreshMsg time.Time

// NewModel creates a new TUI model with storage
func NewModel(store *storage.Store) Model {
	// Initialize with all services visible
	filteredIndices := make([]int, len(store.Services))
	for i := range filteredIndices {
		filteredIndices[i] = i
	}

	return Model{
		store:           store,
		services:        store.Services,
		filteredIndices: filteredIndices,
		totpCodes:       make(map[string]string),
		lastUpdate:      time.Now(),
		remainingTime:   calculateRemainingSeconds(),
		searchMode:      false,
		searchQuery:     "",
	}
}

// calculateRemainingSeconds calculates seconds until next 30s interval
func calculateRemainingSeconds() int {
	now := time.Now().Unix()
	return 30 - int(now%30)
}

// Init implements tea.Model interface
// (T040: Initialize services from storage)
func (m Model) Init() tea.Cmd {
	// Generate initial TOTP codes for all services
	m.generateAllCodes()

	// Start ticker for countdown updates
	return tea.Batch(
		tickCmd(),
		tea.WindowSize(),
	)
}

// generateAllCodes generates TOTP codes for all services
func (m *Model) generateAllCodes() {
	now := time.Now()
	for i := range m.services {
		service := &m.services[i]
		code, err := totp.GenerateCode(service.Secret, now)
		if err != nil {
			m.totpCodes[service.Name] = "ERROR"
			continue
		}
		m.totpCodes[service.Name] = code
	}
	m.remainingTime = calculateRemainingSeconds()
}

// filterServices performs fuzzy search on services
func (m *Model) filterServices() {
	if m.searchQuery == "" {
		// No search query, show all services
		m.filteredIndices = make([]int, len(m.services))
		for i := range m.filteredIndices {
			m.filteredIndices[i] = i
		}
		m.cursor = 0
		return
	}

	// Fuzzy search: match query characters in order (case-insensitive)
	query := strings.ToLower(m.searchQuery)
	m.filteredIndices = []int{}

	for i, service := range m.services {
		// Search in both name and identifier
		searchText := strings.ToLower(service.Name + " " + service.Identifier)
		if fuzzyMatch(searchText, query) {
			m.filteredIndices = append(m.filteredIndices, i)
		}
	}

	// Reset cursor to first result
	if m.cursor >= len(m.filteredIndices) {
		m.cursor = 0
	}
}

// fuzzyMatch checks if all characters in query appear in text in order
func fuzzyMatch(text, query string) bool {
	queryIdx := 0
	for i := 0; i < len(text) && queryIdx < len(query); i++ {
		if text[i] == query[queryIdx] {
			queryIdx++
		}
	}
	return queryIdx == len(query)
}

// tickCmd returns a command that ticks every second
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update implements tea.Model interface
// (T043: Update method with keyboard message handling)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		// T049: Update countdown every second
		m.remainingTime--
		if m.remainingTime <= 0 {
			// T050: Refresh TOTP codes every 30 seconds
			m.remainingTime = 30
			m.generateAllCodes()
		}

		// Clear copy status after 2 seconds
		if !m.copyStatusTime.IsZero() && time.Since(m.copyStatusTime) > 2*time.Second {
			m.copyStatus = ""
			m.copyStatusTime = time.Time{}
		}

		return m, tickCmd()

	case refreshMsg:
		m.generateAllCodes()
		return m, nil
	}

	return m, nil
}
