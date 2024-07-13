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

// Main function
func main() {
	// Load environment variables
	bot.LoadEnv()

	// Load all configurations
	if err := loadAllConfigs(); err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
		os.Exit(1)
	}

	// Register commands
	registerCommands()

	// Initialize and start the bot
	b := bot.NewBot(bot.ConfigData, bot.Users)
	if err := b.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Handle OS signals
	handleOSSignals(b)

	// Start the bot's main loop
	b.Loop()
}

// ============================================================================================================================
// Everything below this line is here to avoid import cycles and to keep the main function clean and readable
// There may be better ways to handle this, but this is a simple and effective solution for now!. (im also lazy at the moment)
// =============================================================================================================================

// Main helper function to load all configurations
func loadAllConfigs() error {
	var err error

	// Load main configuration
	bot.ConfigData, err = config.LoadConfig(ConfigPath)
	if err != nil {
		return err
	}

	// Load command configuration
	bot.CommandConfigData, err = config.LoadCommandConfig(CommandConfigPath)
	if err != nil {
		return err
	}

	// Load URL configuration
	bot.URLConfigData, err = config.LoadURLConfig(URLConfigPath)
	if err != nil {
		return err
	}

	// Load users
	bot.Users, err = bot.LoadUsersAtStart(UserDataPath)
	if err != nil {
		return err
	}

	return nil
}

// Main helper function to register all commands
func registerCommands() {
	commands.RegisterAllCommands()
	commands.RegisterManageCommand(bot.CommandConfigData, CommandConfigPath)
}

// Main helper function to handle OS signals
func handleOSSignals(b *bot.Bot) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %s\n", sig)
		bot.ShutdownBot(b.Connection.Connection)
	}()
}
