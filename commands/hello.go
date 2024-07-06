package commands

import (
	"mbot/bot"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !hello command
func HelloCommand(connection *ircevent.Connection, sender, target, message string) {
	nickname := bot.ExtractNickname(sender)
	connection.Privmsg(target, "Hello, "+nickname+"!")
}

func init() {
	bot.RegisterCommand("!hello", HelloCommand)
}
