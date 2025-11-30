package ssh

import (
	"os/exec"
	"sshbuddy/pkg/models"
	"strings"
)

// PingResult contains the result of a ping operation
type PingResult struct {
	Host     models.Host
	Status   bool   // true if reachable
	PingTime string // ping time in ms
}

// Ping checks if a host is reachable using a simple ping
func Ping(host models.Host) PingResult {
	// Use ping with 1 count and 1 second timeout
	cmd := exec.Command("ping", "-c", "1", "-W", "1", host.Hostname)
	output, err := cmd.CombinedOutput()

	// Parse ping time from output
	pingTime := ""
	if err == nil {
		// Extract time from ping output (e.g., "time=12.3 ms")
		outputStr := string(output)

		// Try to find "time=" pattern
		if idx := strings.Index(outputStr, "time="); idx != -1 {
			timeStr := outputStr[idx+5:]
			// Find the end of the time value (space or newline)
			endIdx := strings.IndexAny(timeStr, " \n\r")
			if endIdx != -1 {
				timeValue := strings.TrimSpace(timeStr[:endIdx])
				pingTime = timeValue
				// Add "ms" if not already present
				if !strings.HasSuffix(pingTime, "ms") {
					pingTime = pingTime + "ms"
				}
			}
		}
	}

	return PingResult{
		Host:     host,
		Status:   err == nil,
		PingTime: pingTime,
	}
}

// GetHostStatus returns a visual indicator for host status
func GetHostStatus(status bool) string {
	if status {
		return "🟢" // Green dot - reachable
	}
	return "🔴" // Red dot - unreachable
}

// GetHostKey creates a unique key for a host (for tracking ping status)
func GetHostKey(host models.Host) string {
	return strings.ToLower(host.Hostname + ":" + host.User)
}
