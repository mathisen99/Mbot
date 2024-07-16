package commands

import (
	"mbot/bot"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

func MemoryWipeCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	args := strings.SplitN(message, " ", 2)
	if len(args) >= 2 && strings.TrimSpace(args[1]) == "wipe" {
		bot.WipeUserMemory(sender)
		connection.Privmsg(target, "Your memory has been wiped. I will no longer remember our conversation.")
	} else {
		connection.Privmsg(target, "Usage: !memory wipe")
	}
}

func RegisterMemoryWipeCommand() {
	bot.RegisterCommand("!memory", MemoryWipeCommand)
}
