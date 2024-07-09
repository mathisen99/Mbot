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

// ExtractHostmask extracts the hostmask from the sender string
func ExtractHostmask(sender string) string {
	// sender is in the format "nickname!username@hostmask"
	parts := strings.Split(sender, "!")
	if len(parts) < 2 {
		return ""
	}
	hostParts := strings.Split(parts[1], "@")
	if len(hostParts) < 2 {
		return ""
	}
	return hostParts[1]
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

// splitMessage splits a message into chunks based on the max length
func SplitMessage(message string, maxLength int) []string {
	var chunks []string

	for len(message) > maxLength {
		cutIndex := strings.LastIndex(message[:maxLength], " ")
		if cutIndex == -1 {
			cutIndex = maxLength
		}
		chunks = append(chunks, message[:cutIndex])
		message = message[cutIndex:]
	}

	chunks = append(chunks, message)
	return chunks
}
