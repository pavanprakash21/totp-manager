package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss styles for TUI
// (T042: Create Lipgloss styles for list, headers, borders)

var (
	// Colors
	colorPrimary   = lipgloss.Color("#00D9FF")
	colorSecondary = lipgloss.Color("#7D56F4")
	colorSuccess   = lipgloss.Color("#04B575")
	colorWarning   = lipgloss.Color("#FFB86C")
	colorMuted     = lipgloss.Color("#BBBBBB")
	colorBorder    = lipgloss.Color("#BBBBBB")

	// Header style
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			BorderBottom(true).
			PaddingBottom(1).
			PaddingLeft(2)

	// Service list item styles - boxed rows
	itemStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			PaddingLeft(2).
			PaddingRight(2).
			Width(80)

	selectedItemStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorSecondary).
				// Background(colorSecondary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				PaddingLeft(2).
				PaddingRight(2).
				Width(80)

	// Service name style
	serviceNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorPrimary).
				Width(50)

	selectedServiceNameStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#FFFFFF")).
					Width(50)

	// TOTP code style
	codeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSuccess).
			Align(lipgloss.Right).
			Width(10)

	selectedCodeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Align(lipgloss.Right).
				Width(10)

	// Global countdown timer style
	timerStyle = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true).
			PaddingLeft(2)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingTop(1).
			PaddingLeft(2)

	// Status message styles
	successStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true).
			PaddingLeft(2)

	warningStyle = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true).
			PaddingLeft(2)

	// Empty state style
	emptyStateStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true).
			PaddingLeft(2).
			PaddingTop(2)

	// Border style
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	// Search query style
	searchQueryStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true).
				PaddingLeft(2)
)
