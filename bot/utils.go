package bot

import (
	"regexp"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

// ExtractNickname extracts the nickname from the full sender string
func ExtractNickname(fullSender string) string {
	if idx := strings.Index(fullSender, "!"); idx != -1 {
		return fullSender[:idx]
	}
	return fullSender
}

// GetBotNickname retrieves the bot's current nickname
func GetBotNickname(connection *ircevent.Connection) string {
	return connection.Nick
}

// FindURLs finds URLs in a given message
func FindURLs(message string) []string {
	urlRegex := `(https?://[^\s]+|http?://[^\s]+|www\.[^\s]+)`
	re := regexp.MustCompile(urlRegex)
	return re.FindAllString(message, -1)
}
