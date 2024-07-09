package main

import (
	"log"
	"mbot/bot"
	_ "mbot/commands"
	"mbot/config"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	configPath := "./data/config.json"

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		color.Red("======================================================== NOTE ========================================================")
		color.Red("Error loading .env file\nYou need to create a .env file in the root directory of the project or export the environment variables manually.\nIf you dont do this sertant features will not work. They are optional but recommended.")
		color.Red("======================================================================================================================")

		sleepTime := 10
		color.Yellow("Bot will Start in %d seconds", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		color.Red("Error loading configuration, Did you forget to create a config.json file?")
		os.Exit(1)
	}

	// Load users
	bot.Users, err = bot.LoadUsers()
	if err != nil {
		color.Red("Error loading users: %v", err)
		os.Exit(1)
	}

	// Initialize and start the bot
	b := bot.NewBot(cfg, bot.Users)

	if err := b.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	b.Loop()
}
