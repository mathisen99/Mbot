package commands

import (
	"fmt"
	"mbot/bot"
	"os/exec"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !kb command
func KBCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	fmt.Println("Received command:", message) // Debug print

	args := strings.Split(message, " ")
	if len(args) != 2 {
		connection.Privmsg(target, "Usage: !kb <KB_NUMBER>")
		return
	}

	kbNumber := args[1]
	fmt.Println("Fetching KB update information for:", kbNumber) // Debug print

	// Command execution
	cmd := exec.Command("python", "./kb/main.py", kbNumber)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error fetching KB update information:", err) // Debug print
		connection.Privmsg(target, "Error fetching KB update information.")
		return
	}

	fmt.Println("Command output:", string(output)) // Debug print

	// Split the output into lines and send each line separately to avoid message length issues
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			connection.Privmsg(target, line)
		}
	}
}

// RegisterKBCommand registers the !kb command
func RegisterKBCommand() {
	bot.RegisterCommand("!kb", KBCommand)
}
