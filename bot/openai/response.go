package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

// Global variables for tracking image creation
var imageTimestamps []time.Time
var imageLimit = 5
var timeWindow = 24 * time.Hour
var mutex = &sync.Mutex{}

// ProcessResponse processes the response from OpenAI.
func ProcessResponse(ctx context.Context, client *openai.Client, resp *openai.ChatCompletionResponse, req openai.ChatCompletionRequest) (string, error) {
	msg := resp.Choices[0].Message
	if msg.FunctionCall != nil {
		color.Cyan(">> Function call detected: %s with params: %v", msg.FunctionCall.Name, msg.FunctionCall.Arguments)

		var functionArgs map[string]string
		err := json.Unmarshal([]byte(msg.FunctionCall.Arguments), &functionArgs)
		if err != nil {
			return "", fmt.Errorf("JSON unmarshal error: %v", err)
		}

		var functionResponse string
		switch msg.FunctionCall.Name {
		case "detect_image_content":
			functionResponse, err = detectImageContent(functionArgs["message"], functionArgs["image_url"])
			if err != nil {
				return "", fmt.Errorf("error detecting image content: %v", err)
			}
		case "create_image":
			if !canCreateImage() {
				return "Image creation limit reached. Please try again later.", nil
			}

			functionResponse, err = createImage(functionArgs["description"])
			if err != nil {
				return "", fmt.Errorf("error generating image: %v", err)
			}
			// Ensure the image URL is included in the response
			functionResponse = fmt.Sprintf("Here’s an image representing your request: %s", functionResponse)

			recordImageCreation()
		case "check_weather":
			location := functionArgs["location"]
			city := extractCity(location)
			functionResponse = checkWeather(city)
		case "search_youtube":
			query := functionArgs["query"]
			functionResponse, err = searchYouTube(query)
			if err != nil {
				return "", fmt.Errorf("error searching YouTube: %v", err)
			}
		default:
			functionResponse = "Unknown function call"
		}

		color.Cyan(">> Function response: %v", functionResponse)

		responseMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: functionResponse,
			Name:    msg.FunctionCall.Name,
		}

		color.Cyan(">> Sending function response back to OpenAI: %v", responseMessage)

		// Update the conversation with the function response
		req.Messages = append(req.Messages, msg)
		req.Messages = append(req.Messages, responseMessage)

		newResp, err := client.CreateChatCompletion(ctx, req)
		if err != nil || len(newResp.Choices) != 1 {
			return "", fmt.Errorf("2nd completion error: %v", err)
		}

		finalMsg := newResp.Choices[0].Message
		color.Cyan(">> Final answer received: %v", finalMsg)
		answer := finalMsg.Content

		// Check answer length if it's too long meaning more than 420 characters then we send the answer to the paste service
		color.Cyan(">> Final answer length: %d", len(answer))
		if len(answer) > 420 {
			pasteURL, err := PasteService(answer)
			if err != nil {
				return "", fmt.Errorf("error calling PasteService: %v", err)
			}
			return pasteURL, nil
		}

		// Flatten the response to remove any newlines or extra spaces
		answer = strings.Join(strings.Fields(answer), " ")

		return answer, nil
	}

	return "No function call was made", nil
}

// canCreateImage checks if an image can be created based on the limit and time window.
func canCreateImage() bool {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()
	// Filter out timestamps older than the time window
	newTimestamps := []time.Time{}
	for _, ts := range imageTimestamps {
		if now.Sub(ts) < timeWindow {
			newTimestamps = append(newTimestamps, ts)
		}
	}
	imageTimestamps = newTimestamps

	// Check if the limit is reached
	return len(imageTimestamps) < imageLimit
}

// recordImageCreation records the timestamp of a new image creation.
func recordImageCreation() {
	mutex.Lock()
	defer mutex.Unlock()

	imageTimestamps = append(imageTimestamps, time.Now())
}

// extractCity extracts the city name from a location string.
func extractCity(location string) string {
	// Split the location by comma and return the first part (assumed to be the city name)
	parts := strings.Split(location, ",")
	return strings.TrimSpace(parts[0])
}
