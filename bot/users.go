package bot

import (
	"encoding/json"
	"fmt"
	"os"
)

// global users map
var Users map[string]User

// Structure to represent a user
type User struct {
	Hostmask string `json:"hostmask"`
	Role     string `json:"role"`
}

// UserRoles is a map of user roles to their respective values
var UserRoles = map[string]int{
	"Owner":   5,
	"Admin":   4,
	"Trusted": 3,
	"Regular": 2,
	"BadBoy":  1,
}

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

// Check if a user has a specific role
func GetUserRole(users map[string]User, hostmask string) string {
	if user, exists := users[hostmask]; exists {
		return user.Role
	}
	return "Regular" // Default role if not found
}

// Role comparison functions
func IsUserOwner(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] == UserRoles["Owner"]
}

func IsUserAdmin(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= UserRoles["Admin"]
}

func IsUserTrusted(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= UserRoles["Trusted"]
}

func IsUserRegular(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] >= UserRoles["Regular"]
}

func IsUserBadBoy(users map[string]User, hostmask string) bool {
	return UserRoles[GetUserRole(users, hostmask)] == UserRoles["BadBoy"]
}
