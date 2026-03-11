# aegis-cli

> [!Note]
> This was written entirely by an LLM.
>
> Just a heads-up if you are uncomfortable using LLM-generated code.

Terminal-based TUI client for Aegis authenticator vaults with passwords.

## Features

- Decrypt your password encrypted vault
- Fuzzy search from issuer, account name, notes
- Color-coded timer (wow eye candy!)

## Building

```bash
go build -o aegis-cli ./cmd/aegis-cli
```

## Usage

```bash
./aegis-cli <path-to-vault.json>
```

You will be prompted to enter your vault password.

## Keybindings

### Password Screen

| Key      | Action          |
|----------|-----------------|
| `enter`  | Submit password |
| `ctrl-c` | Quit            |

### Entry Table

| Key            | Action                 |
|----------------|------------------------|
| `j` / `↓`      | Move down              |
| `k` / `↑`      | Move up                |
| `g`            | Go to top              |
| `g`            | Go to bottom           |
| `/`            | Open search            |
| `c`            | Toggle code display    |
| `y`            | Copy code to clipboard |
| `q` / `ctrl-c` | Quit                   |

### Search

| Key                       | Action                    |
|---------------------------|---------------------------|
| `esc`                     | Cancel search             |
| `enter`                   | Accept search results     |
| `ctrl-a`                  | Move cursor to beginning  |
| `ctrl-e`                  | Move cursor to end        |
| `ctrl-u`                  | Clear input               |
| `ctrl-w`/`ctrl-backspace` | Delete word before cursor |

### Code Display Mode

| Key            | Action                 |
|----------------|------------------------|
| `j` / `down`   | Next entry             |
| `k` / `up`     | Previous entry         |
| `y`            | Copy code to clipboard |
| `c` / `esc`    | Back to table          |
| `q` / `ctrl-c` | Quit                   |


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

## License

MIT

## Acknowledgments

- [Aegis Authenticator](https://getaegis.app) - The original Android authenticator
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [charm.sh](https://charm.sh) - Beautiful terminal UI components
