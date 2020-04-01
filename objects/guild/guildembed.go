package guild

type GuildEmbed struct {
	Enabled   bool   `json:"enabled"`
	ChannelId uint64 `json:"channel_id,string"`
}
