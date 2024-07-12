package bot

import (
	"fmt"
	"strings"

	"mbot/config"

	"github.com/ergochat/irc-go/ircevent"
)

var rateLimiter = NewRateLimiter()
var ConfigData *config.Config

// CommandHandler is a type alias for functions that handle commands
type CommandHandler func(connection *ircevent.Connection, sender, target, message string, users map[string]User)

// Command struct to hold the handler, required role, and group/allowed channels
type Command struct {
	Handler         CommandHandler
	AllowedChannels []string
	RequiredRole    string
}

// Role levels
var roleMapping = map[string]int{
	"everyone": RoleEveryone,
	"badboy":   RoleBadBoy,
	"trusted":  RoleTrusted,
	"admin":    RoleAdmin,
	"owner":    RoleOwner,
}

// Map of commands to their handlers and required roles
var commands = map[string]Command{}

func RegisterCommand(cmd string, handler CommandHandler) {
	for _, group := range ConfigData.CommandGroups {
		for _, command := range group.Commands {
			if command == cmd {
				commands[cmd] = Command{
					Handler:         handler,
					AllowedChannels: group.AllowedChannels,
					RequiredRole:    group.Role,
				}
				return
			}
		}
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

		if !isCommandAllowedInChannel(target, command) {
			connection.Privmsg(target, "This command is not allowed in this channel.")
			return
		}

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

		requiredRoleLevel, ok := roleMapping[command.RequiredRole]
		if !ok {
			connection.Privmsg(target, "Invalid role specified for this command.")
			return
		}

		if userRoleLevel >= requiredRoleLevel {
			command.Handler(connection, sender, target, trimmedMessage, users)
		} else {
			connection.Privmsg(target, "You do not have permission to execute this command.")
		}
	}
}

// Helper function to check if a command is allowed in a channel
func isCommandAllowedInChannel(channel string, command Command) bool {
	for _, allowedChannel := range command.AllowedChannels {
		if allowedChannel == channel {
			return true
		}
	}
	return false
}
