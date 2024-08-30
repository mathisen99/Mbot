package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func detectImageContent(message, imageURL string) (string, error) {
	color.Cyan(">> detectImageContent called with message: %s and imageURL: %s", message, imageURL)

	// Prepare the user message content using raw approach
	messages := []map[string]interface{}{
		{
			"role": "user",
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": message,
				},
				{
					"type": "image_url",
					"image_url": map[string]string{
						"url": imageURL,
					},
				},
			},
		},
	}

	requestBody := OpenAIRequestBody{
		Model:     "gpt-4o-2024-08-06",
		Messages:  messages,
		MaxTokens: 4000,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		color.Red(">> JSON marshal error: %v", err)
		return "", fmt.Errorf("JSON marshal error: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		color.Red(">> OPENAI_API_KEY environment variable not set")
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		color.Red(">> Error creating HTTP request: %v", err)
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red(">> Error sending HTTP request: %v", err)
		return "", fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var result OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		color.Red(">> Error decoding response: %v", err)
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(result.Choices) > 0 {

		color.Green(">> Returning response from OpenAI: %s", result.Choices[0].Message.Content)
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response received from OpenAI")
}

func FlattenMessage(message string) string {
	// Remove all newline characters
	message = strings.ReplaceAll(message, "\n", " ")
	message = strings.ReplaceAll(message, "\r", " ")

	// Remove all multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	message = re.ReplaceAllString(message, " ")

	// Trim leading and trailing spaces
	message = strings.TrimSpace(message)

	return message
}
