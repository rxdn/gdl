package intents

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

func SumIntents(intents ...Intent) int {
	var sum int

	for _, intent := range intents {
		sum += int(intent)
	}

	return sum
}
