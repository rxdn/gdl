package guild

type CachedRole struct {
	Name        string `db:"name"`
	Color       int    `db:"color"`
	Hoist       bool   `db:"hoist"`
	Position    int    `db:"position"`
	Permissions int    `db:"permissions"`
	Managed     bool   `db:"managed"`
	Mentionable bool   `db:"mentionable"`
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
