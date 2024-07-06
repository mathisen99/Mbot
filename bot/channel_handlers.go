package bot

import (
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
)

// Function to handle channel messages
func handleChannelMessage(connection *ircevent.Connection, sender, target, message string) {
	color.Cyan(">> Channel message in %s from %s: %s", target, sender, message)

	// Get the bot's nickname
	botNick := GetBotNickname(connection)

	// Check for URLs in the message and handles them
	urls := FindURLs(message)
	if len(urls) > 0 {
		for _, url := range urls {
			color.Green(">> URL found: %s", url)
			HandleUrl(connection, target, url)
		}
	}

	// Check if the message mentions the bot's nickname
	if strings.Contains(message, botNick) {
		nickname := ExtractNickname(sender)
		response := "Hello, " + nickname + "!"
		connection.Privmsg(target, response)
		return
	}

	// Check for commands
	if strings.HasPrefix(message, "!") {
		handleCommand(connection, sender, target, message)
	}
}
