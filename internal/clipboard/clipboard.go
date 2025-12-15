package clipboard

import (
	"github.com/atotto/clipboard"
)

// Copy copies text to the system clipboard
// (T047: Clipboard copy with visual confirmation)
// (T048: Clipboard error handling)
func Copy(text string) error {
	// Use atotto/clipboard for cross-platform support
	return clipboard.WriteAll(text)
}
