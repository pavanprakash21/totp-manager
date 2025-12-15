# TOTP Manager TUI

A secure, terminal-based Time-based One-Time Password (TOTP) manager built with Go and Bubbletea.

## Features

- 🔐 **Secure Encrypted Storage**: All secrets encrypted with AES-256-GCM
- ⌨️ **Keyboard-Driven**: Navigate with arrow keys or vim bindings (hjkl)
- 📋 **Clipboard Integration**: Copy codes with spacebar
- 🎨 **Modern TUI**: Built with Bubbletea and Lipgloss
- ⚡ **Fast**: Sub-second launch, instant code generation
- 🔄 **Auto-Refresh**: Codes update every 30 seconds with countdown timer

## Installation

```bash
go install github.com/pavanprakash21/totp-manager-go@latest
```

Or build from source:

```bash
git clone https://github.com/pavanprakash21/totp-manager-go.git
cd totp-manager-go
go build -o totp main.go
```

## Usage

### Launch TUI

```bash
totp
```

On first launch, you'll be prompted to create a new passphrase. This passphrase encrypts all your TOTP secrets.

### Add Service via CLI

```bash
# Basic usage
totp add --name "GitHub" --secret "JBSWY3DPEHPK3PXP"

# With optional identifier (e.g., email or username)
totp add --name "GitHub" --identifier "user@example.com" --secret "JBSWY3DPEHPK3PXP"
```

### Change Passphrase

```bash
totp change-passphrase
```

## Keyboard Controls

- **↑/↓ or j/k**: Navigate through services
- **Space**: Copy selected TOTP code to clipboard
- **a**: Add new service (in TUI)
- **q or ESC**: Quit
- **?**: Show help

## Security

- All secrets are encrypted using AES-256-GCM
- Passphrase is never stored on disk
- Encryption keys derived using Argon2id (memory-hard KDF)
- Storage file has 0600 permissions (owner-only read/write)
- No secrets are logged or printed to terminal (except on explicit clipboard failure)

## Storage Location

Encrypted secrets are stored at:
- macOS/Linux: `~/.config/totp-manager/secrets.enc`

## Development

### Prerequisites

- Go 1.21 or later

### Build

```bash
go build -o totp main.go
```

### Test

```bash
go test ./...
```

### Test Coverage

```bash
./scripts/test-coverage.sh
```

## License

MIT

## Contributing

Contributions welcome! Please open an issue or pull request.
