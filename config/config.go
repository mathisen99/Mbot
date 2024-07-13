package config

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server       string      `json:"server"`
	Port         string      `json:"port"`
	Nick         string      `json:"nick"`
	Channels     []string    `json:"channels"`
	NickServUser string      `json:"nick_serv_user"`
	NickServPass string      `json:"nick_serv_pass"`
	UseTLS       bool        `json:"use_tls"`
	TLSConfig    *tls.Config `json:"-"`
	Features     Features    `json:"url_features"`
}

type Features struct {
	EnableYouTubeCheck    bool `json:"enable_youtube_check"`
	EnableWikipediaCheck  bool `json:"enable_wikipedia_check"`
	EnableGithubCheck     bool `json:"enable_github_check"`
	EnableIMDbCheck       bool `json:"enable_imdb_check"`
	EnableVirusTotalCheck bool `json:"enable_virus_total_check"`
}

type CommandPermission struct {
	Channels []string `json:"channels"`
	Role     string   `json:"role"`
}

type CommandConfig struct {
	Commands map[string][]CommandPermission `json:"commands"`
}

// Function to load the configuration from a file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	if config.UseTLS {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return config, nil
}

// Function to save the configuration to a file
func SaveConfig(config *Config, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding config file: %w", err)
	}

	return nil
}

// DefaultCommandConfig provides a default configuration
func DefaultCommandConfig() *CommandConfig {
	return &CommandConfig{
		Commands: map[string][]CommandPermission{
			"!managecmd": {{Role: "Owner", Channels: []string{"*"}}},
		},
	}
}

// Helper function to save the command configuration to a file
func saveCommandConfig(filePath string, config *CommandConfig) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating command config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding command config file: %w", err)
	}

	return nil
}

// Function to load the command configuration from a file
func LoadCommandConfig(filePath string) (*CommandConfig, error) {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		// Create the file with default content if it does not exist
		defaultConfig := DefaultCommandConfig()
		err = saveCommandConfig(filePath, defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("error creating default command config file: %w", err)
		}
		return defaultConfig, nil
	} else if err != nil {
		return nil, fmt.Errorf("error opening command config file: %w", err)
	}
	defer file.Close()

	commandConfig := &CommandConfig{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(commandConfig); err != nil {
		return nil, fmt.Errorf("error decoding command config file: %w", err)
	}

	return commandConfig, nil
}
