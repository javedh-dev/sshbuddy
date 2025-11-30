package cli

import (
	"fmt"
	"os"
)

// GenerateCompletion generates shell completion scripts
func GenerateCompletion(shell string) {
	switch shell {
	case "bash":
		_, _ = os.Stdout.WriteString(bashCompletion)
	case "zsh":
		_, _ = os.Stdout.WriteString(zshCompletion)
	case "fish":
		_, _ = os.Stdout.WriteString(fishCompletion)
	default:
		fmt.Printf("Unsupported shell: %s\n", shell)
		fmt.Println("Supported shells: bash, zsh, fish")
		os.Exit(1)
	}
}

const bashCompletion = `# sshbuddy bash completion script
_sshbuddy_completion() {
    local cur prev commands cmd
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    cmd="${COMP_WORDS[0]}"
    commands="connect|c list|ls completion help"

    # Complete subcommands and flags
    if [ $COMP_CWORD -eq 1 ]; then
        # If starts with -, show flags
        if [[ ${cur} == -* ]]; then
            COMPREPLY=( $(compgen -W "--version --help -v -h" -- ${cur}) )
        else
            local expanded_commands="connect c list ls completion help"
            COMPREPLY=( $(compgen -W "${expanded_commands}" -- ${cur}) )
        fi
        return 0
    fi

    # Complete aliases for connect/c command
    if [ "${prev}" == "connect" ] || [ "${prev}" == "c" ]; then
        local IFS=$'\n'
        local aliases=()
        while IFS= read -r line; do
            # Extract alias: trim leading/trailing spaces, then get text before double-space
            local alias=$(echo "$line" | sed 's/^  *//' | sed 's/  .*//')
            [ -n "$alias" ] && aliases+=("$alias")
        done < <($cmd list 2>/dev/null | tail -n +2)
        COMPREPLY=( $(compgen -W "$(printf '%%q\n' "${aliases[@]}")" -- ${cur}) )
        return 0
    fi

    # Complete shell names for completion command
    if [ "${prev}" == "completion" ]; then
        COMPREPLY=( $(compgen -W "install bash zsh fish" -- ${cur}) )
        return 0
    fi
}

complete -F _sshbuddy_completion sshbuddy
`

const zshCompletion = `#compdef sshbuddy

_sshbuddy() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    _arguments -C \
        '(- *)'{-h,--help}'[Show help]' \
        '(- *)'{-v,--version}'[Show version]' \
        '1: :->command' \
        '*::arg:->args'

    case $state in
        command)
            local -a commands
            commands=(
                {'connect','c'}':Connect to host by alias'
                {'list','ls'}':List all configured hosts'
                'completion:Generate shell completion script'
                'help:Show help'
            )
            _describe 'command' commands
            ;;
        args)
            case $line[1] in
                connect|c)
                    local -a hosts
                    while IFS= read -r line; do
                        # Extract alias: trim leading/trailing spaces, then get text before double-space
                        local alias=$(echo "$line" | sed 's/^  *//' | sed 's/  .*//')
                        [ -n "$alias" ] && hosts+=("${alias}")
                    done < <(sshbuddy list 2>/dev/null | tail -n +2)
                    _describe 'host aliases' hosts
                    ;;
                completion)
                    local -a shells
                    shells=(
                        'install:Auto-install for current shell'
                        'bash:Bash completion script'
                        'zsh:Zsh completion script'
                        'fish:Fish completion script'
                    )
                    _describe 'shells' shells
                    ;;
            esac
            ;;
    esac
}

compdef _sshbuddy sshbuddy
`

const fishCompletion = `# sshbuddy fish completion script

# Disable file completion
complete -c sshbuddy -f

# Flags
complete -c sshbuddy -s h -l help -d "Show help"
complete -c sshbuddy -s v -l version -d "Show version"

# Main commands
complete -c sshbuddy -n "__fish_use_subcommand" -a "connect" -d "Connect to host (or: c)"
complete -c sshbuddy -n "__fish_use_subcommand" -a "c" -d "Connect to host (or: connect)"
complete -c sshbuddy -n "__fish_use_subcommand" -a "list" -d "List all hosts (or: ls)"
complete -c sshbuddy -n "__fish_use_subcommand" -a "ls" -d "List all hosts (or: list)"
complete -c sshbuddy -n "__fish_use_subcommand" -a "completion" -d "Generate shell completion script"
complete -c sshbuddy -n "__fish_use_subcommand" -a "help" -d "Show help"

# Complete aliases for connect/c command
complete -c sshbuddy -n "__fish_seen_subcommand_from connect c" -a "(sshbuddy list 2>/dev/null | tail -n +2 | sed 's/^  *//' | sed 's/  .*//' | string escape)"

# Complete shell names for completion command
complete -c sshbuddy -n "__fish_seen_subcommand_from completion" -a "install" -d "Auto-install for current shell"
complete -c sshbuddy -n "__fish_seen_subcommand_from completion" -a "bash" -d "Bash completion script"
complete -c sshbuddy -n "__fish_seen_subcommand_from completion" -a "zsh" -d "Zsh completion script"
complete -c sshbuddy -n "__fish_seen_subcommand_from completion" -a "fish" -d "Fish completion script"
`

// InstallCompletion automatically installs completion for the detected shell
func InstallCompletion() {
	// Check if sshbuddy is in PATH
	binaryPath := getBinaryPath()
	if binaryPath == "" {
		fmt.Println("⚠ Warning: 'sshbuddy' is not in your PATH")
		fmt.Println("\nTo install sshbuddy to your PATH, run:")
		fmt.Println("  sudo cp sshbuddy /usr/local/bin/")
		fmt.Println("\nOr add the current directory to your PATH.")
		fmt.Println("\nCompletion will be installed but may not work until sshbuddy is in PATH.")
		fmt.Println()
	}

	shell := detectShell()
	if shell == "" {
		fmt.Println("Could not detect your shell.")
		fmt.Println("Please manually install completion using:")
		fmt.Println("  sshbuddy completion bash|zsh|fish")
		os.Exit(1)
	}

	fmt.Printf("Detected shell: %s\n", shell)

	switch shell {
	case "bash":
		installBashCompletion()
	case "zsh":
		installZshCompletion()
	case "fish":
		installFishCompletion()
	default:
		fmt.Printf("Unsupported shell: %s\n", shell)
		os.Exit(1)
	}
}

// detectShell detects the current shell
func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return ""
	}

	// Extract shell name from path
	if len(shell) > 0 {
		parts := []rune(shell)
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] == '/' {
				shell = string(parts[i+1:])
				break
			}
		}
	}

	return shell
}

// installBashCompletion installs bash completion
func installBashCompletion() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	bashrcPath := homeDir + "/.bashrc"
	completionLine := "\n# SSHBuddy completion\nsource <(sshbuddy completion bash)\n"

	// Check if already installed
	content, err := os.ReadFile(bashrcPath)
	if err == nil {
		if containsString(string(content), "sshbuddy completion bash") {
			fmt.Println("✓ Bash completion is already installed in ~/.bashrc")
			fmt.Println("\nTo activate in current shell, run:")
			fmt.Println("  source ~/.bashrc")
			return
		}
	}

	// Append to .bashrc
	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening ~/.bashrc: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString(completionLine); err != nil {
		fmt.Printf("Error writing to ~/.bashrc: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Bash completion installed to ~/.bashrc")
	fmt.Println("\nTo activate in current shell, run:")
	fmt.Println("  source ~/.bashrc")
	fmt.Println("\nOr restart your terminal.")
}

// installZshCompletion installs zsh completion
func installZshCompletion() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	zshrcPath := homeDir + "/.zshrc"
	completionLine := "\n# SSHBuddy completion\nautoload -Uz compinit\ncompinit\nsource <(sshbuddy completion zsh)\n"

	// Check if already installed
	content, err := os.ReadFile(zshrcPath)
	if err == nil {
		if containsString(string(content), "sshbuddy completion zsh") {
			fmt.Println("✓ Zsh completion is already installed in ~/.zshrc")
			fmt.Println("\nTo activate in current shell, run:")
			fmt.Println("  exec zsh")
			return
		}
	}

	// Check if compinit is already in .zshrc
	hasCompinit := false
	if err == nil {
		hasCompinit = containsString(string(content), "compinit")
	}

	// If compinit is already there, don't add it again
	if hasCompinit {
		completionLine = "\n# SSHBuddy completion\nsource <(sshbuddy completion zsh)\n"
	}

	// Append to .zshrc
	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening ~/.zshrc: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString(completionLine); err != nil {
		fmt.Printf("Error writing to ~/.zshrc: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Zsh completion installed to ~/.zshrc")
	fmt.Println("\nTo activate in current shell, run:")
	fmt.Println("  exec zsh")
	fmt.Println("\nOr restart your terminal.")
}

// installFishCompletion installs fish completion
func installFishCompletion() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Create fish completions directory if it doesn't exist
	fishCompletionDir := homeDir + "/.config/fish/completions"
	if err := os.MkdirAll(fishCompletionDir, 0755); err != nil {
		fmt.Printf("Error creating fish completions directory: %v\n", err)
		os.Exit(1)
	}

	completionPath := fishCompletionDir + "/sshbuddy.fish"

	// Check if already installed
	if _, err := os.Stat(completionPath); err == nil {
		fmt.Printf("✓ Fish completion is already installed at %s\n", completionPath)
		fmt.Println("\nCompletion is active in all new fish shells.")
		return
	}

	// Write completion file
	if err := os.WriteFile(completionPath, []byte(fishCompletion), 0644); err != nil {
		fmt.Printf("Error writing fish completion file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Fish completion installed to %s\n", completionPath)
	fmt.Println("\nCompletion is now active in all new fish shells.")
	fmt.Println("To activate in current shell, run:")
	fmt.Println("  source " + completionPath)
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

// findSubstring finds the index of substr in s, returns -1 if not found
func findSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// getBinaryPath checks if sshbuddy is in PATH and returns its path
func getBinaryPath() string {
	// Try to find sshbuddy in PATH
	paths := os.Getenv("PATH")
	if paths == "" {
		return ""
	}

	pathList := []string{}
	currentPath := ""
	for i := 0; i < len(paths); i++ {
		if paths[i] == ':' {
			if currentPath != "" {
				pathList = append(pathList, currentPath)
				currentPath = ""
			}
		} else {
			currentPath += string(paths[i])
		}
	}
	if currentPath != "" {
		pathList = append(pathList, currentPath)
	}

	for _, dir := range pathList {
		fullPath := dir + "/sshbuddy"
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return ""
}
