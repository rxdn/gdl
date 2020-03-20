package events

type GuildIntegrationsUpdate struct {
	GuildId uint64 `json:"guild_id,string"`
}
