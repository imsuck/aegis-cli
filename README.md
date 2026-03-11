# aegis-cli

> [!Note]
> This was written entirely by an LLM.
>
> Just a heads-up if you are uncomfortable with using LLM-generated code.

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
./aegis-cli <path-to-vault.json> [-timeout duration]
```

You will be prompted to enter your vault password.

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-timeout` | `60s` | Auto-exit after duration of inactivity (e.g., `-timeout 30s`, `-timeout 2m`) |

### Examples

```bash
# Default 60 second timeout
./aegis-cli vault.json

# 30 second timeout
./aegis-cli vault.json -timeout 30s

# 2 minute timeout
./aegis-cli vault.json -timeout 2m

# Disable auto-exit
./aegis-cli vault.json -timeout 0
```

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

## Security

### Auto-Exit

The application automatically exits after 60 seconds of inactivity by default to prevent unauthorized access to your decrypted vault. This timeout can be configured using the `-timeout` flag:

- Set a custom duration (e.g., `-timeout 30s` for 30 seconds)
- Disable auto-exit entirely with `-timeout 0`

Any key press resets the inactivity timer, so simply pressing a key will keep the session active. A warning message appears in the last 10 seconds before auto-exit.

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

## TODO

- [Support more OTP types](https://github.com/beemdevelopment/Aegis/blob/59d5c640d6a2d0d16f243dbdd735758eed65bc63/app/src/main/java/com/beemdevelopment/aegis/crypto/otp/OTP.java): HTOP, Steam, MOTP, Yandex

## License

MIT

## Acknowledgments

- [Aegis Authenticator](https://getaegis.app) - The original Android authenticator
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [charm.sh](https://charm.sh) - Beautiful terminal UI components
