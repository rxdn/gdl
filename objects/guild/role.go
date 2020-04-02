package guild

import (
	"fmt"
)

type Role struct {
	Id          uint64 `json:"id,string"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`
}

func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%d>", r.Id)
}
