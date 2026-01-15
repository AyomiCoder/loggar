package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Token     string `json:"token"`
	UserEmail string `json:"user_email"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".loggar", "config.json")
}

// SaveToken saves the JWT token and user email to config file
func SaveToken(token, email string) error {
	configPath := GetConfigPath()

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	config := Config{
		Token:     token,
		UserEmail: email,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// LoadToken loads the JWT token and user email from config file
func LoadToken() (*Config, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// ClearToken removes the config file
func ClearToken() error {
	configPath := GetConfigPath()
	return os.Remove(configPath)
}
