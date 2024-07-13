package bot

import (
	"fmt"
	"mbot/bot/url_features"
	"mbot/config"
	"os"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
)

// Global variable to hold the URL configuration
var URLConfigData *config.URLFeatures

// HandleUrl processes URLs found in messages
func HandleUrl(connection *Connection, sender, target, url string) {
	featureConfig := URLConfigData

	switch {
	case strings.Contains(url, "youtube.com"), strings.Contains(url, "youtu.be"):
		if featureConfig.EnableYouTubeCheck {
			HandleYoutubeLink(connection, target, url)
		} else {
			color.Red(">> YouTube link handling is disabled")
		}
	case strings.Contains(url, "wikipedia.org"):
		if featureConfig.EnableWikipediaCheck {
			HandleWikipediaLink(connection, target, url)
		} else {
			color.Red(">> Wikipedia link handling is disabled")
		}
	case strings.Contains(url, "github.com"):
		if featureConfig.EnableGithubCheck {
			HandleGithubLink(connection, target, url)
		} else {
			color.Red(">> GitHub link handling is disabled")
		}
	case strings.Contains(url, "imdb.com"):
		if featureConfig.EnableIMDbCheck {
			HandleIMDbLink(connection, target, url)
		} else {
			color.Red(">> IMDb link handling is disabled")
		}
	default:
		GetTitle(connection.Connection, target, url)
		if featureConfig.EnableVirusTotalCheck {
			HandleVirusTotalLink(connection, sender, target, url)
		} else {
			color.Red(">> VirusTotal link handling is disabled")
		}
	}
}

// Function to get url title
func GetTitle(connection *ircevent.Connection, target, url string) {
	title, err := url_features.FetchTitle(url)
	if err != nil || title == "" {
		color.Red(">> Error fetching title if <nil>: %v the page does not have a title", err)
	} else {
		connection.Privmsg(target, "^ "+title)
	}
}

// HandleYoutubeLink processes YouTube links
func HandleYoutubeLink(connection *Connection, target, url string) {
	videoID := url_features.ExtractVideoID(url)
	yourAPIKey := os.Getenv("YOUTUBE_API_KEY")
	if yourAPIKey == "" {
		color.Red(">> YouTube API key is not set")
		connection.Privmsg(target, "YouTube API key is not set. Please set it in the environment variable YOUTUBE_API_KEY. or disable the feature in the configuration file.")
		return
	}

	videoInfo, err := url_features.GetYouTubeVideoInfo(videoID, yourAPIKey)
	if err != nil {
		color.Red(">> Error getting video info: %v", err)
		connection.Privmsg(target, "Error getting video info.")
	} else {
		connection.Privmsg(target, videoInfo)
	}
}

// HandleWikipediaLink processes Wikipedia links
func HandleWikipediaLink(connection *Connection, target, url string) {
	connection.Privmsg(target, "Wikipedia links are not supported yet.")
}

// HandleGithubLink processes GitHub links
func HandleGithubLink(connection *Connection, target, url string) {
	info, err := url_features.FetchGithubRepoInfo(url)
	if err != nil {
		color.Red(">> Error fetching GitHub repository info: %v", err)
		connection.Privmsg(target, "Error fetching GitHub repository info.")
	} else {
		connection.Privmsg(target, info)
	}
}

// HandleIMDbLink processes IMDb links
func HandleIMDbLink(connection *Connection, target, url string) {
	// check that API key is set
	if os.Getenv("OMDB_API_KEY") == "" {
		color.Red(">> OMDB API key is not set")
		connection.Privmsg(target, "OMDB API key is not set. Please set it in the environment variable OMDb_API_KEY. or disable the feature in the configuration file.")
		return
	}

	movieID := url_features.ExtractIMDBID(url)
	if movieID == "" {
		color.Red(">> Error extracting IMDb ID from URL")
		connection.Privmsg(target, "Error extracting IMDb ID from URL.")
		return
	}

	info, err := url_features.GetIMDBMovieInfo(movieID)
	if err != nil {
		color.Red(">> Error fetching IMDb movie info: %v", err)
		connection.Privmsg(target, "Error fetching IMDb movie info.")
	} else {
		connection.Privmsg(target, info)
	}
}

// HandleVirusTotalLink processes links using VirusTotal
func HandleVirusTotalLink(connection *Connection, sender, target, url string) {
	if os.Getenv("VIRUSTOTAL_API_KEY") == "" {
		color.Red(">> VirusTotal API key is not set")
		connection.Privmsg(target, "VirusTotal API key is not set. Please set it in the environment variable VIRUSTOTAL_API_KEY. or disable the feature in the configuration file.")
		return
	}

	nick := ExtractNickname(sender)
	reportMessage, err := url_features.CheckAndFetchURLReport(url)
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
