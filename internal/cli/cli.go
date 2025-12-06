package cli

import (
	"fmt"
	"os"
	"sshbuddy/internal/config"
	"sshbuddy/internal/ssh"
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

	case "import":
		if len(args) < 3 {
			fmt.Println("Usage: sshbuddy import termix [--overwrite]")
			fmt.Println("       Import hosts from Termix API to local configuration")
			fmt.Println("\nOptions:")
			fmt.Println("  --overwrite    Overwrite existing hosts with the same alias")
			os.Exit(1)
		}
		if args[2] == "termix" {
			overwrite := false
			if len(args) > 3 && args[3] == "--overwrite" {
				overwrite = true
			}
			ImportFromTermix(overwrite)
		} else {
			fmt.Printf("Unknown import source: %s\n", args[2])
			fmt.Println("Supported sources: termix")
			os.Exit(1)
		}
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
	if err := ssh.ExecuteSSH(*targetHost); err != nil {
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

// ImportFromTermix imports hosts from Termix API to local configuration
func ImportFromTermix(overwrite bool) {
	// Load current config
	cfg, err := config.LoadConfigRaw()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Check if Termix is configured
	if !cfg.Termix.Enabled || cfg.Termix.BaseURL == "" {
		fmt.Println("Error: Termix is not configured")
		fmt.Println("Please configure Termix in the TUI (Settings > Termix)")
		os.Exit(1)
	}

	fmt.Printf("Connecting to Termix API at %s...\n", cfg.Termix.BaseURL)

	// Try to fetch hosts (might need authentication)
	fullCfg, err := config.LoadConfig()
	if err != nil {
		// Check if it's an auth error
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nNote: You may need to authenticate in the TUI first")
		os.Exit(1)
	}

	// Filter Termix hosts
	var termixHosts []models.Host
	for _, host := range fullCfg.Hosts {
		if host.Source == "termix" {
			termixHosts = append(termixHosts, host)
		}
	}

	if len(termixHosts) == 0 {
		fmt.Println("No hosts found in Termix")
		return
	}

	fmt.Printf("Found %d host(s) in Termix\n\n", len(termixHosts))

	// Build map of existing aliases
	existingAliases := make(map[string]bool)
	for _, host := range cfg.Hosts {
		existingAliases[host.Alias] = true
	}

	// Import hosts
	imported := 0
	skipped := 0
	updated := 0

	for _, host := range termixHosts {
		// Change source from termix to manual
		host.Source = "manual"

		if existingAliases[host.Alias] {
			if overwrite {
				// Find and replace the existing host
				for i, existing := range cfg.Hosts {
					if existing.Alias == host.Alias {
						cfg.Hosts[i] = host
						fmt.Printf("✓ Updated: %s (%s@%s)\n", host.Alias, host.User, host.Hostname)
						updated++
						break
					}
				}
			} else {
				fmt.Printf("✗ Skipped: %s (already exists, use --overwrite to replace)\n", host.Alias)
				skipped++
			}
		} else {
			cfg.Hosts = append(cfg.Hosts, host)
			fmt.Printf("✓ Imported: %s (%s@%s)\n", host.Alias, host.User, host.Hostname)
			imported++
		}
	}

	// Save the updated config
	if imported > 0 || updated > 0 {
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("\nError saving config: %v\n", err)
			os.Exit(1)
		}
	}

	// Print summary
	fmt.Printf("\n--- Import Summary ---\n")
	fmt.Printf("Imported: %d\n", imported)
	if updated > 0 {
		fmt.Printf("Updated:  %d\n", updated)
	}
	if skipped > 0 {
		fmt.Printf("Skipped:  %d\n", skipped)
	}
	fmt.Printf("Total:    %d\n", len(termixHosts))
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
	fmt.Println("  sshbuddy import termix      Import hosts from Termix to local config")
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
