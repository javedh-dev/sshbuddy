package config

import (
	"fmt"
	"sort"
	"strings"

	"sshbuddy/internal/ssh"
	"sshbuddy/internal/termix"
	"sshbuddy/pkg/models"
)

// LoadConfig loads configuration and aggregates hosts from all enabled sources
// Priority order: Manual (highest) → SSH Config → Termix (lowest)
// This allows local overrides (Manual) to take precedence over external sources while tracking availability
func LoadConfig() (*models.Config, error) {
	// Load base config from file (Manual hosts)
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

	// Track all hosts in order of processing to maintain sort stability where possible
	// We will process sources in priority order: Manual -> SSH Config -> Termix

	// PRIORITY 1: Manual hosts (HIGHEST PRIORITY - Overrides everything)
	if config.Sources.SSHBuddyEnabled {
		for i := range config.Hosts {
			// Ensure source is set
			if config.Hosts[i].Source == "" {
				config.Hosts[i].Source = "manual"
			}
			host := config.Hosts[i]
			host.AvailableIn = []string{"manual"}

			// Initialize variants
			host.Variants = make(map[string]*models.Host)
			selfCopy := host
			host.Variants["manual"] = &selfCopy

			hostMap[host.Alias] = &host
		}
	} else {
		// Clear hosts if manual source disabled
		config.Hosts = []models.Host{}
	}

	// PRIORITY 2: SSH Config
	if config.Sources.SSHConfigEnabled && config.SSH.Enabled {
		sshHosts, err := ssh.LoadHostsFromSSHConfig()
		if err == nil {
			// Process all hosts: merge duplicates and track sources
			for _, host := range sshHosts { // Changed from allHosts to sshHosts
				// Ensure source is set for SSH hosts
				host.Source = "ssh-config"

				if existing, found := hostMap[host.Alias]; found {
					// Host already exists
					existing.AvailableIn = append(existing.AvailableIn, host.Source)

					// Add this variant
					if existing.Variants == nil {
						existing.Variants = make(map[string]*models.Host)
						// Add the existing/winner host as variant for its source
						winnerCopy := *existing
						existing.Variants[existing.Source] = &winnerCopy
					}
					// Add the new shadowed host as variant
					shadowedCopy := host
					existing.Variants[host.Source] = &shadowedCopy
				} else {
					// New host - initialize AvailableIn with its source
					host.AvailableIn = []string{host.Source}
					// Initialize variants
					host.Variants = make(map[string]*models.Host)
					// Store self as variant
					selfCopy := host
					host.Variants[host.Source] = &selfCopy

					hostMap[host.Alias] = &host
				}
			}
		}
	}

	// PRIORITY 3: Termix API (LOWEST PRIORITY)
	if config.Sources.TermixEnabled && config.Termix.Enabled && config.Termix.BaseURL != "" {
		logError("Termix config loaded", fmt.Errorf("baseUrl=%s", config.Termix.BaseURL))

		client := termix.NewClient(config.Termix.BaseURL, config.Termix.JWT, config.Termix.JWTExpiry)

		// Try to fetch hosts without credentials first
		termixHosts, termixFetchErr := client.FetchHosts("", "")

		// Handle auth errors (same as before)
		if termixFetchErr != nil {
			if _, isAuthError := termixFetchErr.(*termix.AuthError); isAuthError {
				return nil, termixFetchErr // Logic for auth flow
			}
			// Log but don't fail everything if Termix fails
			logError("Termix FetchHosts failed", termixFetchErr)
			// Proceed without Termix hosts (or maybe show error in UI?)
			// For now, we'll just log it and continue
		} else {
			logError("Termix hosts fetched successfully", fmt.Errorf("count=%d", len(termixHosts)))

			for _, termixHost := range termixHosts {
				if existing, found := hostMap[termixHost.Alias]; found {
					// Host exists - just add source availability
					existing.AvailableIn = append(existing.AvailableIn, "termix")

					// Add Termix variant
					if existing.Variants == nil {
						existing.Variants = make(map[string]*models.Host)
					}
					// Ensure existing/winner is stored as variant
					if existing.Variants[existing.Source] == nil {
						winnerCopy := *existing
						existing.Variants[existing.Source] = &winnerCopy
					}
					// Add shadowed variant
					shadowedCopy := termixHost
					existing.Variants["termix"] = &shadowedCopy
				} else {
					// New host
					termixHost.AvailableIn = []string{"termix"}
					// Initialize variants
					termixHost.Variants = make(map[string]*models.Host)
					selfCopy := termixHost
					termixHost.Variants["termix"] = &selfCopy

					hostMap[termixHost.Alias] = &termixHost
				}
			}

			// Save JWT if updated
			if client.GetJWT() != config.Termix.JWT || client.GetJWTExpiry() != config.Termix.JWTExpiry {
				config.Termix.JWT = client.GetJWT()
				config.Termix.JWTExpiry = client.GetJWTExpiry()
				SaveConfig(config)
			}
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

	// Sort hosts: favorites first, then alphabetically
	sortHostsByFavorite(config.Hosts)

	return config, nil
}

// sortHostsByFavorite sorts hosts with favorites at the top, then alphabetically by alias
func sortHostsByFavorite(hosts []models.Host) {
	sort.Slice(hosts, func(i, j int) bool {
		// First priority: favorites come first
		if hosts[i].Favorite != hosts[j].Favorite {
			return hosts[i].Favorite
		}
		// Second priority: alphabetical order by alias (case-insensitive)
		return strings.ToLower(hosts[i].Alias) < strings.ToLower(hosts[j].Alias)
	})
}
