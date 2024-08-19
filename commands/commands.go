package commands

import "mbot/config"

// RegisterAllCommands registers all commands in the package
func RegisterAllCommands() {
	RegisterHelloCommand()        // Hello example command
	RegisterHello2Command()       // Hello2 example command
	RegisterURLCommand()          // URL command to enable/disable URL features (YouTube, Wikipedia, etc.)
	RegisterBaseCommands()        // Base commands (op, deop, kick, etc.)
	RegisterAddUserCommand()      // AddUser command (Used to add users to the bot)
	RegisterRemoveUserCommand()   // RemoveUser command (Used to remove users from the bot)
	RegisterPersonalityCommands() // Personality commands (Used to set the bot's personality for a channel)
	RegisterMemoryWipeCommand()   // MemoryWipe command (Used to wipe a user's memory)
	RegisterYTCommand()           // YouTube search command
	RegisterTriviaCommand()       // Trivia command
	RegisterKBCommand()           // KB search command
	RegisterClaudeCommand()       // Claude command (Anthropic API)
}

// GetDefaultPermissions returns the default command permissions for a given channel
func GetDefaultPermissions(channel string) map[string][]config.CommandPermission {
	return map[string][]config.CommandPermission{
		// User management commands
		"!adduser": {{Role: "Admin", Channels: []string{channel}}},
		"!deluser": {{Role: "Admin", Channels: []string{channel}}},

		// Claude command
		"!claude": {{Role: "Everyone", Channels: []string{channel}}},

		// Trivia command
		"!trivia":     {{Role: "Everyone", Channels: []string{channel}}},
		"!trivia-top": {{Role: "Everyone", Channels: []string{channel}}},

		// kb search command
		"!kb": {{Role: "Everyone", Channels: []string{channel}}},

		// Personality and memory commands
		"!personality": {{Role: "Admin", Channels: []string{channel}}},
		"!memory":      {{Role: "Everyone", Channels: []string{channel}}},

		// URL command
		"!url": {{Role: "Admin", Channels: []string{channel}}},
		"!yt":  {{Role: "Everyone", Channels: []string{channel}}},

		// Base commands
		"!op":      {{Role: "Admin", Channels: []string{channel}}},
		"!deop":    {{Role: "Admin", Channels: []string{channel}}},
		"!voice":   {{Role: "Admin", Channels: []string{channel}}},
		"!devoice": {{Role: "Admin", Channels: []string{channel}}},
		"!kick":    {{Role: "Admin", Channels: []string{channel}}},
		"!ban":     {{Role: "Admin", Channels: []string{channel}}},
		"!unban":   {{Role: "Admin", Channels: []string{channel}}},
		"!invite":  {{Role: "Admin", Channels: []string{channel}}},
		"!topic":   {{Role: "Admin", Channels: []string{channel}}},
		"!join":    {{Role: "Admin", Channels: []string{channel}}},
		"!part":    {{Role: "Admin", Channels: []string{channel}}},

		// Example command
		"!hello": {{Role: "Everyone", Channels: []string{channel}}},

		// Owner commands
		"!shutdown":  {{Role: "Owner", Channels: []string{channel}}},
		"!nick":      {{Role: "Owner", Channels: []string{channel}}},
		"!managecmd": {{Role: "Owner", Channels: []string{channel}}},

		// Trusted commands
		"!hello2": {{Role: "Trusted", Channels: []string{channel}}}, // Example command for testing purposes
	}
}
