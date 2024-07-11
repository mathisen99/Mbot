package main

import (
	"log"
	"mbot/bot"
	"mbot/commands" // Import the commands package
	"mbot/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load environment variables
	bot.LoadEnv()

	// Load configuration
	cfg, err := config.LoadConfig("./data/config.json")
	if err != nil {
		os.Exit(1)
	}
	bot.ConfigData = cfg

	// Register commands after ConfigData is initialized
	commands.RegisterAllCommands()

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

	// Set up channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to listen for shutdown signals
	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %s\n", sig)
		bot.ShutdownBot(b.Connection.Connection)
	}()

	// Start the bot's main loop
	b.Loop()
}
