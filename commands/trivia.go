package commands

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"mbot/bot"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/fatih/color"
	openai "github.com/sashabaranov/go-openai"
)

// TriviaQuestion holds the structure of a trivia question with its topic, question, and answer
type TriviaQuestion struct {
	Topic    string `json:"topic"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

const triviaFile = "./data/trivia_questions.json"
const maxQuestions = 100

var triviaQuestions []TriviaQuestion
var triviaMu sync.Mutex

// LoadTriviaQuestions loads the trivia questions from a JSON file
func loadTriviaQuestions() error {
	file, err := os.ReadFile(triviaFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil // It's okay if the file doesn't exist
		}
		return err
	}
	return json.Unmarshal(file, &triviaQuestions)
}

// SaveTriviaQuestions saves the trivia questions to a JSON file
func saveTriviaQuestions() error {
	data, err := json.Marshal(triviaQuestions)
	if err != nil {
		return err
	}
	return os.WriteFile(triviaFile, data, 0644)
}

// HashQuestion creates a hash of a given question
func hashQuestion(question string) string {
	h := sha256.New()
	h.Write([]byte(question))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HashAnswer creates a hash of a given answer
func hashAnswer(answer string) string {
	h := sha256.New()
	h.Write([]byte(answer))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Clean and validate the question and answer
func cleanAndValidate(question, answer string) (string, string, error) {
	// Remove asterisks and other unwanted characters
	cleanQuestion := strings.TrimSpace(strings.ReplaceAll(question, "**", ""))
	cleanAnswer := strings.TrimSpace(strings.ReplaceAll(answer, "**", ""))

	// Validate the answer is either a single word or two-word proper nouns
	if len(strings.Fields(cleanAnswer)) > 2 {
		return "", "", errors.New("answer is longer than two words")
	}

	// Validate the answer is not empty
	if cleanAnswer == "" {
		return "", "", errors.New("answer is empty")
	}

	return cleanQuestion, cleanAnswer, nil
}

// GenerateTriviaQuestion generates a trivia question and answer based on the given topic using OpenAI
func GenerateTriviaQuestion(topic string) (string, string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	// Gather history of previous questions and answers for the topic
	triviaMu.Lock()
	var history strings.Builder
	for _, q := range triviaQuestions {
		if q.Topic == topic {
			history.WriteString(fmt.Sprintf("Question: %s Answer: %s\n", q.Question, q.Answer))
		}
	}
	triviaMu.Unlock()

	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: `You are a trivia question generator. Provide a trivia question and answer based on the given topic. Ensure the question is clear and the answer is concise. The answer should be a maximum of one word or a proper noun with two words (e.g., names, cities, countries).`,
	}

	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf("Generate a trivia question about %s. Here are previous questions and answers:\n%s", topic, history.String()),
	}

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4o,
		MaxTokens: 150,
		Messages:  []openai.ChatCompletionMessage{systemMessage, userMessage},
	}

	var question, answer string

	for attempts := 0; attempts < 5; attempts++ {
		resp, err := client.CreateChatCompletion(ctx, req)
		if err != nil {
			return "", "", err
		}

		answer = resp.Choices[0].Message.Content
		parts := strings.SplitN(answer, "Answer:", 2)
		if len(parts) < 2 {
			fmt.Println("Failed to split the response into question and answer")
			continue
		}

		question = strings.TrimSpace(parts[0])
		answer = strings.TrimSpace(parts[1])

		// Clean and validate the question and answer
		question, answer, err = cleanAndValidate(question, answer)
		if err != nil {
			fmt.Println("Validation failed:", err)
			continue
		}

		// Check if the question is non-empty and valid
		if question == "" {
			fmt.Println("Generated question is empty")
			continue
		}

		// Check for duplicates
		triviaMu.Lock()
		questionHash := hashQuestion(question)
		answerHash := hashAnswer(answer)
		isDuplicate := false
		for _, q := range triviaQuestions {
			if hashQuestion(q.Question) == questionHash || hashAnswer(q.Answer) == answerHash {
				isDuplicate = true
				break
			}
		}
		triviaMu.Unlock()

		if !isDuplicate {
			// Save the question and answer
			triviaMu.Lock()
			if len(triviaQuestions) >= maxQuestions {
				triviaQuestions = triviaQuestions[1:]
			}
			triviaQuestions = append(triviaQuestions, TriviaQuestion{Topic: topic, Question: question, Answer: answer})
			saveTriviaQuestions()
			triviaMu.Unlock()
			return question, answer, nil
		}
	}

	return "", "", errors.New("failed to generate a unique trivia question after several attempts")
}

// TriviaCommand handles the !trivia command, generating and starting a trivia game
func TriviaCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	bot.TriviaStateInstance.Mu.Lock()
	defer bot.TriviaStateInstance.Mu.Unlock()

	if bot.TriviaStateInstance.Active {
		connection.Privmsg(target, "Trivia is already active!")
		return
	}

	parts := strings.Fields(message)
	if len(parts) < 2 {
		connection.Privmsg(target, "Please provide a topic for the trivia question. Usage: !trivia <topic>")
		return
	}

	topic := strings.Join(parts[1:], " ")

	question, answer, err := GenerateTriviaQuestion(topic)
	if err != nil {
		connection.Privmsg(target, "Error generating trivia question: "+err.Error())
		return
	}

	if question == "" {
		connection.Privmsg(target, "Failed to generate a valid trivia question. Please try again.")
		return
	}

	bot.TriviaStateInstance.Active = true
	bot.TriviaStateInstance.Question = question
	bot.TriviaStateInstance.Answer = answer
	bot.TriviaStateInstance.AnsweredBy = make(map[string]bool)

	connection.Privmsg(target, question)

	// Start the trivia timer
	go StartTriviaTimer(connection, target, answer)
}

// StartTriviaTimer starts a 30-second timer for the trivia game
func StartTriviaTimer(connection *ircevent.Connection, target string, answer string) {
	ctx, cancel := context.WithCancel(context.Background())
	bot.TriviaStateInstance.Mu.Lock()
	bot.TriviaStateInstance.CancelFunc = cancel
	bot.TriviaStateInstance.Mu.Unlock()

	select {
	case <-ctx.Done():
		// Timer was cancelled because the correct answer was given
		return
	case <-time.After(30 * time.Second):
		bot.TriviaStateInstance.Mu.Lock()
		defer bot.TriviaStateInstance.Mu.Unlock()

		if bot.TriviaStateInstance.Active {
			bot.TriviaStateInstance.Active = false
			connection.Privmsg(target, "Time's up! No one got the correct answer. The answer was: "+answer)
		}
	}
}

// Trivia command to display user scores
func ScoresCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	fmt.Println("Scores command triggered")

	bot.ScoresInstance.Mu.Lock()
	defer bot.ScoresInstance.Mu.Unlock()

	if len(bot.ScoresInstance.Scores) == 0 {
		connection.Privmsg(target, "No scores yet!")
		return
	}

	var scoresList []struct {
		Nickname string
		Score    int
	}
	for user, score := range bot.ScoresInstance.Scores {
		nickname := extractNickname(user)
		scoresList = append(scoresList, struct {
			Nickname string
			Score    int
		}{Nickname: nickname, Score: score})
	}

	// Sort scores in descending order
	sort.Slice(scoresList, func(i, j int) bool {
		return scoresList[i].Score > scoresList[j].Score
	})

	// Find the sender's score
	senderNickname := extractNickname(sender)
	var senderScore int
	foundSender := false
	for _, entry := range scoresList {
		if entry.Nickname == senderNickname {
			senderScore = entry.Score
			foundSender = true
			break
		}
	}
	if !foundSender {
		senderScore = 0 // or any default value if the sender doesn't have a score yet
	}

	// Prepare the top 5 scores
	var top5Scores strings.Builder
	topCount := 5
	for i, entry := range scoresList {
		if i >= topCount {
			break
		}
		if i > 0 {
			top5Scores.WriteString(", ")
		}
		top5Scores.WriteString(fmt.Sprintf("%s: %d", entry.Nickname, entry.Score))
	}

	// Construct the message
	message = fmt.Sprintf("Your score is %d and the Top5 is: %s", senderScore, top5Scores.String())

	connection.Privmsg(target, message)
}

// Helper function to extract nickname from the full hostname
func extractNickname(fullHostname string) string {
	exclamationIndex := strings.Index(fullHostname, "!")
	if exclamationIndex != -1 {
		return fullHostname[:exclamationIndex]
	}
	return fullHostname // Return as is if the format is unexpected
}

// RegisterTriviaCommand registers the trivia command
func RegisterTriviaCommand() {
	bot.RegisterCommand("!trivia", TriviaCommand)
	bot.RegisterCommand("!trivia-top", ScoresCommand)
}

func init() {
	err := loadTriviaQuestions()
	if err != nil {
		color.Red("Failed to load trivia questions: %s", err)
	}

	if err := bot.LoadScores("./data/trivia_scores.json"); err != nil {
		log.Fatalf("Error loading scores: %v", err)
	}
}
