package url_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Repository represents a GitHub repository.
type Repository struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
}

// fetchGithubRepoInfo fetches information about a GitHub repository from the GitHub API.
func FetchGithubRepoInfo(repoURL string) (string, error) {
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var repository Repository
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &repository)
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("Github Repository info: %s, Stars: %d, Forks: %d, Open Issues: %d, Description: %s",
		repository.FullName, repository.StargazersCount, repository.ForksCount, repository.OpenIssuesCount, repository.Description)

	return info, nil
}

// parseGitHubURL extracts the owner and repo name from a GitHub URL.
func parseGitHubURL(gitHubURL string) (owner, repo string, err error) {
	parsedURL, err := url.Parse(gitHubURL)
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", errors.New("invalid GitHub URL")
	}
	return parts[0], parts[1], nil
}
