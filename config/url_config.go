package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type URLFeatures struct {
	EnableYouTubeCheck    bool `json:"enable_youtube_check"`
	EnableWikipediaCheck  bool `json:"enable_wikipedia_check"`
	EnableGithubCheck     bool `json:"enable_github_check"`
	EnableIMDbCheck       bool `json:"enable_imdb_check"`
	EnableVirusTotalCheck bool `json:"enable_virus_total_check"`
}

// Function to load the URL configuration from a file
func LoadURLConfig(filePath string) (*URLFeatures, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening URL config file: %w", err)
	}
	defer file.Close()

	urlConfig := &URLFeatures{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(urlConfig); err != nil {
		return nil, fmt.Errorf("error decoding URL config file: %w", err)
	}

	return urlConfig, nil
}

// Function to save the URL configuration to a file
func SaveURLConfig(config *URLFeatures, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating URL config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding URL config file: %w", err)
	}

	return nil
}
