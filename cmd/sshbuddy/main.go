package main

import (
	"fmt"
	"os"
	"sshbuddy/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("sshbuddy version %s\n", version)
		os.Exit(0)
	}
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// Check if we need to connect to a host
	if m, ok := finalModel.(tui.Model); ok {
		if m.GetSelectedHost() != nil {
			host := m.GetSelectedHost()
			fmt.Printf("Connecting to %s@%s...\n", host.User, host.Hostname)
			if err := tui.ExecuteSSH(*host); err != nil {
				fmt.Printf("Error connecting to host: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
