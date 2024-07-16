package commands

import (
	"mbot/bot"
	"mbot/config"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

func PersonalityCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	args := strings.SplitN(message, " ", 2)

	if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
		personality := config.GetPersonality(target)
		connection.Privmsg(target, "Current personality for this channel: "+personality)
	} else {
		personality := args[1]
		config.SetPersonality(target, personality)
		connection.Privmsg(target, "Personality for this channel has been set to: "+personality)
	}
}

func RegisterPersonalityCommands() {
	bot.RegisterCommand("!personality", PersonalityCommand)
}
