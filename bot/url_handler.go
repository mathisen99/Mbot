package bot

import (
	"fmt"
	"mbot/bot/internal"
	"os"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
)

// HandleUrl processes URLs found in messages
func HandleUrl(connection *ircevent.Connection, sender, target, url string) {
	switch {
	case strings.Contains(url, "youtube.com"), strings.Contains(url, "youtu.be"):
		HandleYoutubeLink(connection, target, url)
	case strings.Contains(url, "wikipedia.org"):
		HandleWikipediaLink(connection, target, url)
	case strings.Contains(url, "github.com"):
		HandleGithubLink(connection, target, url)
	case strings.Contains(url, "imdb.com"):
		HandleIMDbLink(connection, target, url)
	default:
		GetTitle(connection, target, url)
		HandleVirusTotalLink(connection, sender, target, url)
	}
}

// Function to get url title
func GetTitle(connection *ircevent.Connection, target, url string) {
	title, err := internal.FetchTitle(url)
	if err != nil {
		color.Red(">> Error fetching title: %v", err)
		connection.Privmsg(target, "Error fetching title.")
	} else {
		connection.Privmsg(target, "^ "+title)
	}
}

// HandleYoutubeLink processes YouTube links
func HandleYoutubeLink(connection *ircevent.Connection, target, url string) {
	videoID := internal.ExtractVideoID(url)
	yourAPIKey := os.Getenv("YOUTUBE_API_KEY")
	videoInfo, err := internal.GetYouTubeVideoInfo(videoID, yourAPIKey)
	if err != nil {
		color.Red(">> Error getting video info: %v", err)
		connection.Privmsg(target, "Error getting video info.")
	} else {
		connection.Privmsg(target, videoInfo)
	}
}

// HandleWikipediaLink processes Wikipedia links
func HandleWikipediaLink(connection *ircevent.Connection, target, url string) {
	connection.Privmsg(target, "Wikipedia links are not supported yet.")
}

// HandleGithubLink processes GitHub links
func HandleGithubLink(connection *ircevent.Connection, target, url string) {
	info, err := internal.FetchGithubRepoInfo(url)
	if err != nil {
		color.Red(">> Error fetching GitHub repository info: %v", err)
		connection.Privmsg(target, "Error fetching GitHub repository info.")
	} else {
		connection.Privmsg(target, info)
	}
}

// HandleIMDbLink processes IMDb links
func HandleIMDbLink(connection *ircevent.Connection, target, url string) {
	movieID := internal.ExtractIMDBID(url)
	if movieID == "" {
		color.Red(">> Error extracting IMDb ID from URL")
		connection.Privmsg(target, "Error extracting IMDb ID from URL.")
		return
	}

	info, err := internal.GetIMDBMovieInfo(movieID)
	if err != nil {
		color.Red(">> Error fetching IMDb movie info: %v", err)
		connection.Privmsg(target, "Error fetching IMDb movie info.")
	} else {
		connection.Privmsg(target, info)
	}
}

// HandleVirusTotalLink processes links using VirusTotal
func HandleVirusTotalLink(connection *ircevent.Connection, sender, target, url string) {
	nick := ExtractNickname(sender)
	reportMessage, err := internal.CheckAndFetchURLReport(url)
	if err != nil {
		color.Red(">> Error checking URL with VirusTotal: %v", err)
		connection.Privmsg(target, "Error checking URL with VirusTotal.")
	} else {
		if strings.Contains(reportMessage, "malicious") {
			color.Red(">> URL is malicious: %s", url)
			connection.Privmsg(target, fmt.Sprintf("⚠️ %s just pasted a link that triggered my automatic defense systems ☢️ %s Here is a VirusTotal report: %s Note: low malicious score may be false positive", nick, url, reportMessage))
		} else {
			color.Green(">> URL is safe: %s", url)
		}
	}
}
