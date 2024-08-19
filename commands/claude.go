package commands

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mbot/bot"
	"net/http"
	"os"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
	"github.com/liushuangls/go-anthropic/v2"
)

// ClaudeCommand handles the !claude command
func ClaudeCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	fmt.Println(">> ClaudeCommand called with sender:", sender, "target:", target, "message:", message)
	// Extract the question from the message
	question := strings.TrimPrefix(message, "!claude ")

	// Create a new Anthropic client
	client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))

	// Call the Claude API
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Dot5Sonnet20240620, // Use Claude 3.5 Sonnet
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(question),
		},
		MaxTokens: 1000,
	})

	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			connection.Privmsg(target, fmt.Sprintf("Claude API error: %s - %s", e.Type, e.Message))
		} else {
			connection.Privmsg(target, fmt.Sprintf("Error calling Claude API: %v", err))
		}
		return
	}

	// Send the response to the IRC channel
	if len(resp.Content) > 0 {
		answer := resp.Content[0].GetText()
		if len(answer) < 450 {
			connection.Privmsg(target, answer)
		} else {

			pasteurl, err := PasteService(answer)
			if err != nil {
				connection.Privmsg(target, "Error sending to paste service.")
				return
			}
			connection.Privmsg(target, pasteurl)
		}
	} else {
		connection.Privmsg(target, "No response from Claude.")
	}
}

// RegisterClaudeCommand registers the !claude command
func RegisterClaudeCommand() {
	bot.RegisterCommand("!claude", ClaudeCommand)
}

func PasteService(content string) (string, error) {
	color.Magenta(">> Sending to paste service...")
	// Load token for the paste service
	token := os.Getenv("VALID_PASTE_TOKEN")
	fmt.Println("Token:", token) // Add this line to print the token and verify it's correct

	// Define the API endpoint
	url := "https://mathizen.net:8787/create"

	// Prepare the request body as JSON
	requestBody, err := json.Marshal(map[string]string{
		"answer": content,
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling request data: %v", err)
	}

	// Create an HTTP client with TLS configuration to skip verification
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to the paste service: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %d, message: %s", resp.StatusCode, responseBody)
	}

	var result map[string]string
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response JSON: %v", err)
	}

	url, ok := result["url"]
	if !ok {
		return "", fmt.Errorf("URL not found in response")
	}

	return url, nil
}
