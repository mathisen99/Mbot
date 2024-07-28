package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/fatih/color"
)

// searchYouTube searches for YouTube videos using the YouTube API
func searchYouTube(query string) (string, error) {
	var apiKey = os.Getenv("YOUTUBE_API_KEY")

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

// Utility function to get string value from map
func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}
