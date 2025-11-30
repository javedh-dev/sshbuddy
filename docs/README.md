# SSHBuddy Documentation

Welcome to the SSHBuddy documentation! This guide will help you get the most out of SSHBuddy's features.

## Quick Links

- **[Getting Started](getting-started.md)** - Installation and first steps
- **[CLI Usage](cli-usage.md)** - Command-line interface and autocomplete
- **[Configuration](configuration.md)** - Detailed configuration options
- **[Data Sources](data-sources.md)** - Working with multiple host sources
- **[Keyboard Shortcuts](keyboard-shortcuts.md)** - Complete shortcut reference
- **[Themes](themes.md)** - Theme customization guide
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions

## Quick Reference

### Interactive Mode

```bash
# Launch the TUI
sshbuddy
```

**Essential Shortcuts:**
- `Enter` - Connect to host
- `n` - Add new host
- `e` - Edit host
- `/` - Search hosts
- `p` - Ping all hosts
- `s` - Settings
- `q` - Quit

### Command-Line Mode

```bash
# Quick connect
sshbuddy connect <alias>
sshbuddy c <alias>

# List hosts
sshbuddy list
sshbuddy ls

# Generate completion
sshbuddy completion bash|zsh|fish
```

### Shell Completion Setup

**Automatic (Recommended)**:
```bash
sshbuddy completion install
```

**Manual Setup:**

**Bash** - Add to `~/.bashrc`:
```bash
source <(sshbuddy completion bash)
```

**Zsh** - Add to `~/.zshrc`:
```bash
source <(sshbuddy completion zsh)
```

**Fish** - Save to completions:
```bash
sshbuddy completion fish > ~/.config/fish/completions/sshbuddy.fish
```

## Features Overview

### Connection Management
- Multiple data sources (Manual, SSH Config, Termix)
- Live ping status indicators
- Tag-based organization
- Quick host duplication

### User Interface
- Two-column grid layout
- Instant search and filtering
- Six color themes
- Keyboard-first navigation

### SSH Features
- Full SSH config support
- SSH key authentication
- ProxyJump support
- Custom ports and options

### Integration
- Termix API support
- SSH config auto-import
- Unified JSON configuration
- Cross-platform compatibility

## Getting Help

- Check the [Troubleshooting Guide](troubleshooting.md) for common issues
- Run `sshbuddy --help` for command-line usage
- Press `?` in the TUI for keyboard shortcuts (if implemented)
- Visit the [GitHub repository](https://github.com/javedh-dev/sshbuddy) for issues and discussions

## Configuration File

SSHBuddy stores its configuration at:
```
~/.config/sshbuddy/config.json
```

You can edit this file directly or use the settings menu (press `s` in the TUI).

## Data Sources

SSHBuddy can load hosts from three sources:

1. **Manual Hosts** (◆) - Added through SSHBuddy
2. **SSH Config** (■) - From `~/.ssh/config`
3. **Termix API** (▲) - From your Termix server

Each source can be enabled or disabled in settings.
