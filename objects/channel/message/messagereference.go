package message

type MessageReference struct {
	MessageId uint64 `json:"message_id,string"`
	ChannelId uint64 `json:"channel_id,string"`
	GuildId   uint64 `json:"guild_id,string"`
}
