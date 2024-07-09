package commands

import (
	"fmt"
	"mbot/bot"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
)

// Handler for the AddUser command
func AddUserCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	nickname := bot.ExtractNickname(sender)

	parts := strings.Fields(message)
	if len(parts) < 3 {
		connection.Privmsg(target, "Usage: !adduser <nickname> <role>")
		color.Red(">> Invalid command format: %s", message)
		return
	}
	nick := parts[1]
	role := parts[2]

	if _, exists := bot.UserRoles[role]; !exists {
		connection.Privmsg(target, "Invalid role. Valid roles are: Owner, Admin, Trusted, Regular, BadBoy")
		color.Red(">> Invalid role: %s", role)
		return
	}

	bot.WhoisMu.Lock()
	bot.PendingWhois[nick] = func(hostmask string) {
		bot.WhoisMu.Unlock()
		defer bot.WhoisMu.Lock()
		if role == "Owner" {
			for _, user := range users {
				if user.Role == "Owner" {
					connection.Privmsg(target, "There is already an Owner. Only one Owner is allowed.")
					color.Red(">> Attempted to add another Owner: %s", nick)
					return
				}
			}
		}

		if existingUser, exists := users[hostmask]; exists {
			if existingUser.Role == "Owner" {
				connection.Privmsg(target, fmt.Sprintf("User %s is the Owner and cannot be demoted.", nick))
				color.Red(">> Attempted to demote Owner: %s", nick)
				return
			}

			if existingUser.Role == role {
				connection.Privmsg(target, fmt.Sprintf("User %s already has the role %s.", nick, role))
				color.Yellow(">> User %s already has role %s", nick, role)
				return
			}

			users[hostmask] = bot.User{Hostmask: hostmask, Role: role}
			if err := bot.SaveUsers(users); err != nil {
				connection.Privmsg(target, "Error updating user: "+err.Error())
				color.Red(">> Error updating user: %s", err.Error())
				return
			}

			color.Green(">> User %s updated to role %s", nick, role)
			connection.Privmsg(target, fmt.Sprintf("User %s's role has been updated to %s.", nick, role))
			return
		}

		user := bot.User{Hostmask: hostmask, Role: role}
		if err := bot.AddUser(users, user); err != nil {
			connection.Privmsg(target, "Error adding user: "+err.Error())
			color.Red(">> Error adding user: %s", err.Error())
			return
		}

		color.Green(">> User %s added with role %s", nick, role)
		connection.Privmsg(target, fmt.Sprintf("User %s has added %s with role %s.", nickname, nick, role))
	}
	bot.WhoisMu.Unlock()

	connection.SendRaw(fmt.Sprintf("WHOIS %s", nick))
}

// Register the command
func init() {
	bot.RegisterCommand("!adduser", AddUserCommand)
}
