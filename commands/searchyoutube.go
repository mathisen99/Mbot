package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"mbot/bot"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
)

// SearchYouTube searches for YouTube videos using the YouTube API
func SearchYouTube(query string, apiKey string) (string, error) {
	color.Magenta(">> Searching YouTube for: %s", query)
	apiURL := "https://www.googleapis.com/youtube/v3/search"
	resp, err := http.Get(fmt.Sprintf("%s?part=snippet&type=video&q=%s&key=%s", apiURL, url.QueryEscape(query), apiKey))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to search for videos")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	items, ok := result["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", errors.New("no videos found")
	}

	video := items[0].(map[string]interface{})
	id := video["id"].(map[string]interface{})
	snippet := video["snippet"].(map[string]interface{})

	videoID := getStringValue(id, "videoId")
	title := getStringValue(snippet, "title")
	channelTitle := getStringValue(snippet, "channelTitle")

	return fmt.Sprintf("\x02\x0300,01► You\x0304,01Tube\x03\x02 :: %s :: Channel: %s :: https://www.youtube.com/watch?v=%s",
		title, channelTitle, videoID), nil
}

// Handler for the !yt command
func YTCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	// Extract the search query from the message
	query := strings.TrimSpace(strings.TrimPrefix(message, "!yt"))
	if query == "" {
		connection.Privmsg(target, "Please provide a search query.")
		return
	}

	// Call the YouTube search function
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	result, err := SearchYouTube(query, apiKey)
	if err != nil {
		connection.Privmsg(target, "Error: "+err.Error())
		return
	}

	// Send the search result to the channel
	connection.Privmsg(target, result)
}

// RegisterYTCommand registers the !yt command
func RegisterYTCommand() {
	bot.RegisterCommand("!yt", YTCommand)
}

// Utility function to get string value from map
func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}
