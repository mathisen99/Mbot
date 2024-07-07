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

type Features struct {
	EnableYouTubeCheck    bool `json:"enable_youtube_check"`
	EnableWikipediaCheck  bool `json:"enable_wikipedia_check"`
	EnableGithubCheck     bool `json:"enable_github_check"`
	EnableIMDbCheck       bool `json:"enable_imdb_check"`
	EnableVirusTotalCheck bool `json:"enable_virus_total_check"`
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

// Function to load the features configuration from a file
func LoadFeatures(filePath string) (*Features, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening features file: %w", err)
	}
	defer file.Close()

	features := &Features{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(features); err != nil {
		return nil, fmt.Errorf("error decoding features file: %w", err)
	}

	return features, nil
}

// Function to save the features configuration to a file
func SaveFeatures(features *Features, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating features file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(features); err != nil {
		return fmt.Errorf("error encoding features file: %w", err)
	}

	return nil
}
