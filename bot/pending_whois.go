package bot

import "sync"

// PendingWhois stores pending WHOIS requests
var (
	PendingWhois = make(map[string]func(string))
	WhoisMu      sync.Mutex
)
