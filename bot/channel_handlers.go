package bot

import (
	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
)

// Function to handle channel messages
func handleChannelMessage(connection *ircevent.Connection, sender, target, message string) {
	color.Cyan(">> Channel message in %s from %s: %s", target, sender, message)

	// Handle commands with the target as the channel
	handleCommand(connection, sender, target, message)
}
