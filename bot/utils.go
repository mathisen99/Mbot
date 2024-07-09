package bot

import (
	"encoding/json"
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

// LoadUsers loads the users from the specified file path
func LoadUsersAtStart(filePath string) (map[string]User, error) {
	users, err := LoadUsersFromFile(filePath)
	if err != nil {
		color.Red("======================================================== NOTE ========================================================")
		color.Red("Error loading users.json file\nYou need to create a users.json file in the data directory of the project.\nIf you dont do this the bot will not work.\nThere is an example file in the data directory named users_example.json.")
		color.Red("Shutting down bot....")
		color.Red("======================================================================================================================")
		return nil, err
	}
	return users, nil
}

// LoadUsersFromFile reads the users from a JSON file
func LoadUsersFromFile(filePath string) (map[string]User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening users file: %w", err)
	}
	defer file.Close()

	var users []User
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("error decoding users file: %w", err)
	}

	userMap := make(map[string]User)
	for _, user := range users {
		userMap[user.Hostmask] = user
	}

	return userMap, nil
}
