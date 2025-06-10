package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	config := &Config{}

	filePath, err := configFilePath()
	if err != nil {
		return config, fmt.Errorf("error getting config file path: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		log.Fatalf("Failed to unmarshal data from config file: %v", err)
	}

	return config, nil
}

func (c *Config) SetUser(name string) {
	c.CurrentUserName = name
	c.write()
}

func (c *Config) String() string {
	s := fmt.Sprintf("Config\n")
	s = s + fmt.Sprintf("  db_url: %s\n", c.DbUrl)
	s = s + fmt.Sprintf("  current_user_name: %s\n", c.CurrentUserName)

	return s
}

func (c *Config) write() error {
	filePath, err := configFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file path: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create or truncate config file %s: %w", filePath, err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	os.WriteFile(filePath, data, 0644)

	return nil
}

func configFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %v", err)
	}

	filePath := homeDir + "/" + CONFIG_FILE_NAME

	return filePath, nil
}
