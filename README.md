# aegis-cli

A terminal-based TUI client for Aegis authenticator vaults.

## Features

- 🔓 Decrypt Aegis vaults with password
- 📱 View all TOTP entries in a clean table
- ⏱️ Real-time code refresh with countdown timer
- 🔍 Fuzzy search across issuer, name, and notes
- ⌨️ Vim-style keybindings (j/k navigation, / search)
- 📋 Clipboard copy with single keypress
- 🎨 Clean, modern TUI with color-coded timers

## Installation

```bash
go build -o aegis-cli ./cmd/aegis-cli
```

## Usage

```bash
./aegis-cli <path-to-vault.json>
```

### Example

```bash
# Open your encrypted Aegis vault
./aegis-cli /path/to/aegis_export.json
```

You will be prompted to enter your vault password.

## Keybindings

### Password Screen

| Key | Action |
|-----|--------|
| `Enter` | Submit password |
| `Ctrl+C` | Quit |

### Entry Table

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` | Go to top |
| `G` | Go to bottom |
| `/` | Open search |
| `c` | Toggle code display mode |
| `y` | Copy code to clipboard |
| `q` / `Ctrl+C` | Quit |

### Search

| Key | Action |
|-----|--------|
| `Type` | Filter entries (matches issuer, name, and note) |
| `Esc` | Cancel search |
| `Enter` | Accept search results |
| `Ctrl+A` | Move cursor to beginning |
| `Ctrl+E` | Move cursor to end |
| `Ctrl+U` | Clear input |
| `Ctrl+W` | Delete word before cursor |

### Code Display Mode

| Key | Action |
|-----|--------|
| `j` / `↓` | Next entry |
| `k` / `↑` | Previous entry |
| `y` | Copy code to clipboard |
| `c` / `Esc` | Back to table |
| `q` / `Ctrl+C` | Quit |

## Timer Colors

The remaining time until code refresh is color-coded:

- 🟢 **Green** (>15 seconds) - Code is fresh
- 🟠 **Orange** (5-15 seconds) - Code will refresh soon
- 🔴 **Red** (<5 seconds) - Code about to refresh

## Requirements

- Go 1.21+
- Linux/macOS/Windows
- Clipboard support:
  - Linux: X11 or Wayland (uses `github.com/atotto/clipboard`)
  - macOS: Native clipboard
  - Windows: Native clipboard

## Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run integration tests with test vault (password: test)
go test ./tests/integration/... -v
```

## Project Structure

```
aegis-cli/
├── cmd/
│   └── aegis-cli/
│       └── main.go          # Application entry point
├── internal/
│   ├── vault/
│   │   ├── types.go         # Vault data structures
│   │   ├── decrypt.go       # Vault decryption (scrypt + AES-GCM)
│   │   └── types_test.go    # Unit tests
│   ├── totp/
│   │   ├── totp.go          # TOTP code generation
│   │   └── totp_test.go     # Unit tests
│   ├── search/
│   │   ├── search.go        # Fuzzy search implementation
│   │   └── search_test.go   # Unit tests
│   └── tui/
│       ├── model.go         # TUI state model
│       ├── update.go        # TUI message handling
│       └── view.go          # TUI rendering
├── test/
│   └── resources/
│       └── aegis_encrypted.json  # Test vault (password: test)
└── docs/
    ├── vault.md             # Aegis vault format documentation
    └── decrypt.py           # Python reference implementation
```

## Security

- Uses scrypt for key derivation (N=32768, r=8, p=1)
- Uses AES-256-GCM for vault decryption
- Password is masked with asterisks during input
- No secrets are logged or written to disk

## License

MIT

## Acknowledgments

- [Aegis Authenticator](https://getaegis.app) - The original Android authenticator
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [charm.sh](https://charm.sh) - Beautiful terminal UI components
