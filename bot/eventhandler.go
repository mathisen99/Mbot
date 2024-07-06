package bot

import (
	"sync"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/fatih/color"
)

var once sync.Once

// Function to register event handlers
func RegisterEventHandlers(connection *ircevent.Connection) {
	once.Do(func() {
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
	})
}

// Function to extract the sender from an IRC message
func getSender(e ircmsg.Message) string {
	return e.Source
}

// Function to handle PRIVMSG events (channel and private messages)
func handlePrivmsg(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	target := e.Params[0]
	message := e.Params[1]

	if target[0] == '#' || target[0] == '&' {
		handleChannelMessage(connection, sender, target, message)
	} else {
		handlePrivateMessage(connection, sender, message)
	}
}

// Function to handle private messages
func handleNotice(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Yellow(">> Notice from %s: %s", sender, e.Params[1])
}

// Function to handle channel messages
func handleJoin(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Green(">> %s joined %s", sender, e.Params[0])
}

// Function to handle channel messages
func handlePart(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Red(">> %s parted %s", sender, e.Params[0])
}

// Function to handle channel messages
func handleQuit(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Magenta(">> %s quit", sender)
}

// Function to handle channel messages
func handleKick(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Red(">> %s was kicked from %s by %s: %s", e.Params[1], e.Params[0], sender, e.Params[2])
}

// Function to handle channel messages
func handleBan(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Red(">> %s was banned from %s by %s", e.Params[1], e.Params[0], sender)
}

// Function to handle channel messages
func handleMode(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Blue(">> %s set mode %s on %s", sender, e.Params[1], e.Params[0])
}

// Function to handle channel messages
func handleNick(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Cyan(">> %s is now known as %s", sender, e.Params[0])
}

// Function to handle channel messages
func handleTopic(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Blue(">> %s changed topic on %s to: %s", sender, e.Params[0], e.Params[1])
}

// Function to handle channel messages
func handleInvite(connection *ircevent.Connection, e ircmsg.Message) {
	sender := getSender(e)
	color.Green(">> %s invited %s to %s", sender, e.Params[0], e.Params[1])
}

// Function to handle channel messages
func handleError(connection *ircevent.Connection, e ircmsg.Message) {
	if len(e.Params) > 0 {
		color.Red(">> ERROR: %s", e.Params[0])
	}
}

// Function to handle channel messages
func handlePing(connection *ircevent.Connection, e ircmsg.Message) {
	color.Green(">> Received PING, sending PONG")
	connection.Send("PONG", e.Params[0])
}
