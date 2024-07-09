package commands

import (
	"mbot/bot"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !hello command
func HelloCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	nickname := bot.ExtractNickname(sender)
	connection.Privmsg(target, "Hello, "+nickname+"!")

	// Print user list
	userList := "Users in the channel: " + bot.GetUserList(users)
	connection.Privmsg(target, userList)
}

func init() {
	bot.RegisterCommand("!hello", HelloCommand)
}
