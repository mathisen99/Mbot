package bot

import (
	"strings"

	"github.com/fatih/color"
)

// Function to handle channel messages
func handleChannelMessage(connection *Connection, sender, target, message string, users map[string]User) {
	color.Cyan(">> Channel message in %s from %s: %s", target, sender, message)

	// Get the bot's nickname
	botNick := GetBotNickname(connection.Connection)

	// Check for URLs in the message and handles them
	urls := FindURLs(message)
	if len(urls) > 0 {
		for _, url := range urls {
			color.Green(">> URL found: %s", url)
			HandleUrl(connection, sender, target, url)
		}
	}

	// Check for commands
	if strings.HasPrefix(message, "!") {
		handleCommand(connection.Connection, sender, target, message, users)
	}

	// Check if the message mentions the bot's nickname
	if strings.Contains(message, botNick) {
		nickname := ExtractNickname(sender)
		response := "Hello, " + nickname + "!"
		connection.Privmsg(target, response)
		return
	}
}
