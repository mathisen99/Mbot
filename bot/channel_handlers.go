package bot

import (
	"strings"

	"github.com/fatih/color"
)

func handleChannelMessage(connection *Connection, sender, target, message string, users map[string]User) {
	color.Cyan(">> Channel message in %s from %s: %s", target, sender, message)

	botNick := GetBotNickname(connection.Connection)

	if strings.HasPrefix(message, "!") {
		handleCommand(connection.Connection, sender, target, message, users)
		return
	}

	if TriviaStateInstance.Active {
		checkTriviaAnswer(sender, message, target, connection)
	}

	if strings.Contains(message, botNick) {
		CallOpenAI(connection, sender, target, message)
		return
	}

	urls := FindURLs(message)
	if len(urls) > 0 {
		for _, url := range urls {
			color.Green(">> URL found: %s", url)
			HandleUrl(connection, sender, target, url)
		}
	}
}
