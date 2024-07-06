package main

import (
	"fmt"
	"log"
	"mbot/bot"
	_ "mbot/commands" // Import the commands package to register the commands
	"mbot/config"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	configPath := "./data/config.json"

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		color.Red("Error loading configuration, Did you forget to create a config.json file?")
		os.Exit(1)
	}

	// Initialize and start the bot
	b := bot.NewBot(cfg)

	// Register the event handlers
	bot.RegisterEventHandlers(b.Connection)

	if err := b.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	b.Loop()
}
