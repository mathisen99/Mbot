package openai

import (
	"context"
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

// Define the struct for the function call response
type FunctionCall struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

// Function to send the main question to OpenAI and get the order of tools
func getToolOrder(ctx context.Context, client *openai.Client, question string) ([]FunctionCall, error) {
	request := openai.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	var functionCalls []FunctionCall
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &functionCalls)
	if err != nil {
		return nil, err
	}

	return functionCalls, nil
}
