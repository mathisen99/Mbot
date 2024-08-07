package commands

import (
	"mbot/bot"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !join command
func JoinCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !join <channel>")
		return
	}
	channel := parts[1]
	connection.Join(channel)
}

// Handler for the !part command
func PartCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !part <channel>")
		return
	}
	channel := parts[1]
	connection.Part(channel)
}

// Handler for the !topic command
func TopicCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 3 {
		connection.Privmsg(target, "Usage: !topic <channel> <new topic>")
		return
	}
	channel := parts[1]
	newTopic := strings.Join(parts[2:], " ")
	connection.Send("TOPIC", channel, newTopic)
}

// Handler for the !nick command
func NickCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !nick <new nickname>")
		return
	}
	newNick := parts[1]
	connection.Send("NICK", newNick)
}

// Handler for the !invite command
func InviteCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 3 {
		connection.Privmsg(target, "Usage: !invite <nickname> <channel>")
		return
	}
	nickname := parts[1]
	channel := parts[2]
	connection.Send("INVITE", nickname, channel)
}

// Handler for the !op command
func OpCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !op <nickname>")
		return
	}
	nickname := parts[1]
	connection.Send("MODE", target, "+o", nickname)
}

// Handler for the !deop command
func DeopCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !deop <nickname>")
		return
	}
	nickname := parts[1]
	connection.Send("MODE", target, "-o", nickname)
}

// Handler for the !voice command
func VoiceCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !voice <nickname>")
		return
	}
	nickname := parts[1]
	connection.Send("MODE", target, "+v", nickname)
}

// Handler for the !devoice command
func DevoiceCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !devoice <nickname>")
		return
	}
	nickname := parts[1]
	connection.Send("MODE", target, "-v", nickname)
}

// Handler for the !kick command
func KickCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !kick <nickname> [reason]")
		return
	}
	nickname := parts[1]
	reason := ""
	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
	}
	connection.Send("KICK", target, nickname, reason)
}

// Handler for the !ban command
func BanCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !ban <nickname>")
		return
	}
	nickname := parts[1]
	connection.Send("MODE", target, "+b", nickname)
}

// Handler for the !unban command
func UnbanCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Usage: !unban <nickname>")
		return
	}
	nickname := parts[1]
	connection.Send("MODE", target, "-b", nickname)
}

// Handler for the !shutdown command
func ShutdownCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	connection.Quit()
	bot.ShutdownBot(connection)
}

// RegisterBaseCommands registers all basic commands
func RegisterBaseCommands() {
	bot.RegisterCommand("!join", JoinCommand)
	bot.RegisterCommand("!part", PartCommand)
	bot.RegisterCommand("!topic", TopicCommand)
	bot.RegisterCommand("!nick", NickCommand)
	bot.RegisterCommand("!invite", InviteCommand)
	bot.RegisterCommand("!op", OpCommand)
	bot.RegisterCommand("!deop", DeopCommand)
	bot.RegisterCommand("!voice", VoiceCommand)
	bot.RegisterCommand("!devoice", DevoiceCommand)
	bot.RegisterCommand("!kick", KickCommand)
	bot.RegisterCommand("!ban", BanCommand)
	bot.RegisterCommand("!unban", UnbanCommand)
	bot.RegisterCommand("!shutdown", ShutdownCommand)
}
