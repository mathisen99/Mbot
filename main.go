package main

import (
	"log"
	"mbot/bot"
	"mbot/commands"
	"mbot/config"
	"os"
	"os/signal"
	"syscall"
)

const (
	// Config paths
	ConfigPath        = "./data/config.json"
	CommandConfigPath = "./data/command_permissions.json"
	UserDataPath      = "./data/users.json"
	URLConfigPath     = "./data/url_config.json"
)

func main() {
	// Load environment variables
	bot.LoadEnv()

	// Load main configuration
	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		os.Exit(1)
	}
	bot.ConfigData = cfg

	// Load command configuration
	cmdCfg, err := config.LoadCommandConfig(CommandConfigPath)
	if err != nil {
		os.Exit(1)
	}
	bot.CommandConfigData = cmdCfg

	// Load URL configuration
	urlCfg, err := config.LoadURLConfig(URLConfigPath)
	if err != nil {
		os.Exit(1)
	}
	bot.URLConfigData = urlCfg

	// Register commands after ConfigData is initialized
	commands.RegisterAllCommands()
	commands.RegisterManageCommand(cmdCfg, CommandConfigPath)

	// Load users
	bot.Users, err = bot.LoadUsersAtStart(UserDataPath)
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
