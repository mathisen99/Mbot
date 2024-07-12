package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
	Hostmask string `json:"hostmask"`
	Role     string `json:"role"`
}

// Path to the users JSON file
var filePath = "./data/users.json"

// function to load users from a JSON file
func LoadUsers() (map[string]User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening users file: %w", err)
	}
	defer file.Close()

	var users []User
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("error decoding users file: %w", err)
	}

	userMap := make(map[string]User)
	for _, user := range users {
		userMap[user.Hostmask] = user
	}

	return userMap, nil
}

// Function to save users to a JSON file
func SaveUsers(users map[string]User) error {
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating users file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(userList); err != nil {
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

// Check if a user has a specific role
func GetUserRole(users map[string]User, hostmask string) string {
	normalizedHostmask := NormalizeHostmask(hostmask)
	for _, user := range users {
		if user.Hostmask == normalizedHostmask {
			return user.Role
		}
	}
	return "Everyone" // Default role if not found
}

// Function to get the role level of a user
func GetUserRoleLevel(users map[string]User, hostmask string) int {
	normalizedHostmask := NormalizeHostmask(hostmask)
	role := GetUserRole(users, normalizedHostmask)
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

	conn.AddCallback("311", func(e ircmsg.Message) {
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
				Role:     "Owner",
			}

			conn.Privmsg(ownerNick, "Hey! If you know me, spill the Setup password. If not, no worries—just laugh and pretend you never saw this. Bot out! 🤖")

			passwordConfirmed := make(chan bool)

			var privmsgCallbackID ircevent.CallbackID

			privmsgCallbackID = conn.AddCallback("PRIVMSG", func(e ircmsg.Message) {
				sourceParts := strings.SplitN(e.Source, "!", 2)
				nick := sourceParts[0]

				color.Yellow(">> Received message from %s: %s", nick, e.Params[1])

				if nick == ownerNick && len(e.Params) > 1 && e.Params[1] == setupPassword {
					color.Green(">> Setup password confirmed")
					conn.Privmsg(ownerNick, "Setup password confirmed. You are now the owner of the bot.")

					err := AddUser(users, owner)
					if err != nil {
						color.Red(">> Failed to add owner: %v", err)
						return
					}
					color.Green(">> Owner set successfully:")
					color.Green(">> Hostmask: %s", hostmask)

					passwordConfirmed <- true

					close(passwordConfirmed)

					conn.RemoveCallback(privmsgCallbackID)
				} else {
					color.Red(">> Setup password incorrect")
					conn.Privmsg(ownerNick, "Setup password was incorrect. Bye!")

					conn.RemoveCallback(privmsgCallbackID)

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
	return UserRoles[GetUserRole(users, hostmask)] == RoleOwner
}

// IsUserAdmin checks if a user is an admin
func IsUserAdmin(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= RoleAdmin
}

// IsUserTrusted checks if a user is trusted
func IsUserTrusted(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= RoleTrusted
}

// IsUserBadBoy checks if a user is a troll :) lol
func IsUserBadBoy(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] == RoleBadBoy
}
