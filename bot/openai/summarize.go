package openai

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

// summarizeWebpage summarizes the content of a webpage.
func summarizeWebpage(url string) (string, error) {
	// Fetch the webpage content
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching the webpage: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received non-200 status code %d", resp.StatusCode)
	}

	// Parse the webpage content using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing the webpage: %v", err)
	}

	// Extract the plain text content
	var contentBuilder strings.Builder
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			contentBuilder.WriteString(text + " ")
		}
	})

	content := contentBuilder.String()
	if content == "" {
		return "", fmt.Errorf("error: no content extracted from the webpage")
	}

	// Initialize the OpenAI client
	client, ctx, err := InitializeClient()
	if err != nil {
		color.Red(err.Error())
		return "", fmt.Errorf("error initializing OpenAI client: %v", err)
	}

	// Prepare the system and user messages for the OpenAI request
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: `You are a helpful assistant. Summarize the following webpage content concisely.`,
	}

	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}

	// Prepare the OpenAI API request
	req := openai.ChatCompletionRequest{
		Model:     "gpt-4o-2024-08-06",
		MaxTokens: 3000, // Adjust based on the desired summary length
		Messages:  []openai.ChatCompletionMessage{systemMessage, userMessage},
	}

	color.Cyan(">> Sending request to OpenAI to summarize webpage content")

	// Send the request to OpenAI and capture the response
	respAI, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error summarizing the content: %v", err)
	}

	// Extract the summary from the response
	summary := respAI.Choices[0].Message.Content

	return strings.TrimSpace(summary), nil
}
