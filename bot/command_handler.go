package bot

import (
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

// CommandHandler is a type alias for functions that handle commands
type CommandHandler func(connection *ircevent.Connection, sender, target, message string, users map[string]User)

// Map of commands to their handlers
var commands = map[string]CommandHandler{}

// Function to register a command
func RegisterCommand(cmd string, handler CommandHandler) {
	commands[cmd] = handler
}

// Function to handle commands
func handleCommand(connection *ircevent.Connection, sender, target, message string, users map[string]User) {
	trimmedMessage := strings.TrimSpace(message)
	parts := strings.Fields(trimmedMessage)
	if len(parts) == 0 {
		return
	}
	cmd := parts[0]
	if handler, exists := commands[cmd]; exists {
		handler(connection, sender, target, trimmedMessage, users)
	}
}
