package web

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Load the valid token from environment variable
var validTokens = map[string]bool{
	os.Getenv("VALID_PASTE_TOKEN"): true,
}

// AuthMiddleware checks for a valid token in the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		if !validTokens[token] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
