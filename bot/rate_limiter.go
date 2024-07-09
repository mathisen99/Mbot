package bot

import (
	"fmt"
	"sync"
	"time"
)

const (
	maxCommandsPerSecond       = 2
	cooldownPeriod             = 10 * time.Second
	shutdownPeriod             = 1 * time.Hour
	globalMaxMessagesPerSecond = 5
)

type RateLimiter struct {
	mu                        sync.Mutex
	userTimestamps            map[string][]time.Time
	cooldowns                 map[string]time.Time
	shutdowns                 map[string]time.Time
	globalTimestamps          []time.Time
	lastSuspensionMessage     map[string]time.Time
	suspensionMessageCooldown time.Duration
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		userTimestamps:            make(map[string][]time.Time),
		cooldowns:                 make(map[string]time.Time),
		shutdowns:                 make(map[string]time.Time),
		globalTimestamps:          []time.Time{},
		lastSuspensionMessage:     make(map[string]time.Time),
		suspensionMessageCooldown: 1 * time.Minute,
	}
}

func (rl *RateLimiter) AllowCommand(user string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Check global rate limit
	globalValidTimestamps := []time.Time{}
	for _, timestamp := range rl.globalTimestamps {
		if now.Sub(timestamp) < time.Second {
			globalValidTimestamps = append(globalValidTimestamps, timestamp)
		}
	}
	rl.globalTimestamps = globalValidTimestamps

	if len(globalValidTimestamps) >= globalMaxMessagesPerSecond {
		return false
	}

	// Check if the user is in shutdown
	if shutdown, exists := rl.shutdowns[user]; exists {
		if now.Before(shutdown) {
			return false
		}
		delete(rl.shutdowns, user)
	}

	// Check if the user is in cooldown
	if cooldown, exists := rl.cooldowns[user]; exists {
		if now.Before(cooldown) {
			rl.shutdowns[user] = now.Add(shutdownPeriod)
			delete(rl.cooldowns, user)
			return false
		}
		delete(rl.cooldowns, user)
	}

	timestamps, exists := rl.userTimestamps[user]
	if !exists {
		rl.userTimestamps[user] = []time.Time{now}
		rl.globalTimestamps = append(rl.globalTimestamps, now)
		return true
	}

	// Filter out timestamps older than one second
	validTimestamps := []time.Time{}
	for _, timestamp := range timestamps {
		if now.Sub(timestamp) < time.Second {
			validTimestamps = append(validTimestamps, timestamp)
		}
	}
	rl.userTimestamps[user] = validTimestamps

	// Check per-second limit
	if len(validTimestamps) >= maxCommandsPerSecond {
		rl.cooldowns[user] = now.Add(cooldownPeriod)
		return false
	}

	rl.userTimestamps[user] = append(validTimestamps, now)
	rl.globalTimestamps = append(rl.globalTimestamps, now)
	return true
}

func (rl *RateLimiter) GetCooldownRemaining(user string) time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if cooldown, exists := rl.cooldowns[user]; exists {
		return cooldown.Sub(now)
	}
	return 0
}

func (rl *RateLimiter) GetShutdownRemaining(user string) time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if shutdown, exists := rl.shutdowns[user]; exists {
		return shutdown.Sub(now)
	}
	return 0
}

func (rl *RateLimiter) CanSendSuspensionMessage(user string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if lastMessage, exists := rl.lastSuspensionMessage[user]; exists {
		if now.Sub(lastMessage) < rl.suspensionMessageCooldown {
			return false
		}
	}
	rl.lastSuspensionMessage[user] = now
	return true
}

func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%d milliseconds", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%d minutes and %d seconds", minutes, seconds)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%d hours, %d minutes and %d seconds", hours, minutes, seconds)
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d days, %d hours, %d minutes and %d seconds", days, hours, minutes, seconds)
}
