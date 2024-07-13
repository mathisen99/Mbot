package commands

import (
	"fmt"
	"mbot/bot"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
)

// Handler for the RemoveUser command
func RemoveUserCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	nickname := bot.ExtractNickname(sender)

	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !deluser <nickname> [<channel>]")
		color.Red(">> Invalid command format: %s", message)
		return
	}
	nick := parts[1]
	channel := target
	if len(parts) == 3 {
		channel = parts[2]
	}

	bot.WhoisMu.Lock()
	bot.PendingWhois[nick] = func(hostmask string) {
		bot.WhoisMu.Unlock()
		defer bot.WhoisMu.Lock()

		if hostmask == "" {
			connection.Privmsg(target, fmt.Sprintf("Could not resolve hostmask for user %s.", nick))
			color.Red(">> Could not resolve hostmask for user: %s", nick)
			return
		}

		if existingUser, exists := users[hostmask]; exists {
			if existingUser.Roles["*"] == "Owner" {
				connection.Privmsg(target, fmt.Sprintf("User %s is the Owner and cannot be removed.", nick))
				color.Red(">> Attempted to remove Owner: %s", nick)
				return
			}

			if _, exists := existingUser.Roles[channel]; exists {
				delete(users[hostmask].Roles, channel)
				if err := bot.SaveUsers(users); err != nil {
					connection.Privmsg(target, "Error removing user: "+err.Error())
					color.Red(">> Error removing user: %s", err.Error())
					return
				}

				color.Green(">> User %s removed from %s", nick, channel)
				connection.Privmsg(target, fmt.Sprintf("User %s has been removed by %s from %s.", nick, nickname, channel))
				return
			}

			connection.Privmsg(target, fmt.Sprintf("User %s does not have any role in %s.", nick, channel))
			color.Yellow(">> User %s does not have any role in %s", nick, channel)
			return
		}

		connection.Privmsg(target, fmt.Sprintf("User %s does not exist.", nick))
		color.Yellow(">> User %s does not exist", nick)
	}
	bot.WhoisMu.Unlock()

	connection.SendRaw(fmt.Sprintf("WHOIS %s", nick))
}

// RegisterRemoveUserCommand registers the !deluser command
func RegisterRemoveUserCommand() {
	bot.RegisterCommand("!deluser", RemoveUserCommand)
}
