package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ServerConfig struct {
	Host        string    `json:"host"`
	Port        string    `json:"port"`
	Username    string    `json:"username"`
	Database    string    `json:"database"`
	Password    string    `json:"password"`
	LastUpdated time.Time `json:"last_updated"`
}

func SaveServerConfig(serverName string, config ServerConfig) error {
	configDir := GetConfigDir()
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	filePath := filepath.Join(configDir, serverName+".json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func LoadServerConfig(serverName string) (ServerConfig, error) {
	var config ServerConfig

	filePath := filepath.Join(GetConfigDir(), serverName+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return config, nil
}

func GetConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".config/do-my-job"
	}
	return filepath.Join(homeDir, ".config", "do-my-job")
}
