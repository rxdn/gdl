package events

type MessageDelete struct {
	Id        uint64 `json:"id,string"`
	ChannelId uint64 `json:"channel_id,string"`
	GuildId   uint64 `json:"guild_id,string"`
}
