// bot/connection.go
package bot

import (
	"mbot/config"

	"github.com/ergochat/irc-go/ircevent"
)

type Connection struct {
	*ircevent.Connection
	Config *config.Config
}
