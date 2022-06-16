package intents

type Intent uint16

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
	MessageContent
)

var AllIntentsWithoutPrivileged = []Intent{
	Guilds, GuildBans, GuildEmojis, GuildIntegrations, GuildWebhooks, GuildInvites, GuildVoiceStates, GuildMessages,
	GuildMessageReactions, GuildMessageTyping, DirectMessages, DirectMessageReactions, DirectMessageTyping,
}

func SumIntents(intents ...Intent) (sum uint16) {
	for _, intent := range intents {
		sum += uint16(intent)
	}

	return
}
