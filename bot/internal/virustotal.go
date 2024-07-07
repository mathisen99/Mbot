package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type VirusTotalResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CheckAndFetchURLReport checks the URL with VirusTotal and fetches the report.
func CheckAndFetchURLReport(urlToCheck string) (string, error) {
	id, err := checkURLWithVirusTotal(urlToCheck)
	if err != nil {
		return "", fmt.Errorf("error checking URL: %v", err)
	}
	color.Green(">> Performing scan on this url: %v", urlToCheck)
	color.Magenta(">> URL check submitted. Analysis ID: %v", id)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	tryCount := 0 // To prevent infinite loops
	for range ticker.C {
		report, err := getVirusTotalReport(id, os.Getenv("VIRUSTOTAL_API_KEY"))
		if err != nil {
			return "", fmt.Errorf("error getting report: %v", err)
		}

		if len(report) != 0 {
			return formatReport(report), nil
		}

		tryCount++
		if tryCount > 5 { // Limit the number of retries
			return "", fmt.Errorf("report taking too long, please try again later")
		}
	}

	return "", fmt.Errorf("no report generated")
}

func formatReport(report map[string]int) string {
	formattedStats := make([]string, 0)

	for stat, value := range report {
		if value > 1 { // Only show stats with values greater than 1 to avoid false positives
			formattedStat := formatStat(stat, value)
			formattedStats = append(formattedStats, formattedStat)
		}
	}

	if len(formattedStats) == 0 {
		return "No significant findings"
	}

	return strings.Join(formattedStats, ", ")
}

func formatStat(stat string, value int) string {
	var color string

	switch stat {
	case "harmless":
		color = "\x033" // Green
	case "malicious":
		color = "\x034" // Red
	case "suspicious":
		color = "\x038" // Yellow
	case "undetected":
		color = "\x036" // Cyan
	case "timeout":
		color = "\x030" // White
	default:
		color = "\x0f" // Default (reset color)
	}

	return fmt.Sprintf("%s%s: %d", color, stat, value)
}

func checkURLWithVirusTotal(urlToCheck string) (string, error) {
	data := url.Values{}
	data.Set("url", urlToCheck)
	apiKey := os.Getenv("VIRUSTOTAL_API_KEY")

	req, err := http.NewRequest("POST", "https://www.virustotal.com/api/v3/urls", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response VirusTotalResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Data.ID, nil
}

func getVirusTotalReport(id string, apiKey string) (map[string]int, error) {
	req, err := http.NewRequest("GET", "https://www.virustotal.com/api/v3/analyses/"+id, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Data struct {
			Attributes struct {
				Stats map[string]int `json:"stats"`
			} `json:"attributes"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Data.Attributes.Stats, nil
}
