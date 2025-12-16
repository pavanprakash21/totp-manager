package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pavanprakash21/totp-manager-go/internal/clipboard"
)

// handleKeyPress handles all keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Search mode handling
	if m.searchMode {
		switch msg.Type {
		case tea.KeyEsc:
			// Exit search mode but keep the filtered results
			m.searchMode = false
			return m, nil

		case tea.KeyBackspace:
			// Remove last character from search query
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filterServices()
			}
			return m, nil

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyCtrlU:
			// Clear search and show all services (vim-style clear line)
			m.searchQuery = ""
			m.filterServices()
			return m, nil

		case tea.KeyUp:
			// Allow navigation in search mode
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.viewportOffset {
					m.viewportOffset = m.cursor
				}
			}
			return m, nil

		case tea.KeyDown:
			// Allow navigation in search mode
			if m.cursor < len(m.filteredIndices)-1 {
				m.cursor++
				maxVisibleItems := (m.height - 9) / 3
				if maxVisibleItems < 1 {
					maxVisibleItems = 1
				}
				if m.cursor >= m.viewportOffset+maxVisibleItems {
					m.viewportOffset = m.cursor - maxVisibleItems + 1
				}
			}
			return m, nil

		case tea.KeySpace, tea.KeyEnter:
			// Allow copying in search mode
			if len(m.filteredIndices) > 0 && m.cursor < len(m.filteredIndices) {
				serviceIdx := m.filteredIndices[m.cursor]
				service := m.services[serviceIdx]
				code := m.totpCodes[service.Name]
				if code != "" {
					if err := clipboard.Copy(code); err != nil {
						m.copyStatus = "⚠ Clipboard unavailable. Code: " + code
					} else {
						m.copyStatus = "✓ Copied to clipboard"
					}
					m.copyStatusTime = time.Now()
					m.store.UpdateLastUsed(service.Name)
					_ = m.store.Save()
				}
			}
			return m, nil

		case tea.KeyRunes:
			// All typed characters are search input in search mode
			// This includes j, k, g, G - only arrow keys navigate
			m.searchQuery += string(msg.Runes)
			m.filterServices()
			return m, nil
		}
		return m, nil
	}

	// Normal mode handling
	switch msg.String() {
	// Enter search mode with '/'
	case "/":
		m.searchMode = true
		m.searchQuery = ""
		return m, nil

	// Clear search filter and show all services
	case "ctrl+u":
		m.searchQuery = ""
		m.filterServices()
		return m, nil

	// T051: Exit on 'q' or ESC
	case "q", "esc", "ctrl+c":
		return m, tea.Quit

	// T044: Arrow key navigation (↑↓)
	case "up", "k": // T045: Vim key 'k' for up
		if m.cursor > 0 {
			m.cursor--
			// Scroll viewport up if cursor goes above visible area
			if m.cursor < m.viewportOffset {
				m.viewportOffset = m.cursor
			}
		}

	case "down", "j": // T045: Vim key 'j' for down
		if m.cursor < len(m.filteredIndices)-1 {
			m.cursor++
			// Scroll viewport down if cursor goes below visible area
			maxVisibleItems := (m.height - 9) / 3
			if maxVisibleItems < 1 {
				maxVisibleItems = 1
			}
			if m.cursor >= m.viewportOffset+maxVisibleItems {
				m.viewportOffset = m.cursor - maxVisibleItems + 1
			}
		}

	// T046: Spacebar to copy code to clipboard
	case " ", "enter":
		if len(m.filteredIndices) > 0 && m.cursor < len(m.filteredIndices) {
			// Get actual service index from filtered list
			serviceIdx := m.filteredIndices[m.cursor]
			service := m.services[serviceIdx]
			code := m.totpCodes[service.Name]
			if code != "" {
				// T047: Copy to clipboard with visual confirmation
				if err := clipboard.Copy(code); err != nil {
					// T048: Clipboard error handling with fallback
					m.copyStatus = "⚠ Clipboard unavailable. Code: " + code
				} else {
					m.copyStatus = "✓ Copied to clipboard"
				}
				m.copyStatusTime = time.Now()

				// Update LastUsed timestamp
				m.store.UpdateLastUsed(service.Name)
				_ = m.store.Save()
			}
		}

	// Home/End keys for quick navigation
	case "home", "g":
		m.cursor = 0
		m.viewportOffset = 0

	case "end", "G":
		if len(m.filteredIndices) > 0 {
			m.cursor = len(m.filteredIndices) - 1
			// Scroll to show last item
			maxVisibleItems := (m.height - 9) / 3
			if maxVisibleItems < 1 {
				maxVisibleItems = 1
			}
			if m.cursor >= maxVisibleItems {
				m.viewportOffset = m.cursor - maxVisibleItems + 1
			}
		}
	}

	return m, nil
}
