package events

type VoiceServerUpdate struct {
	Token   string `json:"token"`
	GuildId uint64 `json:"guild_id,string"`
	Endpoint string `json:"endpoint"`
}
