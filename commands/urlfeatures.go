package commands

import (
	"fmt"
	"mbot/bot"
	"mbot/config"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !url command
func URLCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	// Extract command arguments
	args := strings.Fields(message)
	if len(args) < 3 {
		connection.Privmsg(target, "Usage: !url <feature> <on|off>")
		return
	}

	feature := args[1]
	state := strings.ToLower(args[2])

	// Validate state
	var newState bool
	switch state {
	case "on":
		newState = true
	case "off":
		newState = false
	default:
		connection.Privmsg(target, "Invalid state. Use 'on' or 'off'.")
		return
	}

	// Update the feature configuration
	switch feature {
	case "youtube":
		bot.URLConfigData.EnableYouTubeCheck = newState
	case "wikipedia":
		bot.URLConfigData.EnableWikipediaCheck = newState
	case "github":
		bot.URLConfigData.EnableGithubCheck = newState
	case "imdb":
		bot.URLConfigData.EnableIMDbCheck = newState
	case "virustotal":
		bot.URLConfigData.EnableVirusTotalCheck = newState
	default:
		connection.Privmsg(target, fmt.Sprintf("Unknown feature: %s", feature))
		return
	}

	// Save the updated configuration
	err := config.SaveURLConfig(bot.URLConfigData, "./data/url_config.json")
	if err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to save configuration: %v", err))
		return
	}

	connection.Privmsg(target, fmt.Sprintf("Feature %s has been turned %s.", feature, state))
}

// RegisterURLCommand registers the !url command
func RegisterURLCommand() {
	bot.RegisterCommand("!url", URLCommand)
}
