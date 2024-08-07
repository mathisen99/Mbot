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
		fmt.Println("Sent usage message") // Debug print
		return
	}

	kbNumber := args[1]
	fmt.Println("Fetching KB update information for:", kbNumber) // Debug print

	// Command execution
	connection.Privmsg(target, "Fetching KB update information...")
	cmd := exec.Command("python3", "./kb/main.py", kbNumber)
	fmt.Println("Running Python script...") // Debug print
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error fetching KB update information:", err) // Debug print
		connection.Privmsg(target, "Error fetching KB update information.")
		fmt.Println("Sent error message") // Debug print
		return
	}

	fmt.Println("Python script completed")         // Debug print
	fmt.Println("Command output:", string(output)) // Debug print

	// Process the output to extract description and size
	lines := strings.Split(string(output), "\n")
	var description, size string
	for _, line := range lines {
		if strings.HasPrefix(line, "Description: ") {
			description = strings.TrimPrefix(line, "Description: ")
		} else if strings.HasPrefix(line, "Size: ") {
			size = strings.TrimPrefix(line, "Size: ")
		}
	}

	if description == "" || size == "" {
		connection.Privmsg(target, "No description or size found or failed to retrieve data.")
		fmt.Println("Sent no data message") // Debug print
		return
	}

	// Split the description into chunks of 450 characters
	chunkSize := 450
	for i := 0; i < len(description); i += chunkSize {
		end := i + chunkSize
		if end > len(description) {
			end = len(description)
		}
		chunk := description[i:end]
		connection.Privmsg(target, chunk)
		fmt.Println("Sending description chunk:", chunk) // Debug print
	}

	// Send the size
	connection.Privmsg(target, "Size: "+size)
	fmt.Println("Sending size:", size) // Debug print
}

// RegisterKBCommand registers the !kb command
func RegisterKBCommand() {
	bot.RegisterCommand("!kb", KBCommand)
}
