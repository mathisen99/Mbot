package internal

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

// fetchTitle fetches the final destination URL by following redirects and mimicking a regular browser's request.
func FetchTitle(url string) (string, error) {
	// Create a custom HTTP client for making the initial GET request
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Make a GET request to the URL
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Get the final URL after redirections
	finalURL := resp.Request.URL.String()

	// Check the content type of the response
	contentType := resp.Header.Get("Content-Type")

	// If the content type is HTML, use chromedp to navigate and get the final URL
	if contentType == "text/html" || contentType == "application/xhtml+xml" {
		// Create a custom HTTP client for chromedp
		opts := []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			chromedp.Headless,
			chromedp.DisableGPU,
			chromedp.IgnoreCertErrors, // This option makes chromedp ignore certificate errors
		}

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		// Create a new context from the allocator
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		var chromedpFinalURL string
		err = chromedp.Run(ctx,
			chromedp.Navigate(finalURL),
			chromedp.Sleep(1*time.Second), // Sleep to allow JavaScript execution
			chromedp.Location(&chromedpFinalURL),
		)

		if err != nil {
			return "", err
		}

		return chromedpFinalURL, nil
	}

	// If the content type is not HTML, return the final URL obtained from the GET request
	return finalURL, nil
}
