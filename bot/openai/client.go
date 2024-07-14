package openai

import (
	"context"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

// InitializeClient initializes and returns a new OpenAI client.
func InitializeClient() (*openai.Client, context.Context, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, nil, errors.New("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(openAIKey)
	ctx := context.Background()
	return client, ctx, nil
}
