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

func Read() *Config {
	filePath := configFilePath()

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %v", filePath, err)
	}

	config := Config{}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to unmarshal data from config file: %v", err)
	}

	return &config
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

func (c *Config) write() {
	filePath := configFilePath()

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create or truncate filepath %s: %v", filePath, err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	os.WriteFile(filePath, data, 0644)
}

func configFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	filePath := homeDir + "/" + CONFIG_FILE_NAME

	return filePath
}
