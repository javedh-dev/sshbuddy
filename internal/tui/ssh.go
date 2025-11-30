package tui

import (
	"sshbuddy/internal/ssh"
	"sshbuddy/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
)

// PingHost wraps the SSH ping function as a Bubble Tea command
func PingHost(host models.Host) tea.Cmd {
	return func() tea.Msg {
		result := ssh.Ping(host)
		return PingResultMsg{
			Host:     result.Host,
			Status:   result.Status,
			PingTime: result.PingTime,
		}
	}
}

// StartPingAll starts background ping for all hosts
func StartPingAll(hosts []models.Host) tea.Cmd {
	var cmds []tea.Cmd
	for _, host := range hosts {
		cmds = append(cmds, PingHost(host))
	}
	return tea.Batch(cmds...)
}

// GetHostStatus returns a visual indicator for host status
func GetHostStatus(status bool) string {
	return ssh.GetHostStatus(status)
}

// GetHostKey creates a unique key for a host (for tracking ping status)
func GetHostKey(host models.Host) string {
	return ssh.GetHostKey(host)
}
