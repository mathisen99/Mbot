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

type OpenAIRequestBody struct {
	Model     string                   `json:"model"`
	Messages  []map[string]interface{} `json:"messages"`
	MaxTokens int                      `json:"max_tokens"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs"`
	} `json:"choices"`
}

func OpenAIRequest(message, imageURL, target string) (string, error) {
	color.Cyan(">> OpenAIRequestRaw called with message: %s, imageURL: %s, target: %s", message, imageURL, target)

	// Prepare the user message content
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
		Model:     "gpt-4o",
		Messages:  messages,
		MaxTokens: 300,
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

	color.Cyan(">> Received response from OpenAI: %+v", result)
	if len(result.Choices) > 0 {
		// if lengt is bigger then 420 characters we send it to the paste service
		if len(result.Choices[0].Message.Content) > 420 {
			return PasteService(result.Choices[0].Message.Content)
		}

		// flatten the response to a single string remove any newlines or extra spaces
		result.Choices[0].Message.Content = FlattenMessage(result.Choices[0].Message.Content)

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
