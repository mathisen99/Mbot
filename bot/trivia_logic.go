package bot

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ergochat/irc-go/ircevent"
)

type UserScores struct {
	Mu     sync.Mutex
	Scores map[string]int
}

var ScoresInstance = &UserScores{
	Scores: make(map[string]int),
}

func LoadScores(filename string) error {
	ScoresInstance.Mu.Lock()
	defer ScoresInstance.Mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &ScoresInstance.Scores)
}

func SaveScores(filename string) error {
	ScoresInstance.Mu.Lock()
	defer ScoresInstance.Mu.Unlock()

	bytes, err := json.Marshal(ScoresInstance.Scores)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, bytes, 0644)
}

func checkTriviaAnswer(sender, message, target string, connection *Connection) {
	TriviaStateInstance.Mu.Lock()
	defer TriviaStateInstance.Mu.Unlock()

	if TriviaStateInstance.AnsweredBy[sender] {
		return
	}

	if strings.EqualFold(message, TriviaStateInstance.Answer) {
		TriviaStateInstance.AnsweredBy[sender] = true
		nick := extractNickname(sender)
		connection.Privmsg(target, "Correct answer by "+nick+"!")
		TriviaStateInstance.Active = false
		if TriviaStateInstance.CancelFunc != nil {
			TriviaStateInstance.CancelFunc()
		}

		// Update the user's score
		ScoresInstance.Mu.Lock()
		ScoresInstance.Scores[sender]++
		ScoresInstance.Mu.Unlock()

		// Save the updated scores
		if err := SaveScores("./data/trivia_scores.json"); err != nil {
			connection.Privmsg(target, "Error saving scores: "+err.Error())
		}
	} else {
		// Do not respond to incorrect answers
		return
	}
}

func StartTriviaTimer(connection *ircevent.Connection, target string) {
	ctx, cancel := context.WithCancel(context.Background())
	TriviaStateInstance.Mu.Lock()
	TriviaStateInstance.CancelFunc = cancel
	TriviaStateInstance.Mu.Unlock()

	select {
	case <-ctx.Done():
		// Timer was cancelled because the correct answer was given
		return
	case <-time.After(30 * time.Second):
		TriviaStateInstance.Mu.Lock()
		defer TriviaStateInstance.Mu.Unlock()

		if TriviaStateInstance.Active {
			TriviaStateInstance.Active = false
			connection.Privmsg(target, "Time's up! No one got the correct answer.")
		}
	}
}

// Helper function to extract nickname from the full hostname
func extractNickname(fullHostname string) string {
	exclamationIndex := strings.Index(fullHostname, "!")
	if exclamationIndex != -1 {
		return fullHostname[:exclamationIndex]
	}
	return fullHostname // Return as is if the format is unexpected
}
