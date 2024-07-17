package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// Implement the actual image generation function
func createImage(description string) (string, error) {
	color.Yellow(">> Generating image with description: %v", description)
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}

	// Predefined options for image generation
	body := map[string]interface{}{
		"model":           "dall-e-3",
		"prompt":          description,
		"n":               1,
		"size":            "1024x1024",
		"style":           "vivid",
		"quality":         "standard",
		"response_format": "url",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+openAIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	var respData struct {
		Data []struct {
			Url string `json:"url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	if len(respData.Data) == 0 {
		return "", fmt.Errorf("no image data found in response")
	}

	// Download the image
	imagePath, err := downloadImage(respData.Data[0].Url)
	if err != nil {
		return "", err
	}
	// Ensure that the temporary file is deleted after it's no longer needed
	defer os.Remove(imagePath)

	// Open the downloaded image file
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Upload the image
	uploadedImageUrl, err := upload(file)
	if err != nil {
		return "", err
	}

	color.Yellow(">> Image generated: %v", uploadedImageUrl)
	return uploadedImageUrl, nil
}

func downloadImage(url string) (string, error) {
	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create a temporary file with a .png extension
	tempFile, err := os.CreateTemp("", "temp_image_*.png")
	if err != nil {
		return "", err
	}
	// It's important to close the file
	defer tempFile.Close()

	// Copy the response body to the temporary file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", err
	}

	// The caller should be responsible for deleting this file when done
	return tempFile.Name(), nil
}

func upload(image io.Reader) (string, error) {
	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("image", "irc_generated_image.png")
	if err != nil {
		return "", err
	}

	written, err := io.Copy(part, image)
	if err != nil {
		return "", err
	}

	color.Yellow(">> Image uploaded: %v bytes", written)

	writer.WriteField("key", os.Getenv("IMGBB_API_KEY"))

	if err := writer.Close(); err != nil {
		return "", err
	}

	resp, err := http.Post("https://api.imgbb.com/1/upload", writer.FormDataContentType(), buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result["success"].(bool) {
		return result["data"].(map[string]interface{})["url"].(string), nil
	} else {
		return "", fmt.Errorf("upload failed: %s", result["error"].(map[string]interface{})["message"].(string))
	}
}
