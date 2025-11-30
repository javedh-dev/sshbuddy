package cli

import (
	"fmt"
	"os"
	"sshbuddy/internal/config"
	"sshbuddy/internal/tui"
	"sshbuddy/pkg/models"
	"strings"
)

// HandleCLI processes command-line arguments and returns true if handled
func HandleCLI(args []string, version string) bool {
	if len(args) < 2 {
		return false
	}

	command := args[1]

	switch command {
	case "--version", "-v":
		fmt.Printf("sshbuddy version %s\n", version)
		return true

	case "connect", "c":
		if len(args) < 3 {
			fmt.Println("Usage: sshbuddy connect <alias>")
			os.Exit(1)
		}
		ConnectByAlias(args[2])
		return true

	case "list", "ls":
		ListHosts()
		return true

	case "completion":
		if len(args) < 3 {
			fmt.Println("Usage: sshbuddy completion <shell>")
			fmt.Println("       sshbuddy completion install")
			fmt.Println("Supported shells: bash, zsh, fish")
			os.Exit(1)
		}
		if args[2] == "install" {
			InstallCompletion()
		} else {
			GenerateCompletion(args[2])
		}
		return true

	case "help", "-h", "--help":
		PrintHelp(version)
		return true
	}

	return false
}

// ConnectByAlias connects to a host by its alias
func ConnectByAlias(alias string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find host by alias (case-insensitive)
	var targetHost *models.Host
	for _, host := range cfg.Hosts {
		if strings.EqualFold(host.Alias, alias) {
			targetHost = &host
			break
		}
	}

	if targetHost == nil {
		fmt.Printf("Host with alias '%s' not found\n", alias)
		fmt.Println("\nAvailable hosts:")
		for _, host := range cfg.Hosts {
			fmt.Printf("  - %s (%s@%s)\n", host.Alias, host.User, host.Hostname)
		}
		os.Exit(1)
	}

	fmt.Printf("Connecting to %s (%s@%s)...\n", targetHost.Alias, targetHost.User, targetHost.Hostname)
	if err := tui.ExecuteSSH(*targetHost); err != nil {
		fmt.Printf("Error connecting to host: %v\n", err)
		os.Exit(1)
	}
}

// ListHosts lists all configured hosts
func ListHosts() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Hosts) == 0 {
		fmt.Println("No hosts configured")
		return
	}

	fmt.Println("Available hosts:")
	for _, host := range cfg.Hosts {
		source := ""
		if host.Source != "" && host.Source != "manual" {
			source = fmt.Sprintf(" [%s]", host.Source)
		}
		fmt.Printf("  %-20s %s@%s:%s%s\n", host.Alias, host.User, host.Hostname, host.Port, source)
	}
}

// PrintHelp prints usage information
func PrintHelp(version string) {
	fmt.Printf("sshbuddy version %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  sshbuddy                    Launch interactive TUI")
	fmt.Println("  sshbuddy connect <alias>    Connect to host by alias")
	fmt.Println("  sshbuddy c <alias>          Connect to host by alias (short)")
	fmt.Println("  sshbuddy list               List all configured hosts")
	fmt.Println("  sshbuddy ls                 List all configured hosts (short)")
	fmt.Println("  sshbuddy completion install Auto-install completion for your shell")
	fmt.Println("  sshbuddy completion <shell> Generate shell completion script")
	fmt.Println("  sshbuddy --version          Show version")
	fmt.Println("  sshbuddy --help             Show this help")
	fmt.Println("\nShell Completion:")
	fmt.Println("  Auto:   sshbuddy completion install")
	fmt.Println("  Bash:   source <(sshbuddy completion bash)")
	fmt.Println("  Zsh:    source <(sshbuddy completion zsh)")
	fmt.Println("  Fish:   sshbuddy completion fish | source")
}
