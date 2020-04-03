package cache

type CacheOptions struct {
	Guilds      bool
	Users       bool
	Members     bool // requires Guilds = true
	Channels    bool // requires Guilds = true
	Roles       bool // requires Guilds = true
	Emojis      bool // requires Guilds = true
	VoiceStates bool
}
