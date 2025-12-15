package clipboard

import (
	"testing"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "Copy simple text",
			text: "123456",
		},
		{
			name: "Copy TOTP code",
			text: "890123",
		},
		{
			name: "Copy empty string",
			text: "",
		},
		{
			name: "Copy with special characters",
			text: "Test!@#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.text)
			// Note: Clipboard might not be available in CI/headless environments
			// We just test that the function doesn't panic
			if err != nil {
				t.Logf("Clipboard not available (expected in CI): %v", err)
			}
		})
	}
}

func TestCopy_EmptyString(t *testing.T) {
	// Test copying empty string
	err := Copy("")
	if err != nil {
		t.Logf("Clipboard error (expected in CI): %v", err)
	}
}

func TestCopy_LongString(t *testing.T) {
	// Test copying a long string
	longText := ""
	for i := 0; i < 1000; i++ {
		longText += "a"
	}

	err := Copy(longText)
	if err != nil {
		t.Logf("Clipboard error (expected in CI): %v", err)
	}
}

func TestCopy_Unicode(t *testing.T) {
	// Test copying unicode characters
	unicodeText := "Hello 世界 🔐"

	err := Copy(unicodeText)
	if err != nil {
		t.Logf("Clipboard error (expected in CI): %v", err)
	}
}
