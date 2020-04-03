package events

type MessageReactionRemoveAll struct {
	ChannelId uint64 `json:"channel_id,string"`
	MessageId uint64 `json:"message_id,string"`
	GuildId   uint64 `json:"guild_id,string"`
}
