package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model interface
// (T041: View method for rendering service list)
func (m Model) View() string {
	var b strings.Builder

	// Header
	header := headerStyle.Render("🔐 TOTP Manager")
	b.WriteString(header)
	b.WriteString("\n\n")

	// T052: Empty state view with instructions
	if len(m.services) == 0 {
		emptyMsg := emptyStateStyle.Render(
			"No TOTP services configured yet.\n\n" +
				"To add a service:\n" +
				"  • Use CLI: totp add --name GitHub --secret YOUR_SECRET\n" +
				"  • Optional: totp add --name GitHub --identifier user@example.com --secret YOUR_SECRET\n" +
				"  • Or press 'a' to add via TUI (coming soon)\n",
		)
		b.WriteString(emptyMsg)
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press 'q' to quit"))
		return b.String()
	}

	// Global countdown timer at top
	timerText := timerStyle.Render(fmt.Sprintf("⏱  Refreshing in %ds", m.remainingTime))
	b.WriteString(timerText)
	b.WriteString("\n")

	// Search mode indicator or filter status
	if m.searchMode {
		searchText := searchQueryStyle.Render(fmt.Sprintf("Search: %s_", m.searchQuery))
		b.WriteString(searchText)
		b.WriteString(fmt.Sprintf("  (%d results)", len(m.filteredIndices)))
	} else if m.searchQuery != "" {
		// Show active filter when not in search mode
		filterText := searchQueryStyle.Render(fmt.Sprintf("Filter: %s", m.searchQuery))
		b.WriteString(filterText)
		b.WriteString(fmt.Sprintf("  (%d/%d services)", len(m.filteredIndices), len(m.services)))
	}
	b.WriteString("\n")

	// Service list with boxed rows (filtered)
	if len(m.filteredIndices) == 0 {
		noResultsMsg := emptyStateStyle.Render("No matching services found")
		b.WriteString(noResultsMsg)
		b.WriteString("\n")
	} else {
		for i, serviceIdx := range m.filteredIndices {
			service := m.services[serviceIdx]
			isSelected := i == m.cursor
			code := m.totpCodes[service.Name]
			if code == "" {
				code = "------"
			}

			line := m.renderServiceLine(service.Name, service.Identifier, code, isSelected)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	// Copy status message
	if m.copyStatus != "" {
		b.WriteString("\n")
		if strings.HasPrefix(m.copyStatus, "✓") {
			b.WriteString(successStyle.Render(m.copyStatus))
		} else {
			b.WriteString(warningStyle.Render(m.copyStatus))
		}
		b.WriteString("\n")
	}

	// Help text (context-aware)
	b.WriteString("\n")
	var helpText string
	if m.searchMode {
		helpText = helpStyle.Render("space/enter: copy • ↑/↓: navigate • backspace: delete • ctrl+u: clear • esc: done")
	} else if m.searchQuery != "" {
		// Filtered view (search done but not in search mode)
		helpText = helpStyle.Render("/: search • ctrl+u: clear filter • space/enter: copy • q: quit")
	} else {
		helpText = helpStyle.Render("/: search • ↑/k: up • ↓/j: down • space/enter: copy • q: quit")
	}
	b.WriteString(helpText)

	return b.String()
}

// renderServiceLine renders a single service line with proper alignment
func (m Model) renderServiceLine(name, identifier, code string, selected bool) string {
	// Build full service name with identifier
	fullName := name
	if identifier != "" {
		fullName = fmt.Sprintf("%s (%s)", name, identifier)
	}

	// Truncate name if too long (leave room for code)
	maxNameLen := 50
	if len(fullName) > maxNameLen {
		fullName = fullName[:maxNameLen-3] + "..."
	}

	if selected {
		// Selected row: full-width highlight
		nameStr := selectedServiceNameStyle.Render(fullName)
		codeStr := selectedCodeStyle.Render(code)
		line := lipgloss.JoinHorizontal(lipgloss.Top, nameStr, "  ", codeStr)
		return selectedItemStyle.Render(line)
	}

	// Normal row: colored text in box
	nameStr := serviceNameStyle.Render(fullName)
	codeStr := codeStyle.Render(code)
	line := lipgloss.JoinHorizontal(lipgloss.Top, nameStr, "  ", codeStr)
	return itemStyle.Render(line)
}
