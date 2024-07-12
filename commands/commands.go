package commands

// RegisterAllCommands registers all commands in the package
func RegisterAllCommands() {
	RegisterHelloCommand()      // Hello example command
	RegisterHello2Command()     // Hello2 example command
	RegisterBaseCommands()      // Base commands (op, deop, kick, etc.)
	RegisterAddUserCommand()    // AddUser command (Used to add users to the bot)
	RegisterRemoveUserCommand() // RemoveUser command (Used to remove users from the bot)
}
