package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/fatih/color"
)

// UserRoles is a map of user roles to their respective values
var UserRoles = map[string]int{
	"Owner":    RoleOwner,
	"Admin":    RoleAdmin,
	"Trusted":  RoleTrusted,
	"Everyone": RoleEveryone,
	"BadBoy":   RoleBadBoy,
}

// Role levels
const (
	RoleEveryone = 0
	RoleBadBoy   = -10
	RoleTrusted  = 3
	RoleAdmin    = 5
	RoleOwner    = 10
)

// global users map
var Users map[string]User

// Structure to represent a user
type User struct {
	Hostmask string            `json:"hostmask"`
	Roles    map[string]string `json:"roles"` // map of channel to role
}

// Path to the users JSON file
var filePath = "./data/users.json"

// Mutex to protect access to the owner setup process
var ownerPromptMutex sync.Mutex
var ownerSetupActive bool

// LoadUsers loads the users from the specified file path and creates the file if it does not exist.
func LoadUsers(filePath string) (map[string]User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Ensure the directory exists
			err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("failed to create directory for users.json file: %v", err)
			}

			// Create an empty users.json file if it does not exist
			emptyUsers := make(map[string]User)
			file, err := os.Create(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create users.json file: %v", err)
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			err = encoder.Encode(emptyUsers)
			if err != nil {
				return nil, fmt.Errorf("failed to write to users.json file: %v", err)
			}

			return emptyUsers, nil
		}
		return nil, fmt.Errorf("error opening users file: %w", err)
	}
	defer file.Close()

	var users map[string]User
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("error decoding users file: %w", err)
	}

	return users, nil
}

// Function to save users to a JSON file
func SaveUsers(users map[string]User) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating users file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(users); err != nil {
		return fmt.Errorf("error encoding users file: %w", err)
	}

	return nil
}

// Function to add a user to the list of users
func AddUser(users map[string]User, user User) error {
	users[user.Hostmask] = user
	return SaveUsers(users)
}

// Function to remove a user from the list of users
func RemoveUser(users map[string]User, hostmask string) error {
	delete(users, hostmask)
	return SaveUsers(users)
}

// Function to get a list of users
func GetUserList(users map[string]User) string {
	userList := ""
	for hostmask := range users {
		userList += hostmask + " "
	}
	return userList
}

// Function to update a user in the list of users
func UpdateUser(users map[string]User, user User) error {
	users[user.Hostmask] = user
	return SaveUsers(users)
}

// Normalize the hostmask to ensure consistent format
func NormalizeHostmask(hostmask string) string {
	if !strings.HasPrefix(hostmask, "~") {
		return "~" + hostmask
	}
	return hostmask
}

// Check if a user has a specific role in a channel
func GetUserRole(users map[string]User, hostmask, channel string) string {
	normalizedHostmask := NormalizeHostmask(hostmask)
	for _, user := range users {
		if user.Hostmask == normalizedHostmask {
			if user.Roles["*"] == "Owner" {
				return "Owner"
			}
			if role, exists := user.Roles[channel]; exists {
				return role
			}
		}
	}
	return "Everyone" // Default role if not found
}

// Function to get the role level of a user in a channel
func GetUserRoleLevel(users map[string]User, hostmask, channel string) int {
	normalizedHostmask := NormalizeHostmask(hostmask)
	role := GetUserRole(users, normalizedHostmask, channel)
	return UserRoles[role]
}

// AddOwnerPrompt asks for the owner's nick and adds the owner to the users map
func AddOwnerPrompt(conn *Connection, users map[string]User) {
	color.Cyan("=============================== NO OWNER FOUND ===============================")
	color.Red("No owner was found in the users.json file. Please set an owner.")
	color.Red("The bot will shut down if no owner is set within 1 minute after connecting.")
	color.Red("The bot will message the owner to confirm the Setup password.")
	color.Cyan("==============================================================================")
	color.Blue(">> Please enter the nick of the owner on the network:")
	var ownerNick string
	fmt.Scanln(&ownerNick)
	color.Blue(">> Please enter your Setup password:")
	var setupPassword string
	fmt.Scanln(&setupPassword)

	conn.SendRaw(fmt.Sprintf("WHOIS %s", ownerNick))

	ownerPromptMutex.Lock()
	if ownerSetupActive {
		color.Red(">> Owner setup already active, returning.")
		ownerPromptMutex.Unlock()
		return
	}
	ownerSetupActive = true
	ownerPromptMutex.Unlock()

	defer func() {
		ownerPromptMutex.Lock()
		ownerSetupActive = false
		ownerPromptMutex.Unlock()
	}()

	passwordConfirmed := make(chan bool)
	var privmsgCallbackID ircevent.CallbackID
	var whoisCallbackID ircevent.CallbackID

	whoisCallbackID = conn.AddCallback("311", func(e ircmsg.Message) {
		if len(e.Params) >= 5 {
			user := e.Params[2]
			host := e.Params[3]

			var hostmask string
			if user[0] == '~' {
				hostmask = fmt.Sprintf("%s@%s", user, host)
			} else {
				hostmask = fmt.Sprintf("~%s@%s", user, host)
			}

			owner := User{
				Hostmask: hostmask,
				Roles:    map[string]string{"*": "Owner"},
			}

			conn.Privmsg(ownerNick, "Hey! If you know me, spill the Setup password. If not, no worries—just laugh and pretend you never saw this. Bot out! 🤖")

			privmsgCallbackID = conn.AddCallback("PRIVMSG", func(e ircmsg.Message) {
				sourceParts := strings.SplitN(e.Source, "!", 2)
				nick := sourceParts[0]

				color.Yellow(">> Received message from %s: %s", nick, e.Params[1])

				if strings.EqualFold(nick, ownerNick) && len(e.Params) > 1 && e.Params[1] == setupPassword {
					color.Green(">> Setup password confirmed")
					conn.Privmsg(ownerNick, "Setup password confirmed. You are now the owner of the bot.")
					conn.Privmsg(ownerNick, "Run the command !managecmd setup #channel in your channel where the bot is present to set up all the commands.")
					conn.Privmsg(ownerNick, "If you don't run the setup, no other commands will work except for the !managecmd command.")

					ownerPromptMutex.Lock()
					err := AddUser(users, owner)
					ownerPromptMutex.Unlock()

					if err != nil {
						color.Red(">> Failed to add owner: %v", err)
						return
					}
					color.Green(">> Owner set successfully:")
					color.Green(">> Hostmask: %s", hostmask)

					passwordConfirmed <- true
					close(passwordConfirmed)

					conn.RemoveCallback(privmsgCallbackID)
					conn.RemoveCallback(whoisCallbackID)
				} else {
					color.Red(">> Setup password incorrect")
					conn.Privmsg(ownerNick, "Setup password was incorrect. Bye!")

					conn.RemoveCallback(privmsgCallbackID)
					conn.RemoveCallback(whoisCallbackID)

					// Shutdown the bot
					log.Println("Sending shutdown signal")
					syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}
			})

			// Start a timer for 1 minute to shut down if no response is received
			go func() {
				select {
				case <-time.After(1 * time.Minute):
					color.Red(">> No response within 1 minute, shutting down.")
					conn.Privmsg(ownerNick, "No response within 1 minute. Shutting down. Bye!")
					conn.RemoveCallback(privmsgCallbackID)
					conn.RemoveCallback(whoisCallbackID)
					syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				case <-passwordConfirmed:
					return
				}
			}()
		}
	})
}

// Role comparison functions
func IsUserOwner(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask, "*")] == RoleOwner
}

// IsUserAdmin checks if a user is an admin in a channel
func IsUserAdmin(users map[string]User, hostmask, channel string) bool {
	return UserRoles[GetUserRole(users, hostmask, channel)] >= RoleAdmin
}

// IsUserTrusted checks if a user is trusted in a channel
func IsUserTrusted(users map[string]User, hostmask, channel string) bool {
	return UserRoles[GetUserRole(users, hostmask, channel)] >= RoleTrusted
}

// IsUserBadBoy checks if a user is a troll :) lol
func IsUserBadBoy(users map[string]User, hostmask, channel string) bool {
	return UserRoles[GetUserRole(users, hostmask, channel)] == RoleBadBoy
}
