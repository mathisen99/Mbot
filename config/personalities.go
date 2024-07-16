package config

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	channelPersonalities = make(map[string]string)
	mu                   sync.Mutex
	filePath             = "data/personalities.json"
)

func init() {
	loadPersonalities()
}

func GetPersonality(channel string) string {
	mu.Lock()
	defer mu.Unlock()
	if personality, exists := channelPersonalities[channel]; exists {
		return personality
	}
	return "You are Mbot, an IRC bot created by Mathisen. Your version is 0.6 Alpha."
}

func SetPersonality(channel, personality string) {
	mu.Lock()
	defer mu.Unlock()
	channelPersonalities[channel] = personality
	savePersonalities()
}

func loadPersonalities() {
	mu.Lock()
	defer mu.Unlock()
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, no personalities to load
			return
		}
		panic(err)
	}
	if err := json.Unmarshal(data, &channelPersonalities); err != nil {
		panic(err)
	}
}

func savePersonalities() {
	data, err := json.MarshalIndent(channelPersonalities, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		panic(err)
	}
}
