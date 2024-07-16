package openai

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// ExtractImageURL parses the message to extract an image URL if present
func ExtractImageURL(message string) (string, string) {
	re := regexp.MustCompile(`\bhttps?://\S+\.(jpg|jpeg|png|gif)\b`)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 0 {
		imageURL := matches[0]
		messageWithoutURL := strings.Replace(message, imageURL, "", 1)
		messageWithoutURL = strings.TrimSpace(messageWithoutURL)
		fmt.Println("Image URL found:", imageURL)
		fmt.Println("Message without URL:", messageWithoutURL)
		return messageWithoutURL, imageURL
	}

	fmt.Println("No image URL found")
	return message, ""
}

func PasteService(content string) (string, error) {
	color.Magenta(">> Sending to paste service...")
	// Load token for the paste service
	token := os.Getenv("VALID_PASTE_TOKEN")

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

// detectImageContent uses a regex to detect image content in a URL.
func detectImageContent(imageURL string) string {
	re := regexp.MustCompile(`\bhttps?://\S+\.(jpg|jpeg|png|gif)\b`)
	if re.MatchString(imageURL) {
		return "YES IMAGE DETECTED"
	}
	return "NO IMAGE DETECTED"
}
