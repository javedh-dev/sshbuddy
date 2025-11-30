# Keyboard Shortcuts

SSHBuddy is designed for keyboard-driven workflows. Here's a complete reference of all available shortcuts.

## Main List View

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up (previous row) |
| `↓` / `j` | Move down (next row) |
| `←` / `h` | Move left (previous column) |
| `→` / `l` | Move right (next column) |

The host list uses a two-column layout. Up/down moves between rows, while left/right moves between columns.

### Host Actions

| Key | Action |
|-----|--------|
| `Enter` | Connect to selected host |
| `n` | Add new host |
| `e` | Edit selected host (manual hosts only) |
| `c` | Duplicate selected host |
| `d` | Delete selected host (manual hosts only) |
| `f` | Toggle favorite status (shows ❤ icon beside source) |

### Utility Functions

| Key | Action |
|-----|--------|
| `/` | Search/filter hosts |
| `p` | Ping all hosts to check status |
| `s` | Open settings |
| `q` | Quit application |
| `Ctrl+C` | Force quit |

## Search/Filter Mode

| Key | Action |
|-----|--------|
| Type | Filter hosts by alias or hostname |
| `Esc` | Clear filter and return to full list |

## Host Form (Add/Edit)

### Navigation

| Key | Action |
|-----|--------|
| `Tab` / `↓` | Next field |
| `Shift+Tab` / `↑` | Previous field |
| `←` | Move to left column |
| `→` | Move to right column |

### Actions

| Key | Action |
|-----|--------|
| `Enter` | Save host |
| `Esc` | Cancel and return to list |

## Settings Menu

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Previous option |
| `↓` / `j` | Next option |

### Actions

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Toggle source or cycle theme |
| `e` | Edit configuration (Termix/SSH Config) |
| `Esc` | Return to main list |

## Configuration Edit Screens

### Termix Configuration

| Key | Action |
|-----|--------|
| `Tab` / `↑` / `↓` | Navigate between fields |
| `Enter` | Save configuration |
| `Esc` | Cancel changes |

### SSH Config Path

| Key | Action |
|-----|--------|
| Type | Enter custom SSH config path |
| `Enter` | Save path |
| `Esc` | Cancel |

## Termix Authentication

| Key | Action |
|-----|--------|
| `Tab` / `↑` / `↓` | Switch between username and password |
| `Enter` | Submit credentials |
| `Esc` | Cancel authentication |

## Delete Confirmation

| Key | Action |
|-----|--------|
| `y` / `Y` | Confirm deletion |
| `n` / `N` / `Esc` | Cancel deletion |

## Tips for Efficient Navigation

**Two-Column Layout**: The host list displays in two columns. Use `↑`/`↓` to move between rows and `←`/`→` to switch columns. This layout lets you see more hosts at once.

**Vim-Style Navigation**: If you're comfortable with Vim, you can use `h`, `j`, `k`, `l` for navigation in the main list.

**Quick Search**: Press `/` and start typing to instantly filter your hosts. This is the fastest way to find a specific server when you have many hosts.

**Ping Status**: Press `p` to check which hosts are online. Green dots indicate reachable hosts, red dots show offline hosts, and gray circles mean the status is unknown.

**Duplicate for Similar Hosts**: When adding multiple hosts with similar configurations, use `c` to duplicate an existing host and modify the copy. This saves time compared to filling out the form from scratch.

**Read-Only Indicators**: Hosts from SSH config and Termix can't be edited or deleted through SSHBuddy. The `e` and `d` keys only work on manually added hosts.

**Favorites**: Press `f` to mark a host as favorite. Favorited hosts are shown with a filled heart (❤) icon beside the source indicator and automatically sorted to the top of the list. This works for hosts from all sources (manual, SSH config, and Termix).
