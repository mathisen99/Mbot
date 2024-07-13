package url_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// GetYouTubeVideoInfo fetches information about a Youtube video using the Youtube API
func GetYouTubeVideoInfo(videoID string, apiKey string) (string, error) {
	color.Magenta(">> Fetching Youtube info... Youtube ID: %s", videoID)
	apiURL := "https://www.googleapis.com/youtube/v3/videos"
	resp, err := http.Get(fmt.Sprintf("%s?part=snippet,contentDetails,statistics&id=%s&key=%s", apiURL, videoID, apiKey))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch video information")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	items, ok := result["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", errors.New("no video information found")
	}

	videoInfo := items[0].(map[string]interface{})

	snippet, snippetOk := videoInfo["snippet"].(map[string]interface{})
	contentDetails, contentDetailsOk := videoInfo["contentDetails"].(map[string]interface{})
	statistics, statisticsOk := videoInfo["statistics"].(map[string]interface{})

	if !snippetOk || !contentDetailsOk || !statisticsOk {
		return "", errors.New("incomplete video information")
	}

	title := getStringValue(snippet, "title")
	channelTitle := getStringValue(snippet, "channelTitle")
	uploadDate := getStringValue(snippet, "publishedAt")
	if len(uploadDate) > 10 {
		uploadDate = uploadDate[:10]
	}

	duration := formatDuration(getStringValue(contentDetails, "duration"))

	views := getStringValue(statistics, "viewCount")
	likes := getStringValue(statistics, "likeCount")

	return fmt.Sprintf("\x02\x0300,01► You\x0304,01Tube\x03\x02 :: %s :: Duration: %s :: Views: %s :: Uploader: %s :: Uploaded: %s :: %s likes",
		title, duration, views, channelTitle, uploadDate, likes), nil
}

// ExtractVideoID extracts the video ID from a Youtube URL
func ExtractVideoID(url string) string {
	if strings.Contains(url, "youtu.be/") {
		parts := strings.Split(url, "youtu.be/")
		if len(parts) > 1 {
			return strings.Split(parts[1], "?")[0]
		}
	}

	if strings.Contains(url, "youtube.com/watch?v=") {
		parts := strings.Split(url, "youtube.com/watch?v=")
		if len(parts) > 1 {
			return strings.Split(parts[1], "&")[0]
		}
	}

	return ""
}

// FormatDuration formats the duration string from the Youtube API
func formatDuration(duration string) string {
	re := regexp.MustCompile(`PT(\d+H)?(\d+M)?(\d+S)?`)
	matches := re.FindStringSubmatch(duration)

	hours := strings.TrimSuffix(getOrDefault(matches, 1), "H")
	minutes := strings.TrimSuffix(getOrDefault(matches, 2), "M")
	seconds := strings.TrimSuffix(getOrDefault(matches, 3), "S")

	if hours == "" {
		hours = "0"
	}
	if minutes == "" {
		minutes = "0"
	}
	if seconds == "" {
		seconds = "0"
	}

	if hours != "0" {
		return fmt.Sprintf("%sh %sm %ss", hours, minutes, seconds)
	}
	return fmt.Sprintf("%sm %ss", minutes, seconds)
}

// getStringValue retrieves a string value from a map
func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

// getOrDefault returns the value at the given index or an empty string if the index is out of bounds
func getOrDefault(matches []string, index int) string {
	if index < len(matches) && matches[index] != "" {
		return matches[index]
	}
	return ""
}
