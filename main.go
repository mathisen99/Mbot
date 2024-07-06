package main

import (
	"log"
	"mbot/bot"
	"mbot/config"
)

func main() {
	configPath := "./data/config.json"

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize and start the bot
	b := bot.NewBot(cfg)
	if err := b.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	b.Loop()
}
