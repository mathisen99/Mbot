package main

import (
	"context"
	"log"
	"mbot/bot"
	"mbot/commands"
	"mbot/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	web "mbot/web"
)

var server *http.Server

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
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Run bot and web server concurrently
	go b.Loop()
	go web.StartWebServer()

	// Block until a signal is received
	sig := <-stopChan
	log.Printf("Received signal: %s. Shutting down...", sig)

	// Gracefully shut down the bot
	bot.ShutdownBot(b.Connection.Connection)

	// Gracefully shut down the web server
	shutdownWebServer()
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
	bot.Users, err = bot.LoadUsers(UserDataPath)
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

// Function to gracefully shut down the web server
func shutdownWebServer() {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to gracefully shut down web server: %v", err)
		}
		log.Println("Web server shut down gracefully")
	}
}
