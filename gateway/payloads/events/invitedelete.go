package events

type InviteDelete struct {
	ChannelId uint64 `json:"channel_id,string"`
	GuildId   uint64 `json:"guild_id,string"`
	Code      string `json:"code"`
}
