package bot

import (
	"fmt"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

var rateLimiter = NewRateLimiter()

// CommandHandler is a type alias for functions that handle commands
type CommandHandler func(connection *ircevent.Connection, sender, target, message string, users map[string]User)

// Command struct to hold the handler and required role
type Command struct {
	Handler      CommandHandler
	RequiredRole int
}

// Map of commands to their handlers and required roles
var commands = map[string]Command{}

// Function to register a command
func RegisterCommand(cmd string, handler CommandHandler, requiredRole int) {
	commands[cmd] = Command{
		Handler:      handler,
		RequiredRole: requiredRole,
	}
}

// Function to handle commands
func handleCommand(connection *ircevent.Connection, sender, target, message string, users map[string]User) {
	trimmedMessage := strings.TrimSpace(message)
	parts := strings.Fields(trimmedMessage)

	if len(parts) == 0 {
		return
	}
	cmd := parts[0]
	if command, exists := commands[cmd]; exists {
		nickname := ExtractNickname(sender)
		hostmask := ExtractHostmask(sender)
		userRoleLevel := GetUserRoleLevel(users, hostmask)

		// Check rate limiter
		if !rateLimiter.AllowCommand(nickname) {
			if remaining := rateLimiter.GetCooldownRemaining(nickname); remaining > 0 {
				connection.Privmsg(target, fmt.Sprintf("You are currently in cooldown for %s. Please wait before sending more commands.", FormatDuration(remaining)))
			} else if remaining := rateLimiter.GetShutdownRemaining(nickname); remaining > 0 {
				if rateLimiter.CanSendSuspensionMessage(nickname) {
					connection.Privmsg(target, fmt.Sprintf("You have been temporarily suspended for %s for not reading the warning. Please wait and try again later.", FormatDuration(remaining)))
				}
			}
			return
		}

		if userRoleLevel == RoleBadBoy {
			connection.Privmsg(target, "You do not have permission to execute this command.")
			return
		}
		if userRoleLevel >= command.RequiredRole {
			command.Handler(connection, sender, target, trimmedMessage, users)
		} else {
			connection.Privmsg(target, "You do not have permission to execute this command.")
		}
	}
}
