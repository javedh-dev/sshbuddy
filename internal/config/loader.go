package config

import (
	"fmt"

	"sshbuddy/internal/ssh"
	"sshbuddy/internal/termix"
	"sshbuddy/pkg/models"
)

// LoadConfig loads configuration and aggregates hosts from all enabled sources
func LoadConfig() (*models.Config, error) {
	// Load base config from file
	config, err := LoadConfigRaw()
	if err != nil {
		logError("LoadConfigRaw failed", err)
		return nil, err
	}

	// Initialize favorites map if nil
	if config.Favorites == nil {
		config.Favorites = make(map[string]bool)
	}

	// Mark manual hosts and apply favorites
	for i := range config.Hosts {
		if config.Hosts[i].Source == "" {
			config.Hosts[i].Source = "manual"
		}
		// Apply favorite status from saved config
		if config.Favorites[config.Hosts[i].Alias] {
			config.Hosts[i].Favorite = true
		}
	}

	// Track all aliases to avoid duplicates
	existingAliases := make(map[string]bool)

	// Only add manual hosts if SSHBuddy source is enabled
	if config.Sources.SSHBuddyEnabled {
		for _, host := range config.Hosts {
			existingAliases[host.Alias] = true
		}
	} else {
		// Clear manual hosts if disabled
		config.Hosts = []models.Host{}
	}

	// Load hosts from SSH config if enabled
	if config.Sources.SSHConfigEnabled && config.SSH.Enabled {
		sshHosts, err := ssh.LoadHostsFromSSHConfig()
		if err == nil {
			// Mark SSH config hosts and apply favorites
			for i := range sshHosts {
				sshHosts[i].Source = "ssh-config"
				if config.Favorites[sshHosts[i].Alias] {
					sshHosts[i].Favorite = true
				}
			}

			// Add SSH config hosts that don't conflict
			for _, sshHost := range sshHosts {
				if !existingAliases[sshHost.Alias] {
					config.Hosts = append(config.Hosts, sshHost)
					existingAliases[sshHost.Alias] = true
				}
			}
		}
	}

	// Load hosts from Termix API if enabled
	if config.Sources.TermixEnabled && config.Termix.Enabled && config.Termix.BaseURL != "" {
		logError("Termix config loaded", fmt.Errorf("baseUrl=%s", config.Termix.BaseURL))

		client := termix.NewClient(config.Termix.BaseURL, config.Termix.JWT, config.Termix.JWTExpiry)

		// Try to fetch hosts without credentials first (using cached token)
		termixHosts, termixFetchErr := client.FetchHosts("", "")

		// If auth is required, return a special error that the TUI can handle
		if termixFetchErr != nil {
			if _, isAuthError := termixFetchErr.(*termix.AuthError); isAuthError {
				// Return auth error to trigger credential prompt in TUI
				return nil, termixFetchErr
			}

			// Log other errors
			logError("Termix FetchHosts failed", termixFetchErr)

			// Return error to show in UI with config file hint
			configPath, _ := GetDataPath()
			fullError := fmt.Errorf("%w\n\nCheck your Termix configuration at: %s", termixFetchErr, configPath)
			logError("Returning error to UI", fullError)
			return nil, fullError
		}

		logError("Termix hosts fetched successfully", fmt.Errorf("count=%d", len(termixHosts)))

		// Add Termix hosts that don't conflict and apply favorites
		for _, termixHost := range termixHosts {
			if !existingAliases[termixHost.Alias] {
				if config.Favorites[termixHost.Alias] {
					termixHost.Favorite = true
				}
				config.Hosts = append(config.Hosts, termixHost)
				existingAliases[termixHost.Alias] = true
			}
		}

		// Save the JWT token and expiry if they were updated
		if client.GetJWT() != config.Termix.JWT || client.GetJWTExpiry() != config.Termix.JWTExpiry {
			config.Termix.JWT = client.GetJWT()
			config.Termix.JWTExpiry = client.GetJWTExpiry()
			SaveConfig(config)
		}
	}

	// Sort hosts: favorites first, then by alias
	sortHostsByFavorite(config.Hosts)

	return config, nil
}

// sortHostsByFavorite sorts hosts with favorites at the top
func sortHostsByFavorite(hosts []models.Host) {
	// Simple bubble sort to move favorites to the top while maintaining relative order
	n := len(hosts)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// If current is not favorite but next is, swap them
			if !hosts[j].Favorite && hosts[j+1].Favorite {
				hosts[j], hosts[j+1] = hosts[j+1], hosts[j]
			}
		}
	}
}
