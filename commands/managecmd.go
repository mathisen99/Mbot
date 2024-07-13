package commands

import (
	"encoding/json"
	"fmt"
	"mbot/bot"
	"mbot/config"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ergochat/irc-go/ircevent"
)

// Handler for the !managecmd command
func ManageCommand(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User, cmdCfg *config.CommandConfig, configPath string) {
	args := strings.Fields(message)
	if len(args) < 2 {
		connection.Privmsg(target, "Usage: !managecmd <action> [parameters]")
		return
	}

	action := strings.ToLower(args[1])

	switch action {
	case "edit":
		handleEditCommand(connection, target, args, cmdCfg, configPath)
	case "add":
		handleAddCommand(connection, target, args, cmdCfg, configPath)
	case "remove":
		handleRemoveCommand(connection, target, args, cmdCfg, configPath)
	case "list":
		handleListCommands(connection, target, args, cmdCfg)
	case "setup":
		handleSetupCommand(connection, target, args, cmdCfg, configPath)
	default:
		connection.Privmsg(target, "Unsupported action. Supported actions are: edit, add, remove, list, setup")
	}
}

// NormalizeRole standardizes the role name to ensure consistency
func NormalizeRole(role string) string {
	switch strings.ToLower(role) {
	case "admin":
		return "Admin"
	case "everyone":
		return "Everyone"
	case "owner":
		return "Owner"
	case "trusted":
		return "Trusted"
	default:
		return role
	}
}

// Check if the role is valid
func isValidRole(role string) bool {
	validRoles := []string{"Admin", "Everyone", "Owner", "Trusted"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// Remove duplicate channels from a list
func removeDuplicateChannels(channels []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, channel := range channels {
		if _, exists := seen[channel]; !exists {
			seen[channel] = struct{}{}
			result = append(result, channel)
		}
	}
	return result
}

// Remove a command from all roles in specific channels
func removeCommandFromChannels(cmdCfg *config.CommandConfig, command string, channels []string) {
	for i := 0; i < len(cmdCfg.Commands[command]); i++ {
		perm := &cmdCfg.Commands[command][i]
		newChannels := []string{}
		for _, permChannel := range perm.Channels {
			keep := true
			for _, channel := range channels {
				if permChannel == channel {
					keep = false
					break
				}
			}
			if keep {
				newChannels = append(newChannels, permChannel)
			}
		}
		perm.Channels = newChannels
		if len(perm.Channels) == 0 {
			cmdCfg.Commands[command] = append(cmdCfg.Commands[command][:i], cmdCfg.Commands[command][i+1:]...)
			i-- // adjust index after removal
		}
	}
}

// LimitBackups keeps only the most recent 3 backup files
func limitBackups(configPath string) error {
	dir := filepath.Dir(configPath)
	pattern := fmt.Sprintf("%s.bak.*", filepath.Base(configPath))
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}
	if len(files) > 3 {
		sort.Slice(files, func(i, j int) bool {
			return files[i] > files[j]
		})
		for _, file := range files[3:] {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create a backup of the current configuration
func createBackup(configPath string) error {
	input, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading current config: %w", err)
	}
	backupPath := fmt.Sprintf("%s.bak.%d", configPath, time.Now().Unix())
	err = os.WriteFile(backupPath, input, 0644)
	if err != nil {
		return fmt.Errorf("error creating backup: %w", err)
	}
	return limitBackups(configPath)
}

// Edit an existing command's role and allowed channels
func handleEditCommand(connection *ircevent.Connection, target string, args []string, cmdCfg *config.CommandConfig, configPath string) {
	if len(args) < 5 {
		connection.Privmsg(target, "Usage: !managecmd edit <command> <role> <channels...>")
		return
	}

	command := args[2]
	role := NormalizeRole(args[3])
	channels := removeDuplicateChannels(args[4:])

	if !isValidRole(role) {
		connection.Privmsg(target, fmt.Sprintf("Role %s is invalid.", role))
		return
	}

	if _, exists := cmdCfg.Commands[command]; !exists {
		cmdCfg.Commands[command] = []config.CommandPermission{}
	}

	// Create a backup before making changes
	if err := createBackup(configPath); err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to create backup: %v", err))
		return
	}

	// Remove the command from all roles in the specified channels
	removeCommandFromChannels(cmdCfg, command, channels)

	// Add or update the command's permissions
	updated := false
	for i, perm := range cmdCfg.Commands[command] {
		if perm.Role == role {
			cmdCfg.Commands[command][i].Channels = append(cmdCfg.Commands[command][i].Channels, channels...)
			cmdCfg.Commands[command][i].Channels = removeDuplicateChannels(cmdCfg.Commands[command][i].Channels)
			updated = true
			break
		}
	}

	if !updated {
		cmdCfg.Commands[command] = append(cmdCfg.Commands[command], config.CommandPermission{
			Channels: channels,
			Role:     role,
		})
	}

	connection.Privmsg(target, fmt.Sprintf("Command %s updated to role %s for channels %v", command, role, channels))

	// Save the updated configuration
	err := saveCommandConfig(cmdCfg, configPath)
	if err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to save configuration: %v", err))
	}
}

// Add a new command to a specified role
func handleAddCommand(connection *ircevent.Connection, target string, args []string, cmdCfg *config.CommandConfig, configPath string) {
	if len(args) < 5 {
		connection.Privmsg(target, "Usage: !managecmd add <command> <role> <channels...>")
		return
	}

	command := args[2]
	role := NormalizeRole(args[3])
	channels := removeDuplicateChannels(args[4:])

	if !isValidRole(role) {
		connection.Privmsg(target, fmt.Sprintf("Role %s is invalid.", role))
		return
	}

	if _, exists := cmdCfg.Commands[command]; !exists {
		cmdCfg.Commands[command] = []config.CommandPermission{}
	}

	// Create a backup before making changes
	if err := createBackup(configPath); err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to create backup: %v", err))
		return
	}

	// Remove the command from all roles in the specified channels
	removeCommandFromChannels(cmdCfg, command, channels)

	cmdCfg.Commands[command] = append(cmdCfg.Commands[command], config.CommandPermission{
		Channels: channels,
		Role:     role,
	})

	connection.Privmsg(target, fmt.Sprintf("Command %s added to role %s for channels %v", command, role, channels))

	// Save the updated configuration
	err := saveCommandConfig(cmdCfg, configPath)
	if err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to save configuration: %v", err))
	}
}

// Remove a command from a specified role
func handleRemoveCommand(connection *ircevent.Connection, target string, args []string, cmdCfg *config.CommandConfig, configPath string) {
	if len(args) < 4 {
		connection.Privmsg(target, "Usage: !managecmd remove <command> <role>")
		return
	}

	command := args[2]
	role := NormalizeRole(args[3])

	if permissions, exists := cmdCfg.Commands[command]; exists {
		// Create a backup before making changes
		if err := createBackup(configPath); err != nil {
			connection.Privmsg(target, fmt.Sprintf("Failed to create backup: %v", err))
			return
		}

		for i, perm := range permissions {
			if perm.Role == role {
				cmdCfg.Commands[command] = append(cmdCfg.Commands[command][:i], cmdCfg.Commands[command][i+1:]...)
				connection.Privmsg(target, fmt.Sprintf("Command %s removed from role %s", command, role))

				// Save the updated configuration
				err := saveCommandConfig(cmdCfg, configPath)
				if err != nil {
					connection.Privmsg(target, fmt.Sprintf("Failed to save configuration: %v", err))
				}
				return
			}
		}
		connection.Privmsg(target, fmt.Sprintf("Command %s not found for role %s", command, role))
	} else {
		connection.Privmsg(target, fmt.Sprintf("Command %s not found", command))
	}
}

// List all permissions for a specified command
func handleListCommands(connection *ircevent.Connection, target string, args []string, cmdCfg *config.CommandConfig) {
	if len(args) < 3 {
		connection.Privmsg(target, "Usage: !managecmd list <command>")
		return
	}

	command := args[2]

	if permissions, exists := cmdCfg.Commands[command]; exists {
		for _, perm := range permissions {
			connection.Privmsg(target, fmt.Sprintf("Command: %s, Role: %s, Channels: %v", command, perm.Role, perm.Channels))
		}
	} else {
		connection.Privmsg(target, fmt.Sprintf("Command %s not found", command))
	}
}

// Setup default permissions for a new channel
func handleSetupCommand(connection *ircevent.Connection, target string, args []string, cmdCfg *config.CommandConfig, configPath string) {
	if len(args) < 3 {
		connection.Privmsg(target, "Usage: !managecmd setup <channel>")
		return
	}

	channel := args[2]

	// Get default permissions
	defaultPermissions := GetDefaultPermissions(channel)

	// Create a backup before making changes
	if err := createBackup(configPath); err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to create backup: %v", err))
		return
	}

	// Clear existing permissions for the channel
	for cmd, perms := range cmdCfg.Commands {
		newPerms := []config.CommandPermission{}
		for _, perm := range perms {
			newChannels := []string{}
			for _, ch := range perm.Channels {
				if ch != channel {
					newChannels = append(newChannels, ch)
				}
			}
			if len(newChannels) > 0 {
				newPerms = append(newPerms, config.CommandPermission{
					Role:     perm.Role,
					Channels: newChannels,
				})
			}
		}
		if len(newPerms) > 0 {
			cmdCfg.Commands[cmd] = newPerms
		} else {
			delete(cmdCfg.Commands, cmd)
		}
	}

	// Set default permissions
	for cmd, perms := range defaultPermissions {
		cmdCfg.Commands[cmd] = perms
	}

	connection.Privmsg(target, fmt.Sprintf("Default permissions set up for channel %s", channel))

	// Save the updated configuration
	err := saveCommandConfig(cmdCfg, configPath)
	if err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to save configuration: %v", err))
	}

	// Reload the command configuration
	err = bot.ReloadCommandConfig(configPath)
	if err != nil {
		connection.Privmsg(target, fmt.Sprintf("Failed to reload command configuration: %v", err))
	}

	// Re-register commands to reflect updated permissions
	RegisterAllCommands()
}

// Save the command configuration to a file
func saveCommandConfig(cmdCfg *config.CommandConfig, configPath string) error {
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("error creating command config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cmdCfg); err != nil {
		return fmt.Errorf("error encoding command config file: %w", err)
	}

	return nil
}

// RegisterManageCommand registers the !managecmd command
func RegisterManageCommand(cmdCfg *config.CommandConfig, configPath string) {
	bot.RegisterCommand("!managecmd", func(connection *ircevent.Connection, sender, target, message string, users map[string]bot.User) {
		ManageCommand(connection, sender, target, message, users, cmdCfg, configPath)
	})
}
