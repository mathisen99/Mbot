package bot

import (
	"context"
	"sync"
	"time"
)

type TriviaState struct {
	Active     bool
	Question   string
	Answer     string
	AnsweredBy map[string]bool
	Mu         sync.Mutex
	CancelFunc context.CancelFunc
	StartTime  time.Time
}

var TriviaStateInstance = &TriviaState{
	AnsweredBy: make(map[string]bool),
}
