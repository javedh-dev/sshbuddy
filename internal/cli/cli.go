package cli

import (
	"fmt"
	"os"
	"path/filepath"
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
			fmt.Println("Usage: sshbuddy import <source> [options]")
			fmt.Println("       sshbuddy import termix [--overwrite]")
			fmt.Println("       sshbuddy import ssh-config [--overwrite]")
			fmt.Println("\nOptions:")
			fmt.Println("  --overwrite    Overwrite existing hosts with the same alias")
			os.Exit(1)
		}

		overwrite := false
		// Check for overwrite flag in any remaining position
		for i := 3; i < len(args); i++ {
			if args[i] == "--overwrite" {
				overwrite = true
			}
		}

		if args[2] == "termix" {
			ImportFromTermix(overwrite)
		} else if args[2] == "ssh-config" {
			ImportFromSSHConfig(overwrite)
		} else {
			fmt.Printf("Unknown import source: %s\n", args[2])
			fmt.Println("Supported sources: termix, ssh-config")
			os.Exit(1)
		}
		return true

	case "export":
		if len(args) < 3 {
			fmt.Println("Usage: sshbuddy export <format> [options]")
			fmt.Println("       sshbuddy export ssh-config [--stdout] [--file <path>]")
			fmt.Println("\nOptions:")
			fmt.Println("  --stdout       Print to stdout instead of writing to file")
			fmt.Println("  --file <path>  Write to specific file (defaults to ~/.ssh/config)")
			os.Exit(1)
		}

		outputFile := "~/.ssh/config" // Default path
		toStdout := false

		for i := 3; i < len(args); i++ {
			if args[i] == "--stdout" {
				toStdout = true
				outputFile = ""
			} else if args[i] == "--file" && i+1 < len(args) {
				outputFile = args[i+1]
				toStdout = false
			}
		}

		if toStdout {
			outputFile = ""
		}

		if args[2] == "ssh-config" {
			ExportToSSHConfig(outputFile)
		} else {
			fmt.Printf("Unknown export format: %s\n", args[2])
			fmt.Println("Supported formats: ssh-config")
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

	// Save configuration
	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("\nError saving configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nImport complete! Imported: %d, Updated: %d, Skipped: %d\n", imported, updated, skipped)
}

// ImportFromSSHConfig imports hosts from SSH config to local configuration
func ImportFromSSHConfig(overwrite bool) {
	// Load current config (manual hosts)
	cfg, err := config.LoadConfigRaw()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Reading SSH config...")
	sshHosts, err := ssh.LoadHostsFromSSHConfig()
	if err != nil {
		fmt.Printf("Error loading SSH config: %v\n", err)
		os.Exit(1)
	}

	if len(sshHosts) == 0 {
		fmt.Println("No hosts found in SSH config")
		return
	}

	fmt.Printf("Found %d host(s) in SSH config\n\n", len(sshHosts))

	// Build map of existing aliases
	existingAliases := make(map[string]bool)
	for _, host := range cfg.Hosts {
		existingAliases[host.Alias] = true
	}

	// Import hosts
	imported := 0
	skipped := 0
	updated := 0

	for _, host := range sshHosts {
		// Change source to manual
		host.Source = "manual"
		host.AvailableIn = []string{"manual"}
		// Initialize Variant
		host.Variants = make(map[string]*models.Host)
		selfCopy := host
		host.Variants["manual"] = &selfCopy

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
				fmt.Printf("- Skipped: %s (already exists)\n", host.Alias)
				skipped++
			}
		} else {
			cfg.Hosts = append(cfg.Hosts, host)
			fmt.Printf("+ Imported: %s (%s@%s)\n", host.Alias, host.User, host.Hostname)
			imported++
		}
	}

	// Save configuration
	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("\nError saving configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nImport complete! Imported: %d, Updated: %d, Skipped: %d\n", imported, updated, skipped)
}

// ExportToSSHConfig exports local manual hosts to SSH config format
func ExportToSSHConfig(outputFile string) {
	// Load current config (manual hosts only)
	cfg, err := config.LoadConfigRaw()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Hosts) == 0 {
		fmt.Println("No manual hosts to export")
		return
	}

	var sb strings.Builder
	sb.WriteString("# Generated by SSHBuddy\n\n")

	for _, host := range cfg.Hosts {
		sb.WriteString(fmt.Sprintf("Host %s\n", host.Alias))
		sb.WriteString(fmt.Sprintf("    HostName %s\n", host.Hostname))

		if host.User != "" {
			sb.WriteString(fmt.Sprintf("    User %s\n", host.User))
		}

		if host.Port != "" && host.Port != "22" {
			sb.WriteString(fmt.Sprintf("    Port %s\n", host.Port))
		}

		if host.IdentityFile != "" {
			sb.WriteString(fmt.Sprintf("    IdentityFile %s\n", host.IdentityFile))
		}

		if host.ProxyJump != "" {
			sb.WriteString(fmt.Sprintf("    ProxyJump %s\n", host.ProxyJump))
		}

		sb.WriteString("\n")
	}

	output := sb.String()

	if outputFile != "" {
		// Expand ~ to home directory
		if strings.HasPrefix(outputFile, "~/") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				outputFile = filepath.Join(homeDir, outputFile[2:])
			}
		}

		// Create backup if file exists
		if _, err := os.Stat(outputFile); err == nil {
			backupFile := outputFile + ".bak"
			if err := os.Rename(outputFile, backupFile); err != nil {
				fmt.Printf("Warning: Could not create backup: %v\n", err)
			} else {
				fmt.Printf("Created backup at %s\n", backupFile)
			}
		} else {
			// Ensure directory exists
			if err := os.MkdirAll(filepath.Dir(outputFile), 0700); err != nil {
				fmt.Printf("Error creating directory: %v\n", err)
				os.Exit(1)
			}
		}

		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully exported %d hosts to %s\n", len(cfg.Hosts), outputFile)
	} else {
		fmt.Println(output)
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
	fmt.Println("  sshbuddy import termix      Import hosts from Termix")
	fmt.Println("  sshbuddy import ssh-config  Import hosts from SSH config")
	fmt.Println("  sshbuddy export ssh-config  Export hosts to SSH config format")
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
