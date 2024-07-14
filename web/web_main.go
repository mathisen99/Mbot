package web

import (
	"log"
	"net/http"
	"time"

	mod "mbot/web/backend"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware creates a rate limiter that allows up to maxBurst requests in maxBurst seconds
func RateLimiterMiddleware(maxBurst int, refillTime time.Duration) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(refillTime), maxBurst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}

func StartWebServer() {
	// Start the web server
	r := gin.Default()

	// Apply the rate limiter middleware to all requests
	r.Use(RateLimiterMiddleware(5, time.Second))

	// Serve static files from the web/static directory
	r.Static("/static", "./web/static")

	// Serve the uploads directory (assuming you create this directory)
	r.Static("/web/uploads", "./web/uploads")

	// Routes
	r.POST("/create", mod.HandleCreate)
	r.POST("/pst", mod.HandleCreateSimple)
	r.POST("/images", mod.HandleUploadImage)
	r.GET("/images/:id", mod.HandleViewImage)
	r.GET("/view/:id", mod.HandleView)
	r.GET("/list", mod.HandleListAll)

	// Run HTTPS server using the certificates
	// log.Fatal(r.RunTLS(":8787", "/etc/apache2/ssl/certificate.crt", "/etc/apache2/ssl/private.key"))

	// Run HTTP server for local testing
	log.Fatal(r.Run(":8787"))
}
