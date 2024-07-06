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
