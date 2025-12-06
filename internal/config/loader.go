package config

import (
	"fmt"

	"sshbuddy/internal/ssh"
	"sshbuddy/internal/termix"
	"sshbuddy/pkg/models"
)

// LoadConfig loads configuration and aggregates hosts from all enabled sources
// Priority order: Termix (highest) → SSH Config → Manual (lowest)
// For duplicate aliases, the highest priority source wins, but all sources are tracked in AvailableIn
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

	// Map to track hosts by alias: alias -> host
	hostMap := make(map[string]*models.Host)

	// Track all hosts we'll process, in priority order
	var allHosts []models.Host

	// PRIORITY 1: Load hosts from Termix API if enabled (HIGHEST PRIORITY)
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
		allHosts = append(allHosts, termixHosts...)

		// Save the JWT token and expiry if they were updated
		if client.GetJWT() != config.Termix.JWT || client.GetJWTExpiry() != config.Termix.JWTExpiry {
			config.Termix.JWT = client.GetJWT()
			config.Termix.JWTExpiry = client.GetJWTExpiry()
			SaveConfig(config)
		}
	}

	// PRIORITY 2: Load hosts from SSH config if enabled
	if config.Sources.SSHConfigEnabled && config.SSH.Enabled {
		sshHosts, err := ssh.LoadHostsFromSSHConfig()
		if err == nil {
			// Mark SSH config hosts
			for i := range sshHosts {
				sshHosts[i].Source = "ssh-config"
			}
			allHosts = append(allHosts, sshHosts...)
		}
	}

	// PRIORITY 3: Load manual hosts if enabled (LOWEST PRIORITY)
	if config.Sources.SSHBuddyEnabled {
		for i := range config.Hosts {
			if config.Hosts[i].Source == "" {
				config.Hosts[i].Source = "manual"
			}
		}
		allHosts = append(allHosts, config.Hosts...)
	}

	// Process all hosts: merge duplicates and track sources
	for _, host := range allHosts {
		if existing, found := hostMap[host.Alias]; found {
			// Host already exists - add this source to AvailableIn
			existing.AvailableIn = append(existing.AvailableIn, host.Source)
		} else {
			// New host - initialize AvailableIn with its source
			host.AvailableIn = []string{host.Source}
			hostMap[host.Alias] = &host
		}
	}

	// Convert map back to slice
	config.Hosts = make([]models.Host, 0, len(hostMap))
	for _, host := range hostMap {
		// Apply favorite status from saved config
		if config.Favorites[host.Alias] {
			host.Favorite = true
		}
		config.Hosts = append(config.Hosts, *host)
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
