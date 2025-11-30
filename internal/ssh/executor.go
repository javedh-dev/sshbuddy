package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"sshbuddy/pkg/models"
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

	// Add host
	args = append(args, fmt.Sprintf("%s@%s", host.User, host.Hostname))

	cmd := exec.Command("ssh", args...)

	// Connect to current terminal for interactive SSH session
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run SSH in foreground and wait for it to complete
	return cmd.Run()
}
