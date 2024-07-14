package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

// ProcessResponse processes the response from OpenAI.
func ProcessResponse(ctx context.Context, client *openai.Client, resp *openai.ChatCompletionResponse, req openai.ChatCompletionRequest) (string, error) {
	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		var responses []openai.ChatCompletionMessage
		for _, call := range msg.ToolCalls {
			color.Cyan(">> Function call detected: %s with params: %v", call.Function.Name, call.Function.Arguments)

			var functionArgs map[string]string
			err := json.Unmarshal([]byte(call.Function.Arguments), &functionArgs)
			if err != nil {
				return "", fmt.Errorf("JSON unmarshal error: %v", err)
			}

			var functionResponse string
			switch call.Function.Name {
			case "detect_image_content":
				functionResponse = detectImageContent(functionArgs["image_url"])
			// Add cases for other function calls here
			default:
				functionResponse = "Unknown function call"
			}

			responseMessage := openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    functionResponse,
				Name:       call.Function.Name,
				ToolCallID: call.ID,
			}

			color.Cyan(">> Sending function response back to OpenAI: %v", responseMessage)
			responses = append(responses, responseMessage)
		}

		// Make sure the response messages are part of the conversation history
		req.Messages = append(req.Messages, msg)
		req.Messages = append(req.Messages, responses...)

		resp, err := client.CreateChatCompletion(ctx, req)
		if err != nil || len(resp.Choices) != 1 {
			return "", fmt.Errorf("2nd completion error: %v", err)
		}

		finalMsg := resp.Choices[0].Message
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

		// flatten the response to remove any newlines or extra spaces
		answer = strings.Join(strings.Fields(answer), " ")

		return answer, nil
	}

	return "No function call was made", nil
}
