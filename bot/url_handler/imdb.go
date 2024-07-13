package url_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

// MovieInfo holds the data structure for the OMDb API response.
type MovieInfo struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	IMDbRating string `json:"imdbRating"`
	IMDbVotes  string `json:"imdbVotes"`
	IMDbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}

// ExtractIMDBID extracts the IMDB ID from an IMDB URL.
func ExtractIMDBID(url string) string {
	var imdbRegex = regexp.MustCompile(`https?://(?:www\.)?imdb\.com/title/(tt\d+)/?`)

	matches := imdbRegex.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func GetIMDBMovieInfo(movieID string) (string, error) {
	apiURL := fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&i=%s&r=json", os.Getenv("OMDb_API_KEY"), movieID)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("error making get request to OMDb API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OMDb API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading OMDb API response body: %v", err)
	}

	var movieInfo MovieInfo
	err = json.Unmarshal(body, &movieInfo)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling JSON response from OMDb API: %v", err)
	}

	if movieInfo.Response == "False" {
		return "Movie not found.", nil
	}

	// Format the movie information into a string.
	info := fmt.Sprintf("Title: %s (%s) | Director: %s | Actors: %s | Plot: %s | IMDb Rating: %s | URL: https://www.imdb.com/title/%s/",
		movieInfo.Title, movieInfo.Year, movieInfo.Director, movieInfo.Actors, movieInfo.Plot, movieInfo.IMDbRating, movieInfo.IMDbID)

	// Add color formatting for IRC
	info = "\x02\x0301,08IMDb ► \x03\x02 " + info

	return info, nil
}
