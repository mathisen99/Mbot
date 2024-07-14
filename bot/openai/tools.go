package openai

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// DefineDetectImageContentFunction defines the function for image detection.
func DefineDetectImageContentFunction() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        "detect_image_content",
		Description: "Detects content in the provided image",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"image_url": {
					Type:        jsonschema.String,
					Description: "The URL of the image to analyze",
				},
			},
			Required: []string{"image_url"},
		},
	}
}

// GetTools returns a list of available tools for function calling.
func GetTools() []openai.Tool {
	return []openai.Tool{
		{
			Type:     openai.ToolTypeFunction,
			Function: DefineDetectImageContentFunction(),
		},
		// Add more tools here as needed
	}
}
