package bot

import (
	"mbot/config"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/fatih/color"
)

type Bot struct {
	Connection *Connection
	Config     *config.Config
}

// NewBot creates a new bot instance
func NewBot(cfg *config.Config, users map[string]User) *Bot {
	ircCon := &ircevent.Connection{
		Server:       cfg.Server + ":" + cfg.Port,
		Nick:         cfg.Nick,
		UseTLS:       cfg.UseTLS,
		TLSConfig:    cfg.TLSConfig,
		SASLLogin:    cfg.NickServUser,
		SASLPassword: cfg.NickServPass,
		RequestCaps:  []string{"server-time", "message-tags", "account-tag"},
	}

	conn := &Connection{
		Connection: ircCon,
		Config:     cfg,
	}

	bot := &Bot{
		Connection: conn,
		Config:     cfg,
	}

	bot.Connection.AddConnectCallback(func(e ircmsg.Message) {
		color.Green(">> Connection successful")

		// Check if an owner is set in the users map
		ownerFound := false
		for _, user := range users {
			if user.Role == "Owner" {
				ownerFound = true
				break
			}
		}

		if !ownerFound {
			color.Red(">> Owner not found or invalid in users.json")
			AddOwnerPrompt(bot.Connection, users)
		}

		for _, channel := range cfg.Channels {
			bot.Connection.Join(channel)
		}
	})

	// Registering callbacks and events
	RegisterCallbacks(bot.Connection.Connection)
	RegisterEventHandlers(bot.Connection, users)

	return bot
}

// Connect connects the bot to the server
func (b *Bot) Connect() error {
	return b.Connection.Connect()
}

// Loop starts the bot's main loop
func (b *Bot) Loop() {
	b.Connection.Loop()
}
