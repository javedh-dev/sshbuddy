package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"sshbuddy/pkg/models"
	"strings"
)

// ExecuteSSH executes SSH connection in the foreground
func ExecuteSSH(host models.Host) error {
	port := host.Port
	if port == "" {
		port = "22"
	}

	var args []string

	// Add port
	args = append(args, "-p", port)

	// Add identity file if specified
	if host.IdentityFile != "" {
		args = append(args, "-i", host.IdentityFile)
	}

	// Add proxy jump if specified
	if host.ProxyJump != "" {
		args = append(args, "-J", host.ProxyJump)
	}

	// If a default path is specified, use -t to allocate a pseudo-terminal
	// and execute a command to cd into the directory
	if host.DefaultPath != "" {
		args = append(args, "-t")
		args = append(args, fmt.Sprintf("%s@%s", host.User, host.Hostname))

		// Escape the path for use in double quotes to prevent command injection
		// but allow tilde and variable expansion
		escapedPath := escapeForDoubleQuotes(host.DefaultPath)

		// Use double quotes to allow tilde (~) expansion and variable substitution
		// while still protecting against spaces and most special characters
		args = append(args, fmt.Sprintf("cd \"%s\" && exec $SHELL -l", escapedPath))
	} else {
		// Standard connection without default path
		args = append(args, fmt.Sprintf("%s@%s", host.User, host.Hostname))
	}

	cmd := exec.Command("ssh", args...)

	// Connect to current terminal for interactive SSH session
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run SSH in foreground and wait for it to complete
	return cmd.Run()
}

// escapeForDoubleQuotes escapes characters that need escaping within double quotes
// and converts ~ to $HOME for proper expansion (since ~ doesn't expand in double quotes)
func escapeForDoubleQuotes(path string) string {
	// Convert tilde to $HOME for expansion (~ doesn't expand inside double quotes)
	if strings.HasPrefix(path, "~/") {
		path = "$HOME/" + path[2:]
	} else if path == "~" {
		path = "$HOME"
	}

	// Replace backslash with escaped backslash
	path = strings.ReplaceAll(path, "\\", "\\\\")
	// Replace double quote with escaped double quote
	path = strings.ReplaceAll(path, "\"", "\\\"")
	// Replace backtick with escaped backtick (prevents command substitution)
	path = strings.ReplaceAll(path, "`", "\\`")
	// Don't escape $ as we want variable expansion (including $HOME)
	return path
}
