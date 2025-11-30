package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"sshbuddy/internal/termix"
	"sshbuddy/pkg/models"
)

func GetDataPath() (string, error) {
	// Use XDG_CONFIG_HOME if set, otherwise default to ~/.config
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	
	// Create sshbuddy config directory
	sshbuddyDir := filepath.Join(configDir, "sshbuddy")
	if err := os.MkdirAll(sshbuddyDir, 0755); err != nil {
		return "", err
	}
	
	return filepath.Join(sshbuddyDir, "config.json"), nil
}

// LoadConfig is now in sources.go

func SaveConfig(config *models.Config) error {
	path, err := GetDataPath()
	if err != nil {
		return err
	}

	// Only save manual hosts (not SSH config or termix hosts)
	// But save favorites for all hosts
	saveConfig := &models.Config{
		Theme:     config.Theme,
		Sources:   config.Sources,
		Termix:    config.Termix,
		SSH:       config.SSH,
		Hosts:     []models.Host{},
		Favorites: make(map[string]bool),
	}
	
	for _, host := range config.Hosts {
		if host.Source != "ssh-config" && host.Source != "termix" {
			saveConfig.Hosts = append(saveConfig.Hosts, host)
		}
		// Save favorite status for all hosts (including external sources)
		if host.Favorite {
			saveConfig.Favorites[host.Alias] = true
		}
	}

	data, err := json.MarshalIndent(saveConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}



// logError logs errors to a debug file for troubleshooting
func logError(context string, err error) {
	logPath := "/tmp/sshbuddy-debug.log"
	
	logFile, fileErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		return // Silently fail if we can't log
	}
	defer logFile.Close()
	
	timestamp := fmt.Sprintf("[%s]", os.Getenv("USER"))
	logLine := fmt.Sprintf("%s %s: %v\n", timestamp, context, err)
	logFile.WriteString(logLine)
}



// LoadConfigRaw loads the config file without fetching external sources (SSH config, Termix)
func LoadConfigRaw() (*models.Config, error) {
	path, err := GetDataPath()
	if err != nil {
		return nil, err
	}

	var config models.Config
	
	// Load config from file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Initialize with defaults
		config = models.Config{
			Hosts: []models.Host{},
			Sources: models.SourcesConfig{
				SSHBuddyEnabled:  true,
				SSHConfigEnabled: true,
				TermixEnabled:    false,
			},
			Termix: models.TermixConfig{
				Enabled: false,
			},
			SSH: models.SSHConfig{
				Enabled: true,
			},
		}
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	}
	
	return &config, nil
}

// AuthenticateTermix authenticates with Termix using provided credentials and updates the config
func AuthenticateTermix(username, password string) error {
	// Load config without fetching Termix hosts to avoid circular dependency
	config, err := LoadConfigRaw()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	if !config.Termix.Enabled || config.Termix.BaseURL == "" {
		return fmt.Errorf("termix is not enabled or baseUrl is not configured")
	}
	
	client := termix.NewClient(config.Termix.BaseURL, "", 0)
	jwt, expiry, err := client.Authenticate(username, password)
	if err != nil {
		return err
	}
	
	// Update config with new token and expiry
	config.Termix.JWT = jwt
	config.Termix.JWTExpiry = expiry
	
	return SaveConfig(config)
}
