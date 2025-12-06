# CLI Usage

SSHBuddy can be used both as an interactive TUI and as a command-line tool for quick SSH connections.

## Quick Connect

Connect to a host directly without launching the TUI:

```bash
# Connect using full command
sshbuddy connect <alias>

# Connect using short command
sshbuddy c <alias>
```

**Examples:**
```bash
sshbuddy connect Atlas
sshbuddy c "Pi Zero 2W"
```

The alias matching is case-insensitive, so `sshbuddy c atlas` works the same as `sshbuddy c Atlas`.

## List Hosts

View all configured hosts:

```bash
sshbuddy list
# or
sshbuddy ls
```

This displays all hosts from your configuration, including those from:
- Manual entries (SSHBuddy config)
- SSH config file
- Termix API

## Import from Termix

Import hosts from your Termix server into your local manual configuration:

```bash
sshbuddy import termix
```

This command fetches all hosts from your Termix API and saves them as local manual hosts, allowing you to:
- Keep a local copy of Termix hosts
- Edit Termix hosts locally
- Use Termix hosts even when offline
- Migrate from Termix to SSHBuddy

**Options:**
```bash
# Skip hosts that already exist locally
sshbuddy import termix

# Overwrite existing hosts with the same alias
sshbuddy import termix --overwrite
```

**Example output:**
```
Connecting to Termix API at https://termix.example.com/api...
Found 15 host(s) in Termix

✓ Imported: Production Server (admin@prod.example.com)
✓ Imported: Dev Server (dev@dev.example.com)
✗ Skipped: Database Server (already exists, use --overwrite to replace)

--- Import Summary ---
Imported: 2
Skipped:  1
Total:    3
```

**Requirements:**
- Termix integration must be configured (see [Data Sources](data-sources.md))
- You must authenticate in the TUI at least once before using import

## Shell Completion

SSHBuddy supports autocomplete for bash, zsh, and fish shells. This enables tab completion for:
- Commands (connect, list, completion, etc.)
- Host aliases when using the connect command
- Shell names when generating completion scripts

### Installation

**Completions are automatically installed** when you install SSHBuddy via:
- `brew install sshbuddy` (Homebrew)
- `make install` (from source)

After installation, restart your terminal or source your shell config:

```bash
# Bash
source ~/.bashrc

# Zsh
exec zsh

# Fish (automatic in new shells)
```

### Manual Installation

If you installed via binary download or need to reinstall completions:

```bash
sshbuddy completion install
```

This will detect your shell and add the completion script to the appropriate config file.

### Manual Installation

If you prefer manual installation or need to customize the setup:

#### Bash

Add to your `~/.bashrc`:

```bash
source <(sshbuddy completion bash)
```

Or for one-time use:
```bash
source <(sshbuddy completion bash)
```

#### Zsh

Add to your `~/.zshrc`:

```bash
source <(sshbuddy completion zsh)
```

Or for one-time use:
```bash
source <(sshbuddy completion zsh)
```

#### Fish

Add to your fish config:

```bash
sshbuddy completion fish | source
```

Or save it permanently:
```bash
sshbuddy completion fish > ~/.config/fish/completions/sshbuddy.fish
```

## Interactive TUI

Launch the full interactive interface:

```bash
sshbuddy
```

This is the default behavior when no command is specified.

## Other Commands

```bash
# Show version
sshbuddy --version
sshbuddy -v

# Show help
sshbuddy --help
sshbuddy -h
sshbuddy help
```

## Examples

```bash
# Install completion (one-time setup)
sshbuddy completion install

# Restart your terminal or source your config
source ~/.bashrc  # or ~/.zshrc for zsh

# List all hosts and connect to one
sshbuddy list
sshbuddy c Atlas

# Quick connect with autocomplete
sshbuddy c <TAB>  # Shows all available aliases

# Autocomplete works for commands too
sshbuddy co<TAB>  # Completes to "connect"
sshbuddy connect <TAB>  # Shows all host aliases
```

## Tips

1. **Aliases with spaces**: Use quotes when connecting to hosts with spaces in their alias:
   ```bash
   sshbuddy c "Pi Zero 2W"
   ```

2. **Case insensitive**: Alias matching is case-insensitive for convenience:
   ```bash
   sshbuddy c atlas  # Works for "Atlas"
   ```

3. **Quick access**: Create shell aliases for frequently used hosts:
   ```bash
   alias ssh-atlas='sshbuddy c Atlas'
   ```

4. **Autocomplete setup**: Add the completion source command to your shell's RC file for persistent autocomplete across sessions.
