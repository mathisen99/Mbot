package commands

import "mbot/config"

// RegisterAllCommands registers all commands in the package
func RegisterAllCommands() {
	RegisterHelloCommand()      // Hello example command
	RegisterHello2Command()     // Hello2 example command
	RegisterBaseCommands()      // Base commands (op, deop, kick, etc.)
	RegisterAddUserCommand()    // AddUser command (Used to add users to the bot)
	RegisterRemoveUserCommand() // RemoveUser command (Used to remove users from the bot)
}

// GetDefaultPermissions returns the default command permissions for a given channel
func GetDefaultPermissions(channel string) map[string][]config.CommandPermission {
	return map[string][]config.CommandPermission{
		// User management commands
		"!adduser": {{Role: "Admin", Channels: []string{channel}}},
		"!deluser": {{Role: "Admin", Channels: []string{channel}}},

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
		"!hello":   {{Role: "Everyone", Channels: []string{channel}}},

		// Owner commands
		"!shutdown":  {{Role: "Owner", Channels: []string{channel}}},
		"!nick":      {{Role: "Owner", Channels: []string{channel}}},
		"!managecmd": {{Role: "Owner", Channels: []string{channel}}},

		// Trusted commands
		"!hello2": {{Role: "Trusted", Channels: []string{channel}}},
	}
}
