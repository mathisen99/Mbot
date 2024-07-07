package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// FetchTitle fetches the page title by mimicking a regular browser's request.
func FetchTitle(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()
	contentType := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("error parsing content type: %v", err)
	}

	if strings.HasPrefix(mediaType, "text/html") {
		opts := []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,            // Skip first run tasks
			chromedp.NoDefaultBrowserCheck, // Disable check for default browser
			chromedp.Headless,              // Run in headless mode
			chromedp.DisableGPU,            // Disable hardware acceleration
			chromedp.IgnoreCertErrors,      // Ignore certificate errors
		}

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		var pageTitle string
		err = chromedp.Run(ctx,
			chromedp.Navigate(finalURL),   // Navigate to the final URL
			chromedp.Sleep(1*time.Second), // Sleep to allow JavaScript execution
			chromedp.Title(&pageTitle),    // Get the page title
		)

		if err != nil {
			return "", fmt.Errorf("error running chromedp tasks: %v", err)
		}

		return pageTitle, nil
	}

	return "", fmt.Errorf("content type is not HTML: %s", contentType)
}
