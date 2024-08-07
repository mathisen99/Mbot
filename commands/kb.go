package commands

import (
	"mbot/bot"
	"os/exec"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !kb command
func KBCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
	args := strings.Split(message, " ")
	if len(args) != 2 {
		connection.Privmsg(target, "Usage: !kb <KB_NUMBER>")
		return
	}

	kbNumber := args[1]
	output, err := exec.Command("python3", "./kb/main.py", kbNumber).Output()
	if err != nil {
		connection.Privmsg(target, "Error fetching KB update information.")
		return
	}

	connection.Privmsg(target, string(output))
}

// RegisterKBCommand registers the !kb command
func RegisterKBCommand() {
	bot.RegisterCommand("!kb", KBCommand)
}
