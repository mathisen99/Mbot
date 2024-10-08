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

// DefineCreateImageFunction defines the function for creating images.
func DefineCreateImageFunction() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        "create_image",
		Description: "Creates an image based on the provided description",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"description": {
					Type:        jsonschema.String,
					Description: "The description of the image to create",
				},
			},
			Required: []string{"description"},
		},
	}
}

// DefineCheckWeatherFunction defines the function for checking the weather.
func DefineCheckWeatherFunction() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        "check_weather",
		Description: "Checks the weather for a specified location",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"location": {
					Type:        jsonschema.String,
					Description: "The location to check the weather for",
				},
			},
			Required: []string{"location"},
		},
	}
}

// DefineSearchYouTubeFunction defines the function for searching YouTube.
func DefineSearchYouTubeFunction() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        "search_youtube",
		Description: "Searches YouTube for videos based on the provided query",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"query": {
					Type:        jsonschema.String,
					Description: "The search query to use for the YouTube search",
				},
			},
			Required: []string{"query"},
		},
	}
}

// DefineSummarizeWebpageFunction defines the function for summarizing a webpage.
func DefineSummarizeWebpageFunction() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        "summarize_webpage",
		Description: "Summarizes the content of the provided webpage URL",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"url": {
					Type:        jsonschema.String,
					Description: "The URL of the webpage to summarize",
				},
			},
			Required: []string{"url"},
		},
	}
}

// GetTools returns a list of available tools for function calling.
func GetTools() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		*DefineDetectImageContentFunction(),
		*DefineCreateImageFunction(),
		*DefineCheckWeatherFunction(),
		*DefineSearchYouTubeFunction(),
		*DefineSummarizeWebpageFunction(), // Added the new summarize_webpage function
	}
}
