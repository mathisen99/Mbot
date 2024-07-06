package bot

import "strings"

// ExtractNickname extracts the nickname from the full sender string
func ExtractNickname(fullSender string) string {
	if idx := strings.Index(fullSender, "!"); idx != -1 {
		return fullSender[:idx]
	}
	return fullSender
}
