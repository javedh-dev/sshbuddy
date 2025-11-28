package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"sshbuddy/model"
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

func LoadConfig() (*model.Config, error) {
	path, err := GetDataPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &model.Config{Hosts: []model.Host{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config model.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *model.Config) error {
	path, err := GetDataPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
