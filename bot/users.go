package bot

import (
	"encoding/json"
	"fmt"
	"os"
)

// UserRoles is a map of user roles to their respective values
var UserRoles = map[string]int{
	"Owner":   RoleOwner,
	"Admin":   RoleAdmin,
	"Trusted": RoleTrusted,
	"Regular": RoleRegular,
	"BadBoy":  RoleBadBoy,
}

// Role levels
const (
	RoleEveryone = 0
	RoleBadBoy   = -1
	RoleRegular  = 2
	RoleTrusted  = 3
	RoleAdmin    = 4
	RoleOwner    = 5
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

// Function to get the role level of a user
func GetUserRoleLevel(users map[string]User, hostmask string) int {
	role := GetUserRole(users, hostmask)
	return UserRoles[role]
}

// Check if a user has a specific role
func GetUserRole(users map[string]User, hostmask string) string {
	if user, exists := users[hostmask]; exists {
		return user.Role
	}
	return "Everyone" // Default role if not found
}

// Role comparison functions
func IsUserOwner(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] == RoleOwner
}

func IsUserAdmin(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= RoleAdmin
}

func IsUserTrusted(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= RoleTrusted
}

func IsUserRegular(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= RoleRegular
}

func IsUserBadBoy(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] == RoleBadBoy
}
