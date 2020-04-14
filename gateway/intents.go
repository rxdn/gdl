package gateway

type Intent int

const (
	Guilds Intent = 1 << iota
	GuildMembers
	GuildBans
	GuildEmojis
	GuildIntegrations
	GuildWebhooks
	GuildInvites
	GuildVoiceStates
	GuildPresences
	GuildMessages
	GuildMessageReactions
	GuildMessageTyping
	DirectMessages
	DirectMessageReactions
	DirectMessageTyping
)
