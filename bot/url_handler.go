package bot

import (
	"os"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"

	cmd "mbot/bot/internal"
)

// HandleUrl processes URLs found in messages
func HandleUrl(connection *ircevent.Connection, target, url string) {
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
		connection.Privmsg(target, "You posted a link: "+url)
	}
}

// HandleYoutubeLink processes YouTube links
func HandleYoutubeLink(connection *ircevent.Connection, target, url string) {
	videoID := cmd.ExtractVideoID(url)
	yourAPIKey := os.Getenv("YOUTUBE_API_KEY")
	videoInfo, err := cmd.GetYouTubeVideoInfo(videoID, yourAPIKey)
	if err != nil {
		color.Red(">> Error getting video info: %v", err)
	} else {
		connection.Privmsg(target, videoInfo)
	}
}

// HandleWikipediaLink processes Wikipedia links
func HandleWikipediaLink(connection *ircevent.Connection, target, url string) {
	connection.Privmsg(target, "You shared a Wikipedia link! Wikipedia is a free online encyclopedia, created and edited by volunteers around the world.")
}

// HandleGithubLink processes GitHub links
func HandleGithubLink(connection *ircevent.Connection, target, url string) {
	connection.Privmsg(target, "You shared a GitHub link! GitHub is a web-based platform used for version control and collaboration.")
}

// HandleIMDbLink processes IMDb links
func HandleIMDbLink(connection *ircevent.Connection, target, url string) {
	connection.Privmsg(target, "You shared an IMDb link! IMDb is an online database of information related to films, television programs, home videos, video games, and streaming content.")
}
