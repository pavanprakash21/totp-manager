package tui

import (
	"testing"
)

// TestStyles tests that all styles are initialized
func TestStyles(t *testing.T) {
	// Test that styles render without panicking
	testText := "test"

	// Test each style individually
	if headerStyle.Render(testText) == "" {
		t.Error("headerStyle.Render() returned empty string")
	}
	if itemStyle.Render(testText) == "" {
		t.Error("itemStyle.Render() returned empty string")
	}
	if selectedItemStyle.Render(testText) == "" {
		t.Error("selectedItemStyle.Render() returned empty string")
	}
	if serviceNameStyle.Render(testText) == "" {
		t.Error("serviceNameStyle.Render() returned empty string")
	}
	if codeStyle.Render(testText) == "" {
		t.Error("codeStyle.Render() returned empty string")
	}
	if timerStyle.Render(testText) == "" {
		t.Error("timerStyle.Render() returned empty string")
	}
	if helpStyle.Render(testText) == "" {
		t.Error("helpStyle.Render() returned empty string")
	}
	if successStyle.Render(testText) == "" {
		t.Error("successStyle.Render() returned empty string")
	}
	if warningStyle.Render(testText) == "" {
		t.Error("warningStyle.Render() returned empty string")
	}
	if emptyStateStyle.Render(testText) == "" {
		t.Error("emptyStateStyle.Render() returned empty string")
	}
	if searchQueryStyle.Render(testText) == "" {
		t.Error("searchQueryStyle.Render() returned empty string")
	}
	if borderStyle.Render(testText) == "" {
		t.Error("borderStyle.Render() returned empty string")
	}

	// Test that widths are set for specific styles
	if itemStyle.GetWidth() != 80 {
		t.Errorf("itemStyle width = %d, want 80", itemStyle.GetWidth())
	}

	if selectedItemStyle.GetWidth() != 80 {
		t.Errorf("selectedItemStyle width = %d, want 80", selectedItemStyle.GetWidth())
	}

	if serviceNameStyle.GetWidth() != 50 {
		t.Errorf("serviceNameStyle width = %d, want 50", serviceNameStyle.GetWidth())
	}

	if codeStyle.GetWidth() != 10 {
		t.Errorf("codeStyle width = %d, want 10", codeStyle.GetWidth())
	}
} // TestColorConstants tests that color constants are defined
func TestColorConstants(t *testing.T) {
	colors := []struct {
		name  string
		color interface{}
	}{
		{"colorPrimary", colorPrimary},
		{"colorSecondary", colorSecondary},
		{"colorSuccess", colorSuccess},
		{"colorWarning", colorWarning},
		{"colorMuted", colorMuted},
		{"colorBorder", colorBorder},
	}

	for _, tc := range colors {
		if tc.color == nil {
			t.Errorf("%s should be defined", tc.name)
		}
	}
}
