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
		connection.Privmsg(target, "Usage: !adduser <nickname> <role> [<channel>]")
		color.Red(">> Invalid command format: %s", message)
		return
	}
	nick := parts[1]
	inputRole := parts[2]
	channel := target
	if len(parts) == 4 {
		channel = parts[3]
	}

	// Map for valid roles to make it case-insensitive for better user experience
	validRoles := map[string]string{
		"owner":   "Owner",
		"admin":   "Admin",
		"trusted": "Trusted",
		"regular": "Regular",
		"badboy":  "BadBoy",
	}

	role, exists := validRoles[strings.ToLower(inputRole)]
	if !exists {
		connection.Privmsg(target, "Invalid role. Valid roles are: Owner, Admin, Trusted, Regular, BadBoy")
		color.Red(">> Invalid role: %s", inputRole)
		return
	}

	bot.WhoisMu.Lock()
	bot.PendingWhois[nick] = func(hostmask string) {
		bot.WhoisMu.Unlock()
		defer bot.WhoisMu.Lock()
		if role == "Owner" {
			for _, user := range users {
				if user.Roles["*"] == "Owner" {
					connection.Privmsg(target, "There is already an Owner. Only one Owner is allowed.")
					color.Red(">> Attempted to add another Owner: %s", nick)
					return
				}
			}
		}

		if existingUser, exists := users[hostmask]; exists {
			if existingUser.Roles["*"] == "Owner" {
				connection.Privmsg(target, fmt.Sprintf("User %s is the Owner and cannot be demoted.", nick))
				color.Red(">> Attempted to demote Owner: %s", nick)
				return
			}

			if existingUserRole, exists := existingUser.Roles[channel]; exists && existingUserRole == role {
				connection.Privmsg(target, fmt.Sprintf("User %s already has the role %s in %s.", nick, role, channel))
				color.Yellow(">> User %s already has role %s in %s", nick, role, channel)
				return
			}

			users[hostmask].Roles[channel] = role
			if err := bot.SaveUsers(users); err != nil {
				connection.Privmsg(target, "Error updating user: "+err.Error())
				color.Red(">> Error updating user: %s", err.Error())
				return
			}

			color.Green(">> User %s updated to role %s in %s", nick, role, channel)
			connection.Privmsg(target, fmt.Sprintf("User %s's role has been updated to %s in %s.", nick, role, channel))
			return
		}

		user := bot.User{Hostmask: hostmask, Roles: map[string]string{channel: role}}
		if err := bot.AddUser(users, user); err != nil {
			connection.Privmsg(target, "Error adding user: "+err.Error())
			color.Red(">> Error adding user: %s", err.Error())
			return
		}

		color.Green(">> User %s added with role %s in %s", nick, role, channel)
		connection.Privmsg(target, fmt.Sprintf("User %s has added %s with role %s in %s.", nickname, nick, role, channel))
	}
	bot.WhoisMu.Unlock()

	connection.SendRaw(fmt.Sprintf("WHOIS %s", nick))
}

// RegisterAddUserCommand registers the !adduser command
func RegisterAddUserCommand() {
	bot.RegisterCommand("!adduser", AddUserCommand)
}
