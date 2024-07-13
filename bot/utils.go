package bot

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"mbot/config"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// Function to gracefully shutdown the bot
func ShutdownBot(connection *ircevent.Connection) {
	color.Red("Shutting down bot...")
	connection.Quit()
	os.Exit(0)
}

// ExtractNickname extracts the nickname from the sender string
func ExtractNickname(sender string) string {
	parts := strings.Split(sender, "!")
	if len(parts) > 0 {
		return parts[0]
	}
	return sender
}

// ExtractHostmask extracts the hostmask from the sender string
func ExtractHostmask(sender string) string {
	parts := strings.Split(sender, "!")
	if len(parts) > 1 {
		userHostParts := strings.Split(parts[1], "@")
		if len(userHostParts) > 1 {
			username := userHostParts[0]
			host := userHostParts[1]
			return fmt.Sprintf("%s@%s", username, host)
		}
	}
	return sender
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

// SplitMessage splits a message into chunks based on the max length
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

// LoadEnv loads environment variables from the .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		color.Red("======================================================== NOTE ========================================================")
		color.Red("Error loading .env file\nYou need to create a .env file in the root directory of the project or export the environment variables manually.\nIf you dont do this certain features will not work. They are optional but recommended.")
		color.Red("======================================================================================================================")

		sleepTime := 10
		color.Yellow("Bot will Start in %d seconds", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

// LoadConfig loads the bot configuration from the specified file path
func LoadConfig(configPath string) (*config.Config, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		color.Red("======================================================== NOTE ========================================================")
		color.Red("Error loading config.json file\nYou need to create a config.json file in the data directory of the project.\nIf you dont do this the bot will not work.\nThere is an example file in the data directory named config_example.json.")
		color.Red("Shutting down bot....")
		color.Red("======================================================================================================================")
		return nil, err
	}
	return cfg, nil
}

// Helper function to check if a command is allowed in a channel
func IsCommandAllowedInChannel(channel string, command Command) bool {
	for _, allowedChannel := range command.AllowedChannels {
		if allowedChannel == channel {
			return true
		}
	}
	return false
}
