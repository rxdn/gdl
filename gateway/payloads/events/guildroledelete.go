package events

type GuildRoleDelete struct {
	GuildId uint64 `json:"guild_id,string"`
	RoleId  uint64 `json:"role_id,string"`
}
