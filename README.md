# Aegis CLI

A terminal-based TUI for viewing and managing Aegis 2FA vault entries.

## Installation

```bash
cargo install --path .
```

## Usage

```bash
aegis-cli <path-to-vault.json>
```

You will be prompted for your vault password.

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `/` | Enter search mode |
| `Esc` | Exit search mode |
| `Enter` | Confirm search |
| `c` | Show/hide OTP code |
| `y` | Yank (copy) OTP code |
| `q` | Quit |

## Search

- Default search filters by issuer using fuzzy matching (powered by nucleo)
- Use property prefix for specific fields:
  - `%issuer <query>` - Search by issuer
  - `%name <query>` - Search by name
  - `%note <query>` - Search by note
  - `%favorite <query>` - Search by favorite status
  - `%type <query>` - Search by OTP type

Prefix matching: `%is` matches `issuer`, `%nam` matches `name`, etc.

## Security

- Password is prompted interactively (not stored)
- OTP codes are masked with asterisks by default
- Press `c` to reveal the code for the selected entry only
- Codes are shown as `******` with time remaining (e.g., `[****** | 15s]`)

## Supported OTP Types

- TOTP (RFC 6236)
- HOTP (RFC 4226) - not fully implemented
- Steam - not fully implemented
- MOTP - not fully implemented
- Yandex - not fully implemented

## Example

```bash
# Start with your encrypted vault
aegis-cli /path/to/aegis_encrypted.json

# Enter password when prompted
Password: test

# Navigate with j/k, search with /, show code with c, copy with y
```
