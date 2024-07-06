package bot

import (
	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/fatih/color"
)

// RegisterEventHandlers registers event handlers for the bot
func RegisterEventHandlers(connection *ircevent.Connection) {
	eventHandlers := map[string]func(*ircevent.Connection, ircmsg.Message){
		"PRIVMSG": handlePrivmsg,
		"NOTICE":  handleNotice,
		"JOIN":    handleJoin,
		"PART":    handlePart,
		"QUIT":    handleQuit,
		"KICK":    handleKick,
		"BAN":     handleBan,
		"MODE":    handleMode,
		"NICK":    handleNick,
		"TOPIC":   handleTopic,
		"INVITE":  handleInvite,
		"ERROR":   handleError,
		"PING":    handlePing,
	}

	for event, handler := range eventHandlers {
		connection.AddCallback(event, func(e ircmsg.Message) {
			handler(connection, e)
		})
	}
}

// getSender returns the sender of the message
func getSender(e ircmsg.Message) string {
	return e.Source
}

// Function to handle PRIVMSG events (channel and private messages)
func handlePrivmsg(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	target := e.Params[0]
	message := e.Params[1]

	if target[0] == '#' || target[0] == '&' {
		// Channel message
		color.Cyan(">> Channel message in %s from %s: %s", target, sender, message)
	} else {
		// Private message
		color.Magenta(">> Private message from %s: %s", sender, message)
	}
}

// Function to handle NOTICE events
func handleNotice(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Yellow(">> Notice from %s: %s", sender, e.Params[1])
}

// Function to handle JOIN events
func handleJoin(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	channel := e.Params[0]
	color.Green(">> %s joined %s", sender, channel)
}

// Function to handle PART events
func handlePart(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	channel := e.Params[0]
	color.Red(">> %s parted %s", sender, channel)
}

// Function to handle QUIT events
func handleQuit(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Magenta(">> %s quit", sender)
}

// Function to handle KICK events
func handleKick(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	channel := e.Params[0]
	kickedUser := e.Params[1]
	reason := e.Params[2]
	color.Red(">> %s was kicked from %s by %s: %s", kickedUser, channel, sender, reason)
}

// Function to handle BAN events
func handleBan(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	channel := e.Params[0]
	bannedUser := e.Params[1]
	color.Red(">> %s was banned from %s by %s", bannedUser, channel, sender)
}

// Function to handle MODE events
func handleMode(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	target := e.Params[0]
	mode := e.Params[1]
	color.Blue(">> %s set mode %s on %s", sender, mode, target)
}

// Function to handle NICK events
func handleNick(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	newNick := e.Params[0]
	color.Cyan(">> %s is now known as %s", sender, newNick)
}

// Function to handle TOPIC events
func handleTopic(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	channel := e.Params[0]
	topic := e.Params[1]
	color.Blue(">> %s changed topic on %s to: %s", sender, channel, topic)
}

// Function to handle INVITE events
func handleInvite(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	invitedUser := e.Params[0]
	channel := e.Params[1]
	color.Green(">> %s invited %s to %s", sender, invitedUser, channel)
}

// Function to handle ERROR events
func handleError(connection *ircevent.Connection, e ircmsg.Message) {
	if len(e.Params) > 0 {
		color.Red(">> ERROR: %s", e.Params[0])
	}
}

// Function to handle PING events
func handlePing(connection *ircevent.Connection, e ircmsg.Message) {
	color.Green(">> Received PING, sending PONG")
	connection.Send("PONG", e.Params[0])
}
