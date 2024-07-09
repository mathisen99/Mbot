package main

import (
	"log"
	"mbot/bot"
	_ "mbot/commands"
	"os"
)

func main() {
	// Load environment variables
	bot.LoadEnv()

	// Load configuration
	cfg, err := bot.LoadConfig("./data/config.json")
	if err != nil {
		os.Exit(1)
	}

	// Load users
	bot.Users, err = bot.LoadUsersAtStart("./data/users.json")
	if err != nil {
		os.Exit(1)
	}

	// Initialize and start the bot
	b := bot.NewBot(cfg, bot.Users)

	if err := b.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	b.Loop()
}
