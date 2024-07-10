package bot

import (
	"sync"
	"time"

	"github.com/fatih/color"
)

var lastMessageTime = make(map[string]time.Time)
var mu sync.Mutex

// Function to handle private messages
func handlePrivateMessage(connection *Connection, sender, message string) {
	color.Magenta(">> Private message from %s: %s", sender, message)
	nickname := ExtractNickname(sender)

	mu.Lock()
	defer mu.Unlock()

	// Check the last message time for the user
	if lastTime, ok := lastMessageTime[nickname]; ok {
		if time.Since(lastTime) < 1*time.Minute {
			// If the last message was sent within the last minute, ignore this message
			return
		}
	}

	lastMessageTime[nickname] = time.Now()

	connection.Privmsg(nickname, "Let's keep this between us. I won't tell anyone.")
}
