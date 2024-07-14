package backend

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yuin/goldmark"
)

type Entry struct {
	ID      string
	Content string
	Type    string // "text" or "image"
}

var (
	store = make(map[string]Entry)
	mu    sync.Mutex
	tmpl  = template.Must(template.ParseFiles("./web/template.html")) // Load and parse the template file
)

func HandleCreate(c *gin.Context) {
	var request struct {
		Answer string `json:"answer"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	id := uuid.New().String()
	mu.Lock()
	store[id] = Entry{ID: id, Content: request.Answer, Type: "text"}
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"url": "http://localhost:8787/view/" + id})
}

func HandleView(c *gin.Context) {
	id := c.Param("id")
	mu.Lock()
	entry, ok := store[id]
	mu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	if entry.Type == "image" {
		c.Header("Content-Type", "image/jpeg") // Adjust based on actual image type
		c.File(entry.Content)
		return
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(entry.Content), &buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render response"})
		return
	}

	c.Header("Content-Type", "text/html")
	err := tmpl.Execute(c.Writer, map[string]interface{}{
		"Content": template.HTML(buf.String()),
	})
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
	}
}

// Assuming you have a function to determine the correct URL
func getUrl(entry Entry) string {
	if entry.Type == "image" {
		return "http://localhost:8787/uploads/" + entry.Content
	}
	return "http://localhost:8787/view/" + entry.ID
}

func HandleListAll(c *gin.Context) {
	mu.Lock()
	var pastes []struct {
		URL     string
		Display string
		Type    string
	}
	for _, entry := range store {
		display := entry.Content
		if len(display) > 50 { // Truncate if too long and append ellipsis
			display = display[:50] + "..."
		}
		if entry.Type == "image" {
			// Assuming the Content for images is the file name
			display = "Image: " + filepath.Base(entry.Content) // Show just the file name
		} else {
			display = "Text: " + display // Prefix text pastes with "Text: "
		}
		pastes = append(pastes, struct {
			URL     string
			Display string
			Type    string
		}{
			URL:     getUrl(entry),
			Display: display,
			Type:    entry.Type,
		})
	}
	mu.Unlock()

	c.Header("Content-Type", "text/html")
	err := tmpl.Execute(c.Writer, map[string]interface{}{
		"Pastes": pastes,
	})
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
	}
}

func HandleCreateSimple(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read body"})
		return
	}
	content := string(body)
	id := uuid.New().String()
	mu.Lock()
	store[id] = Entry{ID: id, Content: content, Type: "text"}
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"url": "http://localhost:8787/view/" + id})
}

func HandleUploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get uploaded file"})
		return
	}
	fileID := uuid.New().String()
	fileExtension := filepath.Ext(fileHeader.Filename)
	fileName := fileID + fileExtension
	filePath := "uploads/" + fileName
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	mu.Lock()
	store[fileID] = Entry{ID: fileID, Content: fileName, Type: "image"}
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"url": "http://localhost:8787/uploads/" + fileName})
}

func HandleViewImage(c *gin.Context) {
	imageID := c.Param("imageID")
	filePath := "uploads/" + imageID // Construct the file path where images are stored

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.File(filePath)
}
