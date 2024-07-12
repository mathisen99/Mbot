package commands

import (
	"mbot/bot"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !hello command
func HelloCommand2(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	nickname := bot.ExtractNickname(sender)
	connection.Privmsg(target, "Hello, "+nickname+"!")

	// Print user list
	userList := "Users in the channel: " + bot.GetUserList(users)
	connection.Privmsg(target, userList)
}

// RegisterHelloCommand registers the !hello command
func RegisterHello2Command() {
	bot.RegisterCommand("!hello2", HelloCommand)
}
