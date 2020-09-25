package guild

type CachedRole struct {
	GuildId     uint64 `json:"-"` // Don't include in postgres
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions uint64 `json:"permissions,string"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`
}

func (r *CachedRole) ToRole(roleId uint64) Role {
	return Role{
		Id:          roleId,
		Name:        r.Name,
		Color:       r.Color,
		Hoist:       r.Hoist,
		Position:    r.Position,
		Permissions: r.Permissions,
		Managed:     r.Managed,
		Mentionable: r.Mentionable,
	}
}
